package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	org := "FER"
	token := "iCIavmWDn08-O9C5R4qjR2xrWN-57YluKNJY6HW6NCBEKPXMy_AdwwFmIi0k5TDWKRdkT6f2P4Wpe4QhVOJExQ=="
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	// Create client
	client := influxdb2.NewClient(url, token)
	// Get query client
	queryAPI := client.QueryAPI(org)
	//Set query string
	queryStr := `from(bucket:"apsim-stage-area")
				|> range(start: 0)
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "1" and r["culture_id"] == "1" and r["from"] == "2000" and r["to"] == "2002")
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
			x := result.Record().Value()
			y, _ := x.(float64)
			fmt.Println(y)

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
}

/*
from(bucket:"apsim-stage-area")
				|> range(start: 0)
                |> last()
				|> group(columns: ["from", "to"])
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "1" and r["culture_id"] == "1")


from(bucket:"apsim-stage-area")
				|> range(start: -500h)
				|> group(columns: ["from", "to"])
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "1" and r["culture_id"] == "1")

from(bucket:"apsim-stage-area")
				|> range(start: 0)
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "1" and r["culture_id"] == "1")
                |> group()
                |> sort(columns:["_field"], desc: false)
                |> sort(columns:["_time"], desc: true)
                |> group(columns: ["_time", "from", "to", ])
                |> group()

//works
from(bucket:"apsim-stage-area")
				|> range(start: 0)
				|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "1" and r["culture_id"] == "1" and r["from"] == "2000" and r["to"] == "2002")
                |> sort(columns: ["_time"], desc: true)
                |> first()
                |> group()
*/
