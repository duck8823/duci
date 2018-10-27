package main

import (
	"encoding/json"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/presentation/router"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

func main() {
	serverCmd := &cobra.Command{
		Use: "server",
		Short: "Start server",
		Run: serverCmd,
	}
	serverCmd.PersistentFlags().VarPF(application.Config, "config", "c", "configuration file path")

	configCmd := &cobra.Command{
		Use: "config",
		Short: "Display configuration",
		Run: configCmd,
	}
	configCmd.PersistentFlags().VarPF(application.Config, "config", "c", "configuration file path")

	rootCmd := &cobra.Command{Use: "duci"}
	rootCmd.AddCommand(serverCmd, configCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf(uuid.New(), "Failed to execute command.\n%+v", err)
		os.Exit(1)
	}
}

func serverCmd(_ *cobra.Command, _ []string) {
	mainID := uuid.New()

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

func configCmd(_ *cobra.Command, _ []string) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(application.Config)
}
