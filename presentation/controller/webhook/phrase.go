package webhook

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type phrase string

func (p phrase) Command() docker.Command {
	return strings.Split(string(p), " ")
}

func extractBuildPhrase(comment string) (phrase, SkipBuild) {
	if !regexp.MustCompile("^ci\\s+[^\\s]+").Match([]byte(comment)) {
		return "", SkipBuild(errors.New("Not start with ci."))
	}
	phrase := phrase(regexp.MustCompile("^ci\\s+").ReplaceAllString(comment, ""))
	return phrase, nil
}
