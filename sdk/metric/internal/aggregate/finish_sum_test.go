// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestFinishSumLifecycle(t *testing.T) {
	var c clock
	unregister := c.Register()
	defer unregister()

	t.Run("Cumulative", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{
			Temporality:   metricdata.CumulativeTemporality,
			ReservoirFunc: dropExemplars[int64],
		}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		agg.Finish(alice.Equivalent(), y2kPlus(20))

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(2), points[0].Value)
		assert.Equal(t, y2kPlus(0), points[0].StartTime)
		assert.Equal(t, y2kPlus(10), points[0].Time)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		agg.Measure(t.Context(), 3, alice)
		points = finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(3), points[0].Value)
		assert.Equal(t, y2kPlus(3), points[0].StartTime)
	})

	t.Run("CancelPending", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		agg.Measure(t.Context(), 3, alice)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(5), points[0].Value)
		assert.Equal(t, y2kPlus(0), points[0].StartTime)
		assert.Equal(t, y2kPlus(1), points[0].Time)
		require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 1)
	})

	t.Run("Delta", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{
			Temporality:   metricdata.DeltaTemporality,
			ReservoirFunc: dropExemplars[int64],
		}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(2), points[0].Value)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		// A series retired by a prior delta collection cannot synthesize a
		// zero-valued final point.
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		agg.Measure(t.Context(), 3, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(20))
		points = finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(3), points[0].Value)
		assert.Equal(t, y2kPlus(20), points[0].Time)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})

	t.Run("DeltaIntervalStart", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{
			Temporality:   metricdata.DeltaTemporality,
			ReservoirFunc: dropExemplars[int64],
		}.FinishSum(true)

		agg.Measure(t.Context(), 2, alice)
		first := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, first, 1)
		assert.Equal(t, y2kPlus(0), first[0].StartTime)

		agg.Measure(t.Context(), 3, alice)
		second := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, second, 1)
		assert.Equal(t, first[0].Time, second[0].StartTime)
	})
}

func TestFinishSumFinishTimestamp(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testFinishSumFinishTimestamp(t, int64(7)) })
	t.Run("Float64", func(t *testing.T) { testFinishSumFinishTimestamp(t, 1.25) })
}

func testFinishSumFinishTimestamp[N int64 | float64](t *testing.T, value N) {
	t.Helper()
	var c clock
	t.Cleanup(c.Register())
	for _, monotonic := range []bool{false, true} {
		for _, recreate := range []bool{false, true} {
			for _, offset := range []int64{-1, 0, 1} {
				t.Run(
					fmt.Sprintf("monotonic=%t/recreate=%t/offset=%d", monotonic, recreate, offset),
					func(t *testing.T) {
						c.Reset()
						agg := Builder[N]{}.FinishSum(monotonic)
						collect := func() []metricdata.DataPoint[N] {
							var data metricdata.Aggregation
							agg.ComputeAggregation(&data)
							return data.(metricdata.Sum[N]).DataPoints
						}
						if recreate {
							agg.Measure(t.Context(), value, alice)
							agg.Finish(alice.Equivalent(), y2k)
							require.Len(t, collect(), 1)
							require.Empty(t, collect())
						}
						agg.Measure(t.Context(), value, alice)
						points := collect()
						require.Len(t, points, 1)
						start := points[0].StartTime
						finishTime := start.Add(time.Duration(offset) * time.Nanosecond)
						agg.Finish(alice.Equivalent(), finishTime)
						points = collect()
						require.Len(t, points, 1)
						assert.Equal(t, start, points[0].StartTime)
						assert.Equal(t, value, points[0].Value)
						if offset < 0 {
							assert.False(t, points[0].Time.Before(start))
							require.Len(t, collect(), 1, "stale Finish must not retire the new lifetime")
							agg.Finish(alice.Equivalent(), y2kPlus(100))
							require.Len(t, collect(), 1)
						} else {
							assert.Equal(t, finishTime, points[0].Time)
						}
						assert.Empty(t, collect())
					},
				)
			}
		}
	}
}

func TestFinishSumCollectionTimestamp(t *testing.T) {
	t.Run("Int64", testFinishSumCollectionTimestamp[int64])
	t.Run("Float64", testFinishSumCollectionTimestamp[float64])
}

