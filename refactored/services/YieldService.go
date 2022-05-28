package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
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
