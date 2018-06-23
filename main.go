package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duck8823/minimal-ci/service/runner"
	"github.com/google/go-github/github"
	"github.com/google/logger"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	logger.Init("minimal_ci", false, false, os.Stdout)

	job, err := runner.NewWithEnv()
	if err != nil {
		logger.Fatalf("Failed to create runner.\n%+v", err)
		os.Exit(-1)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read Payload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		event := &github.IssueCommentEvent{}
		if err := json.Unmarshal(body, event); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Trigger build
		githubEvent := r.Header.Get("X-GitHub-Event")
		if githubEvent != "issue_comment" {
			message := fmt.Sprintf("payload event type must be issue_comment. but %s", githubEvent)
			logger.Error(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(event.Comment.GetBody())) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not build."))
			return
		}
		phrase := regexp.MustCompile("^ci\\s+").ReplaceAllString(event.Comment.GetBody(), "")

		if err := job.Run(context.Background(), event.GetRepo(), "master", phrase); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Response
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", nil)
}
