package microclimateReadingService

import (
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"apsim-api/pkg/application"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
	"time"
)

type MicroclimateReadingService struct {
	app *application.Application
}

func GetMicroclimateReadingService(app *application.Application) *MicroclimateReadingService {
	return &MicroclimateReadingService{
		app: app,
	}
}

func (mrs *MicroclimateReadingService) GetMicroclimateReadings(microclimateId, locationId int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error) {
	var err error
	var readings []models.MicroclimateReading

	microclimateReading := models.MicroclimateReading{
		MicroclimateID: uint32(microclimateId),
		LocationID:     uint32(locationId),
		FromDate:       fromDate,
		ToDate:         toDate,
	}

	readings, err = microclimateReading.GetMicroclimateReading(mrs.app)

	return readings, err
}

func (mrs *MicroclimateReadingService) GetLatestMicroclimateReading(locationId int) (models.MicroclimateReading, error) {
	var err error

	microclimateReading := models.MicroclimateReading{
		LocationID: uint32(locationId),
	}

	err = microclimateReading.GetLatestMicroclimateReading(mrs.app)

	return microclimateReading, err
}

//change to float
func (mrs *MicroclimateReadingService) CalculateGrowingDegreeDay(tmaxReadings []models.MicroclimateReading, tminReadings []models.MicroclimateReading, baseTemp int) ([]models.GrowingDegreeDay, error) {
	var err error
	var gdds = []models.GrowingDegreeDay{}

	if len(tmaxReadings) == len(tminReadings) {
		for i, _ := range tmaxReadings {
			date := tmaxReadings[i].Date
			gdd := (tmaxReadings[i].Value+tminReadings[i].Value)/2 - float32(baseTemp)

			gdds = append(gdds, models.GrowingDegreeDay{
				Date:  date,
				Value: gdd,
			})

		}
	} else {
		err = errors.Errorf("tmax and tmin readings are not equal")
	}

	return gdds, err
}

//change to float
func (mrs *MicroclimateReadingService) CalculatePredictedGrowingDegreeDay(tmaxReadings []models.PredictedMicroclimateReading, tminReadings []models.PredictedMicroclimateReading, baseTemp int) ([]models.GrowingDegreeDay, error) {

	return mrs.CalculateGrowingDegreeDay(mrs.ConvertPredictedMicroclimateReadings(tmaxReadings), mrs.ConvertPredictedMicroclimateReadings(tminReadings), baseTemp)
}

func (mrs *MicroclimateReadingService) ConvertPredictedMicroclimateReadings(predictedMicroclimateReadings []models.PredictedMicroclimateReading) []models.MicroclimateReading {
	mr := make([]models.MicroclimateReading, 0, len(predictedMicroclimateReadings))
	for _, j := range predictedMicroclimateReadings {
		mr = append(mr, models.MicroclimateReading{
			ID:             j.ID,
			MicroclimateID: j.MicroclimateID,
			Microclimate:   j.Microclimate,
			LocationID:     j.LocationID,
			Location:       j.Location,
			Date:           j.Date,
			Value:          j.Value,
			FromDate:       j.FromDate,
			ToDate:         j.ToDate,
		})
	}
	return mr
}

func (mrs *MicroclimateReadingService) GetPredictedMicroclimateReadings(microclimateId, locationId int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error) {
	var err error
	var readings []models.PredictedMicroclimateReading

	predictedMicroclimateReading := models.PredictedMicroclimateReading{
		MicroclimateID: uint32(microclimateId),
		LocationID:     uint32(locationId),
		FromDate:       fromDate,
		ToDate:         toDate,
	}

	readings, err = predictedMicroclimateReading.GetPredictedMicroclimateReading(mrs.app)

	return readings, err
}

