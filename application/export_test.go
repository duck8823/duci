package application

import "github.com/tcnksm/go-latest"

type MaskString = maskString

func SetVersion(ver string) {
	version = ver
}

func SetRevision(rev string) {
	revision = rev
}

func SetCheckResponse(chr *latest.CheckResponse) {
	checked = chr
}
