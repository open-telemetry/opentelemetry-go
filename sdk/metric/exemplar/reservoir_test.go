// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
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

		ctx := context.Background()

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

			r.Offer(context.Background(), staticTime, NewValue(N(10)), nil)

			dest := []Exemplar{{}} // Should be reset to empty.
			r.Collect(&dest)
			assert.Empty(t, dest, "no exemplars should be collected")
		})
	}
}
