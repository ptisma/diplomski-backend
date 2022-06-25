package backgroundContainer

import (
	"apsim-api/internal/infra/application"
	"apsim-api/internal/interfaces"
	"apsim-api/internal/repositories"
	"apsim-api/internal/services"
)

type BackgroundWorker interface {
	GetLocationService() interfaces.ILocationService
	GetMicroclimateReadingService() interfaces.IMicroclimateReadingService
}

type bg struct {
	I1 interfaces.IMicroclimateReadingService
	I2 interfaces.ILocationService
}

func (b *bg) GetMicroclimateReadingService() interfaces.IMicroclimateReadingService {
	return b.I1
}

func (b *bg) GetLocationService() interfaces.ILocationService {

	return b.I2

}

func NewBackgroundWorker(app application.Application) BackgroundWorker {
	return &bg{
		I1: &services.MicroclimateReadingService{&repositories.MicroclimateReadingRepository{app.GetDB()}, &repositories.PredictedMicroclimateReadingRepository{app.GetDB()}},
		I2: &services.LocationService{&repositories.LocationRepository{app.GetDB()}},
	}

}
