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
	"go.opentelemetry.io/otel/unit"
)

// InstrumentConfig contains options for instrument descriptors.
type InstrumentConfig struct {
	// Description describes the instrument in human-readable terms.
	Description string
	// Unit describes the measurement unit for a instrument.
	Unit unit.Unit
	// InstrumentationName is the name of the library providing
	// instrumentation.
	InstrumentationName string
	// InstrumentationVersion is the version of the library providing
	// instrumentation.
	InstrumentationVersion string
}

// InstrumentOption is an interface for applying instrument options.
type InstrumentOption interface {
	// ApplyMeter is used to set a InstrumentOption value of a
	// InstrumentConfig.
	ApplyInstrument(*InstrumentConfig)
}

// NewInstrumentConfig creates a new InstrumentConfig
// and applies all the given options.
func NewInstrumentConfig(opts ...InstrumentOption) InstrumentConfig {
	var config InstrumentConfig
	for _, o := range opts {
		o.ApplyInstrument(&config)
	}
	return config
}

// WithDescription applies provided description.
func WithDescription(desc string) InstrumentOption {
	return descriptionOption(desc)
}

type descriptionOption string

func (d descriptionOption) ApplyInstrument(config *InstrumentConfig) {
	config.Description = string(d)
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) InstrumentOption {
	return unitOption(unit)
}

type unitOption unit.Unit

func (u unitOption) ApplyInstrument(config *InstrumentConfig) {
	config.Unit = unit.Unit(u)
}

// WithInstrumentationName sets the instrumentation name.
func WithInstrumentationName(name string) InstrumentOption {
	return instrumentationNameOption(name)
}

type instrumentationNameOption string

func (i instrumentationNameOption) ApplyInstrument(config *InstrumentConfig) {
	config.InstrumentationName = string(i)
}

// MeterConfig contains options for Meters.
type MeterConfig struct {
	// InstrumentationVersion is the version of the library providing
	// instrumentation.
	InstrumentationVersion string
}

// MeterOption is an interface for applying Meter options.
type MeterOption interface {
	// ApplyMeter is used to set a MeterOption value of a MeterConfig.
	ApplyMeter(*MeterConfig)
}

// NewMeterConfig creates a new MeterConfig and applies
// all the given options.
func NewMeterConfig(opts ...MeterOption) MeterConfig {
	var config MeterConfig
	for _, o := range opts {
		o.ApplyMeter(&config)
	}
	return config
}

// TracerConfig is a group of options for a Tracer.
//
// Most users will use the tracer options instead.
type TracerConfig struct {
	// InstrumentationVersion is the version of the instrumentation library.
	InstrumentationVersion string
}

// NewTracerConfig applies all the options to a returned TracerConfig.
// The default value for all the fields of the returned TracerConfig are the
// default zero value of the type. Also, this does not perform any validation
// on the returned TracerConfig (e.g. no uniqueness checking or bounding of
// data), instead it is left to the implementations of the SDK to perform this
// action.
func NewTracerConfig(opts ...TracerOption) *TracerConfig {
	config := new(TracerConfig)
	for _, option := range opts {
		option.ApplyTracer(config)
	}
	return config
}

// TracerOption applies an options to a TracerConfig.
type TracerOption interface {
	ApplyTracer(*TracerConfig)
}

// Option is an interface for applying Instrument, Meter, or Tracer options.
type Option interface {
	InstrumentOption
	MeterOption
	TracerOption
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) Option {
	return instrumentationVersionOption(version)
}

type instrumentationVersionOption string

func (i instrumentationVersionOption) ApplyMeter(config *MeterConfig) {
	config.InstrumentationVersion = string(i)
}

func (i instrumentationVersionOption) ApplyInstrument(config *InstrumentConfig) {
	config.InstrumentationVersion = string(i)
}

func (i instrumentationVersionOption) ApplyTracer(config *TracerConfig) {
	config.InstrumentationVersion = string(i)
}

// ErrorConfig provides options to set properties of an error
// event at the time it is recorded.
//
// Most users will use the error options instead.
type ErrorConfig struct {
	Timestamp  time.Time
	StatusCode codes.Code
}

// ErrorOption applies changes to ErrorConfig that sets options when an error event is recorded.
type ErrorOption func(*ErrorConfig)

// WithErrorTime sets the time at which the error event should be recorded.
func WithErrorTime(t time.Time) ErrorOption {
	return func(c *ErrorConfig) {
		c.Timestamp = t
	}
}

// WithErrorStatus indicates the span status that should be set when recording an error event.
func WithErrorStatus(s codes.Code) ErrorOption {
	return func(c *ErrorConfig) {
		c.StatusCode = s
	}
}

// SpanConfig is a group of options for a Span.
//
// Most users will use span options instead.
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
// The default value for all the fields of the returned SpanConfig are the
// default zero value of the type. Also, this does not perform any validation
// on the returned SpanConfig (e.g. no uniqueness checking or bounding of
// data). Instead, it is left to the implementations of the SDK to perform this
// action.
func NewSpanConfig(opts ...SpanOption) *SpanConfig {
	c := new(SpanConfig)
	for _, option := range opts {
		option.Apply(c)
	}
	return c
}

// SpanOption applies an option to a SpanConfig.
type SpanOption interface {
	Apply(*SpanConfig)
}

type attributeSpanOption []label.KeyValue

func (o attributeSpanOption) Apply(c *SpanConfig) {
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

func (o timestampSpanOption) Apply(c *SpanConfig) { c.Timestamp = time.Time(o) }

// WithTimestamp sets the time of a Span life-cycle moment (e.g. started or
// stopped).
func WithTimestamp(t time.Time) SpanOption {
	return timestampSpanOption(t)
}

type linksSpanOption []Link

func (o linksSpanOption) Apply(c *SpanConfig) { c.Links = append(c.Links, []Link(o)...) }

// WithLinks adds links to a Span. The links are added to the existing Span
// links, i.e. this does not overwrite.
func WithLinks(links ...Link) SpanOption {
	return linksSpanOption(links)
}

type recordSpanOption bool

func (o recordSpanOption) Apply(c *SpanConfig) { c.Record = bool(o) }

// WithRecord specifies that the span should be recorded. It is important to
// note that implementations may override this option, i.e. if the span is a
// child of an un-sampled trace.
func WithRecord() SpanOption {
	return recordSpanOption(true)
}

type newRootSpanOption bool

func (o newRootSpanOption) Apply(c *SpanConfig) { c.NewRoot = bool(o) }

// WithNewRoot specifies that the Span should be treated as a root Span. Any
// existing parent span context will be ignored when defining the Span's trace
// identifiers.
func WithNewRoot() SpanOption {
	return newRootSpanOption(true)
}

type spanKindSpanOption SpanKind

func (o spanKindSpanOption) Apply(c *SpanConfig) { c.SpanKind = SpanKind(o) }

// WithSpanKind sets the SpanKind of a Span.
func WithSpanKind(kind SpanKind) SpanOption {
	return spanKindSpanOption(kind)
}
