package db

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Tag struct {
	Key   string
	Value string
}

type Field struct {
	Key   string
	Value interface{}
}

type Cache interface {
	//Write(ctx context.Context, measurement string, tags []Tag, fields []Field) error
	//Read(ctx context.Context, queryStr string) (*api.QueryTableResult, error)
	GetClient() influxdb2.Client
	Close() error
}

type influxDB struct {
	Client influxdb2.Client
	Bucket string
	Org    string
}

func (i *influxDB) GetClient() influxdb2.Client {
	return i.Client
}

//func (i *influxDB) Write(ctx context.Context, measurement string, tags []Tag, fields []Field) error {
//	var err error
//	p := influxdb2.NewPointWithMeasurement(measurement)
//	for _, tag := range tags {
//		p = p.AddTag(tag.Key, tag.Value)
//	}
//	for _, field := range fields {
//		p = p.AddField(field.Key, field.Value)
//	}
//	writeAPI := i.Client.WriteAPIBlocking(i.Org, i.Bucket)
//	err = writeAPI.WritePoint(ctx, p)
//	return err
//}
//
//func (i *influxDB) Read(ctx context.Context, queryStr string) (*api.QueryTableResult, error) {
//	queryAPI := i.Client.QueryAPI(i.Org)
//	return queryAPI.Query(ctx, queryStr)
//}

func (i *influxDB) Close() error {
	i.Client.Close()
	return nil
}

func GetCache(url, token, bucket, org string) Cache {
	//fmt.Println("INFLUX", url, token, bucket, org)
	client := influxdb2.NewClient(url, token)
	return &influxDB{Client: client, Bucket: bucket, Org: org}

}
