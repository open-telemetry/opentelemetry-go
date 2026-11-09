// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func outputModes(t *testing.T, test func(*testing.T, []stdoutmetric.Option)) {
	t.Helper()
	t.Run("compact", func(t *testing.T) { test(t, nil) })
	t.Run("pretty", func(t *testing.T) { test(t, []stdoutmetric.Option{stdoutmetric.WithPrettyPrint()}) })
}

func exportJSON(t *testing.T, rm *metricdata.ResourceMetrics, opts ...stdoutmetric.Option) []byte {
	t.Helper()
	var out bytes.Buffer
	exp, err := stdoutmetric.New(append(opts, stdoutmetric.WithWriter(&out))...)
	require.NoError(t, err)
	require.NoError(t, exp.Export(t.Context(), rm))
	require.True(t, json.Valid(out.Bytes()), out.String())
	return out.Bytes()
}

func decodeJSON(t *testing.T, data []byte) any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	require.NoError(t, decoder.Decode(&value))
	return value
}

func TestNonFiniteAttributes(t *testing.T) {
	for _, special := range []struct {
		name  string
		value float64
	}{
		{"NaN", math.NaN()}, {"Infinity", math.Inf(1)}, {"-Infinity", math.Inf(-1)},
	} {
		t.Run(special.name, func(t *testing.T) {
			for _, shape := range []struct {
				name  string
				value attribute.Value
				want  string
			}{
				{"scalar", attribute.Float64Value(special.value), `"` + special.name + `"`},
				{"floats", attribute.Float64SliceValue([]float64{1, special.value, 2}), `[1,"` + special.name + `",2]`},
				{"array", attribute.SliceValue(attribute.StringValue("ok"), attribute.Float64Value(special.value)), `["ok","` + special.name + `"]`},
				{"nested", attribute.MapValue(attribute.KeyValue{Key: "array", Value: attribute.SliceValue(attribute.MapValue(attribute.Float64("nested", special.value)))}), `{"array":[{"nested":"` + special.name + `"}]}`},
			} {
				t.Run(shape.name, func(t *testing.T) {
					kv := attribute.KeyValue{Key: "problem", Value: shape.value}
					for _, location := range []string{"resource", "scope", "point", "exemplar"} {
						t.Run(location, func(t *testing.T) {
							outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
								rm := &metricdata.ResourceMetrics{ScopeMetrics: []metricdata.ScopeMetrics{
									{
										Scope: instrumentation.Scope{Name: "scope", Version: "v1", SchemaURL: "schema"},
										Metrics: []metricdata.Metrics{
											{
												Name: "valid",
												Data: metricdata.Gauge[int64]{
													DataPoints: []metricdata.DataPoint[int64]{{Value: 7}},
												},
											},
										},
									},
								}}
								dp := metricdata.DataPoint[int64]{Value: 1}
								switch location {
								case "resource":
									rm.Resource = resource.NewSchemaless(kv)
								case "scope":
									rm.ScopeMetrics[0].Scope.Attributes = attribute.NewSet(kv)
								case "point":
									dp.Attributes = attribute.NewSet(kv)
								case "exemplar":
									dp.Exemplars = []metricdata.Exemplar[int64]{
										{Value: 3, FilteredAttributes: []attribute.KeyValue{kv}},
									}
								}
								rm.ScopeMetrics[0].Metrics = append(
									rm.ScopeMetrics[0].Metrics,
									metricdata.Metrics{
										Name: "affected",
										Data: metricdata.Sum[int64]{DataPoints: []metricdata.DataPoint[int64]{dp}},
									},
								)
								before := kv.Value.String()
								doc := decodeJSON(t, exportJSON(t, rm, opts...)).(map[string]any)
								sm := doc["ScopeMetrics"].([]any)[0].(map[string]any)
								scope := sm["Scope"].(map[string]any)
								assert.Equal(t, "scope", scope["Name"])
								assert.Equal(t, "v1", scope["Version"])
								assert.Equal(t, "schema", scope["SchemaURL"])
								metrics := sm["Metrics"].([]any)
								require.Len(t, metrics, 2)
								assert.Equal(t, "valid", metrics[0].(map[string]any)["Name"])
								assert.Equal(t, "affected", metrics[1].(map[string]any)["Name"])
								data := metrics[1].(map[string]any)["Data"].(map[string]any)
								point := data["DataPoints"].([]any)[0].(map[string]any)
								var gotAttributes any
								switch location {
								case "resource":
									gotAttributes = doc["Resource"]
								case "scope":
									gotAttributes = scope["Attributes"]
								case "point":
									gotAttributes = point["Attributes"]
								case "exemplar":
									gotAttributes = point["Exemplars"].([]any)[0].(map[string]any)["FilteredAttributes"]
								}
								wantJSON := `{"problem":` + shape.want + `}`
								if location == "exemplar" {
									wantJSON = `[` + wantJSON + `]`
								}
								wantAttributes := decodeJSON(t, []byte(wantJSON))
								assert.Equal(t, wantAttributes, gotAttributes)
								var after attribute.Value
								switch location {
								case "resource":
									after, _ = rm.Resource.Set().Value(kv.Key)
								case "scope":
									after, _ = rm.ScopeMetrics[0].Scope.Attributes.Value(kv.Key)
								case "point":
									after, _ = rm.ScopeMetrics[0].Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Attributes.Value(
										kv.Key,
									)
								case "exemplar":
									after = rm.ScopeMetrics[0].Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Exemplars[0].FilteredAttributes[0].Value
								}
								assert.Equal(t, before, after.String())
							})
						})
					}
				})
			}
		})
	}
}

