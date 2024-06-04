package mockapi

import (
	"encoding/json"
	"time"
)

// Structure used to model the MockApi.
//
//	TODO: add mocking payload
type MockApi struct {

	// name of the file without the path and the json suffix
	Name string

	// url where this MockApi will be served
	URL string `json:"url" ,validate:"required"`

	// path where the file is stored, without the file name
	FilePath string `json:"filePath" validate:"dir,required"`

	// timestamp of creation
	Added time.Time `json:"added" validate:"required"`

	// timestamp of last modification
	LastModified time.Time `json:"lastModified" validate:"required"`

	Responses Response `json:"responses" validate:"required"`
}

type Response struct {
	Get    *json.RawMessage `json:"get,omitempty" validate:"json"`
	Patch  *json.RawMessage `json:"patch,omitempty" validate:"json"`
	Post   *json.RawMessage `json:"post,omitempty" validate:"json"`
	Delete *json.RawMessage `json:"delete,omitempty" validate:"json"`
}
