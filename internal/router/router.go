package router

import (
	"apsim-api/internal/handlers/getCultures"
	"apsim-api/internal/handlers/getGrowingDay"
	"apsim-api/internal/handlers/getLocations"
	"apsim-api/internal/handlers/getMicroclimate"
	"apsim-api/internal/handlers/getMicroclimateReading"
	"apsim-api/internal/handlers/getYield"
	"apsim-api/internal/middlewares"
	"apsim-api/pkg/application"
	"github.com/gorilla/mux"
)

func GetRouter(application *application.Application) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/location", getLocations.GetLocations(application)).Methods("GET")

	router.Handle("/culture", getCultures.GetCultures(application)).Methods("GET")

	router.Handle("/microclimate", getMicroclimate.GetMicroclimate(application)).Methods("GET")

	//
	locationRouter := router.PathPrefix("/location/{locationId}").Subrouter()
	locationRouter.Use(middlewares.LocationMiddleware)

	cultureRouter := locationRouter.PathPrefix("/culture/{cultureId}").Subrouter()
	cultureRouter.Use(middlewares.CultureMiddleware)

	microclimateRouter := locationRouter.PathPrefix("/microclimate/{microclimateId}").Subrouter()
	microclimateRouter.Use(middlewares.MicroclimateMiddleware)

	//mux.Handle("/microclimate/{microclimateId}", getMicroclimateReading.GetMicroclimateReading(application)).Methods("GET")
	//
	//mux.Handle("/location/{locationId}/culture/{cultureId}/yield", getYield.GetYield2(application)).Methods("GET")

	cultureRouter.Handle("/gdd", getGrowingDay.GetGrowingDay(application)).Methods("GET")
	microclimateRouter.Handle("", getMicroclimateReading.GetMicroclimateReading(application)).Methods("GET")
	cultureRouter.Handle("/yield", getYield.GetYield2(application)).Methods("GET")

	return router
}
