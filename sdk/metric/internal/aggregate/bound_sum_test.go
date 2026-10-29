// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

type testPinVal struct {
	attrs attribute.Set
}

func TestBoundLimitedSyncMapPin(t *testing.T) {
	newVal := func(attr attribute.Set) testPinVal {
		return testPinVal{attrs: attr}
	}

	newMap := func(limit int) boundLimitedSyncMap[testPinVal] {
		return boundLimitedSyncMap[testPinVal]{
			limitedSyncMap: limitedSyncMap[*boundEntry[testPinVal]]{aggLimit: limit},
		}
	}

	t.Run("IdleUntilMeasured", func(t *testing.T) {
		m := newMap(3)
		attr1 := attribute.NewSet(attribute.String("k", "1"))
		lazy1 := newLazyFilteredAttributes(attr1, nil)

		e1 := m.pin(lazy1, newVal)
		assert.True(t, e1.bound)
		assert.False(t, e1.measured.Load())
		assert.Equal(t, 1, m.Len(), "bound series is resident in map")

		// Duplicate pin returns identical entry.
		e1Dup := m.pin(lazy1, newVal)
		assert.Same(t, e1, e1Dup)
		assert.Equal(t, 1, m.Len())

		// Marking measured and clearing resets measured while keeping entry in map.
		e1.measured.Store(true)
		m.Clear()
		assert.False(t, e1.measured.Load())
		assert.Equal(t, 1, m.Len())
	})

	t.Run("UnboundPromotedByPin", func(t *testing.T) {
		m := newMap(3)
		attr1 := attribute.NewSet(attribute.String("k", "1"))
		lazy1 := newLazyFilteredAttributes(attr1, nil)

		vUnbound := m.LoadOrStoreAttr(lazy1, newVal)
		assert.False(t, vUnbound.bound)
		assert.Equal(t, 1, m.Len())

		eBound := m.pin(lazy1, newVal)
		assert.True(t, eBound.bound)
		assert.Same(t, vUnbound, eBound, "pin must promote existing unbound entry in-place")
		assert.Equal(t, 1, m.Len())
	})

	t.Run("PinnedReusedByUnbound", func(t *testing.T) {
		m := newMap(3)
		attr1 := attribute.NewSet(attribute.String("k", "1"))
		lazy1 := newLazyFilteredAttributes(attr1, nil)

		eBound := m.pin(lazy1, newVal)
		assert.Equal(t, 1, m.Len())

		vUnbound := m.LoadOrStoreAttr(lazy1, newVal)
		assert.Same(t, eBound, vUnbound, "unbound LoadOrStoreAttr must reuse pinned entry")
		assert.Equal(t, 1, m.Len())
	})

	t.Run("CardinalityLimitAndClearPreservesPinnedSlots", func(t *testing.T) {
		m := newMap(3) // 2 concrete + 1 overflow
		attr1 := attribute.NewSet(attribute.String("k", "1"))
		attr2 := attribute.NewSet(attribute.String("k", "2"))
		attr3 := attribute.NewSet(attribute.String("k", "3"))
		attr4 := attribute.NewSet(attribute.String("k", "4"))

		lazy1 := newLazyFilteredAttributes(attr1, nil)
		lazy2 := newLazyFilteredAttributes(attr2, nil)
		lazy3 := newLazyFilteredAttributes(attr3, nil)
		lazy4 := newLazyFilteredAttributes(attr4, nil)

		// Pin 1 concrete series (reserves 1 of 2 concrete slots).
		e1 := m.pin(lazy1, newVal)
		assert.Equal(t, attr1, e1.val.attrs)

		// Record 1 unbound concrete series (uses 2nd concrete slot) and 1 unbound overflow series.
		v2 := m.LoadOrStoreAttr(lazy2, newVal)
		v3 := m.LoadOrStoreAttr(lazy3, newVal)
		assert.Equal(t, attr2, v2.val.attrs)
		assert.Equal(t, overflowSet, v3.val.attrs)

		// An unbound measurement for the pre-bound concrete series (attr1)
		// after overflow has already occurred must still resolve to e1 (not overflowSet).
		v1Unbound := m.LoadOrStoreAttr(lazy1, newVal)
		assert.Same(t, e1, v1Unbound)
		assert.Equal(t, 3, m.Len())

		// Clear deletes unbound entries (attr2 and overflowSet),
		// while keeping attr1's pinned concrete slot resident in the map.
		m.Clear()
		assert.False(t, e1.measured.Load())
		assert.Equal(t, 1, m.Len())

		// In the next cycle, attr4 should get the 1 available concrete slot, while attr3 overflows.
		v4 := m.LoadOrStoreAttr(lazy4, newVal)
		v3Next := m.LoadOrStoreAttr(lazy3, newVal)
		assert.Equal(t, attr4, v4.val.attrs)
		assert.Equal(t, overflowSet, v3Next.val.attrs)

		// Even after overflow in the second cycle, an unbound measurement for attr1
		// must still resolve to e1 rather than routing to overflowSet.
		v1Next := m.LoadOrStoreAttr(lazy1, newVal)
		assert.Same(t, e1, v1Next)
	})

	t.Run("OverflowAttributePinnedBeforeLimit", func(t *testing.T) {
		m := newMap(3) // 2 concrete + 1 overflow
		lazyOverflow := newLazyFilteredAttributes(overflowSet, nil)
		attr2 := attribute.NewSet(attribute.String("k", "2"))
		attr3 := attribute.NewSet(attribute.String("k", "3"))
		attr4 := attribute.NewSet(attribute.String("k", "4"))
		lazy2 := newLazyFilteredAttributes(attr2, nil)
		lazy3 := newLazyFilteredAttributes(attr3, nil)
		lazy4 := newLazyFilteredAttributes(attr4, nil)

		// Pin literal overflowSet before limit (uses 1 of 2 concrete slots).
		eOver := m.pin(lazyOverflow, newVal)
		assert.False(t, m.boundOverflow)

		// Use 2nd concrete slot.
		v2 := m.LoadOrStoreAttr(lazy2, newVal)
		assert.Equal(t, attr2, v2.val.attrs)

		// 3rd series overflows into existing eOver without freeing a concrete slot.
		v3 := m.LoadOrStoreAttr(lazy3, newVal)
		assert.Same(t, eOver, v3)

		// 4th series must also overflow into eOver (not allocate a new concrete slot).
		v4 := m.LoadOrStoreAttr(lazy4, newVal)
		assert.Same(t, eOver, v4)
		e4 := m.pin(lazy4, newVal)
		assert.Same(t, eOver, e4)
	})

	t.Run("UnboundOverflowAttributePromotedByOverflowingPin", func(t *testing.T) {
		m := newMap(3) // 2 concrete + 1 overflow
		lazyOverflow := newLazyFilteredAttributes(overflowSet, nil)
		attr1 := attribute.NewSet(attribute.String("k", "1"))
		attr2 := attribute.NewSet(attribute.String("k", "2"))
		attr3 := attribute.NewSet(attribute.String("k", "3"))
		lazy1 := newLazyFilteredAttributes(attr1, nil)
		lazy2 := newLazyFilteredAttributes(attr2, nil)
		lazy3 := newLazyFilteredAttributes(attr3, nil)

		// Pin 1 concrete series and record literal overflowSet unbound (filling 2 concrete slots).
		_ = m.pin(lazy1, newVal)
		vOver := m.LoadOrStoreAttr(lazyOverflow, newVal)

		// Pinning attr2 at limit must promote vOver as the overflow entry.
		e2 := m.pin(lazy2, newVal)
		assert.Same(t, vOver, e2)
		assert.True(t, m.boundOverflow)

		// After Clear(), the 2nd concrete slot is available again for attr3.
		m.Clear()
		v3 := m.LoadOrStoreAttr(lazy3, newVal)
		assert.Equal(t, attr3, v3.val.attrs)
	})
}

