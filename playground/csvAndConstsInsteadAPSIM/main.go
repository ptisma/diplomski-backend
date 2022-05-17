package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

var LOCATION_ID int = 1

var FROM_DATE string = "1990-01-06"

var TO_DATE string = "1990-01-20"

type Location struct {
	ID        uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name      string `gorm:"size:255;not null;unique" json:"name"`
	Latitude  float32
	Longitude float32
}

type Microclimate struct {
	ID   uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

type MicroclimateReading struct {
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

func main() {
	fmt.Println("Hello world")

	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})

	location := &Location{ID: uint32(LOCATION_ID)}

	db.First(location, location.ID)

	//csv file

	f, _ := os.CreateTemp(".", "csv*.csv")
	fmt.Println(filepath.Base(f.Name()))

	csvFile, _ := os.OpenFile(filepath.Base(f.Name()), os.O_APPEND|os.O_WRONLY, 0644) //Append mode

	results := &[]MicroclimateReading{}
	buff := []MicroclimateReading{}
	counter := 0
	csvFile.WriteString("year,day,radn,maxt,mint,rain,pan,vp,code\n")
	// batch size 100
	result := db.Where("location_id = ? AND date >= ? AND date <= ?", location.ID, FROM_DATE, TO_DATE).FindInBatches(results, 102, func(tx *gorm.DB, batch int) error {
		for _, result := range *results {
			// batch processing found records
			fmt.Println(result)
			counter += 1
			buff = append(buff, result)
			if counter == 6 {
				rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
				csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.Day(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
				csvFile.WriteString(csvRow)
				counter = 0
				buff = nil
			}

		}
		//tx.Save(&results)

		fmt.Println(tx.RowsAffected) // number of records in this batch

		fmt.Println(batch) // Batch 1, 2, 3

		// returns error will stop future batches
		return nil
	})

	fmt.Println(result.Error)

	//consts file

	f2, _ := os.CreateTemp(".", "const*.csv")
	fmt.Println(filepath.Base(f.Name()))

	constFile, _ := os.OpenFile(filepath.Base(f2.Name()), os.O_APPEND|os.O_WRONLY, 0644)
	constFile.WriteString(fmt.Sprintf("location = %s\n", location.Name))
	constFile.WriteString(fmt.Sprintf("latitude = %.2f (DECIMAL DEGREES)\n", location.Latitude))
	constFile.WriteString(fmt.Sprintf("longitude = %.2f (DECIMAL DEGREES)\n", location.Longitude))

}
