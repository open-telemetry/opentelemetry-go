// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpjson

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

var hexPattern = regexp.MustCompile(`^[0-9A-Fa-f]+$`)

func spanForTest() *coltracepb.ExportTraceServiceRequest {
	return &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: []*tracepb.ResourceSpans{{
			Resource: &resourcepb.Resource{
				Attributes: []*commonpb.KeyValue{{
					Key:   "service.name",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "svc"}},
				}},
				DroppedAttributesCount: 1,
				EntityRefs: []*commonpb.EntityRef{{
					SchemaUrl:       "http://example.com",
					Type:            "service.instance",
					IdKeys:          []string{"service.instance.id", "service.name", "service.namespace"},
					DescriptionKeys: []string{"service.version"},
				}},
			},
			ScopeSpans: []*tracepb.ScopeSpans{{
				Scope: &commonpb.InstrumentationScope{Name: "lib", Version: "1.0"},
				Spans: []*tracepb.Span{
					{
						TraceId: []byte{
							0x5B, 0x8E, 0xFF, 0xF7, 0x98, 0x03, 0x81, 0x03,
							0xD2, 0x69, 0xB6, 0x33, 0x81, 0x3F, 0xC6, 0x0C,
						},
						SpanId:                 []byte{0xEE, 0xE1, 0x9B, 0x7E, 0xC3, 0xC1, 0xB1, 0x74},
						ParentSpanId:           []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
						Name:                   "op",
						Kind:                   tracepb.Span_SPAN_KIND_SERVER,
						StartTimeUnixNano:      1_617_187_200_000_000_000,
						EndTimeUnixNano:        1_617_187_201_000_000_000,
						DroppedAttributesCount: 3,
						Attributes: []*commonpb.KeyValue{{
							Key:   "key",
							Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: -42}},
						}},
						Events: []*tracepb.Span_Event{{
							TimeUnixNano: 1_617_187_200_500_000_000,
							Name:         "evt",
						}},
						Links: []*tracepb.Span_Link{{
							TraceId: []byte{
								0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
								0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
							},
							SpanId: []byte{
								0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11,
							},
						}},
						Status: &tracepb.Status{Code: tracepb.Status_STATUS_CODE_OK, Message: "ok"},
					},
				},
			}},
		}},
	}
}

// unmarshalGeneric parses JSON into nested maps, preserving numbers as json.Number.
func unmarshalGeneric(t *testing.T, data []byte) map[string]any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var m map[string]any
	require.NoError(t, dec.Decode(&m))
	return m
}

// firstSpan drills into the generic JSON map and returns the first span object.
func firstSpan(t *testing.T, root map[string]any) map[string]any {
	t.Helper()
	rs := root["resourceSpans"].([]any)
	ss := rs[0].(map[string]any)["scopeSpans"].([]any)
	spans := ss[0].(map[string]any)["spans"].([]any)
	return spans[0].(map[string]any)
}

func TestMarshalTraceAndSpanIDsAreHexStrings(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))

	traceID, ok := span["traceId"].(string)
	require.True(t, ok, "traceId must be a string")
	assert.Len(t, traceID, 32, "traceId must be 32 hex chars (16 bytes)")
	assert.Regexp(t, hexPattern, traceID)
	assert.Equal(t, "5B8EFFF798038103D269B633813FC60C", traceID)

	spanID, ok := span["spanId"].(string)
	require.True(t, ok, "spanId must be a string")
	assert.Len(t, spanID, 16, "spanId must be 16 hex chars (8 bytes)")
	assert.Regexp(t, hexPattern, spanID)
	assert.Equal(t, "EEE19B7EC3C1B174", spanID)

	parentSpanID, ok := span["parentSpanId"].(string)
	require.True(t, ok, "parentSpanId must be a string")
	assert.Len(t, parentSpanID, 16)
	assert.Regexp(t, hexPattern, parentSpanID)
}

func TestMarshalTraceAndSpanIDsInLinksAreHexStrings(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))
	links := span["links"].([]any)
	link := links[0].(map[string]any)

	traceID, ok := link["traceId"].(string)
	require.True(t, ok)
	assert.Len(t, traceID, 32)
	assert.Regexp(t, hexPattern, traceID)

	spanID, ok := link["spanId"].(string)
	require.True(t, ok)
	assert.Len(t, spanID, 16)
	assert.Regexp(t, hexPattern, spanID)
}

func TestMarshalEnumValuesAreIntegers(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))

	kind, ok := span["kind"].(json.Number)
	require.True(t, ok, "kind must be a JSON number, got %T", span["kind"])
	assert.Equal(t, "2", kind.String(), "SPAN_KIND_SERVER = 2")

	status := span["status"].(map[string]any)
	code, ok := status["code"].(json.Number)
	require.True(t, ok, "status.code must be a JSON number, got %T", status["code"])
	assert.Equal(t, "1", code.String(), "STATUS_CODE_OK = 1")
}

