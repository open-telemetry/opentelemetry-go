// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

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
		expectOfferLazy      bool
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
		{
			desc:                 "offer lazy",
			reservoir:            &offerLazyReservoir{},
			expectConcurrentSafe: true,
			expectOfferLazy:      true,
		},
		{
			desc:                 "offer lazy not concurrent safe",
			reservoir:            &notConcurrentSafeOfferLazyReservoir{},
			expectConcurrentSafe: false,
			expectOfferLazy:      true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			reservoir := NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, tc.reservoir)
			var wg sync.WaitGroup
			for range 5 {
				wg.Go(func() {
					reservoir.Offer(t.Context(), 25, *attribute.EmptySet(), nil)
				})
			}
			into := []exemplar.Exemplar{}
			for range 2 {
				reservoir.Collect(&into)
			}
			wg.Wait()
			assert.Len(t, into, 1)
			assert.Equal(t, reservoir.(*filteredExemplarReservoir[int64]).concurrentSafe, tc.expectConcurrentSafe)

			if tc.expectOfferLazy {
				if r, ok := tc.reservoir.(offerLazyReporter); ok {
					assert.True(t, r.OfferLazyCalled())
				} else {
					t.Fatal("reservoir does not implement offerLazyReporter")
				}
			}
		})
	}
}

type notConcurrentSafeReservoir struct {
	ex exemplar.Exemplar
}

func (r *notConcurrentSafeReservoir) Offer(
	_ context.Context,
	t time.Time,
	val exemplar.Value,
	attr []attribute.KeyValue,
) {
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

type offerLazyReporter interface {
	OfferLazyCalled() bool
}

type offerLazyReservoir struct {
	concurrentSafeReservoir
	offerLazyCalled bool
}

func (r *offerLazyReservoir) OfferLazyCalled() bool {
	return r.offerLazyCalled
}

func (r *offerLazyReservoir) OfferLazy(
	ctx context.Context,
	t time.Time,
	val exemplar.Value,
	attr attribute.Set,
	fltr attribute.Filter,
) {
	r.Lock()
	defer r.Unlock()
	r.offerLazyCalled = true
	r.base.Offer(ctx, t, val, getDroppedAttributes(attr, fltr))
}

type notConcurrentSafeOfferLazyReservoir struct {
	notConcurrentSafeReservoir
	offerLazyCalled bool
}

func (r *notConcurrentSafeOfferLazyReservoir) OfferLazyCalled() bool {
	return r.offerLazyCalled
}

func (r *notConcurrentSafeOfferLazyReservoir) OfferLazy(
	_ context.Context,
	t time.Time,
	val exemplar.Value,
	attr attribute.Set,
	fltr attribute.Filter,
) {
	r.offerLazyCalled = true
	r.ex = exemplar.Exemplar{
		FilteredAttributes: getDroppedAttributes(attr, fltr),
		Time:               t,
		Value:              val,
	}
}
