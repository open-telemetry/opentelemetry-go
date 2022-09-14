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
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	expohist "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/structure"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// expoHistValues summarizes a set of measurements as a *expohisto.Structure[N].
type expoHistValues[N int64 | float64] struct {
	config expohist.Config

	valuesMu sync.Mutex
	values   map[attribute.Set]*expohist.Structure[N]
}

func newExpoHistValues[N int64 | float64](maxSize int32) *expoHistValues[N] {
	return &expoHistValues[N]{
		maxSize: maxSize,
		values:  map[attribute.Set]*expohist.Structure[N]{},
	}
}

// Aggregate records the measurement value, scoped by attr, and aggregates it
// into an exponential histogram.
func (h *expoHistValues[N]) Aggregate(value N, attr attribute.Set) {
	h.valuesMu.Lock()
	defer h.valuesMu.Unlock()

	h, ok := h.values[attr]
	if !ok {
		agg := &expohist.Structure[N]{}
		agg.Init(h.config)
		s.values[attr] = agg
	}
	h.Update(value)
}

// NewDeltaHistogram returns an Aggregator that summarizes a set of
// measurements as an histogram. Each histogram is scoped by attributes and
// the aggregation cycle the measurements were made in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregations method is called it will reset all histogram
// counts to zero.
func NewDeltaHistogram[N int64 | float64](cfg aggregation.ExplicitBucketHistogram) Aggregator[N] {
	return &deltaHistogram[N]{
		histValues: newHistValues[N](cfg.Boundaries),
		noMinMax:   cfg.NoMinMax,
		start:      now(),
	}
}

// deltaHistogram summarizes a set of measurements made in a single
// aggregation cycle as an histogram with explicitly defined buckets.
type deltaHistogram[N int64 | float64] struct {
	*histValues[N]

	noMinMax bool
	start    time.Time
}

func (s *deltaHistogram[N]) Aggregation() metricdata.Aggregation {
	h := metricdata.Histogram{Temporality: metricdata.DeltaTemporality}

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	if len(s.values) == 0 {
		return h
	}

	// Do not allow modification of our copy of bounds.
	bounds := make([]float64, len(s.bounds))
	copy(bounds, s.bounds)
	t := now()
	h.DataPoints = make([]metricdata.HistogramDataPoint, 0, len(s.values))
	for a, b := range s.values {
		hdp := metricdata.HistogramDataPoint{
			Attributes:   a,
			StartTime:    s.start,
			Time:         t,
			Count:        b.count,
			Bounds:       bounds,
			BucketCounts: b.counts,
			Sum:          b.sum,
		}
		if !s.noMinMax {
			hdp.Min = &b.min
			hdp.Max = &b.max
		}
		h.DataPoints = append(h.DataPoints, hdp)

		// Unused attribute sets do not report.
		delete(s.values, a)
	}
	// The delta collection cycle resets.
	s.start = t
	return h
}

// NewCumulativeHistogram returns an Aggregator that summarizes a set of
// measurements as an histogram. Each histogram is scoped by attributes.
//
// Each aggregation cycle builds from the previous, the histogram counts are
// the bucketed counts of all values aggregated since the returned Aggregator
// was created.
func NewCumulativeHistogram[N int64 | float64](cfg aggregation.ExplicitBucketHistogram) Aggregator[N] {
	return &cumulativeHistogram[N]{
		histValues: newHistValues[N](cfg.Boundaries),
		noMinMax:   cfg.NoMinMax,
		start:      now(),
	}
}

// cumulativeHistogram summarizes a set of measurements made over all
// aggregation cycles as an histogram with explicitly defined buckets.
type cumulativeHistogram[N int64 | float64] struct {
	*histValues[N]

	noMinMax bool
	start    time.Time
}

func (s *cumulativeHistogram[N]) Aggregation() metricdata.Aggregation {
	h := metricdata.Histogram{Temporality: metricdata.CumulativeTemporality}

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	if len(s.values) == 0 {
		return h
	}

	// Do not allow modification of our copy of bounds.
	bounds := make([]float64, len(s.bounds))
	copy(bounds, s.bounds)
	t := now()
	h.DataPoints = make([]metricdata.HistogramDataPoint, 0, len(s.values))
	for a, b := range s.values {
		// The HistogramDataPoint field values returned need to be copies of
		// the buckets value as we will keep updating them.
		//
		// TODO (#3047): Making copies for bounds and counts incurs a large
		// memory allocation footprint. Alternatives should be explored.
		counts := make([]uint64, len(b.counts))
		copy(counts, b.counts)

		hdp := metricdata.HistogramDataPoint{
			Attributes:   a,
			StartTime:    s.start,
			Time:         t,
			Count:        b.count,
			Bounds:       bounds,
			BucketCounts: counts,
			Sum:          b.sum,
		}
		if !s.noMinMax {
			// Similar to counts, make a copy.
			min, max := b.min, b.max
			hdp.Min = &min
			hdp.Max = &max
		}
		h.DataPoints = append(h.DataPoints, hdp)
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}
	return h
}
