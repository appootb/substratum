package configure

import (
	"github.com/appootb/grc"
	"github.com/appootb/substratum/configure"
)

func Init() {
	if configure.Implementor() == nil {
		debug, _ := grc.New(grc.WithDebugProvider(),
			grc.WithConfigAutoCreation())
		Register(debug)
	}
}

func Register(rc *grc.RemoteConfig) {
	configure.RegisterImplementor(&GRCWrapper{rc})
}

type GRCWrapper struct {
	rc *grc.RemoteConfig
}

// Register the configuration pointer.
func (m *GRCWrapper) Register(component string, v interface{}) error {
	return m.rc.RegisterConfig(component, v)
}
