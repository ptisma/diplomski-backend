package utils

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BackgroundUpdate(app *application.Application) {

	for {
		//Fetch all the locations
		var locations *[]models.Location
		location := &models.Location{}
		locations, err := location.GetAllLLocations(app)
		if err != nil {
			fmt.Println("Error in fetching all locations:", err)
		}
		fmt.Println(*locations)

		//Take the currentTime -2 days since external API updates slow
		currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -2)
		fmt.Println(currentTime)
		//For each location check the last date of any microclimate parameter they are all updated at the same time on external api
		for _, l := range *locations {
			fmt.Println(l)
			microclimateReading := models.MicroclimateReading{LocationID: l.ID}
			err := microclimateReading.GetLatestMicroclimateReading(app)
			if err != nil {
				fmt.Println("Error in fetching latest microclimateReading:", err)
				//Microclimate reading does not exist
				//Get a default one
				microclimateReading = models.MicroclimateReading{
					LocationID: location.ID,
					Date:       "1989-12-31",
				}
			}
			targetTime, err := time.Parse("2006-01-02", microclimateReading.Date)
			if err != nil {
				fmt.Println("Error in parsing date in latest microclimateReading:", err)
				continue
			}
			fmt.Println("targetTime:", targetTime)

			//Difference in days from current time and lastest time in DB
			diff := currentTime.Sub(targetTime) / (24 * time.Hour)
			fmt.Println("diff:", int(diff))
			if diff < 1 {
				fmt.Println("Less than 1 days")
				continue
			}

			start := targetTime.AddDate(0, 0, 1).Format("20060102")
			stop := targetTime.AddDate(0, 0, int(diff)).Format("20060102")
			fmt.Println("start:", start, "stop:", stop)

			//External REST api
			url := fmt.Sprintf("https://worldmodel.csiro.au/gclimate?lat=%f&lon=%f&format=csv&start=%s&stop=%s", l.Latitude, l.Longitude, start, stop)
			fmt.Println("URL:", url)

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
					fmt.Println("Line:", line)
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
						err = newMicroclimateReading.CreateMicroclimateReading(app)
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
				//if len(microclimateReadings)== 100 {
				//	mR.CreateMicroclimateReadingBatch(app, microclimateReadings)
				//	microclimateReadings = nil
				//}
			}
			//if len(microclimateReadings) != 0 {
			//
			//}
			//mR.CreateMicroclimateReading()
			//Skip second location for now
			//break

		}
		time.Sleep(24 * time.Hour)
	}

}
