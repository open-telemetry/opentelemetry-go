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

package core

//go:generate stringer -type=NumberKind

import (
	"fmt"
	"math"
	"sync/atomic"
	"unsafe"
)

// NumberKind describes the data type of the Number.
type NumberKind int8

const (
	// Int64NumberKind means that the Number stores int64.
	Int64NumberKind NumberKind = iota
	// Float64NumberKind means that the Number stores float64.
	Float64NumberKind
)

// Number represents either an integral or a floating point value. It
// needs to be accompanied with a source of NumberKind that describes
// the actual type of the value stored within Number.
type Number uint64

// - constructors

// NewZeroNumber
func NewZeroNumber(kind NumberKind) Number {
	switch kind {
	case Int64NumberKind:
		return NewInt64Number(0)
	case Float64NumberKind:
		return NewFloat64Number(0.)
	}
	return Number(0)
}

// NewNumberFromRaw creates a new Number from a raw value.
func NewNumberFromRaw(r uint64) Number {
	return Number(r)
}

// NewInt64Number creates an integral Number.
func NewInt64Number(i int64) Number {
	return NewNumberFromRaw(int64ToRaw(i))
}

// NewFloat64Number creates a floating point Number.
func NewFloat64Number(f float64) Number {
	return NewNumberFromRaw(float64ToRaw(f))
}

// - as x

// AsNumber gets the raw, uninterpreted raw value. Might be useful for
// some atomic operations.
func (n Number) AsNumber() Number {
	return n
}

// AsRaw gets the raw, uninterpreted raw value. Might be useful for
// some atomic operations.
func (n Number) AsRaw() uint64 {
	return uint64(n)
}

// AsInt64 assumes that the value contains an int64 and returns it as
// such.
func (n Number) AsInt64() int64 {
	return rawToInt64(n.AsRaw())
}

// AsFloat64 assumes that the measurement value contains a float64 and
// returns it as such.
func (n Number) AsFloat64() float64 {
	return rawToFloat64(n.AsRaw())
}

// - as x atomic

// AsNumberAtomic gets the raw, uninterpreted raw value. Might be useful for
// some atomic operations.
func (n *Number) AsNumberAtomic() Number {
	return NewNumberFromRaw(n.AsRawAtomic())
}

// AsRawAtomic gets atomically the raw, uninterpreted raw value. Might
// be useful for some atomic operations.
func (n *Number) AsRawAtomic() uint64 {
	return atomic.LoadUint64(n.AsRawPtr())
}

// AsInt64Atomic assumes that the number contains an int64 and
// atomically returns it as such.
func (n *Number) AsInt64Atomic() int64 {
	return atomic.LoadInt64(n.AsInt64Ptr())
}

// AsFloat64 assumes that the measurement value contains a float64 and
// returns it as such.
func (n *Number) AsFloat64Atomic() float64 {
	return rawToFloat64(n.AsRawAtomic())
}

// - as x ptr

// AsRawPtr gets the pointer to the raw, uninterpreted raw
// value. Might be useful for some atomic operations.
func (n *Number) AsRawPtr() *uint64 {
	return (*uint64)(n)
}

func (n *Number) AsInt64Ptr() *int64 {
	return rawPtrToInt64Ptr(n.AsRawPtr())
}

func (n *Number) AsFloat64Ptr() *float64 {
	return rawPtrToFloat64Ptr(n.AsRawPtr())
}

// - set

func (n *Number) SetNumber(nn Number) {
	*n.AsRawPtr() = nn.AsRaw()
}

func (n *Number) SetRaw(r uint64) {
	*n.AsRawPtr() = r
}

func (n *Number) SetInt64(i int64) {
	*n.AsInt64Ptr() = i
}

func (n *Number) SetFloat64(f float64) {
	*n.AsFloat64Ptr() = f
}

// - set atomic

func (n *Number) SetNumberAtomic(nn Number) {
	atomic.StoreUint64(n.AsRawPtr(), nn.AsRaw())
}

func (n *Number) SetRawAtomic(r uint64) {
	atomic.StoreUint64(n.AsRawPtr(), r)
}

func (n *Number) SetInt64Atomic(i int64) {
	atomic.StoreInt64(n.AsInt64Ptr(), i)
}

func (n *Number) SetFloat64Atomic(f float64) {
	atomic.StoreUint64(n.AsRawPtr(), float64ToRaw(f))
}

// - swap

func (n *Number) SwapNumber(nn Number) Number {
	old := *n
	n.SetNumber(nn)
	return old
}

func (n *Number) SwapRaw(r uint64) uint64 {
	old := n.AsRaw()
	n.SetRaw(r)
	return old
}

func (n *Number) SwapInt64(i int64) int64 {
	old := n.AsInt64()
	n.SetInt64(i)
	return old
}

func (n *Number) SwapFloat64(f float64) float64 {
	old := n.AsFloat64()
	n.SetFloat64(f)
	return old
}

// - swap atomic

func (n *Number) SwapNumberAtomic(nn Number) Number {
	return NewNumberFromRaw(atomic.SwapUint64(n.AsRawPtr(), nn.AsRaw()))
}

func (n *Number) SwapRawAtomic(r uint64) uint64 {
	return atomic.SwapUint64(n.AsRawPtr(), r)
}

func (n *Number) SwapInt64Atomic(i int64) int64 {
	return atomic.SwapInt64(n.AsInt64Ptr(), i)
}

