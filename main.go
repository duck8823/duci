package main

import (
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func main() {
	ctrl, err := controller.New()
	if err != nil {
		logger.Errorf(uuid.UUID{}, "Failed to create controller.\n%+v", err)
		os.Exit(1)
		return
	}

	http.Handle("/", ctrl)

	http.ListenAndServe(":8080", nil)
}
