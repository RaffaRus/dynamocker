package webserver

import (
	"dynamocker/internal/config"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type WebServer struct {
	router  *mux.Router
	version ApiVersion
	webPort string
	apiList []Api
}

func NewServer(apiVersion ApiVersion) (ws *WebServer, err error) {

	// set port
	if webPort, err := config.GetServerPort(); err != nil {
		return nil, fmt.Errorf("error while setting webserver port: %s", err)
	} else {
		ws.webPort = webPort
	}

	// define handlers
	apiList := []Api{}

	ws.router = mux.NewRouter()

	// register handlers
	for _, api := range apiList {
		err := ws.register(api)
		if err != nil {
			return nil, err
		}
	}

	return ws, nil
}

func (s WebServer) Start(version ApiVersion) error {

	var err error

	return http.ListenAndServe(":"+s.webPort, s)
}

func handleMockReq(rw http.ResponseWriter, req *http.Request) {
	return
}

func handleUiReq(rw http.ResponseWriter, req *http.Request) {

}

// register api for each of its versions
func (s WebServer) register(api Api) error {

}
