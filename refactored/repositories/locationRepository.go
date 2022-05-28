package repositories

import (
	"apsim-api/refactored/models"
	"context"

	"gorm.io/gorm"
)

type LocationRepository struct {
	DB *gorm.DB
}

func (r *LocationRepository) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	var err error
	locations := []models.Location{}
	err = r.DB.WithContext(ctx).Debug().Model(&models.Location{}).Find(&locations).Error
	if err != nil {
		return []models.Location{}, err
	}
	return locations, err

}

func (r *LocationRepository) GetLocationById(ctx context.Context, locationId int) (models.Location, error) {
	var err error
	location := models.Location{}
	err = r.DB.WithContext(ctx).Debug().First(&location, locationId).Error
	if err != nil {
		return models.Location{}, err
	}
	return location, err

}
