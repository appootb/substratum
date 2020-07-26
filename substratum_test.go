package substratum

import (
	"os"
	"testing"

	"github.com/appootb/substratum/plugin"
)

func TestMain(m *testing.M) {
	plugin.Register()
	os.Exit(m.Run())
}
