package controllers

import (
	"apsim-api/refactored/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MicroclimateController struct {
	interfaces.IMicroclimateService
}

func (c *MicroclimateController) GetMicroclimates(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	microclimates, err := c.GetAllMicroclimates(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching microclimate parameters")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(&microclimates)
	w.Write(response)
}
