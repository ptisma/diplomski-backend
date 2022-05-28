package services

import (
	"apsim-api/internal/apsimx"
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"apsim-api/refactored/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type CultureService struct {
	interfaces.ICultureRepository
}

func (s *CultureService) FetchAllCultures(ctx context.Context) ([]models.Culture, error) {
	fmt.Println("Sad sam u servisu")
	return s.GetAllCultures(ctx)
}

func (s *CultureService) FetchCultureById(ctx context.Context, id int) (models.Culture, error) {
	return s.GetCultureById(ctx, id)
}

func (s *CultureService) GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error {
	defer fmt.Println("APSIMX ending")
	var err error
	apsimxFile, apsimxFileAbs, err := utils.CreateTempStageFile("apsimxFile*.apsimx")
	if err != nil {
		return err
	}
	mainCh <- models.Message{
		ID:      2,
		Payload: apsimxFileAbs,
	}
	fmt.Println("apsimxFile abs path:", apsimxFileAbs)
	var csvFilePath, constsFilePath string
	counter := 0
	flag := false
	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
			switch msg.ID {
			case 0:
				constsFilePath = msg.Payload
				counter += 1
				break
			case 1:
				csvFilePath = msg.Payload
				counter += 1
				break
			default:
				return errors.New("Received unexpected message")
			}
		//not sure if this is correct to place here and return that
		case <-ctx.Done():
			return ctx.Err()
		default:
			if counter == 2 {
				flag = true
				break
			}

		}
		if flag == true {
			break
		}

	}
	apsimBodyInit, _ := apsimx.InitiateAPSIMCulture[cultureId]
	apsimxBody := apsimBodyInit(fromDate, toDate, csvFilePath, constsFilePath, soil.Data)
	_, err = apsimxFile.WriteString(apsimxBody)
	if err != nil {
		return err
	}
	//can simulation run without apsimx file being closed(it can)
	//	cancel()
	//	return err
	//}
	return err
}
