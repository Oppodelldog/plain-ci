package server

import (
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

	return m
}
