package main

import (
	"github.com/Oppodelldog/simpleci/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := server.Start(); err != nil {
		logrus.Errorf("error running the server: %v", err)
	}
}
