package main

import (
	"dynamocker/internal/config"
	mockapi "dynamocker/internal/mock-api"
	webserver "dynamocker/internal/web-server"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Print("Hello there, this is DynaMocker")

	closeCh := make(chan bool)

	go captureSisCall(closeCh)

	// capture panics
	defer handlePanic(closeCh)

	// read the customized values of the configuration from the env variables
	config.ReadVars()

	// init the mocked api management
	if err := mockapi.Init(closeCh); err != nil {
		log.Errorf("error while reading the existing APIs: %s", err)
	}

	ws, err := webserver.NewServer()
	if err != nil {
		log.Errorf("error while serving the web server: %s", err)
	}

	ws.Start(closeCh)

	// exit after success
	log.Info("Dyanmocker successfully stopped.")
	closeCh <- true
	os.Exit(0)
}

func handlePanic(ch chan bool) {
	if err := recover(); err != nil {
		log.Fatalf("Recovered panic: %f", err)
	}
	log.Info("Dyanmocker stopped after panic.")
	ch <- true
	os.Exit(1)
}

func captureSisCall(closeCh chan bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	select {
	case signal := <-sigChan:
		log.Infof("received signal from OS: %s. Shutting down.", signal)
		closeCh <- true
		return
	case <-closeCh:
		return
	}

}
