package server

import (
	"encoding/json"
	"github.com/Oppodelldog/simpleci/build"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func abortBuild(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(request)
	if id, ok := vars["id"]; ok {
		err := build.AbortBuild(id)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
		return

	} else {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func builds(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	builds := build.GetBuilds()
	if len(builds) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(builds)
	if err != nil {
		logrus.Errorf("error while encoding builds to json: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		logrus.Errorf("error writing builds json to client: %v", err)
	}
}
