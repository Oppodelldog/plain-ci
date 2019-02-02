package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// BuildDir defines a folder in which simple-ci will perform its build actions
var BuildDir string

// ServerPort for API like hooks or REST access
var ServerPort string

// GithubToken is a personal access token that gives the following permissions: repo:status
var GithubToken string

func init() {

	BuildDir = getEnv("SIMPLE_CI_BUILD_DIR", "/tmp/simple-ci")
	ServerPort = getEnv("SIMPLE_CI_SERVER_PORT", "12345")
	GithubToken = getEnv("SIMPLE_CI_GITHUB_TOKEN", "")

	logLevel := getEnv("SIMPLE_CI_LOG_LEVEL", "debug")
	loglevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(loglevel)
	}

	err = os.MkdirAll(BuildDir, 0766)
	if err != nil {
		panic(fmt.Sprintf("could not create buildDir '%s': %v", BuildDir, err))
	}
}

func getEnv(key, def string) string {
	if s, ok := os.LookupEnv(key); ok {
		return s
	} else {
		return def
	}
}
