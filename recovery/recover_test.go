package recovery

import (
	"errors"
	"testing"

	"google.golang.org/grpc/status"
)

func TestRecoverStatusProduction(t *testing.T) {
	t.Setenv("DEVELOPMENT", "")

	err := RecoverStatus("db connection leaked")
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status")
	}
	if st.Message() != publicInternalMessage {
		t.Fatalf("expected generic message, got %q", st.Message())
	}
}

func TestRecoverStatusDevelopment(t *testing.T) {
	t.Setenv("DEVELOPMENT", "local")

	err := RecoverStatus(errors.New("db connection leaked"))
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status")
	}
	if st.Message() != "db connection leaked" {
		t.Fatalf("expected detailed message, got %q", st.Message())
	}
}
