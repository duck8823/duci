package health

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/pkg/errors"
	"net/http"
)

type handler struct {
	docker docker.Docker
}

// NewHandler returns implement of health check handler
func NewHandler() (http.Handler, error) {
	docker, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &handler{docker: docker}, nil
}

// ServeHTTP responses a server status
func (c *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := c.docker.Status(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
