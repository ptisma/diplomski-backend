package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func worker(payload string) {
	fmt.Println("Worker begin")
	task(payload)
	fmt.Println("Worker end")

}

var fileMutex sync.Mutex

func task(payload string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	f, _ := os.Create("test.apsimx")
	_, _ = f.WriteString(payload)
	if payload == "Main" {
		time.Sleep(5 * time.Second)
	}
	f.Close()
}

func main() {
	fmt.Println("Main worker")

	task("Main")

	go worker("worker")

	for {

	}

}
