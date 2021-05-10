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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus/utils"
	"go.opentelemetry.io/otel/codes"
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
		otel.Handle(fmt.Errorf("ignoring custom sampler for span %q created by OpenCensus because OpenTelemetry does not support creating a span with a custom sampler", name))
	}
	return otOpts
}

func (o *otelTracer) StartSpanWithRemoteParent(ctx context.Context, name string, parent octrace.SpanContext, s ...octrace.StartOption) (context.Context, *octrace.Span) {
	// make sure span context is zero'd out so we use the remote parent
	ctx = trace.ContextWithSpan(ctx, nil)
	ctx = trace.ContextWithRemoteSpanContext(ctx, utils.OCSpanContextToOTel(parent))
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
	otel.Handle(fmt.Errorf("unable to create context with span %q, since it was created using a different tracer", s.String()))
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
	return utils.OTelSpanContextToOC(s.otSpan.SpanContext())
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

func convertAttributes(attributes []octrace.Attribute) []attribute.KeyValue {
	otAttributes := make([]attribute.KeyValue, len(attributes))
	for i, a := range attributes {
		otAttributes[i] = attribute.KeyValue{
			Key:   attribute.Key(a.Key()),
			Value: convertValue(a.Value()),
		}
	}
	return otAttributes
}

func convertValue(ocval interface{}) attribute.Value {
	switch v := ocval.(type) {
	case bool:
		return attribute.BoolValue(v)
	case int64:
		return attribute.Int64Value(v)
	case float64:
		return attribute.Float64Value(v)
	case string:
		return attribute.StringValue(v)
	default:
		return attribute.StringValue("unknown")
	}
}

func (s *span) Annotate(attributes []octrace.Attribute, str string) {
	s.otSpan.AddEvent(str, trace.WithAttributes(convertAttributes(attributes)...))
}

func (s *span) Annotatef(attributes []octrace.Attribute, format string, a ...interface{}) {
	s.Annotate(attributes, fmt.Sprintf(format, a...))
}

var (
	uncompressedKey = attribute.Key("uncompressed byte size")
	compressedKey   = attribute.Key("compressed byte size")
)

func (s *span) AddMessageSendEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otSpan.AddEvent("message send",
		trace.WithAttributes(
			attribute.KeyValue{
				Key:   uncompressedKey,
				Value: attribute.Int64Value(uncompressedByteSize),
			},
			attribute.KeyValue{
				Key:   compressedKey,
				Value: attribute.Int64Value(compressedByteSize),
			}),
	)
}

func (s *span) AddMessageReceiveEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otSpan.AddEvent("message receive",
		trace.WithAttributes(
			attribute.KeyValue{
				Key:   uncompressedKey,
				Value: attribute.Int64Value(uncompressedByteSize),
			},
			attribute.KeyValue{
				Key:   compressedKey,
				Value: attribute.Int64Value(compressedByteSize),
			}),
	)
}

func (s *span) AddLink(l octrace.Link) {
	otel.Handle(fmt.Errorf("ignoring OpenCensus link %+v for span %q because OpenTelemetry doesn't support setting links after creation", l, s.String()))
}

func (s *span) String() string {
	return fmt.Sprintf("span %s", s.otSpan.SpanContext().SpanID().String())
}
