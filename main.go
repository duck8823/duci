package main

import (
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/presentation/controller"
	"net/http"
)

func main() {
	logger.Init()

	ctrl, err := controller.New()
	if err != nil {
		logger.Fatalf("Failed to create controller.\n%+v", err)
		return
	}

	http.Handle("/", ctrl)

	http.ListenAndServe(":8080", nil)
}
