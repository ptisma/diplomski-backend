package locationService

import (
	"apsim-api/internal/apsimx"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"apsim-api/pkg/application"
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

type LocationService struct {
	app *application.Application
}

func GetLocationService(app *application.Application) *LocationService {
	return &LocationService{
		app: app,
	}
}

func (ls *LocationService) GetLocation(locationId int) (models.Location, error) {
	var err error
	location := models.Location{ID: uint32(locationId)}
	err = location.GetLocationById(ls.app)

	return location, err

}

func (ls *LocationService) GetLocations() ([]models.Location, error) {
	location := models.Location{}
	return location.GetAllLLocations(ls.app)
}

func (ls *LocationService) GenerateConstsFile(location models.Location, ch chan models.Message, mainCh chan models.Message, ctx context.Context, cancel context.CancelFunc) error {
	var err error
	defer fmt.Println("Consts ending")

	constsFile, constsFileAbs, err := utils.CreateTempStageFile("const*.txt")
	ch <- models.Message{ID: 0, Payload: constsFileAbs}
	mainCh <- models.Message{ID: 0, Payload: constsFileAbs}
	if err != nil {
		cancel()
		return err
	}
	fmt.Println("constsFile abs path:", constsFileAbs)
	//Write into consts file
	constsFileA, err := os.OpenFile(constsFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		cancel()
		return err
	}
	//Write into consts file
	_, err = constsFileA.WriteString(fmt.Sprintf("location = %s\nlatitude = %.2f (DECIMAL DEGREES)\nlongitude = %.2f (DECIMAL DEGREES)\n", location.Name, location.Latitude, location.Longitude))
	if err != nil {
		cancel()
		return err
	}
	err = constsFileA.Close()
	if err != nil {
		cancel()
		return err
	}
	err = constsFile.Close()
	if err != nil {
		cancel()
		return err
	}
	return err

}

func (ls *LocationService) GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context, cancel context.CancelFunc) error {
	defer fmt.Println("APSIMX ending")
	var err error
	apsimxFile, apsimxFileAbs, err := utils.CreateTempStageFile("apsimxFile*.apsimx")
	if err != nil {
		cancel()
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
				cancel()
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
	//TODO
	apsimBodyInit, _ := apsimx.InitiateAPSIMCulture[cultureId]
	apsimxBody := apsimBodyInit(fromDate, toDate, csvFilePath, constsFilePath, soil.Data)
	_, err = apsimxFile.WriteString(apsimxBody)
	if err != nil {
		cancel()
		return err
	}
	//can simulation run without apsimx file being closed(it can)
	//	cancel()
	//	return err
	//}
	return err

}
