package interfaces

import "apsim-api/refactored/models"

type ISoilRepository interface {
	GetSoilByLocationId(locationId int) (models.Soil, error)
}
