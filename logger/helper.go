package logger

import (
	"context"
	"os"

	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/proto/go/common"
)

type Helper struct {
	Logger
	md *common.Metadata
}

func newHelper(ctx context.Context) *Helper {
	md := metadata.RequestMetadata(ctx)
	md.Token = ""
	return &Helper{
		Logger: impl,
		md:     md,
	}
}

func (h *Helper) Debug(msg string, c Content) {
	h.Log(DebugLevel, h.md, msg, c)
}

func (h *Helper) Info(msg string, c Content) {
	h.Log(InfoLevel, h.md, msg, c)
}

func (h *Helper) Warn(msg string, c Content) {
	h.Log(WarnLevel, h.md, msg, c)
}

func (h *Helper) Error(msg string, c Content) {
	h.Log(WarnLevel, h.md, msg, c)
}

func (h *Helper) Fatal(msg string, c Content) {
	h.Log(FatalLevel, h.md, msg, c)
	os.Exit(1)
}
