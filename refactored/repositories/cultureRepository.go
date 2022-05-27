package repositories

import (
	"apsim-api/refactored/models"
	"fmt"
	"gorm.io/gorm"
)

type CultureRepository struct {
	DB *gorm.DB
}

func (r *CultureRepository) GetAllCultures() ([]models.Culture, error) {
	fmt.Println("Sad sam u repozitoriju")
	var err error
	cultures := []models.Culture{}
	err = r.DB.Debug().Model(&models.Culture{}).Find(&cultures).Error
	if err != nil {
		return []models.Culture{}, err
	}
	return cultures, err
}

func (r *CultureRepository) GetCultureById(id int) (models.Culture, error) {
	var err error
	culture := models.Culture{}
	err = r.DB.Debug().Model(&models.Culture{}).First(&culture, id).Error
	return culture, err
}
