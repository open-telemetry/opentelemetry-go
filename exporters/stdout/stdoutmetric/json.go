// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric

import (
	"encoding/json"
	"io"
	"math"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// jsonEncoder applies the non-OTLP attribute representation and handles
// non-finite metric values without changing caller-owned data.
type jsonEncoder struct{ *json.Encoder }

func newJSONEncoder(w io.Writer) *jsonEncoder {
	return &jsonEncoder{json.NewEncoder(w)}
}

func (e *jsonEncoder) Encode(v any) error {
	if rm, ok := v.(*metricdata.ResourceMetrics); ok && rm != nil {
		v = jsonResourceMetrics(rm)
	}
	return e.Encoder.Encode(v)
}

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

func jsonAttributes(kvs []attribute.KeyValue) []json.RawMessage {
	return jsonSlice(kvs, func(kv attribute.KeyValue) json.RawMessage {
		return json.RawMessage(kv.String())
	})
}

func jsonAttributeSet(set attribute.Set) json.RawMessage {
	return json.RawMessage(set.String())
}

type resourceMetricsJSON struct {
	Resource     json.RawMessage
	ScopeMetrics []scopeMetricsJSON
}

type scopeMetricsJSON struct {
	Scope   scopeJSON
	Metrics []metricsJSON
}

type scopeJSON struct {
	instrumentation.Scope
	Attributes json.RawMessage
}

type metricsJSON struct {
	metricdata.Metrics
	Data any
}

func jsonResourceMetrics(rm *metricdata.ResourceMetrics) resourceMetricsJSON {
	var res json.RawMessage
	if rm.Resource != nil {
		res = jsonAttributeSet(*rm.Resource.Set())
	}
	return resourceMetricsJSON{res, jsonSlice(rm.ScopeMetrics, func(sm metricdata.ScopeMetrics) scopeMetricsJSON {
		return scopeMetricsJSON{
			Scope: scopeJSON{sm.Scope, jsonAttributeSet(sm.Scope.Attributes)},
			Metrics: jsonSlice(sm.Metrics, func(m metricdata.Metrics) metricsJSON {
				return metricsJSON{m, jsonAggregation(m.Data)}
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
	return struct{ DataPoints []dataPointJSON[N] }{jsonSlice(a.DataPoints, jsonDataPoint[N])}
}

func jsonSum[N int64 | float64](a metricdata.Sum[N]) any {
	return struct {
		metricdata.Sum[N]
		DataPoints []dataPointJSON[N]
	}{a, jsonSlice(a.DataPoints, jsonDataPoint[N])}
}

type dataPointJSON[N int64 | float64] struct {
	metricdata.DataPoint[N]
	Attributes json.RawMessage
	Value      any
	Exemplars  []exemplarJSON[N] `json:",omitempty"`
}

func jsonDataPoint[N int64 | float64](dp metricdata.DataPoint[N]) dataPointJSON[N] {
	return dataPointJSON[N]{dp, jsonAttributeSet(dp.Attributes), jsonNumber(dp.Value), jsonSlice(dp.Exemplars, jsonExemplar[N])}
}

type exemplarJSON[N int64 | float64] struct {
	metricdata.Exemplar[N]
	FilteredAttributes []json.RawMessage
	Value              any
}

func jsonExemplar[N int64 | float64](ex metricdata.Exemplar[N]) exemplarJSON[N] {
	return exemplarJSON[N]{ex, jsonAttributes(ex.FilteredAttributes), jsonNumber(ex.Value)}
}

func jsonExtrema[N int64 | float64](e metricdata.Extrema[N]) any {
	if value, ok := e.Value(); ok {
		return jsonNumber(value)
	}
	return nil
}

type histogramDataPointJSON[N int64 | float64] struct {
	metricdata.HistogramDataPoint[N]
	Attributes json.RawMessage
	Bounds     []any
	Min        any
	Max        any
	Sum        any
	Exemplars  []exemplarJSON[N] `json:",omitempty"`
}

func jsonHistogram[N int64 | float64](a metricdata.Histogram[N]) any {
	return struct {
		metricdata.Histogram[N]
		DataPoints []histogramDataPointJSON[N]
	}{a, jsonSlice(a.DataPoints, func(dp metricdata.HistogramDataPoint[N]) histogramDataPointJSON[N] {
		return histogramDataPointJSON[N]{dp, jsonAttributeSet(dp.Attributes), jsonSlice(dp.Bounds, jsonFloat), jsonExtrema(dp.Min), jsonExtrema(dp.Max), jsonNumber(dp.Sum), jsonSlice(dp.Exemplars, jsonExemplar[N])}
	})}
}

type exponentialHistogramDataPointJSON[N int64 | float64] struct {
	metricdata.ExponentialHistogramDataPoint[N]
	Attributes    json.RawMessage
	Min           any
	Max           any
	Sum           any
	ZeroThreshold any
	Exemplars     []exemplarJSON[N] `json:",omitempty"`
}

func jsonExponentialHistogram[N int64 | float64](a metricdata.ExponentialHistogram[N]) any {
	return struct {
		metricdata.ExponentialHistogram[N]
		DataPoints []exponentialHistogramDataPointJSON[N]
	}{a, jsonSlice(a.DataPoints, func(dp metricdata.ExponentialHistogramDataPoint[N]) exponentialHistogramDataPointJSON[N] {
		return exponentialHistogramDataPointJSON[N]{dp, jsonAttributeSet(dp.Attributes), jsonExtrema(dp.Min), jsonExtrema(dp.Max), jsonNumber(dp.Sum), jsonFloat(dp.ZeroThreshold), jsonSlice(dp.Exemplars, jsonExemplar[N])}
	})}
}

type summaryDataPointJSON struct {
	metricdata.SummaryDataPoint
	Attributes     json.RawMessage
	Sum            any
	QuantileValues []quantileValueJSON
}

type quantileValueJSON struct{ Quantile, Value any }

func jsonSummary(a metricdata.Summary) any {
	return struct{ DataPoints []summaryDataPointJSON }{jsonSlice(a.DataPoints, func(dp metricdata.SummaryDataPoint) summaryDataPointJSON {
		return summaryDataPointJSON{dp, jsonAttributeSet(dp.Attributes), jsonFloat(dp.Sum), jsonSlice(dp.QuantileValues, func(q metricdata.QuantileValue) quantileValueJSON {
			return quantileValueJSON{jsonFloat(q.Quantile), jsonFloat(q.Value)}
		})}
	})}
}
