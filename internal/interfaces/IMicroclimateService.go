package interfaces

import (
	"apsim-api/internal/models"
	"context"
)

type IMicroclimateService interface {
	GetAllMicroclimates(ctx context.Context) ([]models.Microclimate, error)
	GetMicroclimateByName(ctx context.Context, microclimateName string) (models.Microclimate, error)
}
