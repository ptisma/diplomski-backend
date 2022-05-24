package models

import "apsim-api/pkg/application"

type Soil struct {
	ID         uint32   `gorm:"primary_key;auto_increment" json:"id"`
	Name       string   `gorm:"size:255;not null;unique" json:"name"`
	Data       string   `gorm:"not null" json:"data"`
	LocationID uint32   `gorm:"size:255;not null;index" json:"-"`
	Location   Location `gorm:"foreignKey:LocationID;references:ID" json:"location"`
}

func (l *Soil) GetSoilByLocationId(app *application.Application) error {
	var err error
	err = app.DB.Client.Debug().Model(&Soil{}).Preload("Location").Where("location_id = ?", l.LocationID).Find(l).Error
	return err
}
