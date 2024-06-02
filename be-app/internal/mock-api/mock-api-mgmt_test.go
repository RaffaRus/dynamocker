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

const observeFolderWitingTimeMilliseconds = 200

// reset package map and folderPath variable
func reset(t *testing.T) {
	mockApiList = make(map[string]*MockApi)
	folderPath = ""
	assert.Equal(t, 0, len(mockApiList))

}

func dummyMockApi() *MockApi {
	return &MockApi{
		Name:         "dummy-mock-api",
		URL:          "url.com",
		FilePath:     os.TempDir(),
		Added:        time.Now(),
		LastModified: time.Now(),
	}
}

func dummyMockApiArray() []*MockApi {
	var mockApis []*MockApi
	for i := 0; i < 5; i++ {
		mockApis = append(mockApis,
			&MockApi{
				Name:         fmt.Sprintf("dummy-mock-api-%d", i),
				URL:          "url.com",
				FilePath:     os.TempDir(),
				Added:        time.Now(),
				LastModified: time.Now(),
			})
	}
	return mockApis
}

// write a dummy mock api file to the Temp folder. The temp folder
// comes from os package
func writeDummyMockApiFile(t *testing.T) (*os.File, *MockApi) {
	mockApi := dummyMockApi()
	file, err := os.CreateTemp("", "dummy-mock-api*.json")
	if err != nil {
		file.Close()
		t.Fatal(err)
	}
	defer file.Close()
	name, ok := strings.CutSuffix(path.Base(file.Name()), ".json")
	if ok != true {
		file.Close()
		t.Fatal("malformed string modification")
	}
	mockApi.Name = name
	data, err := json.Marshal(mockApi)
	if err != nil {
		file.Close()
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if _, err := file.Write([]byte(data)); err != nil {
		file.Close()
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}
	return file, mockApi
}

// check that the closeChannel works
func TestMockApiInit(t *testing.T) {
	reset(t)

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

	// create support channel to wait for the wg to be done
	wgDone := make(chan bool)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)

	// send close trigger
	close(closeCh)

	// wait for three seconds
	for counter := 0; counter < 3; counter++ {
		select {
		case <-wgDone:
			// in this case the wg has been emptied, everything worked as it is supposed to
			return
		default:
			// Wait for ch being wg to be done
			time.Sleep(time.Second)
		}
	}
	t.Fatal("wg not emptied, something went wrong")
}

func TestLoadStoredAPIs(t *testing.T) {
	reset(t)

	// call function before defining any folder. This should log only
	err := loadAPIsFromFolder()
	assert.Equal(t, fmt.Errorf("the mock API folder has not been set-up"), err)

	// set temp folder as the one contining the mock api files
	folderPath = os.TempDir()
	assert.Nil(t, loadAPIsFromFolder())

	// add mock api
	dummyMockApiFile, _ := writeDummyMockApiFile(t)
	defer func() {
		dummyMockApiFile.Close()
		os.Remove(dummyMockApiFile.Name())
	}()

	// check that the apis have been loaded
	assert.Nil(t, loadAPIsFromFolder(), "the function should return no error")
	mockApiFileName, _ := strings.CutSuffix(path.Base(dummyMockApiFile.Name()), ".json")
	_, ok := mockApiList[mockApiFileName]
	assert.True(t, ok)
}

func TestGetAPIs(t *testing.T) {
	reset(t)
	assert.Empty(t, GetAPIs(), "GetAPIs() should return empty array")

	// add apis to the map and check length
	mockApis := dummyMockApiArray()
	for i := 0; i < 5; i++ {
		mockApiList[mockApis[i].Name] = mockApis[i]
	}
	assert.Equal(t, 5, len(GetAPIs()))

	// remove apis from the map and check it is empty
	reset(t)
	assert.Equal(t, 0, len(GetAPIs()))

}

func TestGetAPI(t *testing.T) {
	reset(t)

	// check not-existing key
	key := "api"
	_, err := GetAPI(key)
	assert.Equal(t, err, fmt.Errorf("requested mockApi %s has not been found", key))

	// add mock api to the map
	mockApi := dummyMockApi()
	mockApiList[mockApi.Name] = mockApi

	// check the get works
	res, err := GetAPI(mockApi.Name)
	assert.Nil(t, err)
	assert.Equal(t, res, mockApi)
}

func TestObserveFolderNotSet(t *testing.T) {
	reset(t)

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup

	// create support channel to wait for the wg to be done
	wgDone := make(chan bool)
	wg.Add(1)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)

	go observeFolder(closeCh, &wg)
	time.Sleep(observeFolderWitingTimeMilliseconds * time.Millisecond)

	// check wgDone has been closed. it has been closed because the
	// folder has not been setup. this logs an error
	_, ok := <-wgDone
	assert.False(t, ok)
}

func TestObserveFolderNotExisting(t *testing.T) {
	reset(t)

	// set mock api folder as a temp folder
	folderPath = "/asdasd"

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup

	// create support channel to wait for the wg to be done
	wgDone := make(chan bool)
	wg.Add(1)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)

	go observeFolder(closeCh, &wg)
	time.Sleep(observeFolderWitingTimeMilliseconds * time.Millisecond)

	// check wgDone has been closed because the folder was not found
	select {
	case <-wgDone:
		break
	default:
		t.Fatal("the observing goroutine should have been stopped")
	}

}

