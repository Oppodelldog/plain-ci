package server

import (
	"net/http"

	"github.com/Oppodelldog/plainci/build"
	"github.com/Oppodelldog/plainci/webview"
	"github.com/sirupsen/logrus"
)

func abortBuild(aborter BuildAborter, writer http.ResponseWriter, request *http.Request) {

	requestParams := newRequestParams(request)
	if !requestParams.Require("id").Validate() {
		writeBadRequestJson(writer, requestParams)
		return
	}

	err := aborter.AbortBuild(requestParams.GetString("id"))
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

	requestParams := newRequestParams(request)
	if !requestParams.Require("buildId").Validate() {
		writeBadRequestJson(writer, requestParams)
		return
	}

	builds := build.GetBuildLogs(requestParams.GetString("buildId"))
	writeJson(writer, builds)
}

func getBuildLog(writer http.ResponseWriter, request *http.Request) {

	requestParams := newRequestParams(request)
	if !requestParams.Require("buildId", "logId").Integer("logId").Validate() {
		writeBadRequestJson(writer, requestParams)
		return
	}

	buildLog := build.GetBuildLog(requestParams.GetString("buildId"), requestParams.GetInt("logId"))
	writeJson(writer, buildLog)
}

func webViewIndex(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	err := webview.RenderIndexPage(writer)
	if err != nil {
		logrus.Errorf("error rendering about page: %v", err)
	}
}

func webViewQueue(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	err := webview.RenderQueuePage(writer)
	if err != nil {
		logrus.Errorf("error rendering queue page: %v", err)
	}
}

func webViewProjectLog(writer http.ResponseWriter, request *http.Request) {
	requestParams := newRequestParams(request)
	if !requestParams.Require("logID", "projectID").Integer("logID").Validate() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	err := webview.RenderLogPage(writer, requestParams.GetString("projectID"), requestParams.GetInt("logID"))
	if err != nil {
		logrus.Errorf("error rendering log page: %v", err)
	}
}

func webviewAbort(aborter BuildAborter, writer http.ResponseWriter, request *http.Request) {
	requestParams := newRequestParams(request)
	if !requestParams.Require("id").Validate() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err := webview.RenderAbortPage(aborter.AbortBuild, writer, requestParams.GetString("id"))
	if err != nil {
		logrus.Errorf("error rendering abort page: %v", err)
	}
}

func webViewProjects(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	err := webview.RenderProjectsPage(writer)
	if err != nil {
		logrus.Errorf("error rendering projects page: %v", err)
	}
}

func webViewProject(writer http.ResponseWriter, request *http.Request) {
	requestParams := newRequestParams(request)
	if !requestParams.Require("projectID").Validate() {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	err := webview.RenderBuildPage(writer, requestParams.GetString("projectID"))
	if err != nil {
		logrus.Errorf("error rendering build page: %v", err)
	}
}
