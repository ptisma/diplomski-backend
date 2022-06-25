package main

import (
	db2 "apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
	"encoding/json"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
)

type val struct {
	Value float32
	Dates []string
}

func main() {
	//Fetch Yields
	dbFilePath := "C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\apsim-stage-area\\apsimxFile1337869198.db"
	var yields = []models.Yield{}
	fmt.Println("Opening db file:", dbFilePath)
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	var err error
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error

	for i := 0; i < len(yields); i++ {
		query := fmt.Sprintf("SELECT strftime('%%Y-%%m-%%d', date) FROM report WHERE ROUND(yield,2)==%.2f AND yield !=0 AND strftime('%%Y', date) == '%d'", yields[i].Yield, yields[i].Year)
		err = db.Raw(query).Scan(&yields[i].Dates).Error
		if err != nil {
			//return yields, err
			fmt.Println(err)
		}
	}
	fmt.Println("yields:", yields)

	//Writo to influxDB

	bucket := "dates-test"
	org := "FER"
	token := "iCIavmWDn08-O9C5R4qjR2xrWN-57YluKNJY6HW6NCBEKPXMy_AdwwFmIi0k5TDWKRdkT6f2P4Wpe4QhVOJExQ=="
	// Store the URL of your Cache instance
	url := "http://localhost:8086"
	// Create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(url, token)
	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	var fields []db2.Field
	//var sb strings.Builder
	var tags []db2.Tag = []db2.Tag{
		{"location_id", "1"},
		{"culture_id", "1"},
		{"from", "2010"},
		{"to", "2021"},
	}
	for _, yield := range yields {
		v := val{
			Value: yield.Yield,
			Dates: yield.Dates,
		}
		z, err := json.Marshal(&v)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(z))
		fields = append(fields, db2.Field{strconv.FormatInt(int64(yield.Year), 10), string(z)})
	}

	p := influxdb2.NewPointWithMeasurement("simulacija")

	for _, tag := range tags {
		p = p.AddTag(tag.Key, tag.Value)
	}
	for _, field := range fields {
		p = p.AddField(field.Key, field.Value)
	}

	err = writeAPI.WritePoint(context.Background(), p)

	fmt.Println(err)

	yields = nil
	fmt.Println("yields:", yields)
	//Read from Cache
	queryAPI := client.QueryAPI(org)
	//Set query string
	queryStr := `from(bucket:"dates-test")
				|> range(start: 0)
				|> filter(fn: (r) => r._measurement == "simulacija" and r["location_id"] == "1" and r["culture_id"] == "1" and r["from"] == "2010" and r["to"] == "2021")
                |> sort(columns: ["_time"], desc: true)
                |> first()
                |> group()
                `
	//and r["from"]=="2016" and r["to"]=="2018"
	// Get QueryTableResult
	result, err := queryAPI.Query(context.Background(), queryStr)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			//fmt.Println(result.TablePosition())
			//fmt.Println(result.Record().Table())
			// Notice when group key has changed
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
				//fmt.Println("Nova tabela")
			}
			// Access data
			//fmt.Println(result.Record().Result())
			//fmt.Println(result.Record().Values())
			//fmt.Println(result.Record().Value())
			//fmt.Println(result.Record().Field())
			//x := result.Record().Value()
			//y, _ := x.(val)
			//fmt.Println(x)
			//fmt.Println(y)
			year := result.Record().Field()
			yearInt, err := strconv.ParseInt(year, 10, 32)
			if err != nil {
				//flag = false
				//fmt.Println("Cant parse field from influx db", err.Error())
				//break
				fmt.Println(err)
			}
			x := result.Record().Value()
			//fmt.Println(x)
			obj, _ := x.(string)
			target := val{}
			json.Unmarshal([]byte(obj), &target)
			//fmt.Println(yearInt)
			//fmt.Println(obj)
			//fmt.Println(target)

			yields = append(yields, models.Yield{
				Year:  int32(yearInt),
				Yield: target.Value,
				Dates: target.Dates,
			})
			//fmt.Println(result.Record().ValueByKey("table"))
			//fmt.Println(result.Record().ValueByKey("location_id"))
			//fmt.Println(result.Record().Value())
			//fmt.Println(result.Record().String())
			//fmt.Printf("value: %v\n", result.Record().Value())
			fmt.Println("END")
		}
		// Check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		fmt.Println(err)
		// Ensures background processes finishes
		client.Close()
	}

	fmt.Println("yields:", yields)

	// Ensures background processes finishes
	client.Close()
}
