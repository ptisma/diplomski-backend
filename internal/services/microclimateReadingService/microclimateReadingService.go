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

func (mrs *MicroclimateReadingService) GetBatchMicroclimateReadings(locationId int, fromDate, toDate time.Time, ch chan models.MicroclimateReading) error {
	var err error

	microclimateReading := models.MicroclimateReading{
		LocationID: uint32(locationId),
		FromDate:   fromDate,
		ToDate:     toDate,
	}
	err = microclimateReading.GetBatchMicroclimateReading(mrs.app, ch)

	return err
}

func (mrs *MicroclimateReadingService) GenerateCSVFile(locationId int, fromDate, toDate time.Time, ch chan models.Message, mainCh chan models.Message, ctx context.Context, cancel context.CancelFunc) error {
	var err error

	csvFileOrig, csvFileAbs, err := utils.CreateTempStageFile("csv*.txt")

	ch <- models.Message{ID: 1, Payload: csvFileAbs}
	mainCh <- models.Message{ID: 1, Payload: csvFileAbs}
	if err != nil {
		cancel()
		return err
	}
	fmt.Println("CSVFile abs path:", csvFileAbs)

	//append
	csvFile, err := os.OpenFile(csvFileAbs, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		cancel()
		return err
	}
	//first row
	_, err = csvFile.WriteString("year,day,radn,maxt,mint,rain,pan,vp,code\n")
	if err != nil {
		cancel()
		return err
	}

	//launch a batch getter in separate goroutine
	g, ctx := errgroup.WithContext(ctx)
	batchCh := make(chan models.MicroclimateReading)
	endCh := make(chan string)
	buff := []models.MicroclimateReading{}
	counter := 0
	flag := false
	g.Go(func() error {
		return mrs.GetBatchMicroclimateReadings(locationId, fromDate, toDate, batchCh)
	})
	for {
		select {
		case msg := <-batchCh:
			counter += 1
			buff = append(buff, msg)
			if counter == 6 {
				rowDate, _ := time.Parse("2006-01-02", buff[0].Date)
				csvRow := fmt.Sprintf("%d,%d,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%d\n", rowDate.Year(), rowDate.YearDay(), buff[0].Value, buff[1].Value, buff[2].Value, buff[3].Value, buff[4].Value, buff[5].Value, 3000)
				csvFile.WriteString(csvRow)
				counter = 0
				buff = nil
			}
		case _ = <-endCh:
			fmt.Println("Batch goroutine done")
			flag = true
			break
		}
		if flag == true {
			break
		}
	}

	err = g.Wait()
	fmt.Println("Batch error group", err)
	err = csvFileOrig.Close()
	err = csvFile.Close()
	return err
}
