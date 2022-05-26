package main

import (
	"apsim-api/internal/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("./baza.db"), &gorm.Config{})

	microclimates := []models.MicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"
	//err = app.DB.Client.Debug().Model(&MicroclimateReading{}).Where(queryStr, l.MicroclimateID, l.LocationID, l.ToDate, l.FromDate).Find(microclimates).Error

	//app.DB.Client.Debug().Model(&MicroclimateReading{}).Where(queryStr, l.MicroclimateID, l.LocationID, l.ToDate, l.FromDate).Find(microclimates)
	//app.DB.Client.Debug().Model(&MicroclimateReading{}).Find(microclimates)
	//moze samo group by date
	err := db.Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, 1, 1, "1989-01-01", "1989-12-31").Group("microclimate_id,location_id,date").Order("date").Find(&microclimates).Error

	//error vraca samo akd je first, last and take
	fmt.Println(err)
	fmt.Println(microclimates)
}
