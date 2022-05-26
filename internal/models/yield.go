package models

import (
	"apsim-api/pkg/application"
	"apsim-api/pkg/influx"
	"context"
	"fmt"
	"strconv"
	"time"
)

type Yield struct {
	Year  int32   `json:"year"`
	Yield float32 `json:"yield"`

	FromDate   time.Time `json:"-" `
	ToDate     time.Time `json:"-"`
	LocationId int       `json:"-"`
	CultureId  int       `json:"-"`
}

func (y *Yield) GetYields(app *application.Application) ([]Yield, error) {
	var yields = []Yield{}
	var err error
	fluxQueryStr := fmt.Sprintf(`from(bucket:"apsim")
			|> range(start: 0)
			|> filter(fn: (r) => r._measurement == "simulation" and r["location_id"] == "%d" and r["culture_id"] == "%d" and r["from"] == "%d" and r["to"] == "%d")
	       |> sort(columns: ["_time"], desc: true)
	       |> first()
	       |> group()
	       `, y.LocationId, y.CultureId, y.FromDate.Year(), y.ToDate.Year())
	//fmt.Println(fluxQueryStr)

	resultIterator, err := app.Writer.Read(context.TODO(), fluxQueryStr)
	//flag := true
	if err == nil {
		// Iterate over query response
		for resultIterator.Next() {

			// Access data
			fmt.Printf("field: %s, value: %v\n", resultIterator.Record().Field(), resultIterator.Record().Value())
			year := resultIterator.Record().Field()
			yearInt, err := strconv.ParseInt(year, 10, 32)
			if err != nil {
				//flag = false
				fmt.Println("Cant parse field from influx db", err.Error())
				//break
				return nil, err
			}
			yieldFloat, ok := resultIterator.Record().Value().(float64)
			if !ok {
				//flag = false
				fmt.Println("Cant parse value from influx db")
				//break
				return nil, err
			}
			yields = append(yields, Yield{Year: int32(yearInt), Yield: float32(yieldFloat)})
		}
		// Check for an error
		if resultIterator.Err() != nil {
			//flag = false
			fmt.Printf("query parsing error: %s\n", resultIterator.Err().Error())
			return nil, resultIterator.Err()
		}
	} else {
		fmt.Println(err)
	}

	return yields, err
}

func (y *Yield) CreateYields(app *application.Application, yields []Yield) error {

	var fields []influx.Field
	for _, yield := range yields {
		fields = append(fields, influx.Field{strconv.FormatInt(int64(yield.Year), 10), yield.Yield})
	}
	err := app.Writer.Write(
		context.TODO(),
		"simulation",
		[]influx.Tag{
			{"location_id", strconv.FormatInt(int64(y.LocationId), 10)},
			{"culture_id", strconv.FormatInt(int64(y.CultureId), 10)},
			{"from", strconv.FormatInt(int64(y.FromDate.Year()), 10)},
			{"to", strconv.FormatInt(int64(y.ToDate.Year()), 10)},
		},
		fields)

	return err
}
