// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/trace"

import (
	"slices"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// TracerConfig is a group of options for a Tracer.
type TracerConfig struct {
	instrumentationVersion string
	// Schema URL of the telemetry emitted by the Tracer.
	schemaURL     string
	attrs         attribute.Set
	profilingMode ProfilingMode
}

// InstrumentationVersion returns the version of the library providing instrumentation.
func (t *TracerConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}

// InstrumentationAttributes returns the attributes associated with the library
// providing instrumentation.
func (t *TracerConfig) InstrumentationAttributes() attribute.Set {
	return t.attrs
}

// SchemaURL returns the Schema URL of the telemetry emitted by the Tracer.
func (t *TracerConfig) SchemaURL() string {
	return t.schemaURL
}

// ProfilingMode returns the profiling mode for the tracer.
func (t *TracerConfig) ProfilingMode() ProfilingMode {
	return t.profilingMode
}

// NewTracerConfig applies all the options to a returned TracerConfig.
func NewTracerConfig(options ...TracerOption) TracerConfig {
	var config TracerConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}

// TracerOption applies an option to a TracerConfig.
type TracerOption interface {
	apply(TracerConfig) TracerConfig
}

type tracerOptionFunc func(TracerConfig) TracerConfig

func (fn tracerOptionFunc) apply(cfg TracerConfig) TracerConfig {
	return fn(cfg)
}

// SpanConfig is a group of options for a Span.
type SpanConfig struct {
	attributes    []attribute.KeyValue
	timestamp     time.Time
	links         []Link
	newRoot       bool
	spanKind      SpanKind
	stackTrace    bool
	profileRegion ProfilingMode
	profileTask   ProfilingMode
	asyncEnd      bool
}

// Attributes describe the associated qualities of a Span.
func (cfg *SpanConfig) Attributes() []attribute.KeyValue {
	return cfg.attributes
}

// Timestamp is a time in a Span life-cycle.
func (cfg *SpanConfig) Timestamp() time.Time {
	return cfg.timestamp
}

// StackTrace reports whether stack trace capturing is enabled.
func (cfg *SpanConfig) StackTrace() bool {
	return cfg.stackTrace
}

// Links are the associations a Span has with other Spans.
func (cfg *SpanConfig) Links() []Link {
	return cfg.links
}

// NewRoot identifies a Span as the root Span for a new trace. This is
// commonly used when an existing trace crosses trust boundaries and the
// remote parent span context should be ignored for security.
func (cfg *SpanConfig) NewRoot() bool {
	return cfg.newRoot
}

// SpanKind is the role a Span has in a trace.
func (cfg *SpanConfig) SpanKind() SpanKind {
	return cfg.spanKind
}

// ProfileRegion reports whether the span should create a runtime/trace.Region.
// The returned mode may be overridden by tracer-level profiling settings.
func (cfg *SpanConfig) ProfileRegion() ProfilingMode {
	return cfg.profileRegion
}

// ProfileTask reports whether the span should create a runtime/trace.Task.
// The returned mode may be overridden by tracer-level profiling settings.
func (cfg *SpanConfig) ProfileTask() ProfilingMode {
	return cfg.profileTask
}

// AsyncEnd reports whether the span will be ended on a different goroutine
// than the one it was started on.
func (cfg *SpanConfig) AsyncEnd() bool {
	return cfg.asyncEnd
}

// NewSpanStartConfig applies all the options to a returned SpanConfig.
// No validation is performed on the returned SpanConfig (e.g. no uniqueness
// checking or bounding of data), it is left to the SDK to perform this
// action.
func NewSpanStartConfig(options ...SpanStartOption) SpanConfig {
	var c SpanConfig
	for _, option := range options {
		c = option.applySpanStart(c)
	}
	return c
}

// NewSpanEndConfig applies all the options to a returned SpanConfig.
// No validation is performed on the returned SpanConfig (e.g. no uniqueness
// checking or bounding of data), it is left to the SDK to perform this
// action.
func NewSpanEndConfig(options ...SpanEndOption) SpanConfig {
	var c SpanConfig
	for _, option := range options {
		c = option.applySpanEnd(c)
	}
	return c
}

// SpanStartOption applies an option to a SpanConfig. These options are applicable
// only when the span is created.
type SpanStartOption interface {
	applySpanStart(SpanConfig) SpanConfig
}

type spanOptionFunc func(SpanConfig) SpanConfig

func (fn spanOptionFunc) applySpanStart(cfg SpanConfig) SpanConfig {
	return fn(cfg)
}

// ComposeSpanStartOptions combines the given options into one, applying them
// sequentially with later options taking precedence.
func ComposeSpanStartOptions(options ...SpanStartOption) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		for _, option := range options {
			cfg = option.applySpanStart(cfg)
		}
		return cfg
	})
}