func testFinishSumCollectionTimestamp[N int64 | float64](t *testing.T) {
	for _, monotonic := range []bool{false, true} {
		t.Run(fmt.Sprint("monotonic=", monotonic), func(t *testing.T) {
			var c clock
			t.Cleanup(c.Register())
			agg := Builder[N]{}.FinishSum(monotonic)
			agg.Measure(t.Context(), N(1), alice)

			now = func() time.Time {
				timestamp := c.Now()
				now = c.Now
				// Create a series after the collection timestamp has been sampled.
				agg.Measure(t.Context(), N(2), bob)
				return timestamp
			}
			var data metricdata.Aggregation
			agg.ComputeAggregation(&data)
			points := data.(metricdata.Sum[N]).DataPoints
			for _, point := range points {
				assert.False(t, point.StartTime.After(point.Time), "collection emitted an inverted interval")
			}
			require.Len(t, points, 1, "series created after timestamp capture must wait for the next collection")
			assert.Equal(t, alice, points[0].Attributes)

			agg.ComputeAggregation(&data)
			points = data.(metricdata.Sum[N]).DataPoints
			require.Len(t, points, 2)
			for _, point := range points {
				assert.False(t, point.StartTime.After(point.Time))
				if point.Attributes.Equals(&bob) {
					assert.Equal(t, N(2), point.Value)
					assert.Equal(t, y2kPlus(2), point.StartTime)
					assert.Equal(t, y2kPlus(3), point.Time)
				}
			}
		})
	}
}

func TestFinishSumDeltaDeferredSeries(t *testing.T) {
	t.Run("Int64", testFinishSumDeltaDeferredSeries[int64])
	t.Run("Float64", testFinishSumDeltaDeferredSeries[float64])
}

func testFinishSumDeltaDeferredSeries[N int64 | float64](t *testing.T) {
	for _, monotonic := range []bool{false, true} {
		for _, finished := range []bool{false, true} {
			t.Run(fmt.Sprintf("monotonic=%t/finished=%t", monotonic, finished), func(t *testing.T) {
				var c clock
				t.Cleanup(c.Register())
				agg := Builder[N]{Temporality: metricdata.DeltaTemporality}.FinishSum(monotonic)
				agg.Measure(t.Context(), N(1), alice)

				now = func() time.Time {
					now = c.Now
					// Record after snapshotting but before sampling collection time.
					agg.Measure(t.Context(), N(2), bob)
					finishTime := c.Now()
					if finished {
						agg.Finish(bob.Equivalent(), finishTime)
					}
					return c.Now()
				}
				var data metricdata.Aggregation
				agg.ComputeAggregation(&data)
				points := data.(metricdata.Sum[N]).DataPoints
				require.Len(t, points, 1)
				assert.Equal(t, alice, points[0].Attributes)
				assert.Equal(t, y2k, points[0].StartTime)
				assert.Equal(t, y2kPlus(4), points[0].Time)

				agg.ComputeAggregation(&data)
				points = data.(metricdata.Sum[N]).DataPoints
				require.Len(t, points, 1)
				assert.Equal(t, bob, points[0].Attributes)
				assert.Equal(t, N(2), points[0].Value)
				assert.Equal(t, y2kPlus(2), points[0].StartTime)
				end := y2kPlus(5)
				if finished {
					end = y2kPlus(3)
				}
				assert.Equal(t, end, points[0].Time)
				assert.False(t, points[0].StartTime.After(points[0].Time))
				agg.ComputeAggregation(&data)
				assert.Empty(t, data.(metricdata.Sum[N]).DataPoints)
			})
		}
	}
}

func TestFinishSumConcurrentSafeCollectDefersNewSeries(t *testing.T) {
	barrier := &finishSumCollectionBarrier{entered: make(chan struct{}), release: make(chan struct{})}
	agg := Builder[int64]{
		Temporality: metricdata.DeltaTemporality, AggregationLimit: 3,
		ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] { return barrier },
	}.FinishSum(false)
	for _, attrs := range []attribute.Set{alice, bob, attribute.NewSet(userCarol)} {
		agg.Measure(t.Context(), 1, attrs)
	}
	result := make(chan metricdata.Aggregation, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		var data metricdata.Aggregation
		agg.ComputeAggregation(&data)
		result <- data
	}()
	var once sync.Once
	release := func() { once.Do(func() { close(barrier.release) }) }
	defer func() { release(); <-done }()
	select {
	case <-barrier.entered:
	case <-time.After(time.Second):
		t.Fatal("collection did not reach the second series")
	}
	agg.Measure(t.Context(), 123, attribute.NewSet(userDave))
	release()
	points := (<-result).(metricdata.Sum[int64]).DataPoints
	require.Len(t, points, 3, "reclaimed slots must not expand collection beyond the cardinality limit")
	for _, point := range points {
		assert.Equal(t, int64(1), point.Value)
	}
	points = finishSumPoints(t, agg.ComputeAggregation)
	require.Len(t, points, 1)
	assert.Equal(t, int64(123), points[0].Value)
	assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
}

