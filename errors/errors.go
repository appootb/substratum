// Package errors provides gRPC-friendly status errors used across substratum services.
package errors

import (
	"errors"
	"fmt"
	"strconv"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusError is a protobuf Status embedded as a Go error.
// It is recognized by ErrorCode and the response interceptor for client-facing codes.
type StatusError spb.Status

// ErrorCode extracts a numeric gRPC/business code from err.
// Prefer StatusError, then google.golang.org/grpc/status errors; otherwise Unknown.
func ErrorCode(err error) int32 {
	var stErr *StatusError
	if errors.As(err, &stErr) {
		return stErr.Code
	}
	if s, ok := status.FromError(err); ok {
		return int32(s.Code())
	}
	return int32(codes.Unknown)
}

// Error implements the error interface.
func (err *StatusError) Error() string {
	p := (*spb.Status)(err)
	return fmt.Sprintf("status error: code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func parseCode(c interface{}) int32 {
	switch v := c.(type) {
	case int32:
		return v
	case codes.Code:
		return int32(v)
	default:
		i, err := strconv.Atoi(fmt.Sprintf("%d", v))
		if err != nil {
			panic("unknown error code, err: " + err.Error())
		}
		return int32(i)
	}
}

// With wraps err with a static message while preserving errors.Is/As chains.
func With(err error, msg string) error {
	return fmt.Errorf(msg+": %w", err)
}

// Withf wraps err with a formatted message while preserving errors.Is/As chains.
func Withf(err error, format string, a ...interface{}) error {
	args := make([]interface{}, 0, len(a)+1)
	args = append(args, a...)
	args = append(args, err)
	return fmt.Errorf(format+": %w", args...)
}

// New builds a StatusError. code may be int32, codes.Code, or any value fmt-printable as an integer.
func New(code interface{}, msg string) error {
	return &StatusError{
		Code:    parseCode(code),
		Message: msg,
	}
}

// Newf is like New with fmt.Sprintf for the message.
func Newf(code interface{}, format string, a ...interface{}) error {
	return &StatusError{
		Code:    parseCode(code),
		Message: fmt.Sprintf(format, a...),
	}
}
