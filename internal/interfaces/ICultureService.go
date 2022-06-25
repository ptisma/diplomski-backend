package interfaces

import (
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"time"
)

type ICultureService interface {
	FetchAllCultures(ctx context.Context) ([]models.Culture, error)
	FetchCultureById(ctx context.Context, id int) (models.Culture, error)
	GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan utils.Message, mainCh chan utils.Message, ctx context.Context) error
}
