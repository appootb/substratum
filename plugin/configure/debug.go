package configure

import (
	"github.com/appootb/grc"
	"github.com/appootb/substratum/configure"
)

func Init() {
	if configure.Implementor() == nil {
		configure.RegisterImplementor(NewDebug())
	}
}

type Debug struct {
	rc *grc.RemoteConfig
}

func NewDebug() *Debug {
	debug := &Debug{}
	debug.rc, _ = grc.New(
		grc.WithDebugProvider(),
		grc.WithConfigAutoCreation())
	return debug
}

// Register the configuration pointer.
func (m *Debug) Register(component string, v interface{}) error {
	return m.rc.RegisterConfig(component, v)
}
