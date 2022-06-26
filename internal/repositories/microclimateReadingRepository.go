package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
	"time"
)

type MicroclimateReadingRepository struct {
	//DB *gorm.DB
	DB db.DB
}

//Creates new Microclimate reading based on microclimate parameter and location id
func (r *MicroclimateReadingRepository) CreateMicroclimateReading(ctx context.Context, microclimateID, locationID int, date string, value float32) error {
	var err error
	microclimateReading := models.MicroclimateReading{
		ID:             0,
		MicroclimateID: uint32(microclimateID),
		LocationID:     uint32(locationID),
		Date:           date,
		Value:          value,
	}
	err = r.DB.GetClient().Debug().WithContext(ctx).Create(&microclimateReading).Error
	return err
}

//Creates new Microclimate readings (batch)
func (r *MicroclimateReadingRepository) CreateMicroclimateReadings(ctx context.Context, microclimateReadings []models.MicroclimateReading) error {
	var err error
	err = r.DB.GetClient().WithContext(ctx).Create(&microclimateReadings).Error
	return err
}

//Fetch Microclimate readings based on starting, ending date, microclimate parameter and location id
func (r *MicroclimateReadingRepository) GetMicroClimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error) {
	var err error
	microclimates := []models.MicroclimateReading{}
	queryStr := "microclimate_id = ? AND location_id = ? AND date <= ? AND date >= ?"

	//err = r.DB.WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Group("microclimate_id,location_id,date").Order("date").Find(&microclimates).Error
	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, microclimateID, locationID, toDate.Format("2006-01-02"), fromDate.Format("2006-01-02")).Order("date").Find(&microclimates).Error
	if err != nil {
		return nil, err
	}
	return microclimates, err

}

//Fetch first Microclimate reading based on microclimate parameter and location id
//if not found returns ErrRecordNotFound
func (r *MicroclimateReadingRepository) GetFirstMicroClimateReadingByID(ctx context.Context, microclimateID, locationID int) (models.MicroclimateReading, error) {
	var err error
	microclimateReading := models.MicroclimateReading{}
	queryStr := "location_id = ? AND microclimate_id = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID, microclimateID).Order("date").First(&microclimateReading).Error
	if err != nil {
		return models.MicroclimateReading{}, err
	}
	return microclimateReading, err
}

//Fetch latest Microclimate reading based on microclimate parameter and location id
//if not found returns ErrRecordNotFound
func (r *MicroclimateReadingRepository) GetLatestMicroClimateReadingByID(ctx context.Context, microclimateID, locationID int) (models.MicroclimateReading, error) {
	var err error
	microclimateReading := models.MicroclimateReading{}
	queryStr := "location_id = ? AND microclimate_id = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID, microclimateID).Order("date desc").First(&microclimateReading).Error
	if err != nil {
		return models.MicroclimateReading{}, err
	}
	return microclimateReading, err
}

//Fetch first Microclimate reading based on location id
//Since every date will always have 6 parameters, dont care for microclimate parameter id
//if not found returns ErrRecordNotFound
func (r *MicroclimateReadingRepository) GetFirstMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error) {

	var err error
	microclimateReading := models.MicroclimateReading{}
	queryStr := "location_id = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID).Order("date").First(&microclimateReading).Error
	if err != nil {
		return models.MicroclimateReading{}, err
	}
	return microclimateReading, err
}

//Fetch latest Microclimate reading based on location id
//Since every date will always have 6 parameters, dont care for microclimate parameter id
//if not found returns ErrRecordNotFound
func (r *MicroclimateReadingRepository) GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error) {

	var err error
	microclimateReading := models.MicroclimateReading{}
	queryStr := "location_id = ?"

	err = r.DB.GetClient().WithContext(ctx).Debug().Preload("Location").Model(&models.MicroclimateReading{}).Preload("Microclimate").Where(queryStr, locationID).Order("date desc").First(&microclimateReading).Error
	if err != nil {
		return models.MicroclimateReading{}, err
	}
	return microclimateReading, err
}

//Fetch Microclimate readings based on starting, ending date, microclimate parameter and location id in batch mode
//Method runs in gouroutine, fetches batch data and sends them on channel, when over closes channel and returns
func (r *MicroclimateReadingRepository) GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error {
	//cleanup work close chan
	//sender closes channel, sends a default value of the channel's type, if int 0, if bool false, if some struct, struct with default values etc...
	defer close(ch)

	results := []models.MicroclimateReading{}
	//result := r.DB.Debug().WithContext(ctxx).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Order("date, microclimate_id").Group("microclimate_id, location_id, date").FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
	//	for _, result := range results {
	//		//cancel()
	//		//fmt.Println("result:", result)
	//		select {
	//		case ch <- result:
	//			//fmt.Println("Poslao")
	//			fmt.Println(result)
	//		case <-ctxx.Done():
	//			fmt.Println("ctx batch microclimate", ctxx.Err())
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
		result := r.DB.GetClient().WithContext(ctxx).Limit(batchSize).Offset(batch*batchSize).Where("location_id = ? AND date >= ? AND date <= ?", locationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).Order("date, microclimate_id").Find(&results)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			//do something
			for _, result := range results {
				//fmt.Println(result)
				//blocking, waiting until chan is available or timeout
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
			//Last batch, not full size
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
