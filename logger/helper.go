package logger

import (
	"context"
	"os"

	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/proto/go/common"
	"google.golang.org/protobuf/proto"
)

type Helper struct {
	Logger
	md *common.Metadata
}

func newHelper(ctx context.Context) *Helper {
	return &Helper{
		Logger: impl,
		md:     metadataForLog(ctx),
	}
}

// metadataForLog clones request metadata for logging with Token cleared.
func metadataForLog(ctx context.Context) *common.Metadata {
	src := metadata.IncomingMetadata(ctx)
	if src == nil {
		return nil
	}
	cp := proto.Clone(src).(*common.Metadata)
	cp.Token = ""
	return cp
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
	h.Log(ErrorLevel, h.md, msg, c)
}

func (h *Helper) Fatal(msg string, c Content) {
	h.Log(FatalLevel, h.md, msg, c)
	os.Exit(1)
}
