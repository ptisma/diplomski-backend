package config

import (
	"flag"
	"os"
)

type Config interface {
	GetDBConnectionString() string
	GetApiURL() string
	GetApiPort() string
	GetInfluxDbUrl() string
	GetInfluxDbToken() string
	GetInfluxDbBucket() string
	GetInfluxDbOrg() string
	GetInfluxDbMeasurement() string
}

type cfg struct {
	//dbUser         string
	//dbPswd         string
	//dbHost         string
	//dbPort         string
	//dbName         string
	apiURL              string
	apiPort             string
	influxDbUrl         string
	influxDbToken       string
	influxDbBucket      string
	influxDbOrg         string
	influxDbMeasurement string
}

func (c *cfg) GetDBConnectionString() string {
	//for now hardcoded
	return "baza.db"
}
func (c *cfg) GetApiURL() string {
	return c.apiURL
}
func (c *cfg) GetApiPort() string {
	return c.apiPort
}
func (c *cfg) GetInfluxDbUrl() string {
	return c.influxDbUrl
}
func (c *cfg) GetInfluxDbToken() string {
	return c.influxDbToken
}
func (c *cfg) GetInfluxDbBucket() string {
	return c.influxDbBucket
}
func (c *cfg) GetInfluxDbOrg() string {
	return c.influxDbOrg
}
func (c *cfg) GetInfluxDbMeasurement() string {
	return c.influxDbMeasurement
}

func GetConfig() Config {
	conf := &cfg{}
	//flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "DB user name")
	//flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("POSTGRES_PASSWORD"), "DB pass")
	//flag.StringVar(&conf.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "DB port")
	//flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "DB host")
	//flag.StringVar(&conf.dbName, "dbname", os.Getenv("POSTGRES_DB"), "DB name")
	flag.StringVar(&conf.apiURL, "apiURL", os.Getenv("API_URL"), "API URL")
	flag.StringVar(&conf.apiPort, "apiport", os.Getenv("API_PORT"), "API port")
	flag.StringVar(&conf.influxDbUrl, "influxdburl", os.Getenv("INFLUX_DB_URL"), "Influx DB url")
	flag.StringVar(&conf.influxDbToken, "influxdtoken", os.Getenv("INFLUX_DB_TOKEN"), "Influx DB token")
	flag.StringVar(&conf.influxDbBucket, "influxdbbucket", os.Getenv("INFLUX_DB_BUCKET"), "Influx DB bucket")
	flag.StringVar(&conf.influxDbOrg, "influxdborg", os.Getenv("INFLUX_DB_ORG"), "Influx DB org")
	flag.StringVar(&conf.influxDbMeasurement, "influxdbmeasurement", os.Getenv("INFLUX_DB_MEASUREMENT"), "Influx DB Measurement")
	return conf

}
