package main

import (
	"apsim-api/internal/models"
	"bufio"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"time"
)

func loadPredictedMicroClimateReading(microclimateId, locationId uint32, date string, value float32) models.PredictedMicroclimateReading {

	return models.PredictedMicroclimateReading{
		MicroclimateID: microclimateId,
		LocationID:     locationId,
		Date:           date,
		Value:          value,
	}
}

func main() {

	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{}) //baza.db se odnosi na working dir

	firstTime := time.Date(1990, 4, 16, 0, 0, 0, 0, time.UTC)

	lastTime := time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC)

	//formattedFirstTime := firstTime.Format("2006-01-02")
	//
	//fmt.Println(formattedFirstTime)
	//
	//currentTime := firstTime

	file, _ := os.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\files\\petar-dataset\\zagreb-1990.csv")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		lineTime, _ := time.Parse("1/2/2006", line[0])
		if (lineTime.After(firstTime) || lineTime.Equal(firstTime)) && (lineTime.Before(lastTime) || lineTime.Equal(lastTime)) {
			lineTime = lineTime.AddDate(32, 0, 0)
			fmt.Println(line)
			for i := 1; i < len(line); i++ {
				val, _ := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
				res := loadPredictedMicroClimateReading(uint32(i), 1, lineTime.Format("2006-01-02"), float32(val))
				db.Create(&res)
			}

		}
	}

	//for {
	//	if currentTime.Before(lastTime) {
	//
	//		for i := 1; i <= 6; i++ {
	//			res := loadPredictedMicroClimateReading(uint32(i), 1, currentTime.Format("2006-01-02"), 69)
	//			db.Create(&res)
	//		}
	//		currentTime = currentTime.AddDate(0, 0, 1)
	//		fmt.Println(currentTime)
	//	} else {
	//		break
	//	}
	//}

}
