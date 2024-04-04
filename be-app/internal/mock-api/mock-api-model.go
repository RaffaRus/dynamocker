package mockapi

import (
	"net/url"
	"time"
)

type MockApi struct {
	name         string
	URL          url.URL   `json:"url" ,validate:"base64url"`
	FilePath     string    `json:"file_path" validate:"dirpath"`
	Added        time.Time `json:"added" validate:"ltecsfield=InnerStructField.StartDate"`
	LastModified time.Time `json:"last_modified" validate:"ltecsfield=InnerStructField.StartDate"`
}
