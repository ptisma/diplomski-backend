package services

import (
	"apsim-api/internal/interfaces"
	"apsim-api/internal/models"
	"apsim-api/internal/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"time"
)

type MicroclimateReadingService struct {
	I1 interfaces.IMicroclimateReadingRepository
	I2 interfaces.IPredictedMicroclimateReadingRepository
}

//Get min and max date for the location id
//Since every date has 6 microclimate parameters, we dont "need" to check for microclimate id
func (s *MicroclimateReadingService) GetMicroclimateReadingPeriod(ctx context.Context, microclimateID, locationID int) (models.Period, error) {
	period := models.Period{}
	firstMicroclimateReading, err := s.I1.GetFirstMicroClimateReading(ctx, locationID)
	//firstMicroclimateReading, err := s.I1.GetFirstMicroClimateReadingByID(ctx, microclimateID, locationID)
	if err != nil {
		return models.Period{}, err
	}

	period.Min = firstMicroclimateReading.Date

	lastPredictedMicroclimateReading, err := s.I2.GetLatestMicroClimateReading(ctx, locationID)
	//lastPredictedMicroclimateReading, err := s.I2.GetLatestMicroClimateReadingByID(ctx, microclimateID, locationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//no predicted search latest
			lastMicroclimateReading, err := s.I1.GetLatestMicroClimateReading(ctx, locationID)
			//lastMicroclimateReading, err := s.I1.GetLatestMicroClimateReadingByID(ctx, microclimateID, locationID)
			if err != nil {
				return models.Period{}, err
			}
			period.Max = lastMicroclimateReading.Date
			return period, err
		}
		return period, err
	}

	period.Max = lastPredictedMicroclimateReading.Date

	return period, err

}

//calculate GDD
func (s *MicroclimateReadingService) CalculateGrowingDegreeDay(tmaxReadings []models.MicroclimateReading, tminReadings []models.MicroclimateReading, baseTemp float32) ([]models.GrowingDegreeDay, error) {
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
		return nil, err
	}

	return gdds, err
}

//calculate predicted GDD
//convert predicted microclimate readings to microclimate readings
func (s *MicroclimateReadingService) CalculatePredictedGrowingDegreeDay(tmaxReadings []models.PredictedMicroclimateReading, tminReadings []models.PredictedMicroclimateReading, baseTemp float32) ([]models.GrowingDegreeDay, error) {
	//var err error
	//var gdds = []models.GrowingDegreeDay{}
	//
	//if len(tmaxReadings) == len(tminReadings) {
	//	for i, _ := range tmaxReadings {
	//		date := tmaxReadings[i].Date
	//		gdd := (tmaxReadings[i].Value+tminReadings[i].Value)/2 - float32(baseTemp)
	//
	//		gdds = append(gdds, models.GrowingDegreeDay{
	//			Date:  date,
	//			Value: gdd,
	//		})
	//
	//	}
	//} else {
	//	err = errors.Errorf("tmax and tmin readings are not equal")
	//}
	//
	//return gdds, err

	return s.CalculateGrowingDegreeDay(s.ConvertPredictedMicroclimateReadings(tmaxReadings), s.ConvertPredictedMicroclimateReadings(tminReadings), baseTemp)
}

func (s *MicroclimateReadingService) CreateMicroclimateReading(ctx context.Context, microclimateID, locationID int, date string, value float32) error {
	return s.I1.CreateMicroclimateReading(ctx, microclimateID, locationID, date, value)
}

func (s *MicroclimateReadingService) CreateMicroclimateReadings(ctx context.Context, microclimateReadings []models.MicroclimateReading) error {
	return s.I1.CreateMicroclimateReadings(ctx, microclimateReadings)
}

func (s *MicroclimateReadingService) GetMicroclimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.MicroclimateReading, error) {
	return s.I1.GetMicroClimateReadings(ctx, microclimateID, locationID, fromDate, toDate)
}

func (s *MicroclimateReadingService) GetPredictedMicroclimateReadings(ctx context.Context, microclimateID, locationID int, fromDate, toDate time.Time) ([]models.PredictedMicroclimateReading, error) {
	return s.I2.GetMicroClimateReadings(ctx, microclimateID, locationID, fromDate, toDate)
}

