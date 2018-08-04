package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application/service/log"
	"github.com/duck8823/duci/domain/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
		http.Error(w, fmt.Sprintf("Error occurred: %+v", err), http.StatusInternalServerError)
		return
	}

	var read int
	var job *model.Job
	for true {
		job, err = c.LogService.Get(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error occurred: %+v", err), http.StatusInternalServerError)
			return
		}
		for _, msg := range job.Stream[read:] {
			json.NewEncoder(w).Encode(msg)
			flusher.Flush()
			read++
		}
		if job.Finished {
			break
		}
	}
}
