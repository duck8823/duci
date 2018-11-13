package main

import (
	"encoding/json"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/router"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func main() {
	var configPath string

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start server",
		Run:   serverCmd,
	}
	serverCmd.PersistentFlags().StringVarP(&configPath, "config", "c", application.DefaultConfigurationPath, "configuration file path")

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Display configuration",
		Run:   configCmd,
	}
	configCmd.PersistentFlags().StringVarP(&configPath, "config", "c", application.DefaultConfigurationPath, "configuration file path")

	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Health check",
		Run:   healthCmd,
	}
	healthCmd.PersistentFlags().StringVarP(&configPath, "config", "c", application.DefaultConfigurationPath, "configuration file path")

	rootCmd := &cobra.Command{Use: "duci"}
	rootCmd.AddCommand(serverCmd, configCmd, healthCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf(uuid.New(), "Failed to execute command.\n%+v", err)
		os.Exit(1)
	}
}

func serverCmd(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	if err := semaphore.Make(); err != nil {
		logger.Errorf(uuid.New(), "Failed to initialize a semaphore.\n%+v", err)
		os.Exit(1)
		return
	}

	rtr, err := router.New()
	if err != nil {
		logger.Errorf(uuid.New(), "Failed to initialize controllers.\n%+v", err)
		os.Exit(1)
		return
	}

	if err := http.ListenAndServe(application.Config.Addr(), rtr); err != nil {
		logger.Errorf(uuid.New(), "Failed to run server.\n%+v", err)
		os.Exit(1)
		return
	}
}

func configCmd(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(application.Config); err != nil {
		logger.Errorf(uuid.New(), "Failed to display config.\n%+v", err)
		os.Exit(1)
	}
}

func healthCmd(cmd *cobra.Command, _ []string) {
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
		logger.Info(uuid.New(), "Healthy.")
		os.Exit(0)
	}
}

func readConfiguration(cmd *cobra.Command) {
	configFilePath := cmd.Flag("config").Value.String()
	if !exists(configFilePath) && configFilePath == application.DefaultConfigurationPath {
		return
	}

	if err := application.Config.Set(configFilePath); err != nil {
		logger.Errorf(uuid.New(), "Failed to set configuration.\n%+v", err)
		os.Exit(1)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
