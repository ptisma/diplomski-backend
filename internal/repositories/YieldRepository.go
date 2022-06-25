package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"strconv"
	"time"
)

type YieldRepository struct {
	//Client      influxdb2.Client
	DB          db.Cache
	Org         string
	Bucket      string
	Measurement string
}

//type val struct {
//	value float32
//	dates []string
//}

//Fetch all yields based on starting, ending date, culture id and location id
func (r *YieldRepository) GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error) {
	var yields = []models.Yield{}
	var err error

	//flux query
	//filter based on arguments returns tables with records (every pair of field and value with their tags in each new table)
	//sort works on each table, get the latest one with first
	//use group to flat out to single table
	fluxQueryStr := fmt.Sprintf(`from(bucket:"%s")
			|> range(start: 0)
			|> filter(fn: (r) => r._measurement == "%s" and r["location_id"] == "%d" and r["culture_id"] == "%d" and r["from"] == "%d" and r["to"] == "%d")
	       |> sort(columns: ["_time"], desc: true)
	       |> first()
	       |> group()
	       `, r.Bucket, r.Measurement, locationId, cultureId, fromDate.Year(), toDate.Year())
	//fmt.Println(fluxQueryStr)

	queryAPI := r.DB.GetClient().QueryAPI(r.Org)
	resultIterator, err := queryAPI.Query(ctx, fluxQueryStr)
	//flag := true
	if err == nil {
		// Iterate over query response
		for resultIterator.Next() {

			// Access data
			//fmt.Printf("field: %s, value: %v\n", resultIterator.Record().Field(), resultIterator.Record().Value())
			year := resultIterator.Record().Field()
			yearInt, err := strconv.ParseInt(year, 10, 32)
			if err != nil {
				//flag = false
				//fmt.Println("Cant parse field from influx db", err.Error())
				//break
				return nil, err
			}
			//its interface{} so we have to type insert
			yieldFloat, ok := resultIterator.Record().Value().(float64)
			if !ok {
				//flag = false
				//fmt.Println("Cant parse value from influx db")
				//break
				return nil, err
			}
			//x := result.Record().Value()
			////fmt.Println(x)
			//obj, _ := x.(string)
			//target := val{}
			//json.Unmarshal([]byte(obj), &target)
			//yields = append(yields, models.Yield{
			//	Year:  int32(yearInt),
			//	Yield: target.Value,
			//	Dates: target.Dates,
			//})
			yields = append(yields, models.Yield{Year: int32(yearInt), Yield: float32(yieldFloat)})
		}
		// Check for an error
		if resultIterator.Err() != nil {
			//flag = false
			//fmt.Printf("query parsing error: %s\n", resultIterator.Err().Error())
			return nil, resultIterator.Err()
		}
	}

	return yields, err
}

//Create yields based on starting, ending date, culture id and location id
func (r *YieldRepository) CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error {
	var err error
	var fields []db.Field
	//var sb strings.Builder
	//Create tags
	var tags = []db.Tag{
		{"location_id", strconv.FormatInt(int64(locationId), 10)},
		{"culture_id", strconv.FormatInt(int64(cultureId), 10)},
		{"from", strconv.FormatInt(int64(fromDate.Year()), 10)},
		{"to", strconv.FormatInt(int64(toDate.Year()), 10)},
	}
	//Create pairs of field and value
	for _, yield := range yields {
		//v := val{
		//	value: yield.Yield,
		//	dates: yield.Dates,
		//}
		//z, _ := json.Marshal(v)
		fields = append(fields, db.Field{strconv.FormatInt(int64(yield.Year), 10), yield.Yield})
	}

	//Point is a timestamped data
	p := influxdb2.NewPointWithMeasurement(r.Measurement)
	//add tags and pairs of field and value
	for _, tag := range tags {
		p = p.AddTag(tag.Key, tag.Value)
	}
	for _, field := range fields {
		p = p.AddField(field.Key, field.Value)
	}

	//store the simulation results as point
	writeAPI := r.DB.GetClient().WriteAPIBlocking(r.Org, r.Bucket)
	err = writeAPI.WritePoint(ctx, p)

	return err
}
