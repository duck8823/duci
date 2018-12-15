package github_test

import (
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDescription_TrimmedString(t *testing.T) {
	// where
	for _, tt := range []struct {
		in   string
		want string
	}{
		{
			in:   "hello world",
			want: "hello world",
		},
		{
			in:   "123456789012345678901234567890123456789012345678901234567890",
			want: "12345678901234567890123456789012345678901234567...",
		},
		{
			in:   "12345678901234567890123456789012345678901234567890",
			want: "12345678901234567890123456789012345678901234567890",
		},
		{
			in:   "123456789012345678901234567890123456789012345678901",
			want: "12345678901234567890123456789012345678901234567...",
		},
	} {
		// given
		sut := github.Description(tt.in)

		// when
		got := sut.TrimmedString()

		// then
		if got != tt.want {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, tt.want))
		}
	}
}
