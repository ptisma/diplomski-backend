package models

import "apsim-api/pkg/application"

type Culture struct {
	ID                  uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name                string `gorm:"size:255;not null;unique" json:"name"`
	GrowingDegreeDayMin int    `json:"growing_degree_day_min"`
	GrowingDegreeDayMax int    `json:"growing_degree_day_max"`
	BaseTemperature     int    `json:"base_temperature"`
}

func (l *Culture) GetAllCultures(app *application.Application) (*[]Culture, error) {
	var err error
	cultures := &[]Culture{}
	err = app.DB.Client.Debug().Model(&Culture{}).Find(cultures).Error
	if err != nil {
		return &[]Culture{}, err
	}
	return cultures, err
}

func (l *Culture) GetCultureById(app *application.Application) error {
	var err error
	err = app.DB.Client.Debug().Model(&Culture{}).First(l, l.ID).Error
	return err
}
