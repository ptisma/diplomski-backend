package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"time"
)

type IMicroclimateReadingRepository interface {
	GetMicroClimateReadings(microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error)
	GetLatestMicroClimateReading(locationID int) (models.MicroclimateReading, error)
	GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error
}
