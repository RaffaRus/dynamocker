package mockapipkg

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var mu sync.Mutex
var folderPath = ""

// it must act on the file. observer will do its job
func AddNewMockApiFile(body []byte) error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	// unmashal body
	var mockApi MockApi
	err := json.Unmarshal(body, &mockApi)
	if err != nil {
		err := fmt.Errorf("error while unmarshaling body: %s", err)
		return err
	}

	// validate body
	vtor := validator.New(validator.WithRequiredStructEnabled())
	vtorErr := vtor.Struct(mockApi)
	if vtorErr != nil {
		valErrs := vtorErr.(validator.ValidationErrors)
		var valErrsCumulative error
		for _, valErr := range valErrs {
			valErrsCumulative = fmt.Errorf("%s\n%s", valErrsCumulative, valErr.Error())
		}
		if valErrsCumulative != nil {
			err := fmt.Errorf("invalid mock api passed from post request: %s", valErrsCumulative)
			return err
		}
	}

	// check if a mockApi with the same name or URL already exists
	_, found := GetApiByName(mockApi.Name)
	if found {
		return fmt.Errorf("found another file using same name of the one to be added ('%s'). File %s not created", mockApi.Name, err)
	}
	_, found = GetApiByUrl(mockApi.Name)
	if found {
		return fmt.Errorf("found another file using same URL of the one to be added ('%s'). File %s not created", mockApi.URL, err)
	}

	// retrieve file path
	filePath := folderPath + "/" + mockApi.Name + ".json"

	// transform mockApi into []byte
	bytes, err := json.Marshal(mockApi)
	if err != nil {
		return fmt.Errorf("file %s not created. error while marshalling modified mockapi: %s", filePath, err)
	}

	// write mockapi
	if err := os.WriteFile(filePath, bytes, fs.ModePerm); err != nil {
		return fmt.Errorf("file %s not created: %s", filePath, err)
	}

	return nil
}

// it must act on the file. observer will do its job
func RemoveMockApiFile(uuid uint16) error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	mockApi, found := mockApiList[uuid]
	if !found {
		return fmt.Errorf("no mockApi found with the given uuid (%d)", uuid)
	}
	filePath := folderPath + "/" + mockApi.Name + ".json"

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("file %s not removed: %s", filePath, err)
	}

	return nil
}

// it must act on the file. observer will do its job
func RemoveAllMockApisFiles() error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	var files []fs.DirEntry
	var err error

	if files, err = os.ReadDir(folderPath); err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := folderPath + "/" + file.Name()

		if err = os.Remove(filePath); err != nil {
			return fmt.Errorf("file %s not removed: %s", file.Name(), err)
		}

	}

	return nil

}

// it must act on the file. observer will do its job
func ModifyMockApiFile(mockApiUuid uint16, newFile []byte) error {

	if err := RemoveMockApiFile(mockApiUuid); err != nil {
		return err
	}

	if err := AddNewMockApiFile(newFile); err != nil {
		return err
	}

	return nil
}
