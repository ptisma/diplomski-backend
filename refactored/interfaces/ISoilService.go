package interfaces

import "apsim-api/refactored/models"

type ISoilService interface {
	GetSoilByLocationId(locationId int) (models.Soil, error)
}
