package getYield

import (
	"apsim-api/internal/apsimx"
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"apsim-api/pkg/influx"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

//From API
func GetYield(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//remeber current time YYYYMMDD format
		currentDate := time.Now().Format("20060102")
		fmt.Println("Current date:", currentDate)

		//get stuff from URL
		params := mux.Vars(r)
		locationId, _ := strconv.ParseInt(params["locationId"], 10, 32)
		cultureId, _ := strconv.ParseUint(params["cultureId"], 10, 32)

		urlParams := r.URL.Query()
		fromDate, _ := time.Parse("20060102", urlParams.Get("from"))
		toDate, _ := time.Parse("20060102", urlParams.Get("to"))
		fmt.Println("LocationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

		//Position in apsim folder
		os.Chdir("./apsim")
		newDir, _ := os.Getwd()
		fmt.Println("Current working dir:", newDir)

		//Create apsimx file
		//f, _ := os.Create("test.apsimx")
		f, _ := os.CreateTemp(".", "apsimxFile*.apsimx")
		//Abs filepath of apsimx file depending on OS
		absPath := filepath.Join(newDir, f.Name())
		//if runtime.GOOS == "windows" {
		//	absPath = strings.Replace(absPath, "\\", "\\\\", -1)
		//}
		fmt.Println("Absolute path of apsimx file:", absPath)
		//On windows have to use "C:\\etc\\" because of C#

		//Optional if not using met file
		//Create a CSV file
		//csvMicro, _ := os.CreateTemp("./apsim", "csvFile*.csv")
		//Load the data from DB
		//Load in batch
		//Write into CSV

		//Create a consts file

		//Load the data from DB

		//Write into CSV

		//Get a location(latitude and longitude) based on URL ID
		location := &models.Location{ID: uint32(locationId)}
		_ = location.GetLocationById(app)
		fmt.Println("Location:", location)
		//Get a culture nedede later for apsimx file
		culture := &models.Culture{ID: uint32(cultureId)}
		fmt.Println("culture:", culture)
		//Download the met file from external API
		url := fmt.Sprintf("https://worldmodel.csiro.au/gclimate?lat=%g&lon=%g&format=apsim&start=%s&stop=%s", location.Latitude, location.Longitude, fromDate.Format("20060102"), toDate.Format("20060102"))
		//url := ""
		fmt.Println("URL:", url)
		//os.Chdir("C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\apsim")
		//out, _ := os.Create("test.met")
		out, _ := os.CreateTemp(".", "metFile*.met")
		resp, _ := http.Get(url)
		_, _ = io.Copy(out, resp.Body)
		out.Close()
		resp.Body.Close()

		//fmt.Println(filepath.Abs(filepath.Dir(out.Name())))
		rootPath, _ := filepath.Abs(filepath.Dir(out.Name()))
		metPath := filepath.Join(rootPath, out.Name())
		fmt.Println("metPath:", metPath)

		//Get json for apsimx file and write to it
		res := apsimx.NewBarleyApsimx(fromDate, toDate, metPath)
		//res := apsimx.NewBarleyApsimx(fromDate, toDate, )
		f.WriteString(res)
		f.Close()
		////fmt.Println(filepath.toSlash(absPath))
		////fmt.Println(filepath.FromSlash(absPath))

		////Start a simulation
		////dotnet apsim.dll run f /pathToApsimxFile
		//os.Chdir("../../apsim-cli/netcoreapp3.1")
		wd, _ := os.Getwd()
		fmt.Println("Current working dir:", wd)
		cmd := exec.Command("dotnet", "apsim.dll", "run", "--single-threaded", "f", absPath)
		//sad sam u apsimu
		cmd.Dir = "../../apsim-cli/netcoreapp3.1"
		_ = cmd.Run()
		//
		//Read from DB
		os.Chdir("./../../apsim-api/apsim")
		wd, _ = os.Getwd()
		fmt.Println("Current working dir:", wd)
		baseF := filepath.Base(f.Name())
		fNoExt := baseF[:len(baseF)-len(filepath.Ext(baseF))]
		dbFile := fNoExt + ".db"
		fmt.Println("Opening db file:", dbFile)
		db, _ := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		//
		var yields []struct {
			Year  int32
			Yield float32
		}
		db.Raw(`SELECT strftime('%Y', date) as year, sum(yield) as yield FROM report GROUP BY year`).Scan(&yields)
		fmt.Println("yields:", yields)

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(yields)
		w.Write(response)

		////[{1900 60376.05} {1901 58013.273} {1902 849.8151}]

	})
}

