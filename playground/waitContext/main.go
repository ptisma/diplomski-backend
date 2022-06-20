package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("UpdateMicroclimateReadings closing")
		return

	case <-time.After(5 * time.Second):
		fmt.Println("Starting updating again")
	}
}
