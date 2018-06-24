package docker_test

import (
	"github.com/duck8823/minimal-ci/infrastructure/docker"
	"reflect"
	"testing"
)

func TestEnvironments_ToArray(t *testing.T) {
	var empty []string
	for _, testcase := range []struct {
		in       docker.Environments
		expected []string
	}{
		{
			in:       docker.Environments{},
			expected: empty,
		},
		{
			in: docker.Environments{
				"int":    19,
				"string": "hello",
			},
			expected: []string{
				"int=19",
				"string=hello",
			},
		},
	} {
		actual := testcase.in.ToArray()
		if !reflect.DeepEqual(actual, testcase.expected) {
			t.Errorf("must be equal. actual=%+v, wont=%+v", actual, testcase.expected)
		}
	}
}
