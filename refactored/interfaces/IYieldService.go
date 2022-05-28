package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"time"
)

type IYieldService interface {
	GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error)
	CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error
}
