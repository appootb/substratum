package logger

import (
	"os"

	"github.com/appootb/protobuf/go/common"
)

type Helper struct {
	Logger
	md *common.Metadata
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
