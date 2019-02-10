package server

import (
	"context"
	"net/http"

	"github.com/Oppodelldog/plainci/build"

	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, queue *build.Queue) chan bool {
	shutdownChannel := make(chan bool)

	go func() {

		httpWait := make(chan bool)
		httpServer := startHttpServer(httpWait)

		httpsWait := make(chan bool)
		httpsServer := startHttpsServer(queue, httpsWait)

		<-ctx.Done()

		go gracefulShutdown(httpsServer)
		go gracefulShutdown(httpServer)

		<-httpWait
		<-httpsWait

		logrus.Info("Servers have stopped - shutdown.")
		close(shutdownChannel)
	}()

	return shutdownChannel
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
