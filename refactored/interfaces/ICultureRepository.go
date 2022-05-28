package interfaces

import (
	"apsim-api/refactored/models"
	"context"
)

type ICultureRepository interface {
	GetAllCultures(ctx context.Context) ([]models.Culture, error)
	GetCultureById(ctx context.Context, id int) (models.Culture, error)
}