func (n *Number) SwapFloat64Atomic(f float64) float64 {
	return rawToFloat64(atomic.SwapUint64(n.AsRawPtr(), float64ToRaw(f)))
}

// - add

func (n *Number) AddNumber(kind NumberKind, nn Number) {
	switch kind {
	case Int64NumberKind:
		n.AddInt64(nn.AsInt64())
	case Float64NumberKind:
		n.AddFloat64(nn.AsFloat64())
	}
}

func (n *Number) AddRaw(kind NumberKind, r uint64) {
	n.AddNumber(kind, NewNumberFromRaw(r))
}

func (n *Number) AddInt64(i int64) {
	*n.AsInt64Ptr() += i
}

func (n *Number) AddFloat64(f float64) {
	*n.AsFloat64Ptr() += f
}

// - add atomic

func (n *Number) AddNumberAtomic(kind NumberKind, nn Number) {
	switch kind {
	case Int64NumberKind:
		n.AddInt64Atomic(nn.AsInt64())
	case Float64NumberKind:
		n.AddFloat64Atomic(nn.AsFloat64())
	}
}

func (n *Number) AddRawAtomic(kind NumberKind, r uint64) {
	n.AddNumberAtomic(kind, NewNumberFromRaw(r))
}

func (n *Number) AddInt64Atomic(i int64) {
	atomic.AddInt64(n.AsInt64Ptr(), i)
}

func (n *Number) AddFloat64Atomic(f float64) {
	for {
		o := n.AsFloat64Atomic()
		if n.CompareAndSwapFloat64(o, o+f) {
			break
		}
	}
}

// - compare and swap (atomic only)

func (n *Number) CompareAndSwapNumber(on, nn Number) bool {
	return atomic.CompareAndSwapUint64(n.AsRawPtr(), on.AsRaw(), nn.AsRaw())
}

func (n *Number) CompareAndSwapRaw(or, nr uint64) bool {
	return atomic.CompareAndSwapUint64(n.AsRawPtr(), or, nr)
}

func (n *Number) CompareAndSwapInt64(oi, ni int64) bool {
	return atomic.CompareAndSwapInt64(n.AsInt64Ptr(), oi, ni)
}

func (n *Number) CompareAndSwapFloat64(of, nf float64) bool {
	return atomic.CompareAndSwapUint64(n.AsRawPtr(), float64ToRaw(of), float64ToRaw(nf))
}

// - compare

func (n Number) CompareNumber(kind NumberKind, nn Number) int {
	switch kind {
	case Int64NumberKind:
		return n.CompareInt64(nn.AsInt64())
	case Float64NumberKind:
		return n.CompareFloat64(nn.AsFloat64())
	default:
		// you get what you deserve
		return 0
	}
}

func (n Number) CompareRaw(kind NumberKind, r uint64) int {
	return n.CompareNumber(kind, NewNumberFromRaw(r))
}

// Float64Compare assumes that the Number contains a float64 and
// performs a comparison between the value and the other value. It
// returns the typical result of the compare function: -1 if the value
// is less than the other, 0 if both are equal, 1 if the value is
// greater than the other.
func (n Number) CompareFloat64(f float64) int {
	this := n.AsFloat64()
	if this < f {
		return -1
	} else if this > f {
		return 1
	}
	return 0
}

// Int64Compare assumes that the Number contains an int64 and performs
// a comparison between the value and the other value. It returns the
// typical result of the compare function: -1 if the value is less
// than the other, 0 if both are equal, 1 if the value is greater than
// the other.
func (n Number) CompareInt64(i int64) int {
	this := n.AsInt64()
	if this < i {
		return -1
	} else if this > i {
		return 1
	}
	return 0
}

// - relations to zero

// IsPositive returns true if the actual value is greater than zero.
func (n Number) IsPositive(kind NumberKind) bool {
	return n.compareWithZero(kind) > 0
}

// IsNegative returns true if the actual value is less than zero.
func (n Number) IsNegative(kind NumberKind) bool {
	return n.compareWithZero(kind) < 0
}

// IsZero returns true if the actual value is equal to zero.
func (n Number) IsZero(kind NumberKind) bool {
	return n.compareWithZero(kind) == 0
}

// - misc

// Emit returns a string representation of the raw value of the
// Number. A %d is used for integral values, %f for floating point
// values.
func (n Number) Emit(kind NumberKind) string {
	switch kind {
	case Int64NumberKind:
		return fmt.Sprintf("%d", n.AsInt64())
	case Float64NumberKind:
		return fmt.Sprintf("%f", n.AsFloat64())
	default:
		return ""
	}
}

// - private stuff

func (n Number) compareWithZero(kind NumberKind) int {
	switch kind {
	case Int64NumberKind:
		return n.CompareInt64(0)
	case Float64NumberKind:
		return n.CompareFloat64(0.)
	default:
		// you get what you deserve
		return 0
	}
}

func rawToFloat64(r uint64) float64 {
	return math.Float64frombits(r)
}

func float64ToRaw(f float64) uint64 {
	return math.Float64bits(f)
}

func rawToInt64(r uint64) int64 {
	return int64(r)
}

func int64ToRaw(i int64) uint64 {
	return uint64(i)
}

func rawPtrToFloat64Ptr(r *uint64) *float64 {
	return (*float64)(unsafe.Pointer(r))
}

func rawPtrToInt64Ptr(r *uint64) *int64 {
	return (*int64)(unsafe.Pointer(r))
}
