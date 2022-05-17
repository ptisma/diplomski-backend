package main

import (
	"apsim-api/internal/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func main() {

	currentTime := time.Now().AddDate(0, 0, -2)
	fmt.Println(currentTime)
	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})

	var locations []models.Location
	db.Find(&locations)
	var microclimates []models.Microclimate
	db.Find(&microclimates)

	fmt.Println(locations)
	fmt.Println(microclimates)

	queryStr := "location_id = ?"
	//For each location check the last date of any microclimate parameter they are all updated at the same time on external api
	for _, l := range locations {
		var microclimateReading models.MicroclimateReading
		db.Preload("Location").Preload("Microclimate").Where(queryStr, l.ID).Order("date desc").Find(&microclimateReading)
		fmt.Println(microclimateReading)
		targetTime, _ := time.Parse("2006-01-02", microclimateReading.Date)
		fmt.Println(targetTime)

		diff := currentTime.Sub(targetTime) / (24 * time.Hour)
		fmt.Println(int(diff))
		for i := 0; i < int(diff); i++ {
			targetTime = targetTime.AddDate(0, 0, 1)
			fmt.Println("Fetching for:", targetTime)

		}
	}
}
