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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	expohist "go.opentelemetry.io/otel/sdk/metric/metricdata/exponential/structure"
)

// expoHistValues summarizes a set of measurements as an expoHistValues with
// exponentially-defined buckets.
type expoHistValues[N int64 | float64] struct {
	cfg expohist.Config

	values   map[attribute.Set]*expohist.Histogram[N]
	valuesMu sync.Mutex
}

func newExpoHistValues[N int64 | float64](maxSize int32) *expoHistValues[N] {
	cfg := expohist.NewConfig(expohist.WithMaxSize(maxSize))
	cfg, _ = cfg.Validate()
	return &expoHistValues[N]{
		cfg:    cfg,
		values: make(map[attribute.Set]*expohist.Histogram[N]),
	}
}

// Aggregate records the measurement value, scoped by attr, and aggregates it
// into a histogram.
func (hv *expoHistValues[N]) Aggregate(value N, attr attribute.Set) {
	hv.valuesMu.Lock()
	defer hv.valuesMu.Unlock()

	h, ok := hv.values[attr]
	if !ok {
		h = &expohist.Histogram[N]{}
		h.Init(hv.cfg)
		hv.values[attr] = h
	}
	h.Update(value)
}

// NewDeltaExponentialHistogram returns an Aggregator that summarizes a set of
// measurements as an exponential histogram. Each histogram is scoped by attributes and
// the aggregation cycle the measurements were made in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregations method is called it will reset all histogram
// counts to zero.
func NewDeltaExponentialHistogram[N int64 | float64](cfg aggregation.ExponentialHistogram) Aggregator[N] {
	return &deltaExpoHistogram[N]{
		expoHistValues: newExpoHistValues[N](cfg.MaxSize),
		noMinMax:       cfg.NoMinMax,
		start:          now(),
	}
}

// deltaExpoHistogram summarizes a set of measurements made in a single
// aggregation cycle as an histogram with explicitly defined buckets.
type deltaExpoHistogram[N int64 | float64] struct {
	*expoHistValues[N]

	noMinMax bool
	start    time.Time
}

func (s *deltaExpoHistogram[N]) Aggregation() metricdata.Aggregation {
	h := metricdata.ExponentialHistogram{Temporality: metricdata.DeltaTemporality}

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	if len(s.values) == 0 {
		return h
	}

	t := now()
	h.DataPoints = make([]metricdata.ExponentialHistogramDataPoint, 0, len(s.values))
	for a, b := range s.values {
		hdp := metricdata.ExponentialHistogramDataPoint{
			Attributes: a,
			StartTime:  s.start,
			Time:       t,
			Count:      b.Count(),
			Sum:        float64(b.Sum()),
			ZeroCount:  b.ZeroCount(),
			Scale:      b.Scale(),
			Positive:   expoBuckets(b.Positive()),
			// b.Negative() skipped
		}
		if !s.noMinMax && b.Count() != 0 {
			min, max := float64(b.Min()), float64(b.Max())
			hdp.Min = &min
			hdp.Max = &max
		}
		h.DataPoints = append(h.DataPoints, hdp)

		// Unused attribute sets do not report.
		delete(s.values, a)
	}
	// The delta collection cycle resets.
	s.start = t
	return h
}

func expoBuckets(b *expohist.Buckets) metricdata.ExponentialBuckets {
	if b.Len() == 0 {
		return metricdata.ExponentialBuckets{}
	}
	cnts := make([]uint64, b.Len())
	for i := 0; i < len(cnts); i++ {
		cnts[i] = b.At(uint32(i))
	}
	return metricdata.ExponentialBuckets{
		Offset:       b.Offset(),
		BucketCounts: cnts,
	}
}

// NewExponentialCumulativeHistogram returns an Aggregator that summarizes a set of
// measurements as an exponential histogram. Each histogram is scoped by attributes.
//
// Each aggregation cycle builds from the previous, the histogram counts are
// the bucketed counts of all values aggregated since the returned Aggregator
// was created.
func NewCumulativeExponentialHistogram[N int64 | float64](cfg aggregation.ExponentialHistogram) Aggregator[N] {
	return &cumulativeExpoHistogram[N]{
		expoHistValues: newExpoHistValues[N](cfg.MaxSize),
		noMinMax:       cfg.NoMinMax,
		start:          now(),
	}
}

// cumulativeExpoHistogram summarizes a set of measurements made over all
// aggregation cycles as an histogram with explicitly defined buckets.
type cumulativeExpoHistogram[N int64 | float64] struct {
	*expoHistValues[N]

	noMinMax bool
	start    time.Time
}

func (s *cumulativeExpoHistogram[N]) Aggregation() metricdata.Aggregation {
	h := metricdata.ExponentialHistogram{Temporality: metricdata.CumulativeTemporality}

	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	if len(s.values) == 0 {
		return h
	}

	// Do not allow modification of our copy of bounds.
	t := now()
	h.DataPoints = make([]metricdata.ExponentialHistogramDataPoint, 0, len(s.values))
	for a, b := range s.values {
		hdp := metricdata.ExponentialHistogramDataPoint{
			Attributes: a,
			StartTime:  s.start,
			Time:       t,
			Count:      b.Count(),
			Sum:        float64(b.Sum()),
			ZeroCount:  b.ZeroCount(),
			Scale:      b.Scale(),
			Positive:   expoBuckets(b.Positive()),
			// b.Negative() skipped
		}
		if !s.noMinMax {
			min, max := float64(b.Min()), float64(b.Max())
			hdp.Min = &min
			hdp.Max = &max
		}
		h.DataPoints = append(h.DataPoints, hdp)
		// TODO (#3006): This will use an unbounded amount of memory.
	}
	return h
}
