package discovery

import (
	"github.com/appootb/grc"
)

var (
	rc, _ = grc.New(
		grc.WithDebugProvider(),
		grc.WithConfigAutoCreation(),
		grc.WithBasePath("/debug"))

	DefaultConfig  Config  = rc
	DefaultService Service = rc
)
