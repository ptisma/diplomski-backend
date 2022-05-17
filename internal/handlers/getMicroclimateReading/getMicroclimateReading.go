package getMicroclimateReading

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func GetMicroclimateReading(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		//fmt.Println("Hello world")
		//fmt.Println(params["id1"], params["id2"])
		//locationId, err := strconv.ParseInt(params["locationId"], 10, 32)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		microclimateId, _ := strconv.ParseUint(params["microclimateId"], 10, 32)

		//
		//fmt.Println("microclimateId:", microclimateId)
		urlParams := r.URL.Query()
		//fromDate := urlParams.Get("from")
		//toDate := urlParams.Get("to")
		//fmt.Println(fromDate, toDate)

		//second on the January
		//x, _ := time.Parse(fromDate, "01-02-2006")
		//fmt.Println(x.Day(), x.Month(), x.Year())
		//2 January 2006

		//fmt.Println(urlParams.Get("from"), urlParams.Get("to"))
		//fromDate, err := time.Parse(urlParams.Get("from"), "01-02-2006")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//toDate, err := time.Parse(urlParams.Get("to"), "02-01-2006")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println(fromDate.Day(), fromDate.Month(), fromDate.Year())
		//fmt.Println(toDate.Day(), toDate.Month(), toDate.Year())
		//fmt.Println(locationId, microclimateId, fromDate, toDate)
		//
		//microclimateReading := &models.MicroclimateReading{
		//	MicroclimateID: uint32(microclimateId),
		//	LocationID:     uint32(locationId),
		//	FromDate:       fromDate,
		//	ToDate:         toDate,
		//}
		//
		////export or else it will be invisible
		//
		//w.Header().Set("Content-Type", "application/json")
		//response, _ := json.Marshal(microclimateReading)
		//w.Write(response)

		fromDate, _ := time.Parse("20060102", urlParams.Get("from"))
		toDate, _ := time.Parse("20060102", urlParams.Get("to"))
		locationId, _ := strconv.ParseUint(urlParams.Get("locationId"), 10, 32)
		fmt.Println("locationId:", locationId, "microclimateId:", microclimateId, "fromDate:", fromDate, "toDate:", toDate)

		microclimateReading := models.MicroclimateReading{
			MicroclimateID: uint32(microclimateId),
			LocationID:     uint32(locationId),
			FromDate:       fromDate,
			ToDate:         toDate,
		}

		microclimateReadings, _ := microclimateReading.GetMicroclimateReading(app)

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(microclimateReadings)
		w.Write(response)

	})
}
