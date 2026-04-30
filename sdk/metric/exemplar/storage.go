// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/exemplar"

import (
	"context"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// storage is an exemplar storage for [Reservoir] implementations.
type storage struct {
	// measurements are the measurements sampled.
	//
	// This does not use []metricdata.Exemplar because it potentially would
	// require an allocation for trace and span IDs in the hot path of Offer.
	measurements []measurement
}

func newStorage(n int) *storage {
	s := &storage{measurements: make([]measurement, n)}
	s.reset()
	return s
}

// Collect returns all the held exemplars.
//
// The Reservoir state is preserved after this call.
func (r *storage) Collect(dest *[]Exemplar) {
	*dest = reset(*dest, len(r.measurements), len(r.measurements))
	var n int
	for i := range r.measurements {
		if r.measurements[i].exemplar(&(*dest)[n]) {
			n++
		}
	}
	*dest = (*dest)[:n]
}

// reset resets the Algorithm L sampling state for all buckets in the storage.
func (r *storage) reset() {
	for i := range r.measurements {
		r.measurements[i].mux.Lock()
		r.measurements[i].reset()
		r.measurements[i].mux.Unlock()
	}
}

// measurement is a measurement made by a telemetry system.
type measurement struct {
	mux sync.Mutex
	// FilteredAttributes are the attributes dropped during the measurement.
	FilteredAttributes []attribute.KeyValue
	// Time is the time when the measurement was made.
	Time time.Time
	// Value is the value of the measurement.
	Value Value
	// Ctx is the context active when a measurement was made.
	Ctx context.Context

	valid bool

	// Algorithm L state
	// count is the number of measurements offered to this bucket.
	count int64
	// next is the next count that will store a measurement after the first.
	next int64
	// w is the largest random number in a distribution that is used to compute
	// the next next.
	w float64
}

// exemplar returns m as an [Exemplar].
// returns true if it populated the exemplar.
func (m *measurement) exemplar(dest *Exemplar) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	if !m.valid {
		return false
	}

	dest.FilteredAttributes = m.FilteredAttributes
	dest.Time = m.Time
	dest.Value = m.Value

	sc := trace.SpanContextFromContext(m.Ctx)
	if sc.HasTraceID() {
		traceID := sc.TraceID()
		dest.TraceID = traceID[:]
	} else {
		dest.TraceID = dest.TraceID[:0]
	}

	if sc.HasSpanID() {
		spanID := sc.SpanID()
		dest.SpanID = spanID[:]
	} else {
		dest.SpanID = dest.SpanID[:0]
	}
	m.reset()
	return true
}

// The following algorithm is "Algorithm L" from Li, Kim-Hung (4 December
// 1994). "Reservoir-Sampling Algorithms of Time Complexity
// O(n(1+log(N/n)))". ACM Transactions on Mathematical Software. 20 (4):
// 481–493 (https://dl.acm.org/doi/10.1145/198429.198435).
//
// A high-level overview of "Algorithm L":
//  0. Pre-calculate the random count greater than the storage size when
//     an exemplar will be replaced.
//  1. Accept all measurements offered until the configured storage size is
//     reached.
//  2. Loop:
//     a) When the pre-calculate count is reached, replace a random
//     existing exemplar with the offered measurement.
//     b) Calculate the next random count greater than the existing one
//     which will replace another exemplars
//
// The way a "replacement" count is computed is by looking at `n` number of
// independent random numbers each corresponding to an offered measurement.
// Of these numbers the smallest `k` (the same size as the storage
// capacity) of them are kept as a subset. The maximum value in this
// subset, called `w` is used to weight another random number generation
// for the next count that will be considered.
//
// By weighting the next count computation like described, it is able to
// perform a uniformly-weighted sampling algorithm based on the number of
// samples the reservoir has seen so far. The sampling will "slow down" as
// more and more samples are offered so as to reduce a bias towards those
// offered just prior to the end of the collection.
//
// This algorithm is preferred because of its balance of simplicity and
// performance. It will compute three random numbers (the bulk of
// computation time) for each item that becomes part of the reservoir, but
// it does not spend any time on items that do not. In particular it has an
// asymptotic runtime of O(k(1 + log(n/k)) where n is the number of
// measurements offered and k is the reservoir size.
//
// See https://en.wikipedia.org/wiki/Reservoir_sampling for an overview of
// this and other reservoir sampling algorithms. See
// https://github.com/MrAlias/reservoir-sampling for a performance
// comparison of reservoir sampling algorithms.
func (m *measurement) offer(ctx context.Context, ts time.Time, v Value, droppedAttr []attribute.KeyValue) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.count == m.next {
		// Overwrite
		m.FilteredAttributes = droppedAttr
		m.Time = ts
		m.Value = v
		m.Ctx = ctx
		m.valid = true

		m.advance()
	}
	m.count++
}

// reset resets the Algorithm L sampling state of m. It does not clear the
// stored exemplar data. The caller must hold the lock.
func (m *measurement) reset() {
	// This resets the number of exemplars known.
	m.count = 0
	// The first offer always is sampled.
	m.next = 0
	// The advance at the first offer will set the initial random number.
	m.w = 1.0
}

// advance updates the count at which the offered measurement will overwrite an
// existing exemplar.
func (m *measurement) advance() {
	m.w *= randomFloat64()
	// Use the new random number in the series to calculate the count of the
	// next measurement that will be stored.
	//
	// Given 0 < m.w < 1, each iteration will result in subsequent m.w being
	// smaller. This translates here into the next next being selected against
	// a distribution with a higher mean (i.e. the expected value will increase
	// and replacements become less likely)
	//
	// Important to note, the new m.next will always be at least 1 more than
	// the last m.next.
	m.next += int64(math.Log(randomFloat64())/math.Log(1-m.w)) + 1
}

// randomFloat64 returns, as a float64, a uniform pseudo-random number in the
// open interval (0.0,1.0).
func randomFloat64() float64 {
	// TODO: Use an algorithm that avoids rejection sampling. For example:
	//
	//   const precision = 1 << 53 // 2^53
	//   // Generate an integer in [1, 2^53 - 1]
	//   v := rand.Uint64() % (precision - 1) + 1
	//   return float64(v) / float64(precision)
	f := rand.Float64()
	for f == 0 {
		f = rand.Float64()
	}
	return f
}

func reset[T any](s []T, length, capacity int) []T {
	if cap(s) < capacity {
		return make([]T, length, capacity)
	}
	return s[:length]
}
