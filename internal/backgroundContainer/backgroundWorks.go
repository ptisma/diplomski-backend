package backgroundContainer

import (
	"apsim-api/internal/models"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type backgroundWork func(ctx context.Context, waitGroup *sync.WaitGroup, worker BackgroundWorker)

func UpdateMicroclimateReadings(ctx context.Context, waitGroup *sync.WaitGroup, b BackgroundWorker) {
	log.Println("UpdateMicroclimateReadings started")
	for {
		//Fetch all the locations
		locations, err := b.GetLocationService().GetAllLocations(ctx)
		if err != nil {
			log.Println("Error in fetching all locations:", err)
		}
		//log.Println(locations)

		//Take the currentTime -3 days since external API updates slow
		currentTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -3)
		//log.Println(currentTime)
		//For each location check the last date of any microclimate parameter they are all updated at the same time on external api
		for _, l := range locations {
			select {
			case <-ctx.Done():
				log.Println("UpdateMicroclimateReadings closing")
				waitGroup.Done()
				return
			default:
				//do work
			}
			log.Println("Updating location:", l.Name)
			microclimateReading, err := b.GetMicroclimateReadingService().GetLatestMicroClimateReading(ctx, int(l.ID))
			if err != nil {
				log.Println("Error in fetching latest microclimateReading:", err)
				//Microclimate reading does not exist
				//Get a default one
				microclimateReading = models.MicroclimateReading{
					LocationID: l.ID,
					Date:       "1989-12-31",
				}
			}
			targetTime, err := time.Parse("2006-01-02", microclimateReading.Date)
			if err != nil {
				log.Println("Error in parsing date in latest microclimateReading:", err)
				continue
			}
			//log.Println("targetTime:", targetTime)

			//Difference in days from current time and lastest time in DB
			diff := currentTime.Sub(targetTime) / (24 * time.Hour)
			//log.Println("diff:", int(diff))
			if diff < 1 {
				log.Println("Less than 1 days")
				continue
			}

			start := targetTime.AddDate(0, 0, 1).Format("20060102")
			stop := targetTime.AddDate(0, 0, int(diff)).Format("20060102")
			//log.Println("start:", start, "stop:", stop)

			//External REST api
			url := fmt.Sprintf("https://worldmodel.csiro.au/gclimate?lat=%f&lon=%f&format=csv&start=%s&stop=%s", l.Latitude, l.Longitude, start, stop)
			log.Println("URL:", url)

			var buf bytes.Buffer
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Error in HTTP request to external REST endpoint:", err)
				continue
			}
			io.Copy(&buf, resp.Body)

			scanner := bufio.NewScanner(&buf)
			counter := 0
			var doubleBreak = false
			previousTime := targetTime
			//API sometimes changes format, start parsing after see CSV headers
			//date,rad,tmax,tmin,rain,rh,wind
			//-,MJ/m2,oC,oC,mm,%,m/s
			var csvHeader1Flag, csvHeader2Flag bool
			csvHeader1 := "date,rad,tmax,tmin,rain,rh,wind"
			csvHeader2 := "-,MJ/m2,oC,oC,mm,%,m/s"
			for scanner.Scan() {
				//log.Println(scanner.Text())
				if csvHeader1Flag == false {
					if strings.Trim(scanner.Text(), " ") == csvHeader1 {
						csvHeader1Flag = true
						continue
					}
				}
				if csvHeader1Flag == true && csvHeader2Flag == false {
					if strings.Trim(scanner.Text(), " ") == csvHeader2 {
						csvHeader2Flag = true
						continue
					}
				}
				if csvHeader1Flag == true && csvHeader2Flag == true {
					//log.Println("Started parsing")
					//Line example: 2022-01-01, 6.20,11.33, 0.91, 3.06,77.88, 0.99
					line := strings.Split(scanner.Text(), ",")
					//log.Println("Line:", line)
					//Parse date
					rowTime := line[0]
					checkTime, err := time.Parse("2006-01-02", rowTime)

					if err != nil {
						log.Println("Error in parsing time in row:", err)
						break
					}
					//log.Println("checkTime:", checkTime.Format("2006-01-02"))
					//Check time, sometimes API decides to stick random values at the end of the file such as "2022-03-31,,,,12.83,,"
					diffCheck := checkTime.Sub(targetTime) / (24 * time.Hour)
					if diffCheck < 1 {
						log.Println("row date is less than or same as last in DB, checkTime:", checkTime.Format("2006-01-02"), "targetTime:", targetTime.Format("2006-01-02"), diffCheck)
						break
					}
					//Check last time, sometimes API decides to skip some dates
					diffCheck = checkTime.Sub(previousTime) / (24 * time.Hour)
					if diffCheck > 1 {
						log.Println("row date is greater than 1 day than last row date, checkTime:", checkTime.Format("2006-01-02"), "previousTime:", previousTime.Format("2006-01-02"))
						break
					}
					if diffCheck < 1 {
						log.Println("row date is less than last row date, checkTime:", checkTime.Format("2006-01-02"), "previousTime:", previousTime.Format("2006-01-02"))
						break
					}
					//parse microclimate parameters for one date, columns in one row
					rowReadings := []models.MicroclimateReading{}
					for i := 1; i < len(line); i++ {
						//Parse value
						//create the microclimate reading based on the i-th position in row and i-th value
						//row format follows the DB microclimates format
						val, err := strconv.ParseFloat(strings.Trim(line[i], " "), 32)
						if err != nil {
							log.Println("Error in parsing float value in row:", err)
							doubleBreak = true
							break
						}
						//Default non existing value in API
						if val == -999.00 {
							log.Println("Error in parsing float value -999.00")
							doubleBreak = true
							break
						}
						//Append to the slice
						//Make a service method for this?
						rowReading := models.MicroclimateReading{
							MicroclimateID: uint32(i),
							LocationID:     l.ID,
							Date:           rowTime,
							Value:          float32(val),
						}
						rowReadings = append(rowReadings, rowReading)

					}
					if doubleBreak == true {
						log.Println("Stopping due to false line", line)
						break
					}
					//Create the all 6 readings in DB, slice
					//log.Println("Creating:", rowReadings)
					err = b.GetMicroclimateReadingService().CreateMicroclimateReadings(ctx, rowReadings)
					if err != nil {
						log.Println("Error in creating row readings", err)
						break
					}
					previousTime = checkTime

				}
				//Skip first 4 rows to get to the CSV row values
				counter += 1

			}
			log.Println("Finished updating location:", l.Name)
		}
		log.Println("Finished updating all locations, going to sleep now...")
		select {
		case <-ctx.Done():
			log.Println("UpdateMicroclimateReadings closing")
			waitGroup.Done()
			return

		case <-time.After(24 * time.Hour):
			log.Println("Starting updating again")
		}
	}

}

func GetBackgroundWorks() map[int]backgroundWork {
	n := map[int]backgroundWork{0: UpdateMicroclimateReadings}
	return n

}
