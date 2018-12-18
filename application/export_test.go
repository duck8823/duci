package application

import "github.com/tcnksm/go-latest"

type MaskString = maskString

func SetVersion(ver string) (reset func()) {
	tmp := version
	version = ver
	return func() {
		version = tmp
	}
}

func SetCheckResponse(chr *latest.CheckResponse) (reset func()) {
	tmp := checked
	checked = chr
	return func() {
		checked = tmp
	}
}

func CheckLatestVersion() {
	checkLatestVersion()
}

func TrimSuffix(tag string) string {
	return trimSuffix(tag)
}

func GetCtxKey() string {
	return ctxKey
}