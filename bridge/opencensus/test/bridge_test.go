// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
	ocbridge "go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/bridge/opencensus/internal"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestMixedAPIs(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	tracer := tp.Tracer("mixedapitracer")
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))

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

	spans := sr.Ended()

	if len(spans) != 4 {
		for _, span := range spans {
			t.Logf("Span: %s", span.Name())
		}
		t.Fatalf("Got %d spans, expected %d.", len(spans), 4)
	}

	var parent trace.SpanContext
	for i := len(spans) - 1; i >= 0; i-- {
		// Verify that OpenCensus spans and OpenTelemetry spans have each
		// other as parents.
		if psid := spans[i].Parent().SpanID(); psid != parent.SpanID() {
			t.Errorf("Span %v had parent %v. Expected %v", spans[i].Name(), psid, parent.SpanID())
		}
		parent = spans[i].SpanContext()
	}
}

func TestStartOptions(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))

	ctx := context.Background()
	_, span := octrace.StartSpan(ctx, "OpenCensusSpan", octrace.WithSpanKind(octrace.SpanKindClient))
	span.End()

	spans := sr.Ended()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, expected %d", len(spans), 1)
	}

	if spans[0].SpanKind() != trace.SpanKindClient {
		t.Errorf("Got span kind %v, expected %d", spans[0].SpanKind(), trace.SpanKindClient)
	}
}

func TestStartSpanWithRemoteParent(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))
	tracer := tp.Tracer("remoteparent")

	ctx := context.Background()
	ctx, parent := tracer.Start(ctx, "OpenTelemetrySpan1")

	_, span := octrace.StartSpanWithRemoteParent(ctx, "OpenCensusSpan", ocbridge.OTelSpanContextToOC(parent.SpanContext()))
	span.End()

	spans := sr.Ended()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, expected %d", len(spans), 1)
	}

	if psid := spans[0].Parent().SpanID(); psid != parent.SpanContext().SpanID() {
		t.Errorf("Span %v, had parent %v.  Expected %d", spans[0].Name(), psid, parent.SpanContext().SpanID())
	}
}

func TestToFromContext(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))
	tracer := tp.Tracer("tofromcontext")

	func() {
		ctx := context.Background()

		_, otSpan1 := tracer.Start(ctx, "OpenTelemetrySpan1")
		defer otSpan1.End()

		// Use NewContext instead of the context from Start
		ctx = octrace.NewContext(ctx, internal.NewSpan(otSpan1))

		ctx, _ = tracer.Start(ctx, "OpenTelemetrySpan2")

		// Get the opentelemetry span using the OpenCensus FromContext, and end it
		otSpan2 := octrace.FromContext(ctx)
		defer otSpan2.End()
	}()

	spans := sr.Ended()

	if len(spans) != 2 {
		t.Fatalf("Got %d spans, expected %d.", len(spans), 2)
	}

	var parent trace.SpanContext
	for i := len(spans) - 1; i >= 0; i-- {
		// Verify that OpenCensus spans and OpenTelemetry spans have each
		// other as parents.
		if psid := spans[i].Parent().SpanID(); psid != parent.SpanID() {
			t.Errorf("Span %v had parent %v. Expected %v", spans[i].Name(), psid, parent.SpanID())
		}
		parent = spans[i].SpanContext()
	}
}

func TestIsRecordingEvents(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))

	ctx := context.Background()
	_, ocspan := octrace.StartSpan(ctx, "OpenCensusSpan1")
	if !ocspan.IsRecordingEvents() {
		t.Errorf("Got %v, expected true", ocspan.IsRecordingEvents())
	}
}

func attrsMap(s []attribute.KeyValue) map[attribute.Key]attribute.Value {
	m := make(map[attribute.Key]attribute.Value, len(s))
	for _, a := range s {
		m[a.Key] = a.Value
	}
	return m
}

