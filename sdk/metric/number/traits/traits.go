package traits

import (
	"math"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/sdk/metric/number"
)

type Int64 struct{}

func (Int64) ToNumber(x int64) number.Number {
	return number.Number(x)
}

func (Int64) FromNumber(n number.Number) int64 {
	return int64(n)
}

func (Int64) SetAtomic(ptr *int64, value int64) {
	atomic.StoreInt64(ptr, value)
}

func (Int64) SwapAtomic(ptr *int64, value int64) int64 {
	return atomic.SwapInt64(ptr, value)
}

func (Int64) AddAtomic(ptr *int64, value int64) {
	// @@@
	panic("here")
}

func (Int64) IsNaN(_ int64) bool {
	return false
}

type Float64 struct{}

func (Float64) ToNumber(x float64) number.Number {
	return number.Number(math.Float64bits(x))
}

func (Float64) FromNumber(n number.Number) float64 {
	return math.Float64frombits(uint64(n))
}

func (Float64) SetAtomic(ptr *float64, value float64) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(ptr)), math.Float64bits(value))
}

func (Float64) SwapAtomic(ptr *float64, value float64) float64 {
	return math.Float64frombits(atomic.SwapUint64((*uint64)(unsafe.Pointer(ptr)), math.Float64bits(value)))
}

func (Float64) AddAtomic(ptr *float64, value float64) {
	// @@@
	panic("here")
}

func (Float64) IsNaN(value float64) bool {
	return math.IsNaN(value)
}

type Any[N int64|float64] interface {
	FromNumber(number.Number) N
	ToNumber(value N) number.Number
	SetAtomic(ptr *N, value N)
	AddAtomic(ptr *N, value N)
	SwapAtomic(ptr *N, value N) N
	IsNaN(value N) bool
}
