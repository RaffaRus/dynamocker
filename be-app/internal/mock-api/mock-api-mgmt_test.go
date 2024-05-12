package mockapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// reset package map and folderPath variable
func reset(t *testing.T) {
	mockApiList = make(map[string]*MockApi)
	folderPath = ""
	assert.Equal(t, 0, len(mockApiList))

}

func dummyMockApi() *MockApi {
	return &MockApi{
		name:         "dummy-mock-api",
		URL:          "url.com",
		FilePath:     os.TempDir(),
		Added:        time.Now(),
		LastModified: time.Now(),
	}

}

// write a dummy mock api file to the Temp folder. The temp folder
// comes from os package
func writeDummyMockApiFile(t *testing.T) *os.File {
	mockApi := dummyMockApi()
	file, err := os.CreateTemp("", "dummy-mock-api*.json")
	if err != nil {
		file.Close()
		t.Fatal(err)
	}
	data, err := json.Marshal(mockApi)
	if err != nil {
		file.Close()
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if _, err := file.Write([]byte(data)); err != nil {
		file.Close()
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}
	return file
}

// check that the closeChannel works
func TestMockApiInit(t *testing.T) {
	reset(t)

	// TODO copontinue w this test

	// set mock api folder as a temp folder
	if err := os.Setenv("DYNA_MOCK_API_FOLDER", t.TempDir()); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

	// set polling time to 1 second to speed-up testing
	if err := os.Setenv("POLLER_INTERVAL", "1"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup
	// start Init
	err := Init(closeCh, &wg)
	assert.Nil(t, err)
	// send close trigger
	closeCh <- true
	// create support channel to wait for the wg to be done
	wgDone := make(chan bool)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)
	// wait for three seconds
	for counter := 0; counter < 3; counter++ {
		select {
		case <-wgDone:
			// in this case the wg has been emptied, everything worked as it is supposed to
			return
		default:
			// Wait for
			time.Sleep(time.Second)
		}
	}
	t.Fatal("wg not emptied, something went wrong")
}

func TestLoadStoredAPIs(t *testing.T) {

	// call function before defining any folder. This should log only
	err := loadAPIsFromFolder()
	assert.Equal(t, fmt.Errorf("the mock API folder has not been set-up"), err)

	// set temp folder as the one contining the mock api files
	folderPath = os.TempDir()
	assert.Nil(t, loadAPIsFromFolder())

	// add mock api
	dummyMockApiFile := writeDummyMockApiFile(t)
	defer dummyMockApiFile.Close()

	// check that the apis have been loaded
	assert.Nil(t, loadAPIsFromFolder(), "the function should return no error")
	mockApiFileName, _ := strings.CutSuffix(path.Base(dummyMockApiFile.Name()), ".json")
	_, ok := mockApiList[mockApiFileName]
	assert.True(t, ok)
}

func TestGetAPIs(t *testing.T) {

	// add mock-apis to the map

	// check number of a apis returned

}

func TestGetAPI(t *testing.T) {

	// check not-existing key

	// add mock api to the map

	// check the get works

}

func TestObserveFolder(t *testing.T) {

}

func TestStopObserving(t *testing.T) {

}

func TestDetectedNewMockApi(t *testing.T) {
	reset(t)

	// check path not ending with '.json'. This should only log
	fmt.Println("Expected to LOG the Error: suffix '.json' not found in the /internal/testdata/dummy-api file")
	detectedNewMockApi("/internal/testdata/dummy-api")

	// valid mock api
	// detectedNewMockApi(dummyMockApiPath)

	// not unmarshable mock api
	// detectedNewMockApi(dummyMockApiPath)

	// not valid mock api
	// detectedNewMockApi(dummyMockApiPath)

	// change dummy mock api content (same name) and check the new one replaces the old one
	// this tests the ability to capture/handle modifications to the mock api json file

}

func TestDetectedModifiedMockApi(t *testing.T) {
}

func TestDetectedRemovedMockApi(t *testing.T) {
	reset(t)

	// check that in case of path not ending with '.json' the error is logger. This should only log
	fmt.Println("Expected to LOG the Error: suffix '.json' not found in the /internal/testdata/dummy-api file")
	detectedRemovedMockApi("/internal/testdata/dummy-api")

	// check that in case of not existing mock api, the error is logged. This should only log
	fmt.Println("Expected to LOG the INFO: mock api named '/internal/testdata/dummy-api' not found. Probably already removed it")
	detectedRemovedMockApi("/dummy-mock-api.json")

	// add the mock api to the map and to the folder
	dummyMockApiFile := writeDummyMockApiFile(t)
	defer dummyMockApiFile.Close()
	mockApiFileName, _ := strings.CutSuffix(path.Base(dummyMockApiFile.Name()), ".json")
	mockApiList[mockApiFileName] = dummyMockApi()

	// check that is exists
	assert.Equal(t, 1, len(mockApiList))

	// simulate file removal
	detectedRemovedMockApi(dummyMockApiFile.Name())

	// check that it was removed
	assert.Equal(t, 0, len(mockApiList), "mock api 'dummy-mock-api' not deleted in the map")
}
