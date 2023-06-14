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

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type fixedRes[N int64 | float64] struct {
	// store are the measurements sampled.
	//
	// This does not use []metricdata.Exemplar because it potentially would
	// require an allocation for trace and span IDs in the hot path of Offer.
	store []measurement[N]
}

func newFixedRes[N int64 | float64](n int) *fixedRes[N] {
	return &fixedRes[N]{store: make([]measurement[N], n)}
}

func (r *fixedRes[N]) Collect(dest *[]metricdata.Exemplar[N], attrs attribute.Set) {
	*dest = reset(*dest, len(r.store), len(r.store))
	var n int
	for _, m := range r.store {
		if m.Empty() {
			continue
		}

		m.Exemplar(&(*dest)[n], attrs)
		n++
	}
	*dest = (*dest)[:n]
}

func (r *fixedRes[N]) Flush(dest *[]metricdata.Exemplar[N], attrs attribute.Set) {
	*dest = reset(*dest, len(r.store), len(r.store))
	var n int
	for i, m := range r.store {
		if m.Empty() {
			continue
		}

		m.Exemplar(&(*dest)[n], attrs)
		n++

		// Reset.
		r.store[i] = measurement[N]{}
	}
	*dest = (*dest)[:n]
}
