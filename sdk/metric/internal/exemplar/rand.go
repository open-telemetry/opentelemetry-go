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
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// FixedSize returns a [Reservoir] that samples at most n exemplars. If there
// are n or less number of measurements made, the Reservoir will sample each
// one. If there are more than n number of measurements made, the Reservoir
// will then randomly sample all additional measurement with a decreasing
// probability.
func FixedSize[N int64 | float64](n int) Reservoir[N] {
	return &randRes[N]{fixedRes: newFixedRes[N](n)}
}

type randRes[N int64 | float64] struct {
	*fixedRes[N]

	// count is the number of measurement seen.
	count int64
}

func (r *randRes[N]) Offer(ctx context.Context, t time.Time, n N, a attribute.Set) {
	// TODO: fix overflow error.
	r.count++
	if int(r.count) <= cap(r.store) {
		r.store[r.count-1] = newMeasurement(ctx, t, n, a)
		return
	}

	j := int(rng.Int63n(r.count))
	if j < cap(r.store) {
		r.store[j] = newMeasurement(ctx, t, n, a)
	}
}
