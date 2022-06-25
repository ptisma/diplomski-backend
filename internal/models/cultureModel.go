package models

type Culture struct {
	ID              uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Name            string  `gorm:"size:255;not null;unique" json:"name"`
	BaseTemperature float32 `json:"base_temperature"`
}