func TestBoundHotColdMapPin(t *testing.T) {
	newVal := func(attr attribute.Set) testPinVal {
		return testPinVal{attrs: attr}
	}

	var m boundHotColdMap[testPinVal]
	m.init(3) // 2 concrete + 1 overflow

	attr1 := attribute.NewSet(attribute.String("k", "1"))
	attr2 := attribute.NewSet(attribute.String("k", "2"))
	attr3 := attribute.NewSet(attribute.String("k", "3"))
	lazy1 := newLazyFilteredAttributes(attr1, nil)
	lazy2 := newLazyFilteredAttributes(attr2, nil)
	lazy3 := newLazyFilteredAttributes(attr3, nil)

	// Record unbound attr1 and attr2 in the current hot buffer (filling all 2 concrete slots in hot).
	h := m.start()
	u1 := m.hot(h).LoadOrStoreAttr(lazy1, newVal)
	u1.measured.Store(true)
	_ = m.hot(h).LoadOrStoreAttr(lazy2, newVal)
	m.done(h)

	// Pinning attr1 promotes u1 in the hot buffer and creates a fresh pinned entry in the cold buffer.
	entries1 := m.pin(lazy1, newVal)
	assert.Equal(t, attr1, entries1[0].val.attrs)
	assert.Same(t, u1, entries1[h])
	assert.NotSame(t, entries1[0], entries1[1])
	assert.True(t, entries1[h].measured.Load())
	assert.False(t, entries1[1-h].measured.Load())

	// Pinning attr3 while hot has 2 concrete slots occupied must overflow in BOTH buffers.
	entries3 := m.pin(lazy3, newVal)
	assert.Equal(t, overflowSet, entries3[0].val.attrs)
	assert.Equal(t, overflowSet, entries3[1].val.attrs)

	// Swap hot/cold without clearing the now-cold buffer yet; pinning attr2 on the
	// new hot buffer must NOT alias the stale unbound u2 still sitting in the cold buffer.
	readIdx := m.swapHotAndWait()
	staleU2, ok := m.hotColdValMap[readIdx].Load(attr2.Equivalent())
	require.True(t, ok)
	entries2 := m.pin(lazy2, newVal)
	assert.NotSame(t, staleU2, entries2[readIdx], "pin must not alias un-cleared cold buffer entry")
	m.Clear(readIdx)
}

