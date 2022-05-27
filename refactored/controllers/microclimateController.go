package controllers

import (
	"apsim-api/refactored/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

type MicroclimateController struct {
	interfaces.IMicroclimateService
}

func (c *MicroclimateController) GetMicroclimates(w http.ResponseWriter, r *http.Request) {

	microclimates, err := c.GetAllMicroclimates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching microclimate parameters")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimates)
	w.Write(response)
}
