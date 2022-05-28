package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"context"
	"fmt"
)

type MicroclimateService struct {
	I interfaces.IMicroclimateRepository
}

func (s *MicroclimateService) GetAllMicroclimates(ctx context.Context) ([]models.Microclimate, error) {
	fmt.Println("Sad sam u servisu")
	return s.I.GetAllMicroclimates(ctx)
}

func (s *MicroclimateService) GetMicroclimateByName(ctx context.Context, microclimateName string) (models.Microclimate, error) {
	return s.I.GetMicroclimateByName(ctx, microclimateName)
}
