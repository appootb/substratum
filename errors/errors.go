package errors

import (
	"errors"
	"fmt"

	"github.com/appootb/protobuf/go/code"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type StatusError spb.Status

func CodeType(err error) code.Error {
	var stErr *StatusError
	if errors.As(err, &stErr) {
		return code.Error(stErr.Code)
	}
	return code.Error(codes.Unknown)
}

func (err *StatusError) Error() string {
	p := (*spb.Status)(err)
	return fmt.Sprintf("status error: code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func New(code code.Error, msg string) error {
	return &StatusError{
		Code:    int32(code),
		Message: msg,
	}
}

func Newf(code code.Error, format string, a ...interface{}) error {
	return &StatusError{
		Code:    int32(code),
		Message: fmt.Sprintf(format, a...),
	}
}
