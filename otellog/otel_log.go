package otellog

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Logger interface {
	Log(level LogLevel, msg fmt.Stringer)
}

type LogLevel int

const (
	// LogLevelTrace is used to for fine-grained debugging event and disabled in default configurations.
	LogLevelTrace  LogLevel = iota + 1

	// LogLevelDebug is usually only enabled when debugging.
	LogLevelDebug

	// LogLevelInfo is used only for informal event indicates that event has happened.
	LogLevelInfo

	// LogLevelWarn is non-critical entries that deserve eyes.
	LogLevelWarn

	// LogLevelError is used for errors that should definitely be noted.
	LogLevelError

	// LogLevelError is used for fatal errors such as application or system crash.
	LogLevelFatal
)

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("UNKNOWNLOGLEVEL<%d>", ll)
	}
}

func NewDefaultLogger(w io.Writer, minLogLevel LogLevel) Logger {
	return &defaultLogger{w: w, minLevel: minLogLevel}
}

type defaultLogger struct {
	mu       sync.Mutex
	w        io.Writer
	minLevel LogLevel
}

func (l *defaultLogger) Log(ll LogLevel, msg fmt.Stringer) {
	if ll < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintf(l.w, "%s\t[%s]\t%s\n", time.Now().Format(time.RFC3339), ll, msg)
}

var NullLogger = nullLogger{}

type nullLogger struct{}

func (nl nullLogger) Log(ll LogLevel, msg fmt.Stringer) {
}