func (mrs *MicroclimateReadingService) GetBatchMicroclimateReadings(locationId int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error {
	var err error

	microclimateReading := models.MicroclimateReading{
		LocationID: uint32(locationId),
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	fmt.Println("Zovem")
	err = microclimateReading.GetBatchMicroclimateReading(mrs.app, ch, ctxx)

	return err
}

func (mrs *MicroclimateReadingService) GetBatchPredictedMicroclimateReadings(locationId int, fromDate, toDate time.Time, ch chan models.PredictedMicroclimateReading, ctxx context.Context) error {
	var err error

	predictedMicroclimateReading := models.PredictedMicroclimateReading{
		LocationID: uint32(locationId),
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	//fmt.Println("Zovem")
	err = predictedMicroclimateReading.GetBatchMicroclimateReading(mrs.app, ch, ctxx)

	return err
}

func (mrs *MicroclimateReadingService) GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context, cancel context.CancelFunc) error {
	defer fmt.Println("CSV ending")
	var err error

	csvFileOrig, csvFileAbs, err := utils.CreateTempStageFile("csv*.csv")
	if err != nil {
		cancel()
		return err
	}
	ch <- models.Message{ID: 1, Payload: csvFileAbs}
	mainCh <- models.Message{ID: 1, Payload: csvFileAbs}

	fmt.Println("CSVFile abs path:", csvFileAbs)

	//append
	csvFile, err := os.OpenFile(csvFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		cancel()
		return err
	}
	//first row
	_, err = csvFile.WriteString("year,day,radn,maxt,mint,rain,rh,wind,code\n")
	if err != nil {
		cancel()
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
		return mrs.GetBatchMicroclimateReadings(locationId, fromDate, toDate, batchCh, ctxx)
	})
	err = mrs.ReceiveFromBatchAndWrite(csvFile, batchCh, ctxx)
	if err != nil {
		cancel()
	}
	//for {
	//	select {
	//	case msg := <-batchCh:
	//		fmt.Println("msg:", msg)
	//		//check if the chan is closed
	//		if (msg == models.MicroclimateReading{}) {
	//			fmt.Println("empty struct, chan is closed")
	//			flag = true
	//			break
	//		}
	//		counter += 1
	//		buff = append(buff, msg)
	//		if counter == 6 {
	//			rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
	//			csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
	//			_, err := csvFile.WriteString(csvRow)
	//			if err != nil {
	//				//TODO NEW ERROR OR RETURN CONTEXT ERROR CANCELED
	//				fmt.Println("batch microclimate error in writing", err)
	//				cancel()
	//				flag = true
	//				break
	//			}
	//			counter = 0
	//			buff = nil
	//		}
	//
	//	case <-ctxx.Done():
	//		fmt.Println("csv", ctx.Err())
	//		fmt.Println("ctx done in generate csv file")
	//		flag = true
	//		break
	//		//case _ = <-endCh:
	//		//	fmt.Println("Batch goroutine done")
	//		//	flag = true
	//		//	break
	//	}
	//	if flag == true {
	//		break
	//	}
	//}

	err = g.Wait()
	fmt.Println("Batch error group", err)
	if err != nil {
		return err
	}
	//TODO REFACTOR
	//Predicted
	if err == nil && toDate.After(lastDate) {

		g, ctxx := errgroup.WithContext(ctx)
		batchCh := make(chan models.PredictedMicroclimateReading, 6)
		//buff := []models.PredictedMicroclimateReading{}
		//counter := 0
		//flag := false
		g.Go(func() error {
			return mrs.GetBatchPredictedMicroclimateReadings(locationId, lastDate.AddDate(0, 0, 1), toDate, batchCh, ctxx)
		})
		fmt.Println("USAO U PREDICTED")

		err = mrs.ReceiveFromPredictedBatchAndWrite(csvFile, batchCh, ctxx)
		if err != nil {
			cancel()
		}
		//for {
		//	select {
		//	case msg := <-batchCh:
		//		//fmt.Println("msg:", msg)
		//		//check if the chan is closed
		//		if (msg == models.PredictedMicroclimateReading{}) {
		//			fmt.Println("empty struct, chan is closed")
		//			flag = true
		//			break
		//		}
		//		counter += 1
		//		buff = append(buff, msg)
		//		if counter == 6 {
		//			rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
		//			csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
		//			csvFile.WriteString(csvRow)
		//			counter = 0
		//			buff = nil
		//		}
		//
		//	case <-ctx.Done():
		//		fmt.Println("ctx done in predicted generate csv file")
		//		flag = true
		//		break
		//		//case _ = <-endCh:
		//		//	fmt.Println("Batch goroutine done")
		//		//	flag = true
		//		//	break
		//	}
		//	if flag == true {
		//		break
		//	}
		//}

		err = g.Wait()
		fmt.Println("Predicted Batch error group", err)

	}

	//TODO ERROR CANT SHADOW EXISTING ERROR
	_ = csvFileOrig.Close()

	_ = csvFile.Close()

	return err
}

func (mrs *MicroclimateReadingService) ReceiveFromBatchAndWrite(csvFile *os.File, batchCh chan models.MicroclimateReading, ctxx context.Context) error {
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

func (mrs *MicroclimateReadingService) ReceiveFromPredictedBatchAndWrite(csvFile *os.File, batchCh chan models.PredictedMicroclimateReading, ctxx context.Context) error {
	var err error
	buff := []models.PredictedMicroclimateReading{}
	counter := 0
	var sb strings.Builder
	counterWrite := 0
	for {
		select {
		case msg := <-batchCh:
			//fmt.Println("msg:", msg)
			//check if the chan is closed
			if (msg == models.PredictedMicroclimateReading{}) {
				fmt.Println("empty struct, chan is closed")
				fmt.Println("remaining buff", buff)
				fmt.Println("remaining sb", sb.String())
				_, err = csvFile.WriteString(sb.String())
				if err != nil {
					fmt.Println("predicted batch microclimate error in writing", err)
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
						fmt.Println("predicted batch microclimate error in writing", err)
						return err
					}
					sb.Reset()
					counterWrite = 0
				}
				counter = 0
				buff = nil
			}

		case <-ctxx.Done():
			fmt.Println("predictedcsv", ctxx.Err())
			fmt.Println("predicted ctx done in generate csv file")
			return err
			//case _ = <-endCh:
			//	fmt.Println("Batch goroutine done")
			//	flag = true
			//	break
		}

	}

}
