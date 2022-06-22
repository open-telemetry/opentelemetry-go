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
)

// sumAgg summarizes a set of measurements as their arithmetic sum.
type sumAgg[N int64 | float64] struct {
	// zero value used for the base of all new sums.
	newFunc NewAtomicFunc[N]

	// map[attribute.Set]Atomic[N]
	current sync.Map
}

// NewSum returns an Aggregator that summarizes a set of measurements as their
// arithmetic sum. The zero value will be used as the start value for all new
// Aggregations.
func NewSum[N int64 | float64](f NewAtomicFunc[N]) Aggregator[N] {
	return &sumAgg[N]{newFunc: f}
}

func (s *sumAgg[N]) Aggregate(value N, attr *attribute.Set) {
	if v, ok := s.current.Load(*attr); ok {
		v.(Atomic[N]).Add(value)
		return
	}

	v, _ := s.current.LoadOrStore(*attr, s.newFunc())
	v.(Atomic[N]).Add(value)
}

func (s *sumAgg[N]) flush() []Aggregation {
	now := time.Now().UnixNano()
	var aggs []Aggregation
	s.current.Range(func(key, val any) bool {
		attrs := key.(attribute.Set)
		aggs = append(aggs, Aggregation{
			Timestamp:  now,
			Attributes: &attrs,
			Value:      SingleValue[N]{Value: val.(Atomic[N]).Load()},
		})

		// Reset.
		s.current.Delete(key)

		return true
	})

	return aggs
}
