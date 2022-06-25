package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
)

type SoilRepository struct {
	//DB *gorm.DB
	DB db.DB
}

//Fetch Soil based on ID
//if not found returns ErrRecordNotFound
func (r *SoilRepository) GetSoilByLocationId(ctx context.Context, locationId int) (models.Soil, error) {
	var err error
	soil := models.Soil{}
	err = r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Soil{}).Preload("Location").Where("location_id = ?", locationId).Find(&soil).Error
	if err != nil {
		return models.Soil{}, err
	}
	return soil, err
}
