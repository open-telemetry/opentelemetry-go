// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package. User can retrieve an in-memory
// logger to verify the behavior of their integrations.
package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

// embeddedLogger is a type alias so the embedded.Logger type doesn't conflict
// with the Logger method of the recorder when it is embedded.
type embeddedLogger = embedded.Logger // nolint:unused  // Used below.

type enablerKey uint

var enableKey enablerKey

// NewInMemoryRecorder returns a new InMemoryRecorder.
func NewInMemoryRecorder(options ...Option) *InMemoryRecorder {
	cfg := newConfig(options)
	return &InMemoryRecorder{
		minSeverity: cfg.minSeverity,
	}
}

// Scope represents the instrumentation scope.
type Scope struct {
	// Name is the name of the instrumentation scope.
	Name string
	// Version is the version of the instrumentation scope.
	Version string
	// SchemaURL of the telemetry emitted by the scope.
	SchemaURL string
}

// InMemoryRecorder is a recorder that stores all received log records
// in-memory.
type InMemoryRecorder struct {
	embedded.LoggerProvider
	embeddedLogger // nolint:unused  // Used to embed embedded.Logger.

	mu sync.Mutex

	records []log.Record

	// Scope is the Logger scope recorder received when Logger was called.
	Scope Scope

	// minSeverity is the minimum severity the recorder will return true for
	// when Enabled is called (unless enableKey is set).
	minSeverity log.Severity
}

// Logger retrieves acopy of InMemoryRecorder with the provided scope
// information.
func (i *InMemoryRecorder) Logger(name string, opts ...log.LoggerOption) log.Logger {
	cfg := log.NewLoggerConfig(opts...)

	i.Scope = Scope{
		Name:      name,
		Version:   cfg.InstrumentationVersion(),
		SchemaURL: cfg.SchemaURL(),
	}

	return i
}

// Enabled indicates whether a specific record should be stored, according to
// its severity, or context values.
func (i *InMemoryRecorder) Enabled(ctx context.Context, record log.Record) bool {
	return ctx.Value(enableKey) != nil || record.Severity() >= i.minSeverity
}

// Emit stores the log record.
func (i *InMemoryRecorder) Emit(_ context.Context, record log.Record) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.records = append(i.records, record)
}

// GetRecords returns the current in-memory recorder log records.
func (i *InMemoryRecorder) GetRecords() []log.Record {
	i.mu.Lock()
	defer i.mu.Unlock()
	ret := make([]log.Record, len(i.records))
	copy(ret, i.records)
	return ret
}

// Reset the current in-memory recorder log records.
func (i *InMemoryRecorder) Reset() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.records = []log.Record{}
}

// ContextWithEnabledRecorder forces enabling the recorder, no matter the log
// severity level.
func ContextWithEnabledRecorder(ctx context.Context) context.Context {
	return context.WithValue(ctx, enableKey, true)
}
