package interfaces

import (
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
)

type ILocationService interface {
	GetAllLocations(ctx context.Context) ([]models.Location, error)
	GetLocationById(ctx context.Context, locationId int) (models.Location, error)
	GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan utils.Message, mainCh chan utils.Message, ctxx context.Context) error
}
