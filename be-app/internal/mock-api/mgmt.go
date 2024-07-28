package mockapipkg

import (
	"dynamocker/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
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

// TODO: improve the way the mockapi are handled:
// 	- the uniqueness of the uuid is ensured by the map construct whose key is the uuid
// 	- the uniqueness of the name and of the url is checked with a for cycle O(N)*2 --> this must be improved

var mockApiList = make(map[uint16]*MockApi)

func Init(closeAll chan bool, wg *sync.WaitGroup) error {

	mockApiList = make(map[uint16]*MockApi)
	folderPath = config.GetMockApiFolder()

	// load the stored APIs for the first time
	if err := loadAPIsFromFolder(); err != nil {
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

// loading the APIs from the mock api folder at startup
// this function updates the list based on the entried loaded from the folder
func loadAPIsFromFolder() (err error) {

	mu.Lock()
	defer mu.Unlock()

	var files []fs.DirEntry

	// get path from config package
	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	if files, err = os.ReadDir(folderPath); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		// select only *.json files
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		pathToFile := folderPath + "/" + file.Name()

		// open file
		jsonFile, err := os.Open(pathToFile)
		if err != nil {
			log.Errorf("error while opening the file %s: %s", pathToFile, err)
			continue
		}

		// read content
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Errorf("error while reading the file %s: %s", pathToFile, err)
			continue
		}

		// unmarshal content
		var mockApi MockApi
		err = json.Unmarshal(byteValue, &mockApi)
		if err != nil {
			log.Errorf("error while unmarshaling the json file %s into the struct: %s", pathToFile, err)
			continue
		}

		// validate content
		vtor := validator.New(validator.WithRequiredStructEnabled())
		err = vtor.Struct(mockApi)
		if err != nil {
			log.Errorf("invalid mock api saved in the json file %s into the struct: %s", pathToFile, err)
			continue
		}

		addMockApiToMap(mockApi)
	}

	return nil
}

func GetMockAPIs() []*MockApi {
	ret := make([]*MockApi, 0)
	for _, mockApi := range mockApiList {
		ret = append(ret, mockApi)
	}
	return ret
}

func GetMockAPI(uuid uint16) (*MockApi, error) {
	mockApi, found := mockApiList[uuid]
	if !found {
		err := fmt.Errorf("no mockApi with uuid %d found", uuid)
		log.Error(err)
		return nil, err
	}
	return mockApi, nil
}

func GetMockApiList() map[uint16]*MockApi {
	return mockApiList
}

// look for the mockApi whose name mathes the arg passed id. It
// returns the mockApi and true/false if found or not
func GetApiByName(name string) (*MockApi, bool) {
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
func GetApiByUrl(url string) (*MockApi, bool) {
	for _, mockApi := range mockApiList {
		if mockApi.URL == url {
			return mockApi, true
		}
	}
	log.Errorf("no match for mockApi URL '%s'", url)
	return nil, false
}

// pass in the mockApi whose Uuid must be found. It ranges over the
// mockApiList and uses the 'Name' property to find a match with the
// mockApi passed in. It returns the uuid that matches the name. It
// returns a uuid and a boolean, true if found, false otherwise
func GetUuid(mockApiToBeFound *MockApi) (uint16, bool) {
	for uuid, mockApi := range mockApiList {
		if mockApi.Name == mockApiToBeFound.Name {
			return uuid, true
		}
	}
	return 0, false
}

// add mock Api to the list. A random uuid is assigned to the mockapi
func addMockApiToMap(mockApi MockApi) {

	// attempt to retrieve the Uuid
	uuid, found := GetUuid(&mockApi)

	// if not found, assign a UUID
	if !found {
		uuid = generateUuid()
	}

	// add to the mockApi list
	mockApiList[uuid] = &mockApi

	if found {
		log.Info("modified '", mockApi.Name, "' mock API")
	} else {
		log.Info("added '", mockApi.Name, "' mock API")
	}
}

// generate a random uuid
func generateUuid() uint16 {
	var tmp uint16
	var counter uint16 = 0
	for {
		if counter > 100 {
			log.Fatal("Reached maximum numbers of attempts in generating a random uuid")
		}
		tmp = uint16(rand.Intn(_max_size_mockapi_list))
		_, found := mockApiList[tmp]
		if !found {
			break
		}
		counter++
	}
	return tmp
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
				detectedNewMockApi(event.Name)
			}
			// removed api file
			if event.Has(fsnotify.Remove) {
				log.Debug("removed json detected in the folder: ", fileName)
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

	// open file
	jsonFile, err := os.Open(pathToFile)
	if err != nil {
		log.Errorf("error while opening the file %s: %s", pathToFile, err)
		return
	}

	// read content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Errorf("error while reading the file %s: %s", pathToFile, err)
		return
	}

	// unmarshal content
	var mockApi MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		log.Errorf("error while unmarshaling the json file %s into the struct: %s", pathToFile, err)
		return
	}

	// validate structure
	vtor := validator.New(validator.WithRequiredStructEnabled())
	err = vtor.Struct(mockApi)
	if err != nil {
		log.Errorf("invalid mock api saved in the json file %s into the struct: %s", pathToFile, err)
		return
	}

	addMockApiToMap(mockApi)

}

// function called once a json mock api file has been removed from the folder
func detectedRemovedMockApi(pathToFile string) {
	fileName, ok := strings.CutSuffix(path.Base(pathToFile), ".json")
	if !ok {
		log.Error("suffix '.json' not found in the ", pathToFile, " file")
		return
	}
	mockApi, found := GetApiByName(fileName)
	if !found {
		log.Info("mock api named '", fileName, "' not found in the list. Probably already removed it")
		return
	}
	uuid, found := GetUuid(mockApi)
	if !found {
		log.Error("uuid of the mockApi named '", mockApi.Name, "' not found in the list. MockApi not removed from the list")
		return
	}
	delete(mockApiList, uuid)
	if _, ok = mockApiList[uuid]; ok {
		log.Errorf("mock api named %s was not removed", mockApi.Name)
	} else {
		log.Infof("mock api named %s was successfully removed", mockApi.Name)
	}
}
