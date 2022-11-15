package storage

import "github.com/appootb/substratum/v2/configure"

type emptyDialect struct{}

func (d *emptyDialect) Open(configure.Address) (interface{}, error) {
	return nil, nil
}