//From CSV
func GetYield2(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// remember current time im YYYYMMDD format
		currentDate := time.Now().Format("20060102")
		fmt.Println("Current date:", currentDate)

		// ...yield?from=01011990&to=01012000&locationId=1&cultureId=1
		//get stuff from URL- locationId and cultureId
		params := mux.Vars(r)
		locationId, err := strconv.ParseInt(params["locationId"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: invalid locationId")
			return
		}
		cultureId, err := strconv.ParseUint(params["cultureId"], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: invalid cultureId")
			return
		}
		//get stuff from URL- fromDate and toDate
		urlParams := r.URL.Query()
		fromDate, err := time.Parse("20060102", urlParams.Get("from"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: invalid fromDate")
			return
		}
		toDate, err := time.Parse("20060102", urlParams.Get("to"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: invalid toDate")
			return
		}
		fmt.Println("LocationId:", locationId, "cultureId:", cultureId, "fromDate:", fromDate, "toDate:", toDate)

		//Prepare response
		//var yields []struct {
		//	Year  int32
		//	Yield float32
		//}

		yields := []struct {
			Year  int32
			Yield float32
		}{}

		//Check if the result is already cached in InfluxDB before running simulation
		fluxQueryStr := fmt.Sprintf(`from(bucket:"apsim")
				|> range(start: 0)
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "%d" and r["culture_id"] == "%d" and r["from"] == "%d" and r["to"] == "%d")
                |> sort(columns: ["_time"], desc: true)
                |> first()
                |> group()
                `, locationId, cultureId, fromDate.Year(), toDate.Year())

		fmt.Println(fluxQueryStr)

		resultIterator, err := app.Writer.Read(context.TODO(), fluxQueryStr)
		flag := true
		if err == nil {
			// Iterate over query response
			for resultIterator.Next() {

				// Access data
				fmt.Printf("field: %s, value: %v\n", resultIterator.Record().Field(), resultIterator.Record().Value())
				year := resultIterator.Record().Field()
				yearInt, err := strconv.ParseInt(year, 10, 32)
				if err != nil {
					flag = false
					fmt.Println("Cant parse field from influx db", err.Error())
					break
				}
				yieldFloat, ok := resultIterator.Record().Value().(float64)
				if !ok {
					flag = false
					fmt.Println("Cant parse value from influx db")
					break
				}
				yields = append(yields, struct {
					Year  int32
					Yield float32
				}{Year: int32(yearInt), Yield: float32(yieldFloat)})
			}
			// Check for an error
			if resultIterator.Err() != nil {
				flag = false
				fmt.Printf("query parsing error: %s\n", resultIterator.Err().Error())
			}
		} else {
			fmt.Println(err)
		}

		if flag && len(yields) == (toDate.Year()-fromDate.Year()+1) {
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(yields)
			w.Write(response)
			return
		}

		// Get a location(latitude and longitude) based on URL ID and get a culture needed later for apsimx file
		location := &models.Location{ID: uint32(locationId)}
		err = location.GetLocationById(app)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: locationId doesn't exist")
			return
		}
		fmt.Println("Location:", location)
		culture := &models.Culture{ID: uint32(cultureId)}
		fmt.Println("culture:", culture)

		//Position in apsim folder
		//os.Chdir("./apsim")
		//newDir, _ := os.Getwd()
		//fmt.Println("Current working dir:", newDir)

		//Create temp apsimx file in stage area
		f, err := os.CreateTemp("./apsim", "apsimxFile*.apsimx")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching yield: temporary apsimx can't be created")
			return
		}

		////Abs filepath of apsimx file depending on OS
		//absPath := filepath.Join(newDir, f.Name())

		absPathApsimx, _ := filepath.Abs(f.Name())
		fmt.Println("Absolute path of apsimx file:", absPathApsimx)
		////On windows have to use "C:\\etc\\" because of C#
		//
		////Create a CSV file
		f2, _ := os.CreateTemp("./apsim", "csv*.csv")
		absPathCsv, _ := filepath.Abs(f2.Name())
		fmt.Println("Absolute path of csv file:", absPathCsv)
		//fmt.Println(filepath.Base(f2.Name()))
		//fmt.Println(filepath.Join(newDir, f2.Name()))
		csvFile, _ := os.OpenFile(absPathCsv, os.O_APPEND|os.O_WRONLY, 0644)
		//csvFile.WriteString("kek")

		//Load the data from DB
		//Load in batch
		//Write into CSV
		//Write first line
		csvFile.WriteString("year,day,radn,maxt,mint,rain,pan,vp,code\n")
		results := &[]models.MicroclimateReading{}
		buff := []models.MicroclimateReading{}
		counter := 0
		//fix >= doesnt work on start so less than 1 day, fix format myb?
		result := app.DB.Client.Debug().Where("location_id = ? AND date >= ? AND date <= ?", location.ID, fromDate.AddDate(0, 0, -1), toDate).FindInBatches(results, 102, func(tx *gorm.DB, batch int) error {
			for _, result := range *results {
				// batch processing found records
				//fmt.Println(result)
				counter += 1
				buff = append(buff, result)
				if counter == 6 {
					rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
					csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
					csvFile.WriteString(csvRow)
					counter = 0
					buff = nil
				}

			}
			//tx.Save(&results)

			//fmt.Println(tx.RowsAffected) // number of records in this batch

			//fmt.Println(batch) // Batch 1, 2, 3

			// returns error will stop future batches
			return nil
		})
		fmt.Println(result.Error)
		//
		//if toDate is greater than max from db fetch from predicted microclimateReading table
		microclimateReading := models.MicroclimateReading{LocationID: uint32(locationId)}
		_ = microclimateReading.GetLatestMicroclimateReading(app)
		fmt.Println(microclimateReading)
		lastTime, _ := time.Parse("2006-01-02", microclimateReading.Date)
		if toDate.After(lastTime) {
			results := &[]models.PredictedMicroclimateReading{}
			buff := []models.PredictedMicroclimateReading{}
			counter := 0
			result := app.DB.Client.Debug().Where("location_id = ? AND date >= ? AND date <= ?", location.ID, fromDate.AddDate(0, 0, -1), toDate).FindInBatches(results, 102, func(tx *gorm.DB, batch int) error {
				for _, result := range *results {
					// batch processing found records
					//fmt.Println(result)
					counter += 1
					buff = append(buff, result)
					if counter == 6 {
						rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
						csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
						csvFile.WriteString(csvRow)
						counter = 0
						buff = nil
					}

				}
				//tx.Save(&results)

				//fmt.Println(tx.RowsAffected) // number of records in this batch

				//fmt.Println(batch) // Batch 1, 2, 3

				// returns error will stop future batches
				return nil
			})
			fmt.Println(result.Error)

		}

		//Create a consts file
		f3, _ := os.CreateTemp("./apsim", "const*.txt")
		absPathConst, _ := filepath.Abs(f3.Name())
		//fmt.Println(filepath.Base(f3.Name()))
		//fmt.Println(filepath.Join(newDir, f3.Name()))
		constFile, _ := os.OpenFile(absPathConst, os.O_APPEND|os.O_WRONLY, 0644)

		//Write into consts file
		constFile.WriteString(fmt.Sprintf("location = %s\n", location.Name))
		constFile.WriteString(fmt.Sprintf("latitude = %.2f (DECIMAL DEGREES)\n", location.Latitude))
		constFile.WriteString(fmt.Sprintf("longitude = %.2f (DECIMAL DEGREES)\n", location.Longitude))

		//Get json for apsimx file and write to it
		res := apsimx.NewBarleyApsimx2(fromDate, toDate, absPathCsv, absPathConst)

		f.WriteString(res)
		f.Close()

		f2.Close()
		f3.Close()
		constFile.Close()
		csvFile.Close()

		//Start a simulation
		//dotnet apsim.dll run f /pathToApsimxFile
		//os.Chdir("../../apsim-cli/netcoreapp3.1")
		//wd, _ := os.Getwd()
		//fmt.Println("Current working dir:", wd)
		cmd := exec.Command("dotnet", "apsim.dll", "run", "--single-threaded", "f", absPathApsimx)
		//sad sam u apsimu
		cmd.Dir = "../apsim-cli/netcoreapp3.1"
		_ = cmd.Run()
		//
		//Read from DB
		//os.Chdir("./../../apsim-api/apsim")
		//wd, _ = os.Getwd()
		//fmt.Println("Current working dir:", wd)
		baseF := filepath.Base(f.Name())
		fNoExt := baseF[:len(baseF)-len(filepath.Ext(baseF))]
		dbFile := filepath.Join(filepath.Dir(absPathApsimx), fNoExt+".db")
		fmt.Println("Opening db file:", dbFile)
		db, _ := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		////

		db.Raw(`SELECT strftime('%Y', date) as year, sum(yield) as yield FROM report GROUP BY year`).Scan(&yields)
		fmt.Println("yields:", yields)
		dbC, _ := db.DB()
		dbC.Close()
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(yields)
		w.Write(response)

		//[{1900 60376.05} {1901 58013.273} {1902 849.8151}]

		//Write to influx db
		var fields []influx.Field
		for _, yield := range yields {
			fields = append(fields, influx.Field{strconv.FormatInt(int64(yield.Year), 10), yield.Yield})
		}
		err = app.Writer.Write(
			context.TODO(),
			"simulation",
			[]influx.Tag{
				{"location_id", strconv.FormatInt(locationId, 10)},
				{"culture_id", strconv.FormatUint(cultureId, 10)},
				{"from", strconv.FormatInt(int64(fromDate.Year()), 10)},
				{"to", strconv.FormatInt(int64(toDate.Year()), 10)},
			},
			fields)
		fmt.Println(err)

		//Clear stage area
		os.Remove(absPathApsimx)
		os.Remove(absPathConst)
		os.Remove(absPathCsv)
		os.Remove(dbFile)
	})
}
