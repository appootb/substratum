package context

import (
	"context"
)

var (
	Context context.Context
	Cancel  context.CancelFunc
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())
}
