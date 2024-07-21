package mockapi

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewMockApiFile(t *testing.T) {
	reset(t)

	// add file while folderpath == ""
	assert.EqualError(t, AddNewMockApiFile("", []byte{}), "the mock API folder has not been set-up")

	// set mock api folder as a temp folder
	folderPath = os.TempDir()

	// add valid mock api
	api := dummyMockApi(t)
	defer func() {
		filename := api.FilePath + "/" + api.Name
		_, err := os.Stat(filename)
		if err == nil {
			err = os.Remove(filename)
			if err != nil {
				t.Fatal("file not removed")
			}
		}
	}()
	bytes, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, AddNewMockApiFile(api.Name, bytes))

	// add invalid json
	assert.EqualError(t, AddNewMockApiFile("", []byte("invalid json")), "error while unmarshaling body: invalid character 'i' looking for beginning of value")

	// add struct with no Name
	invalidStruct := dummyMockApi(t)
	invalidStruct.Name = string("")
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile("", bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.Name' Error:Field validation for 'Name' failed on the 'required' tag")

	// add struct with no URL
	invalidStruct = dummyMockApi(t)
	invalidStruct.URL = string("")
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile("", bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.URL' Error:Field validation for 'URL' failed on the 'required' tag")

	// add struct with no Response
	invalidStruct = dummyMockApi(t)
	invalidStruct.Responses = Response{}
	bytes, err = json.Marshal(invalidStruct)
	if err != nil {
		t.Fatal("error while marshaling struct")
	}
	assert.EqualError(t, AddNewMockApiFile("", bytes), "invalid mock api passed from post request: %!s(<nil>)\nKey: 'MockApi.Responses' Error:Field validation for 'Responses' failed on the 'required' tag")

}

func TestRemoveMockApiFile(t *testing.T) {
	reset(t)

	// remove file while folderpath == ""
	assert.EqualError(t, RemoveMockApiFile(""), "the mock API folder has not been set-up")

}

func TestRemoveAllMockApisFiles(t *testing.T) {
	reset(t)

	// remove all files while folderpath == ""
	assert.EqualError(t, RemoveAllMockApisFiles(), "the mock API folder has not been set-up")

	// unexisting fodler path
	folderPath = "not_existing_path"
	assert.EqualError(t, RemoveAllMockApisFiles(), "error while getting entries from the mock api folder: open not_existing_path: no such file or directory")

	// it does not remove *json file
	folderPath = os.TempDir()
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
	reset(t)

	// add mock api
	folderPath = os.TempDir()
	api := dummyMockApi(t)
	defer func() {
		filename := api.FilePath + "/" + api.Name + ".json"
		_, err := os.Stat(filename)
		if err == nil {
			err = os.Remove(filename)
			if err != nil {
				t.Fatal("file not removed")
			}
		}
	}()
	bytes, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, AddNewMockApiFile(api.Name, bytes))

	// modify it
	var newApi MockApi
	err = json.Unmarshal(bytes, &newApi)
	if err != nil {
		t.Fatalf("errror while unmarshaling : %s", err)
	}
	if json.Unmarshal([]byte(`{"new_json":true,"new_body":"a new response"}`), newApi.Responses.Get) != nil {
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
	assert.Nil(t, ModifyMockApiFile(newApi.Name, newBytes))

	// check it was modified
	filename := newApi.FilePath + "/" + newApi.Name
	filebytes, err := os.ReadFile(filename + ".json")
	if err != nil {
		t.Fatalf("error file not read :%s", err)
	}
	assert.Equal(t, newBytes, filebytes)

}
