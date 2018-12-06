package controller

import (
	"github.com/duck8823/duci/application/service/docker"
	"net/http"
)

// HealthController is a handler of health check.
type HealthController struct {
	Docker docker.Service
}

// ServeHTTP responses a server status
func (c *HealthController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := c.Docker.Status(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
