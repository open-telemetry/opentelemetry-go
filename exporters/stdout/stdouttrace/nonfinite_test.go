// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type attributeShape struct {
	name  string
	value attribute.Value
	want  any
}

func attributeShapes(value float64, want any) []attributeShape {
	return []attributeShape{
		{"scalar", attribute.Float64Value(value), want},
		{"homogeneous", attribute.Float64SliceValue([]float64{1.5, value, -2.5}), []any{1.5, want, -2.5}},
		{
			"heterogeneous",
			attribute.SliceValue(attribute.StringValue("text"), attribute.Float64Value(value), attribute.Int64Value(42)),
			[]any{"text", want, float64(42)},
		},
		{
			"nested",
			attribute.MapValue(attribute.KeyValue{
				Key: "nested", Value: attribute.SliceValue(attribute.MapValue(attribute.Float64("float", value))),
			}),
			map[string]any{"nested": []any{map[string]any{"float": want}}},
		},
	}
}

func wantAttributes(location string, value any) any {
	attr := map[string]any{"value": value}
	if location == "resource" || location == "scope" {
		return attr
	}
	return []any{attr}
}

func spanWithAttribute(location string, value attribute.Value) tracetest.SpanStub {
	attrs := []attribute.KeyValue{{Key: "value", Value: value}}
	stub := tracetest.SpanStub{Name: "nonfinite"}
	switch location {
	case "span":
		stub.Attributes = attrs
	case "resource":
		stub.Resource = resource.NewSchemaless(attrs...)
	case "event":
		stub.Events = []sdktrace.Event{
			{Name: "event", Attributes: attrs, DroppedAttributeCount: 2, Time: time.Unix(1, 0)},
		}
	case "link":
		stub.Links = []sdktrace.Link{{Attributes: attrs, DroppedAttributeCount: 3}}
	case "scope":
		stub.InstrumentationScope = instrumentation.Scope{
			Name:       "scope",
			Version:    "1.2.3",
			SchemaURL:  "https://example.com/schema",
			Attributes: attribute.NewSet(attrs...),
		}
	}
	return stub
}

func exportedAttributes(t *testing.T, span map[string]any, location string) any {
	t.Helper()
	switch location {
	case "span":
		return span["Attributes"]
	case "resource":
		return span["Resource"]
	case "event":
		return span["Events"].([]any)[0].(map[string]any)["Attributes"]
	case "link":
		return span["Links"].([]any)[0].(map[string]any)["Attributes"]
	case "scope":
		scope := span["InstrumentationScope"].(map[string]any)
		assert.Equal(t, "scope", scope["Name"])
		assert.Equal(t, "1.2.3", scope["Version"])
		assert.Equal(t, "https://example.com/schema", scope["SchemaURL"])
		assert.Equal(t, scope, span["InstrumentationLibrary"])
		return scope["Attributes"]
	default:
		t.Fatalf("unknown location: %s", location)
		return nil
	}
}

func TestExporterNonFiniteAttributes(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		for _, location := range []string{"span", "resource", "event", "link", "scope"} {
			for _, number := range []struct {
				name  string
				value float64
				want  any
			}{
				{"NaN", math.NaN(), "NaN"}, {"Infinity", math.Inf(1), "Infinity"}, {"-Infinity", math.Inf(-1), "-Infinity"}, {"finite", 1.25, 1.25},
			} {
				for _, shape := range attributeShapes(number.value, number.want) {
					name := location + "/" + number.name + "/" + shape.name
					if pretty {
						name += "/pretty"
					}
					t.Run(name, func(t *testing.T) {
						stub := spanWithAttribute(location, shape.value)
						span := stub.Snapshot()
						before := shape.value.String()
						var buf bytes.Buffer
						opts := []stdouttrace.Option{stdouttrace.WithWriter(&buf)}
						if pretty {
							opts = append(opts, stdouttrace.WithPrettyPrint())
						}
						exporter, err := stdouttrace.New(opts...)
						require.NoError(t, err)
						require.NoError(t, exporter.ExportSpans(t.Context(), []sdktrace.ReadOnlySpan{span}))
						var got map[string]any
						require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
						assert.Equal(t, wantAttributes(location, shape.want), exportedAttributes(t, got, location))
						assert.Equal(t, before, shape.value.String(), "the input attribute must remain unchanged")
						if location == "event" {
							assert.Equal(
								t,
								float64(2),
								got["Events"].([]any)[0].(map[string]any)["DroppedAttributeCount"],
							)
						}
						if location == "link" {
							assert.Equal(
								t,
								float64(3),
								got["Links"].([]any)[0].(map[string]any)["DroppedAttributeCount"],
							)
						}
					})
				}
			}
		}
	}
}

