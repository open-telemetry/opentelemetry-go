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
)

// histogramAgg summarizes a set of measurements as an histogram with
// explicitly defined buckets.
type histogramAgg[N int64 | float64] struct {
	// TODO: implement.
}

// NewHistogram returns an Aggregator that summarizes a set of measurements as
// an histogram. The zero value will be used as the start value for all the
// buckets of new Aggregations.
func NewHistogram[N int64 | float64](zero Number[N], cfg aggregation.ExplicitBucketHistogram) Aggregator[N] {
	return &histogramAgg[N]{}
}

func (s *histogramAgg[N]) Record(value N, attr *attribute.Set) {
	// TODO: implement.
}

func (s *histogramAgg[N]) Aggregate() []Aggregation {
	// TODO: implement.
	return []Aggregation{
		{
			Value: HistogramValue{ /* TODO: calculate. */ },
		},
	}
}
