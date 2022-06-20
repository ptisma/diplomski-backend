package interfaces

import (
	"apsim-api/internal/models"
	"context"
	"time"
)

type IYieldRepository interface {
	GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error)
	CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error
}
