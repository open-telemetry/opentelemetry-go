// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package. Users can retrieve an in-memory
// logger to verify the behavior of their integrations.
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
func WithEnabledFunc(fn enabledFn) Option {
	return optFunc(func(c config) config {
		c.enabledFn = fn
		return c
	})
}

// NewRecorder returns a new [Recorder].
func NewRecorder(options ...Option) *Recorder {
	cfg := newConfig(options)

	sr := &ScopeRecords{}

	return &Recorder{
		currentScopeRecord: sr,
		enabledFn:          cfg.enabledFn,
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

	// The log records this instrumentation recorded
	Records []log.Record
}

// Recorder is a recorder that stores all received log records
// in-memory.
type Recorder struct {
	embedded.LoggerProvider
	embeddedLogger // nolint:unused  // Used to embed embedded.Logger.

	mu sync.Mutex

	loggers            []*Recorder
	currentScopeRecord *ScopeRecords

	// enabledFn decides whether the recorder should enable logging of a record or not
	enabledFn enabledFn
}

// Logger retrieves a copy of Recorder with the provided scope
// information.
func (r *Recorder) Logger(name string, opts ...log.LoggerOption) log.Logger {
	cfg := log.NewLoggerConfig(opts...)

	nr := &Recorder{
		currentScopeRecord: &ScopeRecords{
			Name:      name,
			Version:   cfg.InstrumentationVersion(),
			SchemaURL: cfg.SchemaURL(),
		},
		enabledFn: r.enabledFn,
	}
	r.addChildLogger(nr)

	return nr
}

func (r *Recorder) addChildLogger(nr *Recorder) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.loggers = append(r.loggers, nr)
}

// Enabled indicates whether a specific record should be stored.
func (r *Recorder) Enabled(ctx context.Context, record log.Record) bool {
	if r.enabledFn == nil {
		return defaultEnabledFunc(ctx, record)
	}

	return r.enabledFn(ctx, record)
}

// Emit stores the log record.
func (r *Recorder) Emit(_ context.Context, record log.Record) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.currentScopeRecord.Records = append(r.currentScopeRecord.Records, record)
}

// Result returns the current in-memory recorder log records.
func (r *Recorder) Result() []*ScopeRecords {
	r.mu.Lock()
	defer r.mu.Unlock()

	ret := []*ScopeRecords{}
	ret = append(ret, r.currentScopeRecord)
	for _, l := range r.loggers {
		ret = append(ret, l.Result()...)
	}
	return ret
}

// Reset clears the in-memory log records.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentScopeRecord != nil {
		r.currentScopeRecord.Records = nil
	}
	for _, l := range r.loggers {
		l.Reset()
	}
}
