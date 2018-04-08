package handlers

import (
	"encoding/json"
	"github.com/duck8823/webhook-proxy/payloads"
	"github.com/google/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

type SlackNotificator struct {
	Url         string
	ConvertFunc func(body []byte) (*payloads.SlackMessage, error)
}

func (s *SlackNotificator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body."))
		return
	}

	message, err := s.ConvertFunc(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to convert message."))
		return
	}
	jsonStr, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to convert message."))
		return
	}

	resp, err := http.Post(s.Url, "application/json", strings.NewReader(string(jsonStr)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to send message."))
		return
	}

	slackResponse, err := ioutil.ReadAll(resp.Body)
	logger.Infof("%s", slackResponse)
}
