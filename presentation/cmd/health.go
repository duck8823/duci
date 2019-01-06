package cmd

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var healthCmd = createCmd("health", "Health check", healthCheck)

func healthCheck(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	cli, err := docker.New()
	if err != nil {
		logrus.Fatalf("Failed to set configuration.\n%+v", err)
	}

	if err := cli.Status(); err != nil {
		logrus.Fatalf("Unhealth: %s", err)
	} else {
		logrus.Info("ok")
	}
}
