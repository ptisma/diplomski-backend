package controllers

import "C"
import (
	"apsim-api/internal/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	//log.Println("period:", period)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: microclimate readings period")
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
	fromDate, _ := r.Context().Value("from").(time.Time)
	toDate, _ := r.Context().Value("to").(time.Time)

	//log.Println("locationId:", locationId, "microclimateId:", microclimateId, "fromDate:", fromDate, "toDate:", toDate)

	microclimateReadings, err := c.I.GetMicroclimateReadings(ctx, int(microclimateId), int(locationId), fromDate, toDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: microclimate readings")
		return
	}
	//log.Println("microclimateReadings:", microclimateReadings)

	latestMicroclimateReading, err := c.I.GetLatestMicroClimateReading(ctx, int(locationId))
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
	//log.Println("lastDate:", lastDate)

	if toDate.After(lastDate) {
		log.Println("Have to fetch predicted")
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
		cpreadings := c.I.ConvertPredictedMicroclimateReadings(preadings)
		//log.Println("cpreadings:", cpreadings)
		microclimateReadings = append(microclimateReadings, cpreadings...)

	}

	log.Println("microclimateReadings:", microclimateReadings)

	if c.I.ValidateMicroclimateReadings(fromDate, toDate, microclimateReadings) == false {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching microclimate readings, wrong dates")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimateReadings)
	w.Write(response)
}
