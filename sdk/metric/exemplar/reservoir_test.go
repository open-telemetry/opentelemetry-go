// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Sat Jan 01 2000 00:00:00 GMT+0000.
var staticTime = time.Unix(946684800, 0)

type factory func(requestedCap int) (r ReservoirProvider, actualCap int)

func ReservoirTest[N int64 | float64](f factory) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ctx := t.Context()

		t.Run("CaptureSpanContext", func(t *testing.T) {
			t.Helper()

			rp, n := f(1)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}
			r := rp(*attribute.EmptySet())

			tID, sID := trace.TraceID{0x01}, trace.SpanID{0x01}
			sc := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    tID,
				SpanID:     sID,
				TraceFlags: trace.FlagsSampled,
			})
			ctx := trace.ContextWithSpanContext(ctx, sc)

			r.Offer(ctx, staticTime, NewValue(N(10)), nil)

			var dest []Exemplar
			r.Collect(&dest)

			want := Exemplar{
				Time:    staticTime,
				Value:   NewValue(N(10)),
				SpanID:  sID[:],
				TraceID: tID[:],
			}
			require.Len(t, dest, 1, "number of collected exemplars")
			assert.Equal(t, want, dest[0])
		})

		t.Run("FilterAttributes", func(t *testing.T) {
			t.Helper()

			rp, n := f(1)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}
			r := rp(*attribute.EmptySet())

			adminTrue := attribute.Bool("admin", true)
			r.Offer(ctx, staticTime, NewValue(N(10)), []attribute.KeyValue{adminTrue})

			var dest []Exemplar
			r.Collect(&dest)

			want := Exemplar{
				FilteredAttributes: []attribute.KeyValue{adminTrue},
				Time:               staticTime,
				Value:              NewValue(N(10)),
			}
			require.Len(t, dest, 1, "number of collected exemplars")
			assert.Equal(t, want, dest[0])
		})

		t.Run("CollectLessThanN", func(t *testing.T) {
			t.Helper()

			rp, n := f(2)
			if n < 2 {
				t.Skip("skipping, reservoir capacity less than 2:", n)
			}
			r := rp(*attribute.EmptySet())

			r.Offer(ctx, staticTime, NewValue(N(10)), nil)

			var dest []Exemplar
			r.Collect(&dest)
			// No empty exemplars are exported.
			require.Len(t, dest, 1, "number of collected exemplars")
		})

		t.Run("MultipleOffers", func(t *testing.T) {
			t.Helper()

			rp, n := f(3)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}
			r := rp(*attribute.EmptySet())

			for i := 0; i < n+1; i++ {
				v := NewValue(N(i))
				r.Offer(ctx, staticTime, v, nil)
			}

			var dest []Exemplar
			r.Collect(&dest)
			assert.Len(t, dest, n, "multiple offers did not fill reservoir")

			// Ensure the collect reset also resets any counting state.
			for i := 0; i < n+1; i++ {
				v := NewValue(N(i))
				r.Offer(ctx, staticTime, v, nil)
			}

			dest = dest[:0]
			r.Collect(&dest)
			assert.Len(t, dest, n, "internal count state not reset")
		})

		t.Run("DropAll", func(t *testing.T) {
			t.Helper()

			rp, n := f(0)
			if n > 0 {
				t.Skip("skipping, reservoir capacity greater than 0:", n)
			}
			r := rp(*attribute.EmptySet())

			r.Offer(t.Context(), staticTime, NewValue(N(10)), nil)

			dest := []Exemplar{{}} // Should be reset to empty.
			r.Collect(&dest)
			assert.Empty(t, dest, "no exemplars should be collected")
		})
	}
}

func reservoirConcurrentSafeTest[N int64 | float64](f factory) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		rp, n := f(1)
		if n < 1 {
			t.Skip("skipping, reservoir capacity less than 1:", n)
		}
		r := rp(*attribute.EmptySet())

		var wg sync.WaitGroup

		const goroutines = 2

		// Call Offer concurrently with another Offer, and with Collect.
		for i := range goroutines {
			wg.Add(1)
			go func(iteration int) {
				ctx, ts, val, attrs := generateOfferInputs[N](iteration + 1)
				r.Offer(ctx, ts, val, attrs)
				wg.Done()
			}(i)
		}

		// Also test concurrent Collect calls
		wg.Go(func() {
			var dest []Exemplar
			r.Collect(&dest)
		})

		wg.Wait()

		// Final collect to validate state
		var dest []Exemplar
		r.Collect(&dest)
		assert.NotEmpty(t, dest)
		for _, e := range dest {
			validateExemplar[N](t, e)
		}
	}
}

