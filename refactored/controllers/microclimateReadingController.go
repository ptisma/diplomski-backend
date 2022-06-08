package controllers

import "C"
import (
	"apsim-api/refactored/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MicroclimateReadingController struct {
	I interfaces.IMicroclimateReadingService
}

func (c *MicroclimateReadingController) GetMicroclimateReadingPeriod(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

	locationId, _ := r.Context().Value("locationId").(uint64)
	microclimateId, _ := r.Context().Value("microclimateId").(uint64)

	period, err := c.I.GetMicroclimateReadingPeriod(ctx, int(microclimateId), int(locationId))
	fmt.Println("period:", period)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: microclimate readings")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&period)
	w.Write(response)

}

func (c *MicroclimateReadingController) GetMicroclimateReadings(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

	locationId, _ := r.Context().Value("locationId").(uint64)
	microclimateId, _ := r.Context().Value("microclimateId").(uint64)

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
	fmt.Println("locationId:", locationId, "microclimateId:", microclimateId, "fromDate:", fromDate, "toDate:", toDate)

	microclimateReadings, err := c.I.GetMicroclimateReadings(ctx, int(microclimateId), int(locationId), fromDate, toDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: microclimate readings")
		return
	}
	fmt.Println("microclimateReadings:", microclimateReadings)

	latestMicroclimateReading, err := c.I.GetLatestMicroClimateReading(ctx, int(locationId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: latest microclimate reading")
		return
	}
	fmt.Println("latest microclimate reading", latestMicroclimateReading)

	lastDate, err := time.Parse("2006-01-02", latestMicroclimateReading.Date)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in parsing: latest microclimate reading date")
		return
	}
	fmt.Println("lastDate:", lastDate)

	if toDate.After(lastDate) {

		fmt.Println("Have to fetch predicted")
		var newFromDate time.Time
		if fromDate.After(lastDate) {
			newFromDate = fromDate
		} else {
			newFromDate = lastDate.AddDate(0, 0, 1)
		}

		preadings, err := c.I.GetPredictedMicroclimateReadings(ctx, int(microclimateId), int(locationId), newFromDate, toDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: predicted microclimate readings")
			return
		}
		cpreadings := c.I.ConvertPredictedMicroclimateReadings(ctx, preadings)
		fmt.Println("cpreadings:", cpreadings)
		microclimateReadings = append(microclimateReadings, cpreadings...)

	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimateReadings)
	w.Write(response)
}
