package server

import "github.com/gorilla/mux"

func newRouter() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/hook/simple", hookSimple)
	m.HandleFunc("/hook/github", hookGithub)
	m.HandleFunc("/queue", builds)
	m.HandleFunc("/queue/abort/{id}", abortBuild)
	return m
}
