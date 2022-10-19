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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// valueMap is the storage for all sums.
type valueMap[N int64 | float64] struct {
	sync.Mutex
	values map[attribute.Set]N
}

func newValueMap[N int64 | float64]() *valueMap[N] {
	return &valueMap[N]{values: make(map[attribute.Set]N)}
}

func (s *valueMap[N]) set(value N, attr attribute.Set) { // nolint: unused  // This is indeed used.
	s.Lock()
	s.values[attr] = value
	s.Unlock()
}

func (s *valueMap[N]) Aggregate(value N, attr attribute.Set) {
	s.Lock()
	s.values[attr] += value
	s.Unlock()
}

// NewDeltaSum returns an Aggregator that summarizes a set of measurements as
// their arithmetic sum. Each sum is scoped by attributes and the aggregation
// cycle the measurements were made in.
//
// The monotonic value is used to communicate the produced Aggregation is
// monotonic or not. The returned Aggregator does not make any guarantees this
// value is accurate. It is up to the caller to ensure it.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewDeltaSum[N int64 | float64](monotonic bool) Aggregator[N] {
	return newDeltaSum[N](monotonic)
}

func newDeltaSum[N int64 | float64](monotonic bool) *deltaSum[N] {
	return &deltaSum[N]{
		valueMap:  newValueMap[N](),
		monotonic: monotonic,
		start:     now(),
	}
}

// deltaSum summarizes a set of measurements made in a single aggregation
// cycle as their arithmetic sum.
type deltaSum[N int64 | float64] struct {
	*valueMap[N]

	monotonic bool
	start     time.Time
}

func (s *deltaSum[N]) Aggregation() metricdata.Aggregation {
	out := metricdata.Sum[N]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: s.monotonic,
	}

	s.Lock()
	defer s.Unlock()

	if len(s.values) == 0 {
		return out
	}

	t := now()
	out.DataPoints = make([]metricdata.DataPoint[N], 0, len(s.values))
	for attr, value := range s.values {
		out.DataPoints = append(out.DataPoints, metricdata.DataPoint[N]{
			Attributes: attr,
			StartTime:  s.start,
			Time:       t,
			Value:      value,
		})
		// Unused attribute sets do not report.
		delete(s.values, attr)
	}
	// The delta collection cycle resets.
	s.start = t
	return out
}

// NewCumulativeSum returns an Aggregator that summarizes a set of
// measurements as their arithmetic sum. Each sum is scoped by attributes and
// the aggregation cycle the measurements were made in.
//
// The monotonic value is used to communicate the produced Aggregation is
// monotonic or not. The returned Aggregator does not make any guarantees this
// value is accurate. It is up to the caller to ensure it.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewCumulativeSum[N int64 | float64](monotonic bool) Aggregator[N] {
	return newCumulativeSum[N](monotonic)
}

func newCumulativeSum[N int64 | float64](monotonic bool) *cumulativeSum[N] {
	return &cumulativeSum[N]{
		valueMap:  newValueMap[N](),
		monotonic: monotonic,
		start:     now(),
	}
}

// cumulativeSum summarizes a set of measurements made over all aggregation
// cycles as their arithmetic sum.
type cumulativeSum[N int64 | float64] struct {
	*valueMap[N]

	monotonic bool
	start     time.Time
}

func (s *cumulativeSum[N]) Aggregation() metricdata.Aggregation {
	out := metricdata.Sum[N]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: s.monotonic,
	}

	s.Lock()
	defer s.Unlock()

	if len(s.values) == 0 {
		return out
	}

	t := now()
	out.DataPoints = make([]metricdata.DataPoint[N], 0, len(s.values))
	for attr, value := range s.values {
		out.DataPoints = append(out.DataPoints, metricdata.DataPoint[N]{
			Attributes: attr,
			StartTime:  s.start,
			Time:       t,
			Value:      value,
		})
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}
	return out
}

// NewPrecomputedDeltaSum returns an Aggregator that summarizes a set of
// measurements as their pre-computed arithmetic sum. Each sum is scoped by
// attributes and the aggregation cycle the measurements were made in.
//
// The monotonic value is used to communicate the produced Aggregation is
// monotonic or not. The returned Aggregator does not make any guarantees this
// value is accurate. It is up to the caller to ensure it.
//
// The output Aggregation will report recorded values as delta temporality. It
// is up to the caller to ensure this is accurate.
func NewPrecomputedDeltaSum[N int64 | float64](monotonic bool) Aggregator[N] {
	return &precomputedSum[N]{settableSum: newDeltaSum[N](monotonic)}
}

// NewPrecomputedCumulativeSum returns an Aggregator that summarizes a set of
// measurements as their pre-computed arithmetic sum. Each sum is scoped by
// attributes and the aggregation cycle the measurements were made in.
//
// The monotonic value is used to communicate the produced Aggregation is
// monotonic or not. The returned Aggregator does not make any guarantees this
// value is accurate. It is up to the caller to ensure it.
//
// The output Aggregation will report recorded values as cumulative
// temporality. It is up to the caller to ensure this is accurate.
func NewPrecomputedCumulativeSum[N int64 | float64](monotonic bool) Aggregator[N] {
	return &precomputedSum[N]{settableSum: newCumulativeSum[N](monotonic)}
}

type settableSum[N int64 | float64] interface {
	set(value N, attr attribute.Set)
	Aggregation() metricdata.Aggregation
}

// precomputedSum summarizes a set of measurements recorded over all
// aggregation cycles directly as an arithmetic sum.
type precomputedSum[N int64 | float64] struct {
	settableSum[N]
}

// Aggregate records value directly as a sum for attr.
func (s *precomputedSum[N]) Aggregate(value N, attr attribute.Set) {
	s.set(value, attr)
}
