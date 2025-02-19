// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"cmp"
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

// Recorder stores all received log records in-memory.
// Recorder implements [log.LoggerProvider].
type Recorder struct {
	embedded.LoggerProvider

	mu      sync.Mutex
	loggers map[Scope]*logger

	// enabledFn decides whether the recorder should enable logging of a record or not
	enabledFn enabledFn
}

// Compile-time check Recorder implements log.LoggerProvider.
var _ log.LoggerProvider = (*Recorder)(nil)

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
		loggers:   make(map[Scope]*logger),
		enabledFn: cfg.enabledFn,
	}
}

// Result represents the recordered log records.
type Result map[Scope][]Record

// Equal returns if a is equal to b.
func (a Result) Equal(b Result) bool {
	return maps.EqualFunc(a, b, func(x, y []Record) bool {
		return slices.EqualFunc(x, y, func(a, b Record) bool { return a.Equal(b) })
	})
}

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
	Context           context.Context
	EventName         string
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          log.Severity
	SeverityText      string
	Body              log.Value
	Attributes        []log.KeyValue
}

// Equal returns if a is equal to b.
func (a Record) Equal(b Record) bool {
	if a.Context != b.Context {
		return false
	}
	if a.EventName != b.EventName {
		return false
	}
	if !a.Timestamp.Equal(b.Timestamp) {
		return false
	}
	if !a.ObservedTimestamp.Equal(b.ObservedTimestamp) {
		return false
	}
	if a.Severity != b.Severity {
		return false
	}
	if a.SeverityText != b.SeverityText {
		return false
	}
	if !a.Body.Equal(b.Body) {
		return false
	}
	aAttrs := sortKVs(a.Attributes)
	bAttrs := sortKVs(b.Attributes)
	if !slices.EqualFunc(aAttrs, bAttrs, log.KeyValue.Equal) { //nolint:gosimple // We want to use the same pattern.
		return false
	}
	return true
}

// Clone returns a deep copy.
func (a Record) Clone() Record {
	b := a
	attrs := make([]log.KeyValue, len(a.Attributes))
	copy(attrs, a.Attributes)
	b.Attributes = attrs
	return b
}

func sortKVs(kvs []log.KeyValue) []log.KeyValue {
	s := make([]log.KeyValue, len(kvs))
	copy(s, kvs)
	slices.SortFunc(s, func(a, b log.KeyValue) int {
		return cmp.Compare(a.Key, b.Key)
	})
	return s
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

// Result returns a deep copy of the current in-memory recorded log records.
func (r *Recorder) Result() Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make(Result, len(r.loggers))
	for s, l := range r.loggers {
		l.mu.Lock()
		recs := make([]Record, len(l.records))
		for i, r := range l.records {
			recs[i] = r.Clone()
		}
		res[s] = recs
		l.mu.Unlock()
	}
	return res
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

	mu      sync.Mutex
	records []*Record

	// enabledFn decides whether the recorder should enable logging of a record or not.
	enabledFn enabledFn
}

// Enabled indicates whether a specific record should be stored.
func (l *logger) Enabled(ctx context.Context, opts log.EnabledParameters) bool {
	if l.enabledFn == nil {
		return defaultEnabledFunc(ctx, opts)
	}

	return l.enabledFn(ctx, opts)
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
