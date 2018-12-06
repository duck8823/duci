package job

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application/service/job"
	. "github.com/duck8823/duci/domain/model/job"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type Handler struct {
	Service job_service.Service
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := h.logs(w, ID(id)); err != nil {
		http.Error(w, fmt.Sprintf(" Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) logs(w http.ResponseWriter, id ID) error {
	f, ok := w.(http.Flusher)
	if !ok {
		return errors.New("Streaming unsupported!")
	}

	var read int
	var job *Job
	var err error
	for {
		job, err = h.Service.FindBy(id)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, msg := range job.Stream[read:] {
			json.NewEncoder(w).Encode(msg)
			f.Flush()
			read++
		}
		if job.Finished {
			break
		}
	}
	return nil
}
