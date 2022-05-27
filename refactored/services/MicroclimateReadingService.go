package services

import (
	"apsim-api/refactored/interfaces"
	"apsim-api/refactored/models"
	"apsim-api/refactored/utils"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
	"time"
)

type MicroclimateReadingService struct {
	I1 interfaces.IMicroclimateReadingRepository
	I2 interfaces.IPredictedMicroclimateReadingRepository
}

func (s *MicroclimateReadingService) GetMicroclimateReadings(microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error) {
	fmt.Println("Sad sam u servisu")
	return s.I1.GetMicroClimateReadings(microclimateID, locationID, fromDate, toDate)
}
func (s *MicroclimateReadingService) GetLatestMicroClimateReading(locationID int) (models.MicroclimateReading, error) {
	fmt.Println("Sad sam u servisu")
	return s.I1.GetLatestMicroClimateReading(locationID)
}

func (s *MicroclimateReadingService) GetBatchMicroclimateReading(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error {
	return s.I1.GetBatchMicroclimateReading(locationID, fromDate, toDate, ch, ctxx)
}

func (s *MicroclimateReadingService) ReceiveFromBatchAndWrite(csvFile *os.File, batchCh chan models.MicroclimateReading, ctxx context.Context) error {
	var err error
	buff := []models.MicroclimateReading{}
	counter := 0
	var sb strings.Builder
	counterWrite := 0
	for {
		select {
		case msg := <-batchCh:
			//fmt.Println("msg:", msg)
			//check if the chan is closed
			if (msg == models.MicroclimateReading{}) {
				fmt.Println("empty struct, chan is closed")
				fmt.Println("remaining buff", buff)
				fmt.Println("remaining sb", sb.String())
				_, err = csvFile.WriteString(sb.String())
				if err != nil {
					fmt.Println("batch microclimate error in writing", err)
					return err
				}
				return err
			}
			counter += 1
			buff = append(buff, msg)
			if counter == 6 {
				rowDate, err := time.Parse("2006-01-02", buff[0].Date)
				if err != nil {
					return err
				}
				csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
				//_, err = csvFile.WriteString(csvRow)
				sb.WriteString(csvRow)
				counterWrite += 1
				if counterWrite == 6 {
					_, err = csvFile.WriteString(sb.String())
					if err != nil {
						fmt.Println("batch microclimate error in writing", err)
						return err
					}
					sb.Reset()
					counterWrite = 0
				}
				counter = 0
				buff = nil
			}

		case <-ctxx.Done():
			fmt.Println("csv", ctxx.Err())
			fmt.Println("ctx done in generate csv file")
			return err
			//case _ = <-endCh:
			//	fmt.Println("Batch goroutine done")
			//	flag = true
			//	break
		}

	}
}

func (s *MicroclimateReadingService) GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context) error {
	defer fmt.Println("CSV ending")
	var err error

	csvFileOrig, csvFileAbs, err := utils.CreateTempStageFile("csv*.csv")
	if err != nil {
		return err
	}
	ch <- models.Message{ID: 1, Payload: csvFileAbs}
	mainCh <- models.Message{ID: 1, Payload: csvFileAbs}

	fmt.Println("CSVFile abs path:", csvFileAbs)

	//append
	csvFile, err := os.OpenFile(csvFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//first row
	_, err = csvFile.WriteString("year,day,radn,maxt,mint,rain,rh,wind,code\n")
	if err != nil {
		return err
	}

	//launch a batch getter in separate goroutine
	g, ctxx := errgroup.WithContext(ctx)
	batchCh := make(chan models.MicroclimateReading, 50)
	//endCh := make(chan string, 1)
	//buff := []models.MicroclimateReading{}
	//counter := 0
	//flag := false
	g.Go(func() error {
		return s.GetBatchMicroclimateReading(locationId, fromDate, toDate, batchCh, ctxx)
	})
	err = s.ReceiveFromBatchAndWrite(csvFile, batchCh, ctxx)
	if err != nil {
		return err
	}
	err = g.Wait()
	fmt.Println("Batch error group", err)
	if err != nil {
		return err
	}
	_ = csvFileOrig.Close()

	_ = csvFile.Close()

	return err

}
