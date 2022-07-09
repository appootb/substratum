package storage

import (
	"fmt"

	"github.com/appootb/substratum/v2/configure"
)

type redisCache struct {
	configure.Address
}

func (d *redisCache) URL() string {
	return fmt.Sprintf("%s://:%s@%s:%s/%s", d.Schema, d.Password, d.Host, d.Port, d.NameSpace)
}
