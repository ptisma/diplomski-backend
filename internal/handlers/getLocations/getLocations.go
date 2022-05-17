package getLocations

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLocations(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		location := models.Location{}
		locations, err := location.GetAllLLocations(app)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching locations")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(locations)
		w.Write(response)

	})
}