// numericMetrics covers every float field, both number types, and pointer
// aggregations. Distinct finite controls make accidental shape changes visible.
func numericMetrics(value float64) *metricdata.ResourceMetrics {
	attrs := attribute.NewSet(attribute.Int64("large", math.MaxInt64), attribute.String("label", "<&>"))
	tm := time.Unix(1, 2).UTC()
	ex := []metricdata.Exemplar[float64]{
		{
			Value:              value,
			Time:               tm,
			FilteredAttributes: []attribute.KeyValue{attribute.String("ok", "yes")},
			SpanID:             []byte{1},
			TraceID:            []byte{2},
		},
	}
	dp := []metricdata.DataPoint[float64]{{Attributes: attrs, StartTime: tm, Time: tm, Value: value, Exemplars: ex}}
	hdp := []metricdata.HistogramDataPoint[float64]{
		{
			Attributes:   attrs,
			StartTime:    tm,
			Time:         tm,
			Count:        1,
			Bounds:       []float64{value},
			BucketCounts: []uint64{1, 0},
			Min:          metricdata.NewExtrema(value),
			Sum:          value,
			Exemplars:    ex,
		},
	}
	edp := []metricdata.ExponentialHistogramDataPoint[float64]{
		{
			Attributes:    attrs,
			StartTime:     tm,
			Time:          tm,
			Count:         1,
			Min:           metricdata.NewExtrema(value),
			Max:           metricdata.NewExtrema(value),
			Sum:           value,
			ZeroThreshold: value,
			Exemplars:     ex,
		},
	}
	aggs := []metricdata.Aggregation{
		metricdata.Gauge[float64]{DataPoints: dp},
		metricdata.Sum[float64]{DataPoints: dp, IsMonotonic: true, Temporality: metricdata.CumulativeTemporality},
		metricdata.Histogram[float64]{DataPoints: hdp, Temporality: metricdata.DeltaTemporality},
		metricdata.ExponentialHistogram[float64]{DataPoints: edp},
		metricdata.Gauge[int64]{
			DataPoints: []metricdata.DataPoint[int64]{
				{
					Attributes: attrs,
					Value:      math.MaxInt64,
					Exemplars:  []metricdata.Exemplar[int64]{{Value: math.MaxInt64}},
				},
			},
		},
		metricdata.Sum[int64]{DataPoints: []metricdata.DataPoint[int64]{{Attributes: attrs, Value: 1}}},
		metricdata.Histogram[int64]{
			DataPoints: []metricdata.HistogramDataPoint[int64]{
				{
					Attributes: attrs,
					Bounds:     []float64{value},
					Sum:        math.MaxInt64,
					Min:        metricdata.NewExtrema(int64(math.MinInt64)),
					Max:        metricdata.NewExtrema(int64(math.MaxInt64)),
				},
			},
		},
		metricdata.ExponentialHistogram[int64]{
			DataPoints: []metricdata.ExponentialHistogramDataPoint[int64]{
				{Attributes: attrs, ZeroThreshold: value, Sum: math.MaxInt64},
			},
		},
		metricdata.Summary{
			DataPoints: []metricdata.SummaryDataPoint{
				{
					Attributes:     attrs,
					Sum:            value,
					QuantileValues: []metricdata.QuantileValue{{Quantile: value, Value: value}},
				},
			},
		},
	}
	rm := &metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("service", "test")),
		ScopeMetrics: []metricdata.ScopeMetrics{{Scope: instrumentation.Scope{Name: "numeric", Attributes: attrs}}},
	}
	for _, agg := range aggs {
		rm.ScopeMetrics[0].Metrics = append(rm.ScopeMetrics[0].Metrics, metricdata.Metrics{Name: "value", Data: agg})
		ptr := reflect.New(reflect.TypeOf(agg))
		ptr.Elem().Set(reflect.ValueOf(agg))
		rm.ScopeMetrics[0].Metrics = append(
			rm.ScopeMetrics[0].Metrics,
			metricdata.Metrics{Name: "pointer", Data: ptr.Interface().(metricdata.Aggregation)},
		)
	}
	return rm
}

