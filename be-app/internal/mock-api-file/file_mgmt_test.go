package mockapifilepkg

import (
	"dynamocker/internal/common"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// reset package map and folderPath variable
func reset() {
	folderPath = ""
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
		URL:       "url.com",
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

func TestAddNewMockApiFile(t *testing.T) {
	reset()

	// add file while folderpath == ""
	assert.EqualError(t, AddNewMockApiFile([]byte{}), "the mock API folder has not been set-up")

	// set mock api folder as a temp folder
	folderPath = os.TempDir() + "/"

	// add valid mock api
	api := dummyMockApi(t)
	bytes, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, AddNewMockApiFile(bytes))

	// check the mockApi has been added
	var files []fs.DirEntry
	var jsonFilescounter = 0
	var uuidString string
	var found = false
	if files, err = os.ReadDir(folderPath); err != nil {
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

	defer func() {
		filename := folderPath + uuidString + ".json"
		_, err := os.Stat(filename)
		if err == nil {
			err = os.Remove(filename)
			if err != nil {
				t.Fatal("file not removed")
			}
		}
	}()

	// add invalid json
	assert.EqualError(t, AddNewMockApiFile([]byte("invalid json")), "error while unmarshaling body: invalid character 'i' looking for beginning of value")

	// add struct with no Name
	invalidStruct := dummyMockApi(t)
	invalidStruct.Name = string("")
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile(bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.Name' Error:Field validation for 'Name' failed on the 'required' tag")

	// add struct with no URL
	invalidStruct = dummyMockApi(t)
	invalidStruct.URL = string("")
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile(bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.URL' Error:Field validation for 'URL' failed on the 'required' tag")

	// add struct with no Response
	invalidStruct = dummyMockApi(t)
	invalidStruct.Responses = common.Response{}
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile(bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.Responses' Error:Field validation for 'Responses' failed on the 'required' tag")

}

func TestRemoveMockApiFile(t *testing.T) {
	reset()

	// remove file while folderpath == ""
	assert.EqualError(t, RemoveMockApiFile(0), "the mock API folder has not been set-up")

	folderPath = os.TempDir() + "/"

	// add mock api
	uuid, dummyMockApiFile, _ := writeDummyMockApiFile(t)
	defer func() {
		dummyMockApiFile.Close()
		os.Remove(os.TempDir() + "/" + fmt.Sprint(uuid) + ".json")
	}()

	// check that the api has been loaded
	mockApis, err := LoadAPIsFromFolder()
	assert.Nil(t, err, "the function should return no error")
	assert.Equal(t, 1, len(mockApis))
	_, found := mockApis[uuid]
	assert.True(t, found)

	err = RemoveMockApiFile(uuid)
	assert.Nil(t, err)

	// check that the api has been removed
	mockApis, err = LoadAPIsFromFolder()
	assert.Nil(t, err, "the function should return no error")
	assert.Equal(t, 0, len(mockApis))

}

func TestRemoveAllMockApisFiles(t *testing.T) {
	reset()

	// remove all files while folderpath == ""
	assert.EqualError(t, RemoveAllMockApisFiles(), "the mock API folder has not been set-up")

	// unexisting fodler path
	folderPath = "not_existing_path"
	assert.EqualError(t, RemoveAllMockApisFiles(), "error while getting entries from the mock api folder: open not_existing_path: no such file or directory")

	// it does not remove *json file
	folderPath = os.TempDir() + "/"
	file, err := os.CreateTemp(folderPath, "random_file.log")
	defer func() {
		err = os.Remove(file.Name())
		if err != nil {
			t.Fatal("file not removed")
		}
	}()
	assert.Nil(t, RemoveAllMockApisFiles())
	_, err = os.Stat(file.Name())
	assert.Nil(t, err)

}

func TestModifyMockApiFile(t *testing.T) {
	reset()

	// add mock api
	folderPath = os.TempDir() + "/"
	api := dummyMockApi(t)
	bytes, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, AddNewMockApiFile(bytes))

	// check the mockApi has been added
	var files []fs.DirEntry
	var jsonFilescounter = 0
	var uuidString string
	var found = false
	if files, err = os.ReadDir(folderPath); err != nil {
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

	defer func() {
		filename := folderPath + uuidString + ".json"
		_, err := os.Stat(filename)
		if err == nil {
			err = os.Remove(filename)
			if err != nil {
				t.Fatal("file not removed")
			}
		}
	}()

	mockApiUuid64, err := strconv.ParseUint(uuidString, 10, 16)
	if err != nil {
		err := fmt.Errorf("error while parsing uuid '%s' into uint16", uuidString)
		t.Fatal(err)
	}
	uuid := uint16(mockApiUuid64)

	// modify the mockApi file
	newApi := api
	newApi.Responses.Get = &map[string]interface{}{}
	newApi.Responses.Post = &map[string]interface{}{}
	newApi.Responses.Delete = &map[string]interface{}{}
	newApi.Responses.Patch = &map[string]interface{}{}
	if err = json.Unmarshal([]byte(`{"new_json":true,"new_body":"a new response"}`), newApi.Responses.Get); err != nil {
		t.Fatal("error while unmarshalling")
	}
	if json.Unmarshal([]byte(`{"this_is":4}`), newApi.Responses.Patch) != nil {
		t.Fatal("error while unmarshalling")
	}
	if json.Unmarshal([]byte(`{"there_you_go":"maybe","nope":false}`), newApi.Responses.Post) != nil {
		t.Fatal("error while unmarshalling")
	}
	if json.Unmarshal([]byte(`{"still":true,"later":4,"tomorrow":"not sure"}`), newApi.Responses.Delete) != nil {
		t.Fatal("error while unmarshalling")
	}
	newBytes, err := json.Marshal(newApi)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, ModifyMockApiFile(uuid, newBytes))

	// check it was modified
	filename := folderPath + uuidString
	filebytes, err := os.ReadFile(filename + ".json")
	if err != nil {
		t.Fatalf("error file not read :%s", err)
	}
	assert.Equal(t, newBytes, filebytes)

}

func TestLoadStoredAPIs(t *testing.T) {
	reset()

	// call function before defining any folder. This should log only
	_, err := LoadAPIsFromFolder()
	assert.Equal(t, fmt.Errorf("the mock API folder has not been set-up"), err)

	// set temp folder as the one contining the mock api files
	folderPath = os.TempDir() + "/"
	_, err = LoadAPIsFromFolder()
	assert.Nil(t, err)

	// add mock api
	uuid, dummyMockApiFile, _ := writeDummyMockApiFile(t)
	defer func() {
		dummyMockApiFile.Close()
		os.Remove(os.TempDir() + "/" + fmt.Sprint(uuid) + ".json")
	}()

	// check that the apis have been loaded
	mockApis, err := LoadAPIsFromFolder()
	assert.Nil(t, err, "the function should return no error")
	assert.Equal(t, 1, len(mockApis))
	_, found := mockApis[uuid]
	assert.True(t, found)
}
