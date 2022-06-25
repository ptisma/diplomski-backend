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

type MicroclimateController struct {
	interfaces.IMicroclimateService
}

func (c *MicroclimateController) GetMicroclimates(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	microclimates, err := c.GetAllMicroclimates(ctx)
	if err != nil || len(microclimates) == 0 {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching microclimate parameters")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimates)
	w.Write(response)
}
