package router

import (
	"apsim-api/refactored/router/middlewares"
	"apsim-api/refactored/serviceContainer"
	"github.com/gorilla/mux"
	"sync"
)

type IMuxRouter interface {
	InitRouter() *mux.Router
}

type router struct{}

func (r *router) InitRouter() *mux.Router {

	cultureController := serviceContainer.ServiceContainer().InjectCultureController()
	locationController := serviceContainer.ServiceContainer().InjectLocationController()
	microclimateController := serviceContainer.ServiceContainer().InjectMicroclimateController()
	microclimateReadingController := serviceContainer.ServiceContainer().InjectMicroclimateReadingController()
	yieldController := serviceContainer.ServiceContainer().InjectYieldController()
	growingDegreeDayController := serviceContainer.ServiceContainer().InjectGrowingDegreeDayController()
	mux := mux.NewRouter()

	mux.HandleFunc("/culture", cultureController.GetCultures).Methods("GET")
	mux.HandleFunc("/location", locationController.GetLocations).Methods("GET")
	mux.HandleFunc("/microclimate", microclimateController.GetMicroclimates).Methods("GET")

	//middlewares
	locationRouter := mux.PathPrefix("/location/{locationId}").Subrouter()
	locationRouter.Use(middlewares.LocationMiddleware)
	cultureRouter := locationRouter.PathPrefix("/culture/{cultureId}").Subrouter()
	cultureRouter.Use(middlewares.CultureMiddleware)
	microclimateRouter := locationRouter.PathPrefix("/microclimate/{microclimateId}").Subrouter()
	microclimateRouter.Use(middlewares.MicroclimateMiddleware)

	//nested endpoints
	locationRouter.HandleFunc("/microclimate/all/period", microclimateReadingController.GetMicroclimateReadingPeriod).Methods("GET")
	microclimateRouter.HandleFunc("", microclimateReadingController.GetMicroclimateReadings).Methods("GET")
	cultureRouter.HandleFunc("/yield", yieldController.GetYield).Methods("GET")
	cultureRouter.HandleFunc("/gdd", growingDegreeDayController.GetGrowingDegreeDays).Methods("GET")

	return mux

}

var (
	m          *router
	routerOnce sync.Once
)

func MuxRouter() IMuxRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}
