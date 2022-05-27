package repositories

import (
	"apsim-api/refactored/models"
	"gorm.io/gorm"
)

type MicroclimateRepository struct {
	DB *gorm.DB
}

func (r *MicroclimateRepository) GetAllMicroclimates() ([]models.Microclimate, error) {
	var err error
	microclimates := []models.Microclimate{}
	err = r.DB.Debug().Model(&models.Microclimate{}).Find(&microclimates).Error
	if err != nil {
		return []models.Microclimate{}, err
	}
	return microclimates, err

}

func (r *MicroclimateRepository) GetMicroclimateByName(microclimateName string) (models.Microclimate, error) {
	microclimate := models.Microclimate{}
	err := r.DB.Debug().Model(&models.Microclimate{}).Where("name = ?", microclimateName).First(&microclimate).Error
	return microclimate, err

}
