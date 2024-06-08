package webserver

import (
	"encoding/json"
	"net/http"
)

// return dynamocer apis
func (ws WebServer) getHandlers() []Api {
	return apis
}

// encode JSONs in the response and return 200
func encodeJson(data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
