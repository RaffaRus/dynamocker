package mockapi

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
var mockApiList = make(map[string]*MockApi)
var folderPath = ""

// it must act on the file. observer will do its job
func AddNewMockApiFile(fileName string, body []byte) error {

	mu.Lock()
	defer mu.Unlock()

	// unmashal body
	var mockApi MockApi
	err := json.Unmarshal(body, &mockApi)
	if err != nil {
		err := fmt.Errorf("error while unmarshaling body from patch request: %s", err)
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

	// retrieve file path
	filePath := folderPath + "/" + fileName + ".json"

	// write mockapi
	if err := os.WriteFile(filePath, body, fs.ModePerm); err != nil {
		return fmt.Errorf("file %s not created: %s", filePath, err)
	}

	return nil
}

// it must act on the file. observer will do its job
func RemoveMockApiFile(fileName string) error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	filePath := folderPath + "/" + fileName + ".json"

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
func ModifyMockApiFile(fileName string, newFile []byte) error {

	if err := RemoveMockApiFile(fileName); err != nil {
		return err
	}

	if err := AddNewMockApiFile(fileName, newFile); err != nil {
		return err
	}

	return nil
}
