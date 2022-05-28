package repositories

import (
	"apsim-api/refactored/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MicroclimateReadingRepository struct {
	DB *gorm.DB
}

func (r *MicroclimateReadingRepository) CreateMicroclimateReading(ctx context.Context, microclimateID, locationID int, date string, value float32) error {
	var err error
	microclimateReading := models.MicroclimateReading{
		ID:             0,
		MicroclimateID: uint32(microclimateID),
		LocationID:     uint32(locationID),
		Date:           date,
		Value:          value,
	}
	err = r.DB.WithContext(ctx).Create(&microclimateReading).Error
	return err
}
func (r *MicroclimateReadingRepository) GetMicroClimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error) {

	var err error
	microclimates := []models.MicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"

	err = r.DB.WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Group("microclimate_id,location_id,date").Order("date").Find(&microclimates).Error
	if err != nil {
		return []models.MicroclimateReading{}, err
	}
	return microclimates, err

}

func (r *MicroclimateReadingRepository) GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error) {

	var err error
	microclimateReading := models.MicroclimateReading{}
	queryStr := "location_id = ?"

	err = r.DB.WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID).Order("date desc").First(&microclimateReading).Error

	return microclimateReading, err
}

func (r *MicroclimateReadingRepository) GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error {

	results := []models.MicroclimateReading{}

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
