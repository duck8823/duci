package cmd

import (
	"github.com/duck8823/duci/application/service/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
)

var healthCmd = createCmd("health", "Health check", healthCheck)

func healthCheck(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	dockerService, err := docker.New()
	if err != nil {
		logger.Errorf(uuid.New(), "Failed to set configuration.\n%+v", err)
		os.Exit(1)
	}

	if err := dockerService.Status(); err != nil {
		logger.Errorf(uuid.New(), "Unhealthy.\n%s", err)
		os.Exit(1)
	} else {
		logger.Info(uuid.New(), "ok.")
		os.Exit(0)
	}
}
