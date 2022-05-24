package main

import (
	"apsim-api/internal/models"
	"apsim-api/internal/router"
	"apsim-api/pkg/application"
	server "apsim-api/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)
import "github.com/joho/godotenv"

func main() {

	//Load env vars
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed to load env vars")
	}
	//Load application
	app, err := application.GetApplication()
	if err != nil {
		fmt.Println(err.Error())
	}

	//auto migration
	//make tables auto based on the structs from model, if they already exist do nothing
	app.DB.Client.AutoMigrate(models.Location{})
	app.DB.Client.AutoMigrate(models.Microclimate{})
	app.DB.Client.AutoMigrate(models.MicroclimateReading{})
	app.DB.Client.AutoMigrate(models.Culture{})
	app.DB.Client.AutoMigrate(models.PredictedMicroclimateReading{})
	app.DB.Client.AutoMigrate(models.Soil{})

	srv := server.GetServer().WithAddr("localhost:8080").WithRouter(router.GetRouter(app))

	//Handle terminal shutdown
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		//Start the work
		err = srv.Start()
		fmt.Println("Starting up HTTP server", err)
	}()

	go func() {
		//Start the background microclimate parameter fetcher
		//utils.BackgroundUpdate(app)
	}()

	//Wait terminal shutdown
	sig := <-c
	fmt.Printf("Caught SIGTERM %b, gracefully shutting down resources", sig)

	//shutdown DB
	err = app.DB.Close()
	fmt.Println("Shutting down DB value:", err)
	//shutdown Influx
	err = app.Writer.Close()
	fmt.Println("Shutting down Influx value", err)
	//shutdown HTTP server
	err = srv.Close()
	fmt.Println("Shutting down HTTP server", err)

}
