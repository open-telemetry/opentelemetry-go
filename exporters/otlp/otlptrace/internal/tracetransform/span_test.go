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

package tracetransform

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestSpanKind(t *testing.T) {
	for _, test := range []struct {
		kind     trace.SpanKind
		expected tracepb.Span_SpanKind
	}{
		{
			trace.SpanKindInternal,
			tracepb.Span_SPAN_KIND_INTERNAL,
		},
		{
			trace.SpanKindClient,
			tracepb.Span_SPAN_KIND_CLIENT,
		},
		{
			trace.SpanKindServer,
			tracepb.Span_SPAN_KIND_SERVER,
		},
		{
			trace.SpanKindProducer,
			tracepb.Span_SPAN_KIND_PRODUCER,
		},
		{
			trace.SpanKindConsumer,
			tracepb.Span_SPAN_KIND_CONSUMER,
		},
		{
			trace.SpanKind(-1),
			tracepb.Span_SPAN_KIND_UNSPECIFIED,
		},
	} {
		assert.Equal(t, test.expected, spanKind(test.kind))
	}
}

func TestNilSpanEvent(t *testing.T) {
	assert.Nil(t, spanEvents(nil))
}

func TestEmptySpanEvent(t *testing.T) {
	assert.Nil(t, spanEvents([]tracesdk.Event{}))
}

func TestSpanEvent(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.Int("one", 1), attribute.Int("two", 2)}
	eventTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)
	got := spanEvents([]tracesdk.Event{
		{
			Name:       "test 1",
			Attributes: []attribute.KeyValue{},
			Time:       eventTime,
		},
		{
			Name:                  "test 2",
			Attributes:            attrs,
			Time:                  eventTime,
			DroppedAttributeCount: 2,
		},
	})
	if !assert.Len(t, got, 2) {
		return
	}
	eventTimestamp := uint64(1589932800 * 1e9)
	assert.Equal(t, &tracepb.Span_Event{Name: "test 1", Attributes: nil, TimeUnixNano: eventTimestamp}, got[0])
	// Do not test Attributes directly, just that the return value goes to the correct field.
	assert.Equal(t, &tracepb.Span_Event{Name: "test 2", Attributes: KeyValues(attrs), TimeUnixNano: eventTimestamp, DroppedAttributesCount: 2}, got[1])
}

func TestNilLinks(t *testing.T) {
	assert.Nil(t, links(nil))
}

func TestEmptyLinks(t *testing.T) {
	assert.Nil(t, links([]tracesdk.Link{}))
}

func TestLinks(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.Int("one", 1), attribute.Int("two", 2)}
	l := []tracesdk.Link{
		{
			DroppedAttributeCount: 3,
		},
		{
			SpanContext:           trace.SpanContext{},
			Attributes:            attrs,
			DroppedAttributeCount: 3,
		},
	}
	got := links(l)

	// Make sure we get the same number back first.
	if !assert.Len(t, got, 2) {
		return
	}

	// Empty should be empty.
	expected := &tracepb.Span_Link{
		TraceId:                []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		SpanId:                 []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		DroppedAttributesCount: 3,
	}
	assert.Equal(t, expected, got[0])

	// Do not test Attributes directly, just that the return value goes to the correct field.
	expected.Attributes = KeyValues(attrs)
	assert.Equal(t, expected, got[1])

	// Changes to our links should not change the produced links.
	l[1].SpanContext = l[1].SpanContext.WithTraceID(trace.TraceID{})
	assert.Equal(t, expected, got[1])
	assert.Equal(t, l[1].DroppedAttributeCount, int(got[1].DroppedAttributesCount))
}

func TestStatus(t *testing.T) {
	for _, test := range []struct {
		code       codes.Code
		message    string
		otlpStatus tracepb.Status_StatusCode
	}{
		{
			codes.Ok,
			"test Ok",
			tracepb.Status_STATUS_CODE_OK,
		},
		{
			codes.Unset,
			"test Unset",
			tracepb.Status_STATUS_CODE_UNSET,
		},
		{
			message:    "default code is unset",
			otlpStatus: tracepb.Status_STATUS_CODE_UNSET,
		},
		{
			codes.Error,
			"test Error",
			tracepb.Status_STATUS_CODE_ERROR,
		},
	} {
		expected := &tracepb.Status{Code: test.otlpStatus, Message: test.message}
		assert.Equal(t, expected, status(test.code, test.message))
	}
}

func TestNilSpan(t *testing.T) {
	assert.Nil(t, span(nil))
}

func TestNilSpanData(t *testing.T) {
	assert.Nil(t, Spans(nil))
}

func TestEmptySpanData(t *testing.T) {
	assert.Nil(t, Spans(nil))
}

