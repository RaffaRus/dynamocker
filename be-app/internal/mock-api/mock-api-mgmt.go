package mockapi

import (
	"dynamocker/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type MockMgmt struct {
	mockApiList map[string]*MockApi
	folderPath  string
	mutex       *sync.Mutex
}

func Init(closeAll chan bool, wg *sync.WaitGroup) (MockMgmt, error) {

	mockMgmt := MockMgmt{
		mockApiList: make(map[string]*MockApi),
		folderPath:  config.GetMockApiFolder(),
		mutex:       &sync.Mutex{},
	}

	// load the stored APIs for the first time
	if err := mockMgmt.loadAPIsFromFolder(); err != nil {
		return MockMgmt{}, err
	}

	// periodically poll from the folder
	// safe mechanism to recover from not-working observing goroutine
	wg.Add(1)
	go backUpPollingCycle(closeAll, wg)
	wg.Add(1)
	go observeFolder(closeAll, wg)
	time.Sleep(500 * time.Millisecond) // let goroutines start
	return nil
}

// loading the APIs from the mock api folder at startup
// this function empties the mock api map and creates a new one
func (m MockMgmt) loadAPIsFromFolder() (err error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var files []fs.DirEntry

	// get path from config package
	if m.folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	if files, err = os.ReadDir(m.folderPath); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		detectedNewMockApi(m.folderPath + "/" + file.Name())

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

func observeFolder(closeAll chan bool, wg *sync.WaitGroup) {
	defer func() {
		log.Debug("wg.Done observeFolder")
		wg.Done()
	}()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("could not setup new watcher: ", err)
		return
	}
	if folderPath == "" {
		log.Error("the mock API folder has not been set-up")
		return
	} else {
		err := watcher.Add(folderPath)
		if err != nil {
			log.Error("could add folder to the watcher: ", err)
			return
		}
	}
	log.Info("started watching path ", folderPath)
	defer stopObserving(watcher)
detectingCycle:
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Error("returned not ok from watcher Events")
				return
			}
			fileName := path.Base(event.Name)
			// we are interested in modifications to the *.json files
			if !strings.HasSuffix(".json", fileName) {
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
				log.Errorf("returned not ok from watcher Errors. err: %s", err)
				return
			}
			log.Println("error from watcher: ", err)
		case <-closeAll:
			log.Infof("received signal to close the folder observation")
			break detectingCycle
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

func backUpPollingCycle(closeAll chan bool, wg *sync.WaitGroup) {
	defer func() {
		log.Debug("wg.Done backUpPollingCycle")
		wg.Done()
	}()
	pollerInterval, err := strconv.Atoi(config.GetPollingInterval())
	if err != nil {
		log.Error("could not convert poller interval to a number: ", err, ". Safe-polling will not be working")
		return
	}
pollingCycle:
	for {
		select {
		case <-closeAll:
			log.Infof("received signal to close the backup polling cycle")
			break pollingCycle
		default:
			time.Sleep(time.Duration(pollerInterval) * time.Second) // poll each 'config.GetPollingInterval()' seconds
			if err := loadAPIsFromFolder(); err != nil {
				log.Error("error while loading the stored APIs: ", err)
			}
		}
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

	mockApi.Name = fileName
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
