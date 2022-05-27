package models

type Microclimate struct {
	ID   uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
	Unit string `gorm:"" json:"unit"`
}
