package mockapi

import (
	"dynamocker/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var mockApiList map[string]*MockApi = nil

func Init(closeAll chan bool) error {

	// load the stored APIs for the first time
	if err := loadStoredAPIs(); err != nil {
		return nil
	}

	// periodically poll from the folder
	go func(closeAll chan bool) {
		for {
			select {
			case <-closeAll:
				return
			default:
				time.Sleep(time.Minute)
				if err := loadStoredAPIs(); err != nil {
					log.Error("error while loading the stored APIs: ", err)
				}
			}
		}
	}(closeAll)

	go ObserveFolder(closeAll)
	return nil
}

func loadStoredAPIs() (err error) {

	var folderPath string
	var files []fs.DirEntry
	mockApiList = nil

	// get path from config package
	if folderPath, err = config.GetMockApiFolder(); err != nil {
		return fmt.Errorf("error while getting mock api folder: %s", err)
	}

	if files, err = os.ReadDir(folderPath); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		addMockApi(folderPath + "/" + file.Name())

	}

	return nil
}

func updateMockApis(apis []MockApi) (err error) {
	// TO-DO: compare and update mock apis
	return nil
}

func GetAPIs() map[string]*MockApi {
	return mockApiList
}

func ObserveFolder(closeAll chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if path, err := config.GetMockApiFolder(); err != nil {
		log.Error("error getting the mock api folder: ", err, ". no folder watched")
		return
	} else {
		watcher.Add(path)
		log.Info("started watching path ", path)
	}
	if err != nil {
		log.Fatal("could not setup new watcher: ", err)
	}
	defer stopObserving(watcher)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Error("returned not ok from watcher Events")
				return
			}
			fileName := path.Base(event.Name)
			// we are interested in modifications to the *.json files
			if !strings.HasSuffix("*.json", fileName) {
				continue
			}
			// new api
			if event.Has(fsnotify.Create) {
				log.Info("adding mock API ", fileName)
				addMockApi(event.Name)
			}
			// modified api
			if event.Has(fsnotify.Write) {
				log.Info("editing mock API ", fileName)
				editMockApi(event.Name)
			}
			// removed api
			if event.Has(fsnotify.Remove) {
				log.Info("removing mock API ", fileName)
				removeMockApi(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Error("returned not ok from watcher Errors")
				return
			}
			log.Println("error from watcher: ", err)
		case <-closeAll:
			stopObserving(watcher)
		}
	}
}

func stopObserving(watcher *fsnotify.Watcher) {
	if watcher != nil {
		err := watcher.Close()
		if err != nil {
			// Force close the watcher
			watcher = nil
			log.Error("error while closing the watcher: ", err)
		}
		log.Info("stopped watching the mock api folder")
		watcher = nil
	}
}

func addMockApi(pathToFile string) {

	fileName, ok := strings.CutSuffix(path.Base(pathToFile), ".json")
	if !ok {
		log.Error("suffix '.json' not found in the ", pathToFile, " file")
		return
	}

	if _, ok = mockApiList[fileName]; ok {
		log.Info("mock api named '", fileName, "' already present. Replacing the old one with the new one")
	}
	removeMockApi(pathToFile)

	jsonFile, err := os.Open(pathToFile)
	if err != nil {
		log.Errorf("error while opening the file %s: %s", pathToFile, err)
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("error while reading the file %s: %s", pathToFile, err)
	}

	var mockApi MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		log.Errorf("error while unmarshaling the json file %s into the struct: %s", pathToFile, err)
	}

	vtor := validator.New(validator.WithRequiredStructEnabled())
	err = vtor.Struct(mockApi)
	if err != nil {
		log.Errorf("invalid mock api saved in the  json file %s into the struct: %s", pathToFile, err)
	}

	mockApi.name = fileName
	mockApiList[fileName] = &mockApi

	log.Infof("loaded ", fileName, " mock API")

}

func editMockApi(pathToFile string) {
	removeMockApi(pathToFile)
	addMockApi(pathToFile)
}

func removeMockApi(pathToFile string) {
	fileName, ok := strings.CutSuffix(path.Base(pathToFile), ".json")
	if !ok {
		log.Error("suffix '.json' not found in the ", pathToFile, " file")
		return
	}
	_, ok = mockApiList[fileName]
	if !ok {
		log.Error("mock api named '", fileName, "' not found. Probably already removed it")
		return
	}
	delete(mockApiList, fileName)

}
