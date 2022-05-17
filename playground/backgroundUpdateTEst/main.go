package main

import (
	"apsim-api/internal/models"
	"bufio"
	"bytes"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func loadMicroCLimateReading(microclimateId, locationId uint32, date string, value float32) models.MicroclimateReading {

	return models.MicroclimateReading{
		MicroclimateID: microclimateId,
		LocationID:     locationId,
		Date:           date,
		Value:          value,
	}
}

func main() {

	//currentTime := time.Now().AddDate(0, 0, -2)
	currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -2)
	fmt.Println("currentTime:", currentTime)
	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{})

	var locations []models.Location
	db.Find(&locations)
	var microclimates []models.Microclimate
	db.Find(&microclimates)

	//fmt.Println(locations)
	//fmt.Println(microclimates)

	queryStr := "location_id = ?"
	//For each location check the last date of any microclimate parameter they are all updated at the same time on external api
	for _, l := range locations {
		var microclimateReading models.MicroclimateReading
		db.Preload("Location").Preload("Microclimate").Where(queryStr, l.ID).Order("date desc").Find(&microclimateReading)
		//fmt.Println(microclimateReading)
		targetTime, _ := time.Parse("2006-01-02", microclimateReading.Date)
		fmt.Println("targetTime:", targetTime)

		diff := currentTime.Sub(targetTime) / (24 * time.Hour)

		fmt.Println("diff:", int(diff))
		if diff < 1 {
			break
		}
		start := targetTime.AddDate(0, 0, 1).Format("20060102")
		stop := targetTime.AddDate(0, 0, int(diff)).Format("20060102")
		//fmt.Println("start:", start, "stop:", stop)
		url := fmt.Sprintf("https://worldmodel.csiro.au/gclimate?lat=45.815080&lon=15.9819189&format=csv&start=%s&stop=%s", start, stop)
		fmt.Println("URL:", url)

		var buf bytes.Buffer
		resp, _ := http.Get(url)
		io.Copy(&buf, resp.Body)
		scanner := bufio.NewScanner(&buf)
		counter := 0
		for scanner.Scan() {
			if counter > 3 {
				//fmt.Println("Linija:", scanner.Text())
				line := strings.Split(scanner.Text(), ",")
				fmt.Println(line)
				for i := 1; i < len(line); i++ {
					time2 := line[0]
					checkTime, err := time.Parse("2006-01-02", time2)
					if err != nil {
						fmt.Println(err)
						break
					}
					diffCheck := checkTime.Sub(targetTime) / (24 * time.Hour)
					if diffCheck < 1 {
						fmt.Println("checkTime:", checkTime)
						fmt.Println("diffCheck:", diffCheck)
						break
					}
					val, err := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
					if err != nil {
						fmt.Println(err)
						continue
					}
					res := loadMicroCLimateReading(uint32(i), l.ID, time2, float32(val))
					fmt.Println(res.MicroclimateID, res.LocationID, res.Date, res.Value)
					err = db.Create(&res).Error
					if err != nil {
						fmt.Println(err)
					}

				}
			}
			//skipaj prva tri reda jer ne vraca cisti csv, nit json
			counter += 1
		}
	}

}
