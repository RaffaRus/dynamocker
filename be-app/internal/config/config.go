package config

import (
	"fmt"
	"os"
)

// List of the environment variables to be checked. The corresponding
// value is the default value
const (
	logEnv           = "DYNA_LOG_LEVEL"
	portEnv          = "DYNA_SERVER_PORT"
	folderEnv        = "DYNA_MOCK_API_FOLDER"
	logEnvDefault    = "INFO"
	portEnvDefault   = "8150"
	folderEnvDefault = "/var/dynamocker/mocks/"
)

var envVarList map[string]string = map[string]string{
	"DYNA_LOG_LEVEL":       "INFO",
	"DYNA_SERVER_PORT":     "8150",
	"DYNA_MOCK_API_FOLDER": "/var/dynamocker/mocks/",
}

// read all the env variables
func ReadVars() {

	fmt.Println("from function", envVarList)
	for env := range envVarList {
		fmt.Printf("looking for the key %s\n", env)
		if val := os.Getenv(env); val != "" {
			fmt.Printf("found the %s env variable\n", fmt.Sprintf("[env]: %s", val))
			envVarList[env] = val
		}
	}
}

func GetLogLevel() string {

	if val := os.Getenv(logEnv); val != "" {
		return val
	} else {
		return logEnvDefault
	}
}

func GetServerPort() string {
	if val := os.Getenv(portEnv); val != "" {
		return val
	} else {
		return portEnvDefault
	}
}

func GetMockApiFolder() string {
	if val := os.Getenv(folderEnv); val != "" {
		return val
	} else {
		return folderEnvDefault
	}
}