type finishSumCollectionBarrier struct {
	calls   atomic.Int64
	entered chan struct{}
	release chan struct{}
}

func (*finishSumCollectionBarrier) Offer(context.Context, int64, lazyFilteredAttributes) {}

func (r *finishSumCollectionBarrier) Collect(*[]exemplar.Exemplar) {
	if r.calls.Add(1) == 2 {
		close(r.entered)
		<-r.release
	}
}

func TestFinishSumAttributes(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		empty := attribute.NewSet()
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 1, empty)
		agg.Finish(empty.Equivalent(), now())

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, empty, points[0].Attributes)
	})

	t.Run("Unknown", func(t *testing.T) {
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Finish(alice.Equivalent(), y2k)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})

	t.Run("FilteredExemplar", func(t *testing.T) {
		mockRes := &notConcurrentSafeReservoir{}
		agg := Builder[int64]{
			Filter: attrFltr,
			ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
				return NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, mockRes)
			},
		}.FinishSum(true)
		agg.Measure(t.Context(), 7, alice)
		agg.Finish(fltrAlice.Equivalent(), now())

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, fltrAlice, points[0].Attributes)
		require.Len(t, points[0].Exemplars, 1)
		assert.Equal(t, int64(7), points[0].Exemplars[0].Value)
		assert.Equal(t, []attribute.KeyValue{adminTrue}, points[0].Exemplars[0].FilteredAttributes)
	})
}

func TestFinishSumCardinality(t *testing.T) {
	agg := Builder[int64]{
		AggregationLimit: 3,
		ReservoirFunc:    dropExemplars[int64],
	}.FinishSum(true)
	carol := attribute.NewSet(userCarol)
	dave := attribute.NewSet(userDave)

	agg.Measure(t.Context(), 1, alice)
	agg.Measure(t.Context(), 2, bob)
	agg.Measure(t.Context(), 3, carol) // Overflow.
	agg.Finish(alice.Equivalent(), now())
	points := finishSumPoints(t, agg.ComputeAggregation)
	require.Len(t, points, 3)

	// The finalized Alice series released a normal slot even though the
	// shared overflow series remains active.
	agg.Measure(t.Context(), 4, dave)
	points = finishSumPoints(t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	assert.Equal(t, int64(4), pointWithAttrs(t, points, dave).Value)
	assert.Equal(t, int64(3), pointWithAttrs(t, points, overflowSet).Value)

	// The shared overflow series cannot be finalized, even when addressed by
	// its reserved attributes.
	agg.Finish(overflowSet.Equivalent(), now())
	require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 3)
}

func TestFinishSumOverflowProvenance(t *testing.T) {
	t.Run("UnlimitedMarkerSeries", func(t *testing.T) {
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Finish(overflowSet.Equivalent(), now())

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(1), points[0].Value)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})

	t.Run("MarkerSeriesBecomesOverflow", func(t *testing.T) {
		agg := Builder[int64]{
			AggregationLimit: 3,
			ReservoirFunc:    dropExemplars[int64],
		}.FinishSum(true)
		carol := attribute.NewSet(userCarol)

		agg.Measure(t.Context(), 1, overflowSet)
		agg.Measure(t.Context(), 2, alice)
		agg.Measure(t.Context(), 3, bob)
		agg.Finish(overflowSet.Equivalent(), now())
		// Routing this measurement to the marker series promotes it to shared
		// overflow and cancels its pending finish.
		agg.Measure(t.Context(), 4, carol)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 3)
		assert.Equal(t, int64(5), pointWithAttrs(t, points, overflowSet).Value)
		require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 3)
	})
}

