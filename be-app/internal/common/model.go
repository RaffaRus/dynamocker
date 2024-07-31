package common

const MAX_SIZE_MOCKAPI_LIST = 65535

// Structure used to model the MockApi.
type MockApi struct {

	// name of the file without the path and the json suffixs
	Name string `json:"name" validate:"required"`

	// url where this MockApi will be served
	URL string `json:"url" validate:"required"`

	Responses Response `json:"responses" validate:"required"`
}

type Response struct {
	Get    *map[string]interface{} `json:"get,omitempty"`
	Patch  *map[string]interface{} `json:"patch,omitempty"`
	Post   *map[string]interface{} `json:"post,omitempty"`
	Delete *map[string]interface{} `json:"delete,omitempty"`
}
