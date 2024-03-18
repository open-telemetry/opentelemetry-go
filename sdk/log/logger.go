// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// Compile-time check logger implements log.Logger.
var _ log.Logger = (*logger)(nil)

type logger struct {
	embedded.Logger
}

func newLogger(p *LoggerProvider, scope instrumentation.Scope) *logger {
	// TODO (#5061): Implement.
	return &logger{}
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	// TODO (#5061): Implement.
}

func (l *logger) Enabled(ctx context.Context, r log.Record) bool {
	// TODO (#5061): Implement.
	return true
}
