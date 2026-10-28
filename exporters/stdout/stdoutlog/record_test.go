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
	"strings"
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

func TestExporterValueRepresentation(t *testing.T) {
	values := []struct {
		name  string
		value attribute.Value
		want  string
	}{
		{"NaN", attribute.Float64Value(math.NaN()), `"NaN"`},
		{"positive infinity", attribute.Float64Value(math.Inf(1)), `"Infinity"`},
		{"negative infinity", attribute.Float64Value(math.Inf(-1)), `"-Infinity"`},
		{"finite float", attribute.Float64Value(1e-7), `1e-07`},
		{"int64", attribute.Int64Value(math.MaxInt64), `9223372036854775807`},
		{
			"JSON-looking string",
			attribute.StringValue(`{"admin":true}`),
			`"{\"admin\":true}"`,
		},
		{"byte slice", attribute.ByteSliceValue([]byte{0, 1, 255}), `"AAH/"`},
		{"empty", attribute.Value{}, `null`},
		{"boolean", attribute.BoolValue(true), `true`},
		{"minimum int64", attribute.Int64Value(math.MinInt64), `-9223372036854775808`},
		{"empty string", attribute.StringValue(""), `""`},
		{"empty byte slice", attribute.ByteSliceValue(nil), `""`},
		{"empty float slice", attribute.Float64SliceValue(nil), `[]`},
		{"empty slice", attribute.SliceValue(), `[]`},
		{"empty map", attribute.MapValue(), `{}`},
		{
			"escaped string",
			attribute.StringValue("quote\"\\\n\t\x00<>&\u2028\u2029\xff"),
			`"quote\"\\\n\t\u0000<>&\u2028\u2029\ufffd"`,
		},
		{
			"Type and Value map keys",
			attribute.MapValue(attribute.String("Type", "FLOAT64"), attribute.String("Value", "NaN")),
			`{"Type":"FLOAT64","Value":"NaN"}`,
		},
		{
			name: "finite float slice",
			value: attribute.Float64SliceValue([]float64{
				1.5, math.Copysign(0, -1), math.SmallestNonzeroFloat64, math.MaxFloat64,
			}),
			want: `[1.5,-0,5e-324,1.7976931348623157e+308]`,
		},
		{
			name: "float slice",
			value: attribute.Float64SliceValue(
				[]float64{1.5, math.NaN(), math.Inf(1), math.Inf(-1), math.Copysign(0, -1)},
			),
			want: `[1.5,"NaN","Infinity","-Infinity",-0]`,
		},
		{
			name: "heterogeneous slice",
			value: attribute.SliceValue(
				attribute.Float64Value(math.NaN()), attribute.StringValue("NaN"),
				attribute.Float64Value(math.Inf(1)), attribute.Float64Value(math.Inf(-1)),
				attribute.Float64Value(1.5), attribute.Int64Value(math.MaxInt64),
				attribute.ByteSliceValue([]byte{0, 1, 255}), attribute.Value{},
			),
			want: `["NaN","NaN","Infinity","-Infinity",1.5,9223372036854775807,"AAH/",null]`,
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
			want: `{"Duplicate":"first","a":[{"special":"NaN"},["Infinity","-Infinity"]],"duplicate":"second","z":"last"}`,
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
			want: `[true,[true,false],[-9223372036854775808,9223372036854775807],["Infinity","<tag>"],-0,[],[],{},"",null,"NaN"]`,
		},
	}

	for _, tc := range values {
		t.Run(tc.name, func(t *testing.T) {
			testExporterValueRepresentation(t, tc.value, decodeJSON(t, tc.want))
		})
	}
}

func testExporterValueRepresentation(t *testing.T, value attribute.Value, wantValue any) {
	t.Helper()
	for _, location := range []string{"body", "attribute", "resource", "scope"} {
		t.Run(location, func(t *testing.T) {
			factory := recordFactory("special")
			want := expectedRecord("special")
			kv := attribute.KeyValue{Key: "m", Value: value}
			switch location {
			case "body":
				factory.Body = value
				want["Body"] = wantValue
			case "attribute":
				factory.Attributes = append(factory.Attributes, kv)
				want["Attributes"] = append(want["Attributes"].([]any), map[string]any{"m": wantValue})
			case "resource":
				factory.Resource = resource.NewSchemaless(append(factory.Attributes, kv)...)
				want["Resource"].(map[string]any)["m"] = wantValue
			case "scope":
				factory.InstrumentationScope.Attributes = attribute.NewSet(append(factory.Attributes, kv)...)
				want["Scope"].(map[string]any)["Attributes"].(map[string]any)["m"] = wantValue
			}
			records := []sdklog.Record{
				recordFactory("before").NewRecord(),
				factory.NewRecord(),
				recordFactory("after").NewRecord(),
			}
			wantRecords := []map[string]any{expectedRecord("before"), want, expectedRecord("after")}

			for _, format := range []struct {
				name    string
				options []Option
			}{
				{name: "compact"},
				{name: "pretty", options: []Option{WithPrettyPrint()}},
			} {
				t.Run(format.name, func(t *testing.T) {
					var buf bytes.Buffer
					opts := append([]Option{WithWriter(&buf), WithoutTimestamps()}, format.options...)
					exporter, err := New(opts...)
					require.NoError(t, err)
					require.NoError(t, exporter.Export(t.Context(), records))
					assert.Equal(t, wantRecords, decodeExportedJSON(t, &buf))
				})
			}
		})
	}
}

