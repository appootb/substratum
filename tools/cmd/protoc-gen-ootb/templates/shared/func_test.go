package shared

import (
	"strings"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
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

func TestStreamServerImplNaming(t *testing.T) {
	got := streamServerImplName("Example_b", "Test_a")
	want := "exampleBTestAServer"
	if got != want {
		t.Fatalf("streamServerImplName=%q want %q", got, want)
	}
}

func streamServerImplName(service, method string) string {
	svc := pgsgo.PGGUpperCamelCase(pgs.Name(service)).String()
	mtd := pgsgo.PGGUpperCamelCase(pgs.Name(method)).String()
	if svc != "" {
		svc = strings.ToLower(svc[:1]) + svc[1:]
	}
	return svc + mtd + "Server"
}
