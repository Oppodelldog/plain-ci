package server

import (
	"net/http"

	"github.com/Oppodelldog/plainci/webview/assets"
	"github.com/gorilla/mux"
)

func newRouter(queue Queue, aborter BuildAborter) *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/healthcheck", healthcheck).Methods(http.MethodGet)

	m.HandleFunc("/hook/simple", createQueueEnabledHandlerFunc(queue, hookSimple)).Methods(http.MethodPost)
	m.HandleFunc("/hook/github", createQueueEnabledHandlerFunc(queue, hookGithub)).Methods(http.MethodPost)

	m.HandleFunc("/queue", getBuildQueue).Methods(http.MethodGet)
	m.HandleFunc("/queue/abort/{id}", createAbortEnabledHandlerFunc(aborter, abortBuild)).Methods(http.MethodPost)

	m.HandleFunc("/build", getAllBuilds).Methods(http.MethodGet)
	m.HandleFunc("/build/{buildId}", getBuild).Methods(http.MethodGet)
	m.HandleFunc("/build/{buildId}/{logId}", getBuildLog).Methods(http.MethodGet)

	m.HandleFunc("/webview", webViewIndex).Methods(http.MethodGet)
	m.HandleFunc("/webview/queue", webViewQueue).Methods(http.MethodGet)
	m.HandleFunc("/webview/queue/abort/{id}", createAbortEnabledHandlerFunc(aborter, webviewAbort)).Methods(http.MethodGet)
	m.HandleFunc("/webview/project", webViewProjects).Methods(http.MethodGet)
	m.HandleFunc("/webview/project/{projectID}", webViewProject).Methods(http.MethodGet)
	m.HandleFunc("/webview/project/{projectID}/log/{logID}", webViewProjectLog).Methods(http.MethodGet)

	m.PathPrefix("/webview/images/").Handler(http.StripPrefix("/webview", http.FileServer(assets.Images.FS()))).Methods(http.MethodGet)

	return m
}
