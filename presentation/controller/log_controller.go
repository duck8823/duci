package controller

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/infrastructure/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type LogController struct{}

func (c *LogController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occurred: %+v", err), http.StatusInternalServerError)
		return
	}

	var read int
	var job *store.Job
	for true {
		job, err = store.Get(id)
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
