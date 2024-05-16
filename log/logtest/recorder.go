// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

// embeddedLogger is a type alias so the embedded.Logger type doesn't conflict
// with the Logger method of the Recorder when it is embedded.
type embeddedLogger = embedded.Logger // nolint:unused  // Used below.

type enabledFn func(context.Context, log.Record) bool

var defaultEnabledFunc = func(context.Context, log.Record) bool {
	return true
}

type config struct {
	enabledFn enabledFn
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.apply(c)
	}

	return c
}

// Option configures a [Recorder].
type Option interface {
	apply(config) config
}

type optFunc func(config) config

func (f optFunc) apply(c config) config { return f(c) }

// WithEnabledFunc allows configuring whether the [Recorder] is enabled for specific log entries or not.
//
// By default, the Recorder is enabled for every log entry.
func WithEnabledFunc(fn func(context.Context, log.Record) bool) Option {
	return optFunc(func(c config) config {
		c.enabledFn = fn
		return c
	})
}

// NewRecorder returns a new [Recorder].
func NewRecorder(options ...Option) *Recorder {
	cfg := newConfig(options)

	return &Recorder{
		enabledFn: cfg.enabledFn,
	}
}

// ScopeRecords represents the records for a single instrumentation scope.
type ScopeRecords struct {
	// Name is the name of the instrumentation scope.
	Name string
	// Version is the version of the instrumentation scope.
	Version string
	// SchemaURL of the telemetry emitted by the scope.
	SchemaURL string

	// Records are the log records this instrumentation scope recorded.
	Records []log.Record
}

// Recorder is a recorder that stores all received log records
// in-memory.
type Recorder struct {
	embedded.LoggerProvider

	mu      sync.Mutex
	loggers []*logger

	// enabledFn decides whether the recorder should enable logging of a record or not
	enabledFn enabledFn
}

// Logger returns a copy of Recorder as a [log.Logger] with the provided scope
// information.
func (r *Recorder) Logger(name string, opts ...log.LoggerOption) log.Logger {
	cfg := log.NewLoggerConfig(opts...)

	nl := &logger{
		scopeRecord: &ScopeRecords{
			Name:      name,
			Version:   cfg.InstrumentationVersion(),
			SchemaURL: cfg.SchemaURL(),
		},
		enabledFn: r.enabledFn,
	}
	r.addChildLogger(nl)

	return nl
}

func (r *Recorder) addChildLogger(nl *logger) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.loggers = append(r.loggers, nl)
}

// Result returns the current in-memory recorder log records.
func (r *Recorder) Result() []*ScopeRecords {
	r.mu.Lock()
	defer r.mu.Unlock()

	ret := []*ScopeRecords{}
	for _, l := range r.loggers {
		ret = append(ret, l.scopeRecord)
	}
	return ret
}

// Reset clears the in-memory log records for all loggers.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, l := range r.loggers {
		l.Reset()
	}
}

type logger struct {
	embedded.Logger

	mu          sync.Mutex
	scopeRecord *ScopeRecords

	// enabledFn decides whether the recorder should enable logging of a record or not.
	enabledFn enabledFn
}

// Enabled indicates whether a specific record should be stored.
func (l *logger) Enabled(ctx context.Context, record log.Record) bool {
	if l.enabledFn == nil {
		return defaultEnabledFunc(ctx, record)
	}

	return l.enabledFn(ctx, record)
}

// Emit stores the log record.
func (l *logger) Emit(_ context.Context, record log.Record) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.scopeRecord.Records = append(l.scopeRecord.Records, record)
}

// Reset clears the in-memory log records.
func (l *logger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.scopeRecord != nil {
		l.scopeRecord.Records = nil
	}
}
