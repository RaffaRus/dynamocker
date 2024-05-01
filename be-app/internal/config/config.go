package config

import (
	"errors"
	"os"
)

// List of the environment variables to be checked. The corresponding
// value is the default value
var envVarList map[string]string = map[string]string{
	"DYNA_LOG_LEVEL":       "INFO",
	"DYNA_SERVER_PORT":     "8150",
	"DYNA_MOCK_API_FOLDER": "/var/dynamocker/mocks/",
}

// read all the env variables
func ReadVars() {

	for env := range envVarList {
		if val := os.Getenv(env); val != "" {
			envVarList[env] = val
		}
	}
}

func GetLogLevel() (string, error) {

	if val, ok := envVarList["DYNA_LOG_LEVEL"]; ok {
		return val, nil
	} else {
		return "", errors.New("element DYNA_LOG_LEVEL not found in the map")
	}
}

func GetServerPort() (string, error) {

	if val, ok := envVarList["DYNA_SERVER_PORT"]; ok {
		return val, nil
	} else {
		return "", errors.New("element DYNA_SERVER_PORT not found in the map")
	}
}

func GetMockApiFolder() (string, error) {

	if val, ok := envVarList["DYNA_MOCK_API_FOLDER"]; ok {
		return val, nil
	} else {
		return "", errors.New("element DYNA_MOCK_API_FOLDER not found in the map")
	}
}
