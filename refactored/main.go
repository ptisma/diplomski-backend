package main

import (
	"apsim-api/refactored/router"
	"net/http"
)

func main() {

	http.ListenAndServe(":8080", router.MuxRouter().InitRouter())
}
