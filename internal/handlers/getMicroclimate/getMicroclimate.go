package getMicroclimate

import (
	microclimateService2 "apsim-api/internal/services/microclimateService"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetMicroclimate(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		microclimateService := microclimateService2.GetMicroclimateService(app)

		microclimates, err := microclimateService.GetMicroclimates()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching microclimate parameters")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(&microclimates)
		w.Write(response)
		//microclimate := models.Microclimate{}
		//microclimates, err := microclimate.GetAllMicroclimate(app)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	fmt.Fprintf(w, "Error in fetching microclimate parameters")
		//	return
		//}
		//w.Header().Set("Content-Type", "application/json")
		//response, _ := json.Marshal(microclimates)
		//w.Write(response)

	})
}
