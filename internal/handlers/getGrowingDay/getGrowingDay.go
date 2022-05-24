package getGrowingDay

import (
	"apsim-api/internal/models"
	cultureService2 "apsim-api/internal/services/cultureService"
	"apsim-api/internal/services/microclimateReadingService"
	"apsim-api/internal/services/microclimateService"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetGrowingDay(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Retrieving locationId and cultureID from middlewares
		locationId, _ := r.Context().Value("locationId").(uint64)
		cultureId, _ := r.Context().Value("cultureId").(uint64)

		//currentDate := time.Now().Format("20060102")
		//fmt.Println("Current date:", currentDate)

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

		//culture := &models.Culture{ID: uint32(cultureId)}
		//culture.GetCultureById(app)
		cultureService := cultureService2.GetCultureService(app)
		culture, err := cultureService.GetCulture(int(cultureId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: culture")
			return
		}
		fmt.Println("Fetched culture", culture)

		microclimateService := microclimateService.GetMicroclimateService(app)
		tmax, err := microclimateService.GetMicroclimate("tmax")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmax microclimate")
			return
		}

		tmin, err := microclimateService.GetMicroclimate("tmin")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmin microclimate readings")
			return
		}
		//tmax := models.Microclimate{Name: "tmax"}
		//tmax.GetMicroclimateByName(app)
		fmt.Println("tmax:", tmax)
		fmt.Println("tmin:", tmin)

		microclimateReadingService := microclimateReadingService.GetMicroclimateReadingService(app)
		//microclimateReading := models.MicroclimateReading{
		//	MicroclimateID: tmax.ID,
		//	LocationID:     uint32(locationId),
		//	FromDate:       fromDate,
		//	ToDate:         toDate,
		//}
		//
		////fmt.Println("microclimateReading", microclimateReading)
		////fix for all x >= from not included
		x, err := microclimateReadingService.GetMicroclimateReadings(int(tmax.ID), int(locationId), fromDate, toDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmax microclimate readings")
			return
		}
		fmt.Println("x:", x)

		//
		//tmin := models.Microclimate{Name: "tmin"}
		//tmin.GetMicroclimateByName(app)

		//
		//microclimateReading = models.MicroclimateReading{
		//	MicroclimateID: tmin.ID,
		//	LocationID:     uint32(locationId),
		//	FromDate:       fromDate,
		//	ToDate:         toDate,
		//}

		y, err := microclimateReadingService.GetMicroclimateReadings(int(tmin.ID), int(locationId), fromDate, toDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching: tmin microclimate readings")
			return
		}
		fmt.Println("y:", y)
		//
		var gdds = []models.GrowingDegreeDay{}

		//if len(*x) == len(*y) {
		//	for i, _ := range *x {
		//
		//		date := (*x)[i].Date
		//		gdd := ((*x)[i].Value+(*y)[i].Value)/2 - float32(culture.BaseTemperature)
		//
		//		gdds = append(gdds, models.GrowingDegreeDay{
		//			Date:  date,
		//			Value: gdd,
		//		})
		//
		//	}
		//}
		// il
		gdds, err = microclimateReadingService.CalculateGrowingDegreeDay(x, y, culture.BaseTemperature)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in calculating: history growing degree days")
			return
		}
		fmt.Println("gdds:", gdds)
		//err = microclimateReading.GetLatestMicroclimateReading(app)
		latestMicroclimateReading, err := microclimateReadingService.GetLatestMicroclimateReading(int(locationId))
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

		if toDate.After(lastDate) {
			x, err := microclimateReadingService.GetPredictedMicroclimateReadings(int(tmax.ID), int(locationId), lastDate.AddDate(0, 0, 1), toDate)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error in fetching: tmax predicted microclimate readings")
				return
			}

			y, err := microclimateReadingService.GetPredictedMicroclimateReadings(int(tmin.ID), int(locationId), lastDate.AddDate(0, 0, 1), toDate)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error in fetching: tmin predicted microclimate readings")
				return
			}

			predictedGdds, err := microclimateReadingService.CalculatePredictedGrowingDegreeDay(x, y, culture.BaseTemperature)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error in calculating: history growing degree days")
				return
			}
			fmt.Println("predictedGdds:", predictedGdds)
			gdds = append(gdds, predictedGdds...)

			//microclimateReading := models.PredictedMicroclimateReading{
			//	MicroclimateID: tmax.ID,
			//	LocationID:     uint32(locationId),
			//	FromDate:       fromDate,
			//	ToDate:         toDate,
			//}
			//x, err := microclimateReading.GetPredictedMicroclimateReading(app)
			//if err != nil {
			//	w.WriteHeader(http.StatusInternalServerError)
			//	fmt.Fprintf(w, "Error in fetching: tmax predicted microclimate readings")
			//	return
			//}
			//fmt.Println("predicted x:", x)
			//
			//microclimateReading = models.PredictedMicroclimateReading{
			//	MicroclimateID: tmin.ID,
			//	LocationID:     uint32(locationId),
			//	FromDate:       fromDate,
			//	ToDate:         toDate,
			//}
			//y, err := microclimateReading.GetPredictedMicroclimateReading(app)
			//if err != nil {
			//	w.WriteHeader(http.StatusInternalServerError)
			//	fmt.Fprintf(w, "Error in fetching: tmin predicted microclimate readings")
			//	return
			//}
			//fmt.Println("predicted y:", y)
			//
			//if len(*x) == len(*y) {
			//	for i, _ := range *x {
			//
			//		date := (*x)[i].Date
			//		gdd := ((*x)[i].Value+(*y)[i].Value)/2 - float32(culture.BaseTemperature)
			//
			//		gdds = append(gdds, models.GrowingDegreeDay{
			//			Date:  date,
			//			Value: gdd,
			//		})
			//
			//	}
			//}

		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(&gdds)
		w.Write(response)

	})
}
