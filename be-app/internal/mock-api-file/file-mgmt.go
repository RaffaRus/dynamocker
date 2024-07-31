package mockapifilepkg

import (
	"dynamocker/internal/common"
	"dynamocker/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"sync"

	"math/rand"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var mu sync.Mutex
var folderPath = ""

func Init() error {
	folderPath = config.GetMockApiFolder()
	return nil
}

// it must act on the file. observer will do its job
func AddNewMockApiFile(body []byte) error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	// unmashal body
	var mockApi common.MockApi
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
	// _, found := GetApiByName(mockApi.Name)
	// if found {
	// 	return fmt.Errorf("found another file using same name of the one to be added ('%s'). File %s not created", mockApi.Name, err)
	// }
	// _, found = GetApiByUrl(mockApi.Name)
	// if found {
	// 	return fmt.Errorf("found another file using same URL of the one to be added ('%s'). File %s not created", mockApi.URL, err)
	// }

	// retrieve file path
	filePath := folderPath + fmt.Sprint(generateUuid()) + ".json"

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

	file, err := os.Stat(folderPath + fmt.Sprint(uuid) + ".json")

	if err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	if err = os.Remove(folderPath + file.Name()); err != nil {
		return fmt.Errorf("file %s not removed: %s", file.Name(), err)
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

		filePath := folderPath + file.Name()

		if err = os.Remove(filePath); err != nil {
			return fmt.Errorf("file %s not removed: %s", file.Name(), err)
		}

	}

	return nil

}

// it must act on the file. observer will do its job
func ModifyMockApiFile(mockApiUuid uint16, newFile []byte) error {

	mu.Lock()
	defer mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("the mock API folder has not been set-up")
	}

	file, err := os.Stat(folderPath + fmt.Sprint(mockApiUuid) + ".json")

	if err != nil {
		return fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	if err = os.Remove(folderPath + file.Name()); err != nil {
		return fmt.Errorf("file %s not removed: %s", file.Name(), err)
	}

	// unmashal body
	var mockApi common.MockApi
	err = json.Unmarshal(newFile, &mockApi)
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

	// retrieve file path
	filePath := folderPath + fmt.Sprint(mockApiUuid) + ".json"

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

// loading the APIs from the mock api folder at startup
// this function updates the list based on the entried loaded from the folder
func LoadAPIsFromFolder() (map[uint16]*common.MockApi, error) {

	mu.Lock()
	defer mu.Unlock()

	mockApiList := make(map[uint16]*common.MockApi)

	// get path from config package
	if folderPath == "" {
		return nil, fmt.Errorf("the mock API folder has not been set-up")
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error while getting entries from the mock api folder: %s", err)
	}

	for _, file := range files {

		// select only *.json files
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		pathToFile := folderPath + "/" + file.Name()

		uuidString, found := strings.CutSuffix(file.Name(), ".json")
		if !found {
			err = fmt.Errorf("suffix '.json' not found")
			return nil, err
		}
		mockApiUuid64, err := strconv.ParseUint(uuidString, 10, 16)
		if err != nil {
			err := fmt.Errorf("error while parsing uuid of the mockApi file '%s' into uint16", file.Name())
			log.Error(err)
		}
		uuid := uint16(mockApiUuid64)

		// open file
		jsonFile, err := os.Open(pathToFile)
		if err != nil {
			log.Errorf("error while opening the file %s: %s", pathToFile, err)
			continue
		}

		// read content
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Errorf("error while reading the file %s: %s", pathToFile, err)
			continue
		}

		// unmarshal content
		var mockApi common.MockApi
		err = json.Unmarshal(byteValue, &mockApi)
		if err != nil {
			log.Errorf("error while unmarshaling the json file %s into the struct: %s", pathToFile, err)
			continue
		}

		// validate content
		vtor := validator.New(validator.WithRequiredStructEnabled())
		err = vtor.Struct(mockApi)
		if err != nil {
			log.Errorf("invalid mock api saved in the json file %s into the struct: %s", pathToFile, err)
			continue
		}

		// add to the map
		mockApiList[uuid] = &mockApi

	}

	return mockApiList, nil
}

// generate a random uuid
func generateUuid() uint16 {
	var tmp uint16
	var counter uint16 = 0
	for {
		if counter > 100 {
			log.Fatal("Reached maximum numbers of attempts in generating a random uuid")
		}
		tmp = uint16(rand.Intn(common.MAX_SIZE_MOCKAPI_LIST))

		_, err := os.Stat(folderPath + fmt.Sprint(tmp) + "json")

		if errors.Is(err, fs.ErrNotExist) {
			break
		}
		counter++
	}
	return tmp
}
