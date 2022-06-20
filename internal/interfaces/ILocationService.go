package interfaces

import (
	"apsim-api/internal/models"
	"context"
)

type ILocationService interface {
	GetAllLocations(ctx context.Context) ([]models.Location, error)
	GetLocationById(ctx context.Context, locationId int) (models.Location, error)
	GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan models.Message, mainCh chan models.Message, ctxx context.Context) error
}
