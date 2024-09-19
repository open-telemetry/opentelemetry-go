// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

type multiLoggerProvider struct {
	embedded.LoggerProvider

	providers []log.LoggerProvider

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// MultiLoggerProvider returns a composite (fan-out) provider.
// It duplicates its calls to all the provided providers.
// It can be used to set up multiple processing pipelines.
// For instance, you can have separate providers for OTel events
// and application logs.
func MultiLoggerProvider(providers ...log.LoggerProvider) log.LoggerProvider {
	return &multiLoggerProvider{
		providers: providers,
	}
}

// Logger returns a logger delegating to loggers created by all providers.
func (p *multiLoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	var loggers []log.Logger
	for _, p := range p.providers {
		loggers = append(loggers, p.Logger(name, opts...))
	}
	return &multiLogger{loggers: loggers}
}

type multiLogger struct {
	embedded.Logger

	loggers []log.Logger

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

func (l *multiLogger) Emit(ctx context.Context, r log.Record) {
	for _, l := range l.loggers {
		l.Emit(ctx, r)
	}
}

func (l *multiLogger) Enabled(ctx context.Context, param log.EnabledParameters) bool {
	for _, l := range l.loggers {
		if !l.Enabled(ctx, param) {
			return false
		}
	}
	return true
}
