package interfaces

import (
	"apsim-api/internal/models"
	"context"
	"os"
	"time"
)

type IMicroclimateReadingService interface {
	GetMicroclimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error)

	GetPredictedMicroclimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error)

	GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error)

	GetMicroclimateReadingPeriod(ctx context.Context, microclimateID, locationID int) (models.Period, error)

	GetBatchMicroclimateReadings(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error

	GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error

	ReceiveFromBatchAndWrite(csvFile *os.File, batchCh chan models.MicroclimateReading, ctxx context.Context) error

	CreateMicroclimateReading(ctx context.Context, microclimateID, locationID int, date string, value float32) error

	CreateMicroclimateReadings(ctx context.Context, microclimateReadings []models.MicroclimateReading) error

	CalculateGrowingDegreeDay(ctx context.Context, tmaxReadings []models.MicroclimateReading, tminReadings []models.MicroclimateReading, baseTemp float32) ([]models.GrowingDegreeDay, error)

	CalculatePredictedGrowingDegreeDay(ctx context.Context, tmaxReadings []models.PredictedMicroclimateReading, tminReadings []models.PredictedMicroclimateReading, baseTemp float32) ([]models.GrowingDegreeDay, error)

	ConvertPredictedMicroclimateReadings(predictedMicroclimateReadings []models.PredictedMicroclimateReading) []models.MicroclimateReading
}
