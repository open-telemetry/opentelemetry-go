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

func BenchmarkSyncMap(b *testing.B) {
	tests := []struct {
		name    string
		makeMap func() syncMap
	}{
		{"limitedSyncMap", func() syncMap { return limitedSyncMapTestWrapper{&limitedSyncMap{}} }},
		{"lazyLimitedSyncMap", func() syncMap { return lazySyncMapTestWrapper{&lazyLimitedSyncMap[any]{}} }},
	}

	attr := attribute.NewSet(attribute.String("key", "value"))
	newValue := func(attribute.Set) any { return 1 }

	for _, tt := range tests {
		b.Run(tt.name+"/LoadOrStoreNoClear", func(b *testing.B) {
			m := tt.makeMap()
			if w, ok := m.(lazySyncMapTestWrapper); ok {
				w.newValue = newValue
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.LoadOrReuseAttr(attr, newValue)
			}
		})

		b.Run(tt.name+"/LoadOrStoreWithClear", func(b *testing.B) {
			m := tt.makeMap()
			if w, ok := m.(lazySyncMapTestWrapper); ok {
				w.newValue = newValue
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Clear()
				m.LoadOrReuseAttr(attr, newValue)
			}
		})

		b.Run(tt.name+"/OnlyClear", func(b *testing.B) {
			m := tt.makeMap()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Clear()
			}
		})
	}
}

// syncMap represents the shared functionality between limitedSyncMap and lazyLimitedSyncMap.
type syncMap interface {
	LoadOrReuseAttr(fltrAttr attribute.Set, newValue func(attribute.Set) any) any
	Clear()
	Len() int
	Range(f func(key, value any) bool)
}

type limitedSyncMapTestWrapper struct {
	*limitedSyncMap
}

func (w limitedSyncMapTestWrapper) Clear() {
	w.limitedSyncMap.Clear()
}

func (w limitedSyncMapTestWrapper) LoadOrReuseAttr(fltrAttr attribute.Set, newValue func(attribute.Set) any) any {
	return w.LoadOrStoreAttr(fltrAttr, newValue)
}

type lazySyncMapTestWrapper struct {
	*lazyLimitedSyncMap[any]
}

func (w lazySyncMapTestWrapper) Clear() {
	w.lazyLimitedSyncMap.Clear()
}

func (w lazySyncMapTestWrapper) LoadOrReuseAttr(fltrAttr attribute.Set, _ func(attribute.Set) any) any {
	return w.lazyLimitedSyncMap.LoadOrReuseAttr(fltrAttr)
}

func TestSyncMap_Limit(t *testing.T) {
	tests := []struct {
		name    string
		makeMap func(limit int) syncMap
	}{
		{
			"limitedSyncMap",
			func(limit int) syncMap { return limitedSyncMapTestWrapper{&limitedSyncMap{aggLimit: limit}} },
		},
		{
			"lazyLimitedSyncMap",
			func(limit int) syncMap { return lazySyncMapTestWrapper{&lazyLimitedSyncMap[any]{aggLimit: limit}} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We want 2 normal attributes and 1 overflow, so limit is 3.
			m := tt.makeMap(3)

			attr1 := attribute.NewSet(attribute.String("key", "1"))
			attr2 := attribute.NewSet(attribute.String("key", "2"))
			attr3 := attribute.NewSet(attribute.String("key", "3"))
			attr4 := attribute.NewSet(attribute.String("key", "4"))

			newVal := func(attribute.Set) any { return new(int) }
			if w, ok := m.(lazySyncMapTestWrapper); ok {
				w.newValue = newVal
			}

			// Add first (normal)
			v1 := m.LoadOrReuseAttr(attr1, newVal)
			assert.Equal(t, 1, m.Len())

			// Add second (normal)
			v2 := m.LoadOrReuseAttr(attr2, newVal)
			assert.Equal(t, 2, m.Len())

			// Add third (overflow)
			v3 := m.LoadOrReuseAttr(attr3, newVal)
			assert.Equal(t, 3, m.Len()) // Overflow counts as the 3rd entry
			assert.NotSame(t, v1, v3)
			assert.NotSame(t, v2, v3)

			// Add fourth (overflow) - should return same overflow value
			v4 := m.LoadOrReuseAttr(attr4, newVal)
			assert.Same(t, v3, v4)
		})
	}
}

