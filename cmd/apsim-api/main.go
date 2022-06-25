package main

import (
	"apsim-api/internal/backgroundContainer"
	"apsim-api/internal/infra/application"
	"apsim-api/internal/infra/server"
	"apsim-api/internal/router"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//Load env vars
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load env vars")
	}

	//load wrapper application struct
	app, err := application.GetApplication()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Set up a conn str and wrapper server struct
	srvConnStr := fmt.Sprintf(":%s", app.GetConfig().GetApiPort())
	srv := server.GetServer().WithAddr(srvConnStr).WithRouter(router.GetMuxRouter(app).InitRouter())

	//background := backgroundContainer.NewBackground(app)

	//Set up a wrapper struct with services
	backgroundWorker := backgroundContainer.NewBackgroundWorker(app)
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	//set up a background container struct which schedules all the background works where the background worker with services is used
	scheduler := backgroundContainer.GetBackgroundContainer(ctx, backgroundWorker)

	//Handle from terminal shutdown
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	//Start HTTP server
	go func() {
		//Start the work
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server listen: %s\n", err)
		}
	}()
	//Run background works
	scheduler.ScheduleBackgroundWorks()
	//Wait terminal shutdown
	sig := <-c
	log.Printf("Caught SIGTERM %b, gracefully shutting down resources\n", sig)
	//Graceful shutdown of background works
	cancel()
	scheduler.Exit()
	//shutdown HTTP server
	ctxx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	//Cleanup work
	defer func() {
		//extra handling here
		//shutdown DB sqlite
		err = app.GetDB().Close()
		log.Println("Shutting down DB, value:", err)
		//shutdown cache InfluxDB
		err = app.GetCache().Close()
		log.Println("Shutting down InfluxDB, value:", err)
		cancel()

	}()
	//Graceful shutdown of HTTP server
	if err := srv.Close(ctxx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
