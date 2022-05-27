package models

type Culture struct {
	ID                  uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name                string `gorm:"size:255;not null;unique" json:"name"`
	GrowingDegreeDayMin int    `json:"growing_degree_day_min"`
	GrowingDegreeDayMax int    `json:"growing_degree_day_max"`
	BaseTemperature     int    `json:"base_temperature"`
}