func generateOfferInputs[N int64 | float64](
	i int,
) (context.Context, time.Time, Value, []attribute.KeyValue) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID([16]byte{byte(i)}),
		SpanID:     trace.SpanID([8]byte{byte(i)}),
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	ts := time.Unix(int64(i), int64(i))
	val := NewValue(N(i))
	attrs := []attribute.KeyValue{attribute.Int("i", i)}
	return ctx, ts, val, attrs
}

func validateExemplar[N int64 | float64](t *testing.T, e Exemplar) {
	t.Helper()
	i := 0
	switch e.Value.Type() {
	case Int64ValueType:
		i = int(e.Value.Int64())
	case Float64ValueType:
		i = int(e.Value.Float64())
	default:
		t.Fatalf("unexpected value type: %v", e.Value.Type())
	}
	if i == 0 {
		t.Fatal("empty exemplar")
	}
	ctx, ts, _, attrs := generateOfferInputs[N](i)
	sc := trace.SpanContextFromContext(ctx)
	tID := sc.TraceID()
	sID := sc.SpanID()
	assert.Equal(t, tID[:], e.TraceID)
	assert.Equal(t, sID[:], e.SpanID)
	assert.Equal(t, ts, e.Time)
	assert.Equal(t, attrs, e.FilteredAttributes)
}

