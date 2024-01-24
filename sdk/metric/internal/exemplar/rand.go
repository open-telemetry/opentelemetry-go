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
	"math"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// rng is used to make sampling decisions.
//
// Do not use crypto/rand. There is no reason for the decrease in performance
// given this is not a security sensitive decision.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// FixedSize returns a [Reservoir] that samples at most n exemplars. If there
// are n or less measurements made, the Reservoir will sample each one. If
// there are more than n, the Reservoir will then randomly sample all
// additional measurement with a decreasing probability.
func FixedSize[N int64 | float64](n int) Reservoir[N] {
	r := &randRes[N]{storage: newStorage[N](n)}
	r.reset()
	return r
}

type randRes[N int64 | float64] struct {
	*storage[N]

	// count is the number of measurement seen.
	count int64
	// next is the next count that will store a measurement at a random index
	// once the reservoir has been filled.
	next int64
	// w is the largest random number in a distribution that is used to compute
	// the next next.
	w float64
}

func (r *randRes[N]) Offer(ctx context.Context, t time.Time, n N, a []attribute.KeyValue) {
	// The following algorithm is "Algorithm L" from Li, Kim-Hung (4 December
	// 1994). "Reservoir-Sampling Algorithms of Time Complexity
	// O(n(1+log(N/n)))". ACM Transactions on Mathematical Software. 20 (4):
	// 481–493 (https://dl.acm.org/doi/10.1145/198429.198435).
	//
	// It is used because of its balance of simplicity and performance. In
	// particular it has an asymptotic runtime of O(k(1 + log(n/k)) where n is
	// the number of measurements offered and k is the reservoir size.
	//
	// See https://github.com/MrAlias/reservoir-sampling for a comparison of
	// reservoir sampling algorithms (including performance benchmarks).

	if int(r.count) < cap(r.store) {
		r.store[r.count] = newMeasurement(ctx, t, n, a)
	} else {
		if r.count == r.next {
			idx := int(rng.Int63n(int64(cap(r.store))))
			r.store[idx] = newMeasurement(ctx, t, n, a)
			r.advance()
		}
	}
	r.count++
}

// reset resets r to the initial state.
func (r *randRes[N]) reset() {
	// This resets the number of exemplars known.
	r.count = 0
	// Random index inserts should only happen after the storage is full.
	r.next = int64(cap(r.store))

	// Initial random number in the series used to generate r.next.
	//
	// This is set before r.advance to reset or initialize the random number
	// series. Without doing so it would always be 0 or never restart a new
	// random number series.
	r.w = math.Exp(math.Log(rng.Float64()) / float64(cap(r.store)))

	// Calculate the next random index a measurement will become an exemplar.
	r.advance()
}

// advance updates the random number series of r and the next insert index.
func (r *randRes[N]) advance() {
	// Next value in the random number series based on the existing r.w.
	r.w *= math.Exp(math.Log(rng.Float64()) / float64(cap(r.store)))
	// Use the new random number in the series to calculate the index of the
	// next measurement that will be stored.
	r.next += int64(math.Log(rng.Float64())/math.Log(1-r.w)) + 1
}

func (r *randRes[N]) Collect(dest *[]metricdata.Exemplar[N]) {
	r.storage.Collect(dest)
	// Call reset here even though it will reset r.count and restart the random
	// number series. This will persist any old exemplars as long as no new
	// measurements are offered, but it will also prioritize those new
	// measurements that are made over the older collection cycle ones.
	r.reset()
}

func (r *randRes[N]) Flush(dest *[]metricdata.Exemplar[N]) {
	r.storage.Flush(dest)
	r.reset()
}
