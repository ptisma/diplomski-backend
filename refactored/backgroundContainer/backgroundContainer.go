package backgroundContainer

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"apsim-api/refactored/repositories"
	"apsim-api/refactored/services"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Background interface {
	UpdateMicroclimateReadings()
}

type background struct {
	I1 interfaces.IMicroclimateReadingService
	I2 interfaces.ILocationService
}

func NewBackground() Background {
	db, _ := gorm.Open(sqlite.Open("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\baza.db"), &gorm.Config{})
	return &background{
		I1: &services.MicroclimateReadingService{&repositories.MicroclimateReadingRepository{db}, &repositories.PredictedMicroclimateReadingRepository{db}},
		I2: &services.LocationService{&repositories.LocationRepository{db}},
	}

}

func (b *background) UpdateMicroclimateReadings() {

	for {
		//Fetch all the locations
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		locations, err := b.I2.GetAllLocations(ctx)
		if err != nil {
			fmt.Println("Error in fetching all locations:", err)
		}
		//fmt.Println(locations)

		//Take the currentTime -3 days since external API updates slow
		currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -3)
		//fmt.Println(currentTime)
		//For each location check the last date of any microclimate parameter they are all updated at the same time on external api
		for _, l := range locations {
			fmt.Println(l)
			microclimateReading, err := b.I1.GetLatestMicroClimateReading(ctx, int(l.ID))
			if err != nil {
				fmt.Println("Error in fetching latest microclimateReading:", err)
				//Microclimate reading does not exist
				//Get a default one
				microclimateReading = models.MicroclimateReading{
					LocationID: l.ID,
					Date:       "1989-12-31",
				}
			}
			targetTime, err := time.Parse("2006-01-02", microclimateReading.Date)
			if err != nil {
				fmt.Println("Error in parsing date in latest microclimateReading:", err)
				continue
			}
			//fmt.Println("targetTime:", targetTime)

			//Difference in days from current time and lastest time in DB
			diff := currentTime.Sub(targetTime) / (24 * time.Hour)
			fmt.Println("diff:", int(diff))
			if diff < 1 {
				fmt.Println("Less than 1 days")
				continue
			}

			start := targetTime.AddDate(0, 0, 1).Format("20060102")
			stop := targetTime.AddDate(0, 0, int(diff)).Format("20060102")
			//fmt.Println("start:", start, "stop:", stop)

			//External REST api
			url := fmt.Sprintf("https://worldmodel.csiro.au/gclimate?lat=%f&lon=%f&format=csv&start=%s&stop=%s", l.Latitude, l.Longitude, start, stop)
			//fmt.Println("URL:", url)

			var buf bytes.Buffer
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error in HTTP request to external REST endpoint:", err)
			}
			io.Copy(&buf, resp.Body)

			scanner := bufio.NewScanner(&buf)
			counter := 0
			var doubleBreak = false
			//var microclimateReadings []models.MicroclimateReading
			//var mR *models.MicroclimateReading
			for scanner.Scan() {
				if counter > 3 {
					//Line example: 2022-01-01, 6.20,11.33, 0.91, 3.06,77.88, 0.99
					line := strings.Split(scanner.Text(), ",")
					//fmt.Println("Line:", line)
					//Parse date
					rowTime := line[0]
					checkTime, err := time.Parse("2006-01-02", rowTime)
					if err != nil {
						fmt.Println("Error in parsing time in row:", err)
						break
					}
					//Check time, sometimes API decides to stick random values at the end of the file such as "2022-03-31,,,,12.83,,"
					diffCheck := checkTime.Sub(targetTime) / (24 * time.Hour)
					if diffCheck < 1 {
						fmt.Println("checkTime:", checkTime)
						fmt.Println("diffCheck:", diffCheck)
						break
					}
					//parse microclimate parameters
					for i := 1; i < len(line); i++ {

						//Parse value
						//create the microclimate reading based on the i-th position in row and i-th value
						//row format follows the DB microclimates format
						val, err := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
						if err != nil {
							fmt.Println("Error in parsing float value in row:", err)
							doubleBreak = true
							break
						}
						//Default non existing value in REST API
						if val == -999.00 {
							fmt.Println("Error in parsing float value -999.00")
							doubleBreak = true
							break
						}
						newMicroclimateReading := models.MicroclimateReading{
							MicroclimateID: uint32(i),
							LocationID:     l.ID,
							Date:           rowTime,
							Value:          float32(val),
						}

						//microclimateReadings = append(microclimateReadings, newMicroclimateReading)
						fmt.Println("newMicroclimateReading:", newMicroclimateReading.MicroclimateID, newMicroclimateReading.LocationID, newMicroclimateReading.Date, newMicroclimateReading.Value)
						err = b.I1.CreateMicroclimateReading(ctx, i, int(l.ID), rowTime, float32(val))
						if err != nil {
							fmt.Println("Error in creating newMicroclimateReading in DB:", err)
						}

					}
					if doubleBreak == true {
						fmt.Println("Stopping due to false line", line)
						break
					}
				}
				//Skip first 4 rows to get to the CSV row values
				counter += 1

			}
		}
		time.Sleep(24 * time.Hour)
	}
}
