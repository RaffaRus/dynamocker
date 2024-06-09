package webserver

import (
	"bytes"
	mockapi "dynamocker/internal/mock-api"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *WebServer {

	// set folderPath
	if err := os.Setenv("DYNA_MOCK_API_FOLDER", os.TempDir()); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

	// init the mocked api management
	closeCh := make(chan bool)
	var wg sync.WaitGroup
	if err := mockapi.Init(closeCh, &wg); err != nil {
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
	return webServerTest
}

func dummyMockApi() mockapi.MockApi {
	return mockapi.MockApi{
		Name:         fmt.Sprintf("dummy-mock-api-%d", rand.Intn(1000)),
		URL:          "url.com",
		FilePath:     os.TempDir(),
		Added:        time.Now(),
		LastModified: time.Now(),
		Responses: mockapi.Response{
			Get:    ptr(`{"valid_json":true,"body":"this is the response"}`),
			Patch:  ptr(`{"example_patch_body":"this is a string returned from patch operation"}`),
			Post:   ptr(`{"error":"posted an invalid element"}`),
			Delete: ptr(`{"response":"removed the item number 3"}`),
		},
	}
}

// write a dummy mock api file to the Temp folder. The temp folder
// comes from os package
func writeDummyMockApiFile(t *testing.T) (*os.File, mockapi.MockApi) {
	mockApi := dummyMockApi()
	filePath := mockApi.FilePath + "/" + mockApi.Name + ".json"

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		file.Close()
		t.Fatalf("cannot open file :%s", err)
	}

	defer file.Close()
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

func TestGetMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// write dummy mock Api
	_, mockApi := writeDummyMockApiFile(t)
	defer removeMockApiFile(t, mockApi)

	// wait
	time.Sleep(50 * time.Millisecond)

	// test get/{id} api
	r := httptest.NewRecorder()
	url := "/dynamocker/api/mock-api/" + mockApi.Name
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("GET", url, nil))
	assert.Equal(t, http.StatusOK, r.Code)
	bytesMockApi, err := json.Marshal(mockApi)
	if err != nil {
		t.Fatalf("error while marshaling mockApi: %s", err)
	}
	bytesResp, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error while reading from response body: %s", err)
	}
	assert.Equal(t, append(bytesMockApi, []byte("\n")...), bytesResp)
}

func TestGetMockApis(t *testing.T) {

	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	_, mockApi1 := writeDummyMockApiFile(t)
	time.Sleep(50 * time.Millisecond)
	_, mockApi2 := writeDummyMockApiFile(t)
	time.Sleep(50 * time.Millisecond)
	_, mockApi3 := writeDummyMockApiFile(t)

	defer func() {
		removeMockApiFile(t, mockApi1)
		removeMockApiFile(t, mockApi2)
		removeMockApiFile(t, mockApi3)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	// test get api
	url := "/dynamocker/api/mock-apis"
	r := httptest.NewRecorder()
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("GET", url, nil))
	assert.Equal(t, http.StatusOK, r.Code)
	bytesResp, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error while reading from response body: %s", err)
	}
	var mockApisResp []mockapi.MockApi
	err = json.Unmarshal(bytesResp, &mockApisResp)
	if err != nil {
		t.Fatalf("error while unmatshalling the response: %s", err)
	}
	assert.Equal(t, 3, len(mockApisResp))
}

func TestDeleteMockApis(t *testing.T) {

	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	_, mockApi1 := writeDummyMockApiFile(t)
	_, mockApi2 := writeDummyMockApiFile(t)
	_, mockApi3 := writeDummyMockApiFile(t)

	defer func() {
		removeMockApiFile(t, mockApi1)
		removeMockApiFile(t, mockApi2)
		removeMockApiFile(t, mockApi3)
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

	assert.Zero(t, len(mockapi.GetAPIs()))
}

func TestDeleteMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// wait
	time.Sleep(50 * time.Millisecond)

	// write three mock apis
	_, mockApi1 := writeDummyMockApiFile(t)
	_, mockApi2 := writeDummyMockApiFile(t)
	_, mockApi3 := writeDummyMockApiFile(t)

	defer func() {
		removeMockApiFile(t, mockApi1)
		removeMockApiFile(t, mockApi2)
		removeMockApiFile(t, mockApi3)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	// test delete api
	url := "/dynamocker/api/mock-api/" + mockApi2.Name
	r := httptest.NewRecorder()
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("DELETE", url, nil))
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 2, len(mockapi.GetAPIs()))
}