func (s *MicroclimateReadingService) GetLatestMicroClimateReading(ctx context.Context, locationID int) (models.MicroclimateReading, error) {
	//presumption is that all microclimate reading parameters are in database so we check just latest
	return s.I1.GetLatestMicroClimateReading(ctx, locationID)
}

//Runs in goroutine
func (s *MicroclimateReadingService) GetBatchMicroclimateReadings(locationID int, fromDate, toDate time.Time, ch chan models.MicroclimateReading, ctxx context.Context) error {
	var err error
	defer func() {
		log.Println("GetBatchMicroclimateReadings ending")
		//Because err group only logs one first error from whole group
		if err != nil {
			log.Println("GetBatchMicroclimateReadings", err)

		}
	}()
	err = s.I1.GetBatchMicroclimateReading(locationID, fromDate, toDate, ch, ctxx)
	return err
}

//Runs in goroutine
//Fetch readings from DB and send them to chan
func (s *MicroclimateReadingService) GetBatchPredictedMicroclimateReadings(locationID int, fromDate, toDate time.Time, ch chan models.PredictedMicroclimateReading, ctxx context.Context) error {
	var err error
	defer func() {
		log.Println("GetBatchPredictedMicroclimateReadings ending")
		//Because err group only logs one first error from whole group
		if err != nil {
			log.Println("GetBatchPredictedMicroclimateReadings err:", err)

		}
	}()
	err = s.I2.GetBatchMicroclimateReading(locationID, fromDate, toDate, ch, ctxx)
	return err
}

func (s *MicroclimateReadingService) ReceiveFromPredictedBatchAndWrite(csvFile *os.File, batchCh chan models.PredictedMicroclimateReading, ctxx context.Context) error {
	var err error
	defer func() {
		fmt.Println("ReceiveFromPredictedBatchAndWrite ending")
		//Because err group only logs one first error from whole group
		if err != nil {
			fmt.Println("ReceiveFromPredictedBatchAndWrite err:", err)

		}
	}()

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
				//fmt.Println("empty struct, chan is closed")
				//fmt.Println("remaining buff", buff)
				//fmt.Println("remaining sb", sb.String())
				_, err = csvFile.WriteString(sb.String())
				if err != nil {
					//fmt.Println("predicted batch microclimate error in writing", err)
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
						//fmt.Println("predicted batch microclimate error in writing", err)
						return err
					}
					sb.Reset()
					counterWrite = 0
				}
				counter = 0
				buff = nil
			}

		case <-ctxx.Done():
			//fmt.Println("csv", ctxx.Err())
			//fmt.Println("predicted ctx done in generate csv file")
			return ctxx.Err()
			//case _ = <-endCh:
			//	fmt.Println("Batch goroutine done")
			//	flag = true
			//	break
		}
	}
}

//Runs in goroutine
//Receive readings from chan and write to file
func (s *MicroclimateReadingService) ReceiveFromBatchAndWrite(csvFile *os.File, batchCh chan models.MicroclimateReading, ctxx context.Context) error {
	var err error

	defer func() {
		log.Println("ReceiveFromBatchAndWrite ending")
		//Because err group only logs one first error from whole group
		if err != nil {
			log.Println("ReceiveFromBatchAndWrite err:", err)

		}
	}()

	buff := []models.MicroclimateReading{}
	counter := 0
	var sb strings.Builder
	counterWrite := 0
	for {
		select {
		case msg := <-batchCh:
			//fmt.Println("msg:", msg)
			//check if the chan is closed, if true last reading
			if (msg == models.MicroclimateReading{}) {
				//fmt.Println("empty struct, chan is closed")
				//fmt.Println("remaining buff", buff)
				//fmt.Println("remaining sb", sb.String())
				_, err = csvFile.WriteString(sb.String())
				if err != nil {
					//fmt.Println("batch microclimate error in writing", err)
					return err
				}
				return err
			}
			counter += 1
			buff = append(buff, msg)
			//one row in CSV is all microclimate readings from that Date
			if counter == 6 {
				rowDate, err := time.Parse("2006-01-02", buff[0].Date)
				if err != nil {
					return err
				}
				csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
				//_, err = csvFile.WriteString(csvRow)
				sb.WriteString(csvRow)
				counterWrite += 1
				//batch writing with string builder to speed up process
				if counterWrite == 6 {
					_, err = csvFile.WriteString(sb.String())
					if err != nil {
						//fmt.Println("batch microclimate error in writing", err)
						return err
					}
					sb.Reset()
					counterWrite = 0
				}
				counter = 0
				buff = nil
			}

		case <-ctxx.Done():
			//fmt.Println("csv", ctxx.Err())
			//fmt.Println("ctx done in generate csv file")
			return ctxx.Err()
			//case _ = <-endCh:
			//	fmt.Println("Batch goroutine done")
			//	flag = true
			//	break
		}

	}
}

