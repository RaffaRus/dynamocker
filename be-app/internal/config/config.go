package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// List of the environment variables to be checked. The corresponding
// value is the default value
const (
	logEnv                = "DYNA_LOG_LEVEL"
	logEnvDefault         = "INFO"
	portEnv               = "DYNA_SERVER_PORT"
	portEnvDefault        = "8150"
	folderEnv             = "DYNA_MOCK_API_FOLDER"
	folderEnvDefault      = "/var/dynamocker/mocks/"
	pollerIntervalEnv     = "POLLER_INTERVAL"
	pollerIntervalDefault = "60" // seconds
)

var envVarList map[string]string = map[string]string{
	logEnv:            logEnvDefault,
	portEnv:           portEnvDefault,
	folderEnv:         folderEnvDefault,
	pollerIntervalEnv: pollerIntervalDefault,
}

// read all the env variables
func ReadVars() {

	for env := range envVarList {
		log.Debugf("looking for the key %s\n", env)
		if val := os.Getenv(env); val != "" {
			log.Infof("found the %s env variable with value = %s\n", env, val)
			envVarList[env] = val
		} else {
			log.Infof("env variable %s not found, using default value = %s\n", env, envVarList[env])
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

func GetPollingInterval() string {
	if val := os.Getenv(pollerIntervalEnv); val != "" {
		return val
	} else {
		return pollerIntervalDefault
	}
}
