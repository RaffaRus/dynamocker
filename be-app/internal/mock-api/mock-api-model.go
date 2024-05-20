package mockapi

import (
	"time"
)

// Structure used to model the MockApi.
//
//	Name - name of the file without the path and the json suffix
//	URL - url where this MockApi will be served
//	FilePath - path where the file is stored, without the file name
//	Added - timestamp of creation
//	LastModified - timestamp of last modification
type MockApi struct {
	Name         string
	URL          string    `json:"url" ,validate:"required"`
	FilePath     string    `json:"filePath" validate:"dir,required"`
	Added        time.Time `json:"added" validate:"required"`
	LastModified time.Time `json:"lastModified" validate:"required"`
}
