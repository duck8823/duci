package router

import (
	"github.com/duck8823/duci/presentation/controller/health"
	"github.com/duck8823/duci/presentation/controller/job"
	"github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// New returns handler of application.
func New() (http.Handler, error) {
	webhookHandler, err := webhook.NewHandler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jobHandler, err := job.NewHandler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	healthHandler, err := health.NewHandler()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rtr := chi.NewRouter()
	rtr.Post("/", webhookHandler.ServeHTTP)
	rtr.Get("/logs/{uuid}", jobHandler.ServeHTTP)
	rtr.Get("/health", healthHandler.ServeHTTP)

	return rtr, nil
}
