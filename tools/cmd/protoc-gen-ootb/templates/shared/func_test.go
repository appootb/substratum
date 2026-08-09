package shared

import (
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

func TestFullMethodPathMatchesGoGRPC(t *testing.T) {
	cases := []struct {
		pkg, service, method string
		want                 string
	}{
		{"example", "example", "test1", "/example.example/test1"},
		{"example", "Example_b", "test2", "/example.Example_b/test2"},
		{"example", "Example_b", "Test_a", "/example.Example_b/Test_a"},
	}
	for _, c := range cases {
		got := FullMethodPath(pgs.Name(c.pkg), pgs.Name(c.service), pgs.Name(c.method))
		if got != c.want {
			t.Fatalf("FullMethodPath(%q,%q,%q)=%q want %q", c.pkg, c.service, c.method, got, c.want)
		}
	}
}
