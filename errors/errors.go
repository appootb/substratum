package errors

import (
	"errors"
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type StatusError spb.Status

func CodeType(err error) int32 {
	var stErr *StatusError
	if errors.As(err, &stErr) {
		return stErr.Code
	}
	return int32(codes.Unknown)
}

func (err *StatusError) Error() string {
	p := (*spb.Status)(err)
	return fmt.Sprintf("status error: code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func New(code int32, msg string) error {
	return &StatusError{
		Code:    code,
		Message: msg,
	}
}

func Newf(code int32, format string, a ...interface{}) error {
	return &StatusError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
