package errors

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/appootb/protobuf/go/code"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type StatusError spb.Status

func ErrorCode(err error) int32 {
	var stErr *StatusError
	if errors.As(err, &stErr) {
		return stErr.Code
	}
	return int32(codes.Unknown)
}

func CodeType(err error) code.Error {
	return code.Error(ErrorCode(err))
}

func (err *StatusError) Error() string {
	p := (*spb.Status)(err)
	return fmt.Sprintf("status error: code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func parseCode(c interface{}) int32 {
	switch v := c.(type) {
	case int32:
		return v
	case code.Error:
		return int32(v)
	default:
		i, err := strconv.Atoi(fmt.Sprintf("%d", v))
		if err != nil {
			panic("unknown error code, err: " + err.Error())
		}
		return int32(i)
	}
}

func New(code interface{}, msg string) error {
	return &StatusError{
		Code:    parseCode(code),
		Message: msg,
	}
}

func Newf(code interface{}, format string, a ...interface{}) error {
	return &StatusError{
		Code:    parseCode(code),
		Message: fmt.Sprintf(format, a...),
	}
}