func TestBoundSum(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64/DeltaSum", testBoundDeltaSum[int64]())
	c.Reset()

	t.Run("Float64/DeltaSum", testBoundDeltaSum[float64]())
	c.Reset()

	t.Run("Int64/CumulativeSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testBoundCumulativeSum[int64]()(t)
	})
	c.Reset()

	t.Run("Int64/CumulativeSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testBoundCumulativeSum[int64]()(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testBoundCumulativeSum[float64]()(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testBoundCumulativeSum[float64]()(t)
	})
	c.Reset()
}

func testBoundDeltaSum[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()
		mono := false
		agg := Builder[N]{
			Temporality:      metricdata.DeltaTemporality,
			Filter:           attrFltr,
			AggregationLimit: 3,
		}.BoundSum(mono)
		in, out, bind := agg.Measure, agg.ComputeAggregation, agg.Bind

		boundAlice := bind(alice)
		boundBob := bind(bob)

		got := new(metricdata.Aggregation)

		// Step 0: Idle bound handles before any measurements do not report.
		assert.Equal(t, 0, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.DeltaTemporality,
			DataPoints:  []metricdata.DataPoint[N]{},
		}, *got)

		// Step 1: Bound measurements + mixed unbound measurement on alice.
		boundAlice(ctx, 1)
		boundBob(ctx, -1)
		in(ctx, 2, alice)
		boundBob(ctx, -10)

		assert.Equal(t, 2, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.DeltaTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  y2kPlus(1),
					Time:       y2kPlus(2),
					Value:      3,
				},
				{
					Attributes: fltrBob,
					StartTime:  y2kPlus(1),
					Time:       y2kPlus(2),
					Value:      -11,
				},
			},
		}, *got)

		// Step 2: Zero-net delta (+5 and -5) on alice and explicit 0 on bob MUST still report,
		// while a cycle with no measurements in Step 3 does not report.
		boundAlice(ctx, 5)
		boundAlice(ctx, -5)
		boundBob(ctx, 0)

		assert.Equal(t, 2, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.DeltaTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  y2kPlus(2),
					Time:       y2kPlus(3),
					Value:      0,
				},
				{
					Attributes: fltrBob,
					StartTime:  y2kPlus(2),
					Time:       y2kPlus(3),
					Value:      0,
				},
			},
		}, *got)

		// Step 3: Idle cycle resets and reports 0 data points.
		assert.Equal(t, 0, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.DeltaTemporality,
			DataPoints:  []metricdata.DataPoint[N]{},
		}, *got)

		// Step 4: Cardinality limit enforcement across bound and unbound measurements.
		// alice and bob already hold the 2 concrete slots; carol (bound) and dave (unbound) overflow.
		boundCarol := bind(carol)
		boundAlice(ctx, 1)
		boundBob(ctx, 2)
		boundCarol(ctx, 3)
		in(ctx, 4, dave)

		assert.Equal(t, 3, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.DeltaTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  y2kPlus(4),
					Time:       y2kPlus(5),
					Value:      1,
				},
				{
					Attributes: fltrBob,
					StartTime:  y2kPlus(4),
					Time:       y2kPlus(5),
					Value:      2,
				},
				{
					Attributes: overflowSet,
					StartTime:  y2kPlus(4),
					Time:       y2kPlus(5),
					Value:      7,
				},
			},
		}, *got)
	}
}

