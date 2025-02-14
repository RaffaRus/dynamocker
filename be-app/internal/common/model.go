package common

import (
	"encoding/json"
	"fmt"
)

const MAX_SIZE_MOCKAPI_LIST = 65535

// Structure used to model the MockApi.
type MockApi struct {

	// name of the file without the path and the json suffixs
	Name string `json:"name" validate:"required"`

	// url where this MockApi will be served
	URL string `json:"url" validate:"required"`

	Responses Response `json:"responses" validate:"required"`
}

// TODO: investigate why the frontend shows methods that are not defined
// This behavios initaite once I modified the model of the Response
type Response struct {
	Get    *ResponseStruct `json:"get,omitempty"`
	Patch  *ResponseStruct `json:"patch,omitempty"`
	Put    *ResponseStruct `json:"put,omitempty"`
	Post   *ResponseStruct `json:"post,omitempty"`
	Delete *ResponseStruct `json:"delete,omitempty"`
}

type ResponseStruct struct {
	Object *map[string]interface{}
	Array  *[]*map[string]interface{}
}

func (r *ResponseStruct) UnmarshalJSON(data []byte) error {

	var errArr, errObj error

	// try unmashal as an object
	errObj = json.Unmarshal(data, &r.Object)

	// try unmarshal as an array
	errArr = json.Unmarshal(data, &r.Array)

	// return error if both failed
	if errArr != nil && errObj != nil {
		return fmt.Errorf("error unmarshaling the ResponseStruct. Could not Unmarshal into an Array nor Object. errArr: '%s'. errObj: '%s'", errArr, errObj)
	}

	// return error in case there exist both
	if errArr == nil && errObj == nil {
		return fmt.Errorf("error unmarshaling the ResponseStruct. Found both array and object version of the data")
	}

	return nil
}

func (r *ResponseStruct) MarshalJSON() ([]byte, error) {

	var errArr, errObj error
	var dataArr, dataObj []byte

	// try mashal as an object
	dataObj, errObj = json.Marshal(r.Object)

	// try marshal as an array
	dataArr, errArr = json.Marshal(r.Array)

	// return error if both failed
	if string(dataObj) == "null" && string(dataArr) == "null" {
		return nil, fmt.Errorf("error marshaling the ResponseStruct. Both Array and Object failed marshaling. errArr: '%s'. errObj: '%s'", errArr, errObj)
	}

	// return error in case there exist both
	if string(dataObj) != "null" && string(dataArr) != "null" {
		return nil, fmt.Errorf("error marshaling the ResponseStruct. The variable holds both array and object data")
	}

	if string(dataObj) != "null" {
		return dataObj, nil
	} else if string(dataArr) != "null" {
		return dataArr, nil
	}

	return nil, fmt.Errorf("something went wrong marshaling the struct ResponseStruct. errArr: '%s'. errObj: '%s'", errArr, errObj)
}
