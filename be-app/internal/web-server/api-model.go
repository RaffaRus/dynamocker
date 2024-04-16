package webserver

import "net/http"

type Method string

const (
	GET    Method = http.MethodGet
	POST   Method = http.MethodPost
	PATCH  Method = http.MethodPatch
	DELETE Method = http.MethodDelete
)

type ApiInterface interface {
	checkVersion() error
	perform()
}

type Api struct {
	// url of the resource
	resource string
	// map between methods and function handler
	handler map[Method]func(http.ResponseWriter, *http.Request)
}
