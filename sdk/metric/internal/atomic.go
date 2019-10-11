// Copyright 2019, OpenTelemetry Authors
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

package internal

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type (
	AtomicFloat64 struct {
		ptr *uint64
	}

	AtomicInt64 struct {
		ptr *uint64
	}
)

func NewAtomicFloat64(ptr *uint64) AtomicFloat64 {
	return AtomicFloat64{ptr: ptr}
}

func NewAtomicInt64(ptr *uint64) AtomicInt64 {
	return AtomicInt64{ptr: ptr}
}

func (a AtomicInt64) Add(value int64) {
	atomic.AddInt64((*int64)(unsafe.Pointer(a.ptr)), value)
}

func (a AtomicInt64) Load() int64 {
	return atomic.LoadInt64((*int64)(unsafe.Pointer(a.ptr)))
}

func (a AtomicInt64) Store(v int64) {
	atomic.StoreInt64((*int64)(unsafe.Pointer(a.ptr)), v)
}

func (a AtomicFloat64) Add(value float64) {
	for {
		u64 := atomic.LoadUint64(a.ptr)
		f64 := math.Float64frombits(u64)
		nu64 := math.Float64bits(f64 + value)

		if atomic.CompareAndSwapUint64(a.ptr, u64, nu64) {
			return
		}
	}
}

func (a AtomicFloat64) Load() float64 {
	u64 := atomic.LoadUint64(a.ptr)
	return math.Float64frombits(u64)
}

func (a AtomicFloat64) Store(x float64) {
	atomic.StoreUint64(a.ptr, math.Float64bits(x))
}
