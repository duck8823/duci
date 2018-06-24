package main

import (
	"github.com/duck8823/minimal-ci/presentation/controller"
	"github.com/google/logger"
	"net/http"
	"os"
)

func main() {
	logger.Init("minimal_ci", false, false, os.Stdout)

	ctrl, err := controller.New()
	if err != nil {
		logger.Fatalf("Failed to create controller.\n%+v", err)
		return
	}

	http.Handle("/", ctrl)

	http.ListenAndServe(":8080", nil)
}
