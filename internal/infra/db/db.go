package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB interface {
	Close() error
	GetClient() *gorm.DB
}

type database struct {
	Client *gorm.DB
}

func (d *database) GetClient() *gorm.DB {
	return d.Client
}

func (d *database) Close() error {
	db, err := d.Client.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func GetDB(connStr string) (DB, error) {
	var err error
	db, err := gorm.Open(sqlite.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &database{
		Client: db,
	}, err
}
