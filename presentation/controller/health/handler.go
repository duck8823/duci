package health

import (
	"github.com/duck8823/duci/domain/model/docker"
	"net/http"
)

// Handler of health check.
type Handler struct {
	Docker docker.Docker
}

// ServeHTTP responses a server status
func (c *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := c.Docker.Status(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
