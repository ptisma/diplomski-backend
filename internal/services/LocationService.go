package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"fmt"
	"os"
)

type LocationService struct {
	I interfaces.ILocationRepository
}

func (s *LocationService) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	//fmt.Println("Sad sam u servisu")
	return s.I.GetAllLocations(ctx)
}

func (r *LocationService) GetLocationById(ctx context.Context, locationId int) (models.Location, error) {
	return r.I.GetLocationById(ctx, locationId)
}

func (r *LocationService) GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan models.Message, mainCh chan models.Message, ctxx context.Context) error {

	var err error
	defer func() {
		fmt.Println("CONSTS ending")
		//Because err group only logs one error
		if err != nil {
			fmt.Println("GenerateConstsFile err:", err)

		}
	}()

	constsFile, constsFileAbs, err := utils.CreateTempStageFile("const*.txt")
	if err != nil {
		return err
	}
	ch <- models.Message{ID: models.CONSTS_FILE_CODE, Payload: constsFileAbs}
	mainCh <- models.Message{ID: models.CONSTS_FILE_CODE, Payload: constsFileAbs}
	//if err != nil {
	//	return err
	//}
	//fmt.Println("constsFile abs path:", constsFileAbs)
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
	//err = constsFileA.Close()
	//if err != nil {
	//	return err
	//}
	//err = constsFile.Close()
	//if err != nil {
	//	return err
	//}

	_ = constsFileA.Close()
	_ = constsFile.Close()
	return err

	//close .apsimx file?

}
