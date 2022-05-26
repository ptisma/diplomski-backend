package models

import (
	"apsim-api/pkg/application"
)

type Microclimate struct {
	ID   uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
	Unit string `gorm:"" json:"unit"`
}

//add context?
func (l *Microclimate) GetAllMicroclimate(app *application.Application) ([]Microclimate, error) {
	var err error
	microclimates := []Microclimate{}
	err = app.DB.Client.Debug().Model(&Microclimate{}).Find(&microclimates).Error
	if err != nil {
		return []Microclimate{}, err
	}
	return microclimates, err
}

func (l *Microclimate) GetMicroclimateByName(app *application.Application) error {

	err := app.DB.Client.Debug().Model(&Microclimate{}).Where("name = ?", l.Name).First(l).Error
	return err

}
