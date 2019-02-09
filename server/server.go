package server

import (
	"fmt"
	"github.com/Oppodelldog/plainci/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func startHttpServer(wait chan bool) *http.Server {

	addr := fmt.Sprintf("%s:%s", config.ServerAddressBind, config.ServerPort)
	logrus.Infof("binding http on %v", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      http.HandlerFunc(redirectHandler),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("error running the server: %v", err)
			close(wait)
		}
	}()

	return srv
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	host = strings.Replace(host, config.ServerPort, config.ServerPortTLS, 1)
	newURI := "https://" + host + r.URL.String()
	http.Redirect(w, r, newURI, http.StatusFound)
}
