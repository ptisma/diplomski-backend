package config

import (
	"flag"
	"os"
)

type Config struct {
	dbUser         string
	dbPswd         string
	dbHost         string
	dbPort         string
	dbName         string
	apiPort        string
	influxDbUrl    string
	influxDbToken  string
	influxDbBucket string
	influxDbOrg    string
}

func GetConfig() *Config {
	conf := &Config{}

	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("POSTGRES_PASSWORD"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("POSTGRES_DB"), "DB name")
	flag.StringVar(&conf.apiPort, "apiport", os.Getenv("API_PORT"), "API port")
	flag.StringVar(&conf.influxDbUrl, "influxdburl", os.Getenv("INFLUX_DB_URL"), "Influx DB url")
	flag.StringVar(&conf.influxDbToken, "influxdtoken", os.Getenv("INFLUX_DB_TOKEN"), "Influx DB token")
	flag.StringVar(&conf.influxDbBucket, "influxdbbucket", os.Getenv("INFLUX_DB_BUCKET"), "Influx DB bucket")
	flag.StringVar(&conf.influxDbOrg, "influxdborg", os.Getenv("INFLUX_DB_ORG"), "Influx DB org")

	return conf

}

func (cfg *Config) GetDBConnectionString() string {
	//return fmt.Sprintf("")

	//for now hardcoded
	return "baza.db"
}
func (cfg *Config) GetApiPort() string {
	return cfg.apiPort
}
func (cfg *Config) GetInfluxDbUrl() string {
	return cfg.influxDbUrl
}
func (cfg *Config) GetInfluxDbToken() string {
	return cfg.influxDbToken
}
func (cfg *Config) GetInfluxDbBucket() string {
	return cfg.influxDbBucket
}
func (cfg *Config) GetInfluxDbOrg() string {
	return cfg.influxDbOrg
}

