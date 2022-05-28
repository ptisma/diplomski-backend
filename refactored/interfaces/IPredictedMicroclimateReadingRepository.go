package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"time"
)

type IPredictedMicroclimateReadingRepository interface {
	GetMicroClimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error)
	GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.PredictedMicroclimateReading, error)
	GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.PredictedMicroclimateReading, ctxx context.Context) error
}