func TestFinishSumInitialMeasurementAdmission(t *testing.T) {
	point := newFinishSumValue(alice, dropExemplars[int64])
	measurement, ok := point.lifecycle.AcquireMeasurement()
	require.True(t, ok)

	type collectionResult struct {
		point metricdata.DataPoint[int64]
		emit  bool
	}
	started := make(chan struct{})
	collected := make(chan collectionResult, 1)
	go func() {
		close(started)
		point, emit, _ := point.collectDelta(point.start, point.start)
		collected <- collectionResult{point: point, emit: emit}
	}()
	<-started
	select {
	case <-collected:
		t.Fatal("collection completed before the initial measurement")
	case <-time.After(10 * time.Millisecond):
	}

	point.value.add(1)
	measurement.Release()
	result := <-collected
	assert.True(t, result.emit)
	assert.Equal(t, int64(1), result.point.Value)
}

func TestFinishSumShutdown(t *testing.T) {
	agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
	agg.Measure(t.Context(), 1, alice)
	agg.Stop()
	require.NoError(t, agg.Wait(t.Context()))
	agg.Stop()
	require.NoError(t, agg.Wait(t.Context()))
	agg.Measure(t.Context(), 2, alice)
	agg.Finish(alice.Equivalent(), time.Now())
	var data metricdata.Aggregation
	assert.Zero(t, agg.ComputeAggregation(&data))
	assert.Nil(t, data)
}

func TestFinishSumCollectionScratch(t *testing.T) {
	store := newFinishSum(true, metricdata.CumulativeTemporality, 0, dropExemplars[int64])
	store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
	var data metricdata.Aggregation
	require.Equal(t, 1, store.collect(&data))
	require.Len(t, store.points, 1)
	assert.Nil(t, store.points[0], "collection scratch must not retain series")
	require.NoError(t, store.wait(t.Context()))
	assert.Nil(t, store.points)
	assert.Zero(t, store.values.Len())
}

func TestFinishSumRetireAndDelete(t *testing.T) {
	store := newFinishSum(
		true,
		metricdata.CumulativeTemporality,
		0,
		dropExemplars[int64],
	)
	lazy := newLazyFilteredAttributes(alice, nil)
	store.measure(t.Context(), 1, lazy)
	raw, ok := store.values.Load(alice.Equivalent())
	require.True(t, ok)
	point := raw.(*finishSumValue[int64])

	store.retireAndDelete(point)
	store.retireAndDelete(point)
	assert.Zero(t, store.values.Len())
	assert.False(t, point.measure(t.Context(), 1, lazy))

	_, emit, retire := point.collectCumulative(y2k)
	assert.False(t, emit)
	assert.False(t, retire)
}

func TestFinishSumMeasureDeletesRetiredPoint(t *testing.T) {
	store := newFinishSum(
		true,
		metricdata.CumulativeTemporality,
		0,
		dropExemplars[int64],
	)
	lazy := newLazyFilteredAttributes(alice, nil)
	store.measure(t.Context(), 1, lazy)
	raw, ok := store.values.Load(alice.Equivalent())
	require.True(t, ok)
	retired := raw.(*finishSumValue[int64])
	require.NoError(t, retired.shutdown(t.Context()))

	store.measure(t.Context(), 2, lazy)

	raw, ok = store.values.Load(alice.Equivalent())
	require.True(t, ok)
	replacement := raw.(*finishSumValue[int64])
	assert.NotSame(t, retired, replacement)
	assert.Equal(t, int64(2), replacement.value.load())
	assert.Equal(t, 1, store.values.Len())
}

func TestFinishSumConcurrentSafeLifecycle(t *testing.T) {
	agg := Builder[int64]{
		Temporality:   metricdata.DeltaTemporality,
		ReservoirFunc: dropExemplars[int64],
	}.FinishSum(true)

	const measurements = 2_000
	var (
		collected atomic.Int64
		stop      atomic.Bool
		wg        sync.WaitGroup
	)
	wg.Go(func() {
		for !stop.Load() {
			agg.Finish(alice.Equivalent(), time.Now())
		}
	})
	wg.Go(func() {
		for !stop.Load() {
			var data metricdata.Aggregation
			agg.ComputeAggregation(&data)
			for _, point := range data.(metricdata.Sum[int64]).DataPoints {
				collected.Add(point.Value)
			}
		}
	})
	for range measurements {
		agg.Measure(t.Context(), 1, alice)
	}
	stop.Store(true)
	wg.Wait()
	for _, point := range finishSumPoints(t, agg.ComputeAggregation) {
		collected.Add(point.Value)
	}
	assert.Equal(t, int64(measurements), collected.Load())
}

