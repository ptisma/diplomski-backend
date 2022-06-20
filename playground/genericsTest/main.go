package main

import "apsim-api/internal/models"

type MicroclimateReading interface {
	models.MicroclimateReading | models.PredictedMicroclimateReading
}

func agnosticRec[MicroclimateReading](ch chan MicroclimateReading) {

	for  {
		select {
		 msg := <- ch


		}

	}

}

func main() {

	test1 := models.MicroclimateReading{MicroclimateID: uint32(1)}

	ch := make(chan models.MicroclimateReading, 2)

	agnosticRec(ch)

}
