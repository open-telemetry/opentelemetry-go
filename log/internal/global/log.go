// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/log/internal/global"

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

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
	return &logger{name: name, options: options}
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	if del, ok := l.delegate.Load().(log.Logger); ok {
		del.Emit(ctx, r)
	}
}

func (l *logger) Enabled(ctx context.Context, r log.Record) bool {
	var enabled bool
	if del, ok := l.delegate.Load().(log.Logger); ok {
		enabled = del.Enabled(ctx, r)
	}
	return enabled
}

func (l *logger) setDelegate(provider log.LoggerProvider) {
	l.delegate.Store(provider.Logger(l.name, l.options...))
}
