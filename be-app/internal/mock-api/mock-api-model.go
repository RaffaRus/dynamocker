package mockapi

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strings"
	"time"
)

type MockApi struct {
	name         string
	URL          url.URL   `json:"url" ,validate:"base64url"`
	FilePath     string    `json:"filePath" validate:"dirpath"`
	Added        time.Time `json:"added" validate:"ltecsfield=InnerStructField.StartDate"`
	LastModified time.Time `json:"lastModified" validate:"ltecsfield=InnerStructField.StartDate"`
}

// it must act on the file. observer will do its job
func AddNewMockApiFile(fileName string, fileContent []byte) error {

	filePath := folderPath + "/" + fileName + ".json"

	if err := os.WriteFile(filePath, fileContent, fs.ModePerm); err != nil {
		return fmt.Errorf("file %s not created: %s", filePath, err)
	}

	return nil
}

// it must act on the file. observer will do its job
func RemoveMockApiFile(fileName string) error {

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
