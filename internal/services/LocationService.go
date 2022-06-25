package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"fmt"
	"log"
	"os"
)

type LocationService struct {
	I interfaces.ILocationRepository
}

func (s *LocationService) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	return s.I.GetAllLocations(ctx)
}

func (r *LocationService) GetLocationById(ctx context.Context, locationId int) (models.Location, error) {
	return r.I.GetLocationById(ctx, locationId)
}

// Generate consts file in stage area
// Send to the main chan the absolute path of generated file
// Send to the secondary chan the absolute path of generated file
func (r *LocationService) GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan utils.Message, mainCh chan utils.Message, ctxx context.Context) error {

	var err error
	var constsFile, constsFileA *os.File

	//clean up function
	//close file so it can be deleted
	//log error
	defer func() {
		log.Println("CONSTS ending")

		_ = constsFileA.Close()
		_ = constsFile.Close()

		//Because err group only logs one error
		if err != nil {
			log.Println("GenerateConstsFile err:", err)

		}
	}()

	//Create a consts text file
	constsFile, constsFileAbs, err := utils.CreateTempStageFile("const*.txt")
	if err != nil {
		return err
	}
	ch <- utils.Message{ID: utils.CONSTS_FILE_CODE, Payload: constsFileAbs}
	mainCh <- utils.Message{ID: utils.CONSTS_FILE_CODE, Payload: constsFileAbs}

	//Appending to file
	constsFileA, err = os.OpenFile(constsFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//Write into consts file
	//no context method in standard library (byte by byte with select block checking context if there is a huge text)
	_, err = constsFileA.WriteString(fmt.Sprintf("location = %s\nlatitude = %.2f (DECIMAL DEGREES)\nlongitude = %.2f (DECIMAL DEGREES)\n", locationName, locationLatitude, locationLongitude))
	if err != nil {
		return err
	}

	return err

}
