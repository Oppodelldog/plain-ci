package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oppodelldog/plainci/build"
	"github.com/Oppodelldog/plainci/server"
)

func main() {
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	queue := build.NewQueue(ctx)
	serverShutDown := server.Start(ctx, queue)

	<-sigChannel
	cancel()
	<-queue.GetShutDownChannel()
	<-serverShutDown
}
