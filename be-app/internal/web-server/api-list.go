package webserver

import (
	mockapi "dynamocker/internal/mock-api"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// list of apis
var apis []Api = []Api{
	{
		resource: "mock-apis",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     getMockApis,
			OPTIONS: getOptions,
			DELETE:  deleteMockApis,
		},
	},
	{
		resource: "mock-api/{id}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     getMockApi,
			OPTIONS: getOptions,
			POST:    postMockApi,
			PATCH:   patchMockApi,
			DELETE:  deleteMockApi,
		},
	},
	{
		resource: "serve-mock-api/{mockApiName}",
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
	encodeJson(mockapi.GetAPIs(), w)
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

// GET http://<dynamocker-server>/mock-api/{id}
// get mock api by id
func getMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiName, ok := vars["id"]
	if mockApiName == "" || !ok {
		err := fmt.Errorf("no mockApiName provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	if mockApi, err := mockapi.GetAPI(mockApiName); err != nil {
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	} else {
		encodeJson(mockApi, w)
	}
}

// POST http://<dynamocker-server>/mock-api/{id}
// add mock api
func postMockApi(w http.ResponseWriter, r *http.Request) {
	// retrieve id
	vars := mux.Vars(r)
	mockApiName, ok := vars["id"]
	if mockApiName == "" || !ok {
		err := fmt.Errorf("no mockApiName provided in the URL")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
	}

	// return error if the mockApi already exists. Wrong method used (it should be a patch)
	jsonFiles := readJsonFilesFromFolder()
	if slices.Contains(jsonFiles, mockApiName) {
		err := fmt.Errorf("mockApi with name '" + mockApiName + "' already existing'")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}

	// load body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}

	// add mock api file to the folder
	if err = mockapi.AddNewMockApiFile(mockApiName, body); err != nil {
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH http://<dynamocker-server>/mock-api/{id}
// modify existing mock api
func patchMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiName, ok := vars["id"]
	if mockApiName == "" || !ok {
		err := fmt.Errorf("no mockApiName provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusInternalServerError)
		return
	}

	if err := mockapi.ModifyMockApiFile(mockApiName, body); err != nil {
		err := fmt.Errorf("error while modifying existing mock api: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DEL http://<dynamocker-server>/mock-api/{id}
// delete mock api
func deleteMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiName, ok := vars["id"]
	if mockApiName == "" || !ok {
		err := fmt.Errorf("no mockApiName provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	if err := mockapi.RemoveMockApiFile(mockApiName); err != nil {
		err := fmt.Errorf("error while removing the mocking api: %s", err)
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func serveMockApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mockApiName, ok := vars["mockApiName"]
	if mockApiName == "" || !ok {
		err := fmt.Errorf("no mockApiName provided")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	mockApi, err := mockapi.GetAPI(mockApiName)
	if err != nil {
		err := fmt.Errorf("mockApi not found")
		log.Error(err)
		encodeJsonError(err.Error(), w, http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		if mockApi.Responses.Get == nil {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusBadRequest)
			return
		}
		encodeJson(mockApi.Responses.Get, w)
		return
	case "POST":
		if mockApi.Responses.Post == nil {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusBadRequest)
			return
		}
		encodeJson(mockApi.Responses.Post, w)
		return
	case "PATCH":
		if mockApi.Responses.Patch == nil {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusBadRequest)
			return
		}
		encodeJson(mockApi.Responses.Patch, w)
		return
	case "DELETE":
		if mockApi.Responses.Delete == nil {
			err := fmt.Errorf("requested method not defined for this mockApi")
			log.Error(err)
			encodeJsonError(err.Error(), w, http.StatusBadRequest)
			return
		}
		encodeJson(mockApi.Responses.Delete, w)
		return
	}
}
