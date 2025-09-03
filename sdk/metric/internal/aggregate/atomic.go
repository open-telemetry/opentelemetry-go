// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync/atomic"
)

// counter is an efficient way of adding to a number which is either an
// int64 or float64.
type counter[N int64 | float64] struct {
	// nFloatBits contains only the non-integer portion of the counter.
	nFloatBits uint64
	// nInt contains only the integer portion of the counter.
	nInt uint64
}

// value returns the float or integer value.
func (n *counter[N]) value() N {
	fval := math.Float64frombits(atomic.LoadUint64(&n.nFloatBits))
	ival := atomic.LoadUint64(&n.nInt)
	return N(fval + float64(ival))
}

func (n *counter[N]) add(value N) {
	ival := uint64(value)
	// This case is where the value is an int, or if it is a whole-numbered float.
	if float64(ival) == float64(value) {
		atomic.AddUint64(&n.nInt, ival)
		return
	}

	// Value must be a float below.
	for {
		oldBits := atomic.LoadUint64(&n.nFloatBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + float64(value))
		if atomic.CompareAndSwapUint64(&n.nFloatBits, oldBits, newBits) {
			return
		}
	}
}
