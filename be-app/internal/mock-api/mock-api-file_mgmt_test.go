package mockapi

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewMockApiFile(t *testing.T) {
	reset(t)

	// set mock api folder as a temp folder
	folderPath = os.TempDir()

	dummyMockApi := dummyMockApi()
	bytes, err := json.Marshal(dummyMockApi)
	if err != nil {
		t.Fatalf("error while marshaling dummy mock api :%s", err)
	}
	assert.Nil(t, AddNewMockApiFile(dummyMockApi.Name, bytes))
}

func TestRemoveMockApiFile(t *testing.T) {

}

func TestRemoveAllMockApisFiles(t *testing.T) {

}

func TestModifyMockApiFile(t *testing.T) {

}
