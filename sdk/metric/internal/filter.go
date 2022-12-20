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
	filter     attribute.Filter
	aggregator Aggregator[N]

	// Used to aggreagte if an aggregator aggregates values differently for
	// spatically reaggregated attributes.
	aggregateFiltered func(N, attribute.Set)

	sync.Mutex
	seen     map[attribute.Set]attribute.Set
	filtered map[attribute.Set]bool
}

// NewFilter wraps an Aggregator with an attribute filtering function.
func NewFilter[N int64 | float64](agg Aggregator[N], fn attribute.Filter) Aggregator[N] {
	if fn == nil {
		return agg
	}
	af, ok := agg.(interface{ aggregateFiltered(N, attribute.Set) })
	if ok {
		return &filter[N]{
			filter:            fn,
			aggregator:        agg,
			aggregateFiltered: af.aggregateFiltered,
			seen:              make(map[attribute.Set]attribute.Set),
			// Use distinct filtered and seen to ensure un-filtered attributes
			// that match the same previously filtered attributes is treated
			// the same (added, not set).
			filtered: make(map[attribute.Set]bool),
		}
	}
	return &filter[N]{
		filter:     fn,
		aggregator: agg,
		seen:       make(map[attribute.Set]attribute.Set),
		// Don't allocate filtered as it won't be used
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
		var na []attribute.KeyValue
		fAttr, na = attr.Filter(f.filter)
		f.seen[attr] = fAttr
		if f.aggregateFiltered != nil && len(na) != 0 {
			f.filtered[fAttr] = true
		}
	}
	if f.aggregateFiltered != nil && f.filtered[fAttr] {
		f.aggregateFiltered(measurement, fAttr)
	} else {
		f.aggregator.Aggregate(measurement, fAttr)
	}
}

// Aggregation returns an Aggregation, for all the aggregated
// measurements made and ends an aggregation cycle.
func (f *filter[N]) Aggregation() metricdata.Aggregation {
	return f.aggregator.Aggregation()
}