func TestSpanData(t *testing.T) {
	// Full test of span data

	// March 31, 2020 5:01:26 1234nanos (UTC)
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)
	traceState, _ := trace.ParseTraceState("key1=val1,key2=val2")
	spanData := tracetest.SpanStub{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:     trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			TraceState: traceState,
		}),
		Parent: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:     trace.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
			TraceState: traceState,
			Remote:     true,
		}),
		SpanKind:  trace.SpanKindServer,
		Name:      "span data to span data",
		StartTime: startTime,
		EndTime:   endTime,
		Events: []tracesdk.Event{
			{Time: startTime,
				Attributes: []attribute.KeyValue{
					attribute.Int64("CompressedByteSize", 512),
				},
			},
			{Time: endTime,
				Attributes: []attribute.KeyValue{
					attribute.String("EventType", "Recv"),
				},
			},
		},
		Links: []tracesdk.Link{
			{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    trace.TraceID{0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF},
					SpanID:     trace.SpanID{0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7},
					TraceFlags: 0,
				}),
				Attributes: []attribute.KeyValue{
					attribute.String("LinkType", "Parent"),
				},
				DroppedAttributeCount: 0,
			},
			{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    trace.TraceID{0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF},
					SpanID:     trace.SpanID{0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7},
					TraceFlags: 0,
				}),
				Attributes: []attribute.KeyValue{
					attribute.String("LinkType", "Child"),
				},
				DroppedAttributeCount: 0,
			},
		},
		Status: tracesdk.Status{
			Code:        codes.Error,
			Description: "utterly unrecognized",
		},
		Attributes: []attribute.KeyValue{
			attribute.Int64("timeout_ns", 12e9),
		},
		DroppedAttributes: 1,
		DroppedEvents:     2,
		DroppedLinks:      3,
		Resource: resource.NewWithAttributes(
			"http://example.com/custom-resource-schema",
			attribute.String("rk1", "rv1"),
			attribute.Int64("rk2", 5),
			attribute.StringSlice("rk3", []string{"sv1", "sv2"}),
		),
		InstrumentationLibrary: instrumentation.Scope{
			Name:      "go.opentelemetry.io/test/otel",
			Version:   "v0.0.1",
			SchemaURL: semconv.SchemaURL,
		},
	}

	// Not checking resource as the underlying map of our Resource makes
	// ordering impossible to guarantee on the output. The Resource
	// transform function has unit tests that should suffice.
	expectedSpan := &tracepb.Span{
		TraceId:                []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		SpanId:                 []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		ParentSpanId:           []byte{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		TraceState:             "key1=val1,key2=val2",
		Name:                   spanData.Name,
		Kind:                   tracepb.Span_SPAN_KIND_SERVER,
		StartTimeUnixNano:      uint64(startTime.UnixNano()),
		EndTimeUnixNano:        uint64(endTime.UnixNano()),
		Status:                 status(spanData.Status.Code, spanData.Status.Description),
		Events:                 spanEvents(spanData.Events),
		Links:                  links(spanData.Links),
		Attributes:             KeyValues(spanData.Attributes),
		DroppedAttributesCount: 1,
		DroppedEventsCount:     2,
		DroppedLinksCount:      3,
	}

	got := Spans(tracetest.SpanStubs{spanData}.Snapshots())
	require.Len(t, got, 1)

	assert.Equal(t, got[0].GetResource(), Resource(spanData.Resource))
	assert.Equal(t, got[0].SchemaUrl, spanData.Resource.SchemaURL())
	scopeSpans := got[0].GetScopeSpans()
	require.Len(t, scopeSpans, 1)
	assert.Equal(t, scopeSpans[0].SchemaUrl, spanData.InstrumentationLibrary.SchemaURL)
	assert.Equal(t, scopeSpans[0].GetScope(), InstrumentationScope(spanData.InstrumentationLibrary))
	require.Len(t, scopeSpans[0].Spans, 1)
	actualSpan := scopeSpans[0].Spans[0]

	if diff := cmp.Diff(expectedSpan, actualSpan, cmp.Comparer(proto.Equal)); diff != "" {
		t.Fatalf("transformed span differs %v\n", diff)
	}
}

// Empty parent span ID should be treated as root span.
func TestRootSpanData(t *testing.T) {
	sd := Spans(tracetest.SpanStubs{{}}.Snapshots())
	require.Len(t, sd, 1)
	rs := sd[0]
	scopeSpans := rs.GetScopeSpans()
	require.Len(t, scopeSpans, 1)
	got := scopeSpans[0].GetSpans()[0].GetParentSpanId()

	// Empty means root span.
	assert.Nil(t, got, "incorrect transform of root parent span ID")
}

func TestSpanDataNilResource(t *testing.T) {
	assert.NotPanics(t, func() { Spans(tracetest.SpanStubs{{}}.Snapshots()) })
}
