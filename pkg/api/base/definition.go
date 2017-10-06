package base

import (
	"context"
	"net/http"

	"github.com/compsoc-edinburgh/bi-dice-api/pkg/config"
	cosign "github.com/qaisjp/go-cosign"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Config *config.Config
	Log    *logrus.Logger
	Filter *cosign.Filter
	Gin    *gin.Engine

	Server *http.Server
}

// Start binds the API and starts listening.
func (a *API) Start() error {
	a.Server = &http.Server{
		Addr:    a.Config.Address,
		Handler: a.Gin,
	}
	return a.Server.ListenAndServe()
}

// Shutdown shuts down the API
func (a *API) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