func replaceFloat(v any, old float64, replacement string) any {
	switch value := v.(type) {
	case float64:
		if value == old {
			return replacement
		}
	case []any:
		for i := range value {
			value[i] = replaceFloat(value[i], old, replacement)
		}
	case map[string]any:
		for k := range value {
			value[k] = replaceFloat(value[k], old, replacement)
		}
	}
	return v
}

func TestNonFiniteNumericFields(t *testing.T) {
	const sentinel = 1.23456789
	// Keep the original metric-data field representation as the control, with
	// independently specified expectations for the changed attribute fields.
	control, err := json.Marshal(numericMetrics(sentinel))
	require.NoError(t, err)
	var finite map[string]any
	require.NoError(t, json.Unmarshal(control, &finite))
	finite["Resource"] = map[string]any{"service": "test"}
	scope := finite["ScopeMetrics"].([]any)[0].(map[string]any)
	attrs := map[string]any{"large": float64(math.MaxInt64), "label": "<&>"}
	scope["Scope"].(map[string]any)["Attributes"] = attrs
	for _, m := range scope["Metrics"].([]any) {
		for _, p := range m.(map[string]any)["Data"].(map[string]any)["DataPoints"].([]any) {
			point := p.(map[string]any)
			point["Attributes"] = attrs
			if exemplars, ok := point["Exemplars"].([]any); ok {
				for _, e := range exemplars {
					exemplar := e.(map[string]any)
					if exemplar["FilteredAttributes"] != nil {
						exemplar["FilteredAttributes"] = []any{map[string]any{"ok": "yes"}}
					}
				}
			}
		}
	}
	control, err = json.Marshal(finite)
	require.NoError(t, err)
	for name, value := range map[string]float64{"NaN": math.NaN(), "Infinity": math.Inf(1), "-Infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				rm := numericMetrics(value)
				var want, got any
				require.NoError(t, json.Unmarshal(control, &want))
				want = replaceFloat(want, sentinel, name)
				out := exportJSON(t, rm, opts...)
				require.NoError(t, json.Unmarshal(out, &got))
				assert.Equal(t, want, got)
				assert.Contains(t, string(out), "9223372036854775807")
				assert.Contains(t, string(out), "-9223372036854775808")
				assert.JSONEq(t, string(out), string(exportJSON(t, rm, opts...)))
				dp := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[float64]).DataPoints[0]
				assert.Equal(t, math.Float64bits(value), math.Float64bits(dp.Value))
				hdp := rm.ScopeMetrics[0].Metrics[4].Data.(metricdata.Histogram[float64]).DataPoints[0]
				assert.Equal(t, math.Float64bits(value), math.Float64bits(hdp.Bounds[0]))
				minimum, defined := hdp.Min.Value()
				assert.True(t, defined)
				assert.Equal(t, math.Float64bits(value), math.Float64bits(minimum))
			})
		})
	}
}

