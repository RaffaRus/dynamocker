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
	ws.apiList = getHandlers()

	// register handlers
	if err = ws.registerApis(); err != nil {
		return nil, fmt.Errorf("error while registering the APIs: %s", err)
	}

	return &ws, nil
}

func (ws WebServer) Start() error {

	srv := &http.Server{
		Addr:         "127.0.0.1:" + ws.webPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
		Handler:      ws.router,
	}

	return srv.ListenAndServe()
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
