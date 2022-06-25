package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
)

type MicroclimateRepository struct {
	//DB *gorm.DB
	DB db.DB
}

//Fetch all microclimate parameters
func (r *MicroclimateRepository) GetAllMicroclimates(ctx context.Context) ([]models.Microclimate, error) {
	var err error
	microclimates := []models.Microclimate{}
	err = r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Microclimate{}).Find(&microclimates).Error
	if err != nil {
		return nil, err
	}
	return microclimates, err

}

//Fetch Microclimate parameter based on Name
//if not found returns ErrRecordNotFound
func (r *MicroclimateRepository) GetMicroclimateByName(ctx context.Context, microclimateName string) (models.Microclimate, error) {
	microclimate := models.Microclimate{}
	err := r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Microclimate{}).Where("name = ?", microclimateName).First(&microclimate).Error
	if err != nil {
		return models.Microclimate{}, err
	}
	return microclimate, err

}
