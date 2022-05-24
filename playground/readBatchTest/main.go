package main

import (
	"apsim-api/internal/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func main() {
	db, err := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []models.MicroclimateReading
	counter := 0
	//blokirajuca
	//callback nije u zasebnoj gorutini
	result := db.Debug().Where("location_id = ? AND date >= ? AND date <= ?", 1, "1990-01-01", "1990-03-30").FindInBatches(&results, 10, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			fmt.Println(result)

		}
		counter += 1
		if counter == 2 {
			time.Sleep(5 * time.Second)
		}

		return nil
	})
	fmt.Println("Result:", result.Error)
}
