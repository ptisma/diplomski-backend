package router

import (
	"apsim-api/internal/handlers/getCultures"
	"apsim-api/internal/handlers/getGrowingDay"
	"apsim-api/internal/handlers/getLocations"
	"apsim-api/internal/handlers/getMicroclimate"
	"apsim-api/internal/handlers/getMicroclimateReading"
	"apsim-api/internal/handlers/getYield"
	"apsim-api/pkg/application"
	"github.com/gorilla/mux"
)

func GetRouter(application *application.Application) *mux.Router {
	mux := mux.NewRouter()

	mux.Handle("/location", getLocations.GetLocations(application)).Methods("GET")

	mux.Handle("/culture", getCultures.GetCultures(application)).Methods("GET")

	mux.Handle("/microclimate", getMicroclimate.GetMicroclimate(application)).Methods("GET")

	//mux.Handle("/location/{locationId}/microclimate/{microclimateId}", getMicroclimateReading.GetReading(application)).Methods("GET")

	mux.Handle("/microclimate/{microclimateId}", getMicroclimateReading.GetMicroclimateReading(application)).Methods("GET")

	mux.Handle("/location/{locationId}/culture/{cultureId}/yield", getYield.GetYield2(application)).Methods("GET")

	mux.Handle("/culture/{cultureId}/gdd", getGrowingDay.GetGrowingDay(application)).Methods("GET")

	return mux
}
