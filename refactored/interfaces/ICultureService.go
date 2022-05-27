package interfaces

import (
	"apsim-api/refactored/models"
	"context"
	"time"
)

type ICultureService interface {
	FetchAllCultures() ([]models.Culture, error)
	FetchCultureById(id int) (models.Culture, error)
	GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error
}