func TestFinishSumConcurrentSafeShutdown(t *testing.T) {
	reservoir := &blockingFinishSumReservoir{
		offered: make(chan struct{}),
		release: make(chan struct{}),
	}
	agg := Builder[int64]{
		ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
			return reservoir
		},
	}.FinishSum(true)

	measured := make(chan struct{})
	go func() {
		agg.Measure(t.Context(), 1, alice)
		close(measured)
	}()
	<-reservoir.offered

	agg.Stop()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, agg.Wait(ctx), context.Canceled)

	close(reservoir.release)
	<-measured
	require.NoError(t, agg.Wait(t.Context()))
	var data metricdata.Aggregation
	assert.Zero(t, agg.ComputeAggregation(&data))
}

func TestFinishSumWaitCancellationDuringCreation(t *testing.T) {
	for _, populated := range []bool{false, true} {
		t.Run(fmt.Sprint("populated=", populated), func(t *testing.T) {
			entered, release := make(chan struct{}), make(chan struct{})
			store := newFinishSum(true, metricdata.CumulativeTemporality, 0,
				func(attrs attribute.Set) FilteredExemplarReservoir[int64] {
					if attrs.Equals(&bob) {
						close(entered)
						<-release
					}
					return DropReservoir[int64](attrs)
				})
			if populated {
				store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
			}

			measured := make(chan struct{})
			done := make(chan struct{})
			var shutdownErr error
			started := false
			go func() {
				defer close(measured)
				store.measure(t.Context(), 1, newLazyFilteredAttributes(bob, nil))
			}()
			defer func() {
				close(release)
				<-measured
				if started {
					<-done
				}
				if !populated {
					assert.Zero(t, store.values.Len(), "late publication retained a stopped series")
				}
				require.NoError(t, store.wait(t.Context()))
				assert.Zero(t, store.values.Len())
			}()
			select {
			case <-entered:
			case <-time.After(time.Second):
				t.Fatal("measurement did not enter reservoir construction")
			}

			ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
			defer cancel()
			started = true
			go func() {
				defer close(done)
				shutdownErr = store.wait(ctx)
			}()
			select {
			case <-done:
				assert.ErrorIs(t, shutdownErr, context.DeadlineExceeded)
			case <-time.After(time.Second):
				t.Error("shutdown ignored cancellation during reservoir construction")
			}
			assert.True(t, store.stopped.Load())
		})
	}
}

func TestFinishSumWaitCancellationDuringCollection(t *testing.T) {
	reservoir := &blockingFinishSumCollectionReservoir{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
	store := newFinishSum(true, metricdata.CumulativeTemporality, 0,
		func(attribute.Set) FilteredExemplarReservoir[int64] { return reservoir })
	store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
	collected := make(chan struct{})
	done := make(chan struct{})
	var shutdownErr error
	started := false
	go func() {
		defer close(collected)
		var data metricdata.Aggregation
		store.collect(&data)
	}()
	defer func() {
		close(reservoir.release)
		<-collected
		if started {
			<-done
		}
		require.NoError(t, store.wait(t.Context()))
		assert.Zero(t, store.values.Len())
	}()
	select {
	case <-reservoir.entered:
	case <-time.After(time.Second):
		t.Fatal("collection did not enter reservoir")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()
	started = true
	go func() {
		defer close(done)
		shutdownErr = store.wait(ctx)
	}()
	select {
	case <-done:
		assert.ErrorIs(t, shutdownErr, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Error("shutdown ignored cancellation during collection")
	}
	assert.True(t, store.stopped.Load())
}

func TestFinishSumWaitCanceledBeforeCollection(t *testing.T) {
	store := newFinishSum(true, metricdata.CumulativeTemporality, 0, dropExemplars[int64])
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, store.wait(ctx), context.Canceled)
	require.NoError(t, store.wait(t.Context()))
}

type blockingFinishSumCollectionReservoir struct {
	entered chan struct{}
	release chan struct{}
}

func (*blockingFinishSumCollectionReservoir) Offer(context.Context, int64, lazyFilteredAttributes) {}

func (r *blockingFinishSumCollectionReservoir) Collect(*[]exemplar.Exemplar) {
	close(r.entered)
	<-r.release
}

func BenchmarkFinishSumCollect(b *testing.B) {
	for _, count := range []int{1, 100, 1000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
			for i := range count {
				agg.Measure(b.Context(), 1, attribute.NewSet(attribute.Int("id", i)))
			}
			var data metricdata.Aggregation
			agg.ComputeAggregation(&data)
			b.ReportAllocs()
			for b.Loop() {
				agg.ComputeAggregation(&data)
			}
		})
	}
}

