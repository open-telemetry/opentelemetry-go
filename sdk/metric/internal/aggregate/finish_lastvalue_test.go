// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"fmt"
	"slices"
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

func TestFinishLastValueLifecycle(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testFinishLastValueLifecycle(t, int64(7), int64(-3)) })
	t.Run("Float64", func(t *testing.T) { testFinishLastValueLifecycle(t, 1.25, -0.5) })
}

func testFinishLastValueLifecycle[N int64 | float64](t *testing.T, first, last N) {
	t.Helper()
	var c clock
	t.Cleanup(c.Register())
	for _, temporality := range []metricdata.Temporality{
		metricdata.CumulativeTemporality, metricdata.DeltaTemporality,
	} {
		t.Run(temporality.String(), func(t *testing.T) {
			c.Reset()
			agg := Builder[N]{Temporality: temporality}.FinishLastValue()
			agg.Finish(bob.Equivalent(), y2k)
			assert.Empty(t, finishLastValuePoints[N](t, agg.ComputeAggregation))
			agg.Measure(t.Context(), first, alice)
			agg.Measure(t.Context(), last, alice)
			agg.Finish(alice.Equivalent(), y2kPlus(10))
			agg.Finish(alice.Equivalent(), y2kPlus(20))
			points := finishLastValuePoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, last, points[0].Value)
			assert.Equal(t, alice, points[0].Attributes)
			assert.Equal(t, y2kPlus(10), points[0].Time)
			start := points[0].StartTime
			assert.Empty(t, finishLastValuePoints[N](t, agg.ComputeAggregation))

			agg.Measure(t.Context(), first, alice)
			points = finishLastValuePoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, first, points[0].Value)
			assert.True(t, points[0].StartTime.After(start), "recreated series needs a new start")
			agg.Finish(alice.Equivalent(), y2kPlus(30))
			points = finishLastValuePoints[N](t, agg.ComputeAggregation)
			if temporality == metricdata.DeltaTemporality {
				assert.Empty(t, points, "Finish must not synthesize an empty delta point")
			} else {
				require.Len(t, points, 1)
				assert.Equal(t, first, points[0].Value)
				assert.Equal(t, y2kPlus(30), points[0].Time)
			}
		})
		t.Run(temporality.String()+"/Reactivate", func(t *testing.T) {
			c.Reset()
			agg := Builder[N]{Temporality: temporality}.FinishLastValue()
			agg.Measure(t.Context(), first, alice)
			agg.Finish(alice.Equivalent(), y2kPlus(10))
			agg.Measure(t.Context(), last, alice)
			points := finishLastValuePoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, last, points[0].Value)
			assert.NotEqual(t, y2kPlus(10), points[0].Time)
			points = finishLastValuePoints[N](t, agg.ComputeAggregation)
			if temporality == metricdata.CumulativeTemporality {
				require.Len(t, points, 1)
				assert.Equal(t, last, points[0].Value)
			} else {
				assert.Empty(t, points)
			}
		})
	}
}

func TestFinishLastValueAttributes(t *testing.T) {
	t.Run("EmptyZero", func(t *testing.T) {
		agg := Builder[int64]{}.FinishLastValue()
		attrs := attribute.NewSet()
		agg.Measure(t.Context(), 0, attrs)
		agg.Finish(attrs.Equivalent(), y2k)
		points := finishLastValuePoints[int64](t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Zero(t, points[0].Value)
		assert.Equal(t, attrs, points[0].Attributes)
		assert.Empty(t, finishLastValuePoints[int64](t, agg.ComputeAggregation))
	})
	t.Run("FilteredExemplar", func(t *testing.T) {
		mockRes := &notConcurrentSafeReservoir{}
		agg := Builder[int64]{
			Filter: attrFltr,
			ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
				return NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, mockRes)
			},
		}.FinishLastValue()
		agg.Measure(t.Context(), 7, alice)
		agg.Finish(fltrAlice.Equivalent(), y2k)
		points := finishLastValuePoints[int64](t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, fltrAlice, points[0].Attributes)
		require.Len(t, points[0].Exemplars, 1)
		assert.Equal(t, int64(7), points[0].Exemplars[0].Value)
		assert.Equal(t, []attribute.KeyValue{adminTrue}, points[0].Exemplars[0].FilteredAttributes)
		assert.Empty(t, finishLastValuePoints[int64](t, agg.ComputeAggregation))
	})
}

