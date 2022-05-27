package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
)

type SoilService struct {
	I interfaces.ISoilRepository
}

func (s *SoilService) GetSoilByLocationId(locationId int) (models.Soil, error) {

	return s.I.GetSoilByLocationId(locationId)
}
