package main

import (
	"apsim-api/refactored/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbFilePath := "C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\apsim\\apsimxFile3064893327.db"
	var yields = []models.Yield{}
	fmt.Println("Opening db file:", dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error
	fmt.Println(err)
	//for _, j := range yields {
	//	err = db.Raw(`SELECT date FROM report`).Scan(&yields).Error
	//	j.Dates =
	//
	//}
	//fmt.Println("yields:", yields)

	for i, j := range yields {
		dates := []string{}
		query := fmt.Sprintf("SELECT strftime('%%Y-%%m-%%d', date) as day FROM report WHERE ROUND(yield,2)==%.2f AND yield !=0 AND strftime('%%Y', date) == '%d'", j.Yield, j.Year)
		_ = db.Debug().Raw(query).Scan(&dates).Error
		fmt.Println(dates)
		yields[i].Dates = dates
	}

	fmt.Println(yields)

}
