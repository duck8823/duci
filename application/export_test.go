package application

import "github.com/tcnksm/go-latest"

type MaskString = maskString

func SetVersion(ver string) {
	version = ver
}

func SetCheckResponse(chr *latest.CheckResponse) {
	checked = chr
}

func CheckLatestVersion() {
	checkLatestVersion()
}

func TrimSuffix(tag string) string {
	return trimSuffix(tag)
}