func TestExporterNonFiniteBatch(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		var buf bytes.Buffer
		opts := []stdouttrace.Option{stdouttrace.WithWriter(&buf)}
		if pretty {
			opts = append(opts, stdouttrace.WithPrettyPrint())
		}
		exporter, err := stdouttrace.New(opts...)
		require.NoError(t, err)
		bad := spanWithAttribute("span", attribute.Float64Value(math.NaN()))
		require.NoError(
			t,
			exporter.ExportSpans(t.Context(), tracetest.SpanStubs{{Name: "before"}, bad, {Name: "after"}}.Snapshots()),
		)
		decoder := json.NewDecoder(&buf)
		for _, name := range []string{"before", "nonfinite", "after"} {
			var span map[string]any
			require.NoError(t, decoder.Decode(&span))
			assert.Equal(t, name, span["Name"])
		}
		var extra any
		require.ErrorIs(t, decoder.Decode(&extra), io.EOF)
	}
}

func TestExporterFiniteRepresentation(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		for _, location := range []string{"span", "resource", "event", "link", "scope"} {
			for _, shape := range []attributeShape{
				{"bool", attribute.BoolValue(true), true},
				{"boolSlice", attribute.BoolSliceValue([]bool{true, false}), []any{true, false}},
				{"integer", attribute.Int64Value(math.MaxInt64), json.Number("9223372036854775807")},
				{"intSlice", attribute.Int64SliceValue([]int64{1, -2}), []any{json.Number("1"), json.Number("-2")}},
				{"float", attribute.Float64Value(1.5), json.Number("1.5")},
				{"negativeZero", attribute.Float64Value(math.Copysign(0, -1)), json.Number("-0")},
				{"string", attribute.StringValue("\"\\\n<>&\x00\xff\u2028\u2029"), "\"\\\n<>&\x00\ufffd\u2028\u2029"},
				{"stringSlice", attribute.StringSliceValue([]string{"text", "\"\\\n<>&"}), []any{"text", "\"\\\n<>&"}},
				{"bytes", attribute.ByteSliceValue([]byte{0, 1, 2, 255}), "AAEC/w=="},
				{"emptyBytes", attribute.ByteSliceValue(nil), ""},
				{"jsonString", attribute.StringValue(`{"Type":"STRING","Value":[1]}`), `{"Type":"STRING","Value":[1]}`},
				{"typeKeys", attribute.MapValue(attribute.String("Type", "STRING"), attribute.String("Value", `{"nested":42}`)), map[string]any{"Type": "STRING", "Value": `{"nested":42}`}},
				{"array", attribute.SliceValue(), []any{}},
				{"map", attribute.MapValue(), map[string]any{}},
				{"nestedEmpty", attribute.SliceValue(attribute.Value{}), []any{nil}},
			} {
				t.Run(location+"/"+shape.name+"/"+strconv.FormatBool(pretty), func(t *testing.T) {
					var buf bytes.Buffer
					opts := []stdouttrace.Option{stdouttrace.WithWriter(&buf)}
					if pretty {
						opts = append(opts, stdouttrace.WithPrettyPrint())
					}
					exporter, err := stdouttrace.New(opts...)
					require.NoError(t, err)
					stub := spanWithAttribute(location, shape.value)
					require.NoError(t, exporter.ExportSpans(t.Context(), tracetest.SpanStubs{stub}.Snapshots()))
					decoder := json.NewDecoder(&buf)
					decoder.UseNumber()
					var got map[string]any
					require.NoError(t, decoder.Decode(&got))
					assert.Equal(t, wantAttributes(location, shape.want), exportedAttributes(t, got, location))
				})
			}
		}
	}
}

