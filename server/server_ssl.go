package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Oppodelldog/simpleci/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
)

func startHttpsServer(ctx context.Context) error {

	m := newRouter()

	addr := fmt.Sprintf("0.0.0.0:%v", config.ServerPortTLS)
	logrus.Infof("Staring https server on %v", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      m,
		TLSConfig:    getTlsConfig(),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	certFile := path.Join(config.CertDir, config.CertFile)
	keyFile := path.Join(config.CertDir, config.KeyFile)

	go func() {
		<-ctx.Done()
		err := srv.Close()
		if err != nil {
			logrus.Errorf("error closing https server: %v", err)
		}
	}()

	return srv.ListenAndServeTLS(certFile, keyFile)
}

func getTlsConfig() *tls.Config {

	return &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
