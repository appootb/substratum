package storage

import (
	"github.com/appootb/substratum/v2/configure"
)

var (
	commonImpls = make(map[configure.Schema]CommonDialect)
)

// CommonDialectImplementor returns the common storage Dialector implementor.
func CommonDialectImplementor(schema configure.Schema) CommonDialect {
	return commonImpls[schema]
}

// RegisterCommonDialectImplementor registers common storage Dialector.
func RegisterCommonDialectImplementor(schema configure.Schema, dialect CommonDialect) {
	commonImpls[schema] = dialect
}

type CommonDialect interface {
	Open(configure.Address) (interface{}, error)
}
