package discovery

import (
	"github.com/appootb/grc"
)

var (
	rc, _ = grc.New(grc.WithDebugProvider())

	DefaultConfig  Config  = rc
	DefaultService Service = rc
)
