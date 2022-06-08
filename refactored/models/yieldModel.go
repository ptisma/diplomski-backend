package models

type Yield struct {
	Year  int32    `json:"year"`
	Yield float32  `json:"value"`
	Dates []string `gorm:"-:all" json:"-"`
}
