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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// filter is an aggregator that applies attribute filter when Aggregating. filters
// do not have any backing memory, and must be constructed with a backing Aggregator.
type filter[N int64 | float64] struct {
	filter     func(attribute.Set) attribute.Set
	aggregator Aggregator[N]

	sync.Mutex
	seen map[attribute.Set]attribute.Set
}

// NewFilter wraps an Aggregator with an attribute filtering function.
func NewFilter[N int64 | float64](agg Aggregator[N], fn func(attribute.Set) attribute.Set) Aggregator[N] {
	if fn == nil {
		return agg
	}
	return &filter[N]{
		filter:     fn,
		aggregator: agg,
		seen:       map[attribute.Set]attribute.Set{},
	}
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (f *filter[N]) Aggregate(measurement N, attr attribute.Set) {
	// TODO (#3006): drop stale attributes from seen.
	f.Lock()
	defer f.Unlock()
	fAttr, ok := f.seen[attr]
	if !ok {
		fAttr = f.filter(attr)
		f.seen[attr] = fAttr
	}
	f.aggregator.Aggregate(measurement, fAttr)
}

// Aggregation returns an Aggregation, for all the aggregated
// measurements made and ends an aggregation cycle.
func (f *filter[N]) Aggregation() metricdata.Aggregation {
	return f.aggregator.Aggregation()
}
