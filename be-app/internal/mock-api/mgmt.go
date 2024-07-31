package mockapipkg

import (
	"dynamocker/internal/common"
	"dynamocker/internal/config"
	mockapifilepkg "dynamocker/internal/mock-api-file"
	"encoding/json"
	"fmt"
	"io"
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

var folderPath = ""

var mockApiList = make(map[uint16]*common.MockApi)

func Init(closeAll chan bool, wg *sync.WaitGroup) error {

	err := mockapifilepkg.Init()
	if err != nil {
		return err
	}
	mockApiList = make(map[uint16]*common.MockApi)
	folderPath = config.GetMockApiFolder()

	// load the stored APIs for the first time
	mockApiList, err = mockapifilepkg.LoadAPIsFromFolder()
	if err != nil {
		return err
	}

	// periodically poll from the folder
	// safe mechanism to recover from not-working observing goroutine
	wg.Add(1)
	go backUpPollingCycle(closeAll, wg)
	wg.Add(1)
	go observeFolder(closeAll, wg)
	time.Sleep(500 * time.Millisecond) // let goroutines start
	log.Info("mocking-mgmt terminated the initialization phase")
	return nil
}

func GetMockAPIs() []*common.MockApi {
	ret := make([]*common.MockApi, 0)
	for _, mockApi := range mockApiList {
		ret = append(ret, mockApi)
	}
	return ret
}

func GetMockAPI(uuid uint16) (*common.MockApi, error) {
	mockApi, found := mockApiList[uuid]
	if !found {
		err := fmt.Errorf("no mockApi with uuid %d found", uuid)
		log.Error(err)
		return nil, err
	}
	return mockApi, nil
}

func GetMockApiList() map[uint16]*common.MockApi {
	return mockApiList
}

// look for the mockApi whose name mathes the arg passed id. It
// returns the mockApi and true/false if found or not
func GetApiByName(name string) (*common.MockApi, bool) {
	for _, mockApi := range mockApiList {
		if mockApi.Name == name {
			return mockApi, true
		}
	}
	log.Errorf("no match for mockApi name '%s'", name)
	return nil, false
}

// look for the mockApi whose url mathes the arg passed id. It
// returns the mockApi and true/false if found or not
func GetApiByUrl(url string) (*common.MockApi, bool) {
	for _, mockApi := range mockApiList {
		if mockApi.URL == url {
			return mockApi, true
		}
	}
	log.Errorf("no match for mockApi URL '%s'", url)
	return nil, false
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
			if !strings.HasSuffix(fileName, ".json") {
				continue
			}
			// any modification to the api file
			if event.Has(fsnotify.Write) {
				log.Debug("modified json detected in the folder: ", fileName)
				detectedNewMockApi(fileName)
			}
			// removed api file
			if event.Has(fsnotify.Remove) {
				log.Debug("removed json detected in the folder: ", fileName)
				detectedRemovedMockApi(fileName)
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
			if mockApiList, err = mockapifilepkg.LoadAPIsFromFolder(); err != nil {
				log.Error("error while loading the stored APIs: ", err)
			}
		}
	}
}

func detectedNewMockApi(fileName string) {

	// open file
	jsonFile, err := os.Open(folderPath + fileName)
	if err != nil {
		log.Errorf("error while opening the file %s: %s", fileName, err)
		return
	}

	// read content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("error while reading the file %s: %s", fileName, err)
		return
	}

	// unmarshal content
	var mockApi common.MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		log.Errorf("error while unmarshaling the json file %s into the struct: %s", fileName, err)
		return
	}

	// validate structure
	vtor := validator.New(validator.WithRequiredStructEnabled())
	err = vtor.Struct(mockApi)
	if err != nil {
		log.Errorf("invalid mock api saved in the json file %s into the struct: %s", fileName, err)
		return
	}

	// parse uuid into a uint16
	uuidString, found := strings.CutSuffix(fileName, ".json")
	if !found {
		log.Error("suffix '.json' not found")
		return
	}
	mockApiUuid64, err := strconv.ParseUint(uuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid of the mockApi file '%s' into uint16", fileName)
		log.Error(err)
	}
	uuid := uint16(mockApiUuid64)

	// add it to the list
	mockApiList[uuid] = &mockApi

}

// function called once a json mock api file has been removed from the folder
func detectedRemovedMockApi(fileName string) {

	// parse uuid into a uint16
	uuidString, found := strings.CutSuffix(fileName, ".json")
	if !found {
		log.Error("suffix '.json' not found")
		return
	}
	// parse uuid into a uint16
	mockApiUuid64, err := strconv.ParseUint(uuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid of the mockApi file '%s' into uint16", fileName)
		log.Error(err)
	}
	uuid := uint16(mockApiUuid64)

	// search for the mockApi
	mockApi, found := mockApiList[uuid]
	if !found {
		log.Info("mock api named '", fileName, "' not found in the list. Probably already removed it")
		return
	}

	// delete the mockApi from the list
	delete(mockApiList, uuid)

	// check it was removed from the list
	if _, ok := mockApiList[uuid]; ok {
		log.Errorf("mock api named %s was not removed", mockApi.Name)
	} else {
		log.Infof("mock api named %s was successfully removed", mockApi.Name)
	}
}
