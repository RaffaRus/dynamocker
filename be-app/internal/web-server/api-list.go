package webserver

import (
	mockapi "dynamocker/internal/mock-api"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// list of apis
var apis []Api = []Api{
	{
		resource: "mock-api",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			POST:    postMockApi,
			OPTIONS: getOptions,
		},
	},
	{
		resource: "mock-apis",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     getMockApis,
			OPTIONS: getOptions,
			DELETE:  deleteMockApis,
		},
	},
	{
		resource: "mock-api/{uuid}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     getMockApi,
			OPTIONS: getOptions,
			PUT:     putMockApi,
			DELETE:  deleteMockApi,
		},
	},
	{
		resource: "serve-mock-api/{url}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     serveMockApi,
			OPTIONS: getOptions,
			POST:    serveMockApi,
			PATCH:   serveMockApi,
			DELETE:  serveMockApi,
		},
	},
}

// OPTIONS http://<dynamocker-server>/mock-api
// return mock apis
func getOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent)
}

// GET http://<dynamocker-server>/mock-api
// return mock apis
func getMockApis(w http.ResponseWriter, r *http.Request) {
	mockApis := mockapi.GetMockApiList()
	var resourceObjects []ResourceObject = make([]ResourceObject, 0)
	for uuid, mockApi := range mockApis {
		resourceObjects = append(resourceObjects, ResourceObject{ObjId: uuid, ObjType: MockApiArrayType, ObtData: mockApi})
	}
	encodeJson(resourceObjects, w)
}

// DEL http://<dynamocker-server>/mock-api
// delete all the mock apis
func deleteMockApis(w http.ResponseWriter, r *http.Request) {
	if err := mockapi.RemoveAllMockApisFiles(); err != nil {
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET http://<dynamocker-server>/mock-api/{uuid}
// get mock api by uuid
func getMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiUuidString, ok := vars["uuid"]
	if mockApiUuidString == "" || !ok {
		err := fmt.Errorf("no uuid provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	mockApiUuid64, err := strconv.ParseUint(mockApiUuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid into uint16")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	mockApiUuid := uint16(mockApiUuid64)
	if mockApi, err := mockapi.GetMockAPI(mockApiUuid); err != nil {
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	} else {
		encodeJson(ResourceObject{ObjId: mockApiUuid, ObjType: MockApiType, ObtData: mockApi}, w)
	}
}

// PUT http://<dynamocker-server>/mock-api
// add mock api
func postMockApi(w http.ResponseWriter, r *http.Request) {

	// load body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}

	// add mock api file to the folder
	if err = mockapi.AddNewMockApiFile(body); err != nil {
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PUT http://<dynamocker-server>/mock-api/{uuid}
// modify existing mock api
func putMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiUuidString, ok := vars["uuid"]
	if mockApiUuidString == "" || !ok {
		err := fmt.Errorf("no uuid provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}

	mockApiUuid64, err := strconv.ParseUint(mockApiUuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid into uint16")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	mockApiUuid := uint16(mockApiUuid64)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}

	if err := mockapi.ModifyMockApiFile(mockApiUuid, body); err != nil {
		err := fmt.Errorf("error while modifying existing mock api: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DEL http://<dynamocker-server>/mock-api/{uuid}
// delete mock api
func deleteMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiUuidString, ok := vars["uuid"]
	if mockApiUuidString == "" || !ok {
		err := fmt.Errorf("no uuid provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}

	mockApiUuid64, err := strconv.ParseUint(mockApiUuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid into uint16")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	mockApiUuid := uint16(mockApiUuid64)

	if err := mockapi.RemoveMockApiFile(mockApiUuid); err != nil {
		err := fmt.Errorf("error while removing the mocking api: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func serveMockApi(w http.ResponseWriter, r *http.Request) {
	// retrieve the url
	vars := mux.Vars(r)
	mockApiUrl, ok := vars["url"]
	if mockApiUrl == "" || !ok {
		err := fmt.Errorf("no mockApiName provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}

	// find the mockApi mathching the url
	mockApi, found := mockapi.GetApiByUrl(mockApiUrl)
	if !found {
		err := fmt.Errorf("mockApi not found")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		if len(*mockApi.Responses.Get) == 0 {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusNotFound)
			return
		}
		encodeJson(mockApi.Responses.Get, w)
		return
	case "POST":
		if len(*mockApi.Responses.Post) == 0 {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusNotFound)
			return
		}
		encodeJson(mockApi.Responses.Post, w)
		return
	case "PATCH":
		if len(*mockApi.Responses.Patch) == 0 {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusNotFound)
			return
		}
		encodeJson(mockApi.Responses.Patch, w)
		return
	case "DELETE":
		if len(*mockApi.Responses.Delete) == 0 {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusNotFound)
			return
		}
		encodeJson(mockApi.Responses.Delete, w)
		return
	}
}
