package server

import "github.com/gorilla/mux"

func newRouter() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/hook/simple", hookSimple)
	m.HandleFunc("/hook/github", hookGithub)

	m.HandleFunc("/queue", getBuildQueue)
	m.HandleFunc("/queue/abort/{id}", abortBuild)

	m.HandleFunc("/build", getAllBuilds)
	m.HandleFunc("/build/{buildId}", getBuild)
	m.HandleFunc("/build/{buildId}/{logId}", getBuildLog)

	return m
}
