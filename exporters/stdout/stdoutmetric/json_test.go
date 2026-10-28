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
				{"scalar", attribute.Float64Value(special.value), `{"Type":"FLOAT64","Value":"` + special.name + `"}`},
				{"floats", attribute.Float64SliceValue([]float64{1, special.value, 2}), `{"Type":"FLOAT64SLICE","Value":[1,"` + special.name + `",2]}`},
				{"array", attribute.SliceValue(attribute.StringValue("ok"), attribute.Float64Value(special.value)), `{"Type":"SLICE","Value":[{"Type":"STRING","Value":"ok"},{"Type":"FLOAT64","Value":"` + special.name + `"}]}`},
				{"nested", attribute.MapValue(attribute.KeyValue{Key: "array", Value: attribute.SliceValue(attribute.MapValue(attribute.Float64("nested", special.value)))}), `{"Type":"MAP","Value":[{"Key":"array","Value":{"Type":"SLICE","Value":[{"Type":"MAP","Value":[{"Key":"nested","Value":{"Type":"FLOAT64","Value":"` + special.name + `"}}]}]}}]}`},
			} {
				t.Run(shape.name, func(t *testing.T) {
					kv := attribute.KeyValue{Key: "problem", Value: shape.value}
					for _, location := range []string{"resource", "scope", "point", "exemplar"} {
						t.Run(location, func(t *testing.T) {
							outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
								rm := &metricdata.ResourceMetrics{ScopeMetrics: []metricdata.ScopeMetrics{{
									Scope:   instrumentation.Scope{Name: "scope", Version: "v1", SchemaURL: "schema"},
									Metrics: []metricdata.Metrics{{Name: "valid", Data: metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{Value: 7}}}}},
								}}}
								dp := metricdata.DataPoint[int64]{Value: 1}
								switch location {
								case "resource":
									rm.Resource = resource.NewSchemaless(kv)
								case "scope":
									rm.ScopeMetrics[0].Scope.Attributes = attribute.NewSet(kv)
								case "point":
									dp.Attributes = attribute.NewSet(kv)
								case "exemplar":
									dp.Exemplars = []metricdata.Exemplar[int64]{{Value: 3, FilteredAttributes: []attribute.KeyValue{kv}}}
								}
								rm.ScopeMetrics[0].Metrics = append(rm.ScopeMetrics[0].Metrics, metricdata.Metrics{Name: "affected", Data: metricdata.Sum[int64]{DataPoints: []metricdata.DataPoint[int64]{dp}}})
								before := kv.Value.String()
								var compact bytes.Buffer
								require.NoError(t, json.Compact(&compact, exportJSON(t, rm, opts...)))
								assert.Contains(t, compact.String(), `"Value":`+shape.want)
								assert.Contains(t, compact.String(), `"Name":"valid"`)
								assert.Contains(t, compact.String(), `"Name":"affected"`)
								assert.Contains(t, compact.String(), `"Name":"scope","Version":"v1","SchemaURL":"schema"`)
								var after attribute.Value
								switch location {
								case "resource":
									after, _ = rm.Resource.Set().Value(kv.Key)
								case "scope":
									after, _ = rm.ScopeMetrics[0].Scope.Attributes.Value(kv.Key)
								case "point":
									after, _ = rm.ScopeMetrics[0].Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Attributes.Value(kv.Key)
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
	ex := []metricdata.Exemplar[float64]{{Value: value, Time: tm, FilteredAttributes: []attribute.KeyValue{attribute.String("ok", "yes")}, SpanID: []byte{1}, TraceID: []byte{2}}}
	dp := []metricdata.DataPoint[float64]{{Attributes: attrs, StartTime: tm, Time: tm, Value: value, Exemplars: ex}}
	hdp := []metricdata.HistogramDataPoint[float64]{{Attributes: attrs, StartTime: tm, Time: tm, Count: 1, Bounds: []float64{value}, BucketCounts: []uint64{1, 0}, Min: metricdata.NewExtrema(value), Sum: value, Exemplars: ex}}
	edp := []metricdata.ExponentialHistogramDataPoint[float64]{{Attributes: attrs, StartTime: tm, Time: tm, Count: 1, Min: metricdata.NewExtrema(value), Max: metricdata.NewExtrema(value), Sum: value, ZeroThreshold: value, Exemplars: ex}}
	aggs := []metricdata.Aggregation{
		metricdata.Gauge[float64]{DataPoints: dp}, metricdata.Sum[float64]{DataPoints: dp, IsMonotonic: true, Temporality: metricdata.CumulativeTemporality},
		metricdata.Histogram[float64]{DataPoints: hdp, Temporality: metricdata.DeltaTemporality}, metricdata.ExponentialHistogram[float64]{DataPoints: edp},
		metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{Attributes: attrs, Value: math.MaxInt64, Exemplars: []metricdata.Exemplar[int64]{{Value: math.MaxInt64}}}}},
		metricdata.Sum[int64]{DataPoints: []metricdata.DataPoint[int64]{{Attributes: attrs, Value: 1}}},
		metricdata.Histogram[int64]{DataPoints: []metricdata.HistogramDataPoint[int64]{{Attributes: attrs, Bounds: []float64{value}, Sum: math.MaxInt64, Min: metricdata.NewExtrema(int64(math.MinInt64)), Max: metricdata.NewExtrema(int64(math.MaxInt64))}}},
		metricdata.ExponentialHistogram[int64]{DataPoints: []metricdata.ExponentialHistogramDataPoint[int64]{{Attributes: attrs, ZeroThreshold: value, Sum: math.MaxInt64}}},
		metricdata.Summary{DataPoints: []metricdata.SummaryDataPoint{{Attributes: attrs, Sum: value, QuantileValues: []metricdata.QuantileValue{{Quantile: value, Value: value}}}}},
	}
	rm := &metricdata.ResourceMetrics{Resource: resource.NewSchemaless(attribute.String("service", "test")), ScopeMetrics: []metricdata.ScopeMetrics{{Scope: instrumentation.Scope{Name: "numeric", Attributes: attrs}}}}
	for _, agg := range aggs {
		rm.ScopeMetrics[0].Metrics = append(rm.ScopeMetrics[0].Metrics, metricdata.Metrics{Name: "value", Data: agg})
		ptr := reflect.New(reflect.TypeOf(agg))
		ptr.Elem().Set(reflect.ValueOf(agg))
		rm.ScopeMetrics[0].Metrics = append(rm.ScopeMetrics[0].Metrics, metricdata.Metrics{Name: "pointer", Data: ptr.Interface().(metricdata.Aggregation)})
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
	control, err := json.Marshal(numericMetrics(sentinel))
	require.NoError(t, err)
	for name, value := range map[string]float64{"NaN": math.NaN(), "Infinity": math.Inf(1), "-Infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				rm := numericMetrics(value)
				var want, got any
				require.NoError(t, json.Unmarshal(control, &want))
				want = replaceFloat(want, sentinel, name)
				require.NoError(t, json.Unmarshal(exportJSON(t, rm, opts...), &got))
				assert.Equal(t, want, got)
				// Exporting twice also verifies that caller-owned slices and extrema
				// have not been rewritten to produce the JSON representation.
				assert.JSONEq(t, string(exportJSON(t, rm, opts...)), string(exportJSON(t, rm, opts...)))
				dp := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[float64]).DataPoints[0]
				assert.Equal(t, math.Float64bits(value), math.Float64bits(dp.Value))
			})
		})
	}
}

