package webserver

import (
	mockapi "dynamocker/internal/mock-api"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
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

func (ws WebServer) getHandlers() []Api {
	return apis
}

func encodeJson(data any, w http.ResponseWriter, r *http.Request) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("error marshaling data during %s to %s : %s", r.Method, r.URL, err)
		http.Error(w, "marshaling error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(marshaledData)
}

func getMockApis(w http.ResponseWriter, r *http.Request) {
	encodeJson(mockapi.GetAPIs(), w, r)
}

func deleteMockApis(w http.ResponseWriter, r *http.Request) {
	if err := mockapi.RemoveAllMockApisFiles(); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getMockApi(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	if mockApi, err := mockapi.GetAPI(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		encodeJson(mockApi, w, r)
	}
}

func postMockApi(w http.ResponseWriter, r *http.Request) {
	// TO-DO add control over the id
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("no key provided")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var mockApi mockapi.MockApi
	err = json.Unmarshal(body, &mockApi)
	if err != nil {
		log.Errorf("error while unmarshaling body from post request: %s", err)
	}

	vtor := validator.New(validator.WithRequiredStructEnabled())
	err = vtor.Struct(mockApi)
	if err != nil {
		log.Errorf("invalid mock api passed form post request: %s", err)
	}

	if err = mockapi.AddNewMockApiFile(key, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func patchMockApi(http.ResponseWriter, *http.Request) {

}

func deleteMockApi(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	if key == "" {
		err := fmt.Errorf("no key provided")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := mockapi.RemoveMockApiFile(key); err == nil {
		http.Error(w, "requested mockAPI %s not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
