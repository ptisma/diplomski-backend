package application

import (
	"apsim-api/pkg/config"
	"apsim-api/pkg/db"
	"apsim-api/pkg/influx"
)

type Application struct {
	DB     *db.DB
	Config *config.Config
	Writer *influx.InfluxWriter
}

func GetApplication() (*Application, error) {

	config := config.GetConfig()
	database, err := db.GetDB(config.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	influx := influx.GetInfluxWriter(config.GetInfluxDbUrl(), config.GetInfluxDbToken(), config.GetInfluxDbBucket(), config.GetInfluxDbOrg())
	return &Application{
		DB:     database,
		Config: config,
		Writer: influx,
	}, nil
}
