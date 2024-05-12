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

var mockApiList map[string]*MockApi = make(map[string]*MockApi)
var folderPath string = ""

func Init(closeAll chan bool) error {

	folderPath = config.GetMockApiFolder()

	// load the stored APIs for the first time
	if err := loadAPIsFromFolder(); err != nil {
		return nil
	}

	// periodically poll from the folder
	// safe mechanism to recover from not-working observing goroutine
	go func(closeAll chan bool) {
		for {
			select {
			case <-closeAll:
				return
			default:
				time.Sleep(time.Minute)
				if err := loadAPIsFromFolder(); err != nil {
					log.Error("error while loading the stored APIs: ", err)
				}
			}
		}
	}(closeAll)

	go ObserveFolder(closeAll)
	return nil
}

// loading the APIs from the mock api folder at startup
// this function empties the mock api map and creates a new one
func loadAPIsFromFolder() (err error) {

	var files []fs.DirEntry
	mockApiList = make(map[string]*MockApi)

	// get path from config package
	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	if files, err = os.ReadDir(folderPath); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		detectedNewMockApi(folderPath + "/" + file.Name())

	}

	return nil
}

func GetAPIs() []*MockApi {
	ret := make([]*MockApi, 0)
	for _, mockApi := range mockApiList {
		ret = append(ret, mockApi)
	}
	return ret
}

func GetAPI(key string) (*MockApi, error) {
	mockApi, ok := mockApiList[key]
	if !ok {
		err := fmt.Errorf("requested mockApi %s has not been found", key)
		log.Error(err)
		return nil, err
	}
	return mockApi, nil
}

func ObserveFolder(closeAll chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if folderPath == "" {
		log.Error("the mock API folder has not been set-up")
		return
	} else {
		watcher.Add(folderPath)
		log.Info("started watching path ", folderPath)
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
				log.Info("new json detected in the folder: ", fileName)
				detectedNewMockApi(event.Name)
			}
			// modified api
			if event.Has(fsnotify.Write) {
				log.Info("modified json detected in the folder: ", fileName)
				detectedNewMockApi(event.Name)
			}
			// removed api
			if event.Has(fsnotify.Remove) {
				log.Info("removed json detected in the folder: ", fileName)
				detectedRemovedMockApi(event.Name)
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

func detectedNewMockApi(pathToFile string) {

	fileName, ok := strings.CutSuffix(path.Base(pathToFile), ".json")
	if !ok {
		log.Error("suffix '.json' not found in the ", pathToFile, " file")
		return
	}

	if _, ok = mockApiList[fileName]; ok {
		log.Info("mock api named '", fileName, "' already present. Replacing the old one with the new one")
		detectedRemovedMockApi(pathToFile)
	}

	jsonFile, err := os.Open(pathToFile)
	if err != nil {
		log.Errorf("error while opening the file %s: %s", pathToFile, err)
		return
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("error while reading the file %s: %s", pathToFile, err)
		return
	}

	var mockApi MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		log.Errorf("error while unmarshaling the json file %s into the struct: %s", pathToFile, err)
		return
	}

	vtor := validator.New(validator.WithRequiredStructEnabled())
	err = vtor.Struct(mockApi)
	if err != nil {
		log.Errorf("invalid mock api saved in the json file %s into the struct: %s", pathToFile, err)
		return
	}

	mockApi.name = fileName
	mockApiList[fileName] = &mockApi

	log.Info("loaded ", fileName, " mock API")

}

// function called once a json mock api file has been removed from the folder
func detectedRemovedMockApi(pathToFile string) {
	fileName, ok := strings.CutSuffix(path.Base(pathToFile), ".json")
	if !ok {
		log.Error("suffix '.json' not found in the ", pathToFile, " file")
		return
	}
	_, ok = mockApiList[fileName]
	if !ok {
		log.Info("mock api named '", fileName, "' not found. Probably already removed it")
		return
	}
	delete(mockApiList, fileName)
	if _, ok = mockApiList[fileName]; ok {
		log.Errorf("mock api named %s was not removed", fileName)
	} else {
		log.Infof("mock api named %s was successfully removed", fileName)
	}
}
