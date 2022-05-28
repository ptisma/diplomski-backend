package interfaces

import (
	"apsim-api/refactored/models"
	"context"
)

type ISoilService interface {
	GetSoilByLocationId(ctx context.Context, locationId int) (models.Soil, error)
}
