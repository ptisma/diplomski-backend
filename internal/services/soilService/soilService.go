package soilService

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
)

type SoilService struct {
	app *application.Application
}

func GetSoilService(app *application.Application) *SoilService {
	return &SoilService{
		app: app,
	}
}

func (ss *SoilService) GetSoil(locationId int) (models.Soil, error) {
	var err error
	soil := models.Soil{LocationID: uint32(locationId)}
	err = soil.GetSoilByLocationId(ss.app)

	return soil, err

}
