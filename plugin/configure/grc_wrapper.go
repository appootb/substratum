package configure

import (
	"github.com/appootb/grc"
	"github.com/appootb/substratum/configure"
)

func Init() {
	if configure.Implementor() == nil {
		configure.RegisterImplementor(NewGRCWrapper(nil))
	}
}

type GRCWrapper struct {
	rc *grc.RemoteConfig
}

func NewGRCWrapper(rc *grc.RemoteConfig) *GRCWrapper {
	if rc == nil {
		rc, _ = grc.New(grc.WithDebugProvider(),
			grc.WithConfigAutoCreation())
	}
	return &GRCWrapper{
		rc: rc,
	}
}

// Register the configuration pointer.
func (m *GRCWrapper) Register(component string, v interface{}) error {
	return m.rc.RegisterConfig(component, v)
}
