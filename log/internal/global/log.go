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
// Do not use sdk/instrumentation (API cannot depend on the SDK).
type instLib struct{ name, version string }

type loggerProvider struct {
	embedded.LoggerProvider

	mu       sync.Mutex
	loggers  map[instLib]*logger
	delegate log.LoggerProvider
}

// Compile-time guarantee loggerProvider implements LoggerProvider.
var _ log.LoggerProvider = (*loggerProvider)(nil)

func (p *loggerProvider) Logger(cfg log.LoggerConfig) log.Logger {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.delegate != nil {
		return p.delegate.Logger(cfg)
	}

	key := instLib{cfg.Name, cfg.Version}

	if p.loggers == nil {
		l := &logger{cfg: cfg}
		p.loggers = map[instLib]*logger{key: l}
		return l
	}

	if l, ok := p.loggers[key]; ok {
		return l
	}

	l := &logger{cfg: cfg}
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
	p.loggers = nil // Only set logger delegates once.
}

type logger struct {
	embedded.Logger

	cfg log.LoggerConfig

	delegate atomic.Value // log.Logger
}

// Compile-time guarantee logger implements Logger.
var _ log.Logger = (*logger)(nil)

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
	l.delegate.Store(provider.Logger(l.cfg))
}
