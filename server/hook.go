package server

import (
	"net/http"

	"github.com/Oppodelldog/plainci/build"
)

type Queue interface {
	NewBuild(repoURL string, commit string, originator string, customPreBuildNotify build.NotificationFunc, customPostBuildNotify build.NotificationFunc)
}

type queueEnabledHandlerFunc func(queue Queue, writer http.ResponseWriter, request *http.Request)

func createQueueEnabledHandlerFunc(queue Queue, f queueEnabledHandlerFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		f(queue, writer, request)
	}
}

type BuildAborter interface {
	AbortBuild(id string) error
}

type abortEnabledHandlerFunc func(aborter BuildAborter, writer http.ResponseWriter, request *http.Request)

func createAbortEnabledHandlerFunc(aborter BuildAborter, f abortEnabledHandlerFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		f(aborter, writer, request)
	}
}