func TestExporterEmptyAttributes(t *testing.T) {
	for _, stub := range []tracetest.SpanStub{
		{},
		{Attributes: []attribute.KeyValue{}, Events: []sdktrace.Event{}, Links: []sdktrace.Link{}, Resource: resource.Empty()},
		{InstrumentationScope: instrumentation.Scope{Name: "scope", Attributes: attribute.NewSet()}},
		{Attributes: []attribute.KeyValue{{Key: "empty", Value: attribute.Value{}}}},
	} {
		var buf bytes.Buffer
		exporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf))
		require.NoError(t, err)
		require.NoError(t, exporter.ExportSpans(t.Context(), tracetest.SpanStubs{stub}.Snapshots()))
		var got struct {
			Attributes             json.RawMessage
			Events                 json.RawMessage
			Links                  json.RawMessage
			Resource               json.RawMessage
			InstrumentationScope   struct{ Attributes json.RawMessage }
			InstrumentationLibrary struct{ Attributes json.RawMessage }
		}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		assert.JSONEq(t, `{}`, string(got.InstrumentationScope.Attributes))
		assert.Equal(t, got.InstrumentationScope.Attributes, got.InstrumentationLibrary.Attributes)
		if stub.Resource == nil {
			assert.Equal(t, "null", string(got.Resource))
		} else {
			assert.JSONEq(t, `{}`, string(got.Resource))
		}
		if stub.Attributes == nil {
			assert.Equal(t, "null", string(got.Attributes))
		} else if len(stub.Attributes) == 0 {
			assert.JSONEq(t, `[]`, string(got.Attributes))
		} else {
			assert.JSONEq(t, `[{"empty":null}]`, string(got.Attributes))
		}
		if stub.Events == nil {
			assert.Equal(t, "null", string(got.Events))
		} else {
			assert.JSONEq(t, `[]`, string(got.Events))
		}
		if stub.Links == nil {
			assert.Equal(t, "null", string(got.Links))
		} else {
			assert.JSONEq(t, `[]`, string(got.Links))
		}
	}
}

func TestExporterAttributeLists(t *testing.T) {
	for _, pretty := range []bool{false, true} {
		for _, tc := range []struct {
			name  string
			attrs []attribute.KeyValue
			want  string
		}{
			{"nil", nil, `null`},
			{"empty", []attribute.KeyValue{}, `[]`},
			{
				"ordered",
				[]attribute.KeyValue{attribute.String("z", "first"), attribute.Int64("a", 2), attribute.String("z", "last"), attribute.String("Type", "user type"), attribute.String("Value", `{"nested":42}`)},
				`[{"z":"first"},{"a":2},{"z":"last"},{"Type":"user type"},{"Value":"{\"nested\":42}"}]`,
			},
		} {
			t.Run(tc.name+"/"+strconv.FormatBool(pretty), func(t *testing.T) {
				stub := tracetest.SpanStub{
					Attributes: tc.attrs,
					Events:     []sdktrace.Event{{Name: "event", Attributes: tc.attrs, Time: time.Unix(1, 0), DroppedAttributeCount: 2}},
					Links:      []sdktrace.Link{{Attributes: tc.attrs, DroppedAttributeCount: 3}},
				}
				span := stub.Snapshot()
				var buf bytes.Buffer
				opts := []stdouttrace.Option{stdouttrace.WithWriter(&buf)}
				if pretty {
					opts = append(opts, stdouttrace.WithPrettyPrint())
				}
				exporter, err := stdouttrace.New(opts...)
				require.NoError(t, err)
				require.NoError(t, exporter.ExportSpans(t.Context(), []sdktrace.ReadOnlySpan{span}))
				var got struct {
					Attributes json.RawMessage
					Events     []struct {
						Name                  string
						Attributes            json.RawMessage
						Time                  time.Time
						DroppedAttributeCount int
					}
					Links []struct {
						Attributes            json.RawMessage
						DroppedAttributeCount int
					}
				}
				require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
				assert.JSONEq(t, tc.want, string(got.Attributes))
				require.Len(t, got.Events, 1)
				assert.JSONEq(t, tc.want, string(got.Events[0].Attributes))
				assert.Equal(t, "event", got.Events[0].Name)
				assert.Equal(t, time.Unix(1, 0).UTC(), got.Events[0].Time)
				assert.Equal(t, 2, got.Events[0].DroppedAttributeCount)
				require.Len(t, got.Links, 1)
				assert.JSONEq(t, tc.want, string(got.Links[0].Attributes))
				assert.Equal(t, 3, got.Links[0].DroppedAttributeCount)
				assert.Equal(t, tc.attrs, span.Attributes())
				assert.Equal(t, tc.attrs, span.Events()[0].Attributes)
				assert.Equal(t, tc.attrs, span.Links()[0].Attributes)
			})
		}
	}
}

