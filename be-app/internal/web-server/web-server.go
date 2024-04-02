package webserver

import (
	"dynamocker/internal/config"
	"fmt"
	"net/http"
)

func Start() error {

	var webPort string
	var err error

	if webPort, err = config.GetServerPort(); err != nil {
		return fmt.Errorf("error while setting webserver port: %s", err)
	}

	s := http.NewServeMux()
	s.HandleFunc("/dynamocker/api", handleUiReq)
	s.HandleFunc("/", handleMockReq)
	return http.ListenAndServe(":"+webPort, s)
}

func handleMockReq(rw http.ResponseWriter, req *http.Request) {

}

func handleUiReq(rw http.ResponseWriter, req *http.Request) {

}
