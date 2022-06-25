package router

import (
	"apsim-api/internal/infra/application"
	"apsim-api/internal/router/middlewares"
	"apsim-api/internal/serviceContainer"
	"github.com/gorilla/mux"
	"sync"
)

type IMuxRouter interface {
	InitRouter() *mux.Router
}

//router helper kernel function
type rkernel struct {
	app application.Application
}

func (r *rkernel) InitRouter() *mux.Router {

	//get a Singleton instance of service container
	serviceContainer := serviceContainer.GetServiceContainer(r.app)

	//controllers
	cultureController := serviceContainer.InjectCultureController()
	locationController := serviceContainer.InjectLocationController()
	microclimateController := serviceContainer.InjectMicroclimateController()
	microclimateReadingController := serviceContainer.InjectMicroclimateReadingController()
	yieldController := serviceContainer.InjectYieldController()
	growingDegreeDayController := serviceContainer.InjectGrowingDegreeDayController()
	mux := mux.NewRouter()

	//endpoints
	mux.HandleFunc("/culture", cultureController.GetCultures).Methods("GET")
	mux.HandleFunc("/location", locationController.GetLocations).Methods("GET")
	mux.HandleFunc("/microclimate", microclimateController.GetMicroclimates).Methods("GET")

	//middlewares
	locationRouter := mux.PathPrefix("/location/{locationId}").Subrouter()
	locationRouter.Use(middlewares.LocationMiddleware)
	cultureRouter := locationRouter.PathPrefix("/culture/{cultureId}").Subrouter()
	cultureRouter.Use(middlewares.CultureMiddleware, middlewares.DatesMiddleware)
	microclimateRouter := locationRouter.PathPrefix("/microclimate/{microclimateId}").Subrouter()
	microclimateRouter.Use(middlewares.MicroclimateMiddleware, middlewares.DatesMiddleware)
	//endpoints
	locationRouter.HandleFunc("/microclimate/all/period", microclimateReadingController.GetMicroclimateReadingPeriod).Methods("GET")
	microclimateRouter.HandleFunc("", microclimateReadingController.GetMicroclimateReadings).Methods("GET")
	cultureRouter.HandleFunc("/yield", yieldController.GetYield).Methods("GET")
	cultureRouter.HandleFunc("/gdd", growingDegreeDayController.GetGrowingDegreeDays).Methods("GET")

	return mux

}

// Singleton
var (
	m          *rkernel
	routerOnce sync.Once
)

func GetMuxRouter(app application.Application) IMuxRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &rkernel{app: app}
		})
	}
	return m
}
