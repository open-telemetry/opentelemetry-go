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

import "go.opentelemetry.io/otel/attribute"

// sumAgg summarizes a set of measurements as their arithmetic sum.
type sumAgg[N int64 | float64] struct {
	// TODO(#2972): implement.
}

// NewSum returns an Aggregator that summarizes a set of
// measurements as their arithmetic sum. The zero value will be used as the
// start value for all new Aggregations.
func NewSum[N int64 | float64]() Aggregator[N] {
	// TODO(#2972): implement.
	return &sumAgg[N]{}
}

func (s *sumAgg[N]) Aggregate(value N, attr *attribute.Set) {
	// TODO(#2972): implement.
}

func (s *sumAgg[N]) flush() []Aggregation {
	// TODO(#2972): implement.
	return []Aggregation{
		{
			Value: SingleValue[N]{ /* TODO(#2972): calculate */ },
		},
	}
}