// SpanEndOption applies an option to a SpanConfig. These options are
// applicable only when the span is ended.
type SpanEndOption interface {
	applySpanEnd(SpanConfig) SpanConfig
}

// EventConfig is a group of options for an Event.
type EventConfig struct {
	attributes []attribute.KeyValue
	timestamp  time.Time
	stackTrace bool
}

// Attributes describe the associated qualities of an Event.
func (cfg *EventConfig) Attributes() []attribute.KeyValue {
	return cfg.attributes
}

// Timestamp is a time in an Event life-cycle.
func (cfg *EventConfig) Timestamp() time.Time {
	return cfg.timestamp
}

// StackTrace reports whether stack trace capturing is enabled.
func (cfg *EventConfig) StackTrace() bool {
	return cfg.stackTrace
}

// NewEventConfig applies all the EventOptions to a returned EventConfig. If no
// timestamp option is passed, the returned EventConfig will have a Timestamp
// set to the call time, otherwise no validation is performed on the returned
// EventConfig.
func NewEventConfig(options ...EventOption) EventConfig {
	var c EventConfig
	for _, option := range options {
		c = option.applyEvent(c)
	}
	if c.timestamp.IsZero() {
		c.timestamp = time.Now()
	}
	return c
}

// EventOption applies span event options to an EventConfig.
type EventOption interface {
	applyEvent(EventConfig) EventConfig
}

// SpanOption are options that can be used at both the beginning and end of a span.
type SpanOption interface {
	SpanStartOption
	SpanEndOption
}

// SpanStartEventOption are options that can be used at the start of a span, or with an event.
type SpanStartEventOption interface {
	SpanStartOption
	EventOption
}

// SpanEndEventOption are options that can be used at the end of a span, or with an event.
type SpanEndEventOption interface {
	SpanEndOption
	EventOption
}

type attributeOption []attribute.KeyValue

func (o attributeOption) applySpan(c SpanConfig) SpanConfig {
	c.attributes = append(c.attributes, []attribute.KeyValue(o)...)
	return c
}
func (o attributeOption) applySpanStart(c SpanConfig) SpanConfig { return o.applySpan(c) }
func (o attributeOption) applyEvent(c EventConfig) EventConfig {
	c.attributes = append(c.attributes, []attribute.KeyValue(o)...)
	return c
}

var _ SpanStartEventOption = attributeOption{}

// WithAttributes adds the attributes related to a span life-cycle event.
// These attributes are used to describe the work a Span represents when this
// option is provided to a Span's start event. Otherwise, these
// attributes provide additional information about the event being recorded
// (e.g. error, state change, processing progress, system event).
//
// If multiple of these options are passed the attributes of each successive
// option will extend the attributes instead of overwriting. There is no
// guarantee of uniqueness in the resulting attributes.
func WithAttributes(attributes ...attribute.KeyValue) SpanStartEventOption {
	return attributeOption(attributes)
}

// SpanEventOption are options that can be used with an event or a span.
type SpanEventOption interface {
	SpanOption
	EventOption
}

type timestampOption time.Time

