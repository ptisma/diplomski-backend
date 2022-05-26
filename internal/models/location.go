package models

import (
	"apsim-api/pkg/application"
)

type Location struct {
	ID        uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Name      string  `gorm:"size:255;not null;unique" json:"name"`
	Latitude  float32 `gorm:"not null;" json:"latitude"`
	Longitude float32 `gorm:"not null;" json:"longitude"`
}

//add context?
func (l *Location) GetAllLLocations(app *application.Application) ([]Location, error) {
	var err error
	locations := []Location{}
	err = app.DB.Client.Debug().Model(&Location{}).Find(&locations).Error
	if err != nil {
		return []Location{}, err
	}
	return locations, err
}

func (l *Location) GetLocationById(app *application.Application) error {
	var err error
	err = app.DB.Client.Debug().First(l, l.ID).Error
	if err != nil {
		return err
	}
	return err
}
