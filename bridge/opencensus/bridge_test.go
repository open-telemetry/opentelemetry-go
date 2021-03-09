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
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

func TestMixedAPIs(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	tracer := tp.Tracer("mixedapitracer")
	octrace.DefaultTracer = NewTracer(tracer)

	func() {
		ctx := context.Background()
		var ocspan1 *octrace.Span
		ctx, ocspan1 = octrace.StartSpan(ctx, "OpenCensusSpan1")
		defer ocspan1.End()

		var otspan1 trace.Span
		ctx, otspan1 = tracer.Start(ctx, "OpenTelemetrySpan1")
		defer otspan1.End()

		var ocspan2 *octrace.Span
		ctx, ocspan2 = octrace.StartSpan(ctx, "OpenCensusSpan2")
		defer ocspan2.End()

		var otspan2 trace.Span
		_, otspan2 = tracer.Start(ctx, "OpenTelemetrySpan2")
		defer otspan2.End()
	}()

	spans := sr.Completed()

	if len(spans) != 4 {
		for _, span := range spans {
			t.Logf("Span: %s", span.Name())
		}
		t.Fatalf("Got %d spans, exepected %d.", len(spans), 4)
	}

	parent := &oteltest.Span{}
	for i := range spans {
		// Reverse the order we look at the spans in, since they are listed in last-to-first order.
		i = len(spans) - i - 1
		// Verify that OpenCensus spans and opentelemetry spans have each other as parents.
		if spans[i].ParentSpanID() != parent.SpanContext().SpanID() {
			t.Errorf("Span %v had parent %v.  Expected %d", spans[i].Name(), spans[i].ParentSpanID(), parent.SpanContext().SpanID())
		}
		parent = spans[i]
	}
}

func TestStartOptions(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	octrace.DefaultTracer = NewTracer(tp.Tracer("startoptionstracer"))

	ctx := context.Background()
	_, span := octrace.StartSpan(ctx, "OpenCensusSpan", octrace.WithSpanKind(octrace.SpanKindClient))
	span.End()

	spans := sr.Completed()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, exepected %d", len(spans), 1)
	}

	if spans[0].SpanKind() != trace.SpanKindClient {
		t.Errorf("Got span kind %v, exepected %d", spans[0].SpanKind(), trace.SpanKindClient)
	}
}

func TestStartSpanWithRemoteParent(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	tracer := tp.Tracer("remoteparent")
	octrace.DefaultTracer = NewTracer(tracer)

	ctx := context.Background()
	ctx, parent := tracer.Start(ctx, "OpenTelemetrySpan1")

	_, span := octrace.StartSpanWithRemoteParent(ctx, "OpenCensusSpan", utils.OTelSpanContextToOC(parent.SpanContext()))
	span.End()

	spans := sr.Completed()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, exepected %d", len(spans), 1)
	}

	if spans[0].ParentSpanID() != parent.SpanContext().SpanID() {
		t.Errorf("Span %v, had parent %v.  Expected %d", spans[0].Name(), spans[0].ParentSpanID(), parent.SpanContext().SpanID())
	}
}

func TestToFromContext(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	tracer := tp.Tracer("tofromcontext")
	octrace.DefaultTracer = NewTracer(tracer)

	func() {
		ctx := context.Background()

		_, otSpan1 := tracer.Start(ctx, "OpenTelemetrySpan1")
		defer otSpan1.End()

		// Use NewContext instead of the context from Start
		ctx = octrace.NewContext(ctx, octrace.NewSpan(&span{otSpan: otSpan1}))

		ctx, _ = tracer.Start(ctx, "OpenTelemetrySpan2")

		// Get the opentelemetry span using the OpenCensus FromContext, and end it
		otSpan2 := octrace.FromContext(ctx)
		defer otSpan2.End()

	}()

	spans := sr.Completed()

	if len(spans) != 2 {
		t.Fatalf("Got %d spans, exepected %d.", len(spans), 2)
	}

	parent := &oteltest.Span{}
	for i := range spans {
		// Reverse the order we look at the spans in, since they are listed in last-to-first order.
		i = len(spans) - i - 1
		// Verify that OpenCensus spans and opentelemetry spans have each other as parents.
		if spans[i].ParentSpanID() != parent.SpanContext().SpanID() {
			t.Errorf("Span %v had parent %v.  Expected %d", spans[i].Name(), spans[i].ParentSpanID(), parent.SpanContext().SpanID())
		}
		parent = spans[i]
	}
}

func TestIsRecordingEvents(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	octrace.DefaultTracer = NewTracer(tp.Tracer("isrecordingevents"))

	ctx := context.Background()
	_, ocspan := octrace.StartSpan(ctx, "OpenCensusSpan1")
	if !ocspan.IsRecordingEvents() {
		t.Errorf("Got %v, expected true", ocspan.IsRecordingEvents())
	}
}

