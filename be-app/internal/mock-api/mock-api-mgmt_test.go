package mockapi

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dummyMockApiPath string = fmt.Sprintf("%s/internal/test-features/dummy_mock_api.json", execFolder())

func execFolder() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

// reset package map and folderPath variable
func resetPackage(t *testing.T) {
	mockApiList = make(map[string]*MockApi)
	folderPath = ""
	assert.Equal(t, 0, len(mockApiList))
}

func dummyMockApi(t *testing.T) *MockApi {
	dummyUrl, err := url.Parse("/dummy-url")
	if err != nil {
		t.Fatalf("error while setting-up the summy URL :%s", err)
	}
	return &MockApi{
		name:         "dummy_mock_api",
		URL:          *dummyUrl,
		FilePath:     "/path/mock_api_folder/dummy_api.json",
		Added:        time.Now(),
		LastModified: time.Now(),
	}

}
func setupDummyMockApiFile(t *testing.T) {
	mockApi := dummyMockApi(t)
	data, err := json.Marshal(mockApi)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if err := os.WriteFile(dummyMockApiPath, data, fs.ModePerm); err != nil {
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}
}

func TestMockApiInit(t *testing.T) {

}

func TestLoadStoredAPIs(t *testing.T) {

	// add mock apis to the map

	// check that the apis have been loaded
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
	resetPackage(t)

	// check path not ending with '.json'. This should only log
	fmt.Println("Expected to LOG the Error: suffix '.json' not found in the /internal/test-features/dummy-api file")
	detectedNewMockApi("/internal/test-features/dummy-api")

	// valid mock api
	detectedNewMockApi(dummyMockApiPath)

	// not unmarshable mock api
	detectedNewMockApi(dummyMockApiPath)

	// not valid mock api
	detectedNewMockApi(dummyMockApiPath)

	// change dummy mock api content (same name) and check the new one replaces the old one
	// this tests the ability to capture/handle modifications to the mock api json file

}

func TestDetectedModifiedMockApi(t *testing.T) {
}

func TestDetectedRemovedMockApi(t *testing.T) {
	resetPackage(t)

	// check path not ending with '.json'. This should only log
	fmt.Println("Expected to LOG the Error: suffix '.json' not found in the /internal/test-features/dummy-api file")
	detectedRemovedMockApi("/internal/test-features/dummy-api")

	// check the not existing mock api in the package map
	fmt.Println("Expected to LOG the INFO: mock api named '/internal/test-features/dummy-api' not found. Probably already removed it")
	detectedRemovedMockApi(dummyMockApiPath)

	// add the mock api to the map
	mockApiList["dummy_mock_api"] = dummyMockApi(t)

	// check that is exists
	assert.Equal(t, 1, len(mockApiList))

	// simulate file removal
	detectedRemovedMockApi(dummyMockApiPath)

	// check that it was removed
	assert.Equal(t, 0, len(mockApiList), "mock api 'dummy_api' not deleted in the map")
}
