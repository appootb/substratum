package discovery

import (
	"github.com/appootb/grc"
	"github.com/appootb/substratum/discovery"
)

func Init() {
	if discovery.Implementor() != nil && discovery.ConfigImplementor() != nil {
		return
	}
	remoteConfig, _ := grc.New(
		grc.WithDebugProvider(),
		grc.WithConfigAutoCreation(),
		grc.WithBasePath("/debug"))
	if discovery.Implementor() == nil {
		discovery.RegisterImplementor(remoteConfig)
	}
	if discovery.ConfigImplementor() == nil {
		discovery.RegisterConfigImplementor(remoteConfig)
	}
}
