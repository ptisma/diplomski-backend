package interfaces

import (
	"apsim-api/refactored/models"
	"context"
)

type ISoilRepository interface {
	GetSoilByLocationId(ctx context.Context, locationId int) (models.Soil, error)
}
