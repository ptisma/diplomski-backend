package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
	"time"
)

//same workflow as microclimateRepository
type PredictedMicroclimateReadingRepository struct {
	//DB *gorm.DB
	DB db.DB
}

func (r *PredictedMicroclimateReadingRepository) GetMicroClimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error) {

	var err error
	microclimates := []models.PredictedMicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"

	//err = r.DB.WithContext(ctx).Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Group("microclimate_id,location_id,date").Order("date").Find(&microclimates).Error
	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Order("date").Find(&microclimates).Error
	if err != nil {
		return nil, err
	}
	return microclimates, err

}

func (r *PredictedMicroclimateReadingRepository) GetLatestMicroClimateReadingByID(ctx context.Context, microclimateID, locationID int) (models.PredictedMicroclimateReading, error) {

	var err error
	microclimateReading := models.PredictedMicroclimateReading{}
	queryStr := "location_id = ? AND microclimate_ID = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID, microclimateID).Order("date desc").First(&microclimateReading).Error
	if err != nil {
		return models.PredictedMicroclimateReading{}, err
	}
	return microclimateReading, err
}

func (r *PredictedMicroclimateReadingRepository) GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.PredictedMicroclimateReading, error) {

	var err error
	microclimateReading := models.PredictedMicroclimateReading{}
	queryStr := "location_id = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.PredictedMicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID).Order("date desc").First(&microclimateReading).Error
	if err != nil {
		return models.PredictedMicroclimateReading{}, err
	}
	return microclimateReading, err
}

func (r *PredictedMicroclimateReadingRepository) GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.PredictedMicroclimateReading, ctxx context.Context) error {
	defer close(ch)
	results := []models.PredictedMicroclimateReading{}
	//result := r.DB.Debug().WithContext(ctxx).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Order("date, microclimate_id").FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
	//	for _, result := range results {
	//		//cancel()
	//		//fmt.Println("result:", result)
	//		select {
	//		case ch <- result:
	//			//fmt.Println("Poslao")
	//		case <-ctxx.Done():
	//			fmt.Println("ctx batch predicted microclimate reading", ctxx.Err())
	//			//close(ch)
	//			return ctxx.Err()
	//		}
	//	}
	//
	//	return nil
	//}).Error
	//
	//close(ch)
	//
	//return result
	var err error
	batchSize := 100
	var (
		rowsAffected int64
		batch        int
	)
	for {
		results = nil
		result := r.DB.GetClient().Limit(batchSize).Offset(batch*batchSize).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Order("date, microclimate_id").Find(&results)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			//do something
			for _, result := range results {
				//fmt.Println(result)
				select {
				case ch <- result:
					//fmt.Println("Poslao")
					//fmt.Println(result)
				case <-ctxx.Done():
					//fmt.Println("ctx batch microclimate", ctxx.Err())
					//close(ch)
					return ctxx.Err()

				}
			}

			//fmt.Println("GOTOV BATCH")
			if int(result.RowsAffected) < batchSize {
				break
			}

		} else if result.Error != nil {
			err = result.Error
			break
		}

	}

	return err
}
