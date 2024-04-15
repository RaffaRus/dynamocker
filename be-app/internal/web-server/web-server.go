package webserver

import (
	"dynamocker/internal/config"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type WebServer struct {
	router  *mux.Router
	version uint16
	webPort string
	apiList []Api
}

func NewServer() (ws *WebServer, err error) {

	// set api version of the server
	if version, err := config.GetApiVersion(); err != nil {
		return nil, fmt.Errorf("error while setting webserver port: %s", err)
	} else {
		ws.version = version
	}

	// set port
	if webPort, err := config.GetServerPort(); err != nil {
		return nil, fmt.Errorf("error while setting webserver port: %s", err)
	} else {
		ws.webPort = webPort
	}

	// define handlers
	apiList := []Api{
		{resource: "/mock-api", versions: []uint16{1, 2}},
	}

	ws.router = mux.NewRouter()

	// register handlers
	for _, api := range apiList {
		ws.register(api)
	}

	return ws, nil
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

func handleMockReq(rw http.ResponseWriter, req *http.Request) {
	return
}

func handleUiReq(rw http.ResponseWriter, req *http.Request) {

}

// register api for each of its versions
func (ws WebServer) register(api Api) {
	for _, ver := range api.versions {
		ws.router.HandleFunc("/dynamocker/api/"+strconv.Itoa(int(ver))+api.resource, api.handler[ver])
	}
}
