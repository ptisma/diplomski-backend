package repositories

import (
	"apsim-api/refactored/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PredictedMicroclimateReadingRepository struct {
	DB *gorm.DB
}

func (r *PredictedMicroclimateReadingRepository) GetMicroClimateReadings(microclimateID, locationID int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error) {

	var err error
	microclimates := []models.PredictedMicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"

	err = r.DB.Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Group("microclimate_id,location_id,date").Order("date").Find(&microclimates).Error
	if err != nil {
		return []models.PredictedMicroclimateReading{}, err
	}
	return microclimates, err

}

func (r *PredictedMicroclimateReadingRepository) GetLatestMicroClimateReading(locationID int) (models.PredictedMicroclimateReading, error) {

	var err error
	microclimateReading := models.PredictedMicroclimateReading{}
	queryStr := "location_id = ?"

	err = r.DB.Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID).Order("date desc").First(&microclimateReading).Error

	return microclimateReading, err
}

func (r *PredictedMicroclimateReadingRepository) GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.PredictedMicroclimateReading, ctxx context.Context) error {

	results := []models.PredictedMicroclimateReading{}

	result := r.DB.Debug().WithContext(ctxx).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			//cancel()
			//fmt.Println("result:", result)
			select {
			case ch <- result:
				//fmt.Println("Poslao")
			case <-ctxx.Done():
				fmt.Println("ctx batch microclimate", ctxx.Err())
				//close(ch)
				return ctxx.Err()
			}
		}

		return nil
	}).Error

	close(ch)

	return result
}
