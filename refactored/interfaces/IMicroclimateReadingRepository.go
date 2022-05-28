package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"time"
)

type IMicroclimateReadingRepository interface {
	CreateMicroclimateReading(ctx context.Context, microclimateID, locationID int, date string, value float32) error
	GetMicroClimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error)
	GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error)
	GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error
}
