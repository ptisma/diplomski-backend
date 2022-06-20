package controllers

import (
	"apsim-api/internal/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CultureController struct {
	interfaces.ICultureService
}

func (c *CultureController) GetCultures(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

	cultures, err := c.FetchAllCultures(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching cultures")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(cultures)
	w.Write(response)
}