func TestExporterNonFiniteWithoutTimestamps(t *testing.T) {
	stub := spanWithAttribute("event", attribute.Float64Value(math.NaN()))
	stub.StartTime, stub.EndTime = time.Unix(1, 0), time.Unix(2, 0)
	span := stub.Snapshot()
	var buf bytes.Buffer
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf), stdouttrace.WithoutTimestamps())
	require.NoError(t, err)
	require.NoError(t, exporter.ExportSpans(t.Context(), []sdktrace.ReadOnlySpan{span}))
	var got struct {
		StartTime, EndTime time.Time
		Events             []struct{ Time time.Time }
	}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	assert.Zero(t, got.StartTime)
	assert.Zero(t, got.EndTime)
	require.Len(t, got.Events, 1)
	assert.Zero(t, got.Events[0].Time)
	assert.Equal(t, stub.StartTime, span.StartTime())
	assert.Equal(t, stub.EndTime, span.EndTime())
	assert.Equal(t, stub.Events[0].Time, span.Events()[0].Time)
	assert.True(t, math.IsNaN(span.Events()[0].Attributes[0].Value.AsFloat64()))
}

func TestExporterNonFiniteSDK(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.Float64("value", math.NaN())}
	var buf bytes.Buffer
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf))
	require.NoError(t, err)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.NewSchemaless(attrs...)),
	)
	tracer := provider.Tracer(
		"scope",
		trace.WithInstrumentationVersion("1.2.3"),
		trace.WithSchemaURL("https://example.com/schema"),
		trace.WithInstrumentationAttributes(attrs...),
	)
	_, span := tracer.Start(
		t.Context(),
		"SDK",
		trace.WithAttributes(attrs...),
		trace.WithLinks(trace.Link{Attributes: attrs}),
	)
	span.AddEvent("event", trace.WithAttributes(attrs...))
	span.End()
	require.NoError(t, provider.Shutdown(t.Context()))
	var got map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	for _, location := range []string{"span", "resource", "event", "link", "scope"} {
		assert.Equal(t, wantAttributes(location, "NaN"), exportedAttributes(t, got, location))
	}
}

type errorWriter struct {
	err   error
	calls int
}

func (w *errorWriter) Write([]byte) (int, error) { w.calls++; return 0, w.err }

func TestExporterErrorsContinue(t *testing.T) {
	var buf bytes.Buffer
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf))
	require.NoError(t, err)
	bad := spanWithAttribute("span", attribute.Float64Value(math.NaN()))
	bad.StartTime = time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	err = exporter.ExportSpans(t.Context(), tracetest.SpanStubs{bad, {Name: "after"}, bad}.Snapshots())
	require.ErrorContains(t, err, "failed to encode span 0")
	require.ErrorContains(t, err, "failed to encode span 2")
	assert.NotContains(t, err.Error(), "failed to encode span 1")
	var got map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
	assert.Equal(t, "after", got["Name"])

	writeErr := errors.New("write failed")
	writer := &errorWriter{err: writeErr}
	exporter, err = stdouttrace.New(stdouttrace.WithWriter(writer))
	require.NoError(t, err)
	err = exporter.ExportSpans(
		t.Context(),
		tracetest.SpanStubs{spanWithAttribute("span", attribute.Float64Value(math.NaN())), {Name: "after"}}.Snapshots(),
	)
	require.ErrorIs(t, err, writeErr)
	require.ErrorContains(t, err, "failed to encode span 0")
	require.ErrorContains(t, err, "failed to encode span 1")
	assert.Equal(t, 1, writer.calls, "encoding/json retains the first writer error")

	// A destination may itself return an encoding/json error. Even this must
	// keep the encoder's ordinary sticky writer-error behavior.
	unsupported := &json.UnsupportedValueError{Str: "NaN"}
	writer = &errorWriter{err: unsupported}
	exporter, err = stdouttrace.New(stdouttrace.WithWriter(writer))
	require.NoError(t, err)
	err = exporter.ExportSpans(t.Context(), tracetest.SpanStubs{{Name: "first"}, {Name: "second"}}.Snapshots())
	require.ErrorIs(t, err, unsupported)
	assert.Equal(t, 1, writer.calls)
}

func BenchmarkExporterAttributes(b *testing.B) {
	large := make([]float64, 1024)
	for i := range large {
		large[i] = float64(i)
	}
	for _, tc := range []struct {
		name  string
		value attribute.Value
	}{
		{"finite", attribute.Float64Value(1.5)},
		{"nonfinite", attribute.Float64Value(math.NaN())},
		{"finiteSlice", attribute.Float64SliceValue(large)},
		{"nested", attribute.MapValue(attribute.KeyValue{Key: "nested", Value: attribute.SliceValue(attribute.Float64Value(1.5))})},
	} {
		b.Run(tc.name, func(b *testing.B) {
			spans := tracetest.SpanStubs{spanWithAttribute("span", tc.value)}.Snapshots()
			exporter, err := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
			require.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = exporter.ExportSpans(b.Context(), spans)
			}
		})
	}
}
