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
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// now is used to return the current local time while allowing tests to
// override the default time.Now function.
var now = time.Now

// Input receives measurements to be aggregated.
type Input[N int64 | float64] func(context.Context, N, attribute.Set)

var bgCtx = context.Background()

// Async receives asynchronous measurements that do not have a context
// associated with them.
func (f Input[N]) Async(v N, a attribute.Set) { f(bgCtx, v, a) }

// Output produces the aggregate of measurements.
type Output func(dest *metricdata.Aggregation)

// Builder builds an aggregate function.
type Builder[N int64 | float64] struct {
	// TODO: doc cumulative default temporality.
	// Temporality is the temporality used for the returned aggregate function.
	//
	// If this is not provided a default of cumulative will be used (except for
	// the last-value aggregate function where delta is the only appropriate
	// temporality).
	Temporality metricdata.Temporality
	// Filter is the attribute filter the aggregate function will use on the
	// input of measurements.
	Filter attribute.Filter
	// ReservoirFunc is the factory function used by aggregate functions to
	// create new exemplar reservoirs for a new seen attribute set.
	//
	// If this is not provided a default factory function that returns an
	// exemplar.Drop reservoir will be used.
	ReservoirFunc func() exemplar.Reservoir[N]
}

func (b Builder[N]) resFunc() func() exemplar.Reservoir[N] {
	if b.ReservoirFunc != nil {
		return b.ReservoirFunc
	}

	return exemplar.Drop[N]
}

func (b Builder[N]) input(f func(context.Context, N, attribute.Set, attribute.Set)) Input[N] {
	if b.Filter == nil {
		return func(ctx context.Context, n N, a attribute.Set) {
			f(ctx, n, a, a)
		}
	}
	return func(ctx context.Context, n N, a attribute.Set) {
		fltr, _ := a.Filter(b.Filter)
		f(ctx, n, a, fltr)
	}
}

// LastValue returns a last-value aggregate function input and output.
//
// The Builder.Temporality is ignored and delta is use always.
func (b Builder[N]) LastValue() (Input[N], Output) {
	// Delta temporality is the only temporality that makes semantic sense for
	// a last-value aggregate.
	lv := newLastValue[N](b.resFunc())

	return b.input(lv.input), func(dest *metricdata.Aggregation) {
		// Ignore if dest is not a metricdata.Gauge. The chance for memory
		// reuse of the DataPoints is missed (better luck next time).
		gData, _ := (*dest).(metricdata.Gauge[N])
		lv.output(&gData.DataPoints)
		*dest = gData
	}
}

// PrecomputedSum returns a sum aggregate function input and output. The
// arguments passed to the input are expected to be the precomputed sum values.
func (b Builder[N]) PrecomputedSum(monotonic bool) (Input[N], Output) {
	s := newPrecomputedSum[N](b.resFunc())

	var setData func(dest *[]metricdata.DataPoint[N])
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		setData = s.delta
	default:
		setData = s.cumulative
	}

	setData = b.fltrSumDPts(setData)

	return s.input, func(dest *metricdata.Aggregation) {
		// Ignore if dest is not a metricdata.Sum. The chance for memory
		// reuse of the DataPoints is missed (better luck next time).
		sData, _ := (*dest).(metricdata.Sum[N])
		sData.Temporality = b.Temporality
		sData.IsMonotonic = monotonic
		setData(&sData.DataPoints)
		*dest = sData
	}
}

func (b Builder[N]) fltrSumDPts(out func(*[]metricdata.DataPoint[N])) func(*[]metricdata.DataPoint[N]) {
	if b.Filter == nil {
		return out
	}
	f := b.Filter
	return func(dest *[]metricdata.DataPoint[N]) {
		out(dest)

		index := make(map[attribute.Distinct]int)
		var n int
		for _, dpt := range *dest {
			filtered, dropped := dpt.Attributes.Filter(f)
			key := filtered.Equivalent()

			dpt = dropExemplarAttrs[N](dpt, dropped)

			idx, ok := index[key]
			if !ok {
				// First appearance. Update with filtered dpt.
				dpt.Attributes = filtered
				(*dest)[n] = dpt
				index[key] = n
				n++
				continue
			}

			// Attributes previously recorded
			base := (*dest)[idx]
			(*dest)[idx] = foldSum[N](base, dpt)
			*dest = append((*dest)[:n], (*dest)[n+1:]...)
		}
	}
}

// Sum returns a sum aggregate function input and output.
func (b Builder[N]) Sum(monotonic bool) (Input[N], Output) {
	s := newSum[N](b.resFunc())

	var setData func(dest *[]metricdata.DataPoint[N])
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		setData = s.delta
	default:
		setData = s.cumulative
	}
	return b.input(s.input), func(dest *metricdata.Aggregation) {
		// Ignore if dest is not a metricdata.Sum. The chance for memory
		// reuse of the DataPoints is missed (better luck next time).
		sData, _ := (*dest).(metricdata.Sum[N])
		sData.Temporality = b.Temporality
		sData.IsMonotonic = monotonic
		setData(&sData.DataPoints)
		*dest = sData
	}
}

// ExplicitBucketHistogram returns a histogram aggregate function input and
// output.
func (b Builder[N]) ExplicitBucketHistogram(cfg aggregation.ExplicitBucketHistogram) (Input[N], Output) {
	h := newHistogram[N](b.resFunc(), cfg)

	var setData func(dest *[]metricdata.HistogramDataPoint[N])
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		setData = h.delta
	default:
		setData = h.cumulative
	}
	return b.input(h.input), func(dest *metricdata.Aggregation) {
		// Ignore if dest is not a metricdata.Histogram. The chance for memory
		// reuse of the DataPoints is missed (better luck next time).
		hData, _ := (*dest).(metricdata.Histogram[N])
		hData.Temporality = b.Temporality
		setData(&hData.DataPoints)
		*dest = hData
	}
}

func reset[T any](s []T, length, capacity int) []T {
	if cap(s) < capacity {
		return make([]T, length, capacity)
	}
	return s[:length]
}

func foldSum[N int64 | float64](base, overlay metricdata.DataPoint[N]) metricdata.DataPoint[N] {
	// Assumes attributes and time are the same given these are assumed sums
	// from the same collection cycle.
	if base.StartTime.After(overlay.StartTime) {
		base.StartTime = overlay.StartTime
	}
	base.Value += overlay.Value
	base.Exemplars = append(base.Exemplars, overlay.Exemplars...)
	return base
}

func dropExemplarAttrs[N int64 | float64](dpt metricdata.DataPoint[N], drop []attribute.KeyValue) metricdata.DataPoint[N] {
	if len(drop) == 0 {
		return dpt
	}

	for i, e := range dpt.Exemplars {
		e.FilteredAttributes = append(e.FilteredAttributes, drop...)
		dpt.Exemplars[i] = e
	}
	return dpt
}
