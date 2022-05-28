package repositories

import (
	"apsim-api/refactored/models"
	"context"
	"gorm.io/gorm"
)

type MicroclimateRepository struct {
	DB *gorm.DB
}

func (r *MicroclimateRepository) GetAllMicroclimates(ctx context.Context) ([]models.Microclimate, error) {
	var err error
	microclimates := []models.Microclimate{}
	err = r.DB.WithContext(ctx).Debug().Model(&models.Microclimate{}).Find(&microclimates).Error
	if err != nil {
		return []models.Microclimate{}, err
	}
	return microclimates, err

}

func (r *MicroclimateRepository) GetMicroclimateByName(ctx context.Context, microclimateName string) (models.Microclimate, error) {
	microclimate := models.Microclimate{}
	err := r.DB.WithContext(ctx).Debug().Model(&models.Microclimate{}).Where("name = ?", microclimateName).First(&microclimate).Error
	return microclimate, err

}
