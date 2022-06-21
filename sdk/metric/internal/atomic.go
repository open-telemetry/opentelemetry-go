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

//go:build go1.18
// +build go1.18

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"math"
	"sync/atomic"
)

// Atomic provides atomic access to a generic value type.
type Atomic[N int64 | float64] interface {
	// Store value atomically.
	Store(value N)

	// Add value atomically.
	Add(value N)

	// Load returns the stored value.
	Load() N

	// Clone creates an independent copy of the current value.
	Clone() Atomic[N]
}

type Int64 struct {
	value *int64
}

var _ Atomic[int64] = Int64{}

func NewInt64(v int64) Int64 {
	return Int64{value: &v}
}

func (v Int64) Store(value int64) { atomic.StoreInt64(v.value, value) }
func (v Int64) Add(value int64)   { atomic.AddInt64(v.value, value) }
func (v Int64) Load() int64       { return atomic.LoadInt64(v.value) }
func (v Int64) Clone() Atomic[int64] {
	return NewInt64(v.Load())
}

type Float64 struct {
	value *uint64
}

var _ Atomic[float64] = Float64{}

func NewFloat64(v float64) Float64 {
	u := math.Float64bits(v)
	return Float64{value: &u}
}

func (v Float64) Store(value float64) {
	atomic.StoreUint64(v.value, math.Float64bits(value))
}

func (v Float64) Add(value float64) {
	for {
		old := atomic.LoadUint64(v.value)
		sum := math.Float64bits(math.Float64frombits(old) + value)
		if atomic.CompareAndSwapUint64(v.value, old, sum) {
			return
		}
	}
}

func (v Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(v.value))
}

func (v Float64) Clone() Atomic[float64] {
	return NewFloat64(v.Load())
}
