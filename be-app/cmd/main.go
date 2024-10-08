package main

import (
	"dynamocker/internal/config"
	mockapipkg "dynamocker/internal/mock-api"
	mockapifilepkg "dynamocker/internal/mock-api-file"
	webserver "dynamocker/internal/web-server"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// TODO: complete Tests

func main() {
	log.Info("Hello there, this is DynaMocker")

	closeCh := make(chan bool)
	var wg sync.WaitGroup

	go captureSisCall(closeCh, &wg)

	// capture panics
	defer handlePanic(closeCh)

	// read the customized values of the configuration from the env variables
	config.ReadVars()

	// init the mocked api management
	if err := mockapipkg.Init(closeCh, &wg); err != nil {
		log.Errorf("error initiating mockapi: %s", err)
		panic("panic during mockapi initiations")
	}

	// init the mocked api file management
	if err := mockapifilepkg.Init(); err != nil {
		log.Errorf("error initiating mockapi: %s", err)
		panic("panic during mockapi initiations")
	}

	ws, err := webserver.NewServer()
	if err != nil {
		log.Errorf("error while serving the web server: %s", err)
	}

	ws.Start(closeCh, &wg)

	// attempt exit after success. wait some time for all waiting group to be done.
	// force exit after that
	<-closeCh
	log.Info("dyanmocker successfully stopped.")
	close(closeCh)
	log.Info("waiting for all the goroutines to be closed.")
	wgDone := make(chan bool)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)
waitingCycle:
	for counter := 0; counter < 3; counter++ {
		select {
		case <-wgDone:
			break waitingCycle // leading app exit
		default:
			log.Info("waiting one more sec before the waiting group is done")
			time.Sleep(time.Second)
		}
	}
	os.Exit(0)
}

func handlePanic(ch chan bool) {
	// TODO: stopping the app with a syscall will fill the recovered err with:
	// 		Recovered panic: %!f(runtime.plainError=close of closed channel)
	// try to fix this
	if err := recover(); err != nil {
		log.Fatalf("Recovered panic: %f", err)
	}
	log.Info("dyanmocker stopped after panic.")
	close(ch)
	os.Exit(1)
}

func captureSisCall(closeCh chan bool, wg *sync.WaitGroup) {
	wg.Add(1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	select {
	case signal := <-sigChan:
		log.Infof("received signal from OS: %s. Shutting down.", signal)
		close(closeCh)
		return
	case <-closeCh:
		return
	}

}
