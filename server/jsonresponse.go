package server

import (
	"encoding/binary"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

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

func writeBadRequestJson(writer http.ResponseWriter, requestParms *RequestParams) {
	writer.WriteHeader(http.StatusBadRequest)
	writeJson(writer, requestParms.GetValidationErrors())
}
