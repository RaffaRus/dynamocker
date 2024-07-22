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
			OPTIONS: getMockApis,
			DELETE:  deleteMockApis,
		},
	},
	{
		resource: "mock-api/{id}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:     getMockApi,
			OPTIONS: getMockApi,
			POST:    postMockApi,
			PATCH:   patchMockApi,
			DELETE:  deleteMockApi,
		},
	},
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mockApi, err := mockapi.GetAPI(mockApiName); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// return error if the mockApi already exists. Wrong method used (it should be a patch)
	jsonFiles := readJsonFilesFromFolder()
	if slices.Contains(jsonFiles, mockApiName) {
		err := fmt.Errorf("mockApi with name '" + mockApiName + "' already existing'")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// load body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add mock api file to the folder
	if err = mockapi.AddNewMockApiFile(mockApiName, body); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := mockapi.ModifyMockApiFile(mockApiName, body); err != nil {
		err := fmt.Errorf("error while modifying existing mock api: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := mockapi.RemoveMockApiFile(mockApiName); err != nil {
		err := fmt.Errorf("error while removing the mocking api: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
