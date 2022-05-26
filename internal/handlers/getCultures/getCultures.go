package getCultures

import (
	cultureService2 "apsim-api/internal/services/cultureService"
	"apsim-api/pkg/application"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCultures(app *application.Application) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cultureService := cultureService2.GetCultureService(app)
		cultures, err := cultureService.GetCultures()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error in fetching cultures")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(cultures)
		w.Write(response)

		//culture := models.Culture{}
		//cultures, err := culture.GetAllCultures(app)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	fmt.Fprintf(w, "Error in fetching cultures")
		//	return
		//}
		//w.Header().Set("Content-Type", "application/json")
		//response, _ := json.Marshal(cultures)
		//w.Write(response)

	})
}