func TestJSONAttributeRepresentation(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.Bool("bool", true),
		attribute.String("A", "upper"),
		attribute.String("a", "lower"),
		attribute.String("json-string", `{"Type":"STRING","Value":"user"}`),
		{Key: "envelope-keys", Value: attribute.MapValue(attribute.String("Type", "user"), attribute.String("Value", "payload"))},
		attribute.BoolSlice("bools", []bool{true, false}),
		attribute.Int64("int", math.MinInt64),
		attribute.Int64Slice("ints", []int64{math.MaxInt64}),
		attribute.Float64("float", 1.5),
		attribute.Float64Slice("floats", []float64{1.5, math.Copysign(0, -1)}),
		attribute.String("string", "<&>\\\"\n\xff"),
		attribute.StringSlice("strings", []string{"<&>", "quote\""}),
		{Key: "bytes", Value: attribute.ByteSliceValue([]byte{0, 1, 255})},
		{Key: "empty", Value: attribute.Value{}},
		{Key: "map", Value: attribute.MapValue(attribute.String("nested", "value"))},
		{Key: "array", Value: attribute.SliceValue(attribute.StringValue("x"), attribute.Int64Value(3), attribute.Value{})},
		{Key: "empty-map", Value: attribute.MapValue()},
		{Key: "empty-array", Value: attribute.SliceValue()},
	}
	const wantAttrs = `{
		"bool":true,"bools":[true,false],"int":-9223372036854775808,"ints":[9223372036854775807],
		"A":"upper","a":"lower","json-string":"{\"Type\":\"STRING\",\"Value\":\"user\"}","envelope-keys":{"Type":"user","Value":"payload"},
		"float":1.5,"floats":[1.5,-0],"string":"<&>\\\"\n\ufffd","strings":["<&>","quote\""],
		"bytes":"AAH/","empty":null,"map":{"nested":"value"},"array":["x",3,null],"empty-map":{},"empty-array":[]
	}`
	outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
		rm := &metricdata.ResourceMetrics{
			Resource: resource.NewSchemaless(attrs...),
			ScopeMetrics: []metricdata.ScopeMetrics{{
				Scope: instrumentation.Scope{Name: "scope", Attributes: attribute.NewSet(attrs...)},
				Metrics: []metricdata.Metrics{{Name: "finite", Data: metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{
					Attributes: attribute.NewSet(attrs...),
					Value:      math.MaxInt64,
					Exemplars: []metricdata.Exemplar[int64]{
						{FilteredAttributes: []attribute.KeyValue{attribute.String("z", "first"), attribute.String("a", "middle"), attribute.String("z", "last")}},
						{FilteredAttributes: []attribute.KeyValue{}},
						{},
					},
				}}}}},
			}},
		}
		out := exportJSON(t, rm, opts...)
		doc := decodeJSON(t, out).(map[string]any)
		want := decodeJSON(t, []byte(wantAttrs))
		assert.Equal(t, want, doc["Resource"])
		scope := doc["ScopeMetrics"].([]any)[0].(map[string]any)
		assert.Equal(t, want, scope["Scope"].(map[string]any)["Attributes"])
		point := scope["Metrics"].([]any)[0].(map[string]any)["Data"].(map[string]any)["DataPoints"].([]any)[0].(map[string]any)
		assert.Equal(t, want, point["Attributes"])
		assert.Equal(t, json.Number("9223372036854775807"), point["Value"])
		exemplars := point["Exemplars"].([]any)
		assert.Equal(t, decodeJSON(t, []byte(`[{"z":"first"},{"a":"middle"},{"z":"last"}]`)), exemplars[0].(map[string]any)["FilteredAttributes"])
		assert.Equal(t, []any{}, exemplars[1].(map[string]any)["FilteredAttributes"])
		assert.Nil(t, exemplars[2].(map[string]any)["FilteredAttributes"])
		assert.NotContains(t, string(out), `"FLOAT64"`)
	})
}

