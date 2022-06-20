package repositories

import (
	"apsim-api/internal/models"
	"context"
	"gorm.io/gorm"
)

type SoilRepository struct {
	DB *gorm.DB
}

func (r *SoilRepository) GetSoilByLocationId(ctx context.Context, locationId int) (models.Soil, error) {
	var err error
	soil := models.Soil{}
	err = r.DB.WithContext(ctx).Debug().Model(&models.Soil{}).Preload("Location").Where("location_id = ?", locationId).Find(&soil).Error
	return soil, err
}
