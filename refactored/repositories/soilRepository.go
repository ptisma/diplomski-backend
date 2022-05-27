package repositories

import (
	"apsim-api/refactored/models"
	"gorm.io/gorm"
)

type SoilRepository struct {
	DB *gorm.DB
}

func (r *SoilRepository) GetSoilByLocationId(locationId int) (models.Soil, error) {
	var err error
	soil := models.Soil{}
	err = r.DB.Debug().Model(&models.Soil{}).Preload("Location").Where("location_id = ?", locationId).Find(&soil).Error
	return soil, err
}
