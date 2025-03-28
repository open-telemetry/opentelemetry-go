// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

type enabledFn func(context.Context, log.EnabledParameters) bool

var defaultEnabledFunc = func(context.Context, log.EnabledParameters) bool {
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
func WithEnabledFunc(fn func(context.Context, log.EnabledParameters) bool) Option {
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

// Recording represents the recorded log records snapshot.
type Recording map[Scope][]Record

// Scope represents the instrumentation scope.
type Scope struct {
	// Name is the name of the instrumentation scope. This should be the
	// Go package name of that scope.
	Name string
	// Version is the version of the instrumentation scope.
	Version string
	// SchemaURL of the telemetry emitted by the scope.
	SchemaURL string
	// Attributes of the telemetry emitted by the scope.
	Attributes attribute.Set
}

// Record represents the record alongside its context.
type Record struct {
	// Ensure forward compatibility by explicitly making this not comparable.
	_ [0]func()

	Context           context.Context
	EventName         string
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          log.Severity
	SeverityText      string
	Body              log.Value
	Attributes        []log.KeyValue
}

// Recorder stores all received log records in-memory.
// Recorder implements [log.LoggerProvider].
type Recorder struct {
	// Ensure forward compatibility by explicitly making this not comparable.
	_ [0]func()

	embedded.LoggerProvider

	mu      sync.Mutex
	loggers map[Scope]*logger

	// enabledFn decides whether the recorder should enable logging of a record or not
	enabledFn enabledFn
}

// Compile-time check Recorder implements log.LoggerProvider.
var _ log.LoggerProvider = (*Recorder)(nil)

// Clone returns a deep copy.
func (a Record) Clone() Record {
	b := a
	attrs := make([]log.KeyValue, len(a.Attributes))
	copy(attrs, a.Attributes)
	b.Attributes = attrs
	return b
}

// Logger returns a copy of Recorder as a [log.Logger] with the provided scope
// information.
func (r *Recorder) Logger(name string, opts ...log.LoggerOption) log.Logger {
	cfg := log.NewLoggerConfig(opts...)
	scope := Scope{
		Name:       name,
		Version:    cfg.InstrumentationVersion(),
		SchemaURL:  cfg.SchemaURL(),
		Attributes: cfg.InstrumentationAttributes(),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.loggers == nil {
		r.loggers = make(map[Scope]*logger)
	}

	l, ok := r.loggers[scope]
	if ok {
		return l
	}
	l = &logger{
		enabledFn: r.enabledFn,
	}
	r.loggers[scope] = l
	return l
}

// Reset clears the in-memory log records for all loggers.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, l := range r.loggers {
		l.Reset()
	}
}

// Result returns a deep copy of the current in-memory recorded log records.
func (r *Recorder) Result() Recording {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make(Recording, len(r.loggers))
	for s, l := range r.loggers {
		func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			if l.records == nil {
				res[s] = nil
				return
			}
			recs := make([]Record, len(l.records))
			for i, r := range l.records {
				recs[i] = r.Clone()
			}
			res[s] = recs
		}()
	}
	return res
}

type logger struct {
	embedded.Logger

	mu      sync.Mutex
	records []*Record

	// enabledFn decides whether the recorder should enable logging of a record or not.
	enabledFn enabledFn
}

// Enabled indicates whether a specific record should be stored.
func (l *logger) Enabled(ctx context.Context, param log.EnabledParameters) bool {
	if l.enabledFn == nil {
		return defaultEnabledFunc(ctx, param)
	}

	return l.enabledFn(ctx, param)
}

// Emit stores the log record.
func (l *logger) Emit(ctx context.Context, record log.Record) {
	l.mu.Lock()
	defer l.mu.Unlock()

	attrs := make([]log.KeyValue, 0, record.AttributesLen())
	record.WalkAttributes(func(kv log.KeyValue) bool {
		attrs = append(attrs, kv)
		return true
	})

	r := &Record{
		Context:           ctx,
		EventName:         record.EventName(),
		Timestamp:         record.Timestamp(),
		ObservedTimestamp: record.ObservedTimestamp(),
		Severity:          record.Severity(),
		SeverityText:      record.SeverityText(),
		Body:              record.Body(),
		Attributes:        attrs,
	}

	l.records = append(l.records, r)
}

// Reset clears the in-memory log records.
func (l *logger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.records = nil
}
