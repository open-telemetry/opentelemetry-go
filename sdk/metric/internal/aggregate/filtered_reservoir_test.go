// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/reservoir"
)

func TestConcurrentSafeFilteredReservoir(t *testing.T) {
	for _, tc := range []struct {
		desc                 string
		reservoir            exemplar.Reservoir
		expectConcurrentSafe bool
	}{
		{
			desc:                 "concurrent safe",
			reservoir:            &concurrentSafeReservoir{},
			expectConcurrentSafe: true,
		},
		{
			desc:                 "not concurrent safe",
			reservoir:            &notConcurrentSafeReservoir{},
			expectConcurrentSafe: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			reservoir := NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, tc.reservoir)
			var wg sync.WaitGroup
			for range 5 {
				wg.Go(func() {
					reservoir.Offer(t.Context(), 25, lazyFilteredAttributes{})
				})
			}
			into := []exemplar.Exemplar{}
			for range 2 {
				reservoir.Collect(&into)
			}
			wg.Wait()
			assert.Len(t, into, 1)
			assert.Equal(t, reservoir.(*filteredExemplarReservoir[int64]).concurrentSafe, tc.expectConcurrentSafe)
		})
	}
}

func TestFilteredExemplarReservoir_Offer(t *testing.T) {
	orig := attribute.NewSet(attribute.String("k1", "v1"), attribute.String("k2", "v2"))
	for _, tc := range []struct {
		desc           string
		filter         exemplar.Filter
		attrFilter     attribute.Filter
		wantOffered    bool
		wantAttributes []attribute.KeyValue
	}{
		{
			desc:           "sampled with dropped attributes",
			filter:         exemplar.AlwaysOnFilter,
			attrFilter:     func(kv attribute.KeyValue) bool { return kv.Key == "k1" },
			wantOffered:    true,
			wantAttributes: []attribute.KeyValue{attribute.String("k2", "v2")},
		},
		{
			desc:           "sampled with no dropped attributes",
			filter:         exemplar.AlwaysOnFilter,
			attrFilter:     nil,
			wantOffered:    true,
			wantAttributes: nil,
		},
		{
			desc:           "not sampled",
			filter:         exemplar.AlwaysOffFilter,
			attrFilter:     nil,
			wantOffered:    false,
			wantAttributes: nil,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			mockRes := &notConcurrentSafeReservoir{}
			res := NewFilteredExemplarReservoir[int64](tc.filter, mockRes)
			lazy := newLazyFilteredAttributes(orig, tc.attrFilter)
			res.Offer(t.Context(), 10, lazy)
			assert.Equal(t, tc.wantOffered, mockRes.offered)
			assert.Equal(t, tc.wantAttributes, mockRes.ex.FilteredAttributes)
		})
	}
}

func TestAggregators_OfferToReservoir(t *testing.T) {
	DropReservoir[int64](attribute.NewSet()).Offer(t.Context(), 1, lazyFilteredAttributes{})

	mockRes := &notConcurrentSafeReservoir{}
	newRes := func(attribute.Set) FilteredExemplarReservoir[int64] {
		return NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, mockRes)
	}
	ctx := t.Context()
	lazy := newLazyFilteredAttributes(attribute.NewSet(attribute.String("k", "v")), nil)

	for _, tc := range []struct {
		desc    string
		measure func()
	}{
		{
			desc:    "DeltaSum",
			measure: func() { newDeltaSum[int64](true, 10, newRes).measure(ctx, 10, lazy) },
		},
		{
			desc:    "DeltaLastValue",
			measure: func() { newDeltaLastValue[int64](10, newRes).measure(ctx, 10, lazy) },
		},
		{
			desc:    "DeltaHistogram",
			measure: func() { newDeltaHistogram[int64]([]float64{1, 10}, false, false, 10, newRes).measure(ctx, 10, lazy) },
		},
		{
			desc: "CumulativeHistogram",
			measure: func() {
				newCumulativeHistogram[int64]([]float64{1, 10}, false, false, 10, newRes).measure(ctx, 10, lazy)
			},
		},
		{
			desc:    "ExponentialHistogram",
			measure: func() { newExponentialHistogram(160, 20, false, false, 2, newRes).measure(ctx, 10, lazy) },
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			mockRes.offered = false
			tc.measure()
			assert.True(t, mockRes.offered)
		})
	}
}

type notConcurrentSafeReservoir struct {
	offered bool
	ex      exemplar.Exemplar
}

func (r *notConcurrentSafeReservoir) Offer(
	_ context.Context,
	t time.Time,
	val exemplar.Value,
	attr []attribute.KeyValue,
) {
	r.offered = true
	r.ex = exemplar.Exemplar{
		FilteredAttributes: attr,
		Time:               t,
		Value:              val,
	}
}

func (r *notConcurrentSafeReservoir) Collect(dest *[]exemplar.Exemplar) {
	*dest = make([]exemplar.Exemplar, 1)
	(*dest)[0].FilteredAttributes = r.ex.FilteredAttributes
	(*dest)[0].Time = r.ex.Time
	(*dest)[0].Value = r.ex.Value
	*dest = (*dest)[:1]
}

type concurrentSafeReservoir struct {
	base notConcurrentSafeReservoir
	sync.Mutex
	reservoir.ConcurrentSafe
}

func (r *concurrentSafeReservoir) Offer(
	ctx context.Context,
	t time.Time,
	val exemplar.Value,
	attr []attribute.KeyValue,
) {
	r.Lock()
	defer r.Unlock()
	r.base.Offer(ctx, t, val, attr)
}

func (r *concurrentSafeReservoir) Collect(dest *[]exemplar.Exemplar) {
	r.Lock()
	defer r.Unlock()
	r.base.Collect(dest)
}
