package microclimateService

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
)

type MicroclimateService struct {
	app *application.Application
}

func (ms *MicroclimateService) GetMicroclimate(name string) (models.Microclimate, error) {
	tmax := models.Microclimate{Name: name}
	err := tmax.GetMicroclimateByName(ms.app)
	return tmax, err

}
func (ms *MicroclimateService) GetMicroclimates() ([]models.Microclimate, error) {
	m := models.Microclimate{}
	return m.GetAllMicroclimate(ms.app)
}

func GetMicroclimateService(app *application.Application) *MicroclimateService {
	return &MicroclimateService{
		app: app,
	}
}