type blockingFinishSumReservoir struct {
	offered chan struct{}
	release chan struct{}
}

func (r *blockingFinishSumReservoir) Offer(
	context.Context,
	int64,
	lazyFilteredAttributes,
) {
	close(r.offered)
	<-r.release
}

func (*blockingFinishSumReservoir) Collect(*[]exemplar.Exemplar) {}

func BenchmarkFinishSumLifecycle(b *testing.B) {
	b.Run("Int64", benchmarkFinishSumLifecycle[int64])
	b.Run("Float64", benchmarkFinishSumLifecycle[float64])
}

func benchmarkFinishSumLifecycle[N int64 | float64](b *testing.B) {
	for _, name := range []string{"MeasureAndFinish", "RepeatedFinish"} {
		b.Run(name, func(b *testing.B) {
			agg := Builder[N]{}.FinishSum(true)
			agg.Measure(b.Context(), N(1), alice)
			finishTime := time.Now()
			distinct := alice.Equivalent()
			b.ReportAllocs()
			for b.Loop() {
				if name == "MeasureAndFinish" {
					agg.Measure(b.Context(), N(1), alice)
				}
				agg.Finish(distinct, finishTime)
			}
		})
	}
}

func BenchmarkFinishSum(b *testing.B) {
	b.Run("Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		agg := Builder[int64]{Temporality: metricdata.CumulativeTemporality}.FinishSum(true)
		return agg.Measure, agg.ComputeAggregation
	}))
	b.Run("Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		agg := Builder[int64]{Temporality: metricdata.DeltaTemporality}.FinishSum(true)
		return agg.Measure, agg.ComputeAggregation
	}))
}

func BenchmarkFinishSumMeasure(b *testing.B) {
	benchmarks := []struct {
		name    string
		factory func() Measure[int64]
	}{
		{"false", func() Measure[int64] {
			measure, _ := Builder[int64]{Temporality: metricdata.CumulativeTemporality}.Sum(true)
			return measure
		}},
		{"true", func() Measure[int64] {
			return Builder[int64]{Temporality: metricdata.CumulativeTemporality}.FinishSum(true).Measure
		}},
	}
	for _, benchmark := range benchmarks {
		b.Run("finish="+benchmark.name, func(b *testing.B) {
			b.Run("mode=serial", func(b *testing.B) {
				measure := benchmark.factory()
				b.ReportAllocs()
				for b.Loop() {
					measure(b.Context(), 1, alice)
				}
			})
			b.Run("mode=parallel", func(b *testing.B) {
				measure := benchmark.factory()
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						measure(b.Context(), 1, alice)
					}
				})
			})
		})
	}
}

func BenchmarkFinishSumOverflowMeasurement(b *testing.B) {
	lazy := newLazyFilteredAttributes(overflowSet, nil)
	b.Run("state=promotion", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			point := newFinishSumValue(overflowSet, dropExemplars[int64])
			measurement, ok := point.lifecycle.AcquireMeasurement()
			if !ok {
				b.Fatal("initial measurement was not admitted")
			}
			measurement.Release()
			if !point.measureOverflow(b.Context(), 1, lazy) {
				b.Fatal("overflow measurement was not admitted")
			}
		}
	})
	b.Run("state=promoted", func(b *testing.B) {
		point := newFinishSumValue(overflowSet, dropExemplars[int64])
		measurement, ok := point.lifecycle.AcquireSharedMeasurement()
		if !ok {
			b.Fatal("shared measurement was not admitted")
		}
		measurement.Release()
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			if !point.measureOverflow(b.Context(), 1, lazy) {
				b.Fatal("overflow measurement was not admitted")
			}
		}
	})
}

func finishSumPoints(
	t testing.TB,
	compute ComputeAggregation,
) []metricdata.DataPoint[int64] {
	t.Helper()
	var data metricdata.Aggregation
	compute(&data)
	sum, ok := data.(metricdata.Sum[int64])
	require.True(t, ok, "aggregation type %T", data)
	return sum.DataPoints
}

func pointWithAttrs(
	t testing.TB,
	points []metricdata.DataPoint[int64],
	attrs attribute.Set,
) metricdata.DataPoint[int64] {
	t.Helper()
	for _, point := range points {
		if point.Attributes.Equals(&attrs) {
			return point
		}
	}
	t.Fatalf("no point with attributes %v in %#v", attrs, points)
	return metricdata.DataPoint[int64]{}
}
