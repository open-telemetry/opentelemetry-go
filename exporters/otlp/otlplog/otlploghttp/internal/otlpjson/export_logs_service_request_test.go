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
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

var hexPattern = regexp.MustCompile(`^[0-9A-Fa-f]+$`)

func logsForTest() *collogpb.ExportLogsServiceRequest {
	return &collogpb.ExportLogsServiceRequest{
		ResourceLogs: []*logspb.ResourceLogs{{
			Resource: &resourcepb.Resource{
				Attributes: []*commonpb.KeyValue{{
					Key:   "service.name",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "svc"}},
				}},
			},
			ScopeLogs: []*logspb.ScopeLogs{{
				Scope: &commonpb.InstrumentationScope{Name: "lib", Version: "1.0"},
				LogRecords: []*logspb.LogRecord{{
					TimeUnixNano:         1_617_187_200_000_000_000,
					ObservedTimeUnixNano: 1_617_187_200_500_000_000,
					SeverityNumber:       logspb.SeverityNumber_SEVERITY_NUMBER_INFO,
					SeverityText:         "INFO",
					Body: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_StringValue{StringValue: "hello"},
					},
					Attributes: []*commonpb.KeyValue{{
						Key:   "key",
						Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: -42}},
					}},
					Flags: 1,
					TraceId: []byte{
						0x5B, 0x8E, 0xFF, 0xF7, 0x98, 0x03, 0x81, 0x03,
						0xD2, 0x69, 0xB6, 0x33, 0x81, 0x3F, 0xC6, 0x0C,
					},
					SpanId:    []byte{0xEE, 0xE1, 0x9B, 0x7E, 0xC3, 0xC1, 0xB1, 0x74},
					EventName: "evt",
				}},
			}},
		}},
	}
}

func unmarshalGeneric(t *testing.T, data []byte) map[string]any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var m map[string]any
	require.NoError(t, dec.Decode(&m))
	return m
}

func firstLogRecord(t *testing.T, root map[string]any) map[string]any {
	t.Helper()
	rl := root["resourceLogs"].([]any)
	sl := rl[0].(map[string]any)["scopeLogs"].([]any)
	records := sl[0].(map[string]any)["logRecords"].([]any)
	return records[0].(map[string]any)
}

func TestMarshalTraceAndSpanIDsAreHexStrings(t *testing.T) {
	data, err := MarshalExportLogsServiceRequest(logsForTest())
	require.NoError(t, err)

	rec := firstLogRecord(t, unmarshalGeneric(t, data))

	traceID, ok := rec["traceId"].(string)
	require.True(t, ok, "traceId must be a string")
	assert.Len(t, traceID, 32)
	assert.Regexp(t, hexPattern, traceID)
	assert.Equal(t, "5B8EFFF798038103D269B633813FC60C", traceID)

	spanID, ok := rec["spanId"].(string)
	require.True(t, ok, "spanId must be a string")
	assert.Len(t, spanID, 16)
	assert.Regexp(t, hexPattern, spanID)
	assert.Equal(t, "EEE19B7EC3C1B174", spanID)
}

func TestMarshalEnumValuesAreIntegers(t *testing.T) {
	data, err := MarshalExportLogsServiceRequest(logsForTest())
	require.NoError(t, err)

	rec := firstLogRecord(t, unmarshalGeneric(t, data))

	sev, ok := rec["severityNumber"].(json.Number)
	require.True(t, ok, "severityNumber must be a JSON number, got %T", rec["severityNumber"])
	assert.Equal(t, "9", sev.String(), "SEVERITY_NUMBER_INFO = 9")
}

func TestMarshal64BitIntegersAreQuotedStrings(t *testing.T) {
	data, err := MarshalExportLogsServiceRequest(logsForTest())
	require.NoError(t, err)

	rec := firstLogRecord(t, unmarshalGeneric(t, data))
	assert.Equal(t, "1617187200000000000", rec["timeUnixNano"])
	assert.Equal(t, "1617187200500000000", rec["observedTimeUnixNano"])

	attrs := rec["attributes"].([]any)
	val := attrs[0].(map[string]any)["value"].(map[string]any)
	assert.Equal(t, "-42", val["intValue"])
}

func TestMarshalRoundTrip(t *testing.T) {
	original := logsForTest()
	data, err := MarshalExportLogsServiceRequest(original)
	require.NoError(t, err)

	var decoded collogpb.ExportLogsServiceRequest
	require.NoError(t, UnmarshalExportLogsServiceRequest(data, &decoded))
	require.Len(t, decoded.ResourceLogs, 1)
	require.Len(t, decoded.ResourceLogs[0].ScopeLogs, 1)
	require.Len(t, decoded.ResourceLogs[0].ScopeLogs[0].LogRecords, 1)

	got := decoded.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	want := original.ResourceLogs[0].ScopeLogs[0].LogRecords[0]
	assert.Equal(t, want.TimeUnixNano, got.TimeUnixNano)
	assert.Equal(t, want.SeverityNumber, got.SeverityNumber)
	assert.Equal(t, want.SeverityText, got.SeverityText)
	assert.Equal(t, want.TraceId, got.TraceId)
	assert.Equal(t, want.SpanId, got.SpanId)
	assert.Equal(t, want.EventName, got.EventName)
	assert.Equal(t, want.Body.GetStringValue(), got.Body.GetStringValue())
}
