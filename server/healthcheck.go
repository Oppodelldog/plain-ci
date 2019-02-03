package server

import (
	"net/http"
)

func healthcheck(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("simple-ci"))
}
