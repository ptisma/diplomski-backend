package interfaces

import (
	"apsim-api/internal/models"
	"context"
	"time"
)

type IYieldService interface {
	GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error)

	CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error

	ValidateYields(fromDate, toDate time.Time, yields []models.Yield) bool

	ReadYields(dbFilePath string) ([]models.Yield, error)

	ReadYieldErrors(dbFilePath string) ([]string, error)
}
