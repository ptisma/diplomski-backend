package repositories

import (
	"apsim-api/refactored/models"

	"gorm.io/gorm"
)

type LocationRepository struct {
	DB *gorm.DB
}

func (r *LocationRepository) GetAllLocations() ([]models.Location, error) {
	var err error
	locations := []models.Location{}
	err = r.DB.Debug().Model(&models.Location{}).Find(&locations).Error
	if err != nil {
		return []models.Location{}, err
	}
	return locations, err

}

func (r *LocationRepository) GetLocationById(locationId int) (models.Location, error) {
	var err error
	location := models.Location{}
	err = r.DB.Debug().First(&location, locationId).Error
	if err != nil {
		return models.Location{}, err
	}
	return location, err

}
