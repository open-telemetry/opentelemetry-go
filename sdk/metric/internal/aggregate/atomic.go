// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync/atomic"
)

// atomicSum is an efficient way of adding to a number which is either an
// int64 or float64.
type atomicSum[N int64 | float64] struct {
	// nFloatBits contains only the non-integer portion of the counter.
	nFloatBits atomic.Uint64
	// nInt contains only the integer portion of the counter.
	nInt atomic.Int64
}

// load returns the float or integer value.
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

type atomicIntOrFloat[N int64 | float64] struct {
	// nFloatBits contains the float bits if N is float64.
	nFloatBits atomic.Uint64
	// nInt contains the int64 if N is int64
	nInt atomic.Int64
}

func (n *atomicIntOrFloat[N]) store(value N) {
	switch v := any(value).(type) {
	case int64:
		n.nInt.Store(v)
	case float64:
		n.nFloatBits.Store(math.Float64bits(v))
	}
}

func (n *atomicIntOrFloat[N]) load() (value N) {
	switch any(value).(type) {
	case int64:
		value = N(n.nInt.Load())
	case float64:
		value = N(math.Float64frombits(n.nFloatBits.Load()))
	}
	return
}

func (n *atomicIntOrFloat[N]) compareAndSwap(oldVal, newVal N) bool {
	switch any(oldVal).(type) {
	case float64:
		return n.nFloatBits.CompareAndSwap(math.Float64bits(float64(oldVal)), math.Float64bits(float64(newVal)))
	default:
		return n.nInt.CompareAndSwap(int64(oldVal), int64(newVal))
	}
}

type atomicMinMax[N int64 | float64] struct {
	min   atomicIntOrFloat[N]
	max   atomicIntOrFloat[N]
	isSet atomic.Bool
}

func (n *atomicMinMax[N]) observe(value N) {
	for {
		minLoaded := n.min.load()
		if (!n.isSet.Load() || value < minLoaded) && !n.min.compareAndSwap(minLoaded, value) {
			// We got a new min value, but lost the race. Try again.
			continue
		}
		maxLoaded := n.max.load()
		if (!n.isSet.Load() || value > maxLoaded) && !n.max.compareAndSwap(maxLoaded, value) {
			// We got a new max value, but lost the race. Try again.
			continue
		}
		break
	}
	n.isSet.Store(true)
}

func (n *atomicMinMax[N]) loadMin() (value N) {
	return n.min.load()
}

func (n *atomicMinMax[N]) loadMax() (value N) {
	return n.max.load()
}
