package main

import (
	"apsim-api/internal/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	batchSize := 100
	results := []models.MicroclimateReading{}
	locationID := 2
	fromDate := "1989-01-01"
	toDate := "1989-04-01"
	var (
		rowsAffected int64
		batch        int
	)

	fmt.Println("Krecem")

	for {
		results = nil
		result := db.Limit(batchSize).Offset(batch*batchSize).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate, toDate).Order("date, microclimate_id").Find(&results)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			//do something
			for _, result := range results {
				fmt.Println(result.Date)
			}
			fmt.Println("Batch rezultati")
			fmt.Println()

			if int(result.RowsAffected) < batchSize {
				break
			}

		} else if result.Error != nil {
			err = result.Error
			break
		} else if result.RowsAffected == 0 {
			break
		}

		//if int(result.RowsAffected) < batchSize {
		//	for _, result := range results {
		//		fmt.Println(result.Date)
		//	}
		//	fmt.Println("Batch rezultati")
		//	fmt.Println("Zadnji batch")
		//	break
		//}
	}

	fmt.Println(err)
	fmt.Println("Gotov")

}
