package cmd

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "duci"}

func init() {
	rootCmd.AddCommand(serverCmd, configCmd, healthCmd, versionCmd)
}

// Execute command
func Execute(args []string) {
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type command func(cmd *cobra.Command, args []string)

func createCmd(use string, short string, run command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Run:   run,
	}
	cmd.PersistentFlags().StringP("config", "c", application.DefaultConfigurationPath, "configuration file path")
	return cmd
}

func readConfiguration(cmd *cobra.Command) {
	configFilePath := cmd.Flag("config").Value.String()
	if !exists(configFilePath) && configFilePath == application.DefaultConfigurationPath {
		return
	}

	if err := application.Config.Set(configFilePath); err != nil {
		println(fmt.Sprintf("Failed to set configuration.\n%+v", err))
		os.Exit(1)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
