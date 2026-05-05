// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/trace"
)

var staticTime = time.Unix(946684800, 0)

func TestFixedSizeRoundRobinReservoir(t *testing.T) {
	t.Run("CaptureSpanContext", func(t *testing.T) {
		r := NewFixedSizeRoundRobinReservoir(1)
		ctx := t.Context()

		tID, sID := trace.TraceID{0x01}, trace.SpanID{0x01}
		sc := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tID,
			SpanID:     sID,
			TraceFlags: trace.FlagsSampled,
		})
		ctx = trace.ContextWithSpanContext(ctx, sc)

		r.Offer(ctx, staticTime, exemplar.NewValue[int64](10), nil)

		var dest []exemplar.Exemplar
		r.Collect(&dest)

		want := exemplar.Exemplar{
			Time:    staticTime,
			Value:   exemplar.NewValue[int64](10),
			SpanID:  sID[:],
			TraceID: tID[:],
		}
		require.Len(t, dest, 1)
		assert.Equal(t, want, dest[0])
	})

	t.Run("FilterAttributes", func(t *testing.T) {
		r := NewFixedSizeRoundRobinReservoir(1)
		ctx := t.Context()

		adminTrue := attribute.Bool("admin", true)
		r.Offer(ctx, staticTime, exemplar.NewValue[int64](10), []attribute.KeyValue{adminTrue})

		var dest []exemplar.Exemplar
		r.Collect(&dest)

		want := exemplar.Exemplar{
			FilteredAttributes: []attribute.KeyValue{adminTrue},
			Time:               staticTime,
			Value:              exemplar.NewValue[int64](10),
		}
		require.Len(t, dest, 1)
		assert.Equal(t, want, dest[0])
	})

	t.Run("CollectLessThanN", func(t *testing.T) {
		r := NewFixedSizeRoundRobinReservoir(2)
		ctx := t.Context()

		r.Offer(ctx, staticTime, exemplar.NewValue[int64](10), nil)

		var dest []exemplar.Exemplar
		r.Collect(&dest)
		require.Len(t, dest, 1)
	})

	t.Run("MultipleOffers", func(t *testing.T) {
		n := 3
		r := NewFixedSizeRoundRobinReservoir(n)
		ctx := t.Context()

		for i := 0; i < n+1; i++ {
			r.Offer(ctx, staticTime, exemplar.NewValue[int64](int64(i)), nil)
		}

		var dest []exemplar.Exemplar
		r.Collect(&dest)
		assert.Len(t, dest, n, "multiple offers did not fill reservoir")

		// Ensure the collect reset also resets any counting state.
		for i := 0; i < n+1; i++ {
			r.Offer(ctx, staticTime, exemplar.NewValue[int64](int64(i)), nil)
		}

		dest = dest[:0]
		r.Collect(&dest)
		assert.Len(t, dest, n, "internal count state not reset")
	})

	t.Run("DropAll", func(t *testing.T) {
		r := NewFixedSizeRoundRobinReservoir(0)
		ctx := t.Context()

		r.Offer(ctx, staticTime, exemplar.NewValue[int64](10), nil)

		dest := []exemplar.Exemplar{{}} // Should be reset to empty.
		r.Collect(&dest)
		assert.Empty(t, dest, "no exemplars should be collected")
	})
}

func TestFixedSizeRoundRobinReservoirDistribution(t *testing.T) {
	size := 3
	r := NewFixedSizeRoundRobinReservoir(size)

	ctx := t.Context()
	now := time.Now()

	offers := 900
	for i := range offers {
		r.Offer(ctx, now, exemplar.NewValue[int64](int64(i)), nil)
	}

	var dest []exemplar.Exemplar
	r.Collect(&dest)

	assert.Len(t, dest, size, "Should have filled the reservoir")
}

func TestResetHelperReuseSlice(t *testing.T) {
	dest := make([]exemplar.Exemplar, 2, 5)
	res := reset(dest, 1, 3)
	assert.Len(t, res, 1)
	assert.Equal(t, 5, cap(res))
}
