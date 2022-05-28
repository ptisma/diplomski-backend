package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"context"
)

type SoilService struct {
	I interfaces.ISoilRepository
}

func (s *SoilService) GetSoilByLocationId(ctx context.Context, locationId int) (models.Soil, error) {

	return s.I.GetSoilByLocationId(ctx, locationId)
}
