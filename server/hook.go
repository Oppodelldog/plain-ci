package server

import (
	"net/http"

	"github.com/Oppodelldog/plainci/build"
)

type Queue interface {
	NewBuild(repoURL string, commit string, originator string, customPreBuildNotify build.NotificationFunc, customPostBuildNotify build.NotificationFunc)
}

type queueEnabledHookFunc func(queue Queue, writer http.ResponseWriter, request *http.Request)

func enableQueueHook(queue Queue, f queueEnabledHookFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		f(queue, writer, request)
	}
}
