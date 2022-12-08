// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TimelessJsonEncoder is an encoder that outputs JSON without timestamps.
type timelessJsonEncoder struct {
	enc Encoder
}

func (e timelessJsonEncoder) Encode(v any) error {
	orig, ok := v.(metricdata.ResourceMetrics)
	if !ok {
		return e.enc.Encode(v)
	}
	rm := metricdata.ResourceMetrics{
		Resource:     orig.Resource,
		ScopeMetrics: timlessCopyScopeMetrics(orig.ScopeMetrics),
	}
	return e.enc.Encode(rm)
}

func timlessCopyScopeMetrics(sms []metricdata.ScopeMetrics) []metricdata.ScopeMetrics {
	out := make([]metricdata.ScopeMetrics, 0, len(sms))
	for _, sm := range sms {
		out = append(out, metricdata.ScopeMetrics{
			Scope:   sm.Scope,
			Metrics: timelessCopyMetrics(sm.Metrics),
		})
	}
	return out
}

func timelessCopyMetrics(ms []metricdata.Metrics) []metricdata.Metrics {
	out := make([]metricdata.Metrics, 0, len(ms))
	for _, m := range ms {
		out = append(out, metricdata.Metrics{
			Name:        m.Name,
			Description: m.Description,
			Unit:        m.Unit,
			Data:        timelessCopyAggregation(m.Data),
		})
	}
	return out
}

func timelessCopyAggregation(a metricdata.Aggregation) metricdata.Aggregation {
	switch a := a.(type) {
	case metricdata.Histogram:
		return metricdata.Histogram{
			Temporality: a.Temporality,
			DataPoints:  timelessCopyHistogramDataPoints(a.DataPoints),
		}
	case metricdata.Sum[int64]:
		return metricdata.Sum[int64]{
			Temporality: a.Temporality,
			DataPoints:  timelessCopyDataPoints(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Sum[float64]:
		return metricdata.Sum[float64]{
			Temporality: a.Temporality,
			DataPoints:  timelessCopyDataPoints(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Gauge[int64]:
		return metricdata.Gauge[int64]{
			DataPoints: timelessCopyDataPoints(a.DataPoints),
		}
	case metricdata.Gauge[float64]:
		return metricdata.Gauge[float64]{
			DataPoints: timelessCopyDataPoints(a.DataPoints),
		}
	}
	return a
}

func timelessCopyHistogramDataPoints(hdp []metricdata.HistogramDataPoint) []metricdata.HistogramDataPoint {
	out := make([]metricdata.HistogramDataPoint, 0, len(hdp))
	for _, dp := range hdp {
		out = append(out, metricdata.HistogramDataPoint{
			Attributes:   dp.Attributes,
			Count:        dp.Count,
			Sum:          dp.Sum,
			Bounds:       dp.Bounds,
			BucketCounts: dp.BucketCounts,
			Min:          dp.Min,
			Max:          dp.Max,
		})
	}
	return out
}

func timelessCopyDataPoints[T int64 | float64](sdp []metricdata.DataPoint[T]) []metricdata.DataPoint[T] {
	out := make([]metricdata.DataPoint[T], 0, len(sdp))
	for _, dp := range sdp {
		out = append(out, metricdata.DataPoint[T]{
			Attributes: dp.Attributes,
			Value:      dp.Value,
		})
	}
	return out
}

func WithTimelessEncoder(enc Encoder) Option {
	return WithEncoder(timelessJsonEncoder{enc})
}
