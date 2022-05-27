package interfaces

import "apsim-api/refactored/models"

type ILocationRepository interface {
	GetAllLocations() ([]models.Location, error)
	GetLocationById(locationId int) (models.Location, error)
}
