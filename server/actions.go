package server

import (
	"encoding/binary"
	"encoding/json"
	"github.com/Oppodelldog/simpleci/build"
	"github.com/sirupsen/logrus"
	"net/http"
)

func abortBuild(writer http.ResponseWriter, request *http.Request) {

	requestParms := newRequestParams(request)
	if !requestParms.Require("id").Validate() {
		writeBadRequest(writer, requestParms)
		return
	}

	err := build.AbortBuild(requestParms.GetString("id"))
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

	requestParms := newRequestParams(request)
	if !requestParms.Require("buildId").Validate() {
		writeBadRequest(writer, requestParms)
		return
	}

	builds := build.GetBuild(requestParms.GetString("buildId"))
	writeJson(writer, builds)
}

func getBuildLog(writer http.ResponseWriter, request *http.Request) {

	requestParms := newRequestParams(request)
	if !requestParms.Require("buildId", "logId").Integer("logId").Validate() {
		writeBadRequest(writer, requestParms)
		return
	}

	buildLog := build.GetBuildLog(requestParms.GetString("buildId"), requestParms.GetInt("logId"))
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

func writeBadRequest(writer http.ResponseWriter, requestParms *RequestParams) {
	writer.WriteHeader(http.StatusBadRequest)
	writeJson(writer, requestParms.GetValidationErrors())
}
