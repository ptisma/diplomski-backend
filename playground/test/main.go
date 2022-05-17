package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Car struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func main() {
	db, err := gorm.Open(sqlite.Open("baza.db"), &gorm.Config{}) //baza.db se odnosi na working dir
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)
	car := Car{Name: "chevy"}
	result := db.Create(&car)
	fmt.Println(result)
	//db.Commit()

	wd, err := os.Getwd()
	fmt.Println(wd)
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	fmt.Println(parent)
}
