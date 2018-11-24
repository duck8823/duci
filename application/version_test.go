package application_test

import (
	"github.com/duck8823/duci/application"
	"github.com/tcnksm/go-latest"
	"testing"
)

func TestVersionString(t *testing.T) {
	// where
	for _, tt := range []struct {
		version  string
		revision string
		expected string
	}{
		{"hoge", "fuga", "hoge (fuga)"},
		{"foo", "bar", "foo (bar)"},
	} {
		// given
		application.SetVersion(tt.version)
		application.SetRevision(tt.revision)

		// when
		actual := application.VersionString()

		// then
		if actual != tt.expected {
			t.Errorf("wont '%+v', but got '%+v'", tt.expected, actual)
		}
	}
}

func TestVersionStringShort(t *testing.T) {
	// given
	expected := "dev"
	application.SetVersion(expected)

	// when
	actual := application.VersionStringShort()

	// then
	if actual != expected {
		t.Errorf("wont '%+v', but got '%+v'", expected, actual)
	}
}

func TestIsOutdatedVersion(t *testing.T) {
	// where
	for _, tt := range []struct {
		outdated bool
		expected bool
	}{
		{true, true},
		{false, false},
	} {
		// given
		application.SetCheckResponse(&latest.CheckResponse{Outdated: tt.outdated})

		// when
		actual := application.IsOutdatedVersion()

		// then
		if actual != tt.expected {
			t.Errorf("wont '%+v', but got '%+v'", tt.expected, actual)
		}
	}
}

func TestCurrentVersion(t *testing.T) {
	// where
	for _, tt := range []struct {
		current  string
		expected string
	}{
		{"hoge", "hoge"},
		{"fuga", "fuga"},
	} {
		// given
		application.SetCheckResponse(&latest.CheckResponse{Current: tt.current})

		// when
		actual := application.CurrentVersion()

		// then
		if actual != tt.expected {
			t.Errorf("wont '%+v', but got '%+v'", tt.expected, actual)
		}
	}
}
