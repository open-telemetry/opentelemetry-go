// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/log/global"

/*
Package global provides a global implementation of the OpenTelemetry Logs
Bridge API.

This package is experimental. It will be deprecated and removed when the [log]
package becomes stable. Its functionality will be migrated to
go.opentelemetry.io/otel.
*/

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
)

// Logger returns a [log.Logger] configured with the provided name and options
// from the globally configured [log.LoggerProvider].
//
// If this is called before a global LoggerProvider is configured, the returned
// Logger will be a No-Op implementation of a Logger. When a global
// LoggerProvider is registered for the first time, the returned Logger is
// updated in-place to report to this new LoggerProvider. There is no need to
// call this function again for an updated instance.
//
// This is a convenience function. It is equivalent to:
//
//	GetLoggerProvider().Logger(name, options...)
func Logger(name string, options ...log.LoggerOption) log.Logger {
	return GetLoggerProvider().Logger(name, options...)
}

// GetLoggerProvider returns the globally configured [log.LoggerProvider].
//
// If a global LoggerProvider has not been configured with [SetLoggerProvider], the returned
// Logger will be a No-Op implementation of a LoggerProvider. When a global
// LoggerProvider is registered for the first time, the returned LoggerProvider
// and all of its created Loggers are updated in-place. There is no need to
// call this function again for an updated instance.
func GetLoggerProvider() log.LoggerProvider {
	// TODO: implement.
	return nil
}

// SetLoggerProvider configures provider as the global [log.LoggerProvider].
func SetLoggerProvider(provider log.LoggerProvider) {
	// TODO: implement.
}

// instLib defines the instrumentation library a logger is created for.
//
// Do not use the sdk/instrumentation package. The API cannot depend on the
// SDK.
type instLib struct {
	name    string
	version string
}

// loggerProvider is a placeholder for a configured SDK LoggerProvider.
//
// All LoggerProvider functionality is forwarded to a delegate once configured.
type loggerProvider struct {
	embedded.LoggerProvider

	mu       sync.Mutex
	loggers  map[instLib]*logger
	delegate log.LoggerProvider
}

// Compile-time guarantee that loggerProvider implements the LoggerProvider
// interface.
var _ log.LoggerProvider = (*loggerProvider)(nil)

func (p *loggerProvider) Logger(name string, options ...log.LoggerOption) log.Logger {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.delegate != nil {
		return p.delegate.Logger(name, options...)
	}

	cfg := log.NewLoggerConfig(options...)
	key := instLib{
		name:    name,
		version: cfg.InstrumentationVersion(),
	}

	if p.loggers == nil {
		l := newLogger(name, options)
		p.loggers = map[instLib]*logger{key: l}
		return l
	}

	if l, ok := p.loggers[key]; ok {
		return l
	}

	l := newLogger(name, options)
	p.loggers[key] = l
	return l
}

func (p *loggerProvider) setDelegate(provider log.LoggerProvider) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.delegate = provider

	for _, l := range p.loggers {
		l.setDelegate(provider)
	}

	// Only set logger delegates once.
	p.loggers = nil
}

type logger struct {
	embedded.Logger

	name    string
	options []log.LoggerOption

	delegate atomic.Value // log.Logger
}

// Compile-time guarantee that logger implements the trace.Tracer interface.
var _ log.Logger = (*logger)(nil)

func newLogger(name string, options []log.LoggerOption) *logger {
	var base log.Logger = noop.Logger{}
	l := &logger{name: name, options: options}
	l.delegate.Store(base)
	return l
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	l.delegate.Load().(log.Logger).Emit(ctx, r)
}

func (l *logger) Enabled(ctx context.Context, r log.Record) bool {
	return l.delegate.Load().(log.Logger).Enabled(ctx, r)
}

func (l *logger) setDelegate(provider log.LoggerProvider) {
	l.delegate.Store(provider.Logger(l.name, l.options...))
}
