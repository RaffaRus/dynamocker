package mockapi

import (
	"time"
)

// Structure used to model the MockApi.
//
//	TODO: add mocking payload
type MockApi struct {

	// name of the file without the path and the json suffix
	Name string `json:"name" validate:"required"`

	// url where this MockApi will be served
	URL string `json:"url" validate:"required"`

	// path where the file is stored, without the file name
	FilePath string `json:"filePath" validate:"dir,required"`

	// timestamp of creation
	Added time.Time `json:"added" validate:"required"`

	// timestamp of last modification
	LastModified time.Time `json:"lastModified" validate:"required"`

	Responses Response `json:"responses" validate:"required"`
}

type Response struct {
	Get    *map[string]interface{} `json:"get,omitempty" validate:"omitempty,json"`
	Patch  *map[string]interface{} `json:"patch,omitempty" validate:"omitempty,json"`
	Post   *map[string]interface{} `json:"post,omitempty" validate:"omitempty,json"`
	Delete *map[string]interface{} `json:"delete,omitempty" validate:"omitempty,json"`
}
