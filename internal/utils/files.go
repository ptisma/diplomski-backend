package utils

import (
	"os"
	"path/filepath"
)

func CreateTempStageFile(pattern string) (*os.File, string, error) {
	var f *os.File
	var err error
	var absPath string
	f, err = os.CreateTemp("./apsim", pattern)
	if err != nil {
		return f, absPath, err
	}
	absPath, err = filepath.Abs(f.Name())
	if err != nil {
		return f, absPath, err
	}
	return f, absPath, nil
}

//func GenerateConstsFile(locationName string, locationLatitude, locationLongitude float32, ch chan models.Message, mainCh chan models.Message, ctxx context.Context) error {
//	var err error
//	defer fmt.Println("Consts ending")
//
//	constsFile, constsFileAbs, err := utils.CreateTempStageFile("const*.txt")
//	ch <- models.Message{ID: 0, Payload: constsFileAbs}
//	mainCh <- models.Message{ID: 0, Payload: constsFileAbs}
//	if err != nil {
//		return err
//	}
//	fmt.Println("constsFile abs path:", constsFileAbs)
//	//Write into consts file
//	constsFileA, err := os.OpenFile(constsFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	//Write into consts file
//	_, err = constsFileA.WriteString(fmt.Sprintf("location = %s\nlatitude = %.2f (DECIMAL DEGREES)\nlongitude = %.2f (DECIMAL DEGREES)\n", locationName, locationLatitude, locationLongitude))
//	if err != nil {
//		return err
//	}
//	err = constsFileA.Close()
//	if err != nil {
//		return err
//	}
//	err = constsFile.Close()
//	if err != nil {
//		return err
//	}
//	return err
//
//}
//
//func GenerateAPSIMXFile(cultureId int, fromDate, toDate time.Time, soil models.Soil, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error {
//	defer fmt.Println("APSIMX ending")
//	var err error
//	apsimxFile, apsimxFileAbs, err := utils.CreateTempStageFile("apsimxFile*.apsimx")
//	if err != nil {
//		return err
//	}
//	mainCh <- models.Message{
//		ID:      2,
//		Payload: apsimxFileAbs,
//	}
//	fmt.Println("apsimxFile abs path:", apsimxFileAbs)
//	var csvFilePath, constsFilePath string
//	counter := 0
//	flag := false
//	for {
//		select {
//		case msg := <-ch:
//			fmt.Println(msg)
//			switch msg.ID {
//			case 0:
//				constsFilePath = msg.Payload
//				counter += 1
//				break
//			case 1:
//				csvFilePath = msg.Payload
//				counter += 1
//				break
//			default:
//				return errors.New("Received unexpected message")
//			}
//		//not sure if this is correct to place here and return that
//		case <-ctx.Done():
//			return ctx.Err()
//		default:
//			if counter == 2 {
//				flag = true
//				break
//			}
//
//		}
//		if flag == true {
//			break
//		}
//	}
//	//TODO
//	apsimBodyInit, _ := apsimx.InitiateAPSIMCulture[cultureId]
//	apsimxBody := apsimBodyInit(fromDate, toDate, csvFilePath, constsFilePath, soil.Data)
//	_, err = apsimxFile.WriteString(apsimxBody)
//	if err != nil {
//		return err
//	}
//	//can simulation run without apsimx file being closed(it can)
//	//	cancel()
//	//	return err
//	//}
//	return err
//}
//
//func GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error {
//
//}
