package webserver

import (
	"bytes"
	"dynamocker/internal/common"
	mockapipkg "dynamocker/internal/mock-api"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (chan bool, *WebServer) {

	// set folderPath
	if err := os.Setenv("DYNA_MOCK_API_FOLDER", os.TempDir()+"/"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

	// init the mocked api management
	closeCh := make(chan bool)
	var wg sync.WaitGroup
	if err := mockapipkg.Init(closeCh, &wg); err != nil {
		t.Errorf("error initiating mockapi: %s", err)
		panic("panic during mockapi initiations")
	}

	// start web server
	webServerTest, err := NewServer()
	if err != nil {
		t.Fatal(t, "error while initiating the test web server")
	}
	err = webServerTest.registerApis()
	if err != nil {
		t.Fatal(t, "error while registering the APIs of the the test web server")
	}
	return closeCh, webServerTest
}

func dummyMockApi(t *testing.T) common.MockApi {
	var response common.Response
	if json.Unmarshal([]byte(`{"valid_json":true,"body":"this is the response"}`), &response.Get) != nil {
		t.Fatal("error while unmashaling")
	}
	if json.Unmarshal([]byte(`{"example_patch_body":"this is a string returned from patch operation"}`), &response.Patch) != nil {
		t.Fatal("error while unmashaling")
	}
	if json.Unmarshal([]byte(`{"error":"posted an invalid element"}`), &response.Post) != nil {
		t.Fatal("error while unmashaling")
	}
	if json.Unmarshal([]byte(`{"response":"removed the item number 3"}`), &response.Delete) != nil {
		t.Fatal("error while unmashaling")
	}
	return common.MockApi{
		Name:      fmt.Sprintf("dummy-mock-api-%d", rand.Intn(1000)),
		URL:       fmt.Sprintf("url-%d.com", rand.Intn(1000)),
		Responses: response,
	}
}

// write a dummy mock api file to the Temp folder. The temp folder
// comes from os package
func writeDummyMockApiFile(t *testing.T) (uint16, *os.File, common.MockApi) {
	mockApi := dummyMockApi(t)
	uuid := uint16(rand.Intn(1000))
	filename := fmt.Sprintf("%d", uuid) + ".json"
	filePath := os.TempDir() + "/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		file.Close()
		t.Fatal(err)
	}
	defer file.Close()
	_, ok := strings.CutSuffix(filename, ".json")
	if !ok {
		file.Close()
		t.Fatal("malformed string modification")
	}
	mockApi.Name = fmt.Sprintf("dummy-mock-api-%d", uuid)
	data, err := json.Marshal(mockApi)
	if err != nil {
		file.Close()
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	if _, err := file.Write([]byte(data)); err != nil {
		file.Close()
		t.Fatalf("error while writing dummy mock api to file :%s", err)
	}
	return uuid, file, mockApi
}

func TestGetMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)

	// write dummy mock Api
	uuid, _, mockApi := writeDummyMockApiFile(t)
	defer func() {
		closeCh <- true
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	// test get/{uuid} api
	r := httptest.NewRecorder()
	url := "/dynamocker/api/mock-api/" + fmt.Sprint(uuid)
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("GET", url, nil))
	assert.Equal(t, http.StatusOK, r.Code)
	resObj := ResourceObject{
		ObjId:   uuid,
		ObjType: MockApiType,
		ObtData: mockApi,
	}
	resObjBytes, err := json.Marshal(resObj)
	if err != nil {
		t.Fatalf("error while marshaling mockApi: %s", err)
	}
	bytesResp, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error while reading from response body: %s", err)
	}
	assert.Equal(t, append(resObjBytes, []byte("\n")...), bytesResp)
}

