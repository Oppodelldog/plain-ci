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

func getBuildQueue(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	builds := build.GetBuildQueueList()
	if len(builds) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(builds)
	if err != nil {
		logrus.Errorf("error while encoding getBuildQueue to json: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		logrus.Errorf("error writing getBuildQueue json to client: %v", err)
	}
}

func getAllBuilds(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	builds := build.GetAllBuilds()
	if len(builds) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(builds)
	if err != nil {
		logrus.Errorf("error while encoding all builds to json: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		logrus.Errorf("error writing all builds json to client: %v", err)
	}
}

func getBuild(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	if buildId, ok := vars["buildId"]; ok {

		builds := build.GetBuild(buildId)
		if len(builds) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		jsonBytes, err := json.Marshal(builds)
		if err != nil {
			logrus.Errorf("error while encoding build %s to json: %v", buildId, err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(jsonBytes)
		if err != nil {
			logrus.Errorf("error writing build %s json to client: %v", buildId, err)
		}
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func getBuildLog(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	buildId, hasBuildId := vars["buildId"]
	logId, hasLogId := vars["logId"]
	if !hasBuildId || !hasLogId {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	buildLog, err := build.GetBuildLog(buildId, logId)
	if err != nil {
		logrus.Errorf("error while reading build log %s %s: %v", buildId, logId, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(buildLog) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(struct{ Log string }{Log: buildLog})
	if err != nil {
		logrus.Errorf("error while encoding build log to json %s %s: %v", buildId, logId, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		logrus.Errorf("error writing build log json to client %s %s: %v", buildId, logId, err)
	}
}
