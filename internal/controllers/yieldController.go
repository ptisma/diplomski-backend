package controllers

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

type YieldController struct {
	LocationService            interfaces.ILocationService
	CultureService             interfaces.ICultureService
	MicroclimateReadingService interfaces.IMicroclimateReadingService
	SoilService                interfaces.ISoilService
	YieldService               interfaces.IYieldService
}

func (c *YieldController) GetYield(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 40*time.Second)

	//Retrieving locationId and cultureID from middlewares
	locationId, _ := r.Context().Value("locationId").(uint64)
	cultureId, _ := r.Context().Value("cultureId").(uint64)
	fromDate, _ := r.Context().Value("from").(time.Time)
	toDate, _ := r.Context().Value("to").(time.Time)

	//log.Println("locationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

	location, err := c.LocationService.GetLocationById(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching location: locationId doesn't exist")
		return
	}
	//log.Println("Location:", location)

	latestMicroclimateReading, err := c.MicroclimateReadingService.GetLatestMicroClimateReading(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: latest microclimate reading")
		return
	}
	//log.Println("latest microclimate reading", latestMicroclimateReading)

	lastTime, err := time.Parse("2006-01-02", latestMicroclimateReading.Date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: latest microclimate reading date")
		return
	}
	//log.Println("lastTime", lastTime)

	soil, err := c.SoilService.GetSoilByLocationId(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching soil: soil doesn't exist for given locationID %d", locationId)
		return
	}
	//log.Println("Soil:", soil)

	log.Println("Checking cache")

	yields, err := c.YieldService.GetYields(ctx, int(locationId), int(cultureId), fromDate, toDate)
	//remove nil in yield service
	if err == nil && yields != nil {
		if c.YieldService.ValidateYields(fromDate, toDate, yields) == true {
			log.Println("Found cached results in InfluxDB")
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(&yields)
			w.Write(response)
			return
		}
	} else {
		log.Println("Error reading from influx db", err)
	}
	log.Println("Did not found cached results in InfluxDB")
	log.Println("Starting work")
	//ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	//Prepare goroutines, each for one file(apsimx, csv and consts)
	g, ctxx := errgroup.WithContext(ctx)
	//Make a buffered channel for shared communication because apsimx goroutine needs the filepaths of other two files
	//ID = 0 for CSV and ID = 1 for Consts file
	//chans closed in goroutines, move here?
	ch := make(chan utils.Message, 2)
	mainCh := make(chan utils.Message, 3)

	var csvAbsPath, constsAbsPath, apsimxAbsPath, dbAbsPath string
	var debug bool //if error in running simulation dont delete .db file
	defer func() {
		//Cleanup function which cleans only after this point
		//clean up work deleting, not handled
		log.Println("Deleting stage files")
		utils.DeleteStageFile(csvAbsPath)
		utils.DeleteStageFile(constsAbsPath)
		utils.DeleteStageFile(apsimxAbsPath)
		if debug == false {
			utils.DeleteStageFile(dbAbsPath) //commented for now due to debugging reasons
		}

	}()

	g.Go(func() error {
		return c.LocationService.GenerateConstsFile(location.Name, location.Latitude, location.Longitude, ch, mainCh, ctxx)
	})

	//Create a apsimx file, fetch the
	g.Go(func() error {
		return c.CultureService.GenerateAPSIMXFile(int(cultureId), fromDate, toDate, soil, ch, mainCh, ctxx)
	})

	//Create a CSV file
	g.Go(func() error {
		return c.MicroclimateReadingService.GenerateCSVFile(int(locationId), fromDate, toDate, lastTime, ch, mainCh, ctxx)
	})
	log.Println("Waiting for workers to finish")
	err = g.Wait()
	//log.Println("Workers done, closing channel")
	//close(ch)

	csvAbsPath, constsAbsPath, apsimxAbsPath = utils.GetStageFilesAbsPaths(mainCh)
	log.Println("absolute paths of generated files:", csvAbsPath, constsAbsPath, apsimxAbsPath)
	if err != nil {
		log.Println("errGroup error:", err)
		log.Println("Finishing")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in generating files for simulation")
		return
	}
	//log.Println("Continue work")

	log.Println("Starting simulation")
	err = utils.RunAPSIMSimulation(ctx, apsimxAbsPath)

	//Disregarding the outcome of simulation the .db file will be created
	dbAbsPath = utils.ConstructDBAbsPath(apsimxAbsPath)
	log.Println("dbAbsPath:", dbAbsPath)

	if err != nil {
		//log.Println(err)
		errLogs, err := c.YieldService.ReadYieldErrors(dbAbsPath)
		if err != nil {
			log.Println(err)
		}
		for index, message := range errLogs {
			log.Printf("Error %d: %s\n", index, message)
		}
		debug = true

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in running simulation, pick a different from & to")
		return
	}

	yields, err = c.YieldService.ReadYields(dbAbsPath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing db yield results")
		return
	}

	if c.YieldService.ValidateYields(fromDate, toDate, yields) == false {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching db yield results")
		return
	}

	log.Println("yields:", yields)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&yields)
	w.Write(response)

	err = c.YieldService.CreateYields(ctx, int(locationId), int(cultureId), fromDate, toDate, yields)
	if err != nil {
		log.Println("Error writing in influx db", err)
	}
	log.Println("Ended")
}
