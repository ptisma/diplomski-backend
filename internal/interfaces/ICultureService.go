package interfaces

import (
	"apsim-api/internal/models"
	"context"
	"time"
)

type ICultureService interface {
	FetchAllCultures(ctx context.Context) ([]models.Culture, error)
	FetchCultureById(ctx context.Context, id int) (models.Culture, error)
	GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error
}
