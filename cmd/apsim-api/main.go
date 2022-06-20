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
		fmt.Println("failed to load env vars")
	}

	app, err := application.GetApplication()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//srvConnStr := fmt.Sprintf("%s:%s", app.GetConfig().GetApiURL(), app.GetConfig().GetApiPort())
	srvConnStr := fmt.Sprintf(":%s", app.GetConfig().GetApiPort())
	srv := server.GetServer().WithAddr(srvConnStr).WithRouter(router.GetMuxRouter(app).InitRouter())

	//background := backgroundContainer.NewBackground(app)

	//new
	backgroundWorker := backgroundContainer.NewBackgroundWorker(app)
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := backgroundContainer.NewScheduler(ctx, backgroundWorker)

	//Handle terminal shutdown
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		//Start the work
		//err = srv.Start()
		//fmt.Println("HTTP server", err)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server listen: %s\n", err)
		}
	}()
	//scheduler.ScheduleBackgroundWorks()
	//Wait terminal shutdown
	sig := <-c
	fmt.Printf("Caught SIGTERM %b, gracefully shutting down resources", sig)
	cancel()
	scheduler.Exit()
	//shutdown HTTP server
	ctxx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer func() {
		//extra handling here
		//shutdown DB sqlite
		err = app.GetDB().Close()
		fmt.Println("Shutting down DB, value:", err)
		//shutdown cache InfluxDB
		err = app.GetCache().Close()
		fmt.Println("Shutting down InfluxDB, value:", err)
		cancel()

	}()
	if err := srv.Close(ctxx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
