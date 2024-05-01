package webserver

import (
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

var mu sync.Mutex // single client (UI), multiple requests

// return dynamocer apis
func (ws WebServer) getHandlers() []Api {
	return apis
}

// encode JSONs in the response and return 200
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
