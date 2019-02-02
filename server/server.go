package server

import (
	"fmt"
	"github.com/Oppodelldog/simpleci/config"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Start() error {

	m := newRouter()

	addr := fmt.Sprintf("0.0.0.0:%v", config.ServerPort)
	logrus.Infof("Staring server on %v", addr)

	return http.ListenAndServe(addr, m)
}
