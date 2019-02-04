package server

import (
	"github.com/Oppodelldog/simpleci/webview/assets"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/healthcheck", healthcheck).Methods(http.MethodGet)

	m.HandleFunc("/hook/simple", hookSimple).Methods(http.MethodPost)
	m.HandleFunc("/hook/github", hookGithub).Methods(http.MethodPost)

	m.HandleFunc("/queue", getBuildQueue).Methods(http.MethodGet)
	m.HandleFunc("/queue/abort/{id}", abortBuild).Methods(http.MethodPost)

	m.HandleFunc("/build", getAllBuilds).Methods(http.MethodGet)
	m.HandleFunc("/build/{buildId}", getBuild).Methods(http.MethodGet)
	m.HandleFunc("/build/{buildId}/{logId}", getBuildLog).Methods(http.MethodGet)

	m.HandleFunc("/webview", webViewIndex).Methods(http.MethodGet)
	m.HandleFunc("/webview/queue", webViewQueue).Methods(http.MethodGet)
	m.HandleFunc("/webview/queue/abort/{id}", webviewAbort).Methods(http.MethodGet)
	m.HandleFunc("/webview/build", webViewBuilds).Methods(http.MethodGet)
	m.HandleFunc("/webview/build/{buildId}", webViewBuild).Methods(http.MethodGet)
	m.HandleFunc("/webview/build/{buildId}/{logId}", webViewLog).Methods(http.MethodGet)

	m.PathPrefix("/webview/images/").Handler(http.StripPrefix("/webview", http.FileServer(assets.Images.FS()))).Methods(http.MethodGet)

	return m
}
