package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

// BuildDir defines a folder in which plain-ci will perform its build actions
var BuildDir string

// ServerPort for API like hooks or REST access
var ServerPort string

// ServerPortTLS for secure API like hooks or REST access
var ServerPortTLS string

// GithubToken is a personal access token that gives the following permissions: repo:status
var GithubToken string

//Domain the domain the service will be hosted at
var Domain string

//CertDir defines where to put certificate files
var CertDir string

//CertFile filepath to certificate
var CertFile string

//KeyFile filepath to certificate key
var KeyFile string

//ReadTimeout used as http server config
var ReadTimeout time.Duration

//WriteTimeout used as http server config
var WriteTimeout time.Duration

//IdleTimeout used as http server config
var IdleTimeout time.Duration

//GracefulShutdownPeriod seconds to wait for api server to shutdown while shutting down the service
var GracefulShutdownPeriod time.Duration

func init() {

	BuildDir = getEnv("PLAIN_CI_BUILD_DIR", "/tmp/plain-ci")
	ServerPort = getEnv("PLAIN_CI_SERVER_PORT", "10080")
	ServerPortTLS = getEnv("PLAIN_CI_SERVER_PORT_TLS", "10443")
	GithubToken = getEnv("PLAIN_CI_GITHUB_TOKEN", "")

	Domain = getEnv("PLAIN_CI_DOMAIN", "www.plain-ci.test")
	CertDir = getEnv("PLAIN_CI_CERT_DIR", ".cert")
	CertFile = getEnv("PLAIN_CI_CERT_FILE", "sample.crt")
	KeyFile = getEnv("PLAIN_CI_KEY_FILE", "sample.key")

	ReadTimeout = getDurationFromEnv("PLAIN_CI_READ_TIMEOUT", time.Second*5)
	WriteTimeout = getDurationFromEnv("PLAIN_CI_WRITE_TIMEOUT", time.Second*5)
	IdleTimeout = getDurationFromEnv("PLAIN_CI_READ_TIMEOUT", time.Second*120)
	GracefulShutdownPeriod = getDurationFromEnv("PLAIN_CI_SHUTDOWN_TIMEOUT", time.Second*6)

	IdleTimeout = getDurationFromEnv("PLAIN_CI_READ_TIMEOUT", time.Second*120)

	logLevel := getEnv("PLAIN_CI_LOG_LEVEL", "debug")
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

func getDurationFromEnv(key string, def time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		iVal, err := strconv.Atoi(value)
		if err != nil {
			logrus.Warnf("could not parse duration from %s: %v", key, err)
			return def
		} else {
			return time.Duration(iVal) * time.Second
		}
	}

	return def
}