func testBoundCumulativeSum[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		ctx := t.Context()
		mono := false
		agg := Builder[N]{
			Temporality:      metricdata.CumulativeTemporality,
			Filter:           attrFltr,
			AggregationLimit: 3,
		}.BoundSum(mono)
		in, out, bind := agg.Measure, agg.ComputeAggregation, agg.Bind

		boundAlice := bind(alice)
		boundBob := bind(bob)

		aliceStartTime := y2kPlus(0)
		bobStartTime := y2kPlus(0)
		overflowStartTime := y2kPlus(0)
		if x.PerSeriesStartTimestamps.Enabled() {
			aliceStartTime = y2kPlus(2)
			bobStartTime = y2kPlus(3)
			overflowStartTime = y2kPlus(6)
		}

		got := new(metricdata.Aggregation)

		// Step 0: Idle bound handles before any measurements do not report.
		assert.Equal(t, 0, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.DataPoint[N]{},
		}, *got)

		// Step 1: Bound + unbound measurements accumulate into the same series.
		boundAlice(ctx, 1)
		boundBob(ctx, -1)
		in(ctx, 3, alice)
		boundBob(ctx, -10)

		assert.Equal(t, 2, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  aliceStartTime,
					Time:       y2kPlus(4),
					Value:      4,
				},
				{
					Attributes: fltrBob,
					StartTime:  bobStartTime,
					Time:       y2kPlus(4),
					Value:      -11,
				},
			},
		}, *got)

		// Step 2: Subsequent cumulative measurements add to existing totals.
		boundAlice(ctx, 10)
		boundBob(ctx, 3)

		assert.Equal(t, 2, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  aliceStartTime,
					Time:       y2kPlus(5),
					Value:      14,
				},
				{
					Attributes: fltrBob,
					StartTime:  bobStartTime,
					Time:       y2kPlus(5),
					Value:      -8,
				},
			},
		}, *got)

		// Step 3: Exceed cardinality limit via bound + unbound measurements.
		boundCarol := bind(carol)
		boundCarol(ctx, 1)
		in(ctx, 1, dave)

		assert.Equal(t, 3, out(got))
		metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[N]{
			IsMonotonic: mono,
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[N]{
				{
					Attributes: fltrAlice,
					StartTime:  aliceStartTime,
					Time:       y2kPlus(7),
					Value:      14,
				},
				{
					Attributes: fltrBob,
					StartTime:  bobStartTime,
					Time:       y2kPlus(7),
					Value:      -8,
				},
				{
					Attributes: overflowSet,
					StartTime:  overflowStartTime,
					Time:       y2kPlus(7),
					Value:      2,
				},
			},
		}, *got)
	}
}

