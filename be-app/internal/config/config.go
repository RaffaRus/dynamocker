package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// List of the environment variables to be checked. The corresponding
// value is the default value
var envVarList map[string]string = map[string]string{
	"DYNA_LOG_LEVEL":          "INFO",
	"DYNA_SERVER_PORT":        "8150",
	"DYNA_MOCK_API_FOLDER":    "/var/dynamocker/mocks/",
	"DYNA_WEB_SERVER_VERSION": "1",
}

// read all the env variables. If found, replaces the default value with the one provided
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

func GetApiVersion() (uint16, error) {

	if val, ok := envVarList["DYNA_WEB_SERVER_VERSION"]; ok {
		if ver, err := strconv.ParseInt(val, 10, 16); err != nil {
			return 0, fmt.Errorf("error while parsing the API version to uint16: %s", err)
		} else {
			return uint16(ver), nil
		}
	} else {
		return 0, errors.New("element DYNA_WEB_SERVER_VERSION not found in the map")
	}
}