func TestGetMockApis(t *testing.T) {

	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)
	defer func() { closeCh <- true }()
	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	uuid1, _, _ := writeDummyMockApiFile(t)
	time.Sleep(50 * time.Millisecond)
	uuid2, _, _ := writeDummyMockApiFile(t)
	time.Sleep(50 * time.Millisecond)
	uuid3, _, _ := writeDummyMockApiFile(t)

	defer func() {
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid1)
		removeMockApiFile(t, uuid2)
		removeMockApiFile(t, uuid3)
	}()

	// wait
	time.Sleep(100 * time.Millisecond)

	// test get api
	url := "/dynamocker/api/mock-apis"
	r := httptest.NewRecorder()
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("GET", url, nil))
	assert.Equal(t, http.StatusOK, r.Code)
	bytesResp, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error while reading from response body: %s", err)
	}
	var resObjResp []ResourceObject
	err = json.Unmarshal(bytesResp, &resObjResp)
	if err != nil {
		t.Fatalf("error while unmatshalling the response: %s", err)
	}
	assert.Equal(t, 3, len(resObjResp))
}

func TestDeleteMockApis(t *testing.T) {

	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)
	defer func() { closeCh <- true }()

	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	uuid1, _, _ := writeDummyMockApiFile(t)
	uuid2, _, _ := writeDummyMockApiFile(t)
	uuid3, _, _ := writeDummyMockApiFile(t)

	defer func() {
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid1)
		removeMockApiFile(t, uuid2)
		removeMockApiFile(t, uuid3)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	// test delete api
	url := "/dynamocker/api/mock-apis"
	r := httptest.NewRecorder()
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("DELETE", url, nil))
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(50 * time.Millisecond)

	assert.Zero(t, len(mockapipkg.GetMockAPIs()))
}

func TestDeleteMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)
	defer func() { closeCh <- true }()

	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	uuid1, _, _ := writeDummyMockApiFile(t)
	uuid2, _, _ := writeDummyMockApiFile(t)
	uuid3, _, _ := writeDummyMockApiFile(t)

	defer func() {
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid1)
		removeMockApiFile(t, uuid2)
		removeMockApiFile(t, uuid3)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	// test delete api
	url := "/dynamocker/api/mock-api/" + fmt.Sprint(uuid1)
	r := httptest.NewRecorder()
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("DELETE", url, nil))
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 2, len(mockapipkg.GetMockApiList()))
}

func TestPostMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)
	defer func() { closeCh <- true }()

	// wait
	time.Sleep(50 * time.Millisecond)

	r := httptest.NewRecorder()

	// test post api without key
	mockApiPost := dummyMockApi(t)
	bytesPost, err := json.Marshal(mockApiPost)
	if err != nil {
		t.Fatalf("error while marshalign object : %s", err)
	}

	// create POST request
	postReqUrl := "/dynamocker/api/mock-api"
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("POST", postReqUrl, bytes.NewBuffer(bytesPost)))

	// check status code
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(100 * time.Millisecond)

	// check content of file
	var files []fs.DirEntry
	var jsonFilescounter = 0
	var uuidString string
	var found = false
	// retrieve uuid of the mockApi just written in the temp file after the POST request
	if files, err = os.ReadDir(os.TempDir() + "/"); err != nil {
		t.Fatalf("error while getting entries from the mock api folder: %s", err)
	}
	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		jsonFilescounter++
		uuidString, found = strings.CutSuffix(file.Name(), ".json")
	}
	assert.True(t, found)
	assert.Equal(t, 1, jsonFilescounter, "this means that some other json file is present in the test folder, jeopardizing the test result")

	mockApiUuid64, err := strconv.ParseUint(uuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid '%s' into uint16", uuidString)
		t.Fatal(err)
	}
	uuid := uint16(mockApiUuid64)

	defer func() {
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid)
	}()

	jsonFile, err := os.Open(os.TempDir() + "/" + uuidString + ".json")
	if err != nil {
		t.Fatalf("cannot open the file: %s", err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("cannot read from the file: %s", err)
	}
	var mockApi common.MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		t.Fatalf("cannot unmarshal bytes: %s", err)
	}

	// check that the mockApi passed to the POST is equal to the one just read from the file
	assert.Equal(t, mockApiPost.URL, mockApi.URL)
	assert.Equal(t, mockApiPost.Name, mockApi.Name)
	assert.Equal(t, mockApiPost.Responses.Get, mockApi.Responses.Get)
	assert.Equal(t, mockApiPost.Responses.Post, mockApi.Responses.Post)
	assert.Equal(t, mockApiPost.Responses.Delete, mockApi.Responses.Delete)
	assert.Equal(t, mockApiPost.Responses.Patch, mockApi.Responses.Patch)
}

