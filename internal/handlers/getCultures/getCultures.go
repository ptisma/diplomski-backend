package getCultures

import (
	"apsim-api/internal/models"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCultures(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		culture := models.Culture{}
		cultures, err := culture.GetAllCultures(app)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching cultures")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(cultures)
		w.Write(response)

	})
}
