package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
	"os"
)

var version string
var versionCmd = createCmd("version", "Display version", displayVersion)

func displayVersion(cmd *cobra.Command, _ []string) {
	readConfiguration(cmd)
	if len(version) == 0 {
		println("Version: unknown")
		os.Exit(0)
		return
	}
	fmt.Printf("Version: %s\n", version)
	res, _ := latest.Check(&latest.GithubTag{Owner: "duck8823", Repository: "duci"}, version)
	if res.Outdated {
		msg := fmt.Sprintf("%s is not latest, you should upgrade to v%s", version, res.Current)
		println(msg)
		os.Exit(0)
		return
	}
}
