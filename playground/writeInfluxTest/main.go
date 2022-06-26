package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	//Petar
	//Sifra123
	bucket := "demo"
	org := "FER"
	token := "iCIavmWDn08-O9C5R4qjR2xrWN-57YluKNJY6HW6NCBEKPXMy_AdwwFmIi0k5TDWKRdkT6f2P4Wpe4QhVOJExQ=="
	// Store the URL of your InfluxDB instance
	url := "http://localhost:8086"
	// Create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(url, token)
	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	//// Create point using full params constructor
	//p := influxdb2.NewPoint("simulation",
	//	map[string]string{"2014": "350", "2015": "213"},
	//	map[string]interface{}{},
	//	time.Now())

	p := influxdb2.NewPointWithMeasurement("simulation")
	p = p.AddTag("location_id", "1")
	p = p.AddTag("culture_id", "1")
	p = p.AddTag("from", "2000")
	p = p.AddTag("to", "2002")
	p = p.AddField("2000", 31.27)
	p = p.AddField("2001", 45.67)
	p = p.AddField("2002", 87.12)
	// Write point immediately
	err := writeAPI.WritePoint(context.Background(), p)
	fmt.Println(err)
	// Ensures background processes finishes
	client.Close()
}
