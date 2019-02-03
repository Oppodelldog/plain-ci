package server

import (
	"encoding/binary"
	"encoding/json"
	"github.com/Oppodelldog/simpleci/build"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func abortBuild(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id, hasId := vars["id"]
	if !hasId {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := build.AbortBuild(id)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	return

}

func getBuildQueue(writer http.ResponseWriter, request *http.Request) {
	builds := build.GetBuildQueueList()
	writeJson(writer, builds)
}

func getAllBuilds(writer http.ResponseWriter, request *http.Request) {
	builds := build.GetAllBuilds()
	writeJson(writer, builds)
}

func getBuild(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	buildId, hasBuildId := vars["buildId"]
	if !hasBuildId {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	builds := build.GetBuild(buildId)
	writeJson(writer, builds)

}

func getBuildLog(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	buildId, hasBuildId := vars["buildId"]
	logId, hasLogId := vars["logId"]

	if !hasBuildId || !hasLogId {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	buildLog := build.GetBuildLog(buildId, logId)
	writeJson(writer, buildLog)
}

func writeJson(writer http.ResponseWriter, data interface{}) {
	if binary.Size(data) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("error while encoding to json: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		logrus.Errorf("error writing json to client: %v", err)
	}
}