func TestMarshalUint64AsDecimalString(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))

	start, ok := span["startTimeUnixNano"].(string)
	require.True(t, ok, "startTimeUnixNano must be a JSON string, got %T", span["startTimeUnixNano"])
	assert.Equal(t, "1617187200000000000", start)

	end, ok := span["endTimeUnixNano"].(string)
	require.True(t, ok, "endTimeUnixNano must be a JSON string, got %T", span["endTimeUnixNano"])
	assert.Equal(t, "1617187201000000000", end)

	events := span["events"].([]any)
	evt := events[0].(map[string]any)
	evtTime, ok := evt["timeUnixNano"].(string)
	require.True(t, ok, "event timeUnixNano must be a JSON string")
	assert.Equal(t, "1617187200500000000", evtTime)
}

func TestMarshalInt64AsDecimalString(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))
	attrs := span["attributes"].([]any)
	attr := attrs[0].(map[string]any)
	val := attr["value"].(map[string]any)
	intVal, ok := val["intValue"].(string)
	require.True(t, ok, "intValue must be a JSON string, got %T", val["intValue"])
	assert.Equal(t, "-42", intVal)
}

func TestMarshalFieldNamesAreCamelCase(t *testing.T) {
	data, err := MarshalExportTraceServiceRequest(spanForTest())
	require.NoError(t, err)

	root := unmarshalGeneric(t, data)

	_, ok := root["resourceSpans"]
	assert.True(t, ok, "must use resourceSpans, not resource_spans")

	rs := root["resourceSpans"].([]any)[0].(map[string]any)
	_, ok = rs["scopeSpans"]
	assert.True(t, ok, "must use scopeSpans, not scope_spans")

	span := firstSpan(t, root)

	// Verify present fields use camelCase.
	for _, key := range []string{"traceId", "spanId", "parentSpanId", "startTimeUnixNano", "endTimeUnixNano"} {
		_, ok = span[key]
		assert.True(t, ok, "expected camelCase field %q", key)
	}

	// Verify droppedAttributesCount uses camelCase (set to 3 in test data).
	_, ok = span["droppedAttributesCount"]
	assert.True(t, ok, "expected camelCase field droppedAttributesCount")

	// Verify no snake_case fields are present.
	for _, snakeCase := range []string{"trace_id", "span_id", "parent_span_id", "start_time_unix_nano", "dropped_attributes_count"} {
		_, ok = span[snakeCase]
		assert.False(t, ok, "must not use snake_case field %q", snakeCase)
	}
}

func TestMarshalOmitsDefaultValues(t *testing.T) {
	req := &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: []*tracepb.ResourceSpans{{
			ScopeSpans: []*tracepb.ScopeSpans{{
				Spans: []*tracepb.Span{{
					TraceId: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					SpanId:  []byte{0, 0, 0, 0, 0, 0, 0, 1},
					Name:    "minimal",
				}},
			}},
		}},
	}

	data, err := MarshalExportTraceServiceRequest(req)
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))

	for _, key := range []string{
		"parentSpanId", "status", "events", "links", "attributes",
		"traceState", "flags", "droppedAttributesCount",
		"droppedEventsCount", "droppedLinksCount",
	} {
		_, present := span[key]
		assert.False(t, present, "zero-value field %q should be omitted", key)
	}
}

func TestMarshalRootSpanOmitsParentSpanID(t *testing.T) {
	req := &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: []*tracepb.ResourceSpans{{
			ScopeSpans: []*tracepb.ScopeSpans{{
				Spans: []*tracepb.Span{{
					TraceId: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					SpanId:  []byte{0, 0, 0, 0, 0, 0, 0, 1},
				}},
			}},
		}},
	}

	data, err := MarshalExportTraceServiceRequest(req)
	require.NoError(t, err)

	span := firstSpan(t, unmarshalGeneric(t, data))
	_, present := span["parentSpanId"]
	assert.False(t, present, "root spans must not have parentSpanId")
}

func TestUnmarshalIgnoresUnknownFields(t *testing.T) {
	input := `{
		"resourceSpans": [{
			"unknownTopLevel": true,
			"scopeSpans": [{
				"spans": [{
					"traceId": "5B8EFFF798038103D269B633813FC60C",
					"spanId": "EEE19B7EC3C1B174",
					"name": "op",
					"futureField": {"nested": 123}
				}]
			}]
		}]
	}`

	var req coltracepb.ExportTraceServiceRequest
	require.NoError(t, UnmarshalExportTraceServiceRequest([]byte(input), &req))

	require.Len(t, req.ResourceSpans, 1)
	s := req.ResourceSpans[0].ScopeSpans[0].Spans[0]
	assert.Equal(t, "op", s.Name)
	assert.Equal(t,
		[]byte{0x5B, 0x8E, 0xFF, 0xF7, 0x98, 0x03, 0x81, 0x03, 0xD2, 0x69, 0xB6, 0x33, 0x81, 0x3F, 0xC6, 0x0C},
		s.TraceId,
	)
}

