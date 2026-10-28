// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric

import (
	"encoding/json"
	"errors"
	"io"
	"math"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// jsonEncoder preserves the usual encoding path for finite data. A marshal
// failure writes no bytes, so non-finite floats can be retried with JSON-safe
// representations without losing the other metrics in the collection.
type jsonEncoder struct{ *json.Encoder }

func newJSONEncoder(w io.Writer) *jsonEncoder {
	return &jsonEncoder{json.NewEncoder(jsonWriter{w})}
}

func (e *jsonEncoder) Encode(v any) error {
	err := e.Encoder.Encode(v)
	if err == nil {
		return nil
	}
	var writeErr *jsonWriteError
	if errors.As(err, &writeErr) {
		return writeErr.err
	}
	var unsupported *json.UnsupportedValueError
	if !errors.As(err, &unsupported) {
		return err
	}
	// Some encoding/json paths leave Value unset. Str also distinguishes
	// non-finite values from other unsupported values, such as cycles.
	switch unsupported.Str {
	case "NaN", "+Inf", "-Inf":
	default:
		return err
	}
	rm, ok := v.(*metricdata.ResourceMetrics)
	if !ok || rm == nil {
		return err
	}
	err = e.Encoder.Encode(jsonResourceMetrics(rm))
	if errors.As(err, &writeErr) {
		return writeErr.err
	}
	return err
}

// Identify writer errors separately from marshal errors, including errors
// cached by json.Encoder after a failed write.
type jsonWriter struct{ io.Writer }

func (w jsonWriter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err != nil {
		return n, &jsonWriteError{err}
	}
	return n, nil
}

type jsonWriteError struct{ err error }

func (e *jsonWriteError) Error() string { return e.err.Error() }

func jsonFloat(f float64) any {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	default:
		return f
	}
}

func jsonNumber[N int64 | float64](n N) any {
	if f, ok := any(n).(float64); ok {
		return jsonFloat(f)
	}
	return n
}

func jsonSlice[T, U any](in []T, convert func(T) U) []U {
	if in == nil {
		return nil
	}
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = convert(v)
	}
	return out
}

type jsonAttribute struct {
	Key   attribute.Key
	Value any
}

func jsonAttributes(kvs []attribute.KeyValue) []jsonAttribute {
	return jsonSlice(kvs, func(kv attribute.KeyValue) jsonAttribute {
		return jsonAttribute{Key: kv.Key, Value: jsonAttributeValue(kv.Value)}
	})
}

func jsonAttributeSet(set attribute.Set) any {
	if set.Len() == 0 {
		// An initialized empty Set encodes as [], while the zero Set is null.
		return &set
	}
	return jsonAttributes(set.ToSlice())
}

func jsonAttributeValue(v attribute.Value) any {
	var value any
	switch v.Type() {
	case attribute.FLOAT64:
		value = jsonFloat(v.AsFloat64())
	case attribute.FLOAT64SLICE:
		value = jsonSlice(v.AsFloat64Slice(), jsonFloat)
	case attribute.SLICE:
		value = jsonSlice(v.AsSlice(), jsonAttributeValue)
	case attribute.MAP:
		value = jsonAttributes(v.AsMap())
	default:
		value = v.AsInterface()
	}
	return struct {
		Type  string
		Value any
	}{v.Type().String(), value}
}

func jsonResourceMetrics(rm *metricdata.ResourceMetrics) any {
	var res any
	if rm.Resource != nil {
		res = jsonAttributeSet(*rm.Resource.Set())
	}
	return struct {
		Resource     any
		ScopeMetrics []any
	}{res, jsonSlice(rm.ScopeMetrics, func(sm metricdata.ScopeMetrics) any {
		return struct {
			Scope   any
			Metrics []any
		}{
			Scope: struct {
				instrumentation.Scope
				Attributes any
			}{sm.Scope, jsonAttributeSet(sm.Scope.Attributes)},
			Metrics: jsonSlice(sm.Metrics, func(m metricdata.Metrics) any {
				return struct {
					metricdata.Metrics
					Data any
				}{m, jsonAggregation(m.Data)}
			}),
		}
	})}
}