func TestSyncMap_Concurrent(t *testing.T) {
	tests := []struct {
		name    string
		makeMap func() syncMap
	}{
		{"limitedSyncMap", func() syncMap { return limitedSyncMapTestWrapper{&limitedSyncMap{aggLimit: 5}} }},
		{"lazyLimitedSyncMap", func() syncMap { return lazySyncMapTestWrapper{&lazyLimitedSyncMap[any]{aggLimit: 5}} }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.makeMap()
			if w, ok := m.(lazySyncMapTestWrapper); ok {
				w.newValue = func(attribute.Set) any { return 1 }
			}
			attr := attribute.NewSet(attribute.String("k", "v"))

			var wg sync.WaitGroup
			// 100 routines trying to read/write the same key
			for range 100 {
				wg.Go(func() {
					m.LoadOrReuseAttr(attr, func(attribute.Set) any { return 1 })
				})
			}
			wg.Wait()
			assert.Equal(t, 1, m.Len())

			// 100 routines clearing and loading
			for range 100 {
				wg.Go(func() {
					m.Clear()
					m.LoadOrReuseAttr(attr, func(attribute.Set) any { return 1 })
				})
			}
			wg.Wait()
			// At the end, len should be 1
			assert.Equal(t, 1, m.Len())
		})
	}
}

// Specific tests for lazyLimitedSyncMap's cycle and GC behavior.
func TestLazyLimitedSyncMap_ClearAndReuse(t *testing.T) {
	var m lazyLimitedSyncMap[any]
	m.aggLimit = 10
	attr1 := attribute.NewSet(attribute.String("k", "v"))

	allocCount := 0
	newVal := func(attribute.Set) any {
		allocCount++
		return allocCount
	}

	// Cycle 0
	m.newValue = func(attribute.Set) any { return newVal(attr1) }
	v1 := m.LoadOrReuseAttr(attr1)
	assert.Equal(t, 1, v1)
	assert.Equal(t, 1, m.Len())

	// Clear -> moves to Cycle 1
	m.Clear()
	assert.Equal(t, 0, m.Len())

	// Re-inserting the same key should reuse the map entry, without calling newVal
	v2 := m.LoadOrReuseAttr(attr1)
	assert.Equal(t, 1, v2, "Value should be reused without calling newVal")
	assert.Equal(t, 1, m.Len())

	// Check underlying map length to ensure no growth
	physLen := 0
	m.m.Range(func(_, _ any) bool {
		physLen++
		return true
	})
	assert.Equal(t, 1, physLen, "Underlying map should not grow when reusing keys")
}

func TestLazyLimitedSyncMap_RangeAndGC(t *testing.T) {
	var m lazyLimitedSyncMap[any]
	m.aggLimit = 10
	attr1 := attribute.NewSet(attribute.String("k", "v"))

	newVal := func(attribute.Set) any { return 1 }

	// Cycle 0: add item
	m.newValue = newVal
	m.LoadOrReuseAttr(attr1)

	// Cycle 1: item is stale
	m.Clear()

	// Range should yield nothing
	yieldCount := 0
	m.Range(func(_, _ any) bool {
		yieldCount++
		return true
	})
	assert.Equal(t, 0, yieldCount, "Stale items should not be yielded")

	// Move to Cycle 4 (Cycle 0 + 4)
	m.Clear() // Cycle 2
	m.Clear() // Cycle 3
	m.Clear() // Cycle 4

	// Underlying map should still have the item
	physLen := 0
	m.m.Range(func(_, _ any) bool { physLen++; return true })
	assert.Equal(t, 1, physLen)

	// Range at Cycle 4 should trigger GC for Cycle 0 item (4 > 0 + 3)
	m.Range(func(_, _ any) bool { return true })

	// Underlying map should now be empty
	physLen = 0
	m.m.Range(func(_, _ any) bool { physLen++; return true })
	assert.Equal(t, 0, physLen, "Item should be GC'd after 4 cycles")
}
