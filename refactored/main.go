package main

import (
	"apsim-api/refactored/backgroundContainer"
	"apsim-api/refactored/router"
	"net/http"
)

func main() {

	background := backgroundContainer.NewBackground()

	go background.UpdateMicroclimateReadings()

	http.ListenAndServe(":8080", router.MuxRouter().InitRouter())

}
