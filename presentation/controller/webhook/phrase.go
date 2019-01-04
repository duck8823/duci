package webhook

import (
	"github.com/duck8823/duci/domain/model/docker"
	"regexp"
	"strings"
)

type phrase string

// Command returns command of docker
func (p phrase) Command() docker.Command {
	return strings.Split(string(p), " ")
}

func extractBuildPhrase(comment string) (phrase, error) {
	if !regexp.MustCompile(`^ci\s+[^\\s]+`).Match([]byte(comment)) {
		return "", SkipBuild
	}
	phrase := phrase(regexp.MustCompile(`^ci\s+`).ReplaceAllString(comment, ""))
	return phrase, nil
}
