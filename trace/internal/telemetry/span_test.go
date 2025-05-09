// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"testing"
	"time"
)

func TestSpanEncoding(t *testing.T) {
	span := &Span{
		TraceID:      [16]byte{0x1},
		SpanID:       [8]byte{0x2},
		TraceState:   "test=a",
		ParentSpanID: [8]byte{0x1},
		Flags:        1,
		Name:         "span.a",
		Kind:         SpanKindClient,
		StartTime:    y2k,
		EndTime:      y2k.Add(time.Second),
		Attrs:        []Attr{String("key", "val")},
		DroppedAttrs: 2,
		Events: []*SpanEvent{{
			Name: "name",
		}},
		DroppedEvents: 3,
		Links: []*SpanLink{{
			TraceID: TraceID{0x2},
			SpanID:  SpanID{0x1},
		}},
		DroppedLinks: 4,
		Status: &Status{
			Message: "okay",
			Code:    StatusCodeOK,
		},
	}

	t.Run("CamelCase", runJSONEncodingTests(span, []byte(`{
		"traceId": "01000000000000000000000000000000",
		"spanId": "0200000000000000",
		"traceState": "test=a",
		"flags": 1,
		"name": "span.a",
		"kind": 3,
		"attributes": [
			{
				"key": "key",
				"value": {
					"stringValue": "val"
				}
			}
		],
		"droppedAttributesCount": 2,
		"events": [
			{
				"name": "name"
			}
		],
		"droppedEventsCount": 3,
		"links": [
			{
				"traceId": "02000000000000000000000000000000",
				"spanId": "0100000000000000"
			}
		],
		"droppedLinksCount": 4,
		"status": {
			"message": "okay",
			"code": 1
		},
		"parentSpanId": "0100000000000000",
		"startTimeUnixNano": 946684800000000000,
		"endTimeUnixNano": 946684801000000000
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(span, []byte(`{
		"trace_id": "01000000000000000000000000000000",
		"span_id": "0200000000000000",
		"trace_state": "test=a",
		"flags": 1,
		"name": "span.a",
		"kind": 3,
		"attributes": [
			{
				"key": "key",
				"value": {
					"string_value": "val"
				}
			}
		],
		"dropped_attributes_count": 2,
		"events": [
			{
				"name": "name"
			}
		],
		"dropped_events_count": 3,
		"links": [
			{
				"trace_id": "02000000000000000000000000000000",
				"span_id": "0100000000000000"
			}
		],
		"dropped_links_count": 4,
		"status": {
			"message": "okay",
			"code": 1
		},
		"parent_span_id": "0100000000000000",
		"start_time_unix_nano": 946684800000000000,
		"end_time_unix_nano": 946684801000000000
	}`)))

	t.Run("RequiredFields", runJSONMarshalTest(new(Span), []byte(`{
		"traceId": "",
		"spanId": "",
		"name": ""
	}`)))
}

func TestSpanEventEncoding(t *testing.T) {
	event := &SpanEvent{
		Time:         y2k.Add(10 * time.Microsecond),
		Name:         "span.event",
		Attrs:        []Attr{Float64("impact", 0.4372)},
		DroppedAttrs: 2,
	}

	t.Run("CamelCase", runJSONEncodingTests(event, []byte(`{
		"name": "span.event",
		"attributes": [
			{
				"key": "impact",
				"value": {
					"doubleValue": 0.4372
				}
			}
		],
		"droppedAttributesCount": 2,
		"timeUnixNano": 946684800000010000
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(event, []byte(`{
		"name": "span.event",
		"attributes": [
			{
				"key": "impact",
				"value": {
					"double_value": 0.4372
				}
			}
		],
		"dropped_attributes_count": 2,
		"time_unix_nano": 946684800000010000
	}`)))
}

func TestSpanLinkEncoding(t *testing.T) {
	link := &SpanLink{
		TraceID:      TraceID{0x2},
		SpanID:       SpanID{0x1},
		TraceState:   "test=green",
		Attrs:        []Attr{Int("queue", 17)},
		DroppedAttrs: 8,
		Flags:        1,
	}

	t.Run("CamelCase", runJSONEncodingTests(link, []byte(`{
		"traceId": "02000000000000000000000000000000",
		"spanId": "0100000000000000",
		"traceState": "test=green",
		"attributes": [
			{
				"key": "queue",
				"value": {
					"intValue": "17"
				}
			}
		],
		"droppedAttributesCount": 8,
		"flags": 1
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(link, []byte(`{
		"trace_id": "02000000000000000000000000000000",
		"span_id": "0100000000000000",
		"trace_state": "test=green",
		"attributes": [
			{
				"key": "queue",
				"value": {
					"int_value": "17"
				}
			}
		],
		"dropped_attributes_count": 8,
		"flags": 1
	}`)))
}
