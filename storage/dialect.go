package storage

import (
	"github.com/appootb/substratum/v2/configure"
	"gorm.io/gorm"
)

var (
	sqlImpl SQLDialect
)

// SQLDialectImplementor returns the gorm sql Dialector implementor.
func SQLDialectImplementor() SQLDialect {
	return sqlImpl
}

// RegisterSQLDialectImplementor registers gorm sql Dialector.
func RegisterSQLDialectImplementor(sql SQLDialect) {
	if sqlImpl != nil {
		panic("storage: sql dialect implementor already registered")
	}
	sqlImpl = sql
}

type SQLDialect interface {
	Open(configure.Address) gorm.Dialector
}
