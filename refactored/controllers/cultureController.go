package controllers

import (
	"apsim-api/refactored/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

type CultureController struct {
	interfaces.ICultureService
}

func (c *CultureController) GetCultures(w http.ResponseWriter, r *http.Request) {

	cultures, err := c.FetchAllCultures()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in fetching cultures")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(cultures)
	w.Write(response)
}
