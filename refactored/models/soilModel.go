package models

type Soil struct {
	ID         uint32   `gorm:"primary_key;auto_increment" json:"id"`
	Name       string   `gorm:"size:255;not null;unique" json:"name"`
	Data       string   `gorm:"not null" json:"data"`
	LocationID uint32   `gorm:"size:255;not null;index" json:"-"`
	Location   Location `gorm:"foreignKey:LocationID;references:ID" json:"location"`
}