//Generate CSV file in stage area
func (s *MicroclimateReadingService) GenerateCSVFile(locationId int, fromDate, toDate, lastDate time.Time, ch chan utils.Message, mainCh chan utils.Message, ctx context.Context) error {
	var err error
	var csvFileOrig, csvFile *os.File
	defer func() {
		fmt.Println("CSV ending")
		_ = csvFileOrig.Close()

		_ = csvFile.Close()
		//Because err group only logs one error
		if err != nil {
			fmt.Println("GenerateCSVFile err:", err)

		}
	}()

	csvFileOrig, csvFileAbs, err := utils.CreateTempStageFile("csv*.csv")
	if err != nil {
		return err
	}
	ch <- utils.Message{ID: utils.CSV_FILE_CODE, Payload: csvFileAbs}
	mainCh <- utils.Message{ID: utils.CSV_FILE_CODE, Payload: csvFileAbs}

	fmt.Println("CSVFile abs path:", csvFileAbs)

	//append
	csvFile, err = os.OpenFile(csvFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
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
	batchCh := make(chan models.MicroclimateReading, 100)
	//endCh := make(chan string, 1)
	//buff := []models.MicroclimateReading{}
	//counter := 0
	//flag := false
	g.Go(func() error {
		return s.GetBatchMicroclimateReadings(locationId, fromDate, toDate, batchCh, ctxx)
	})
	g.Go(func() error {
		return s.ReceiveFromBatchAndWrite(csvFile, batchCh, ctxx)
	})
	//err = s.ReceiveFromBatchAndWrite(csvFile, batchCh, ctxx)
	//if err != nil {
	//	return err
	//}
	err = g.Wait()
	//fmt.Println("Batch error group", err)
	if err != nil {
		return err
	}
	//PREDICTED
	if err == nil && toDate.After(lastDate) {
		var newFromDate time.Time
		if fromDate.After(lastDate) {
			newFromDate = fromDate
		} else {
			newFromDate = lastDate.AddDate(0, 0, 1)
		}
		g, ctxx := errgroup.WithContext(ctx)
		batchCh := make(chan models.PredictedMicroclimateReading, 100)
		//buff := []models.PredictedMicroclimateReading{}
		//counter := 0
		//flag := false
		g.Go(func() error {
			return s.GetBatchPredictedMicroclimateReadings(locationId, newFromDate, toDate, batchCh, ctxx)
		})
		//fmt.Println("USAO U PREDICTED")
		g.Go(func() error {
			return s.ReceiveFromPredictedBatchAndWrite(csvFile, batchCh, ctxx)
		})
		//err = s.ReceiveFromPredictedBatchAndWrite(csvFile, batchCh, ctxx)
		//if err != nil {
		//	return err
		//}
		err = g.Wait()
		if err != nil {
			return err
		}
	}

	return err

}

//helper function to convert PredictedMicroclimateReading to MicroclimateReading
func (s *MicroclimateReadingService) ConvertPredictedMicroclimateReadings(predictedMicroclimateReadings []models.PredictedMicroclimateReading) []models.MicroclimateReading {
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
		})
	}
	return mr
}

//helper function to validate microclimate readings if they fit the interval
func (s *MicroclimateReadingService) ValidateMicroclimateReadings(fromDate, toDate time.Time, microclimateReadings []models.MicroclimateReading) bool {
	difference := toDate.Sub(fromDate)
	if len(microclimateReadings) != int(difference.Hours()/24)+1 {
		return false
	}
	return true
}

//helper function to validate growing degree days if they fit the interval
func (s *MicroclimateReadingService) ValidateGrowingDegreeDays(fromDate, toDate time.Time, gdds []models.GrowingDegreeDay) bool {
	difference := toDate.Sub(fromDate)

	if len(gdds) != int(difference.Hours()/24)+1 {
		return false
	}
	return true
}
