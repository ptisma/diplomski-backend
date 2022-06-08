package main

import (
	"apsim-api/refactored/router"
	"net/http"
)

func main() {

	//background := backgroundContainer.NewBackground()
	//TODO SOME CSV EXTERNAL API NOT CORRECT MISSING SOME FIRST ROWS
	//go background.UpdateMicroclimateReadings()

	http.ListenAndServe(":8080", router.MuxRouter().InitRouter())

}
