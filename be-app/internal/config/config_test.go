package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getterConfTest struct {
	functionRes  func() string
	defaultValue string
	customValue  string
	relevantKey  string
}

func TestConfig(t *testing.T) {
	getterTests := []getterConfTest{
		{
			GetLogLevel,
			logEnvDefault,
			"WARN",
			logEnv,
		},
		{
			GetServerPort,
			portEnvDefault,
			"9999",
			portEnv,
		},
		{
			GetMockApiFolder,
			folderEnvDefault,
			"test_folder",
			folderEnv,
		},
	}
	for _, test := range getterTests {
		test.Tester(t)
	}

}

func (c getterConfTest) Tester(t *testing.T) {

	// test that Getter function returns default value
	result := c.functionRes()
	assert.NotEqual(t, result, "", "there should be no error")
	assert.Equal(t, c.defaultValue, result, "getter returned an unexpected result")

	// set custom env var
	os.Setenv(c.relevantKey, c.customValue)

	// test that Getter function returns custom value
	result = c.functionRes()
	assert.Equal(t, c.customValue, result, "getter does not return not-default folder", c.functionRes)

}

// it should read all the env variables set in the system and update the env variables saved in the package
func TestReadVars(t *testing.T) {

	fmt.Println("from test", envVarList)

	// unset env variables if they exist already
	for _, key := range []string{logEnv, portEnv, folderEnv, pollerIntervalEnv} {
		if os.Getenv(key) != "" {
			err := os.Unsetenv(key)
			if err != nil {
				t.Fatalf("impossible to unset env variable %s", key)
			}
		}
	}

	// if no matching env var has been set, let the default ones
	ReadVars()
	assert.Equal(t, logEnvDefault, envVarList[logEnv], "incorrect default env variable")
	assert.Equal(t, portEnvDefault, envVarList[portEnv], "incorrect default env variable")
	assert.Equal(t, folderEnvDefault, envVarList[folderEnv], "incorrect default env variable")

	// check the custom ones
	if err := os.Setenv(logEnv, "DEBUG"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}
	if err := os.Setenv(portEnv, "5555"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}
	if err := os.Setenv(folderEnv, "mock_folder"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}
	if err := os.Setenv(pollerIntervalEnv, "50"); err != nil {
		t.Fatalf("cannot set env variable: %s", err)
	}

	ReadVars()
	assert.Equal(t, "DEBUG", envVarList[logEnv], "incorrect default env variable")
	assert.Equal(t, "5555", envVarList[portEnv], "incorrect default env variable")
	assert.Equal(t, "mock_folder", envVarList[folderEnv], "incorrect default env variable")
	assert.Equal(t, "50", envVarList[pollerIntervalEnv], "incorrect default env variable")
}
