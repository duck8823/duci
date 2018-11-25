package application_test

import (
	"github.com/duck8823/duci/application"
	"github.com/google/go-github/github"
	"github.com/tcnksm/go-latest"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestVersionString(t *testing.T) {
	// where
	for _, tt := range []struct {
		version  string
		expected string
	}{
		{"hoge", "hoge"},
		{"foo", "foo"},
	} {
		// given
		application.SetVersion(tt.version)

		// when
		actual := application.VersionString()

		// then
		if actual != tt.expected {
			t.Errorf("wont '%+v', but got '%+v'", tt.expected, actual)
		}
	}
}

func TestIsOutdatedVersion(t *testing.T) {
	// where
	for _, tt := range []struct {
		latest   bool
		expected bool
	}{
		{true, true},
		{false, false},
	} {
		// given
		application.SetCheckResponse(&latest.CheckResponse{Latest: tt.latest})

		// when
		actual := application.IsLatestVersion()

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

func TestCheckLatestVersion(t *testing.T) {
	// given
	expected := "0.0.2"
	application.SetVersion("0.0.1")

	// and
	gock.New("https://api.github.com").
		Get("/repos/duck8823/duci/tags").
		Reply(200).
		JSON([]github.RepositoryTag{
			{Name: github.String("0.0.1")},
			{Name: github.String(expected)},
		})

	// when
	application.CheckLatestVersion()

	// and
	actual := application.CurrentVersion()

	// then
	if actual != expected {
		t.Errorf("wont %+v, but got %+v", expected, actual)
	}
}

func TestTrimSuffix(t *testing.T) {
	// where
	for _, tt := range []struct {
		tag      string
		expected string
	}{
		{"v0.0.1", "v0.0.1"},
		{"v0.0.1-suffix", "v0.0.1"},
		{"-suffix", ""},
		{"v0.0.1-tag-suffix", "v0.0.1"},
	} {
		// when
		actual := application.TrimSuffix(tt.tag)

		// then
		if actual != tt.expected {
			t.Errorf("wont %+v, but got %+v", tt.expected, actual)
		}
	}
}
