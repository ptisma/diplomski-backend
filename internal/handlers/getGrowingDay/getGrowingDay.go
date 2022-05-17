package getGrowingDay

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

func GetGrowingDay(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		currentDate := time.Now().Format("20060102")
		fmt.Println("Current date:", currentDate)

		params := mux.Vars(r)

		cultureId, _ := strconv.ParseUint(params["cultureId"], 10, 32)

		urlParams := r.URL.Query()

		fromDate, _ := time.Parse("20060102", urlParams.Get("from"))
		toDate, _ := time.Parse("20060102", urlParams.Get("to"))
		locationId, _ := strconv.ParseUint(urlParams.Get("locationId"), 10, 32)
		fmt.Println("locationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

		culture := models.Culture{ID: uint32(cultureId)}

		tmax := models.Microclimate{Name: "tmax"}
		tmax.GetMicroclimateByName(app)
		fmt.Println("tmax:", tmax)
		microclimateReading := models.MicroclimateReading{
			MicroclimateID: tmax.ID,
			LocationID:     uint32(locationId),
			FromDate:       fromDate,
			ToDate:         toDate,
		}
		//fix for all x >= from not included
		x, _ := microclimateReading.GetMicroclimateReading(app)
		fmt.Println("x:", x)

		tmin := models.Microclimate{Name: "tmin"}
		tmin.GetMicroclimateByName(app)
		fmt.Println("tmin:", tmin)

		microclimateReading = models.MicroclimateReading{
			MicroclimateID: tmin.ID,
			LocationID:     uint32(locationId),
			FromDate:       fromDate,
			ToDate:         toDate,
		}

		y, _ := microclimateReading.GetMicroclimateReading(app)
		fmt.Println("y:", y)

		var gdds []struct {
			Date string
			Gdd  float32
		}

		if len(*x) == len(*y) {
			for i, _ := range *x {

				date := (*x)[i].Date
				gdd := ((*x)[i].Value+(*y)[i].Value)/2 - float32(culture.BaseTemperature)

				gdds = append(gdds, struct {
					Date string
					Gdd  float32
				}{Date: date, Gdd: gdd})

			}
		}

		lastTime, _ := time.Parse("2006-01-02", microclimateReading.Date)
		if toDate.After(lastTime) {
			microclimateReading := models.PredictedMicroclimateReading{
				MicroclimateID: tmax.ID,
				LocationID:     uint32(locationId),
				FromDate:       fromDate,
				ToDate:         toDate,
			}
			x, _ := microclimateReading.GetPredictedMicroclimateReading(app)
			fmt.Println("x:", x)

			microclimateReading = models.PredictedMicroclimateReading{
				MicroclimateID: tmin.ID,
				LocationID:     uint32(locationId),
				FromDate:       fromDate,
				ToDate:         toDate,
			}
			y, _ := microclimateReading.GetPredictedMicroclimateReading(app)
			fmt.Println("y:", y)

			if len(*x) == len(*y) {
				for i, _ := range *x {

					date := (*x)[i].Date
					gdd := ((*x)[i].Value+(*y)[i].Value)/2 - float32(culture.BaseTemperature)

					gdds = append(gdds, struct {
						Date string
						Gdd  float32
					}{Date: date, Gdd: gdd})

				}
			}

		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(gdds)
		w.Write(response)

	})
}
