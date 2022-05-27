package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"fmt"
)

type MicroclimateService struct {
	I interfaces.IMicroclimateRepository
}

func (s *MicroclimateService) GetAllMicroclimates() ([]models.Microclimate, error) {
	fmt.Println("Sad sam u servisu")
	return s.I.GetAllMicroclimates()
}

func (s *MicroclimateService) GetMicroclimateByName(microclimateName string) (models.Microclimate, error) {
	return s.I.GetMicroclimateByName(microclimateName)
}
