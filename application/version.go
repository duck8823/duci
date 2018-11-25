package application

import (
	"github.com/tcnksm/go-latest"
	"regexp"
)

var (
	version            = "dev"
	versionSuffixRegex = regexp.MustCompile("-.+$")
	checked            = &latest.CheckResponse{Latest: true, Current: version}
)

func init() {
	checkLatestVersion()
}

// VersionString returns application version
func VersionString() string {
	return version
}

// IsLatestVersion return witch latest stable version or not
func IsLatestVersion() bool {
	return checked.Latest
}

// CurrentVersion returns current version string
func CurrentVersion() string {
	return checked.Current
}

func checkLatestVersion() {
	checkSrc := &latest.GithubTag{Owner: "duck8823", Repository: "duci", FixVersionStrFunc: trimSuffix}
	if res, err := latest.Check(checkSrc, trimSuffix(version)); err == nil {
		checked = res
	}
}

func trimSuffix(tag string) string {
	return versionSuffixRegex.ReplaceAllString(tag, "")
}
