package interfaces

import (
	"apsim-api/internal/models"
	"context"
)

type ILocationRepository interface {
	GetAllLocations(ctx context.Context) ([]models.Location, error)
	GetLocationById(ctx context.Context, locationId int) (models.Location, error)
}
