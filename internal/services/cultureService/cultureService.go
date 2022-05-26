package cultureService

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
)

type CultureService struct {
	app *application.Application
}

func GetCultureService(app *application.Application) *CultureService {
	return &CultureService{
		app: app,
	}
}

func (cs *CultureService) GetCulture(cultureId int) (models.Culture, error) {
	var err error
	culture := models.Culture{ID: uint32(cultureId)}
	err = culture.GetCultureById(cs.app)

	return culture, err

}

func (cs *CultureService) GetCultures() ([]models.Culture, error) {

	culture := models.Culture{}
	return culture.GetAllCultures(cs.app)

}
