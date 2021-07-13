package logger

import (
	"os"

	"github.com/appootb/substratum/proto/go/common"
)

var (
	impl Logger
)

// Implementor returns the logger service implementor.
func Implementor() Logger {
	return impl
}

// RegisterImplementor registers the logger service implementor.
func RegisterImplementor(log Logger) {
	impl = log
}

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

func Debug(msg string, c Content) {
	impl.Log(DebugLevel, nil, msg, c)
}

func Info(msg string, c Content) {
	impl.Log(InfoLevel, nil, msg, c)
}

func Warn(msg string, c Content) {
	impl.Log(WarnLevel, nil, msg, c)
}

func Error(msg string, c Content) {
	impl.Log(ErrorLevel, nil, msg, c)
}

func Fatal(msg string, c Content) {
	impl.Log(FatalLevel, nil, msg, c)
	os.Exit(1)
}
