package main

import (
	"apsim-api/internal/models"
	"bufio"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var MICROCLIMATE_ID uint32 = 1
var LOCATION_ID uint32 = 1

func loadMicroCLimateReading(microclimateId, locationId uint32, date string, value float32) models.MicroclimateReading {

	return models.MicroclimateReading{
		MicroclimateID: microclimateId,
		LocationID:     locationId,
		Date:           date,
		Value:          value,
	}
}

func main() {

	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{}) //baza.db se odnosi na working dir

	file, err := os.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\files\\petar-dataset\\zagreb-1990.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	counter := 0
	f, _ := os.CreateTemp(".", "logLoading*.txt")
	logFile, _ := os.OpenFile(filepath.Base(f.Name()), os.O_APPEND|os.O_WRONLY, 0644)
	for scanner.Scan() {
		if counter > 3 {
			//date, rad, tmax, tmin, rain, rh,wind
			//mm/dd/yyyy
			line := strings.Split(scanner.Text(), ",")
			if len(line) == 1 {
				continue
			}
			fmt.Println(line)
			time2, err := time.Parse("1/2/2006", line[0])
			if err != nil {
				fmt.Println(err)
				logFile.WriteString(err.Error() + "\n")
				continue
			}
			for i := 1; i < len(line); i++ {
				val, err := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
				if err != nil {
					fmt.Println(err)
					logFile.WriteString(err.Error() + "\n")
					continue
				}
				res := loadMicroCLimateReading(uint32(i), LOCATION_ID, time2.String(), float32(val))
				fmt.Println(res.MicroclimateID, res.LocationID, res.Date, res.Value)
				//db.Create(&res)
				fmt.Println("Creating")
				tx := db.Begin()
				tx.Create(&res)
				//time.Sleep(1 * time.Minute)
				tx.Commit()
				fmt.Println("Finished")
			}
		}

		counter += 1

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
