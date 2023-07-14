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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Input receives measurements to be aggregated.
type Input[N int64 | float64] func(context.Context, N, attribute.Set)

// Output stores the aggregate of measurements into dest and returns the number
// of aggregate data-points output.
type Output func(dest *metricdata.Aggregation) int

// Builder builds an aggregate function.
type Builder[N int64 | float64] struct {
	// Temporality is the temporality used for the returned aggregate function.
	//
	// If this is not provided a default of cumulative will be used (except for
	// the last-value aggregate function where delta is the only appropriate
	// temporality).
	Temporality metricdata.Temporality
	// Filter is the attribute filter the aggregate function will use on the
	// input of measurements.
	Filter attribute.Filter
}

func (b Builder[N]) input(agg aggregator[N]) Input[N] {
	if b.Filter != nil {
		agg = newFilter[N](agg, b.Filter)
	}
	return func(_ context.Context, n N, a attribute.Set) {
		agg.Aggregate(n, a)
	}
}

// LastValue returns a last-value aggregate function input and output.
//
// The Builder.Temporality is ignored and delta is use always.
func (b Builder[N]) LastValue() (Input[N], Output) {
	// Delta temporality is the only temporality that makes semantic sense for
	// a last-value aggregate.
	lv := newLastValue[N]()

	return b.input(lv), func(dest *metricdata.Aggregation) int {
		// TODO (#4220): optimize memory reuse here.
		*dest = lv.Aggregation()

		gData, _ := (*dest).(metricdata.Gauge[N])
		return len(gData.DataPoints)
	}
}

// PrecomputedSum returns a sum aggregate function input and output. The
// arguments passed to the input are expected to be the precomputed sum values.
func (b Builder[N]) PrecomputedSum(monotonic bool) (Input[N], Output) {
	var s aggregator[N]
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		s = newPrecomputedDeltaSum[N](monotonic)
	default:
		s = newPrecomputedCumulativeSum[N](monotonic)
	}

	return b.input(s), func(dest *metricdata.Aggregation) int {
		// TODO (#4220): optimize memory reuse here.
		*dest = s.Aggregation()

		sData, _ := (*dest).(metricdata.Sum[N])
		return len(sData.DataPoints)
	}
}

// Sum returns a sum aggregate function input and output.
func (b Builder[N]) Sum(monotonic bool) (Input[N], Output) {
	var s aggregator[N]
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		s = newDeltaSum[N](monotonic)
	default:
		s = newCumulativeSum[N](monotonic)
	}

	return b.input(s), func(dest *metricdata.Aggregation) int {
		// TODO (#4220): optimize memory reuse here.
		*dest = s.Aggregation()

		sData, _ := (*dest).(metricdata.Sum[N])
		return len(sData.DataPoints)
	}
}

// ExplicitBucketHistogram returns a histogram aggregate function input and
// output.
func (b Builder[N]) ExplicitBucketHistogram(cfg aggregation.ExplicitBucketHistogram) (Input[N], Output) {
	var h aggregator[N]
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		h = newDeltaHistogram[N](cfg)
	default:
		h = newCumulativeHistogram[N](cfg)
	}
	return b.input(h), func(dest *metricdata.Aggregation) int {
		// TODO (#4220): optimize memory reuse here.
		*dest = h.Aggregation()

		hData, _ := (*dest).(metricdata.Histogram[N])
		return len(hData.DataPoints)
	}
}
