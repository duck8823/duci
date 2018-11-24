package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
	"os"
)

var (
	version  = "dev"
	revision = "unknown"
)
var versionCmd = createCmd("version", "Display version", displayVersion)

func displayVersion(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)

	fmt.Printf("Version: %s (%s)\n", version, revision)
	res, _ := latest.Check(&latest.GithubTag{Owner: "duck8823", Repository: "duci"}, version)
	if res != nil && res.Outdated {
		msg := fmt.Sprintf("%s is not latest, you should upgrade to v%s", version, res.Current)
		println(msg)
		os.Exit(0)
		return
	}
}
