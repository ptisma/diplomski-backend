package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
)

type LocationRepository struct {
	//DB *gorm.DB
	DB db.DB
}

//Fetch all locations
func (r *LocationRepository) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	var err error
	locations := []models.Location{}
	err = r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Location{}).Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, err

}

//Fetch Location based on ID
//if not found returns ErrRecordNotFound
func (r *LocationRepository) GetLocationById(ctx context.Context, locationId int) (models.Location, error) {
	var err error
	location := models.Location{}
	err = r.DB.GetClient().WithContext(ctx).Debug().First(&location, locationId).Error
	if err != nil {
		return models.Location{}, err
	}
	return location, err

}
