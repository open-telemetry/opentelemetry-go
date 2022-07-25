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
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var errNegVal = errors.New("monotonic increasing sum: negative value")

// now is used to return the current local time while allowing tests to
// override the the default time.Now function.
var now = time.Now

// valueMap is the sum aggregator storage.
type valueMap[N int64 | float64] struct {
	sync.Mutex

	values map[attribute.Set]N
}

// newValueMap returns an instantiated valueMap.
func newValueMap[N int64 | float64]() valueMap[N] {
	return valueMap[N]{values: make(map[attribute.Set]N)}
}

func (v *valueMap[N]) add(value N, attr attribute.Set) {
	v.Lock()
	v.values[attr] += value
	v.Unlock()
}

// nonMonotonicSum summarizes a set of measurements as their arithmetic sum.
type nonMonotonicSum[N int64 | float64] struct {
	valueMap[N]
}

func (s *nonMonotonicSum[N]) Aggregate(value N, attr attribute.Set) {
	s.add(value, attr)
}

// monotonicSum summarizes a set of monotonically increasing measurements as
// their arithmetic sum.
type monotonicSum[N int64 | float64] struct {
	valueMap[N]
}

func (s *monotonicSum[N]) Aggregate(value N, attr attribute.Set) {
	if value < 0 {
		otel.Handle(fmt.Errorf("%w: %v", errNegVal, value))
	}
	s.add(value, attr)
}

// deltaSum summarizes a set of measurements made in a single aggregation
// cycle as their arithmetic sum.
type deltaSum[N int64 | float64] struct {
	valueMap[N]

	start time.Time
}

func (s *deltaSum[N]) dataPoints() []metricdata.DataPoint[N] {
	s.Lock()
	defer s.Unlock()

	if len(s.values) == 0 {
		return nil
	}

	t := now()

	data := make([]metricdata.DataPoint[N], 0, len(s.values))
	for attr, value := range s.values {
		data = append(data, metricdata.DataPoint[N]{
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
	return data
}

// cumulativeSum summarizes a set of measurements made over all aggregation
// cycles as their arithmetic sum.
type cumulativeSum[N int64 | float64] struct {
	valueMap[N]

	start time.Time
}

func (s *cumulativeSum[N]) dataPoints() []metricdata.DataPoint[N] {
	s.Lock()
	defer s.Unlock()

	if len(s.values) == 0 {
		return nil
	}

	t := now()

	data := make([]metricdata.DataPoint[N], 0, len(s.values))
	for attr, value := range s.values {
		data = append(data, metricdata.DataPoint[N]{
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

	return data
}

func NewNonMonotonicDeltaSum[N int64 | float64]() Aggregator[N] {
	v := newValueMap[N]()
	return &nonMonotonicDeltaSum[N]{
		nonMonotonicSum: nonMonotonicSum[N]{v},
		deltaSum: deltaSum[N]{
			valueMap: v,
			start:    now(),
		},
	}
}

type nonMonotonicDeltaSum[N int64 | float64] struct {
	nonMonotonicSum[N]
	deltaSum[N]
}

func (s *nonMonotonicDeltaSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: false,
		DataPoints:  s.deltaSum.dataPoints(),
	}
}

// NewMonotonicDeltaSum returns an Aggregator that summarizes a set of
// monotonically increasing measurements as their arithmetic sum. Each sum is
// scoped by attributes and the aggregation cycle the measurements were made
// in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewMonotonicDeltaSum[N int64 | float64]() Aggregator[N] {
	v := newValueMap[N]()
	return &monotonicDeltaSum[N]{
		monotonicSum: monotonicSum[N]{v},
		deltaSum: deltaSum[N]{
			valueMap: v,
			start:    now(),
		},
	}
}

type monotonicDeltaSum[N int64 | float64] struct {
	monotonicSum[N]
	deltaSum[N]
}

func (s *monotonicDeltaSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: true,
		DataPoints:  s.deltaSum.dataPoints(),
	}
}

func NewNonMonotonicCumulativeSum[N int64 | float64]() Aggregator[N] {
	v := newValueMap[N]()
	return &nonMonotonicCumulativeSum[N]{
		nonMonotonicSum: nonMonotonicSum[N]{v},
		cumulativeSum: cumulativeSum[N]{
			valueMap: v,
			start:    now(),
		},
	}
}

type nonMonotonicCumulativeSum[N int64 | float64] struct {
	nonMonotonicSum[N]
	cumulativeSum[N]
}

func (s *nonMonotonicCumulativeSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  s.cumulativeSum.dataPoints(),
	}
}

// NewMonotonicCumulativeSum returns an Aggregator that summarizes a set of
// monotonically increasing measurements as their arithmetic sum. Each sum is
// scoped by attributes and the aggregation cycle the measurements were made
// in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewMonotonicCumulativeSum[N int64 | float64]() Aggregator[N] {
	v := newValueMap[N]()
	return &monotonicCumulativeSum[N]{
		monotonicSum: monotonicSum[N]{v},
		cumulativeSum: cumulativeSum[N]{
			valueMap: v,
			start:    now(),
		},
	}
}

type monotonicCumulativeSum[N int64 | float64] struct {
	monotonicSum[N]
	cumulativeSum[N]
}

func (s *monotonicCumulativeSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  s.cumulativeSum.dataPoints(),
	}
}
