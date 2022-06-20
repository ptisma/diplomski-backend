package main

import (
	"apsim-api/internal/models"
	"bufio"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func loadPredictedMicroCLimateReading(microclimateId, locationId uint32, date string, value float32) models.PredictedMicroclimateReading {

	return models.PredictedMicroclimateReading{
		MicroclimateID: microclimateId,
		LocationID:     locationId,
		Date:           date,
		Value:          value,
	}
}

func main() {

	LOCATION_ID := 2

	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})

	file, err := os.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\files\\petar-dataset\\weather\\osijek\\osijek-2010-predicted.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	flag := false
	counter := 0
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		if counter > 0 {
			line := strings.Split(scanner.Text(), ",")
			//fmt.Println(line[0])

			date, err := time.Parse("1/2/2006 3:04", line[0])
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Printf("%s", date.Format("2006-01-02"))
			for i := 1; i < len(line); i++ {
				val, err := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
				if err != nil {
					fmt.Println(err)
					flag = true
					break
				}
				fmt.Printf(",%f", val)
				reading := loadPredictedMicroCLimateReading(uint32(i), uint32(LOCATION_ID), date.Format("2006-01-02"), float32(val))
				db.CreateInBatches(&reading, 100)
			}
			if flag {
				break
			}
			fmt.Println()
		}

		counter += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
