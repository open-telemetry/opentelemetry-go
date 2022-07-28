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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"sort"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var (
	bounds   = []float64{1, 5}
	histConf = aggregation.ExplicitBucketHistogram{
		Boundaries: bounds,
		NoMinMax:   false,
	}
)

func TestHistogram(t *testing.T) {
	t.Cleanup(mockTime(now))
	t.Run("Int64", testHistogram[int64])
	t.Run("Float64", testHistogram[float64])
}

func testHistogram[N int64 | float64](t *testing.T) {
	tester := &aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}

	incr := monoIncr
	eFunc := deltaHistExpecter(incr)
	t.Run("Delta", tester.Run(NewDeltaHistogram[N](histConf), incr, eFunc))
	eFunc = cumuHistExpecter(incr)
	t.Run("Cumulative", tester.Run(NewCumulativeHistogram[N](histConf), incr, eFunc))
}

func deltaHistExpecter(incr setMap) expectFunc {
	h := metricdata.Histogram{Temporality: metricdata.DeltaTemporality}
	return func(m int) metricdata.Aggregation {
		h.DataPoints = make([]metricdata.HistogramDataPoint, 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, hPoint(a, float64(v), uint64(m)))
		}
		return h
	}
}

func cumuHistExpecter(incr setMap) expectFunc {
	var cycle int
	h := metricdata.Histogram{Temporality: metricdata.CumulativeTemporality}
	return func(m int) metricdata.Aggregation {
		cycle++
		h.DataPoints = make([]metricdata.HistogramDataPoint, 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, hPoint(a, float64(v), uint64(cycle*m)))
		}
		return h
	}
}

// hPoint returns an HistogramDataPoint that started and ended now with multi
// number of measurements values v. It includes a min and max (set to v).
func hPoint(a attribute.Set, v float64, multi uint64) metricdata.HistogramDataPoint {
	idx := sort.SearchFloat64s(bounds, v)
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint{
		Attributes:   a,
		StartTime:    now(),
		Time:         now(),
		Count:        multi,
		Bounds:       bounds,
		BucketCounts: counts,
		Min:          &v,
		Max:          &v,
		Sum:          v * float64(multi),
	}
}