func TestJSONAggregationShapes(t *testing.T) {
	for _, tc := range []struct {
		name string
		data metricdata.Aggregation
		want string
	}{
		{"nil", nil, `null`},
		{"typed-nil", (*metricdata.Gauge[float64])(nil), `null`},
		{"nil-points", metricdata.Gauge[int64]{}, `{"DataPoints":null}`},
		{"empty-points", metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{}}, `{"DataPoints":[]}`},
		{"empty-attributes", metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{}, {Attributes: attribute.NewSet(), Exemplars: []metricdata.Exemplar[int64]{}}}}, `{"DataPoints":[{"Attributes":{},"StartTime":"0001-01-01T00:00:00Z","Time":"0001-01-01T00:00:00Z","Value":0},{"Attributes":{},"StartTime":"0001-01-01T00:00:00Z","Time":"0001-01-01T00:00:00Z","Value":0}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				rm := &metricdata.ResourceMetrics{ScopeMetrics: []metricdata.ScopeMetrics{{Metrics: []metricdata.Metrics{{Data: tc.data}}}}}
				doc := decodeJSON(t, exportJSON(t, rm, opts...)).(map[string]any)
				scope := doc["ScopeMetrics"].([]any)[0].(map[string]any)
				got := scope["Metrics"].([]any)[0].(map[string]any)["Data"]
				assert.Equal(t, decodeJSON(t, []byte(tc.want)), got)
			})
		})
	}
}

func TestJSONEmptyCollections(t *testing.T) {
	for _, tc := range []struct {
		name    string
		metrics *metricdata.ResourceMetrics
		want    string
	}{
		{"nil", nil, `null`},
		{"zero", &metricdata.ResourceMetrics{}, `{"Resource":null,"ScopeMetrics":null}`},
		{"empty-scopes", &metricdata.ResourceMetrics{ScopeMetrics: []metricdata.ScopeMetrics{}}, `{"Resource":null,"ScopeMetrics":[]}`},
		{"zero-attributes", &metricdata.ResourceMetrics{Resource: resource.Empty(), ScopeMetrics: []metricdata.ScopeMetrics{{}}}, `{"Resource":{},"ScopeMetrics":[{"Scope":{"Name":"","Version":"","SchemaURL":"","Attributes":{}},"Metrics":null}]}`},
		{"empty-attributes", &metricdata.ResourceMetrics{Resource: resource.NewSchemaless(), ScopeMetrics: []metricdata.ScopeMetrics{{Scope: instrumentation.Scope{Attributes: attribute.NewSet()}, Metrics: []metricdata.Metrics{}}}}, `{"Resource":{},"ScopeMetrics":[{"Scope":{"Name":"","Version":"","SchemaURL":"","Attributes":{}},"Metrics":[]}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				assert.JSONEq(t, tc.want, string(exportJSON(t, tc.metrics, opts...)))
			})
		})
	}
}

