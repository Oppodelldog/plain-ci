package server

import (
	"context"
	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start() {

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	httpWait := make(chan bool)
	httpServer := startHttpServer(httpWait)

	httpsWait := make(chan bool)
	httpsServer := startHttpsServer(httpsWait)

	<-sigChannel

	go gracefulShutdown(httpsServer)
	go gracefulShutdown(httpServer)

	<-httpWait
	<-httpsWait

	logrus.Info("Servers have stopped - shutdown.")
}

func gracefulShutdown(server *http.Server) {
	ctxShutDown, _ := context.WithTimeout(context.Background(), config.GracefulShutdownPeriod)
	err := server.Shutdown(ctxShutDown)
	if err != nil {
		logrus.Errorf("error shutting down server (%s): %v", server.Addr, err)
		err = server.Close()
		if err != nil {
			logrus.Errorf("error closing server (%s): %v", server.Addr, err)
		}
	}
}
