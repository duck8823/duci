package cmd

import (
	"github.com/duck8823/duci/application"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var versionCmd = createCmd("version", "Display version", displayVersion)

func displayVersion(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	logrus.Infof("Version: %s", application.VersionString())
	if application.IsLatestVersion() {
		os.Exit(0)
		return
	}

	logrus.Warnf(
		"%s is not latest, you should upgrade to v%s",
		application.VersionString(),
		application.CurrentVersion(),
	)
}
