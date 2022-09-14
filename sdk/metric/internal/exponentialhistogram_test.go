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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/exponential/mapping/logarithm"
)

var (
	expoHistConf = aggregation.ExponentialHistogram{
		MaxSize:  4,
		NoMinMax: false,
	}
)

func TestExponentialHistogram(t *testing.T) {
	t.Cleanup(mockTime(now))
	t.Run("Int64", testExponentialHistogram[int64])
	t.Run("Float64", testExponentialHistogram[float64])
}

func testExponentialHistogram[N int64 | float64](t *testing.T) {
	tester := &aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}

	incr := monoIncr
	eFunc := deltaExpoHistExpecter(incr)
	t.Run("Delta", tester.Run(NewDeltaExponentialHistogram[N](expoHistConf), incr, eFunc))
	eFunc = cumuExpoHistExpecter(incr)
	t.Run("Cumulative", tester.Run(NewCumulativeExponentialHistogram[N](expoHistConf), incr, eFunc))
}

func deltaExpoHistExpecter(incr setMap) expectFunc {
	h := metricdata.ExponentialHistogram{Temporality: metricdata.DeltaTemporality}
	return func(m int) metricdata.Aggregation {
		h.DataPoints = make([]metricdata.ExponentialHistogramDataPoint, 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, ehPoint(a, float64(v), uint64(m)))
		}
		return h
	}
}

func cumuExpoHistExpecter(incr setMap) expectFunc {
	var cycle int
	h := metricdata.ExponentialHistogram{Temporality: metricdata.CumulativeTemporality}
	return func(m int) metricdata.Aggregation {
		cycle++
		h.DataPoints = make([]metricdata.ExponentialHistogramDataPoint, 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, ehPoint(a, float64(v), uint64(cycle*m)))
		}
		return h
	}
}

// ehPoint returns an ExponentialHistogramDataPoint that started and ended now with multi
// number of measurements values v. It includes a min and max (set to v).
func ehPoint(a attribute.Set, v float64, multi uint64) metricdata.ExponentialHistogramDataPoint {
	mapping, _ := logarithm.NewMapping(logarithm.MaxScale)
	offset := mapping.MapToIndex(v)
	counts := make([]uint64, 1)
	counts[0] = multi
	return metricdata.ExponentialHistogramDataPoint{
		Attributes: a,
		StartTime:  now(),
		Time:       now(),
		Count:      multi,
		Scale:      logarithm.MaxScale,
		ZeroCount:  0,
		Positive: metricdata.ExponentialBuckets{
			Offset:       offset,
			BucketCounts: counts,
		},
		Min: &v,
		Max: &v,
		Sum: v * float64(multi),
	}
}

// TODO: more tests like histogram_test.go has

func BenchmarkExponentialHistogram(b *testing.B) {
	b.Run("Int64", benchmarkExponentialHistogram[int64])
	b.Run("Float64", benchmarkExponentialHistogram[float64])
}

func benchmarkExponentialHistogram[N int64 | float64](b *testing.B) {
	factory := func() Aggregator[N] { return NewDeltaExponentialHistogram[N](expoHistConf) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() Aggregator[N] { return NewCumulativeExponentialHistogram[N](expoHistConf) }
	b.Run("Cumulative", benchmarkAggregator(factory))
}
