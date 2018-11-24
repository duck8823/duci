package application

import (
	"fmt"
	"github.com/tcnksm/go-latest"
)

var (
	version  = "dev"
	revision = "unknown"
	checked  = &latest.CheckResponse{Latest: true, Current: "dev"}
)

func init() {
	checkLatestVersion()
}

// VersionString returns application version with revision (commit hash)
func VersionString() string {
	return fmt.Sprintf("%s (%s)", version, revision)
}

// VersionStringShort returns application version
func VersionStringShort() string {
	return version
}

// IsLatestVersion return witch latest version or not
func IsLatestVersion() bool {
	return checked.Latest
}

// CurrentVersion returns current version string
func CurrentVersion() string {
	return checked.Current
}

func checkLatestVersion() {
	if res, err := latest.Check(&latest.GithubTag{Owner: "duck8823", Repository: "duci"}, version); err == nil {
		checked = res
	}
}
