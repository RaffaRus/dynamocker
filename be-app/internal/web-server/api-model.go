package webserver

import "net/http"

type Method string

const (
	GET              Method = http.MethodGet
	POST             Method = http.MethodPost
	PUT              Method = http.MethodPut
	PATCH            Method = http.MethodPatch
	DELETE           Method = http.MethodDelete
	OPTIONS          Method = http.MethodOptions
	MockApiArrayType string = "arrayOfMockApis"
	MockApiType      string = "mockApi"
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

type ResourceObject struct {
	ObjId   uint16 `json:"id"`
	ObjType string `json:"type"`
	ObtData any    `json:"data"`
}