func TestBoundSumExemplarDroppedAttributes(t *testing.T) {
	for _, temp := range []metricdata.Temporality{
		metricdata.DeltaTemporality,
		metricdata.CumulativeTemporality,
	} {
		t.Run(temp.String(), func(t *testing.T) {
			agg := Builder[int64]{
				Temporality: temp,
				Filter:      attrFltr,
				ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
					return NewFilteredExemplarReservoir[int64](
						exemplar.AlwaysOnFilter,
						exemplar.NewFixedSizeReservoir(1),
					)
				},
			}.BoundSum(true)

			bm := agg.Bind(alice) // alice has userAlice (kept) and adminTrue (dropped by attrFltr)
			bm(t.Context(), 42)

			var dest metricdata.Aggregation
			require.Equal(t, 1, agg.ComputeAggregation(&dest))
			s, ok := dest.(metricdata.Sum[int64])
			require.True(t, ok)
			require.Len(t, s.DataPoints, 1)
			assert.Equal(t, fltrAlice, s.DataPoints[0].Attributes)
			assert.Equal(t, int64(42), s.DataPoints[0].Value)
			require.Len(t, s.DataPoints[0].Exemplars, 1)
			assert.Equal(t, int64(42), s.DataPoints[0].Exemplars[0].Value)
			assert.Equal(t, []attribute.KeyValue{adminTrue}, s.DataPoints[0].Exemplars[0].FilteredAttributes)
		})
	}
}

func TestBoundCumulativeSumDelayedFirstMeasurement(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	for _, enabled := range []bool{false, true} {
		name := "PerSeriesStartTimeDisabled"
		if enabled {
			name = "PerSeriesStartTimeEnabled"
		}
		t.Run(name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", strconv.FormatBool(enabled))
			c.Reset()

			agg := Builder[int64]{
				Temporality: metricdata.CumulativeTemporality,
				Filter:      attrFltr,
			}.BoundSum(true)

			// Step 0: Bind handle while clock is at tick 0 (agg.start is tick 0).
			// Binding must not record a start timestamp before the first measurement.
			boundAlice := agg.Bind(alice)

			// Step 1: Multiple collection cycles while the handle is idle.
			var got metricdata.Aggregation
			require.Equal(t, 0, agg.ComputeAggregation(&got))
			require.Equal(t, 0, agg.ComputeAggregation(&got))

			// Step 2: First measurement occurs after delayed idle cycles.
			measurementTime := y2kPlus(c.ticks.Load())
			boundAlice(t.Context(), 42)

			// Step 3: Collect aggregation.
			require.Equal(t, 1, agg.ComputeAggregation(&got))
			s, ok := got.(metricdata.Sum[int64])
			require.True(t, ok)
			require.Len(t, s.DataPoints, 1)

			expectedStartTime := y2kPlus(0)
			if enabled {
				expectedStartTime = measurementTime
			}
			assert.Equal(t, expectedStartTime, s.DataPoints[0].StartTime)
			assert.Equal(t, int64(42), s.DataPoints[0].Value)
		})
	}
}

