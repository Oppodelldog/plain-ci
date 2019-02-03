package server

import (
	"context"
	"fmt"
	"github.com/Oppodelldog/simpleci/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func startHttpServer(ctx context.Context) error {

	addr := fmt.Sprintf("0.0.0.0:%v", config.ServerPort)
	logrus.Infof("binding http on %v", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      http.HandlerFunc(redirectHandler),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		<-ctx.Done()
		err := srv.Close()
		if err != nil {
			logrus.Errorf("error closing http server: %v", err)
		}
	}()

	return srv.ListenAndServe()
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	host = strings.Replace(host, config.ServerPort, config.ServerPortTLS, 1)
	newURI := "https://" + host + r.URL.String()
	http.Redirect(w, r, newURI, http.StatusFound)
}
