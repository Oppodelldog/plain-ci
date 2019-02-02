package server

import (
	"encoding/json"
	"github.com/Oppodelldog/simpleci/build"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SimpleHookPayload struct {
	URL string
	REV string
}

func hookSimple(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload SimpleHookPayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		logrus.Errorf("error decoding request payload: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	go build.New(payload.URL, payload.REV, "hook:simple", nil, nil)

	writer.WriteHeader(http.StatusNoContent)
}
