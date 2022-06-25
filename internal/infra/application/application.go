package application

import (
	"apsim-api/internal/infra/config"
	"apsim-api/internal/infra/db"
)

type Application interface {
	GetDB() db.DB
	GetConfig() config.Config
	GetCache() db.Cache
}

type application struct {
	DB     db.DB
	Config config.Config
	Cache  db.Cache
}

func (a *application) GetDB() db.DB {
	return a.DB
}
func (a *application) GetConfig() config.Config {
	return a.Config
}
func (a *application) GetCache() db.Cache {
	return a.Cache
}

func GetApplication() (Application, error) {
	var err error
	config := config.GetConfig()
	database, err := db.GetDB(config.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	influxDB := db.GetCache(config.GetInfluxDbUrl(), config.GetInfluxDbToken(), config.GetInfluxDbBucket(), config.GetInfluxDbOrg())
	return &application{
		DB:     database,
		Config: config,
		Cache:  influxDB,
	}, err
}
