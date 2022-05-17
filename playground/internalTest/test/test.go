package test

import (
	"apsim-api/internal/models"
	"fmt"
)

func PrintTest() {

	x := models.Location{
		ID:        0,
		Name:      "",
		Latitude:  0,
		Longitude: 0,
	}

	fmt.Println(x)

}
