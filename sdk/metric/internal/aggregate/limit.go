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

import "go.opentelemetry.io/otel/attribute"

// overflowSet is the attribute set used to record a measurement when adding
// another distinct attribute set to the aggregate would exceed the aggregate
// limit.
var overflowSet = attribute.NewSet(attribute.Bool("otel.metric.overflow", true))

// limtAttr checks if adding a measurement for a will exceed the limit of the
// already measured values in m. If it will, overflowSet is returned.
// Otherwise, if it will not exceed the limit, or the limit is not set (limit
// <= 0), a is returned.
func limitAttr[V any](a attribute.Set, m map[attribute.Set]V, limit int) attribute.Set {
	if limit > 0 {
		_, exists := m[a]
		if !exists && len(m) >= limit-1 {
			return overflowSet
		}
	}

	return a
}
