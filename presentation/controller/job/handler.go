package job

import (
	"encoding/json"
	"fmt"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type handler struct {
	service jobService.Service
}

// NewHandler returns implement of job
func NewHandler() (http.Handler, error) {
	service, err := jobService.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &handler{service: service}, nil
}

// ServeHTTP responses log stream
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := h.logs(w, job.ID(id)); err != nil {
		http.Error(w, fmt.Sprintf(" Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *handler) logs(w http.ResponseWriter, id job.ID) error {
	f, ok := w.(http.Flusher)
	if !ok {
		return errors.New("Streaming unsupported")
	}

	// TODO: add timeout
	var read int
	for {
		job, err := h.service.FindBy(id)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, msg := range job.Stream[read:] {
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				logrus.Errorf("%+v", err)
			}
			f.Flush()
			read++
		}
		if job.Finished {
			break
		}
	}
	return nil
}
