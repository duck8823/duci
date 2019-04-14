package cmd

import (
	"github.com/duck8823/duci/application"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "duci"}

func init() {
	rootCmd.AddCommand(serverCmd, runCmd, configCmd, healthCmd, versionCmd, updateCmd)
}

// Execute command
func Execute(args []string) {
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
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
		logrus.Fatalf("Failed to set configuration.\n%+v", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
