package logger

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel/otellog"
)

// This internal package hides the actual logging functions from the user.

// Logger instance used by xray to log. Set via xray.SetLogger().
var Logger otellog.Logger = otellog.NewDefaultLogger(os.Stdout, otellog.LogLevelInfo)

func Debugf(format string, args ...interface{}) {
	Logger.Log(otellog.LogLevelDebug, printfArgs{format, args})
}

func Debug(args ...interface{}) {
	Logger.Log(otellog.LogLevelDebug, printArgs(args))
}

func DebugDeferred(fn func() string) {
	Logger.Log(otellog.LogLevelDebug, stringerFunc(fn))
}

func Infof(format string, args ...interface{}) {
	Logger.Log(otellog.LogLevelInfo, printfArgs{format, args})
}

func Info(args ...interface{}) {
	Logger.Log(otellog.LogLevelInfo, printArgs(args))
}

func Warnf(format string, args ...interface{}) {
	Logger.Log(otellog.LogLevelWarn, printfArgs{format, args})
}

func Warn(args ...interface{}) {
	Logger.Log(otellog.LogLevelWarn, printArgs(args))
}

func Errorf(format string, args ...interface{}) {
	Logger.Log(otellog.LogLevelError, printfArgs{format, args})
}

func Error(args ...interface{}) {
	Logger.Log(otellog.LogLevelError, printArgs(args))
}

type printfArgs struct {
	format string
	args   []interface{}
}

func (p printfArgs) String() string {
	return fmt.Sprintf(p.format, p.args...)
}

type printArgs []interface{}

func (p printArgs) String() string {
	return fmt.Sprint([]interface{}(p)...)
}

type stringerFunc func() string

func (sf stringerFunc) String() string {
	return sf()
}

