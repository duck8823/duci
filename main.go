package main

import (
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
		Run: serverCmd,
	}
	serverCmd.PersistentFlags().VarP(application.Config, "config", "c", "configuration file path")

	rootCmd := &cobra.Command{Use: "duci"}
	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf(uuid.New(), "Failed to execute command.\n%+v", err)
		os.Exit(1)
	}
}

func serverCmd(cmd *cobra.Command, _ []string) {
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
