package main

import (
	"flag"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/router"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func main() {
	mainID := uuid.New()

	flag.Var(application.Config, "c", "configuration file path")
	flag.Parse()

	if err := semaphore.Make(); err != nil {
		logger.Errorf(mainID, "Failed to initialize a semaphore.\n%+v", err)
		os.Exit(1)
		return
	}

	rtr, err := router.New()
	if err != nil {
		logger.Errorf(mainID, "Failed to initialize controllers.\n%+v", err)
		os.Exit(1)
		return
	}

	if err := http.ListenAndServe(application.Config.Addr(), rtr); err != nil {
		logger.Errorf(mainID, "Failed to run server.\n%+v", err)
		os.Exit(1)
		return
	}
}
