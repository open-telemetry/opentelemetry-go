// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"testing"
	"time"
)

func TestTracesEncoding(t *testing.T) {
	traces := &Traces{
		ResourceSpans: []*ResourceSpans{{}},
	}

	t.Run("CamelCase", runJSONEncodingTests(traces, []byte(`{
		"resourceSpans": [
			{
				"resource": {}
			}
		]
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(traces, []byte(`{
		"resource_spans": [
			{
				"resource": {}
			}
		]
	}`)))
}

func TestResourceSpansEncoding(t *testing.T) {
	rs := &ResourceSpans{
		Resource: Resource{
			Attrs: []Attr{String("key", "val")},
		},
		ScopeSpans: []*ScopeSpans{{}},
		SchemaURL:  schema100,
	}

	t.Run("CamelCase", runJSONEncodingTests(rs, []byte(`{
		"resource": {
			"attributes": [
				{
					"key": "key",
					"value": {
						"stringValue": "val"
					}
				}
			]
		},
		"scopeSpans": [
			{
				"scope": null
			}
		],
		"schemaUrl": "http://go.opentelemetry.io/schema/v1.0.0"
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(rs, []byte(`{
		"resource": {
			"attributes": [
				{
					"key": "key",
					"value": {
						"string_value": "val"
					}
				}
			]
		},
		"scope_spans": [
			{
				"scope": null
			}
		],
		"schema_url": "http://go.opentelemetry.io/schema/v1.0.0"
	}`)))
}

func TestScopeSpansEncoding(t *testing.T) {
	ss := &ScopeSpans{
		Scope: &Scope{Name: "scope"},
		Spans: []*Span{{
			TraceID:   [16]byte{0x1},
			SpanID:    [8]byte{0x2},
			Name:      "A",
			StartTime: y2k,
			EndTime:   y2k.Add(time.Second),
		}, {
			TraceID:   [16]byte{0x1},
			SpanID:    [8]byte{0x3},
			Name:      "B",
			StartTime: y2k.Add(time.Second),
			EndTime:   y2k.Add(2 * time.Second),
		}},
		SchemaURL: schema100,
	}

	t.Run("CamelCase", runJSONEncodingTests(ss, []byte(`{
		"scope": {
			"name": "scope"
		},
		"spans": [
			{
				"traceId": "01000000000000000000000000000000",
				"spanId": "0200000000000000",
				"name": "A",
				"startTimeUnixNano": 946684800000000000,
				"endTimeUnixNano": 946684801000000000
			},
			{
				"traceId": "01000000000000000000000000000000",
				"spanId": "0300000000000000",
				"name": "B",
				"startTimeUnixNano": 946684801000000000,
				"endTimeUnixNano": 946684802000000000
			}
		],
		"schemaUrl": "http://go.opentelemetry.io/schema/v1.0.0"
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(ss, []byte(`{
		"scope": {
			"name": "scope"
		},
		"spans": [
			{
				"trace_id": "01000000000000000000000000000000",
				"span_id": "0200000000000000",
				"name": "A",
				"start_time_unix_nano": 946684800000000000,
				"end_time_unix_nano": 946684801000000000
			},
			{
				"trace_id": "01000000000000000000000000000000",
				"span_id": "0300000000000000",
				"name": "B",
				"start_time_unix_nano": 946684801000000000,
				"end_time_unix_nano": 946684802000000000
			}
		],
		"schema_url": "http://go.opentelemetry.io/schema/v1.0.0"
	}`)))
}
