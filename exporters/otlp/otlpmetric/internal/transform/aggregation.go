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

//go:build go1.18
// +build go1.18

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

func Gauge[N int64 | float64](g metricdata.Gauge[N]) *mpb.Metric_Gauge {
	return &mpb.Metric_Gauge{
		Gauge: &mpb.Gauge{
			DataPoints: DataPoints(g.DataPoints),
		},
	}
}

func Sum[N int64 | float64](s metricdata.Sum[N]) (*mpb.Metric_Sum, error) {
	t, err := Temporality(s.Temporality)
	return &mpb.Metric_Sum{
		Sum: &mpb.Sum{
			AggregationTemporality: t,
			IsMonotonic:            s.IsMonotonic,
			DataPoints:             DataPoints(s.DataPoints),
		},
	}, err
}

func DataPoints[N int64 | float64](dPts []metricdata.DataPoint[N]) []*mpb.NumberDataPoint {
	out := make([]*mpb.NumberDataPoint, 0, len(dPts))
	for _, dPt := range dPts {
		ndp := &mpb.NumberDataPoint{
			Attributes:        Iterator(dPt.Attributes.Iter()),
			StartTimeUnixNano: uint64(dPt.StartTime.UnixNano()),
			TimeUnixNano:      uint64(dPt.Time.UnixNano()),
		}
		switch v := any(dPt.Value).(type) {
		case int64:
			ndp.Value = &mpb.NumberDataPoint_AsInt{
				AsInt: v,
			}
		case float64:
			ndp.Value = &mpb.NumberDataPoint_AsDouble{
				AsDouble: v,
			}
		}
		out = append(out, ndp)
	}
	return out
}

func Histogram(h metricdata.Histogram) (*mpb.Metric_Histogram, error) {
	t, err := Temporality(h.Temporality)
	return &mpb.Metric_Histogram{
		Histogram: &mpb.Histogram{
			AggregationTemporality: t,
			DataPoints:             HistogramDataPoints(h.DataPoints),
		},
	}, err
}

func HistogramDataPoints(dPts []metricdata.HistogramDataPoint) []*mpb.HistogramDataPoint {
	out := make([]*mpb.HistogramDataPoint, 0, len(dPts))
	for _, dPt := range dPts {
		out = append(out, &mpb.HistogramDataPoint{
			Attributes:        Iterator(dPt.Attributes.Iter()),
			StartTimeUnixNano: uint64(dPt.StartTime.UnixNano()),
			TimeUnixNano:      uint64(dPt.Time.UnixNano()),
			Count:             dPt.Count,
			Sum:               &dPt.Sum,
			BucketCounts:      dPt.BucketCounts,
			ExplicitBounds:    dPt.Bounds,
			Min:               dPt.Min,
			Max:               dPt.Max,
		})
	}
	return out
}

func Temporality(t metricdata.Temporality) (mpb.AggregationTemporality, error) {
	switch t {
	case metricdata.DeltaTemporality:
		return mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, nil
	case metricdata.CumulativeTemporality:
		return mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, nil
	default:
		err := fmt.Errorf("unknown temporality: %s", t)
		return mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED, err
	}
}
