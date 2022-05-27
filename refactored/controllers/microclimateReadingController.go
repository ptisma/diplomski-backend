package controllers

import (
	"apsim-api/refactored/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MicroclimateReadingController struct {
	I interfaces.IMicroclimateReadingService
}

func (c *MicroclimateReadingController) GetMicroclimateReadings(w http.ResponseWriter, r *http.Request) {

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

	microclimateReadings, err := c.I.GetMicroclimateReadings(int(microclimateId), int(locationId), fromDate, toDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching: microclimate readings")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimateReadings)
	w.Write(response)
}
