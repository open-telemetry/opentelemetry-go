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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// newFilter returns an Aggregator that wraps an agg with an attribute
// filtering function. Both pre-computed non-pre-computed Aggregators can be
// passed for agg. An appropriate Aggregator will be returned for the detected
// type.
func newFilter[N int64 | float64](agg aggregator[N], fn attribute.Filter) aggregator[N] {
	if fn == nil {
		return agg
	}
	return &filter[N]{
		filter:     fn,
		aggregator: agg,
	}
}

// filter wraps an aggregator with an attribute filter. All recorded
// measurements will have their attributes filtered before they are passed to
// the underlying aggregator's Aggregate method.
//
// This should not be used to wrap a pre-computed Aggregator. Use a
// precomputedFilter instead.
type filter[N int64 | float64] struct {
	filter     attribute.Filter
	aggregator aggregator[N]
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (f *filter[N]) Aggregate(measurement N, attr attribute.Set) {
	fAttr, _ := attr.Filter(f.filter)
	f.aggregator.Aggregate(measurement, fAttr)
}

// Aggregation returns an Aggregation, for all the aggregated
// measurements made and ends an aggregation cycle.
func (f *filter[N]) Aggregation() metricdata.Aggregation {
	return f.aggregator.Aggregation()
}
