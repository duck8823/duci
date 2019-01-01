package cmd

import (
	"fmt"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/spf13/cobra"
	"os"
)

var healthCmd = createCmd("health", "Health check", healthCheck)

func healthCheck(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	cli, err := docker.New()
	if err != nil {
		msg := fmt.Sprintf("Failed to set configuration.\n%+v", err)
		if _, err := fmt.Fprint(os.Stderr, msg); err != nil {
			println(err)
		}
		os.Exit(1)
	}

	if err := cli.Status(); err != nil {
		println(fmt.Sprintf("Unhealth\n%s", err))
		os.Exit(1)
	} else {
		println("ok")
		os.Exit(0)
	}
}
