package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type YieldService struct {
	I interfaces.IYieldRepository
}

//Retreive Yields from cache in InfluxDB
func (s *YieldService) GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error) {

	return s.I.GetYields(ctx, locationId, cultureId, fromDate, toDate)
}

//Create Yields for cache in InfluxDB
func (s *YieldService) CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error {

	return s.I.CreateYields(ctx, locationId, cultureId, fromDate, toDate, yields)
}

//Read Yields from the created .db file after simulation run
func (s *YieldService) ReadYields(dbFilePath string) ([]models.Yield, error) {
	//Move to repository or keep here
	var yields = []models.Yield{}
	//fmt.Println("Opening db file:", dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Raw(`SELECT strftime('%Y', date) as year, max(yield) as yield FROM report GROUP BY year`).Scan(&yields).Error

	//for i := 0; i < len(yields); i++ {
	//	query := fmt.Sprintf("SELECT strftime('%%Y-%%m-%%d', date) FROM report WHERE ROUND(yield,2)==%.2f AND yield !=0 AND strftime('%%Y', date) == '%d'", yields[i].Yield, yields[i].Year)
	//	err = db.Raw(query).Scan(&yields[i].Dates).Error
	//	if err != nil {
	//		//fmt.Println(err)
	//		return yields, err
	//
	//	}
	//}
	//fmt.Println("yields:", yields)
	client, _ := db.DB()
	//not handled
	client.Close()
	return yields, err

}

//check error log messages if there have been any errors in simulation
func (s *YieldService) ReadYieldErrors(dbFilePath string) ([]string, error) {
	var err error
	errLogs := []string{}
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//MessageType 0 for errors
	err = db.Raw(`SELECT Message FROM _Messages WHERE MessageType == 0`).Scan(&errLogs).Error
	if err != nil {
		return nil, err
	}
	return errLogs, err

}

//helper function to validate yield readings if they fit the interval
func (s *YieldService) ValidateYields(fromDate, toDate time.Time, yields []models.Yield) bool {
	if len(yields) != (toDate.Year() - fromDate.Year() + 1) {
		return false
	}
	return true
}
