package main

import (
	"apsim-api/internal/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// getJSON fetches the contents of the given URL
// and decodes it as JSON into the given result,
// which should be a pointer to the expected data.
func getJSON(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}
	// We could check the resulting content type
	// here if desired.
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %v", err)
	}
	return nil
}

func loadMicroCLimateReading(microclimateId, locationId uint32, date string, value float32) models.MicroclimateReading {

	return models.MicroclimateReading{
		MicroclimateID: microclimateId,
		LocationID:     locationId,
		Date:           date,
		Value:          value,
	}
}

var MICROCLIMATE_ID uint32 = 1
var LOCATION_ID uint32 = 1

func main() {
	db, _ := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{}) //baza.db se odnosi na working dir
	var buf bytes.Buffer
	url := "https://worldmodel.csiro.au/gclimate?lat=45.815080&lon=15.9819189&format=csv&start=20220414&stop=20220415"

	resp, _ := http.Get(url)

	io.Copy(&buf, resp.Body)

	scanner := bufio.NewScanner(&buf)
	counter := 0
	for scanner.Scan() {
		if counter > 3 {
			fmt.Println("Linija:", scanner.Text())
			line := strings.Split(scanner.Text(), ",")
			fmt.Println(line)
			for i := 1; i < len(line); i++ {
				time2 := line[0]
				val, err := strconv.ParseFloat(line[1], 32)
				if err != nil {
					continue
				}
				res := loadMicroCLimateReading(MICROCLIMATE_ID, LOCATION_ID, time2, float32(val))
				fmt.Println(res.MicroclimateID, res.LocationID, res.Date, res.Value)
				db.Create(&res)

			}
		}
		counter += 1
	}

}