func TestPostMockApi(t *testing.T) {

	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// wait
	time.Sleep(50 * time.Millisecond)

	r := httptest.NewRecorder()

	// test post api without key
	mockApiPost := dummyMockApi()
	bytesPost, err := json.Marshal(mockApiPost)
	if err != nil {
		t.Fatalf("error while marshalign object : %s", err)
	}

	// POST request
	postReqUrl := "/dynamocker/api/mock-api/" + mockApiPost.Name
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("POST", postReqUrl, bytes.NewBuffer(bytesPost)))

	// check status code
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(100 * time.Millisecond)

	// check content of file
	file, err := os.Stat(os.TempDir() + "/" + mockApiPost.Name + ".json")
	assert.Nil(t, err)

	// retrieve mockApi just written in the temp file after the POST request
	jsonFile, err := os.Open(os.TempDir() + "/" + file.Name())
	if err != nil {
		t.Fatalf("cannot open the file: %s", err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("cannot read from the file: %s", err)
	}
	var mockApi mockapi.MockApi
	err = json.Unmarshal(byteValue, &mockApi)
	if err != nil {
		t.Fatalf("cannot unmarshal bytes: %s", err)
	}

	// check that the mockApi passed to the POST is equal to the one just read from the file
	assert.Equal(t, mockApiPost.URL, mockApi.URL)
	assert.Equal(t, mockApiPost.Name, mockApi.Name)
	assert.Equal(t, os.TempDir(), mockApi.FilePath)
	assert.Equal(t, mockApiPost.Responses.Get, mockApi.Responses.Get)
	assert.Equal(t, mockApiPost.Responses.Post, mockApi.Responses.Post)
	assert.Equal(t, mockApiPost.Responses.Delete, mockApi.Responses.Delete)
	assert.Equal(t, mockApiPost.Responses.Patch, mockApi.Responses.Patch)
	// time cannot be compared using the "==" operator
	assert.True(t, mockApi.Added.Equal(mockApiPost.Added))
	assert.True(t, mockApi.LastModified.Equal(mockApiPost.LastModified))
}

func TestPatch(t *testing.T) {
	// setup server and mockApi mgmt
	webServerTest := setup(t)

	// wait
	time.Sleep(50 * time.Millisecond)

	// write mock api
	_, mockApi := writeDummyMockApiFile(t)

	defer func() {
		removeMockApiFile(t, mockApi)
	}()

	// wait
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, len(mockapi.GetAPIs()))

	// test patch api
	mockApi.URL = "new-url.it"
	newTime := time.Now()
	mockApi.LastModified = newTime
	mockApi.Responses.Get = ptr(`{"new_get":true,"body":"new response"}`)
	mockApi.Responses.Delete = ptr(`{"new_delete":"deleted"}`)
	mockApi.Responses.Post = ptr(`{"new_post":1,"success":true}`)
	mockApi.Responses.Patch = ptr(`{"patched":"yes","id":3}`)
	url := "/dynamocker/api/mock-api/" + mockApi.Name
	r := httptest.NewRecorder()
	bytesPatch, err := json.Marshal(mockApi)
	if err != nil {
		t.Fatalf("error while marshalign object : %s", err)
	}
	webServerTest.router.ServeHTTP(r, httptest.NewRequest("PATCH", url, bytes.NewBuffer(bytesPatch)))
	assert.Equal(t, http.StatusNoContent, r.Code)

	// wait
	time.Sleep(50 * time.Millisecond)

	// check that the mockApi has been modified
	currentMockApi, err := mockapi.GetAPI(mockApi.Name)
	assert.Nil(t, err)
	assert.Equal(t, currentMockApi.URL, mockApi.URL)
	assert.Equal(t, currentMockApi.Name, mockApi.Name)
	assert.Equal(t, os.TempDir(), mockApi.FilePath)
	assert.Equal(t, currentMockApi.Responses.Get, mockApi.Responses.Get)
	assert.Equal(t, currentMockApi.Responses.Post, mockApi.Responses.Post)
	assert.Equal(t, currentMockApi.Responses.Delete, mockApi.Responses.Delete)
	assert.Equal(t, currentMockApi.Responses.Patch, mockApi.Responses.Patch)
	// time cannot be compared using the "==" operator
	assert.True(t, mockApi.Added.Equal(currentMockApi.Added))
	assert.True(t, mockApi.LastModified.Equal(currentMockApi.LastModified))
}

func ptr[A any](a A) *A {
	return &a
}

func removeMockApiFile(t *testing.T, mockApi mockapi.MockApi) {

	filename := mockApi.FilePath + "/" + mockApi.Name + ".json"
	_, err := os.Stat(filename)
	if err == nil {
		err = os.Remove(filename)
		if err != nil {
			t.Fatal("file not removed")
		}
	}
}
