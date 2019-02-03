package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	httpWait := make(chan bool)
	httpsWait := make(chan bool)
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := startHttpServer(ctx); err != nil {
			logrus.Errorf("error running the server: %v", err)
		}
		close(httpWait)
	}()

	go func() {
		if err := startHttpsServer(ctx); err != nil {
			logrus.Errorf("error running the server: %v", err)
		}
		close(httpsWait)
	}()

	<-sigChannel

	cancel()

	<-httpWait
	<-httpsWait

	logrus.Info("Servers have stopped - shutdown.")
}
