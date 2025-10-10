// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicSumAddFloatConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	var aSum atomicCounter[float64]
	for _, in := range []float64{
		0.2,
		0.25,
		1.6,
		10.55,
		42.4,
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			aSum.add(in)
		}()
	}
	wg.Wait()
	assert.Equal(t, float64(55), math.Round(aSum.load()))
}

func TestAtomicSumAddIntConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	var aSum atomicCounter[int64]
	for _, in := range []int64{
		1,
		2,
		3,
		4,
		5,
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			aSum.add(in)
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(15), aSum.load())
}

func BenchmarkAtomicCounter(b *testing.B) {
	b.Run("Int64", benchmarkAtomicCounter[int64])
	b.Run("Float64", benchmarkAtomicCounter[float64])
}

func benchmarkAtomicCounter[N int64 | float64](b *testing.B) {
	b.Run("add", func(b *testing.B) {
		var a atomicCounter[N]
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				a.add(2)
			}
		})
	})
	b.Run("load", func(b *testing.B) {
		var a atomicCounter[N]
		a.add(2)
		var v N
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v = a.load()
			}
		})
		assert.Equal(b, N(2), v)
	})
}

func TestHotColdWaitGroupConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	hcwg := &hotColdWaitGroup{}
	var data [2]uint64
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hotIdx := hcwg.start()
			defer hcwg.done(hotIdx)
			atomic.AddUint64(&data[hotIdx], 1)
		}()
	}
	for range 2 {
		readIdx := hcwg.swapHotAndWait()
		assert.NotPanics(t, func() {
			// reading without using atomics should not panic since we are
			// reading from the cold element, and have waited for all writes to
			// finish.
			t.Logf("read value %+v", data[readIdx])
		})
	}
	wg.Wait()
}

func TestAtomicN(t *testing.T) {
	t.Run("Int64", testAtomicN[int64])
	t.Run("Float64", testAtomicN[float64])

}

func testAtomicN[N int64 | float64](t *testing.T) {
	var v atomicN[N]
	assert.Equal(t, N(0), v.Load())
	assert.True(t, v.CompareAndSwap(0, 6))
	assert.Equal(t, N(6), v.Load())
	assert.False(t, v.CompareAndSwap(0, 6))
	v.Store(22)
	assert.Equal(t, N(22), v.Load())
}

func TestAtomicNConcurrentSafe(t *testing.T) {
	t.Run("Int64", testAtomicNConcurrentSafe[int64])
	t.Run("Float64", testAtomicNConcurrentSafe[float64])
}

func testAtomicNConcurrentSafe[N int64 | float64](t *testing.T) {
	var wg sync.WaitGroup
	var v atomicN[N]

	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got := v.Load()
			assert.Equal(t, int64(0), int64(got)%6)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.Store(12)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.CompareAndSwap(0, 6)
		}()
	}
	wg.Wait()
}

func BenchmarkAtomicN(b *testing.B) {
	b.Run("Int64", benchmarkAtomicN[int64])
	b.Run("Float64", benchmarkAtomicN[float64])
}

func benchmarkAtomicN[N int64 | float64](b *testing.B) {
	b.Run("Load", func(b *testing.B) {
		var a atomicN[N]
		a.Store(2)
		var v N
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v = a.Load()
			}
		})
		assert.Equal(b, N(2), v)
	})
	b.Run("Store", func(b *testing.B) {
		var a atomicN[N]
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				a.Store(3)
			}
		})
	})
	b.Run("CompareAndSwap", func(b *testing.B) {
		var a atomicN[N]
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				// Make sure we swap back and forth, in-case that matters.
				if i%2 == 0 {
					a.CompareAndSwap(0, 1)
				} else {
					a.CompareAndSwap(1, 0)
				}
				i++
			}
		})
	})
}

func TestAtomicMinMaxConcurrentSafe(t *testing.T) {
	t.Run("Int64", testAtomicMinMaxConcurrentSafe[int64])
	t.Run("Float64", testAtomicMinMaxConcurrentSafe[float64])
}

func testAtomicMinMaxConcurrentSafe[N int64 | float64](t *testing.T) {
	var wg sync.WaitGroup
	var minMax atomicMinMax[N]

	assert.False(t, minMax.set.Load())
	for _, i := range []float64{2, 4, 6, 8, -3, 0, 8, 0} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			minMax.Update(N(i))
		}()
	}
	wg.Wait()

	assert.True(t, minMax.set.Load())
	assert.Equal(t, N(-3), minMax.minimum.Load())
	assert.Equal(t, N(8), minMax.maximum.Load())
}

func BenchmarkAtomicMinMax(b *testing.B) {
	b.Run("Int64", benchmarkAtomicMinMax[int64])
	b.Run("Float64", benchmarkAtomicMinMax[float64])
}

func benchmarkAtomicMinMax[N int64 | float64](b *testing.B) {
	b.Run("UpdateIncreasing", func(b *testing.B) {
		var a atomicMinMax[N]
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				a.Update(N(i))
				i++
			}
		})
	})
	b.Run("UpdateDecreasing", func(b *testing.B) {
		var a atomicMinMax[N]
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				a.Update(N(i))
				i--
			}
		})
	})
	b.Run("UpdateConstant", func(b *testing.B) {
		var a atomicMinMax[N]
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				a.Update(N(5))
			}
		})
	})
}
