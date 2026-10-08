// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func TestExporterNonFiniteValues(t *testing.T) {
	values := []struct {
		name  string
		value attribute.Value
		want  string
	}{
		{"NaN", attribute.Float64Value(math.NaN()), `{"Type":"FLOAT64","Value":"NaN"}`},
		{"positive infinity", attribute.Float64Value(math.Inf(1)), `{"Type":"FLOAT64","Value":"Infinity"}`},
		{"negative infinity", attribute.Float64Value(math.Inf(-1)), `{"Type":"FLOAT64","Value":"-Infinity"}`},
		{
			name:  "float slice",
			value: attribute.Float64SliceValue([]float64{1.5, math.NaN(), math.Inf(1), math.Inf(-1), math.Copysign(0, -1)}),
			want:  `{"Type":"FLOAT64SLICE","Value":[1.5,"NaN","Infinity","-Infinity",-0]}`,
		},
		{
			name: "heterogeneous slice",
			value: attribute.SliceValue(
				attribute.Float64Value(math.NaN()), attribute.StringValue("NaN"),
				attribute.Float64Value(math.Inf(1)), attribute.Float64Value(math.Inf(-1)),
				attribute.Float64Value(1.5), attribute.Int64Value(math.MaxInt64),
				attribute.ByteSliceValue([]byte{0, 1, 255}), attribute.Value{},
			),
			want: `{"Type":"SLICE","Value":[
				{"Type":"FLOAT64","Value":"NaN"},{"Type":"STRING","Value":"NaN"},
				{"Type":"FLOAT64","Value":"Infinity"},{"Type":"FLOAT64","Value":"-Infinity"},
				{"Type":"FLOAT64","Value":1.5},{"Type":"INT64","Value":9223372036854775807},
				{"Type":"BYTESLICE","Value":"AAH/"},{"Type":"EMPTY","Value":null}]}`,
		},
		{
			name: "nested map and slice",
			value: attribute.MapValue(
				attribute.String("z", "last"),
				attribute.KeyValue{Key: "a", Value: attribute.SliceValue(
					attribute.MapValue(attribute.Float64("special", math.NaN())),
					attribute.Float64SliceValue([]float64{math.Inf(1), math.Inf(-1)}),
				)},
				attribute.String("Duplicate", "first"), attribute.String("duplicate", "second"),
			),
			want: `{"Type":"MAP","Value":[
				{"Key":"Duplicate","Value":{"Type":"STRING","Value":"first"}},
				{"Key":"a","Value":{"Type":"SLICE","Value":[
					{"Type":"MAP","Value":[{"Key":"special","Value":{"Type":"FLOAT64","Value":"NaN"}}]},
					{"Type":"FLOAT64SLICE","Value":["Infinity","-Infinity"]}]}},
				{"Key":"duplicate","Value":{"Type":"STRING","Value":"second"}},
				{"Key":"z","Value":{"Type":"STRING","Value":"last"}}]}`,
		},
		{
			name: "other types",
			value: attribute.SliceValue(
				attribute.BoolValue(true), attribute.BoolSliceValue([]bool{true, false}),
				attribute.Int64SliceValue([]int64{math.MinInt64, math.MaxInt64}),
				attribute.StringSliceValue([]string{"Infinity", "<tag>"}),
				attribute.Float64Value(math.Copysign(0, -1)),
				attribute.Float64SliceValue(nil), attribute.SliceValue(), attribute.MapValue(),
				attribute.ByteSliceValue(nil), attribute.Value{}, attribute.Float64Value(math.NaN()),
			),
			want: `{"Type":"SLICE","Value":[
				{"Type":"BOOL","Value":true},{"Type":"BOOLSLICE","Value":[true,false]},
				{"Type":"INT64SLICE","Value":[-9223372036854775808,9223372036854775807]},
				{"Type":"STRINGSLICE","Value":["Infinity","<tag>"]},{"Type":"FLOAT64","Value":-0},
				{"Type":"FLOAT64SLICE","Value":[]},{"Type":"SLICE","Value":[]},{"Type":"MAP","Value":[]},
				{"Type":"BYTESLICE","Value":""},{"Type":"EMPTY","Value":null},{"Type":"FLOAT64","Value":"NaN"}]}`,
		},
	}

	for _, pretty := range []bool{false, true} {
		for _, location := range []string{"body", "attribute", "resource", "scope"} {
			for _, tc := range values {
				t.Run(fmt.Sprintf("pretty=%t/%s/%s", pretty, location, tc.name), func(t *testing.T) {
					var buf bytes.Buffer
					opts := []Option{WithWriter(&buf), WithoutTimestamps()}
					if pretty {
						opts = append(opts, WithPrettyPrint())
					}
					exporter, err := New(opts...)
					require.NoError(t, err)

					factory := nonFiniteRecordFactory("special")
					body, attrs, res, scope := ordinaryJSONFields("special")
					kv := attribute.KeyValue{Key: "m", Value: tc.value}
					wantKV := `{"Key":"m","Value":` + tc.want + `}`
					switch location {
					case "body":
						factory.Body, body = tc.value, tc.want
					case "attribute":
						factory.Attributes = append(factory.Attributes, kv)
						attrs = attrs[:len(attrs)-1] + "," + wantKV + "]"
					case "resource":
						factory.Resource = resource.NewSchemaless(append(factory.Attributes, kv)...)
						res = sortedJSONAttributes(wantKV)
					case "scope":
						factory.InstrumentationScope.Attributes = attribute.NewSet(append(factory.Attributes, kv)...)
						scope = sortedJSONAttributes(wantKV)
					}
					records := []sdklog.Record{
						nonFiniteRecordFactory("before").NewRecord(), factory.NewRecord(), nonFiniteRecordFactory("after").NewRecord(),
					}
					require.NoError(t, exporter.Export(t.Context(), records))
					got := decodeExportedJSON(t, &buf)
					require.Len(t, got, 3)
					assertDebugJSON(t, expectedOrdinaryRecordJSON("before"), got[0])
					assertDebugJSON(t, expectedRecordJSON("special", body, attrs, res, scope), got[1])
					assertDebugJSON(t, expectedOrdinaryRecordJSON("after"), got[2])
				})
			}
		}
	}
}

