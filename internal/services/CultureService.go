package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"
)

type CultureService struct {
	interfaces.ICultureRepository
}

func (s *CultureService) FetchAllCultures(ctx context.Context) ([]models.Culture, error) {
	//fmt.Println("Sad sam u servisu")
	return s.GetAllCultures(ctx)
}

func (s *CultureService) FetchCultureById(ctx context.Context, id int) (models.Culture, error) {
	return s.GetCultureById(ctx, id)
}

func (s *CultureService) GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error {
	var err error
	defer func() {
		fmt.Println("APSIMX ending")
		//Because err group only logs one first error from whole group
		if err != nil {
			fmt.Println("GenerateAPSIMXFile err:", err)

		}
	}()
	apsimxFile, apsimxFileAbs, err := utils.CreateTempStageFile("apsimxFile*.apsimx")
	if err != nil {
		return err
	}
	mainCh <- models.Message{
		ID:      models.APSIMX_FILE_CODE,
		Payload: apsimxFileAbs,
	}
	//fmt.Println("apsimxFile abs path:", apsimxFileAbs)
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
	apsimBodyInit, _ := models.InitiateAPSIMCulture[cultureId]
	apsimxBody := apsimBodyInit(fromDate, toDate, csvFilePath, constsFilePath, soil.Data)
	_, err = apsimxFile.WriteString(apsimxBody)
	if err != nil {
		return err
	}

	_ = apsimxFile.Close()
	//can simulation run without apsimx file being closed(it can)
	//	cancel()
	//	return err
	//}
	return err
}
