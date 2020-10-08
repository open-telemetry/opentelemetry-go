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
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

const (
	// FlagsSampled is a bitmask with the sampled bit set. A SpanContext
	// with the sampling bit set means the span is sampled.
	FlagsSampled = byte(0x01)
	// FlagsDeferred is a bitmask with the deferred bit set. A SpanContext
	// with the deferred bit set means the sampling decision has been
	// defered to the receiver.
	FlagsDeferred = byte(0x02)
	// FlagsDebug is a bitmask with the debug bit set.
	FlagsDebug = byte(0x04)

	ErrInvalidHexID errorConst = "trace-id and span-id can only contain [0-9a-f] characters, all lowercase"

	ErrInvalidTraceIDLength errorConst = "hex encoded trace-id must have length equals to 32"
	ErrNilTraceID           errorConst = "trace-id can't be all zero"

	ErrInvalidSpanIDLength errorConst = "hex encoded span-id must have length equals to 16"
	ErrNilSpanID           errorConst = "span-id can't be all zero"
)

type errorConst string

func (e errorConst) Error() string {
	return string(e)
}

// TraceID is a unique identity of a trace.
type TraceID [16]byte

var nilTraceID TraceID
var _ json.Marshaler = nilTraceID

// IsValid checks whether the trace TraceID is valid. A valid trace ID does
// not consist of zeros only.
func (t TraceID) IsValid() bool {
	return !bytes.Equal(t[:], nilTraceID[:])
}

// MarshalJSON implements a custom marshal function to encode TraceID
// as a hex string.
func (t TraceID) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String returns the hex string representation form of a TraceID
func (t TraceID) String() string {
	return hex.EncodeToString(t[:])
}

// SpanID is a unique identity of a span in a trace.
type SpanID [8]byte

var nilSpanID SpanID
var _ json.Marshaler = nilSpanID

// IsValid checks whether the SpanID is valid. A valid SpanID does not consist
// of zeros only.
func (s SpanID) IsValid() bool {
	return !bytes.Equal(s[:], nilSpanID[:])
}

// MarshalJSON implements a custom marshal function to encode SpanID
// as a hex string.
func (s SpanID) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// String returns the hex string representation form of a SpanID
func (s SpanID) String() string {
	return hex.EncodeToString(s[:])
}

// TraceIDFromHex returns a TraceID from a hex string if it is compliant with
// the W3C trace-context specification.  See more at
// https://www.w3.org/TR/trace-context/#trace-id
func TraceIDFromHex(h string) (TraceID, error) {
	t := TraceID{}
	if len(h) != 32 {
		return t, ErrInvalidTraceIDLength
	}

	if err := decodeHex(h, t[:]); err != nil {
		return t, err
	}

	if !t.IsValid() {
		return t, ErrNilTraceID
	}
	return t, nil
}

// SpanIDFromHex returns a SpanID from a hex string if it is compliant
// with the w3c trace-context specification.
// See more at https://www.w3.org/TR/trace-context/#parent-id
func SpanIDFromHex(h string) (SpanID, error) {
	s := SpanID{}
	if len(h) != 16 {
		return s, ErrInvalidSpanIDLength
	}

	if err := decodeHex(h, s[:]); err != nil {
		return s, err
	}

	if !s.IsValid() {
		return s, ErrNilSpanID
	}
	return s, nil
}

func decodeHex(h string, b []byte) error {
	for _, r := range h {
		switch {
		case 'a' <= r && r <= 'f':
			continue
		case '0' <= r && r <= '9':
			continue
		default:
			return ErrInvalidHexID
		}
	}

	decoded, err := hex.DecodeString(h)
	if err != nil {
		return err
	}

	copy(b, decoded)
	return nil
}

// SpanContext contains identifying trace information about a Span.
type SpanContext struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags byte
}

// IsValid returns if the SpanContext is valid. A valid span context has a
// valid TraceID and SpanID.
func (sc SpanContext) IsValid() bool {
	return sc.HasTraceID() && sc.HasSpanID()
}

// HasTraceID checks if the SpanContext has a valid TraceID.
func (sc SpanContext) HasTraceID() bool {
	return sc.TraceID.IsValid()
}

// HasSpanID checks if the SpanContext has a valid SpanID.
func (sc SpanContext) HasSpanID() bool {
	return sc.SpanID.IsValid()
}

// IsDeferred returns if the deferred bit is set in the trace flags.
func (sc SpanContext) IsDeferred() bool {
	return sc.TraceFlags&FlagsDeferred == FlagsDeferred
}

// IsDebug returns if the debug bit is set in the trace flags.
func (sc SpanContext) IsDebug() bool {
	return sc.TraceFlags&FlagsDebug == FlagsDebug
}

// IsSampled returns if the sampling bit is set in the trace flags.
func (sc SpanContext) IsSampled() bool {
	return sc.TraceFlags&FlagsSampled == FlagsSampled
}

type traceContextKeyType int

const (
	currentSpanKey traceContextKeyType = iota
	remoteContextKey
)

// ContextWithSpan returns a copy of parent with span set to current.
func ContextWithSpan(parent context.Context, span Span) context.Context {
	return context.WithValue(parent, currentSpanKey, span)
}

