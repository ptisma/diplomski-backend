package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"context"
	"time"
)

type YieldService struct {
	I interfaces.IYieldRepository
}

func (s *YieldService) GetYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time) ([]models.Yield, error) {

	return s.I.GetYields(ctx, locationId, cultureId, fromDate, toDate)
}

func (s *YieldService) CreateYields(ctx context.Context, locationId, cultureId int, fromDate, toDate time.Time, yields []models.Yield) error {

	return s.I.CreateYields(ctx, locationId, cultureId, fromDate, toDate, yields)
}

func (s *YieldService) ValidateYields(fromDate, toDate time.Time, yields []models.Yield) bool {
	if len(yields) != (toDate.Year() - fromDate.Year() + 1) {
		return false
	}
	return true
}
