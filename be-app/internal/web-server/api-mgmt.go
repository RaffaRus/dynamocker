package webserver

import (
	mockapi "dynamocker/internal/mock-api"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getHandlers() []Api {
	return apis
}

// list of apis
var apis []Api = []Api{
	{
		resource: "mock-apis",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:    GetMockApis,
			DELETE: DeleteMockApis,
		},
	},
	{
		resource: "mock-api/{id}",
		handler: map[Method]func(http.ResponseWriter, *http.Request){
			GET:    GetMockApi,
			POST:   PostMockApi,
			PATCH:  PatchMockApi,
			DELETE: DeleteMockApi,
		},
	},
}

func GetMockApis(w http.ResponseWriter, r *http.Request) {
	encodeJson(mockapi.GetAPIs(), w, r)
}

func DeleteMockApis(http.ResponseWriter, *http.Request) {

}

func GetMockApi(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	if mockApi := mockapi.GetAPI(key); mockApi == nil {
		http.Error(w, "requested mockAPI %s not found", http.StatusNotFound)
		return
	} else {
		encodeJson(mockApi, w, r)
	}
}

func PostMockApi(http.ResponseWriter, *http.Request) {

}

func PatchMockApi(http.ResponseWriter, *http.Request) {

}

func DeleteMockApi(http.ResponseWriter, *http.Request) {

}

func encodeJson(data any, w http.ResponseWriter, r *http.Request) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("error marshaling data during %s to %s : %s", r.Method, r.URL, err)
		http.Error(w, "marshaling error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(marshaledData)

}
