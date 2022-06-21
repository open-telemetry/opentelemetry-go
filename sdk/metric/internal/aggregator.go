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

// Aggregator forms an aggregation from a collection of recorded measurements.
// Aggregators are use with Cyclers to collect and produce metrics from
// instrument measurements. Aggregators handle the collection (and
// aggregation) of measurements, while Cyclers handle how those aggregated
// measurements are combined and then produced to the telemetry pipeline.
type Aggregator[N int64 | float64] interface {
	// Aggregate records the measurement, scoped by attr, and aggregates it
	// into an aggregation.
	Aggregate(measurement N, attr *attribute.Set)

	// flush clears aggregations that have been recorded. The Aggregator
	// resets itself for a new aggregation period when called, it does not
	// carry forward any state. If aggregation periods need to be combined it
	// is the callers responsibility to achieve this.
	flush() []Aggregation
}