func TestExporterBodyEscaping(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		t.Run(fmt.Sprintf("pretty=%t", pretty), func(t *testing.T) {
			var buf bytes.Buffer
			opts := []Option{WithWriter(&buf), WithoutTimestamps()}
			if pretty {
				opts = append(opts, WithPrettyPrint())
			}
			exporter, err := New(opts...)
			require.NoError(t, err)
			factory := recordFactory("escaping")
			factory.Body = attribute.StringValue("quote\"\\\n\t\x00<>&\u2028\u2029")
			require.NoError(t, exporter.Export(t.Context(), []sdklog.Record{factory.NewRecord()}))
			var got struct{ Body json.RawMessage }
			require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
			assert.Equal(t, `"quote\"\\\n\t\u0000\u003c\u003e\u0026\u2028\u2029"`, string(got.Body))
		})
	}
}

func recordFactory(name string) logtest.RecordFactory {
	attrs := []attribute.KeyValue{attribute.String("z", "last"), attribute.String("a", "first")}
	return logtest.RecordFactory{
		EventName:    name,
		Body:         attribute.StringValue(name),
		Severity:     log.SeverityWarn,
		SeverityText: "WARN",
		Attributes:   attrs,
		Resource:     resource.NewSchemaless(append([]attribute.KeyValue(nil), attrs...)...),
		InstrumentationScope: &instrumentation.Scope{
			Name:       "scope",
			Version:    "1.0",
			SchemaURL:  "https://example.com/schema",
			Attributes: attribute.NewSet(append([]attribute.KeyValue(nil), attrs...)...),
		},
		TraceID:           trace.TraceID{1},
		SpanID:            trace.SpanID{2},
		TraceFlags:        trace.FlagsSampled,
		DroppedAttributes: 3,
	}
}

func expectedRecord(name string) map[string]any {
	return map[string]any{
		"EventName":    name,
		"Severity":     json.Number("13"),
		"SeverityText": "WARN",
		"Body":         name,
		"Attributes":   []any{map[string]any{"z": "last"}, map[string]any{"a": "first"}},
		"TraceID":      "01000000000000000000000000000000",
		"SpanID":       "0200000000000000",
		"TraceFlags":   "01",
		"Resource":     map[string]any{"a": "first", "z": "last"},
		"Scope": map[string]any{
			"Name":       "scope",
			"Version":    "1.0",
			"SchemaURL":  "https://example.com/schema",
			"Attributes": map[string]any{"a": "first", "z": "last"},
		},
		"DroppedAttributes": json.Number("3"),
	}
}

func decodeExportedJSON(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var records []map[string]any
	dec := json.NewDecoder(buf)
	dec.UseNumber() // Preserve int64 precision and negative zero in comparisons.
	for {
		var record map[string]any
		err := dec.Decode(&record)
		if errors.Is(err, io.EOF) {
			return records
		}
		require.NoError(t, err)
		records = append(records, record)
	}
}

func decodeJSON(t *testing.T, data string) any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(data))
	dec.UseNumber()
	var result any
	require.NoError(t, dec.Decode(&result))
	return result
}

func TestExporterEmptyCollections(t *testing.T) {
	for _, tc := range []struct {
		name string
		res  *resource.Resource
		want any
	}{
		{"nil resource", nil, nil},
		{"empty resource", resource.Empty(), map[string]any{}},
		{"allocated empty resource", resource.NewSchemaless(), map[string]any{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			exporter, err := New(WithWriter(&buf), WithoutTimestamps())
			require.NoError(t, err)
			factory := recordFactory("empty")
			factory.Attributes, factory.Resource = nil, tc.res
			factory.InstrumentationScope.Attributes = attribute.Set{}
			require.NoError(t, exporter.Export(t.Context(), []sdklog.Record{factory.NewRecord()}))
			want := expectedRecord("empty")
			want["Attributes"] = []any{}
			want["Resource"] = tc.want
			want["Scope"].(map[string]any)["Attributes"] = map[string]any{}
			assert.Equal(t, []map[string]any{want}, decodeExportedJSON(t, &buf))
		})
	}
}

func TestExporterTimestampMarshalError(t *testing.T) {
	for _, observed := range []bool{false, true} {
		t.Run(fmt.Sprintf("observed=%t", observed), func(t *testing.T) {
			var buf bytes.Buffer
			exporter, err := New(WithWriter(&buf))
			require.NoError(t, err)
			before := recordFactory("before").NewRecord()
			bad := recordFactory("bad")
			invalid := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
			if observed {
				bad.ObservedTimestamp = invalid
			} else {
				bad.Timestamp = invalid
			}
			after := recordFactory("after").NewRecord()
			err = exporter.Export(t.Context(), []sdklog.Record{before, bad.NewRecord(), after})
			var marshalErr *json.MarshalerError
			require.ErrorAs(t, err, &marshalErr)
			assert.Contains(t, marshalErr.Error(), "year outside of range")
			got := decodeExportedJSON(t, &buf)
			require.Len(t, got, 1)
			assert.Equal(t, "before", got[0]["EventName"])
			require.NoError(t, exporter.Export(t.Context(), []sdklog.Record{after}))
			got = decodeExportedJSON(t, &buf)
			require.Len(t, got, 1)
			assert.Equal(t, "after", got[0]["EventName"])
		})
	}
}
