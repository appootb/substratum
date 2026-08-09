package generator

import (
	"testing"
)

func TestNormalizeOutFile(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "export.dart"},
		{"barrel", "barrel.dart"},
		{"api.dart", "api.dart"},
	}
	for _, c := range cases {
		if got := normalizeOutFile(c.in); got != c.want {
			t.Fatalf("normalizeOutFile(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizePrefix(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "./src/generated"},
		{"./src/generated", "./src/generated"},
		{"./lib/gen/", "./lib/gen"},
		{"lib/src", "lib/src"},
	}
	for _, c := range cases {
		if got := normalizePrefix(c.in); got != c.want {
			t.Fatalf("normalizePrefix(%q)=%q want %q", c.in, got, c.want)
		}
	}
}
