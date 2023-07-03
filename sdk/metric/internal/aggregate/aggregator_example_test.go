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

package aggregate

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type meter struct {
	// When a reader initiates a collection, the meter would collect
	// aggregations from each of these functions.
	aggregations []metricdata.Aggregation
}

func (p *meter) Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	// This is an example of how a meter would create an aggregator for a new
	// counter. At this point the provider would determine the aggregation and
	// temporality to used based on the Reader and View configuration. Assume
	// here these are determined to be a cumulative sum.

	aggregator := NewCumulativeSum[int64](true)
	count := inst{aggregateFunc: aggregator.Aggregate}

	p.aggregations = append(p.aggregations, aggregator.Aggregation())

	fmt.Printf("using %T aggregator for counter\n", aggregator)

	return count, nil
}

func (p *meter) Int64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	// This is an example of how a meter would create an aggregator for a new
	// up-down counter. At this point the provider would determine the
	// aggregation and temporality to used based on the Reader and View
	// configuration. Assume here these are determined to be a last-value
	// aggregation (the temporality does not affect the produced aggregations).

	aggregator := NewLastValue[int64]()
	upDownCount := inst{aggregateFunc: aggregator.Aggregate}

	p.aggregations = append(p.aggregations, aggregator.Aggregation())

	fmt.Printf("using %T aggregator for up-down counter\n", aggregator)

	return upDownCount, nil
}

func (p *meter) Int64Histogram(string, ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	// This is an example of how a meter would create an aggregator for a new
	// histogram. At this point the provider would determine the aggregation
	// and temporality to used based on the Reader and View configuration.
	// Assume here these are determined to be a delta explicit-bucket
	// histogram.

	aggregator := NewDeltaHistogram[int64](aggregation.ExplicitBucketHistogram{
		Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
		NoMinMax:   false,
	})
	hist := inst{aggregateFunc: aggregator.Aggregate}

	p.aggregations = append(p.aggregations, aggregator.Aggregation())

	fmt.Printf("using %T aggregator for histogram\n", aggregator)

	return hist, nil
}

// inst is a generalized int64 synchronous counter, up-down counter, and
// histogram used for demonstration purposes only.
type inst struct {
	aggregateFunc func(int64, attribute.Set)

	embedded.Int64Counter
	embedded.Int64UpDownCounter
	embedded.Int64Histogram
}

func (inst) Add(context.Context, int64, ...metric.AddOption)       {}
func (inst) Record(context.Context, int64, ...metric.RecordOption) {}

func Example() {
	m := meter{}

	_, _ = m.Int64Counter("counter example")
	_, _ = m.Int64UpDownCounter("up-down counter example")
	_, _ = m.Int64Histogram("histogram example")

	// Output:
	// using *aggregate.cumulativeSum[int64] aggregator for counter
	// using *aggregate.lastValue[int64] aggregator for up-down counter
	// using *aggregate.deltaHistogram[int64] aggregator for histogram
}
