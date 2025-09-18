// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync/atomic"
)

// atomicSum is an efficient way of adding to a number which is either an
// int64 or float64. It is designed to be efficient when adding whole
// numbers, regardless of whether N is an int64 or float64.
//
// Inspired by the Prometheus counter implementation:
// https://github.com/prometheus/client_golang/blob/14ccb93091c00f86b85af7753100aa372d63602b/prometheus/counter.go#L108
type atomicSum[N int64 | float64] struct {
	// nFloatBits contains only the non-integer portion of the counter.
	nFloatBits atomic.Uint64
	// nInt contains only the integer portion of the counter.
	nInt atomic.Int64
}

// load returns the current value. The caller must ensure all calls to add have
// returned prior to calling load.
func (n *atomicSum[N]) load() N {
	fval := math.Float64frombits(n.nFloatBits.Load())
	ival := n.nInt.Load()
	return N(fval + float64(ival))
}

func (n *atomicSum[N]) add(value N) {
	ival := int64(value)
	// This case is where the value is an int, or if it is a whole-numbered float.
	if float64(ival) == float64(value) {
		n.nInt.Add(ival)
		return
	}

	// Value must be a float below.
	for {
		oldBits := n.nFloatBits.Load()
		newBits := math.Float64bits(math.Float64frombits(oldBits) + float64(value))
		if n.nFloatBits.CompareAndSwap(oldBits, newBits) {
			return
		}
	}
}
