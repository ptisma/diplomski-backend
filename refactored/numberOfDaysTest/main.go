package main

import (
	"fmt"
	"time"
)

func main() {

	currentDate := time.Date(2022, 5, 23, 0, 0, 0, 0, time.UTC)

	targetDate := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)

	x := targetDate.Sub(currentDate) / (24 * time.Hour)

	fmt.Println(x)
}
