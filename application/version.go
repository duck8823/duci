package application

import (
	"fmt"
	"github.com/tcnksm/go-latest"
)

var (
	version  = "dev"
	revision = "unknown"
	checked  = &latest.CheckResponse{Outdated: false, Current: "dev"}
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

// IsOutdatedVersion return witch outdated version or not
func IsOutdatedVersion() bool {
	return checked.Outdated
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
