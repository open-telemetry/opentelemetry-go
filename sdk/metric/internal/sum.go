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

// valueMap is the aggregator storage for all sums.
type valueMap[N int64 | float64] struct {
	sync.Mutex

	values map[attribute.Set]N
}

// newValueMap returns an instantiated valueMap.
func newValueMap[N int64 | float64]() *valueMap[N] {
	return &valueMap[N]{values: make(map[attribute.Set]N)}
}

func (v *valueMap[N]) add(value N, attr attribute.Set) {
	v.Lock()
	v.values[attr] += value
	v.Unlock()
}

// nonMonotonicSum summarizes a set of measurements as their arithmetic sum.
type nonMonotonicSum[N int64 | float64] struct {
	*valueMap[N]
}

func (s *nonMonotonicSum[N]) Aggregate(value N, attr attribute.Set) {
	s.add(value, attr)
}

// monotonicSum summarizes a set of monotonically increasing measurements as
// their arithmetic sum.
type monotonicSum[N int64 | float64] struct {
	*valueMap[N]
}

var errNegVal = errors.New("monotonic increasing sum: negative value")

func (s *monotonicSum[N]) Aggregate(value N, attr attribute.Set) {
	if value < 0 {
		otel.Handle(fmt.Errorf("%w: %v", errNegVal, value))
	}
	s.add(value, attr)
}

// deltaSum summarizes a set of measurements made in a single aggregation
// cycle as their arithmetic sum.
type deltaSum[N int64 | float64] struct {
	*valueMap[N]

	monotonic bool
	start     time.Time
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

func (s *deltaSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: s.monotonic,
		DataPoints:  s.dataPoints(),
	}
}

// cumulativeSum summarizes a set of measurements made over all aggregation
// cycles as their arithmetic sum.
type cumulativeSum[N int64 | float64] struct {
	*valueMap[N]

	monotonic bool
	start     time.Time
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

func (s *cumulativeSum[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[N]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: s.monotonic,
		DataPoints:  s.dataPoints(),
	}
}

// NewDeltaSum returns an Aggregator that summarizes a set of measurements as
// their arithmetic sum. Each sum is scoped by attributes and the aggregation
// cycle the measurements were made in.
//
// If monotonic is true, the returned Aggregator will only accept increments
// zero or greater. Otherwise, an error is sent to the OTel ErrorHandler an
// no value is aggregated.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewDeltaSum[N int64 | float64](monotonic bool) Aggregator[N] {
	v := newValueMap[N]()
	ds := deltaSum[N]{
		valueMap:  v,
		monotonic: monotonic,
		start:     now(),
	}
	if monotonic {
		return &monotonicDeltaSum[N]{
			monotonicSum: monotonicSum[N]{v},
			deltaSum:     ds,
		}
	}
	return &nonMonotonicDeltaSum[N]{
		nonMonotonicSum: nonMonotonicSum[N]{v},
		deltaSum:        ds,
	}
}

type nonMonotonicDeltaSum[N int64 | float64] struct {
	nonMonotonicSum[N]
	deltaSum[N]
}

type monotonicDeltaSum[N int64 | float64] struct {
	monotonicSum[N]
	deltaSum[N]
}

// NewCumulativeSum returns an Aggregator that summarizes a set of
// measurements as their arithmetic sum. Each sum is scoped by attributes and
// the aggregation cycle the measurements were made in.
//
// If monotonic is true, the returned Aggregator will only accept increments
// zero or greater. Otherwise, an error is sent to the OTel ErrorHandler an
// no value is aggregated.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregation method is called it will reset all sums to zero.
func NewCumulativeSum[N int64 | float64](monotonic bool) Aggregator[N] {
	v := newValueMap[N]()
	cs := cumulativeSum[N]{
		valueMap:  v,
		monotonic: monotonic,
		start:     now(),
	}
	if monotonic {
		return &monotonicCumulativeSum[N]{
			monotonicSum:  monotonicSum[N]{v},
			cumulativeSum: cs,
		}
	}
	return &nonMonotonicCumulativeSum[N]{
		nonMonotonicSum: nonMonotonicSum[N]{v},
		cumulativeSum:   cs,
	}
}

type nonMonotonicCumulativeSum[N int64 | float64] struct {
	nonMonotonicSum[N]
	cumulativeSum[N]
}

type monotonicCumulativeSum[N int64 | float64] struct {
	monotonicSum[N]
	cumulativeSum[N]
}
