package interfaces

import (
	"apsim-api/refactored/models"
	"context"
)

type ILocationService interface {
	GetAllLocations() ([]models.Location, error)
	GetLocationById(locationId int) (models.Location, error)
	GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan models.Message, mainCh chan models.Message, ctxx context.Context) error
}
