package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/substratum/util/jsonpb"
)

var Default = newConsole()

type Level int

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}
}

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger interface.
type Logger interface {
	UpdateLevel(level Level)
	Log(level Level, md *common.Metadata, msg string, c Content)
}

func newConsole() Logger {
	return &Console{}
}

type Console struct {
	level int32
}

func (log *Console) UpdateLevel(level Level) {
	atomic.StoreInt32(&log.level, int32(level))
}

func (log *Console) Log(level Level, md *common.Metadata, msg string, c Content) {
	if int32(level) < atomic.LoadInt32(&log.level) {
		return
	}
	var (
		meta    []byte
		content []byte
	)
	if md != nil {
		meta, _ = jsonpb.Marshal(md)
	}
	if c != nil && len(c) > 0 {
		content, _ = json.Marshal(c)
	}
	fmt.Println(fmt.Sprintf("%v metadata: %v, %v: %v", level.String(), string(meta), msg, string(content)))
}

func Debug(msg string, c Content) {
	Default.Log(DebugLevel, nil, msg, c)
}

func Info(msg string, c Content) {
	Default.Log(InfoLevel, nil, msg, c)
}

func Warn(msg string, c Content) {
	Default.Log(WarnLevel, nil, msg, c)
}

func Error(msg string, c Content) {
	Default.Log(ErrorLevel, nil, msg, c)
}

func Fatal(msg string, c Content) {
	Default.Log(FatalLevel, nil, msg, c)
	os.Exit(1)
}
