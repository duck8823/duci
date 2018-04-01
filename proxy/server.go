package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type server struct {
	handlerMap map[string]func(http.ResponseWriter, *http.Request)
}

func New() *server {
	return &server{make(map[string]func(http.ResponseWriter, *http.Request))}
}

func (s *server) Register(pattern string, url string, convertFunc func(incomingPayload io.Reader) (interface{}, error)) {
	s.handlerMap[pattern] = func(w http.ResponseWriter, r *http.Request) {
		outgoingPayload, err := convertFunc(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req, err := json.Marshal(outgoingPayload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res, err := http.Post(url, "application/json", bytes.NewReader(req))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%s", body)

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) Start() {
	for pattern, handleFunc := range s.handlerMap {
		http.HandleFunc(pattern, handleFunc)
	}
	http.ListenAndServe(":8080", nil)
}
