package main

import (
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/presentation/controller"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"github.com/google/uuid"
)

func main() {
	logger.Init(os.Stdout, logging.DEBUG)

	ctrl, err := controller.New()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to create controller.\n%+v", err)
		os.Exit(1)
		return
	}

	http.Handle("/", ctrl)

	http.ListenAndServe(":8080", nil)
}
