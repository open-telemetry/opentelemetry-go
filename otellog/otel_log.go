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
	// LogLevelDebug is usually only enabled when debugging.
	LogLevelDebug LogLevel = iota + 1

	// LogLevelInfo is general operational entries about what's going on inside the application.
	LogLevelInfo

	// LogLevelWarn is non-critical entries that deserve eyes.
	LogLevelWarn

	// LogLevelError is used for errors that should definitely be noted.
	LogLevelError
)

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
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
	_, _ = fmt.Fprintf(l.w, "%s [%s] %s\n", time.Now().Format(time.RFC3339), ll, msg)
}

var NullLogger = nullLogger{}

type nullLogger struct{}

func (nl nullLogger) Log(ll LogLevel, msg fmt.Stringer) {
}
