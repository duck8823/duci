package cmd

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/presentation/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

var (
	serverCmd = createCmd("server", "Start server", runServer)
	logo      = `
 ___    __ __    __  ____ 
|   \  |  |  |  /  ]|    |
|    \ |  |  | /  /  |  | 
|  D  ||  |  |/  /   |  | 
|     ||  :  /   \_  |  | 
|     ||     \     | |  | 
|_____| \__,_|\____||____|
`
)

func runServer(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	if err := application.Initialize(); err != nil {
		logrus.Fatal(fmt.Sprintf("Failed to initialize a semaphore.\n%+v", err))
		return
	}

	if err := semaphore.Make(); err != nil {
		logrus.Fatal(fmt.Sprintf("Failed to initialize a semaphore.\n%+v", err))
		return
	}

	rtr, err := router.New()
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Failed to initialize controllers.\n%+v", err))
		return
	}

	for _, l := range strings.Split(logo, "\n") {
		logrus.Info(l)
	}
	if err := http.ListenAndServe(application.Config.Addr(), rtr); err != nil {
		logrus.Fatal(fmt.Sprintf("Failed to run server.\n%+v", err))
		return
	}
}