func TestBoundSumConcurrentSafe(t *testing.T) {
	t.Run("Int64/DeltaSum", testBoundSumConcurrentSafe[int64](metricdata.DeltaTemporality))
	t.Run("Float64/DeltaSum", testBoundSumConcurrentSafe[float64](metricdata.DeltaTemporality))
	t.Run("Int64/CumulativeSum", testBoundSumConcurrentSafe[int64](metricdata.CumulativeTemporality))
	t.Run("Float64/CumulativeSum", testBoundSumConcurrentSafe[float64](metricdata.CumulativeTemporality))
}

func testBoundSumConcurrentSafe[N int64 | float64](temp metricdata.Temporality) func(t *testing.T) {
	return func(t *testing.T) {
		agg := Builder[N]{
			Temporality:      temp,
			Filter:           attrFltr,
			AggregationLimit: 3,
		}.BoundSum(false)
		in, out, bind := agg.Measure, agg.ComputeAggregation, agg.Bind

		ctx := t.Context()
		var wg sync.WaitGroup

		attrs := make([]attribute.Set, concurrentNumGoroutines)
		for i := range attrs {
			attrs[i] = attribute.NewSet(attribute.String(keyUser, strconv.Itoa(i)), adminTrue)
		}

		vals := getConcurrentVals[N]()

		// Pre-bind half the attribute sets before concurrent writes, and bind the
		// other half dynamically from inside the worker goroutines while also
		// mixing bound and unbound measurements.
		prebound := make([]BoundMeasure[N], concurrentNumGoroutines/2)
		for i := range prebound {
			prebound[i] = bind(attrs[i])
		}

		wg.Add(concurrentNumGoroutines)
		for i := range concurrentNumGoroutines {
			go func(id int) {
				defer wg.Done()
				attr := attrs[id]
				var bm BoundMeasure[N]
				if id < len(prebound) {
					bm = prebound[id]
				}

				for j := range concurrentNumRecords {
					v := vals[j%len(vals)]
					if bm == nil && j == len(vals) {
						bm = bind(attr)
					}
					if bm != nil && j%2 == 0 {
						bm(ctx, v)
					} else {
						in(ctx, v, attr)
					}
				}
			}(i)
		}

		var results []metricdata.Aggregation
		wg.Go(func() {
			for range concurrentNumRecords {
				got := new(metricdata.Aggregation)
				out(got)
				results = append(results, *got)
			}
		})

		wg.Wait()

		got := new(metricdata.Aggregation)
		out(got)
		results = append(results, *got)

		validateSum[N](false)(t, results)
	}
}

func BenchmarkBoundSum(b *testing.B) {
	b.Run(
		"Int64/Cumulative",
		benchmarkBoundAggregate(func() (Measure[int64], ComputeAggregation, Binder[int64]) {
			s := Builder[int64]{
				Temporality: metricdata.CumulativeTemporality,
			}.BoundSum(false)
			return s.Measure, s.ComputeAggregation, s.Bind
		}),
	)
	b.Run(
		"Int64/Delta",
		benchmarkBoundAggregate(func() (Measure[int64], ComputeAggregation, Binder[int64]) {
			s := Builder[int64]{
				Temporality: metricdata.DeltaTemporality,
			}.BoundSum(false)
			return s.Measure, s.ComputeAggregation, s.Bind
		}),
	)
	b.Run(
		"Float64/Cumulative",
		benchmarkBoundAggregate(func() (Measure[float64], ComputeAggregation, Binder[float64]) {
			s := Builder[float64]{
				Temporality: metricdata.CumulativeTemporality,
			}.BoundSum(false)
			return s.Measure, s.ComputeAggregation, s.Bind
		}),
	)
	b.Run(
		"Float64/Delta",
		benchmarkBoundAggregate(func() (Measure[float64], ComputeAggregation, Binder[float64]) {
			s := Builder[float64]{
				Temporality: metricdata.DeltaTemporality,
			}.BoundSum(false)
			return s.Measure, s.ComputeAggregation, s.Bind
		}),
	)
}