func TestObserveFolderCorrectlyClosing(t *testing.T) {
	reset(t)

	// set mock api folder as a temp folder
	folderPath = os.TempDir()

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup

	// create support channel to wait for the wg to be done
	wgDone := make(chan bool)
	wg.Add(1)
	go func(wgDone chan bool) {
		wg.Wait()
		close(wgDone)
	}(wgDone)

	go observeFolder(closeCh, &wg)
	close(closeCh)
	time.Sleep(200 * time.Millisecond)

	// check wgDone has been closed after the 'close' cmd
	select {
	case <-wgDone:
		break
	default:
		t.Fatal("the observing goroutine should have been stopped")
	}
}

func TestObserveFolderNoJson(t *testing.T) {
	reset(t)

	// set mock api folder as a temp folder
	folderPath = os.TempDir()

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup

	// start observing
	go observeFolder(closeCh, &wg)
	defer close(closeCh)

	time.Sleep(200 * time.Millisecond)

	// write file without the *.json suffix
	notJsonFile, err := os.CreateTemp("", "dummy-mock-api*.tmp")
	if err != nil {
		notJsonFile.Close()
		t.Fatal(err)
	}
	defer os.Remove(notJsonFile.Name())
	data, err := json.Marshal(dummyMockApi())
	if err != nil {
		notJsonFile.Close()
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if _, err := notJsonFile.Write([]byte(data)); err != nil {
		notJsonFile.Close()
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}
	// check the mock api has not been loaded
	assert.Empty(t, mockApiList)
}

func TestObserveFolder(t *testing.T) {
	reset(t)

	// set mock api folder as a temp folder
	folderPath = os.TempDir()

	// make channel and waiting group
	closeCh := make(chan bool)
	var wg sync.WaitGroup

	// start observing
	go observeFolder(closeCh, &wg)
	defer close(closeCh)

	time.Sleep(1000 * time.Millisecond)

	// write proper mock api file
	file, mockApi := writeDummyMockApiFile(t)
	defer func() {
		if _, err := os.Stat(file.Name()); err == nil {
			os.Remove(file.Name())
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// check the mock api has been loaded
	assert.Equal(t, 1, len(mockApiList))
	retrievedMockApi, ok := mockApiList[mockApi.Name]
	assert.True(t, ok)
	assert.Equal(t, mockApi.Name, retrievedMockApi.Name)
	assert.Equal(t, mockApi.URL, retrievedMockApi.URL)
	assert.Equal(t, mockApi.FilePath, retrievedMockApi.FilePath)
	// time cannot be compared using the "==" operator
	assert.True(t, mockApi.Added.Equal(retrievedMockApi.Added))
	assert.True(t, mockApi.LastModified.Equal(retrievedMockApi.LastModified))

	// modify the file
	nowTime := time.Now()
	mockApi.LastModified = nowTime
	mockApi.URL = "newUrl.com"
	data, err := json.Marshal(mockApi)
	if err != nil {
		file.Close()
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if _, err = os.Stat(file.Name()); err != nil {
		t.Fatalf("error while querying for file info: %s", err)
	}
	file, err = os.OpenFile(file.Name(), os.O_WRONLY, 0644)
	if err != nil {
		file.Close()
		t.Fatalf("cannot open file :%s", err)
	}
	if _, err := file.Write([]byte(data)); err != nil {
		file.Close()
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}

	// attend for the modifications to be loaded by the observing gorutine
	time.Sleep(200 * time.Millisecond)

	// check the mock api has been modified
	assert.Equal(t, 1, len(mockApiList))
	retrievedMockApi, ok = mockApiList[mockApi.Name]
	assert.True(t, ok)
	assert.Equal(t, mockApi.Name, retrievedMockApi.Name)
	assert.Equal(t, mockApi.URL, retrievedMockApi.URL)
	assert.Equal(t, mockApi.FilePath, retrievedMockApi.FilePath)
	// time cannot be compared using the "==" operator
	assert.True(t, mockApi.Added.Equal(retrievedMockApi.Added))
	assert.True(t, mockApi.LastModified.Equal(retrievedMockApi.LastModified))

	// remove file
	os.Remove(file.Name())

	// attend for the modifications to be loaded by the observing gorutine
	time.Sleep(200 * time.Millisecond)

	// check the mock api has been removed
	assert.Equal(t, 0, len(mockApiList))
	_, ok = mockApiList[mockApi.Name]
	assert.False(t, ok)

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
	dummyMockApiFile, _ := writeDummyMockApiFile(t)
	defer func() {
		dummyMockApiFile.Close()
		os.Remove(dummyMockApiFile.Name())
	}()
	mockApiFileName, _ := strings.CutSuffix(path.Base(dummyMockApiFile.Name()), ".json")
	mockApiList[mockApiFileName] = dummyMockApi()

	// check that is exists
	assert.Equal(t, 1, len(mockApiList))

	// simulate file removal
	detectedRemovedMockApi(dummyMockApiFile.Name())

	// check that it was removed
	assert.Equal(t, 0, len(mockApiList), "mock api 'dummy-mock-api' not deleted in the map")
}