func TestNonFinitePublicSDK(t *testing.T) {
	for name, value := range map[string]float64{"NaN": math.NaN(), "Infinity": math.Inf(1), "-Infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				reader := metric.NewManualReader()
				provider := metric.NewMeterProvider(
					metric.WithReader(reader),
					metric.WithResource(resource.NewSchemaless(attribute.Float64("resource", value))),
				)
				t.Cleanup(func() { require.NoError(t, provider.Shutdown(t.Context())) })
				meter := provider.Meter(
					"public",
					api.WithInstrumentationVersion("v1"),
					api.WithSchemaURL("schema"),
					api.WithInstrumentationAttributes(attribute.Float64("scope", value)),
				)
				counter, err := meter.Int64Counter("valid")
				require.NoError(t, err)
				counter.Add(t.Context(), 5)
				counter.Add(t.Context(), 7, api.WithAttributes(attribute.Float64("point", value)))
				gauge, err := meter.Float64Gauge("gauge")
				require.NoError(t, err)
				gauge.Record(t.Context(), value)
				floatCounter, err := meter.Float64Counter("counter")
				require.NoError(t, err)
				floatCounter.Add(t.Context(), value)
				updown, err := meter.Float64UpDownCounter("updown")
				require.NoError(t, err)
				updown.Add(t.Context(), value)
				hist, err := meter.Float64Histogram("histogram")
				require.NoError(t, err)
				hist.Record(t.Context(), value)
				var rm metricdata.ResourceMetrics
				require.NoError(t, reader.Collect(t.Context(), &rm))
				out := string(exportJSON(t, &rm, opts...))
				for _, metricName := range []string{"valid", "gauge", "counter", "updown", "histogram"} {
					assert.Contains(t, out, `"`+metricName+`"`)
				}
				assert.Contains(t, out, `"`+name+`"`)
			})
		})
	}
}

type inspectEncoder struct{ encode func(any) error }

func (e inspectEncoder) Encode(v any) error { return e.encode(v) }

func TestNonFiniteCustomEncoder(t *testing.T) {
	rm := numericMetrics(math.NaN())
	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(inspectEncoder{func(got any) error {
		assert.Same(t, rm, got)
		return errEnc
	}}), stdoutmetric.WithPrettyPrint())
	require.NoError(t, err)
	assert.ErrorIs(t, exp.Export(t.Context(), rm), errEnc)
	var out bytes.Buffer
	exp, err = stdoutmetric.New(stdoutmetric.WithEncoder(json.NewEncoder(&out)), stdoutmetric.WithPrettyPrint())
	require.NoError(t, err)
	assert.Error(t, exp.Export(t.Context(), rm))
	assert.Empty(t, out.String())
}

type errorWriter struct {
	err   error
	calls int
}

func (w *errorWriter) Write([]byte) (int, error) {
	w.calls++
	return 0, w.err
}

func TestJSONErrors(t *testing.T) {
	for _, writeErr := range []error{errors.New("write failed"), &json.UnsupportedValueError{Value: reflect.ValueOf(math.NaN()), Str: "NaN"}} {
		for _, value := range []float64{1, math.NaN()} {
			w := &errorWriter{err: writeErr}
			exp, err := stdoutmetric.New(stdoutmetric.WithWriter(w))
			require.NoError(t, err)
			for range 2 {
				assert.ErrorIs(t, exp.Export(t.Context(), numericMetrics(value)), writeErr)
			}
			assert.Equal(t, 1, w.calls, "a writer error must not trigger a second write")
		}
	}
	rm := numericMetrics(math.NaN())
	badTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	dp := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[float64]).DataPoints
	dp[0].Time = badTime
	exp, err := stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
	require.NoError(t, err)
	assert.ErrorContains(t, exp.Export(t.Context(), rm), "year outside of range")
	// An invalid timestamp late in the collection must still be reported.
	rm = numericMetrics(math.NaN())
	last := len(rm.ScopeMetrics[0].Metrics) - 2
	rm.ScopeMetrics[0].Metrics[last].Data.(metricdata.Summary).DataPoints[0].Time = badTime
	assert.ErrorContains(t, exp.Export(t.Context(), rm), "year outside of range")
}

func BenchmarkJSONExport(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value float64
	}{
		{"Finite", 1.25}, {"NonFinite", math.NaN()},
	} {
		b.Run(tc.name, func(b *testing.B) {
			rm := numericMetrics(tc.value)
			exp, err := stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
			require.NoError(b, err)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = exp.Export(b.Context(), rm)
			}
		})
	}
}
