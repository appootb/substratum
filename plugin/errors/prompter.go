package errors

import (
	"github.com/appootb/substratum/v2/errors"
)

func Init() {
	if errors.Implementor() == nil {
		errors.RegisterImplementor(&Debug{})
	}
}

type Debug struct{}

func (m Debug) Translate(_ int32) string {
	return ""
}
