package controllers

import (
	"apsim-api/internal/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GrowingDegreeDayController struct {
	CultureService             interfaces.ICultureService
	MicroclimateReadingService interfaces.IMicroclimateReadingService
	MicroclimateService        interfaces.IMicroclimateService
}

func (c *GrowingDegreeDayController) GetGrowingDegreeDays(w http.ResponseWriter, r *http.Request) {

	//Retrieving locationId and cultureID from middlewares
	locationId, _ := r.Context().Value("locationId").(uint64)
	cultureId, _ := r.Context().Value("cultureId").(uint64)
	fromDate, _ := r.Context().Value("from").(time.Time)
	toDate, _ := r.Context().Value("to").(time.Time)

	//currentDate := time.Now().Format("20060102")
	//log.Println("Current date:", currentDate)

	////Parsing URL params
	//urlParams := r.URL.Query()
	////parse direct into string YYYY-MM-DD
	//fromDate, err := time.Parse("20060102", urlParams.Get("from"))
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintf(w, "Error in parsing: fromDate is not in YYYYMMDD format")
	//	return
	//}
	//toDate, err := time.Parse("20060102", urlParams.Get("to"))
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintf(w, "Error in parsing: toDate is not in YYYYMMDD format")
	//	return
	//}
	//log.Println("locationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

	culture, err := c.CultureService.FetchCultureById(ctx, int(cultureId))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: culture")
		return
	}
	//log.Println("Fetched culture", culture)

	tmax, err := c.MicroclimateService.GetMicroclimateByName(ctx, "tmax")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: tmax microclimate")
		return
	}

	tmin, err := c.MicroclimateService.GetMicroclimateByName(ctx, "tmin")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: tmin microclimate readings")
		return
	}

	//log.Println("tmax:", tmax)
	//log.Println("tmin:", tmin)

	x, err := c.MicroclimateReadingService.GetMicroclimateReadings(ctx, int(tmax.ID), int(locationId), fromDate, toDate)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: tmax microclimate readings")
		return
	}
	//log.Println("x:", x)

	y, err := c.MicroclimateReadingService.GetMicroclimateReadings(ctx, int(tmin.ID), int(locationId), fromDate, toDate)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: tmin microclimate readings")
		return
	}
	//log.Println("y:", y)

	gdds, err := c.MicroclimateReadingService.CalculateGrowingDegreeDay(x, y, culture.BaseTemperature)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in calculating: history growing degree days")
		return
	}
	//log.Println("gdds:", gdds)

	latestMicroclimateReading, err := c.MicroclimateReadingService.GetLatestMicroClimateReading(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: latest microclimate reading")
		return
	}
	//log.Println("latest microclimate reading", latestMicroclimateReading)
	lastDate, err := time.Parse("2006-01-02", latestMicroclimateReading.Date)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: latest microclimate reading date")
		return
	}

	if toDate.After(lastDate) {
		var newFromDate time.Time
		if fromDate.After(lastDate) {
			newFromDate = fromDate
		} else {
			newFromDate = lastDate.AddDate(0, 0, 1)
		}
		x, err := c.MicroclimateReadingService.GetPredictedMicroclimateReadings(ctx, int(tmax.ID), int(locationId), newFromDate, toDate)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmax predicted microclimate readings")
			return
		}

		//log.Println("predicted x:", x)

		y, err := c.MicroclimateReadingService.GetPredictedMicroclimateReadings(ctx, int(tmin.ID), int(locationId), newFromDate, toDate)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmin predicted microclimate readings")
			return
		}

		//log.Println("predicted y:", y)

		predictedGdds, err := c.MicroclimateReadingService.CalculatePredictedGrowingDegreeDay(x, y, culture.BaseTemperature)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in calculating: predicted growing degree days")
			return
		}
		//log.Println("predictedGdds:", predictedGdds)
		gdds = append(gdds, predictedGdds...)

	}

	log.Println("gdds:", gdds)
	if c.MicroclimateReadingService.ValidateGrowingDegreeDays(fromDate, toDate, gdds) == false {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching growing degree days, wrong dates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&gdds)

	w.Write(response)

}