func jsonAggregation(agg metricdata.Aggregation) any {
	switch a := agg.(type) {
	case metricdata.Gauge[int64]:
		return jsonGauge(a)
	case metricdata.Gauge[float64]:
		return jsonGauge(a)
	case metricdata.Sum[int64]:
		return jsonSum(a)
	case metricdata.Sum[float64]:
		return jsonSum(a)
	case metricdata.Histogram[int64]:
		return jsonHistogram(a)
	case metricdata.Histogram[float64]:
		return jsonHistogram(a)
	case metricdata.ExponentialHistogram[int64]:
		return jsonExponentialHistogram(a)
	case metricdata.ExponentialHistogram[float64]:
		return jsonExponentialHistogram(a)
	case metricdata.Summary:
		return jsonSummary(a)
	case *metricdata.Gauge[int64]:
		if a != nil {
			return jsonGauge(*a)
		}
	case *metricdata.Gauge[float64]:
		if a != nil {
			return jsonGauge(*a)
		}
	case *metricdata.Sum[int64]:
		if a != nil {
			return jsonSum(*a)
		}
	case *metricdata.Sum[float64]:
		if a != nil {
			return jsonSum(*a)
		}
	case *metricdata.Histogram[int64]:
		if a != nil {
			return jsonHistogram(*a)
		}
	case *metricdata.Histogram[float64]:
		if a != nil {
			return jsonHistogram(*a)
		}
	case *metricdata.ExponentialHistogram[int64]:
		if a != nil {
			return jsonExponentialHistogram(*a)
		}
	case *metricdata.ExponentialHistogram[float64]:
		if a != nil {
			return jsonExponentialHistogram(*a)
		}
	case *metricdata.Summary:
		if a != nil {
			return jsonSummary(*a)
		}
	}
	return agg
}

func jsonGauge[N int64 | float64](a metricdata.Gauge[N]) any {
	return struct{ DataPoints []any }{jsonSlice(a.DataPoints, jsonDataPoint[N])}
}

func jsonSum[N int64 | float64](a metricdata.Sum[N]) any {
	return struct {
		metricdata.Sum[N]
		DataPoints []any
	}{a, jsonSlice(a.DataPoints, jsonDataPoint[N])}
}

func jsonDataPoint[N int64 | float64](dp metricdata.DataPoint[N]) any {
	return struct {
		metricdata.DataPoint[N]
		Attributes any
		Value      any
		Exemplars  []any `json:",omitempty"`
	}{dp, jsonAttributeSet(dp.Attributes), jsonNumber(dp.Value), jsonSlice(dp.Exemplars, jsonExemplar[N])}
}

func jsonExemplar[N int64 | float64](ex metricdata.Exemplar[N]) any {
	return struct {
		metricdata.Exemplar[N]
		FilteredAttributes []jsonAttribute
		Value              any
	}{ex, jsonAttributes(ex.FilteredAttributes), jsonNumber(ex.Value)}
}

func jsonExtrema[N int64 | float64](e metricdata.Extrema[N]) any {
	if value, ok := e.Value(); ok {
		return jsonNumber(value)
	}
	return nil
}

func jsonHistogram[N int64 | float64](a metricdata.Histogram[N]) any {
	return struct {
		metricdata.Histogram[N]
		DataPoints []any
	}{a, jsonSlice(a.DataPoints, func(dp metricdata.HistogramDataPoint[N]) any {
		return struct {
			metricdata.HistogramDataPoint[N]
			Attributes any
			Bounds     []any
			Min        any
			Max        any
			Sum        any
			Exemplars  []any `json:",omitempty"`
		}{dp, jsonAttributeSet(dp.Attributes), jsonSlice(dp.Bounds, jsonFloat), jsonExtrema(dp.Min), jsonExtrema(dp.Max), jsonNumber(dp.Sum), jsonSlice(dp.Exemplars, jsonExemplar[N])}
	})}
}

func jsonExponentialHistogram[N int64 | float64](a metricdata.ExponentialHistogram[N]) any {
	return struct {
		metricdata.ExponentialHistogram[N]
		DataPoints []any
	}{a, jsonSlice(a.DataPoints, func(dp metricdata.ExponentialHistogramDataPoint[N]) any {
		return struct {
			metricdata.ExponentialHistogramDataPoint[N]
			Attributes    any
			Min           any
			Max           any
			Sum           any
			ZeroThreshold any
			Exemplars     []any `json:",omitempty"`
		}{dp, jsonAttributeSet(dp.Attributes), jsonExtrema(dp.Min), jsonExtrema(dp.Max), jsonNumber(dp.Sum), jsonFloat(dp.ZeroThreshold), jsonSlice(dp.Exemplars, jsonExemplar[N])}
	})}
}

func jsonSummary(a metricdata.Summary) any {
	return struct{ DataPoints []any }{jsonSlice(a.DataPoints, func(dp metricdata.SummaryDataPoint) any {
		return struct {
			metricdata.SummaryDataPoint
			Attributes     any
			Sum            any
			QuantileValues []any
		}{dp, jsonAttributeSet(dp.Attributes), jsonFloat(dp.Sum), jsonSlice(dp.QuantileValues, func(q metricdata.QuantileValue) any {
			return struct{ Quantile, Value any }{jsonFloat(q.Quantile), jsonFloat(q.Value)}
		})}
	})}
}