func TestUnmarshalAcceptsBothHexCases(t *testing.T) {
	for _, tc := range []struct {
		name    string
		traceID string
		spanID  string
	}{
		{"uppercase", "5B8EFFF798038103D269B633813FC60C", "EEE19B7EC3C1B174"},
		{"lowercase", "5b8efff798038103d269b633813fc60c", "eee19b7ec3c1b174"},
		{"mixed", "5b8eFFF798038103d269B633813fc60c", "eeE19B7eC3c1b174"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := `{"resourceSpans":[{"scopeSpans":[{"spans":[{
				"traceId":"` + tc.traceID + `",
				"spanId":"` + tc.spanID + `"
			}]}]}]}`

			var req coltracepb.ExportTraceServiceRequest
			require.NoError(t, UnmarshalExportTraceServiceRequest([]byte(input), &req))

			s := req.ResourceSpans[0].ScopeSpans[0].Spans[0]
			assert.Equal(t,
				[]byte{0x5B, 0x8E, 0xFF, 0xF7, 0x98, 0x03, 0x81, 0x03, 0xD2, 0x69, 0xB6, 0x33, 0x81, 0x3F, 0xC6, 0x0C},
				s.TraceId,
			)
			assert.Equal(t, []byte{0xEE, 0xE1, 0x9B, 0x7E, 0xC3, 0xC1, 0xB1, 0x74}, s.SpanId)
		})
	}
}

func TestUnmarshalAcceptsIntAndUint64AsStringOrNumber(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
	}{
		{"uint64_as_string", `"1617187200000000000"`},
		{"uint64_as_number", `1617187200000000000`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var i Uint64
			require.NoError(t, json.Unmarshal([]byte(tc.value), &i))
			assert.Equal(t, uint64(1617187200000000000), uint64(i))
		})
	}

	for _, tc := range []struct {
		name  string
		value string
	}{
		{"int64_as_string", `"-42"`},
		{"int64_as_number", `-42`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var i Int64
			require.NoError(t, json.Unmarshal([]byte(tc.value), &i))
			assert.Equal(t, int64(-42), int64(i))
		})
	}
}

func TestAnyValueRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proto *commonpb.AnyValue
		check func(t *testing.T, got *commonpb.AnyValue)
	}{
		{
			name:  "string",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "hello"}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				assert.Equal(t, "hello", got.GetStringValue())
			},
		},
		{
			name:  "bool_true",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_BoolValue{BoolValue: true}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				assert.True(t, got.GetBoolValue())
			},
		},
		{
			name:  "bool_false",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_BoolValue{BoolValue: false}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				_, ok := got.Value.(*commonpb.AnyValue_BoolValue)
				require.True(t, ok, "expected BoolValue variant")
				assert.False(t, got.GetBoolValue())
			},
		},
		{
			name:  "int",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 42}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				assert.Equal(t, int64(42), got.GetIntValue())
			},
		},
		{
			name:  "double",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 3.14}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				assert.InEpsilon(t, 3.14, got.GetDoubleValue(), 1e-9)
			},
		},
		{
			name: "bytes",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_BytesValue{
				BytesValue: []byte{0xDE, 0xAD, 0xBE, 0xEF},
			}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				assert.Equal(t, []byte{0xDE, 0xAD, 0xBE, 0xEF}, got.GetBytesValue())
			},
		},
		{
			name: "array_of_mixed_values",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_ArrayValue{
				ArrayValue: &commonpb.ArrayValue{Values: []*commonpb.AnyValue{
					{Value: &commonpb.AnyValue_StringValue{StringValue: "a"}},
					{Value: &commonpb.AnyValue_IntValue{IntValue: 1}},
					{Value: &commonpb.AnyValue_BoolValue{BoolValue: true}},
				}},
			}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				arr := got.GetArrayValue()
				require.NotNil(t, arr)
				require.Len(t, arr.Values, 3)
				assert.Equal(t, "a", arr.Values[0].GetStringValue())
				assert.Equal(t, int64(1), arr.Values[1].GetIntValue())
				assert.True(t, arr.Values[2].GetBoolValue())
			},
		},
		{
			name: "kvlist",
			proto: &commonpb.AnyValue{Value: &commonpb.AnyValue_KvlistValue{
				KvlistValue: &commonpb.KeyValueList{Values: []*commonpb.KeyValue{
					{
						Key:   "nested_str",
						Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "v"}},
					},
					{
						Key:   "nested_int",
						Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 99}},
					},
				}},
			}},
			check: func(t *testing.T, got *commonpb.AnyValue) {
				kvl := got.GetKvlistValue()
				require.NotNil(t, kvl)
				require.Len(t, kvl.Values, 2)
				assert.Equal(t, "nested_str", kvl.Values[0].Key)
				assert.Equal(t, "v", kvl.Values[0].Value.GetStringValue())
				assert.Equal(t, "nested_int", kvl.Values[1].Key)
				assert.Equal(t, int64(99), kvl.Values[1].Value.GetIntValue())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &coltracepb.ExportTraceServiceRequest{
				ResourceSpans: []*tracepb.ResourceSpans{{
					ScopeSpans: []*tracepb.ScopeSpans{{
						Spans: []*tracepb.Span{{
							TraceId: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
							SpanId:  []byte{0, 0, 0, 0, 0, 0, 0, 1},
							Attributes: []*commonpb.KeyValue{{
								Key:   "test",
								Value: tc.proto,
							}},
						}},
					}},
				}},
			}

			data, err := MarshalExportTraceServiceRequest(req)
			require.NoError(t, err)

			var decoded coltracepb.ExportTraceServiceRequest
			require.NoError(t, UnmarshalExportTraceServiceRequest(data, &decoded))

			attrs := decoded.ResourceSpans[0].ScopeSpans[0].Spans[0].Attributes
			require.Len(t, attrs, 1)
			assert.Equal(t, "test", attrs[0].Key)
			tc.check(t, attrs[0].Value)
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := spanForTest()

	data, err := MarshalExportTraceServiceRequest(original)
	require.NoError(t, err)

	var decoded coltracepb.ExportTraceServiceRequest
	require.NoError(t, UnmarshalExportTraceServiceRequest(data, &decoded))
	require.NotEmpty(t, original.ResourceSpans)

	orig := original.ResourceSpans[0].ScopeSpans[0].Spans[0]
	got := decoded.ResourceSpans[0].ScopeSpans[0].Spans[0]

	assert.Equal(t, orig.TraceId, got.TraceId)
	assert.Equal(t, orig.SpanId, got.SpanId)
	assert.Equal(t, orig.ParentSpanId, got.ParentSpanId)
	assert.Equal(t, orig.Name, got.Name)
	assert.Equal(t, orig.Kind, got.Kind)
	assert.Equal(t, orig.StartTimeUnixNano, got.StartTimeUnixNano)
	assert.Equal(t, orig.EndTimeUnixNano, got.EndTimeUnixNano)
	assert.Equal(t, orig.DroppedAttributesCount, got.DroppedAttributesCount)

	require.Len(t, got.Attributes, 1)
	assert.Equal(t, "key", got.Attributes[0].Key)
	assert.Equal(t, int64(-42), got.Attributes[0].Value.GetIntValue())

	require.Len(t, got.Events, 1)
	assert.Equal(t, "evt", got.Events[0].Name)
	assert.Equal(t, orig.Events[0].TimeUnixNano, got.Events[0].TimeUnixNano)

	require.Len(t, got.Links, 1)
	assert.Equal(t, orig.Links[0].TraceId, got.Links[0].TraceId)
	assert.Equal(t, orig.Links[0].SpanId, got.Links[0].SpanId)

	assert.Equal(t, orig.Status.Code, got.Status.Code)
	assert.Equal(t, orig.Status.Message, got.Status.Message)

	rs := decoded.ResourceSpans[0]
	assert.Equal(t, "service.name", rs.Resource.Attributes[0].Key)
	assert.Equal(t, "svc", rs.Resource.Attributes[0].Value.GetStringValue())
	assert.Equal(t, uint32(1), rs.Resource.DroppedAttributesCount)

	require.NotEmpty(t, rs.Resource.EntityRefs)
	assert.Equal(t, "http://example.com", rs.Resource.EntityRefs[0].SchemaUrl)
	assert.Equal(t, "service.instance", rs.Resource.EntityRefs[0].Type)
	assert.Equal(t,
		[]string{"service.instance.id", "service.name", "service.namespace"},
		rs.Resource.EntityRefs[0].IdKeys,
	)
	assert.Equal(t, []string{"service.version"}, rs.Resource.EntityRefs[0].DescriptionKeys)

	assert.Equal(t, "lib", rs.ScopeSpans[0].Scope.Name)
	assert.Equal(t, "1.0", rs.ScopeSpans[0].Scope.Version)
}
