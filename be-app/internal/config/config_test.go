package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getterConfTest struct {
	functionRes       func() (string, error)
	defaultResult     string
	overwrittenResult string
	relevantKey       string
}

func TestConfig(t *testing.T) {
	getterTests := []getterConfTest{
		{
			GetLogLevel,
			"INFO",
			"WARN",
			"DYNA_LOG_LEVEL",
		},
		{
			GetServerPort,
			"8150",
			"9999",
			"DYNA_SERVER_PORT",
		},
		{
			GetMockApiFolder,
			"/var/dynamocker/mocks/",
			"test_folder",
			"DYNA_MOCK_API_FOLDER",
		},
	}
	for _, test := range getterTests {
		test.Tester(t)
	}

}

func (c getterConfTest) Tester(t *testing.T) {

	// test that Getter function returns default value
	result, err := c.functionRes()
	assert.Nil(t, err, "there should be no error")
	assert.Equal(t, c.defaultResult, result, "getter returned an unexpected result")

	// test that Getter function returns custom value
	envVarList[c.relevantKey] = c.overwrittenResult
	result, err = c.functionRes()
	assert.Nil(t, err, "there should be no error")
	assert.Equal(t, c.overwrittenResult, result, "getter does not return not-default folder", c.functionRes)

	// test that Getter function returns the correct error when no key is found in the map
	delete(envVarList, c.relevantKey)
	_, err = c.functionRes()
	assert.Equal(t, fmt.Errorf(fmt.Sprintf("element %s not found in the map", c.relevantKey)), err, "incorrect error returned from the getter")

}

// it should read all the env variables set in the system and update the env variables saved in the package
func TestReadVars(t *testing.T) {

	// unset env variables if they exist already
	for _, key := range []string{"DYNA_LOG_LEVEL", "DYNA_SERVER_PORT", "DYNA_MOCK_API_FOLDER"} {
		if os.Getenv(key) != "" {
			err := os.Unsetenv(key)
			if err != nil {
				t.Fatalf("impossible to unset env variable %s", key)
			}
		}
	}

	// if no matching env var has been set, let the default ones
	ReadVars()
	assert.Equal(t, "INFO", envVarList["DYNA_LOG_LEVEL"], "incorrect default env variable")
	assert.Equal(t, "8150", envVarList["DYNA_SERVER_PORT"], "incorrect default env variable")
	assert.Equal(t, "/var/dynamocker/mocks/", envVarList["DYNA_MOCK_API_FOLDER"], "incorrect default env variable")

	// check the custom ones
	if err := os.Setenv("DYNA_LOG_LEVEL", "DEBUG"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}
	if err := os.Setenv("DYNA_SERVER_PORT", "5555"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}
	if err := os.Setenv("DYNA_MOCK_API_FOLDER", "mock_folder"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

}
