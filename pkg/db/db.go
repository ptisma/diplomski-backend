package db

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Client *gorm.DB
}

func GetDB(connStr string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &DB{
		Client: db,
	}, nil
}

func (d *DB) Close() error {
	db, err := d.Client.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
