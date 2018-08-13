package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application/service/log"
	"github.com/duck8823/duci/domain/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type LogController struct {
	LogService log.StoreService
}

func (c *LogController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := c.logs(w, flusher, id); err != nil {
		http.Error(w, fmt.Sprintf("Sorry, Error occurred: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (c *LogController) logs(w http.ResponseWriter, f http.Flusher, id uuid.UUID) error {
	var read int
	var job *model.Job
	var err error
	for true {
		job, err = c.LogService.Get(id)
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
