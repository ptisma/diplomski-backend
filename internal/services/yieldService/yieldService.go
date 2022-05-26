package yieldService

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"time"
)

type YieldService struct {
	app *application.Application
}

func GetYieldService(app *application.Application) *YieldService {
	return &YieldService{
		app: app,
	}
}

func (ys *YieldService) CheckCachedYield(locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error) {

	yield := models.Yield{
		LocationId: locationId,
		CultureId:  cultureId,
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	yields, err := yield.GetYields(ys.app)
	//fmt.Println("yields:", yields)
	//fmt.Println("err:", err)
	//fmt.Println(yields == nil)
	//fmt.Println(len(yields), (toDate.Year() - fromDate.Year() + 1))
	if len(yields) != (toDate.Year() - fromDate.Year() + 1) {
		return nil, err
	}
	return yields, err

}

func (ys *YieldService) CacheYield(locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error {

	yield := models.Yield{
		LocationId: locationId,
		CultureId:  cultureId,
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	err := yield.CreateYields(ys.app, yields)
	//fmt.Println("CacheYield error", err)
	return err

}
