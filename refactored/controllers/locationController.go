package controllers

import (
	"apsim-api/refactored/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LocationController struct {
	interfaces.ILocationService
}

func (c *LocationController) GetLocations(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	locations, err := c.GetAllLocations(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching cultures")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&locations)
	w.Write(response)
}
