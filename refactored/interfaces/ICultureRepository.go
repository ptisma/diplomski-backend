package interfaces

import "apsim-api/refactored/models"

type ICultureRepository interface {
	GetAllCultures() ([]models.Culture, error)
	GetCultureById(id int) (models.Culture, error)
}
