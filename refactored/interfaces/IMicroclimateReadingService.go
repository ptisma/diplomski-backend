package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"os"
	"time"
)

type IMicroclimateReadingService interface {
	GetMicroclimateReadings(microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error)
	GetLatestMicroClimateReading(locationID int) (models.MicroclimateReading, error)
	GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error
	GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error
	ReceiveFromBatchAndWrite(csvFile *os.File, batchCh chan models.MicroclimateReading, ctxx context.Context) error
}
