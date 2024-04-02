package main

import (
	"dynamocker/internal/config"
	mockapi "dynamocker/internal/mock-api"
	webserver "dynamocker/internal/web-server"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Print("Hello there, this is DynaMocker")

	// capture panics
	defer handlePanic()

	// read the customized values of the configuration from the env variables
	config.ReadVars()

	// init the mocked api management
	if err := mockapi.Init(); err != nil {
		log.Errorf("error while reading the existing APIs: %s", err)
	}

	// start server
	if err := webserver.Start(); err != nil {
		log.Errorf("error while serving the web server: %s", err)
	}

	// exit after success
	log.Info("Dyanmocker successfully stopped.")
	os.Exit(0)
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Fatalf("Recovered panic: %f", err)
	}
	log.Info("Dyanmocker stopped after panic.")
	os.Exit(1)
}
