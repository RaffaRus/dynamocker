package mockapi

import (
	"dynamocker/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var mockApiList *[]MockApi = nil

var watcher *fsnotify.Watcher = nil

func Init() error {

	if err := loadStoredAPIs(); err != nil {
		return nil
	}
	if err := startObservingFolder(); err != nil {
		return nil
	}
	return nil
}

func loadStoredAPIs() (err error) {

	var path string
	var files []fs.DirEntry
	var mockApis []MockApi

	// get path from config package
	if path, err = config.GetMockApiFolder(); err != nil {
		return fmt.Errorf("error while getting mock api folder: %s", err)
	}

	if files, err = os.ReadDir(path); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		jsonFile, err := os.Open(path + "/" + file.Name())
		if err != nil {
			log.Errorf("error while opening the file %s: %s", path+"/"+file.Name(), err)
			continue
		}

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Errorf("error while reading the file %s: %s", path+"/"+file.Name(), err)
			continue
		}

		var mockApi MockApi
		err = json.Unmarshal(byteValue, &mockApi)
		if err != nil {
			log.Errorf("error while unmarshaling the json file %s into the struct: %s", path+"/"+file.Name(), err)
			continue
		}

		vtor := validator.New(validator.WithRequiredStructEnabled())
		err = vtor.Struct(mockApi)
		if err != nil {
			log.Errorf("invalid mock api saved in the  json file %s into the struct: %s", path+"/"+file.Name(), err)
			continue
		}

		mockApis = append(mockApis, mockApi)

	}

	return updateMockApis(mockApis)
}

func updateMockApis(apis []MockApi) (err error) {
	// TO-DO: compare and update mock apis
	return nil
}

func GetAPIs() *[]MockApi {
	return mockApiList
}

func startObservingFolder() error {
	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stopObserving()
	return nil
}

func stopObserving() error {
	if watcher != nil {
		return watcher.Close()
	}
	return nil
}
