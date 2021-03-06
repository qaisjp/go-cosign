package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/qaisjp/gosign"
)

func newGosignClient(cfg CoSignConfig) (*gosign.Client, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read certfile+keyfile")
	}

	// Read CAFile containing multiple certs
	certs, err := ioutil.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read CAFile")
	}

	// Build a cert pool based from the CAFile
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)

	filter, err := gosign.Dial(&gosign.Config{
		Host:    cfg.DaemonHost,
		Port:    cfg.DaemonPort,
		Service: cfg.Service,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
			ServerName:         cfg.ServerName,
			Certificates:       []tls.Certificate{cert},
			RootCAs:            pool,
		},
	})
	if err != nil {
		return nil, err
	}

	return filter, nil
}
