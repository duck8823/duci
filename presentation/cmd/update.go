package cmd

import (
	"github.com/blang/semver"
	"github.com/duck8823/duci/application"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = createCmd("update", "Update binary", doUpdate)

func doUpdate(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	v, err := semver.ParseTolerant(application.VersionString())
	if err != nil {
		logrus.Fatalf("Error occurred: %+v", err)
	}

	latest, err := selfupdate.UpdateSelf(v, "duck8823/duci")
	if err != nil {
		logrus.Fatalf("Binary update failed: %+v", err)
	}
	if latest.Version.Equals(v) {
		logrus.Infof("Current binary is the latest %s", application.VersionString())
	} else {
		logrus.Infof("Successfully updated to %s", application.CurrentVersion())
	}
}