func TestFixedSizeReservoir_Merge(t *testing.T) {
	ctx := t.Context()
	rCum := NewFixedSizeReservoir(2)
	rDelta := NewFixedSizeReservoir(2)

	rDelta.Offer(ctx, staticTime, NewValue(int64(10)), nil)
	rDelta.Offer(ctx, staticTime, NewValue(int64(20)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	var dest []Exemplar
	rCum.Collect(&dest)
	require.Len(t, dest, 2)
	assert.Equal(t, NewValue(int64(10)), dest[0].Value)
	assert.Equal(t, NewValue(int64(20)), dest[1].Value)

	dest = nil
	rDelta.Collect(&dest)
	assert.Empty(t, dest)

	t2 := staticTime.Add(time.Second)
	// Interval 2: offer 1 new measurement (which overwrites slot 0).
	rDelta.Offer(ctx, t2, NewValue(int64(30)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	dest = nil
	rCum.Collect(&dest)
	require.Len(t, dest, 2)
	assert.Equal(t, NewValue(int64(30)), dest[0].Value)
	assert.Equal(t, NewValue(int64(20)), dest[1].Value)
}

func TestFixedSizeReservoir_Merge_EdgeCases(t *testing.T) {
	ctx := t.Context()

	t.Run("NilSelfAndTypeMismatch", func(t *testing.T) {
		r := NewFixedSizeReservoir(2)
		r.Offer(ctx, staticTime, NewValue(int64(1)), nil)
		assert.NotPanics(t, func() { r.Merge(nil) })
		assert.NotPanics(t, func() { r.Merge((*FixedSizeReservoir)(nil)) })
		assert.NotPanics(t, func() { (*FixedSizeReservoir)(nil).Merge(r) })
		assert.NotPanics(t, func() { r.Merge(r) })
		assert.NotPanics(t, func() { r.Merge(NewHistogramReservoir([]float64{1})) })

		var dest []Exemplar
		r.Collect(&dest)
		require.Len(t, dest, 1)
	})

	t.Run("ZeroCapacity", func(t *testing.T) {
		r0 := NewFixedSizeReservoir(0)
		rDelta := NewFixedSizeReservoir(2)
		rDelta.Offer(ctx, staticTime, NewValue(int64(10)), nil)
		assert.NotPanics(t, func() { r0.Merge(rDelta) })
	})

	t.Run("FullReservoirSamplingBranch", func(t *testing.T) {
		rCum := NewFixedSizeReservoir(2)
		rCum.Offer(ctx, staticTime, NewValue(int64(1)), nil)
		rCum.Offer(ctx, staticTime, NewValue(int64(2)), nil)

		rDelta := NewFixedSizeReservoir(2)
		rDelta.Offer(ctx, staticTime, NewValue(int64(3)), nil)
		rDelta.Offer(ctx, staticTime, NewValue(int64(4)), nil)

		rCum.Merge(rDelta)
		var dest []Exemplar
		rCum.Collect(&dest)
		assert.Len(t, dest, 2)
	})
}

func TestFixedSizeReservoir_Reset(t *testing.T) {
	ctx := t.Context()
	r := NewFixedSizeReservoir(2)
	r.Offer(ctx, staticTime, NewValue(int64(10)), nil)

	r.Reset()

	var dest []Exemplar
	r.Collect(&dest)
	assert.Empty(t, dest)

	r.Offer(ctx, staticTime, NewValue(int64(20)), nil)
	r.Collect(&dest)
	require.Len(t, dest, 1)
	assert.Equal(t, NewValue(int64(20)), dest[0].Value)
}

func TestHistogramReservoir_Merge(t *testing.T) {
	ctx := t.Context()
	bounds := []float64{5, 10}
	rCum := NewHistogramReservoir(bounds)
	rDelta := NewHistogramReservoir(bounds)

	rDelta.Offer(ctx, staticTime, NewValue(int64(2)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	var dest []Exemplar
	rCum.Collect(&dest)
	require.Len(t, dest, 1)
	assert.Equal(t, NewValue(int64(2)), dest[0].Value)

	dest = nil
	rDelta.Collect(&dest)
	assert.Empty(t, dest)

	t2 := staticTime.Add(time.Second)
	rDelta.Offer(ctx, t2, NewValue(int64(7)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	dest = nil
	rCum.Collect(&dest)
	require.Len(t, dest, 2)
	assert.Equal(t, NewValue(int64(2)), dest[0].Value)
	assert.Equal(t, NewValue(int64(7)), dest[1].Value)
}

func TestHistogramReservoir_Merge_EdgeCases(t *testing.T) {
	ctx := t.Context()

	t.Run("NilSelfAndMismatchedBounds", func(t *testing.T) {
		r := NewHistogramReservoir([]float64{5, 10})
		r.Offer(ctx, staticTime, NewValue(int64(2)), nil)
		assert.NotPanics(t, func() { r.Merge(nil) })
		assert.NotPanics(t, func() { r.Merge(r) })
		assert.NotPanics(t, func() { r.Merge(NewHistogramReservoir([]float64{5})) })
		assert.NotPanics(t, func() { r.Merge(NewHistogramReservoir([]float64{10, 20})) })

		var dest []Exemplar
		r.Collect(&dest)
		require.Len(t, dest, 1)
	})

	t.Run("BucketCollision", func(t *testing.T) {
		bounds := []float64{5, 10}
		rCum := NewHistogramReservoir(bounds)
		rDelta := NewHistogramReservoir(bounds)

		rCum.Offer(ctx, staticTime, NewValue(int64(2)), nil)
		rDelta.Offer(ctx, staticTime, NewValue(int64(3)), nil)

		rCum.Merge(rDelta)
		var dest []Exemplar
		rCum.Collect(&dest)
		require.Len(t, dest, 1)
		assert.Equal(t, NewValue(int64(3)), dest[0].Value)
	})
}

func TestHistogramReservoir_Reset(t *testing.T) {
	ctx := t.Context()
	bounds := []float64{5, 10}
	r := NewHistogramReservoir(bounds)
	r.Offer(ctx, staticTime, NewValue(int64(2)), nil)

	r.Reset()

	var dest []Exemplar
	r.Collect(&dest)
	assert.Empty(t, dest)

	r.Offer(ctx, staticTime, NewValue(int64(7)), nil)
	r.Collect(&dest)
	require.Len(t, dest, 1)
	assert.Equal(t, NewValue(int64(7)), dest[0].Value)
}

func TestFixedSizeReservoir_EquivalenceWithSingleReservoir(t *testing.T) {
	ctx := t.Context()
	k := 4

	// 1. Single reservoir across 2 collection intervals.
	rSingle := NewFixedSizeReservoir(k)
	// Interval 1: Make K offer calls.
	t1 := staticTime
	for i := 1; i <= k; i++ {
		rSingle.Offer(ctx, t1, NewValue(int64(i)), nil)
	}
	var destSingle1 []Exemplar
	rSingle.Collect(&destSingle1)
	require.Len(t, destSingle1, k)

	// Interval 2: Make K/2 offer calls.
	t2 := t1.Add(time.Second)
	for i := 1; i <= k/2; i++ {
		rSingle.Offer(ctx, t2, NewValue(int64(10+i)), nil)
	}
	var destSingle2 []Exemplar
	rSingle.Collect(&destSingle2)
	require.Len(t, destSingle2, k)

	// 2. Double-buffered cumulative reservoir setup.
	rCum := NewFixedSizeReservoir(k)
	rDelta := NewFixedSizeReservoir(k)

	// Interval 1: Make K offer calls on delta, merge into cumulative, reset delta, collect cumulative.
	for i := 1; i <= k; i++ {
		rDelta.Offer(ctx, t1, NewValue(int64(i)), nil)
	}
	rCum.Merge(rDelta)
	rDelta.Reset()

	var destDouble1 []Exemplar
	rCum.Collect(&destDouble1)
	assert.Equal(t, destSingle1, destDouble1)

	// Interval 2: Make K/2 offer calls on delta, merge into cumulative, reset delta, collect cumulative.
	for i := 1; i <= k/2; i++ {
		rDelta.Offer(ctx, t2, NewValue(int64(10+i)), nil)
	}
	rCum.Merge(rDelta)
	rDelta.Reset()

	var destDouble2 []Exemplar
	rCum.Collect(&destDouble2)
	assert.Equal(t, destSingle2, destDouble2)
}

func TestHistogramReservoir_EquivalenceWithSingleReservoir(t *testing.T) {
	ctx := t.Context()
	bounds := []float64{5, 10, 15, 20}

	// 1. Single reservoir across 2 collection intervals.
	rSingle := NewHistogramReservoir(bounds)
	t1 := staticTime
	// Interval 1: Offer values in buckets 0, 1, 2, 3.
	rSingle.Offer(ctx, t1, NewValue(int64(2)), nil)
	rSingle.Offer(ctx, t1, NewValue(int64(7)), nil)
	rSingle.Offer(ctx, t1, NewValue(int64(12)), nil)
	rSingle.Offer(ctx, t1, NewValue(int64(17)), nil)

	var destSingle1 []Exemplar
	rSingle.Collect(&destSingle1)
	require.Len(t, destSingle1, 4)

	// Interval 2: Offer values in buckets 0, 1.
	t2 := t1.Add(time.Second)
	rSingle.Offer(ctx, t2, NewValue(int64(3)), nil)
	rSingle.Offer(ctx, t2, NewValue(int64(8)), nil)

	var destSingle2 []Exemplar
	rSingle.Collect(&destSingle2)
	require.Len(t, destSingle2, 4)

	// 2. Double-buffered cumulative reservoir setup.
	rCum := NewHistogramReservoir(bounds)
	rDelta := NewHistogramReservoir(bounds)

	// Interval 1: Offer on delta, merge into cumulative, reset delta, collect.
	rDelta.Offer(ctx, t1, NewValue(int64(2)), nil)
	rDelta.Offer(ctx, t1, NewValue(int64(7)), nil)
	rDelta.Offer(ctx, t1, NewValue(int64(12)), nil)
	rDelta.Offer(ctx, t1, NewValue(int64(17)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	var destDouble1 []Exemplar
	rCum.Collect(&destDouble1)
	assert.Equal(t, destSingle1, destDouble1)

	// Interval 2: Offer on delta, merge into cumulative, reset delta, collect.
	rDelta.Offer(ctx, t2, NewValue(int64(3)), nil)
	rDelta.Offer(ctx, t2, NewValue(int64(8)), nil)

	rCum.Merge(rDelta)
	rDelta.Reset()

	var destDouble2 []Exemplar
	rCum.Collect(&destDouble2)
	assert.Equal(t, destSingle2, destDouble2)
}