func TestJSONFiniteCompatibility(t *testing.T) {
	for _, rm := range []*metricdata.ResourceMetrics{nil, {}, {ScopeMetrics: []metricdata.ScopeMetrics{}}, numericMetrics(1.25)} {
		outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
			var want bytes.Buffer
			enc := json.NewEncoder(&want)
			if len(opts) > 0 {
				enc.SetIndent("", "\t")
			}
			require.NoError(t, enc.Encode(rm))
			assert.Equal(t, want.String(), string(exportJSON(t, rm, opts...)))
		})
	}
}

func TestNonFinitePublicSDK(t *testing.T) {
	for name, value := range map[string]float64{"NaN": math.NaN(), "Infinity": math.Inf(1), "-Infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			outputModes(t, func(t *testing.T, opts []stdoutmetric.Option) {
				reader := metric.NewManualReader()
				provider := metric.NewMeterProvider(metric.WithReader(reader), metric.WithResource(resource.NewSchemaless(attribute.Float64("resource", value))))
				t.Cleanup(func() { require.NoError(t, provider.Shutdown(t.Context())) })
				meter := provider.Meter("public", api.WithInstrumentationVersion("v1"), api.WithSchemaURL("schema"), api.WithInstrumentationAttributes(attribute.Float64("scope", value)))
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
		w := &errorWriter{err: writeErr}
		exp, err := stdoutmetric.New(stdoutmetric.WithWriter(w))
		require.NoError(t, err)
		for range 2 {
			assert.ErrorIs(t, exp.Export(t.Context(), numericMetrics(1)), writeErr)
		}
		assert.Equal(t, 1, w.calls, "a writer error must not trigger a second write")
	}
	rm := numericMetrics(math.NaN())
	badTime := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	dp := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[float64]).DataPoints
	dp[0].Time = badTime
	exp, err := stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
	require.NoError(t, err)
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
