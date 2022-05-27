package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"apsim-api/refactored/utils"
	"context"
	"fmt"
	"os"
)

type LocationService struct {
	I interfaces.ILocationRepository
}

func (s *LocationService) GetAllLocations() ([]models.Location, error) {
	fmt.Println("Sad sam u servisu")
	return s.I.GetAllLocations()
}

func (r *LocationService) GetLocationById(locationId int) (models.Location, error) {
	return r.I.GetLocationById(locationId)
}

func (r *LocationService) GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan models.Message, mainCh chan models.Message, ctxx context.Context) error {

	var err error
	defer fmt.Println("Consts ending")

	constsFile, constsFileAbs, err := utils.CreateTempStageFile("const*.txt")
	ch <- models.Message{ID: 0, Payload: constsFileAbs}
	mainCh <- models.Message{ID: 0, Payload: constsFileAbs}
	if err != nil {
		return err
	}
	fmt.Println("constsFile abs path:", constsFileAbs)
	//Write into consts file
	constsFileA, err := os.OpenFile(constsFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//Write into consts file
	_, err = constsFileA.WriteString(fmt.Sprintf("location = %s\nlatitude = %.2f (DECIMAL DEGREES)\nlongitude = %.2f (DECIMAL DEGREES)\n", locationName, locationLatitude, locationLongitude))
	if err != nil {
		return err
	}
	err = constsFileA.Close()
	if err != nil {
		return err
	}
	err = constsFile.Close()
	if err != nil {
		return err
	}
	return err

}
