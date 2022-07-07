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
	mu     sync.Mutex
	values map[attribute.Set]N
}

// NewSum returns an Aggregator that summarizes a set of measurements as their
// arithmetic sum. The zero value will be used as the start value for all new
// Aggregations.
func NewSum[N int64 | float64]() Aggregator[N] {
	return &sumAgg[N]{
		values: map[attribute.Set]N{},
	}
}

func (s *sumAgg[N]) Aggregate(value N, attr *attribute.Set) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[*attr] += value
}

func (s *sumAgg[N]) flush() []Aggregation {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	aggs := make([]Aggregation, 0, len(s.values))

	for attr, value := range s.values {
		attr := attr
		aggs = append(aggs, Aggregation{
			Timestamp:  now,
			Attributes: &attr,
			Value:      SingleValue[N]{Value: value},
		})
		delete(s.values, attr)
	}
	return aggs
}
