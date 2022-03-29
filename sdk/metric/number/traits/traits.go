package traits

import (
	"math"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/sdk/metric/number"
)

// Any is the generic traits interface for numbers used in sdk/metric.
// Two concrete implementations are given with this interface for
// int64 and float64.
type Any[N int64 | float64] interface {
	// FromNumber turns a generic 64bits into the correct machine type.
	FromNumber(number.Number) N

	// ToNumber turns this type into a generic 64bit value.
	ToNumber(value N) number.Number

	// SetAtomic sets `ptr` to `value`.
	SetAtomic(ptr *N, value N)

	// AddAtomic sets `ptr` to `value+*ptr`.
	AddAtomic(ptr *N, value N)

	// AddAtomic sets `ptr` to `value` and returns the former value.
	SwapAtomic(ptr *N, value N) N

	// IsNaN indicates whether `math.IsNaN()` is true (impossible for int64).
	IsNaN(value N) bool

	// IsInf indicates whether `math.IsInf()` is true (impossible for int64).
	IsInf(value N) bool
}

// Int64 implements Traits[int64].
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
	atomic.AddInt64(ptr, value)
}

func (Int64) IsNaN(_ int64) bool {
	return false
}

func (Int64) IsInf(_ int64) bool {
	return false
}

// Int64 implements Traits[float64].
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
	for {
		oldBits := atomic.LoadUint64((*uint64)(unsafe.Pointer(ptr)))
		sum := math.Float64frombits(oldBits) + value
		newBits := math.Float64bits(sum)

		if atomic.CompareAndSwapUint64((*uint64)(unsafe.Pointer(ptr)), oldBits, newBits) {
			return
		}
	}
}

func (Float64) IsNaN(value float64) bool {
	return math.IsNaN(value)
}

func (Float64) IsInf(value float64) bool {
	return math.IsInf(value, 0)
}
