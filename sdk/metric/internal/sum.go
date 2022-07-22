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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// sum summarizes a set of measurements as their arithmetic sum.
type sum[N int64 | float64] struct {
	// TODO(#2972): implement.
}

func (s *sum[N]) Aggregate(value N, attr attribute.Set) {
	// TODO(#2972): implement.
}

// NewDeltaSum returns an Aggregator that summarizes a set of measurements as
// their arithmetic sum. Each sum is scoped by attributes and the aggregation
// cycle the measurements were made in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregations method is called it will reset all sums to zero.
func NewDeltaSum[N int64 | float64]() Aggregator[N] {
	// TODO(#2972): implement.
	return &deltaSum[N]{}
}

// deltaSum summarizes a set of measurements made in a single aggregation
// cycle as their arithmetic sum.
type deltaSum[N int64 | float64] struct {
	sum[N]

	// TODO(#2972): implement.
}

func (s *deltaSum[N]) Aggregation() metricdata.Aggregation {
	// TODO(#2972): implement.
	return nil
}

// NewCumulativeSum returns an Aggregator that summarizes a set of
// measurements as their arithmetic sum. Each sum is scoped by attributes.
//
// Each aggregation cycle builds from the previous, the sums are the
// arithmetic sum of all values aggregated since the returned Aggregator was
// created.
func NewCumulativeSum[N int64 | float64]() Aggregator[N] {
	// TODO(#2972): implement.
	return &cumulativeSum[N]{}
}

// cumulativeSum summarizes a set of measurements made over all aggregation
// cycles as their arithmetic sum.
type cumulativeSum[N int64 | float64] struct {
	sum[N]

	// TODO(#2972): implement.
}

func (s *cumulativeSum[N]) Aggregation() metricdata.Aggregation {
	// TODO(#2972): implement.
	return nil
}