func nonFiniteRecordFactory(name string) logtest.RecordFactory {
	attrs := []attribute.KeyValue{attribute.String("z", "last"), attribute.String("a", "first")}
	return logtest.RecordFactory{
		EventName: name, Body: attribute.StringValue(name), Severity: log.SeverityWarn, SeverityText: "WARN",
		Attributes: attrs, Resource: resource.NewSchemaless(append([]attribute.KeyValue(nil), attrs...)...),
		InstrumentationScope: &instrumentation.Scope{
			Name: "scope", Version: "1.0", SchemaURL: "https://example.com/schema",
			Attributes: attribute.NewSet(append([]attribute.KeyValue(nil), attrs...)...),
		},
		TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled,
		DroppedAttributes: 3,
	}
}

func ordinaryJSONFields(name string) (body, attrs, res, scope string) {
	body = `{"Type":"STRING","Value":"` + name + `"}`
	attrs = `[{"Key":"z","Value":{"Type":"STRING","Value":"last"}},{"Key":"a","Value":{"Type":"STRING","Value":"first"}}]`
	res = sortedJSONAttributes("")
	return body, attrs, res, res
}

func sortedJSONAttributes(middle string) string {
	if middle != "" {
		middle += ","
	}
	return `[{"Key":"a","Value":{"Type":"STRING","Value":"first"}},` + middle + `{"Key":"z","Value":{"Type":"STRING","Value":"last"}}]`
}

func expectedRecordJSON(name, body, attrs, res, scope string) string {
	return fmt.Sprintf(`{"EventName":%q,"Severity":13,"SeverityText":"WARN","Body":%s,"Attributes":%s,
		"TraceID":"01000000000000000000000000000000","SpanID":"0200000000000000","TraceFlags":"01",
		"Resource":%s,"Scope":{"Name":"scope","Version":"1.0","SchemaURL":"https://example.com/schema","Attributes":%s},
		"DroppedAttributes":3}`, name, body, attrs, res, scope)
}

func expectedOrdinaryRecordJSON(name string) string {
	body, attrs, res, scope := ordinaryJSONFields(name)
	return expectedRecordJSON(name, body, attrs, res, scope)
}

func decodeExportedJSON(t *testing.T, buf *bytes.Buffer) []json.RawMessage {
	t.Helper()
	var records []json.RawMessage
	dec := json.NewDecoder(buf)
	for {
		var record json.RawMessage
		err := dec.Decode(&record)
		if errors.Is(err, io.EOF) {
			return records
		}
		require.NoError(t, err)
		records = append(records, record)
	}
}

func assertDebugJSON(t *testing.T, want string, got json.RawMessage) {
	t.Helper()
	decode := func(data []byte) any {
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber() // Preserve int64 precision and negative zero in comparisons.
		var result any
		require.NoError(t, dec.Decode(&result))
		return result
	}
	assert.Equal(t, decode([]byte(want)), decode(got))
}

func TestExporterTimestampMarshalError(t *testing.T) {
	for _, observed := range []bool{false, true} {
		t.Run(fmt.Sprintf("observed=%t", observed), func(t *testing.T) {
			var buf bytes.Buffer
			exporter, err := New(WithWriter(&buf))
			require.NoError(t, err)
			before := nonFiniteRecordFactory("before").NewRecord()
			bad := nonFiniteRecordFactory("bad")
			invalid := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
			if observed {
				bad.ObservedTimestamp = invalid
			} else {
				bad.Timestamp = invalid
			}
			after := nonFiniteRecordFactory("after").NewRecord()
			err = exporter.Export(t.Context(), []sdklog.Record{before, bad.NewRecord(), after})
			var marshalErr *json.MarshalerError
			require.ErrorAs(t, err, &marshalErr)
			assert.Contains(t, marshalErr.Error(), "year outside of range")
			got := decodeExportedJSON(t, &buf)
			require.Len(t, got, 1)
			assert.Contains(t, string(got[0]), `"EventName":"before"`)
			require.NoError(t, exporter.Export(t.Context(), []sdklog.Record{after}))
			got = decodeExportedJSON(t, &buf)
			require.Len(t, got, 1)
			assert.Contains(t, string(got[0]), `"EventName":"after"`)
		})
	}
}

type nonFiniteFailingWriter struct {
	err   error
	calls int
}

func (w *nonFiniteFailingWriter) Write([]byte) (int, error) {
	w.calls++
	return 0, w.err
}

func TestExporterWriterError(t *testing.T) {
	want := errors.New("writer unavailable")
	writer := &nonFiniteFailingWriter{err: want}
	exporter, err := New(WithWriter(writer))
	require.NoError(t, err)
	err = exporter.Export(t.Context(), []sdklog.Record{
		nonFiniteRecordFactory("before").NewRecord(), nonFiniteRecordFactory("after").NewRecord(),
	})
	assert.ErrorIs(t, err, want)
	assert.Equal(t, 1, writer.calls)
}
