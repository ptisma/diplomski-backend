package getMicroclimateReading

import (
	"apsim-api/internal/services/microclimateReadingService"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetMicroclimateReading(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		microclimateReadingService := microclimateReadingService.GetMicroclimateReadingService(app)
		latestMicroclimateReading, err := microclimateReadingService.GetLatestMicroclimateReading(int(locationId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: latest microclimate reading")
			return
		}
		fmt.Println("latest microclimate reading", latestMicroclimateReading)

		readings, err := microclimateReadingService.GetMicroclimateReadings(int(microclimateId), int(locationId), fromDate, toDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: microclimate readings")
			return
		}
		fmt.Println("readings:", readings)

		lastDate, err := time.Parse("2006-01-02", latestMicroclimateReading.Date)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in parsing: latest microclimate reading date")
			return
		}
		if toDate.After(lastDate) {
			preadings, err := microclimateReadingService.GetPredictedMicroclimateReadings(int(microclimateId), int(locationId), lastDate.AddDate(0, 0, 1), toDate)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error in fetching: predicted microclimate readings")
				return
			}
			cpreadings := microclimateReadingService.ConvertPredictedMicroclimateReadings(preadings)
			fmt.Println("cpreadings:", cpreadings)
			readings = append(readings, cpreadings...)

		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(&readings)
		w.Write(response)
		//params := mux.Vars(r)
		//fmt.Println("Hello world")
		//fmt.Println(params["id1"], params["id2"])
		//locationId, err := strconv.ParseInt(params["locationId"], 10, 32)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//microclimateId, _ := strconv.ParseUint(params["microclimateId"], 10, 32)

		//
		//fmt.Println("microclimateId:", microclimateId)
		//urlParams := r.URL.Query()
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

		//fromDate, _ := time.Parse("20060102", urlParams.Get("from"))
		//toDate, _ := time.Parse("20060102", urlParams.Get("to"))
		//locationId, _ := strconv.ParseUint(urlParams.Get("locationId"), 10, 32)
		//fmt.Println("locationId:", locationId, "microclimateId:", microclimateId, "fromDate:", fromDate, "toDate:", toDate)

		//microclimateReading := models.MicroclimateReading{
		//	MicroclimateID: uint32(microclimateId),
		//	LocationID:     uint32(locationId),
		//	FromDate:       fromDate,
		//	ToDate:         toDate,
		//}
		//
		//microclimateReadings, _ := microclimateReading.GetMicroclimateReading(app)
		//
		//w.Header().Set("Content-Type", "application/json")
		//response, _ := json.Marshal(microclimateReadings)
		//w.Write(response)

	})
}
