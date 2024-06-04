package webserver

import (
	mockapi "dynamocker/internal/mock-api"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// list of apis
var apis []Api = []Api{
	{
		resource: "mock-apis",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:    getMockApis,
			DELETE: deleteMockApis,
		},
	},
	{
		resource: "mock-api/{id}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:    getMockApi,
			POST:   postMockApi,
			PATCH:  patchMockApi,
			DELETE: deleteMockApi,
		},
	},
}

// GET http://<dynamocker-server>/mock-api
// return mock apis
func getMockApis(w http.ResponseWriter, r *http.Request) {
	encodeJson(mockapi.GetAPIs(), w, r)
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
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("key to be deleted not recovered")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mockApi, err := mockapi.GetAPI(key); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		encodeJson(mockApi, w, r)
	}
}

// POST http://<dynamocker-server>/mock-api/{id}
// add mock api
func postMockApi(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("no key provided")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("error while reading request body: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = mockapi.AddNewMockApiFile(key, body); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PATCH http://<dynamocker-server>/mock-api/{id}
// modify existing mock api
func patchMockApi(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("no key provided")
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

	if err := mockapi.ModifyMockApiFile(key, body); err != nil {
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
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("no key provided")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := mockapi.RemoveMockApiFile(key); err == nil {
		err := fmt.Errorf("error while removing the mocking api: %s", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
