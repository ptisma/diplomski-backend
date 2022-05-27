package controllers

import (
	"apsim-api/refactored/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationController struct {
	interfaces.ILocationService
}

func (c *LocationController) GetLocations(w http.ResponseWriter, r *http.Request) {

	locations, err := c.GetAllLocations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching cultures")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&locations)
	w.Write(response)
}
