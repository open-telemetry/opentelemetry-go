// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package global provides access to a global implementation of the OpenTelemetry
Logs API.

Deprecated: Use the equivalent APIs in [go.opentelemetry.io/otel].
*/
package global

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
)

// Logger returns a [log.Logger] configured with the provided name and options
// from the globally configured [log.LoggerProvider].
//
// If this is called before a global LoggerProvider is configured, the returned
// Logger will be a No-Op implementation of a Logger. When a global
// LoggerProvider is registered for the first time, the returned Logger is
// updated in place to report to this new LoggerProvider. There is no need to
// call this function again for an updated instance.
//
// This is a convenience function. It is equivalent to:
//
//	GetLoggerProvider().Logger(name, options...)
//
// Deprecated: Use [otel.Logger] instead.
func Logger(name string, options ...log.LoggerOption) log.Logger {
	return otel.Logger(name, options...)
}

// GetLoggerProvider returns the globally configured [log.LoggerProvider].
//
// If a global LoggerProvider has not been configured with [SetLoggerProvider],
// the returned LoggerProvider will be a No-Op implementation. When
// a global LoggerProvider is registered for the first time, the returned
// LoggerProvider and all Loggers it has created are updated in place. There is
// no need to call this function again for an updated instance.
//
// Deprecated: Use [otel.GetLoggerProvider] instead.
func GetLoggerProvider() log.LoggerProvider {
	return otel.GetLoggerProvider()
}

// SetLoggerProvider configures provider as the global [log.LoggerProvider].
//
// Deprecated: Use [otel.SetLoggerProvider] instead.
func SetLoggerProvider(provider log.LoggerProvider) {
	otel.SetLoggerProvider(provider)
}
