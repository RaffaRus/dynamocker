package mockapi

import (
	"time"
)

type MockApi struct {
	name         string
	URL          string    `json:"url" ,validate:"required"`
	FilePath     string    `json:"filePath" validate:"dir,required"`
	Added        time.Time `json:"added" validate:"required"`
	LastModified time.Time `json:"lastModified" validate:"required"`
}
