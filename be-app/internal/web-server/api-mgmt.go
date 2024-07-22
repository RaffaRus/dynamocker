package webserver

import (
	"dynamocker/internal/config"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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

func readJsonFilesFromFolder() []string {

	array := []string{}

	var files []fs.DirEntry

	folderPath := config.GetMockApiFolder()

	// get path from config package
	if folderPath == "" {
		log.Errorf("could not retrieve stored files: the mock API folder has not been set-up")
	}

	files, err := os.ReadDir(folderPath)

	if err != nil {
		log.Errorf("could not retrieve stored files: %s", err)
	}

	for _, file := range files {

		// select only *.json files
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		fileName, ok := strings.CutSuffix(file.Name(), ".json")
		if !ok {
			log.Errorf("could not retrieve stored files: could not CutSuffix")
		}

		array = append(array, fileName)
	}
	return array
}