func (o timestampOption) applySpan(c SpanConfig) SpanConfig {
	c.timestamp = time.Time(o)
	return c
}
func (o timestampOption) applySpanStart(c SpanConfig) SpanConfig { return o.applySpan(c) }
func (o timestampOption) applySpanEnd(c SpanConfig) SpanConfig   { return o.applySpan(c) }
func (o timestampOption) applyEvent(c EventConfig) EventConfig {
	c.timestamp = time.Time(o)
	return c
}

var _ SpanEventOption = timestampOption{}

// WithTimestamp sets the time of a Span or Event life-cycle moment (e.g.
// started, stopped, errored).
func WithTimestamp(t time.Time) SpanEventOption {
	return timestampOption(t)
}

type stackTraceOption bool

func (o stackTraceOption) applyEvent(c EventConfig) EventConfig {
	c.stackTrace = bool(o)
	return c
}

func (o stackTraceOption) applySpan(c SpanConfig) SpanConfig {
	c.stackTrace = bool(o)
	return c
}
func (o stackTraceOption) applySpanEnd(c SpanConfig) SpanConfig { return o.applySpan(c) }

// WithStackTrace sets the flag to capture the error with stack trace (e.g. true, false).
func WithStackTrace(b bool) SpanEndEventOption {
	return stackTraceOption(b)
}

// WithProfileRegion controls whether the span should create a
// runtime/trace.Region.
//   - ProfilingManual: the span always creates a Region.
//   - ProfilingDisabled: the span does not create a Region.
//   - ProfilingDefault: equivalent to ProfilingDisabled.
//   - ProfilingAuto: the tracer decides whether to create a Region. This is
//     typically configured at the tracer level via WithAutoProfiling, but it
//     may also be set per span.
//
// Note: when profiling is set to ProfilingAuto, spans that end on a different
// goroutine than they started must be annotated with AsyncEnd().
func WithProfileRegion(profileRegion ProfilingMode) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.profileRegion = profileRegion
		return cfg
	})
}

// ProfileRegion is equivalent to WithProfileRegion(ProfilingManual).
func ProfileRegion() SpanStartOption {
	return WithProfileRegion(ProfilingManual)
}

// WithProfileTask controls whether the span should create a runtime/trace.Task.
//   - ProfilingManual: the span always creates a Task.
//   - ProfilingDisabled: the span does not create a Task.
//   - ProfilingDefault: only root local spans create a Task.
//   - ProfilingAuto: the tracer decides whether to create a Task. This is
//     typically configured at the tracer level via WithAutoProfiling, but it
//     may also be set per span.
func WithProfileTask(profileTask ProfilingMode) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.profileTask = profileTask
		return cfg
	})
}

// ProfileTask is equivalent to WithProfileTask(ProfilingManual).
func ProfileTask() SpanStartOption {
	return WithProfileTask(ProfilingManual)
}

// NoProfiling is equivalent to applying both
// WithProfileRegion(ProfilingDisabled) and WithProfileTask(ProfilingDisabled).
func NoProfiling() SpanStartOption {
	return ComposeSpanStartOptions(WithProfileRegion(ProfilingDisabled), WithProfileTask(ProfilingDisabled))
}

// WithAsyncEnd hints the tracer that the span will be ended on a different
// goroutine than the one it was started on.
func WithAsyncEnd(asyncEnd bool) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.asyncEnd = asyncEnd
		return cfg
	})
}

// AsyncEnd is equivalent to WithAsyncEnd(true).
func AsyncEnd() SpanStartOption {
	return WithAsyncEnd(true)
}

// WithLinks adds links to a Span. The links are added to the existing Span
// links, i.e. this does not overwrite. Links with invalid span context are ignored.
func WithLinks(links ...Link) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.links = append(cfg.links, links...)
		return cfg
	})
}

// WithNewRoot specifies that the Span should be treated as a root Span. Any
// existing parent span context will be ignored when defining the Span's trace
// identifiers.
func WithNewRoot() SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.newRoot = true
		return cfg
	})
}

// WithSpanKind sets the SpanKind of a Span.
func WithSpanKind(kind SpanKind) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.spanKind = kind
		return cfg
	})
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) TracerOption {
	return tracerOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.instrumentationVersion = version
		return cfg
	})
}

