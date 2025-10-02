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
				wg.Add(1)
				go func() {
					reservoir.Offer(t.Context(), 25, []attribute.KeyValue{})
					wg.Done()
				}()
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
