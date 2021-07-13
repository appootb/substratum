package errors

import (
	"errors"
	"fmt"
	"strconv"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatusError spb.Status

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

func (err *StatusError) Error() string {
	p := (*spb.Status)(err)
	return fmt.Sprintf("status error: code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func parseCode(c interface{}) int32 {
	switch v := c.(type) {
	case int32:
		return v
	default:
		i, err := strconv.Atoi(fmt.Sprintf("%v", v))
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
