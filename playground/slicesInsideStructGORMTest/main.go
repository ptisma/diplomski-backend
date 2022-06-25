package main

import (
	"apsim-api/internal/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	//dbFilePath := "C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-stage-area-api\\apsim-stage-area\\apsimxFile3064893327.db"
	//var yields = []models.Yield{}
	//fmt.Println("Opening db file:", dbFilePath)
	//db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	//err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error
	//fmt.Println(err)
	////for _, j := range yields {
	////	err = db.Raw(`SELECT date FROM report`).Scan(&yields).Error
	////	j.Dates =
	////
	////}
	////fmt.Println("yields:", yields)
	//
	//for i, j := range yields {
	//	dates := []string{}
	//	query := fmt.Sprintf("SELECT strftime('%%Y-%%m-%%d', date) as day FROM report WHERE ROUND(yield,2)==%.2f AND yield !=0 AND strftime('%%Y', date) == '%d'", j.Yield, j.Year)
	//	_ = db.Debug().Raw(query).Scan(&dates).Error
	//	fmt.Println(dates)
	//	yields[i].Dates = dates
	//}
	//
	//fmt.Println(yields)

	dbFilePath := "C:\\Users\\gulas\\Desktop\\faks\\peta\\diplomski\\backend\\apsim-api\\apsim-stage-area\\apsimxFile1337869198.db"
	var yields = []models.Yield{}
	fmt.Println("Opening db file:", dbFilePath)
	db, _ := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	var err error
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error

	for i := 0; i < len(yields); i++ {
		query := fmt.Sprintf("SELECT strftime('%%Y-%%m-%%d', date) FROM report WHERE ROUND(yield,2)==%.2f AND yield !=0 AND strftime('%%Y', date) == '%d'", yields[i].Yield, yields[i].Year)
		err = db.Raw(query).Scan(&yields[i].Dates).Error
		if err != nil {
			//return yields, err
			fmt.Println(err)
		}
	}
	fmt.Println("yields:", yields)

}
