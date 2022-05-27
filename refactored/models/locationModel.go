package models

type Location struct {
	ID        uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Name      string  `gorm:"size:255;not null;unique" json:"name"`
	Latitude  float32 `gorm:"not null;" json:"latitude"`
	Longitude float32 `gorm:"not null;" json:"longitude"`
}
