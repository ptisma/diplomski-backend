package models

import (
	"apsim-api/pkg/application"
	"fmt"
	"time"
)

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

func (l *PredictedMicroclimateReading) GetPredictedMicroclimateReading(app *application.Application) ([]PredictedMicroclimateReading, error) {
	var err error

	microclimates := []PredictedMicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"
	//err = app.DB.Client.Debug().Model(&MicroclimateReading{}).Where(queryStr, l.MicroclimateID, l.LocationID, l.ToDate, l.FromDate).Find(microclimates).Error

	//app.DB.Client.Debug().Model(&MicroclimateReading{}).Where(queryStr, l.MicroclimateID, l.LocationID, l.ToDate, l.FromDate).Find(microclimates)
	//app.DB.Client.Debug().Model(&MicroclimateReading{}).Find(microclimates)
	//moze samo group by date
	app.DB.Client.Debug().Preload("Location").Model(&PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, l.MicroclimateID, l.LocationID, l.ToDate.Format("2006-01-02"), l.FromDate.Format("2006-01-02")).Group("microclimate_id,location_id,date").Order("date").Find(&microclimates)
	if err != nil {
		return []PredictedMicroclimateReading{}, err
	}

	fmt.Println(microclimates)

	//for _, value := range *microclimates {
	//	fmt.Println(value.LocationID)
	//	fmt.Println(value.Location)
	//	fmt.Println(value.MicroclimateID)
	//	fmt.Println(value.Microclimate)
	//	fmt.Println(value.ID)
	//	fmt.Println(value.Date)
	//	fmt.Println(value.Value)
	//
	//}
	return microclimates, err
}

func (l *PredictedMicroclimateReading) GetLatestPredictedMicroclimateReading(app *application.Application) error {
	var err error

	microclimates := &[]PredictedMicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ?"

	err = app.DB.Client.Debug().Preload("Location").Model(&PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, l.MicroclimateID, l.LocationID).Order("date desc").Find(microclimates).Error

	return err

}
