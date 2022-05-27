package interfaces

import "apsim-api/refactored/models"

type IMicroclimateRepository interface {
	GetAllMicroclimates() ([]models.Microclimate, error)
	GetMicroclimateByName(microclimateName string) (models.Microclimate, error)
}