func TestSetThings(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))

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

	spans := sr.Ended()

	if len(spans) != 1 {
		t.Fatalf("Got %d spans, expected %d.", len(spans), 1)
	}
	s := spans[0]

	if s.Name() != "span-foo" {
		t.Errorf("Got name %v, expected span-foo", s.Name())
	}

	if s.Status().Code != codes.Error {
		t.Errorf("Got code %v, expected %v", s.Status().Code, codes.Error)
	}

	if s.Status().Description != "foo" {
		t.Errorf("Got code %v, expected foo", s.Status().Description)
	}

	attrs := attrsMap(s.Attributes())
	if v := attrs[attribute.Key("bool")]; !v.AsBool() {
		t.Errorf("Got attributes[bool] %v, expected true", v.AsBool())
	}
	if v := attrs[attribute.Key("int64")]; v.AsInt64() != 12345 {
		t.Errorf("Got attributes[int64] %v, expected 12345", v.AsInt64())
	}
	if v := attrs[attribute.Key("float64")]; v.AsFloat64() != 12.345 {
		t.Errorf("Got attributes[float64] %v, expected 12.345", v.AsFloat64())
	}
	if v := attrs[attribute.Key("string")]; v.AsString() != "stringval" {
		t.Errorf("Got attributes[string] %v, expected stringval", v.AsString())
	}

	if len(s.Events()) != 4 {
		t.Fatalf("Got len(events) = %v, expected 4", len(s.Events()))
	}
	annotateEvent := s.Events()[0]
	aeAttrs := attrsMap(annotateEvent.Attributes)
	annotatefEvent := s.Events()[1]
	afeAttrs := attrsMap(annotatefEvent.Attributes)
	sendEvent := s.Events()[2]
	receiveEvent := s.Events()[3]
	if v := aeAttrs[attribute.Key("string")]; v.AsString() != "annotateval" {
		t.Errorf("Got annotateEvent.Attributes[string] = %v, expected annotateval", v.AsString())
	}
	if annotateEvent.Name != "annotate" {
		t.Errorf("Got annotateEvent.Name = %v, expected annotate", annotateEvent.Name)
	}
	if v := afeAttrs[attribute.Key("int64")]; v.AsInt64() != 12345 {
		t.Errorf("Got annotatefEvent.Attributes[int64] = %v, expected 12345", v.AsInt64())
	}
	if v := afeAttrs[attribute.Key("float64")]; v.AsFloat64() != 12.345 {
		t.Errorf("Got annotatefEvent.Attributes[float64] = %v, expected 12.345", v.AsFloat64())
	}
	if annotatefEvent.Name != "annotate67890" {
		t.Errorf("Got annotatefEvent.Name = %v, expected annotate67890", annotatefEvent.Name)
	}
	if v := aeAttrs[attribute.Key("string")]; v.AsString() != "annotateval" {
		t.Errorf("Got annotateEvent.Attributes[string] = %v, expected annotateval", v.AsString())
	}
	seAttrs := attrsMap(sendEvent.Attributes)
	reAttrs := attrsMap(receiveEvent.Attributes)
	if sendEvent.Name != internal.MessageSendEvent {
		t.Errorf("Got sendEvent.Name = %v, expected message send", sendEvent.Name)
	}
	if v := seAttrs[internal.UncompressedKey]; v.AsInt64() != 456 {
		t.Errorf("Got sendEvent.Attributes[uncompressedKey] = %v, expected 456", v.AsInt64())
	}
	if v := seAttrs[internal.CompressedKey]; v.AsInt64() != 789 {
		t.Errorf("Got sendEvent.Attributes[compressedKey] = %v, expected 789", v.AsInt64())
	}
	if receiveEvent.Name != internal.MessageReceiveEvent {
		t.Errorf("Got receiveEvent.Name = %v, expected message receive", receiveEvent.Name)
	}
	if v := reAttrs[internal.UncompressedKey]; v.AsInt64() != 135 {
		t.Errorf("Got receiveEvent.Attributes[uncompressedKey] = %v, expected 135", v.AsInt64())
	}
	if v := reAttrs[internal.CompressedKey]; v.AsInt64() != 369 {
		t.Errorf("Got receiveEvent.Attributes[compressedKey] = %v, expected 369", v.AsInt64())
	}
}