func TestSetThings(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	octrace.DefaultTracer = NewTracer(tp.Tracer("setthings"))

	ctx := context.Background()
	_, ocspan := octrace.StartSpan(ctx, "OpenCensusSpan1")
	ocspan.SetName("span-foo")
	ocspan.SetStatus(octrace.Status{Code: 1, Message: "foo"})
	ocspan.AddAttributes(
		octrace.BoolAttribute("bool", true),
		octrace.Int64Attribute("int64", 12345),
		octrace.Float64Attribute("float64", 12.345),
		octrace.StringAttribute("string", "stringval"),
	)
	ocspan.Annotate(
		[]octrace.Attribute{octrace.StringAttribute("string", "annotateval")},
		"annotate",
	)
	ocspan.Annotatef(
		[]octrace.Attribute{
			octrace.Int64Attribute("int64", 12345),
			octrace.Float64Attribute("float64", 12.345),
		},
		"annotate%d", 67890,
	)
	ocspan.AddMessageSendEvent(123, 456, 789)
	ocspan.AddMessageReceiveEvent(246, 135, 369)
	ocspan.End()

	spans := sr.Completed()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, exepected %d.", len(spans), 1)
	}
	s := spans[0]

	if s.Name() != "span-foo" {
		t.Errorf("Got name %v, expected span-foo", s.Name())
	}

	if s.StatusCode().String() != codes.Error.String() {
		t.Errorf("Got code %v, expected 1", s.StatusCode().String())
	}

	if s.StatusMessage() != "foo" {
		t.Errorf("Got code %v, expected foo", s.StatusMessage())
	}

	if v := s.Attributes()[attribute.Key("bool")]; !v.AsBool() {
		t.Errorf("Got attributes[bool] %v, expected true", v.AsBool())
	}
	if v := s.Attributes()[attribute.Key("int64")]; v.AsInt64() != 12345 {
		t.Errorf("Got attributes[int64] %v, expected 12345", v.AsInt64())
	}
	if v := s.Attributes()[attribute.Key("float64")]; v.AsFloat64() != 12.345 {
		t.Errorf("Got attributes[float64] %v, expected 12.345", v.AsFloat64())
	}
	if v := s.Attributes()[attribute.Key("string")]; v.AsString() != "stringval" {
		t.Errorf("Got attributes[string] %v, expected stringval", v.AsString())
	}

	if len(s.Events()) != 4 {
		t.Fatalf("Got len(events) = %v, expected 4", len(s.Events()))
	}
	annotateEvent := s.Events()[0]
	annotatefEvent := s.Events()[1]
	sendEvent := s.Events()[2]
	receiveEvent := s.Events()[3]
	if v := annotateEvent.Attributes[attribute.Key("string")]; v.AsString() != "annotateval" {
		t.Errorf("Got annotateEvent.Attributes[string] = %v, expected annotateval", v.AsString())
	}
	if annotateEvent.Name != "annotate" {
		t.Errorf("Got annotateEvent.Name = %v, expected annotate", annotateEvent.Name)
	}
	if v := annotatefEvent.Attributes[attribute.Key("int64")]; v.AsInt64() != 12345 {
		t.Errorf("Got annotatefEvent.Attributes[int64] = %v, expected 12345", v.AsInt64())
	}
	if v := annotatefEvent.Attributes[attribute.Key("float64")]; v.AsFloat64() != 12.345 {
		t.Errorf("Got annotatefEvent.Attributes[float64] = %v, expected 12.345", v.AsFloat64())
	}
	if annotatefEvent.Name != "annotate67890" {
		t.Errorf("Got annotatefEvent.Name = %v, expected annotate67890", annotatefEvent.Name)
	}
	if v := annotateEvent.Attributes[attribute.Key("string")]; v.AsString() != "annotateval" {
		t.Errorf("Got annotateEvent.Attributes[string] = %v, expected annotateval", v.AsString())
	}
	if sendEvent.Name != "message send" {
		t.Errorf("Got sendEvent.Name = %v, expected message send", sendEvent.Name)
	}
	if v := sendEvent.Attributes[uncompressedKey]; v.AsInt64() != 456 {
		t.Errorf("Got sendEvent.Attributes[uncompressedKey] = %v, expected 456", v.AsInt64())
	}
	if v := sendEvent.Attributes[compressedKey]; v.AsInt64() != 789 {
		t.Errorf("Got sendEvent.Attributes[compressedKey] = %v, expected 789", v.AsInt64())
	}
	if receiveEvent.Name != "message receive" {
		t.Errorf("Got receiveEvent.Name = %v, expected message receive", receiveEvent.Name)
	}
	if v := receiveEvent.Attributes[uncompressedKey]; v.AsInt64() != 135 {
		t.Errorf("Got receiveEvent.Attributes[uncompressedKey] = %v, expected 135", v.AsInt64())
	}
	if v := receiveEvent.Attributes[compressedKey]; v.AsInt64() != 369 {
		t.Errorf("Got receiveEvent.Attributes[compressedKey] = %v, expected 369", v.AsInt64())
	}
}
