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

package opencensus

import (
	"context"
	"fmt"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// NewTracer returns an implementation of the OpenCensus Tracer interface which
// uses OpenTelemetry APIs.  Using this implementation of Tracer "upgrades"
// libraries that use OpenCensus to OpenTelemetry to facilitate a migration.
func NewTracer(tracer trace.Tracer) octrace.Tracer {
	return &otelTracer{tracer: tracer}
}

type otelTracer struct {
	tracer trace.Tracer
}

var _ octrace.Tracer = (*otelTracer)(nil)

func (o *otelTracer) StartSpan(ctx context.Context, name string, s ...octrace.StartOption) (context.Context, *octrace.Span) {
	ctx, sp := o.tracer.Start(ctx, name, convertStartOptions(s, name)...)
	return ctx, octrace.NewSpan(&span{otSpan: sp})
}

func convertStartOptions(optFns []octrace.StartOption, name string) []trace.SpanOption {
	var ocOpts octrace.StartOptions
	for _, fn := range optFns {
		fn(&ocOpts)
	}
	otOpts := []trace.SpanOption{}
	switch ocOpts.SpanKind {
	case octrace.SpanKindClient:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindClient))
	case octrace.SpanKindServer:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindServer))
	case octrace.SpanKindUnspecified:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindUnspecified))
	}

	if ocOpts.Sampler != nil {
		global.Handle(fmt.Errorf("ignoring custom sampler for span %q created by OpenCensus because OpenTelemetry does not support creating a span with a custom sampler", name))
	}
	return otOpts
}

func (o *otelTracer) StartSpanWithRemoteParent(ctx context.Context, name string, parent octrace.SpanContext, s ...octrace.StartOption) (context.Context, *octrace.Span) {
	// make sure span context is zero'd out so we use the remote parent
	ctx = trace.ContextWithSpan(ctx, nil)
	ctx = trace.ContextWithRemoteSpanContext(ctx, ocSpanContextToOTel(parent))
	return o.StartSpan(ctx, name, s...)
}

func (o *otelTracer) FromContext(ctx context.Context) *octrace.Span {
	otSpan := trace.SpanFromContext(ctx)
	return octrace.NewSpan(&span{otSpan: otSpan})
}

func (o *otelTracer) NewContext(parent context.Context, s *octrace.Span) context.Context {
	if otSpan, ok := s.Internal().(*span); ok {
		return trace.ContextWithSpan(parent, otSpan.otSpan)
	}
	global.Handle(fmt.Errorf("unable to create context with span %q, since it was created using a different tracer", s.String()))
	return parent
}

type span struct {
	otSpan trace.Span
}

func (s *span) IsRecordingEvents() bool {
	return s.otSpan.IsRecording()
}

func (s *span) End() {
	s.otSpan.End()
}

func (s *span) SpanContext() octrace.SpanContext {
	return otelSpanContextToOc(s.otSpan.SpanContext())
}

func (s *span) SetName(name string) {
	s.otSpan.SetName(name)
}

func (s *span) SetStatus(status octrace.Status) {
	s.otSpan.SetStatus(codes.Code(status.Code), status.Message)
}

func (s *span) AddAttributes(attributes ...octrace.Attribute) {
	s.otSpan.SetAttributes(convertAttributes(attributes)...)
}

func convertAttributes(attributes []octrace.Attribute) []label.KeyValue {
	otAttributes := make([]label.KeyValue, len(attributes))
	for i, a := range attributes {
		otAttributes[i] = label.KeyValue{
			Key:   label.Key(a.Key()),
			Value: convertValue(a.Value()),
		}
	}
	return otAttributes
}

func convertValue(ocval interface{}) label.Value {
	switch v := ocval.(type) {
	case bool:
		return label.BoolValue(v)
	case int64:
		return label.Int64Value(v)
	case float64:
		return label.Float64Value(v)
	case string:
		return label.StringValue(v)
	default:
		return label.StringValue("unknown")
	}
}

func (s *span) Annotate(attributes []octrace.Attribute, str string) {
	s.otSpan.AddEvent(str, trace.WithAttributes(convertAttributes(attributes)...))
}

func (s *span) Annotatef(attributes []octrace.Attribute, format string, a ...interface{}) {
	s.Annotate(attributes, fmt.Sprintf(format, a...))
}

var (
	uncompressedKey = label.Key("uncompressed byte size")
	compressedKey   = label.Key("compressed byte size")
)

func (s *span) AddMessageSendEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otSpan.AddEvent("message send",
		trace.WithAttributes(
			label.KeyValue{
				Key:   uncompressedKey,
				Value: label.Int64Value(uncompressedByteSize),
			},
			label.KeyValue{
				Key:   compressedKey,
				Value: label.Int64Value(compressedByteSize),
			}),
	)
}

func (s *span) AddMessageReceiveEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otSpan.AddEvent("message receive",
		trace.WithAttributes(
			label.KeyValue{
				Key:   uncompressedKey,
				Value: label.Int64Value(uncompressedByteSize),
			},
			label.KeyValue{
				Key:   compressedKey,
				Value: label.Int64Value(compressedByteSize),
			}),
	)
}

func (s *span) AddLink(l octrace.Link) {
	global.Handle(fmt.Errorf("ignoring OpenCensus link %+v for span %q because OpenTelemetry doesn't support setting links after creation", l, s.String()))
}

func (s *span) String() string {
	return fmt.Sprintf("span %s", s.otSpan.SpanContext().SpanID.String())
}

func otelSpanContextToOc(sc trace.SpanContext) octrace.SpanContext {
	if sc.IsDebug() || sc.IsDeferred() {
		global.Handle(fmt.Errorf("ignoring OpenTelemetry Debug or Deferred trace flags for span %q because they are not supported by OpenCensus", sc.SpanID))
	}
	var to octrace.TraceOptions
	if sc.IsSampled() {
		// OpenCensus doesn't expose functions to directly set sampled
		to = 0x1
	}
	return octrace.SpanContext{
		TraceID:      octrace.TraceID(sc.TraceID),
		SpanID:       octrace.SpanID(sc.SpanID),
		TraceOptions: to,
	}
}

func ocSpanContextToOTel(sc octrace.SpanContext) trace.SpanContext {
	var traceFlags byte
	if sc.IsSampled() {
		traceFlags = trace.FlagsSampled
	}
	return trace.SpanContext{
		TraceID:    trace.TraceID(sc.TraceID),
		SpanID:     trace.SpanID(sc.SpanID),
		TraceFlags: traceFlags,
	}
}
