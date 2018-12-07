package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/spf13/cobra"
	"os"
)

var configCmd = createCmd("config", "Display configuration", displayConfig)

func displayConfig(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(application.Config); err != nil {
		println(fmt.Sprintf("Failed to display config.\n%+v", err))
		os.Exit(1)
	}
}
