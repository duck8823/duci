package cmd

import (
	"encoding/json"
	"github.com/duck8823/duci/application"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var configCmd = createCmd("config", "Display configuration", displayConfig)

func displayConfig(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(application.Config); err != nil {
		logrus.Fatalf("Failed to display config.\n%+v", err)
	}
}
