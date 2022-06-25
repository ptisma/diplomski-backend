package controllers

import (
	"apsim-api/internal/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LocationController struct {
	interfaces.ILocationService
}

func (c *LocationController) GetLocations(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	locations, err := c.GetAllLocations(ctx)
	if err != nil || len(locations) == 0 {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching cultures")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&locations)
	w.Write(response)
}
