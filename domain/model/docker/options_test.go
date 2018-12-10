package docker_test

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/google/go-cmp/cmp"
	"sort"
	"testing"
)

func TestEnvironments_ToArray(t *testing.T) {
	var empty []string
	for _, tt := range []struct {
		in   docker.Environments
		want []string
	}{
		{
			in:   docker.Environments{},
			want: empty,
		},
		{
			in: docker.Environments{
				"int":    19,
				"string": "hello",
			},
			want: []string{
				"int=19",
				"string=hello",
			},
		},
	} {
		// when
		got := tt.in.Array()
		want := tt.want
		sort.Strings(got)
		sort.Strings(want)

		// then
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal. but %+v", cmp.Diff(got, want))
		}
	}
}

func TestVolumes_Volumes(t *testing.T) {
	for _, tt := range []struct {
		in   docker.Volumes
		want map[string]struct{}
	}{
		{
			in:   docker.Volumes{},
			want: make(map[string]struct{}),
		},
		{
			in: docker.Volumes{
				"/hoge/fuga:/hoge/hoge",
			},
			want: map[string]struct{}{
				"/hoge/fuga": {},
			},
		},
	} {
		// when
		got := tt.in.Map()
		want := tt.want

		// then
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal. but %+v", cmp.Diff(got, want))
		}
	}
}
