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
	"sort"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func newHistogram[N int64 | float64](r func() exemplar.Reservoir[N], cfg aggregation.ExplicitBucketHistogram) *hist[N] {
	// The responsibility of keeping all buckets correctly associated with the
	// passed boundaries is ultimately this type's responsibility. Make a copy
	// here so we can always guarantee this. Or, in the case of failure, have
	// complete control over the fix.
	b := make([]float64, len(cfg.Boundaries))
	copy(b, cfg.Boundaries)
	sort.Float64s(b)
	return &hist[N]{
		newRes:   r,
		noMinMax: cfg.NoMinMax,
		start:    now(),
		bounds:   b,
		values:   make(map[attribute.Distinct]*buckets[N]),
	}
}

type buckets[N int64 | float64] struct {
	attr attribute.Set
	res  exemplar.Reservoir[N]

	counts   []uint64
	count    uint64
	sum      N
	min, max N
}

// newBuckets returns buckets with n bins.
func newBuckets[N int64 | float64](n int) *buckets[N] {
	return &buckets[N]{counts: make([]uint64, n)}
}

func (b *buckets[N]) bin(idx int, value N) {
	b.counts[idx]++
	b.count++
	b.sum += value
	if value < b.min {
		b.min = value
	} else if value > b.max {
		b.max = value
	}
}

// hist summarizes a set of measurements as an histogram with explicitly
// defined buckets.
type hist[N int64 | float64] struct {
	noMinMax bool
	start    time.Time

	newRes func() exemplar.Reservoir[N]

	bounds   []float64
	valuesMu sync.Mutex
	values   map[attribute.Distinct]*buckets[N]
}

// Aggregate records the measurement value, scoped by attr, and aggregates it
// into a histogram.
func (s *hist[N]) input(ctx context.Context, value N, origAttr, fltrAttr attribute.Set) {
	// This search will return an index in the range [0, len(s.bounds)], where
	// it will return len(s.bounds) if value is greater than the last element
	// of s.bounds. This aligns with the buckets in that the length of buckets
	// is len(s.bounds)+1, with the last bucket representing:
	// (s.bounds[len(s.bounds)-1], +∞).
	idx := sort.SearchFloat64s(s.bounds, float64(value))

	t := now()
	key := fltrAttr.Equivalent()

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	b, ok := s.values[key]
	if !ok {
		b.attr = fltrAttr
		b.res = s.newRes()

		// N+1 buckets. For example:
		//
		//   bounds = [0, 5, 10]
		//
		// Then,
		//
		//   buckets = (-∞, 0], (0, 5.0], (5.0, 10.0], (10.0, +∞)
		b = newBuckets[N](len(s.bounds) + 1)
		// Ensure min and max are recorded values (not zero), for new buckets.
		b.min, b.max = value, value
		s.values[key] = b
	}
	b.bin(idx, value)
	b.res.Offer(ctx, t, value, origAttr)
}

func (s *hist[N]) delta(dest *[]metricdata.HistogramDataPoint[N]) {
	t := now()

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	nBounds := len(s.bounds)
	n := len(s.values)
	*dest = reset(*dest, n, n)
	var i int
	for key, buckets := range s.values {
		(*dest)[i].Attributes = buckets.attr
		(*dest)[i].StartTime = s.start
		(*dest)[i].Time = t
		(*dest)[i].Count = buckets.count
		// TODO: It is inefficient to not pool b.counts and just copy
		// values here similar to the cumulative case.
		(*dest)[i].BucketCounts = buckets.counts
		(*dest)[i].Sum = buckets.sum
		buckets.res.Flush(&(*dest)[i].Exemplars, buckets.attr)

		// Do not allow modification of our copy of bounds.
		reset((*dest)[i].Bounds, nBounds, nBounds)
		copy((*dest)[i].Bounds, s.bounds)

		if !s.noMinMax {
			(*dest)[i].Min = metricdata.NewExtrema(buckets.min)
			(*dest)[i].Max = metricdata.NewExtrema(buckets.max)
		} else {
			(*dest)[i].Min = metricdata.Extrema[N]{}
			(*dest)[i].Max = metricdata.Extrema[N]{}
		}
		i++

		// Unused attribute sets do not report.
		delete(s.values, key)
	}
	// The delta collection cycle resets.
	s.start = t
}

func (s *hist[N]) cumulative(dest *[]metricdata.HistogramDataPoint[N]) {
	t := now()

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	nBounds := len(s.bounds)
	n := len(s.values)
	*dest = reset(*dest, n, n)
	var i int
	for _, buckets := range s.values {
		(*dest)[i].Attributes = buckets.attr
		(*dest)[i].StartTime = s.start
		(*dest)[i].Time = t
		(*dest)[i].Count = buckets.count
		(*dest)[i].Sum = buckets.sum
		buckets.res.Collect(&(*dest)[i].Exemplars, buckets.attr)

		// The HistogramDataPoint field values returned need to be copies of
		// the buckets value as we will keep updating them.
		reset((*dest)[i].BucketCounts, nBounds+1, nBounds+1)
		copy((*dest)[i].BucketCounts, buckets.counts)

		// Do not allow modification of our copy of bounds.
		reset((*dest)[i].Bounds, nBounds, nBounds)
		copy((*dest)[i].Bounds, s.bounds)

		if !s.noMinMax {
			(*dest)[i].Min = metricdata.NewExtrema(buckets.min)
			(*dest)[i].Max = metricdata.NewExtrema(buckets.max)
		} else {
			(*dest)[i].Min = metricdata.Extrema[N]{}
			(*dest)[i].Max = metricdata.Extrema[N]{}
		}
		i++
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}
}
