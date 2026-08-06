package errors

import (
	stderrors "errors"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewAndErrorCode(t *testing.T) {
	err := New(codes.NotFound, "missing")
	if ErrorCode(err) != int32(codes.NotFound) {
		t.Fatalf("code=%d want NotFound", ErrorCode(err))
	}
	if err.Error() == "" {
		t.Fatal("empty error string")
	}
}

func TestNewfAndInt32Code(t *testing.T) {
	err := Newf(int32(7), "user %d", 42)
	if ErrorCode(err) != 7 {
		t.Fatalf("code=%d want 7", ErrorCode(err))
	}
	se := err.(*StatusError)
	if se.Message != "user 42" {
		t.Fatalf("message=%q", se.Message)
	}
	if !strings.Contains(err.Error(), "user 42") {
		t.Fatalf("Error()=%q", err.Error())
	}
}

func TestErrorCodeFromGRPCStatus(t *testing.T) {
	err := status.Error(codes.PermissionDenied, "nope")
	if ErrorCode(err) != int32(codes.PermissionDenied) {
		t.Fatalf("code=%d", ErrorCode(err))
	}
}

func TestErrorCodeUnknown(t *testing.T) {
	err := stderrors.New("plain")
	if ErrorCode(err) != int32(codes.Unknown) {
		t.Fatalf("code=%d want Unknown", ErrorCode(err))
	}
}

func TestWithWraps(t *testing.T) {
	base := stderrors.New("root")
	err := With(base, "wrap")
	if !stderrors.Is(err, base) {
		t.Fatal("expected errors.Is to match base")
	}
	errf := Withf(base, "id=%d", 1)
	if !stderrors.Is(errf, base) {
		t.Fatal("expected errors.Is to match base for Withf")
	}
}
