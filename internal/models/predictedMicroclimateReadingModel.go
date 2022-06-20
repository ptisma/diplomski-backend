package models

import "time"

type PredictedMicroclimateReading struct {
	ID             uint32       `gorm:"primary_key;auto_increment" json:"-"`
	MicroclimateID uint32       `gorm:"size:255;not null;index" json:"-"`
	Microclimate   Microclimate `gorm:"foreignKey:MicroclimateID;references:ID" json:"-"`
	LocationID     uint32       `gorm:"size:255;not null;index" json:"-"`
	Location       Location     `gorm:"foreignKey:LocationID;references:ID" json:"-"`
	//Date           time.Time `gorm:"not null" json:"date"`
	Date  string  `gorm:"not null" json:"date"`
	Value float32 `gorm:"not null" json:"value"`
	//strings or time
	FromDate time.Time `gorm:"-:all" json:"-" `
	ToDate   time.Time `gorm:"-:all" json:"-"`
}