func TestPutMockApi(t *testing.T) {
	// setup server and mockApi mgmt
	closeCh, webServerTest := setup(t)
	defer func() { closeCh <- true }()

	// wait
	time.Sleep(50 * time.Millisecond)

	// write mock api
	uuid, _, mockApi := writeDummyMockApiFile(t)

	defer func() {
		// wait
		time.Sleep(50 * time.Millisecond)
		removeMockApiFile(t, uuid)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, len(mockapipkg.GetMockApiList()))

	url := "/dynamocker/api/mock-api/" + fmt.Sprint(uuid)
	r := httptest.NewRecorder()

	// test patch api
	mockApi.Name = "new_name"
	mockApi.URL = "new-url.it"
	if json.Unmarshal([]byte(`{"new_get":true,"body":"new response"}`), &mockApi.Responses.Get) != nil {
		t.Fatalf("error while unmarshalling")
	}
	if json.Unmarshal([]byte(`{"new_delete":"deleted"}`), &mockApi.Responses.Delete) != nil {
		t.Fatalf("error while unmarshalling")
	}
	if json.Unmarshal([]byte(`{"new_post":1,"success":true}`), &mockApi.Responses.Post) != nil {
		t.Fatalf("error while unmarshalling")
	}
	mockApi.Responses.Patch = nil
	if json.Unmarshal([]byte(`{"patched":"yes","id":3}`), &mockApi.Responses.Patch) != nil {
		t.Fatalf("error while unmarshalling")
	}
	bytesPatch, err := json.Marshal(mockApi)
	if err != nil {
		t.Fatalf("error while marshalign object : %s", err)
	}
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("PUT", url, bytes.NewBuffer(bytesPatch)))
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(50 * time.Millisecond)

	// check that the mockApi has been modified
	currentMockApi, found := mockapipkg.GetApiByName(mockApi.Name)
	assert.True(t, found)
	assert.Equal(t, currentMockApi.URL, mockApi.URL)
	assert.Equal(t, currentMockApi.Name, mockApi.Name)
	assert.Equal(t, currentMockApi.Responses.Get, mockApi.Responses.Get)
	assert.Equal(t, currentMockApi.Responses.Post, mockApi.Responses.Post)
	assert.Equal(t, currentMockApi.Responses.Delete, mockApi.Responses.Delete)
	assert.Equal(t, currentMockApi.Responses.Patch, mockApi.Responses.Patch)
}

func TestServeMockApi(t *testing.T) {
	assert.True(t, true)
	// setup server and mockApi mgmt
	// closeCh, webServerTest := setup(t)
	// defer func() { closeCh <- true }()

	// // wait
	// time.Sleep(50 * time.Millisecond)

	// // write mock api
	// uuid, _, mockApi := writeDummyMockApiFile(t)

	// defer func() {
	// 	// wait
	// 	time.Sleep(50 * time.Millisecond)
	// 	removeMockApiFile(t, uuid)
	// }()

	// // wait
	// time.Sleep(50 * time.Millisecond)

	// assert.Equal(t, 1, len(mockapi.GetMockApiList()))

	// // generate request
	// url := mockApi.Name
	// r := httptest.NewRecorder()

	// // test get response of the MockApi
	// webServerTest.router.ServeHTTP(r, httptest.NewRequest("GET", url, nil))
	// assert.Equal(t, http.StatusNoContent, r.Code)

	// // wait
	// time.Sleep(50 * time.Millisecond)

	// TODO: complete the test

}

func removeMockApiFile(t *testing.T, uuid uint16) {

	filePath := os.TempDir() + "/" + fmt.Sprintf("%d", uuid) + ".json"
	_, err := os.Stat(filePath)
	if err == nil {
		err = os.Remove(filePath)
		if err != nil {
			t.Fatal("file not removed")
		}
	}
}
