package influx

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Tag struct {
	Key   string
	Value string
}

type Field struct {
	Key   string
	Value interface{}
}
type InfluxWriter struct {
	Client influxdb2.Client
	Bucket string
	Org    string
}

func (i *InfluxWriter) Write(ctx context.Context, measurement string, tags []Tag, fields []Field) error {
	var err error
	p := influxdb2.NewPointWithMeasurement(measurement)
	for _, tag := range tags {
		p = p.AddTag(tag.Key, tag.Value)
	}
	for _, field := range fields {
		p = p.AddField(field.Key, field.Value)
	}
	writeAPI := i.Client.WriteAPIBlocking(i.Org, i.Bucket)
	err = writeAPI.WritePoint(ctx, p)
	return err
}

func (i *InfluxWriter) Read(ctx context.Context, queryStr string) (*api.QueryTableResult, error) {
	queryAPI := i.Client.QueryAPI(i.Org)
	return queryAPI.Query(ctx, queryStr)
}

func (i *InfluxWriter) Close() error {
	i.Client.Close()
	return nil
}
func GetInfluxWriter(url, token, bucket, org string) *InfluxWriter {
	client := influxdb2.NewClient(url, token)

	return &InfluxWriter{Client: client, Bucket: bucket, Org: org}

}
