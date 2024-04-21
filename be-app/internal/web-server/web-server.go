package webserver

import (
	"dynamocker/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type WebServer struct {
	router  *mux.Router
	webPort string
	apiList []Api
}

func NewServer() (*WebServer, error) {

	var ws = WebServer{}
	var err error

	// set port
	if webPort, err := config.GetServerPort(); err != nil {
		return nil, fmt.Errorf("error while setting webserver port: %s", err)
	} else {
		ws.webPort = webPort
	}

	ws.router = mux.NewRouter()

	// setup logger for the server
	ws.router.Use(loggingMiddleware)

	// load handlers
	ws.apiList = ws.getHandlers()

	// register handlers
	if err = ws.registerApis(); err != nil {
		return nil, fmt.Errorf("error while registering the APIs: %s", err)
	}

	return &ws, nil
}

func (ws WebServer) Start(closeCh chan bool) {

	srv := &http.Server{
		Addr:         "127.0.0.1:" + ws.webPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
		Handler:      ws.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Debugf("web server closed: %s", err)
		}
	}()

	go monitorAndCloseWebServer(srv, closeCh)
}

// register apis
func (ws WebServer) registerApis() error {
	for _, api := range ws.apiList {
		for method, handler := range api.handler {
			ws.router.HandleFunc("/dynamocker/api/"+api.resource, handler).Methods(string(method))
		}
	}
	return nil
}

// middleware used for logging the incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func monitorAndCloseWebServer(ws *http.Server, closeCh chan bool) {
	select {
	case <-closeCh:
		if err := ws.Close(); err != nil {
			log.Fatalf("error while closing web server: %s", err)
		}
	}
}