// mergeSets returns the union of keys between a and b. Any duplicate keys will
// use the value associated with b.
func mergeSets(a, b attribute.Set) attribute.Set {
	// NewMergeIterator uses the first value for any duplicates.
	iter := attribute.NewMergeIterator(&b, &a)
	merged := make([]attribute.KeyValue, 0, a.Len()+b.Len())
	for iter.Next() {
		merged = append(merged, iter.Attribute())
	}
	return attribute.NewSet(merged...)
}

// WithInstrumentationAttributes adds the instrumentation attributes.
//
// This is equivalent to calling [WithInstrumentationAttributeSet] with an
// [attribute.Set] created from a clone of the passed attributes.
// [WithInstrumentationAttributeSet] is recommended for more control.
//
// If multiple [WithInstrumentationAttributes] or [WithInstrumentationAttributeSet]
// options are passed, the attributes will be merged together in the order
// they are passed. Attributes with duplicate keys will use the last value passed.
func WithInstrumentationAttributes(attr ...attribute.KeyValue) TracerOption {
	set := attribute.NewSet(slices.Clone(attr)...)
	return WithInstrumentationAttributeSet(set)
}

// WithInstrumentationAttributeSet adds the instrumentation attributes.
//
// If multiple [WithInstrumentationAttributes] or [WithInstrumentationAttributeSet]
// options are passed, the attributes will be merged together in the order
// they are passed. Attributes with duplicate keys will use the last value passed.
func WithInstrumentationAttributeSet(set attribute.Set) TracerOption {
	if set.Len() == 0 {
		return tracerOptionFunc(func(config TracerConfig) TracerConfig {
			return config
		})
	}

	return tracerOptionFunc(func(config TracerConfig) TracerConfig {
		if config.attrs.Len() == 0 {
			config.attrs = set
		} else {
			config.attrs = mergeSets(config.attrs, set)
		}
		return config
	})
}

// WithSchemaURL sets the schema URL for the Tracer.
func WithSchemaURL(schemaURL string) TracerOption {
	return tracerOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.schemaURL = schemaURL
		return cfg
	})
}

// WithProfilingMode controls the profiling behavior of the tracer.
//
// Tracer-level profiling settings can be overridden by span-level settings,
// following the hierarchy:
//
//	ProfilingDefault < ProfilingAuto < ProfilingManual < ProfilingDisabled
//
// A span may override the tracer’s mode only by selecting a *less strict*
// setting. For example, if the tracer is configured with ProfilingAuto,
// a span configured with ProfilingManual can take over and perform its own
// instrumentation. However, if the tracer is configured with ProfilingDisabled,
// no span-level option can re-enable profiling.
//
// When profiling is set to ProfilingAuto, spans that end on a different
// goroutine than they started must be annotated with AsyncEnd().
func WithProfilingMode(mode ProfilingMode) TracerOption {
	return tracerOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.profilingMode = mode
		return cfg
	})
}

// AutoProfiling is equivalent to WithProfilingMode(ProfilingAuto).
func AutoProfiling() TracerOption {
	return WithProfilingMode(ProfilingAuto)
}

// ProfilingMode defines the profiling behavior that can be applied at the
// tracer or span level.
type ProfilingMode int

const (
	// ProfilingDefault is the default profiling mode. Root local spans create
	// a runtime/trace.Task, while child spans do not create a Task or Region
	// unless overridden at the span level.
	ProfilingDefault ProfilingMode = iota

	// ProfilingAuto delegates instrumentation decisions to the tracer. Because
	// tasks must begin and end on the same goroutine, spans that end on a
	// different goroutine than they started must be annotated with AsyncEnd()
	// when using this mode.
	ProfilingAuto

	// ProfilingManual disables default and automatic profiling. The user is
	// responsible for explicitly creating runtime/trace.Task and
	// runtime/trace.Region instrumentation at the span level.
	ProfilingManual

	// ProfilingDisabled disables profiling entirely. All profiler-related settings
	// at the span level are ignored, and no runtime/trace instrumentation is
	// produced.
	ProfilingDisabled
)