// SpanFromContext returns the current span from ctx, or nil if none set.
func SpanFromContext(ctx context.Context) Span {
	if span, ok := ctx.Value(currentSpanKey).(Span); ok {
		return span
	}
	return noopSpan{}
}

// ContextWithRemoteSpanContext returns a copy of parent with a remote set as
// the remote span context.
func ContextWithRemoteSpanContext(parent context.Context, remote SpanContext) context.Context {
	return context.WithValue(parent, remoteContextKey, remote)
}

// RemoteSpanContextFromContext returns the remote span context from ctx.
func RemoteSpanContextFromContext(ctx context.Context) SpanContext {
	if sc, ok := ctx.Value(remoteContextKey).(SpanContext); ok {
		return sc
	}
	return SpanContext{}
}

// Span is the individual component of a trace. It represents a single named
// and timed operation of a workflow that is traced. A Tracer is used to
// create a Span and it is then up to the operation the Span represents to
// properly end the Span when the operation itself ends.
type Span interface {
	// Tracer returns the Tracer that created the Span. Tracer MUST NOT be
	// nil.
	Tracer() Tracer

	// End completes the Span. Updates are not allowed the Span after End is
	// called other than setting the status.
	End(options ...SpanOption)

	// AddEvent adds an event to the span.
	AddEvent(ctx context.Context, name string, attrs ...label.KeyValue)
	// AddEventWithTimestamp adds an event that occurred at timestamp to the
	// Span.
	AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue)

	// IsRecording returns the recording state of the Span. It will return
	// true if the Span is active and events can be recorded.
	IsRecording() bool

	// RecordError records an error as a Span event.
	RecordError(ctx context.Context, err error, opts ...ErrorOption)

	// SpanContext returns the SpanContext of the Span. The returned
	// SpanContext is usable even after the End has been called for the Span.
	SpanContext() SpanContext

	// SetStatus sets the status of the Span in the form of a code and a
	// message. SetStatus overrides the value of previous calls to SetStatus
	// on the Span.
	SetStatus(code codes.Code, msg string)

	// SetName sets the Span name.
	SetName(name string)

	// SetAttributes sets kv as attributes of the Span. If a key from kv
	// already exists for an attribute of the Span it will be overwritten with
	// the value contained in kv.
	SetAttributes(kv ...label.KeyValue)
}

// Link is the relationship between two Spans. The relationship can be within
// the same Trace or across different Traces.
//
// For example, a Link is used in the following situations:
//
//   1. Batch Processing: A batch of operations may contain operations
//      associated with one or more traces/spans. Since there can only be one
//      parent SpanContext, a Link is used to keep reference to the
//      SpanContext of all operations in the batch.
//   2. Public Endpoint: A SpanContext for an in incoming client request on a
//      public endpoint should be considered untrusted. In such a case, a new
//      trace with its own identity and sampling decision needs to be created,
//      but this new trace needs to be related to the original trace in some
//      form. A Link is used to keep reference to the original SpanContext and
//      track the relationship.
type Link struct {
	SpanContext
	Attributes []label.KeyValue
}

// SpanKind is the role a Span plays in a Trace.
type SpanKind int

const (
	// As a convenience, these match the proto definition, see
	// opentelemetry/proto/trace/v1/trace.proto
	//
	// The unspecified value is not a valid `SpanKind`.  Use
	// `ValidateSpanKind()` to coerce a span kind to a valid
	// value.
	SpanKindUnspecified SpanKind = 0
	SpanKindInternal    SpanKind = 1
	SpanKindServer      SpanKind = 2
	SpanKindClient      SpanKind = 3
	SpanKindProducer    SpanKind = 4
	SpanKindConsumer    SpanKind = 5
)

// ValidateSpanKind returns a valid span kind value.  This will coerce
// invalid values into the default value, SpanKindInternal.
func ValidateSpanKind(spanKind SpanKind) SpanKind {
	switch spanKind {
	case SpanKindInternal,
		SpanKindServer,
		SpanKindClient,
		SpanKindProducer,
		SpanKindConsumer:
		// valid
		return spanKind
	default:
		return SpanKindInternal
	}
}

// String returns the specified name of the SpanKind in lower-case.
func (sk SpanKind) String() string {
	switch sk {
	case SpanKindInternal:
		return "internal"
	case SpanKindServer:
		return "server"
	case SpanKindClient:
		return "client"
	case SpanKindProducer:
		return "producer"
	case SpanKindConsumer:
		return "consumer"
	default:
		return "unspecified"
	}
}

// Tracer is the creator of Spans.
type Tracer interface {
	// Start creates a span.
	Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span)
}

// TracerProvider provides access to instrumentation Tracers.
type TracerProvider interface {
	// Tracer creates an implementation of the Tracer interface.
	// The instrumentationName must be the name of the library providing
	// instrumentation. This name may be the same as the instrumented code
	// only if that code provides built-in instrumentation. If the
	// instrumentationName is empty, then a implementation defined default
	// name will be used instead.
	Tracer(instrumentationName string, opts ...TracerOption) Tracer
}
