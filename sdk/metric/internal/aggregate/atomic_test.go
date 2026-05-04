// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
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
		wg.Go(func() {
			aSum.add(in)
		})
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
		wg.Go(func() {
			aSum.add(in)
		})
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
	var data [2]atomic.Uint64
	for range 5 {
		wg.Go(func() {
			hotIdx := hcwg.start()
			defer hcwg.done(hotIdx)
			data[hotIdx].Add(1)
		})
	}
	for range 2 {
		readIdx := hcwg.swapHotAndWait()
		assert.NotPanics(t, func() {
			// reading without using atomics should not panic since we are
			// reading from the cold element, and have waited for all writes to
			// finish.
			t.Logf("read value %+v", data[readIdx].Load())
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
		wg.Go(func() {
			got := v.Load()
			assert.Equal(t, int64(0), int64(got)%6)
		})
		wg.Go(func() {
			v.Store(12)
		})
		wg.Go(func() {
			v.CompareAndSwap(0, 6)
		})
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
		wg.Go(func() {
			minMax.Update(N(i))
		})
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

func TestLimitedSyncMapLimit(t *testing.T) {
	m := newLimitedSyncMap[any](3)
	newValue := func(attribute.Set) any { return new(int) }

	attr1 := attribute.NewSet(attribute.String("key", "1"))
	attr2 := attribute.NewSet(attribute.String("key", "2"))
	attr3 := attribute.NewSet(attribute.String("key", "3"))
	attr4 := attribute.NewSet(attribute.String("key", "4"))

	// Add first (normal)
	v1 := m.LoadOrStoreAttr(attr1, newValue)
	assert.Equal(t, 1, m.Len())

	// Add second (normal)
	v2 := m.LoadOrStoreAttr(attr2, newValue)
	assert.Equal(t, 2, m.Len())

	// Add third (overflow)
	v3 := m.LoadOrStoreAttr(attr3, newValue)
	assert.Equal(t, 3, m.Len()) // Overflow counts as the 3rd entry
	assert.NotSame(t, v1, v3)
	assert.NotSame(t, v2, v3)

	// Add fourth (overflow) - should return same overflow value
	v4 := m.LoadOrStoreAttr(attr4, newValue)
	assert.Same(t, v3, v4)

	// Clear the map. Should be able to add new keys up to limit again.
	m.Clear()
	assert.Equal(t, 0, m.Len())

	attr5 := attribute.NewSet(attribute.String("key", "5"))
	attr6 := attribute.NewSet(attribute.String("key", "6"))
	attr7 := attribute.NewSet(attribute.String("key", "7"))
	attr8 := attribute.NewSet(attribute.String("key", "8"))

	v5 := m.LoadOrStoreAttr(attr5, newValue)
	assert.Equal(t, 1, m.Len())

	v6 := m.LoadOrStoreAttr(attr6, newValue)
	assert.Equal(t, 2, m.Len())

	assert.NotSame(t, v5, v6, "Different keys should return different values")

	v7 := m.LoadOrStoreAttr(attr7, newValue)
	assert.Equal(t, 3, m.Len()) // Overflow counts as 3rd entry
	assert.NotSame(t, v5, v7, "Overflow should be different from normal values")
	assert.NotSame(t, v6, v7, "Overflow should be different from normal values")

	v8 := m.LoadOrStoreAttr(attr8, newValue)
	assert.Same(t, v7, v8, "Subsequent keys should return same overflow value")
}

func TestLimitedSyncMapConcurrentSafe(t *testing.T) {
	m := newLimitedSyncMap[any](5)
	newValue := func(attribute.Set) any { return 1 }
	attr := attribute.NewSet(attribute.String("k", "v"))

	var wg sync.WaitGroup
	// 100 routines trying to read/write the same key
	for range 100 {
		wg.Go(func() {
			m.LoadOrStoreAttr(attr, newValue)
		})
	}
	wg.Wait()
	assert.Equal(t, 1, m.Len())

	// 10 routines trying to read/write DIFFERENT keys exceeding limit
	var wg2 sync.WaitGroup
	attrs := []attribute.Set{
		attribute.NewSet(attribute.String("k", "1")),
		attribute.NewSet(attribute.String("k", "2")),
		attribute.NewSet(attribute.String("k", "3")),
		attribute.NewSet(attribute.String("k", "4")),
		attribute.NewSet(attribute.String("k", "5")),
		attribute.NewSet(attribute.String("k", "6")),
		attribute.NewSet(attribute.String("k", "7")),
		attribute.NewSet(attribute.String("k", "8")),
		attribute.NewSet(attribute.String("k", "9")),
		attribute.NewSet(attribute.String("k", "10")),
	}
	for _, a := range attrs {
		attrCopy := a
		wg2.Go(func() {
			m.LoadOrStoreAttr(attrCopy, newValue)
		})
	}
	wg2.Wait()
	// Map should be at limit (5)
	assert.Equal(t, 5, m.Len())
}

func BenchmarkSyncMap(b *testing.B) {
	attr := attribute.NewSet(attribute.String("key", "value"))
	newValue := func(attribute.Set) any { return 1 }

	b.Run("limitedSyncMap/LoadOrStoreNoClear", func(b *testing.B) {
		m := newLimitedSyncMap[any](10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.LoadOrStoreAttr(attr, newValue)
		}
	})

	b.Run("limitedSyncMap/LoadOrStoreWithClear", func(b *testing.B) {
		m := newLimitedSyncMap[any](10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Clear()
			m.LoadOrStoreAttr(attr, newValue)
		}
	})

	b.Run("limitedSyncMap/OnlyClear", func(b *testing.B) {
		m := newLimitedSyncMap[any](10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Clear()
		}
	})
}
