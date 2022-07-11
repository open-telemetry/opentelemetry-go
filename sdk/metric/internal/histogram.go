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
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

// histogram summarizes a set of measurements as an histogram with
// explicitly defined buckets.
type histogram[N int64 | float64] struct {
	// TODO(#2970): implement.
}

func (s *histogram[N]) Aggregate(value N, attr attribute.Set) {
	// TODO(#2970): implement.
}

// NewDeltaHistogram returns an Aggregator that summarizes a set of
// measurements as an histogram. Each histogram is scoped by attributes and
// the aggregation cycle the measurements were made in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregations method is called it will reset all histogram
// counts to zero.
func NewDeltaHistogram[N int64 | float64](cfg aggregation.ExplicitBucketHistogram) Aggregator[N] {
	return &deltaHistogram[N]{}
}

// deltaHistogram summarizes a set of measurements made in a single
// aggregation cycle as an histogram with explicitly defined buckets.
type deltaHistogram[N int64 | float64] struct {
	histogram[N]

	// TODO(#2970): implement.
}

func (s *deltaHistogram[N]) Aggregations() []export.Aggregation {
	// TODO(#2970): implement.
	return nil
}

// NewCumulativeHistogram returns an Aggregator that summarizes a set of
// measurements as an histogram. Each histogram is scoped by attributes.
//
// Each aggregation cycle builds from the previous, the histogram counts are
// the bucketed counts of all values aggregated since the returned Aggregator
// was created.
func NewCumulativeHistogram[N int64 | float64](cfg aggregation.ExplicitBucketHistogram) Aggregator[N] {
	return &cumulativeHistogram[N]{}
}

// cumulativeHistogram summarizes a set of measurements made over all
// aggregation cycles as an histogram with explicitly defined buckets.
type cumulativeHistogram[N int64 | float64] struct {
	histogram[N]

	// TODO(#2970): implement.
}

func (s *cumulativeHistogram[N]) Aggregations() []export.Aggregation {
	// TODO(#2970): implement.
	return nil
}
