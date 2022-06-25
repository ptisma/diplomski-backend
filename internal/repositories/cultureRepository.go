package repositories

import (
	"apsim-api/internal/infra/db"
	"apsim-api/internal/models"
	"context"
)

type CultureRepository struct {
	//DB *gorm.DB
	DB db.DB
}

//Fetch all cultures
func (r *CultureRepository) GetAllCultures(ctx context.Context) ([]models.Culture, error) {
	//fmt.Println("Sad sam u repozitoriju")
	var err error
	cultures := []models.Culture{}
	err = r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Culture{}).Find(&cultures).Error
	if err != nil {
		//nil for referenced types, empty value for non ref
		return nil, err
	}
	return cultures, err
}

//Fetch Culture based on ID
//if not found returns ErrRecordNotFound
func (r *CultureRepository) GetCultureById(ctx context.Context, id int) (models.Culture, error) {
	var err error
	culture := models.Culture{}
	err = r.DB.GetClient().WithContext(ctx).Debug().Model(&models.Culture{}).First(&culture, id).Error
	if err != nil {
		return models.Culture{}, err
	}
	return culture, err
}
