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

//go:build go1.17
// +build go1.17

package aggregation_test

import "go.opentelemetry.io/otel/sdk/metric/aggregation"

func ExampleAggregation() {
	// An aggregation that drops measurements.
	_ = aggregation.Aggregation{Operation: aggregation.Drop{}}

	// An aggregation that sums measurements.
	_ = aggregation.Aggregation{Operation: aggregation.Sum{}}

	// An aggregation that uses the last value for measurements.
	_ = aggregation.Aggregation{Operation: aggregation.LastValue{}}

	// An aggregation that bins measurements in a histogram.
	_ = aggregation.Aggregation{Operation: aggregation.ExplicitBucketHistogram{
		Boundaries:   []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
		RecordMinMax: true,
	}}
}