func TestFinishLastValueCardinality(t *testing.T) {
	agg := Builder[int64]{AggregationLimit: 3}.FinishLastValue()
	carol := attribute.NewSet(userCarol)
	dave := attribute.NewSet(userDave)
	agg.Measure(t.Context(), 1, alice)
	agg.Measure(t.Context(), 2, bob)
	agg.Measure(t.Context(), 3, carol)
	agg.Measure(t.Context(), 4, carol)
	agg.Finish(carol.Equivalent(), y2k)
	agg.Finish(alice.Equivalent(), y2k)
	points := finishLastValuePoints[int64](t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	assert.Equal(t, int64(4), pointWithAttrs(t, points, overflowSet).Value)

	agg.Measure(t.Context(), 5, dave)
	agg.Finish(overflowSet.Equivalent(), y2k)
	points = finishLastValuePoints[int64](t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	assert.Equal(t, int64(5), pointWithAttrs(t, points, dave).Value)
	assert.Equal(t, int64(4), pointWithAttrs(t, points, overflowSet).Value)
}

func TestFinishLastValueOverflowProvenance(t *testing.T) {
	t.Run("UnlimitedMarkerSeries", func(t *testing.T) {
		agg := Builder[int64]{}.FinishLastValue()
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Finish(overflowSet.Equivalent(), y2k)
		require.Len(t, finishLastValuePoints[int64](t, agg.ComputeAggregation), 1)
		assert.Empty(t, finishLastValuePoints[int64](t, agg.ComputeAggregation))
	})
	t.Run("MarkerSeriesBecomesOverflow", func(t *testing.T) {
		agg := Builder[int64]{AggregationLimit: 3}.FinishLastValue()
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Measure(t.Context(), 2, alice)
		agg.Measure(t.Context(), 3, bob)
		agg.Finish(overflowSet.Equivalent(), y2k)
		agg.Measure(t.Context(), 4, attribute.NewSet(userCarol))
		points := finishLastValuePoints[int64](t, agg.ComputeAggregation)
		require.Len(t, points, 3)
		assert.Equal(t, int64(4), pointWithAttrs(t, points, overflowSet).Value)
		require.Len(t, finishLastValuePoints[int64](t, agg.ComputeAggregation), 3)
	})
}

func TestFinishLastValueShutdown(t *testing.T) {
	store := newFinishLastValue(metricdata.CumulativeTemporality, 0, dropExemplars[int64])
	store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
	var data metricdata.Aggregation
	require.Equal(t, 1, store.collect(&data))
	require.Len(t, store.points, 1)
	assert.Nil(t, store.points[0], "collection scratch must not retain series")
	store.stop()
	require.NoError(t, store.wait(t.Context()))
	require.NoError(t, store.wait(t.Context()))
	assert.Zero(t, store.values.Len())
	assert.Nil(t, store.points)
	store.measure(t.Context(), 2, newLazyFilteredAttributes(alice, nil))
	store.finish(alice.Equivalent(), y2k)
	data = nil
	assert.Zero(t, store.collect(&data))
	assert.Nil(t, data)
}

func TestFinishLastValueInitialMeasurementAdmission(t *testing.T) {
	point := newFinishLastValuePoint(alice, dropExemplars[int64])
	measurement, ok := point.lifecycle.AcquireMeasurement()
	require.True(t, ok)
	started := make(chan struct{})
	collected := make(chan metricdata.DataPoint[int64], 1)
	go func() {
		close(started)
		dp, emit, _ := point.collectCumulative(y2k)
		if !emit {
			collected <- metricdata.DataPoint[int64]{}
			return
		}
		collected <- dp
	}()
	<-started
	select {
	case <-collected:
		t.Fatal("collection completed before the initial measurement")
	case <-time.After(10 * time.Millisecond):
	}
	point.value.Store(7)
	measurement.Release()
	assert.Equal(t, int64(7), (<-collected).Value)
}

func TestFinishLastValueCollectionTimestamp(t *testing.T) {
	t.Run("Int64", testFinishLastValueCollectionTimestamp[int64])
	t.Run("Float64", testFinishLastValueCollectionTimestamp[float64])
}

func testFinishLastValueCollectionTimestamp[N int64 | float64](t *testing.T) {
	t.Helper()
	var c clock
	t.Cleanup(c.Register())
	agg := Builder[N]{}.FinishLastValue()
	agg.Measure(t.Context(), N(1), alice)

	now = func() time.Time {
		timestamp := c.Now()
		now = c.Now
		// Create a series after the collection timestamp has been sampled.
		agg.Measure(t.Context(), N(2), bob)
		return timestamp
	}
	points := finishLastValuePoints[N](t, agg.ComputeAggregation)
	for _, point := range points {
		assert.False(t, point.StartTime.After(point.Time), "collection emitted an inverted interval")
	}
	require.Len(t, points, 1, "series created after timestamp capture must wait for the next collection")
	assert.Equal(t, alice, points[0].Attributes)
	assert.Equal(t, y2k, points[0].StartTime)
	assert.Equal(t, y2kPlus(1), points[0].Time)

	points = finishLastValuePoints[N](t, agg.ComputeAggregation)
	require.Len(t, points, 2)
	i := slices.IndexFunc(points, func(point metricdata.DataPoint[N]) bool {
		return point.Attributes.Equals(&bob)
	})
	require.NotEqual(t, -1, i)
	assert.Equal(t, y2kPlus(2), points[i].StartTime)
	assert.Equal(t, y2kPlus(3), points[i].Time)
	assert.Equal(t, N(2), points[i].Value)
	assert.Empty(t, points[i].Exemplars)
}

func TestFinishLastValueCollectDefersNewSeries(t *testing.T) {
	barrier := &lastValueCollectionBarrier{entered: make(chan struct{}), release: make(chan struct{})}
	agg := Builder[int64]{
		Temporality: metricdata.DeltaTemporality, AggregationLimit: 4,
		ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] { return barrier },
	}.FinishLastValue()
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
	dave := attribute.NewSet(userDave)
	agg.Measure(t.Context(), 123, dave)
	release()
	gauge := (<-result).(metricdata.Gauge[int64])
	require.Len(t, gauge.DataPoints, 3)
	points := finishLastValuePoints[int64](t, agg.ComputeAggregation)
	require.Len(t, points, 1)
	assert.Equal(t, dave, points[0].Attributes)
	assert.Equal(t, int64(123), points[0].Value)
}

type lastValueCollectionBarrier struct {
	calls   atomic.Int64
	entered chan struct{}
	release chan struct{}
}

func (*lastValueCollectionBarrier) Offer(context.Context, int64, lazyFilteredAttributes) {}

func (r *lastValueCollectionBarrier) Collect(*[]exemplar.Exemplar) {
	if r.calls.Add(1) == 2 {
		close(r.entered)
		<-r.release
	}
}

func TestFinishLastValueConcurrentSafe(t *testing.T) {
	t.Run("Int64", testFinishLastValueConcurrentSafe[int64])
	t.Run("Float64", testFinishLastValueConcurrentSafe[float64])
}

func testFinishLastValueConcurrentSafe[N int64 | float64](t *testing.T) {
	t.Helper()
	for _, temporality := range []metricdata.Temporality{
		metricdata.CumulativeTemporality, metricdata.DeltaTemporality,
	} {
		t.Run(temporality.String(), func(t *testing.T) {
			agg := Builder[N]{Temporality: temporality, AggregationLimit: 3}.FinishLastValue()
			done := make(chan struct{})
			var wg sync.WaitGroup
			wg.Go(func() {
				attrs := attribute.NewSet(attribute.String(keyUser, "0"))
				for {
					select {
					case <-done:
						return
					default:
						agg.Finish(attrs.Equivalent(), time.Now())
					}
				}
			})
			defer func() { close(done); wg.Wait() }()
			testAggregationConcurrentSafe(agg.Measure, agg.ComputeAggregation, validateGauge[N])(t)
		})
	}
}

func finishLastValuePoints[N int64 | float64](t *testing.T, compute ComputeAggregation) []metricdata.DataPoint[N] {
	t.Helper()
	var data metricdata.Aggregation
	compute(&data)
	gauge, ok := data.(metricdata.Gauge[N])
	require.True(t, ok, "aggregation type %T", data)
	return gauge.DataPoints
}

func TestFinishLastValueWaitCancellationDuringCreation(t *testing.T) {
	for _, populated := range []bool{false, true} {
		t.Run(fmt.Sprint("populated=", populated), func(t *testing.T) {
			entered, release := make(chan struct{}), make(chan struct{})
			store := newFinishLastValue(metricdata.CumulativeTemporality, 0,
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

func TestFinishLastValueWaitCancellationDuringCollection(t *testing.T) {
	reservoir := &blockingFinishSumCollectionReservoir{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
	store := newFinishLastValue(metricdata.CumulativeTemporality, 0,
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
