// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otel

import (
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// TracerConfig is a group of options for a Tracer.
type TracerConfig struct {
	// InstrumentationVersion is the version of the library providing
	// instrumentation.
	InstrumentationVersion string
}

// NewTracerConfig applies all the options to a returned TracerConfig.
func NewTracerConfig(options ...TracerOption) *TracerConfig {
	config := new(TracerConfig)
	for _, option := range options {
		option.ApplyTracer(config)
	}
	return config
}

// TracerOption applies an option to a TracerConfig.
type TracerOption interface {
	ApplyTracer(*TracerConfig)
}

type instVersionTracerOption string

func (o instVersionTracerOption) ApplyTracer(c *TracerConfig) { c.InstrumentationVersion = string(o) }

// WithInstrumentationVersion sets the instrumentation version for a Tracer.
func WithInstrumentationVersion(version string) TracerOption {
	return instVersionTracerOption(version)
}

// ErrorConfig is a group of options for an error event.
type ErrorConfig struct {
	Timestamp  time.Time
	StatusCode codes.Code
}

// NewErrorConfig applies all the options to a returned ErrorConfig.
func NewErrorConfig(options ...ErrorOption) *ErrorConfig {
	c := new(ErrorConfig)
	for _, option := range options {
		option.ApplyError(c)
	}
	return c
}

// ErrorOption applies an option to a ErrorConfig.
type ErrorOption interface {
	ApplyError(*ErrorConfig)
}

type errorTimeOption time.Time

func (o errorTimeOption) ApplyError(c *ErrorConfig) { c.Timestamp = time.Time(o) }

// WithErrorTime sets the time at which the error event should be recorded.
func WithErrorTime(t time.Time) ErrorOption {
	return errorTimeOption(t)
}

type errorStatusOption struct{ c codes.Code }

func (o errorStatusOption) ApplyError(c *ErrorConfig) { c.StatusCode = o.c }

// WithErrorStatus indicates the span status that should be set when recording an error event.
func WithErrorStatus(c codes.Code) ErrorOption {
	return errorStatusOption{c}
}

// SpanConfig is a group of options for a Span.
type SpanConfig struct {
	// Attributes describe the associated qualities of a Span.
	Attributes []label.KeyValue
	// Timestamp is a time in a Span life-cycle.
	Timestamp time.Time
	// Links are the associations a Span has with other Spans.
	Links []Link
	// Record is the recording state of a Span.
	Record bool
	// NewRoot identifies a Span as the root Span for a new trace. This is
	// commonly used when an existing trace crosses trust boundaries and the
	// remote parent span context should be ignored for security.
	NewRoot bool
	// SpanKind is the role a Span has in a trace.
	SpanKind SpanKind
}

// NewSpanConfig applies all the options to a returned SpanConfig.
// No validation is performed on the returned SpanConfig (e.g. no uniqueness
// checking or bounding of data), it is left to the SDK to perform this
// action.
func NewSpanConfig(options ...SpanOption) *SpanConfig {
	c := new(SpanConfig)
	for _, option := range options {
		option.ApplySpan(c)
	}
	return c
}

// SpanOption applies an option to a SpanConfig.
type SpanOption interface {
	ApplySpan(*SpanConfig)
}

type attributeSpanOption []label.KeyValue

func (o attributeSpanOption) ApplySpan(c *SpanConfig) {
	c.Attributes = append(c.Attributes, []label.KeyValue(o)...)
}

// WithAttributes adds the attributes to a span. These attributes are meant to
// provide additional information about the work the Span represents. The
// attributes are added to the existing Span attributes, i.e. this does not
// overwrite.
func WithAttributes(attributes ...label.KeyValue) SpanOption {
	return attributeSpanOption(attributes)
}

type timestampSpanOption time.Time

func (o timestampSpanOption) ApplySpan(c *SpanConfig) { c.Timestamp = time.Time(o) }

// WithTimestamp sets the time of a Span life-cycle moment (e.g. started or
// stopped).
func WithTimestamp(t time.Time) SpanOption {
	return timestampSpanOption(t)
}

type linksSpanOption []Link

func (o linksSpanOption) ApplySpan(c *SpanConfig) { c.Links = append(c.Links, []Link(o)...) }

// WithLinks adds links to a Span. The links are added to the existing Span
// links, i.e. this does not overwrite.
func WithLinks(links ...Link) SpanOption {
	return linksSpanOption(links)
}

type recordSpanOption bool

func (o recordSpanOption) ApplySpan(c *SpanConfig) { c.Record = bool(o) }

// WithRecord specifies that the span should be recorded. It is important to
// note that implementations may override this option, i.e. if the span is a
// child of an un-sampled trace.
func WithRecord() SpanOption {
	return recordSpanOption(true)
}

type newRootSpanOption bool

func (o newRootSpanOption) ApplySpan(c *SpanConfig) { c.NewRoot = bool(o) }

// WithNewRoot specifies that the Span should be treated as a root Span. Any
// existing parent span context will be ignored when defining the Span's trace
// identifiers.
func WithNewRoot() SpanOption {
	return newRootSpanOption(true)
}

type spanKindSpanOption SpanKind

func (o spanKindSpanOption) ApplySpan(c *SpanConfig) { c.SpanKind = SpanKind(o) }

// WithSpanKind sets the SpanKind of a Span.
func WithSpanKind(kind SpanKind) SpanOption {
	return spanKindSpanOption(kind)
}
