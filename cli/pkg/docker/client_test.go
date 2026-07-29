package docker

import (
	"regexp"
	"strings"
	"testing"
)

func TestSanitizeContainerNamePart(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "stuntdouble", want: "stuntdouble"},
		{input: "My Repo!", want: "my-repo"},
		{input: "...___---", want: "workspace"},
		{input: "repo@prod/db", want: "repo-prod-db"},
	}

	for _, tc := range tests {
		got := sanitizeContainerNamePart(tc.input)
		if got != tc.want {
			t.Fatalf("sanitizeContainerNamePart(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestSanitizeContainerNamePartTruncatesLongInput(t *testing.T) {
	got := sanitizeContainerNamePart(strings.Repeat("a", 80))
	if len(got) != 40 {
		t.Fatalf("len(sanitizeContainerNamePart(long)) = %d, want 40", len(got))
	}
}

func TestNewSidecarNameIsDockerSafeAndUnique(t *testing.T) {
	first, err := newSidecarName("/tmp/My Repo!")
	if err != nil {
		t.Fatalf("newSidecarName() error = %v", err)
	}
	second, err := newSidecarName("/tmp/My Repo!")
	if err != nil {
		t.Fatalf("newSidecarName() error = %v", err)
	}
	if first == second {
		t.Fatalf("newSidecarName() returned duplicate names: %q", first)
	}

	validDockerName := regexp.MustCompile(`^stunt-keploy-my-repo-[0-9a-f]{8}$`)
	for _, name := range []string{first, second} {
		if !validDockerName.MatchString(name) {
			t.Fatalf("sidecar name %q does not match %s", name, validDockerName.String())
		}
	}
}
