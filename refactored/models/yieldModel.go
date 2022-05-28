package models

type Yield struct {
	Year  int32    `json:"year"`
	Yield float32  `json:"yield"`
	Dates []string `gorm:"-:all" json:"-"`
}
