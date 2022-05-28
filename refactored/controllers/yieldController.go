package controllers

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"apsim-api/refactored/utils"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
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

	ctx, _ := context.WithTimeout(r.Context(), 20*time.Second)

	//Retrieving locationId and cultureID from middlewares
	locationId, _ := r.Context().Value("locationId").(uint64)
	cultureId, _ := r.Context().Value("cultureId").(uint64)

	//Parsing URL params
	urlParams := r.URL.Query()
	//parse direct into string YYYY-MM-DD
	fromDate, err := time.Parse("20060102", urlParams.Get("from"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: fromDate is not in YYYYMMDD format")
		return
	}
	toDate, err := time.Parse("20060102", urlParams.Get("to"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: toDate is not in YYYYMMDD format")
		return
	}
	fmt.Println("locationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

	location, err := c.LocationService.GetLocationById(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching location: locationId doesn't exist")
		return
	}
	fmt.Println("Location:", location)

	latestMicroclimateReading, err := c.MicroclimateReadingService.GetLatestMicroClimateReading(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: latest microclimate reading")
		return
	}
	fmt.Println("latest microclimate reading", latestMicroclimateReading)

	lastTime, err := time.Parse("2006-01-02", latestMicroclimateReading.Date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: latest microclimate reading date")
		return
	}
	fmt.Println("lastTime", lastTime)

	soil, err := c.SoilService.GetSoilByLocationId(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching soil: soil doesn't exist for given locationID %d", locationId)
		return
	}
	fmt.Println("Soil:", soil)

	fmt.Println("Checking cache")

	yields, err := c.YieldService.GetYields(ctx, int(locationId), int(cultureId), fromDate, toDate)
	if err == nil && yields != nil {
		fmt.Println("Found cached results in InfluxDB")
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(&yields)
		w.Write(response)
		return
	} else {
		fmt.Println("Did not found cached results in InfluxDB")
	}

	fmt.Println("Starting work")
	//ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	//Prepare goroutines, each for one file(apsimx, csv and consts)
	g, ctx := errgroup.WithContext(ctx)
	//Make a buffered channel for shared communication because apsimx goroutine needs the filepaths of other two files
	//ID = 0 for CSV and ID = 1 for Consts file

	ch := make(chan models.Message, 2)
	mainCh := make(chan models.Message, 3)

	g.Go(func() error {
		return c.LocationService.GenerateConstsFile(location.Name, location.Latitude, location.Longitude, ch, mainCh, ctx)
	})

	//Create a apsimx file, fetch the
	g.Go(func() error {
		return c.CultureService.GenerateAPSIMXFile(int(cultureId), fromDate, toDate, soil, ch, mainCh, ctx)
	})

	//Create a CSV file
	g.Go(func() error {
		return c.MicroclimateReadingService.GenerateCSVFile(int(locationId), fromDate, toDate, lastTime, ch, mainCh, ctx)
	})

	fmt.Println("Waiting for workers to finish")
	err = g.Wait()
	fmt.Println("errgroup error", err) //prints done
	fmt.Println("Workers done, closing channel")
	close(ch)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Continue work")

	csvAbsPath, constsAbsPath, apsimxAbsPath := utils.GetStageFilesAbsPaths(mainCh)

	fmt.Println("absolute paths of generated files:", csvAbsPath, constsAbsPath, apsimxAbsPath)

	fmt.Println("Starting simulation")
	err = utils.RunAPSIMSimulation(apsimxAbsPath)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in running simulation, pick a different from & to")
		return
	}

	dbAbsPath := utils.ConstructDBAbsPath(apsimxAbsPath)
	fmt.Println("dbAbsPath:", dbAbsPath)

	yields, err = utils.ReadAPSIMSimulationResults(dbAbsPath)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing db yield results")
		return
	}

	if len(yields) != (toDate.Year() - fromDate.Year() + 1) {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching db yield results")
		return
	}

	fmt.Println("yields:", yields)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(yields)
	w.Write(response)

	err = c.YieldService.CreateYields(ctx, int(locationId), int(cultureId), fromDate, toDate, yields)
	if err != nil {
		fmt.Println("Error writing in influx db", err)
	}

}
