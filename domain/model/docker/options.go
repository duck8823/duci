package docker

import (
	"fmt"
	"strings"
)

// RuntimeOptions is a docker options.
type RuntimeOptions struct {
	Environments Environments
	Volumes      Volumes
}

// Environments represents a docker `-e` option.
type Environments map[string]interface{}

// ToArray returns string array of environments
func (e Environments) ToArray() []string {
	var a []string
	for key, val := range e {
		a = append(a, fmt.Sprintf("%s=%v", key, val))
	}
	return a
}

// Volumes represents a docker `-v` option.
type Volumes []string

// ToMap returns map of volumes.
func (v Volumes) ToMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, volume := range v {
		key := strings.Split(volume, ":")[0]
		m[key] = struct{}{}
	}
	return m
}
