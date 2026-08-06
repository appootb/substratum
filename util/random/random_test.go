package random

import (
	"regexp"
	"testing"
)

func TestNumberLengthAndCharset(t *testing.T) {
	re := regexp.MustCompile(`^[0-9]+$`)
	for _, n := range []int{0, 1, 8, 32} {
		s := Number(n)
		if len(s) != n {
			t.Fatalf("len=%d want %d", len(s), n)
		}
		if n > 0 && !re.MatchString(s) {
			t.Fatalf("invalid charset: %q", s)
		}
	}
}

func TestStringLengthAndCharset(t *testing.T) {
	re := regexp.MustCompile(`^[0-9a-zA-Z]+$`)
	for _, n := range []int{0, 1, 16, 64} {
		s := String(n)
		if len(s) != n {
			t.Fatalf("len=%d want %d", len(s), n)
		}
		if n > 0 && !re.MatchString(s) {
			t.Fatalf("invalid charset: %q", s)
		}
	}
}

func TestStringNotConstant(t *testing.T) {
	// Extremely unlikely all equal for n=16 over many draws.
	seen := map[string]struct{}{}
	for i := 0; i < 20; i++ {
		seen[String(16)] = struct{}{}
	}
	if len(seen) < 2 {
		t.Fatal("String appears non-random")
	}
}
