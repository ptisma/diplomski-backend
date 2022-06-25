package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"errors"
	"log"
	"os"
	"time"
)

type CultureService struct {
	interfaces.ICultureRepository
}

func (s *CultureService) FetchAllCultures(ctx context.Context) ([]models.Culture, error) {
	return s.GetAllCultures(ctx)
}

func (s *CultureService) FetchCultureById(ctx context.Context, id int) (models.Culture, error) {
	return s.GetCultureById(ctx, id)
}

// Generate .apsimxfile in stage area
// Send to the main chan the absolute path of generated file
// Receive from secondary chan the absolute paths of other files which are coming from concurent goroutines creating their own files
func (s *CultureService) GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan utils.Message, mainCh chan utils.Message, ctx context.Context) error {
	var err error
	var apsimxFile *os.File

	//clean up function
	//close file so it can be deleted
	//log error
	//TODO instead of printing err, make error processing component which stores errors from all goroutines
	defer func() {
		log.Println("APSIMX ending")
		_ = apsimxFile.Close()
		//Because err group only logs one first error from whole group
		if err != nil {
			log.Println("GenerateAPSIMXFile err:", err)

		}
	}()
	apsimxFile, apsimxFileAbs, err := utils.CreateTempStageFile("apsimxFile*.apsimx")
	if err != nil {
		return err
	}
	mainCh <- utils.Message{
		ID:      utils.APSIMX_FILE_CODE,
		Payload: apsimxFileAbs,
	}
	//fmt.Println("apsimxFile abs path:", apsimxFileAbs)
	var csvFilePath, constsFilePath string
	counter := 0
	flag := false
	for {
		//blocking
		select {
		case msg := <-ch:
			//fmt.Println(msg)
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
		case <-ctx.Done():
			return ctx.Err()
		default:
			if counter == 2 {
				//TODO cleanup work for secondary chan, receiver closing chan if others are all over?
				close(ch)
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

	//can simulation run without apsimx file being closed(it can)
	return err
}
