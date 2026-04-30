// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestFixedSizeExemplarConcurrentSafe(t *testing.T) {
	// Tests https://github.com/open-telemetry/opentelemetry-go/issues/5814

	t.Setenv("OTEL_METRICS_EXEMPLAR_FILTER", "always_on")

	r := NewManualReader()
	m := NewMeterProvider(WithReader(r)).Meter("exemplar-concurrency")
	// Use two instruments to get concurrent access to any shared globals.
	i0, err := m.Int64Counter("counter.0")
	require.NoError(t, err)
	i1, err := m.Int64Counter("counter.1")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())

	add := func() {
		i0.Add(ctx, 1)
		i1.Add(ctx, 2)
	}

	goRoutines := max(10, runtime.NumCPU())

	var wg sync.WaitGroup
	for range goRoutines {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					require.NotPanics(t, add)
				}
			}
		})
	}

	const collections = 100
	var rm metricdata.ResourceMetrics
	for range collections {
		require.NotPanics(t, func() { _ = r.Collect(ctx, &rm) })
	}

	cancel()
	wg.Wait()
}

func TestReservoirFunc(t *testing.T) {
	type testCase struct {
		name       string
		kind       InstrumentKind
		filter     exemplar.Filter
		expectDrop bool
	}

	testCases := []testCase{
		{
			name:       "AlwaysOff",
			kind:       InstrumentKindCounter,
			filter:     exemplar.AlwaysOffFilter,
			expectDrop: true,
		},
		{
			name:       "AlwaysOn",
			kind:       InstrumentKindCounter,
			filter:     exemplar.AlwaysOnFilter,
			expectDrop: false,
		},
		{
			name:       "TraceBasedSync",
			kind:       InstrumentKindCounter,
			filter:     exemplar.TraceBasedFilter,
			expectDrop: false,
		},
		{
			name:       "TraceBasedAsyncCounter",
			kind:       InstrumentKindObservableCounter,
			filter:     exemplar.TraceBasedFilter,
			expectDrop: true,
		},
		{
			name:       "TraceBasedAsyncUpDownCounter",
			kind:       InstrumentKindObservableUpDownCounter,
			filter:     exemplar.TraceBasedFilter,
			expectDrop: true,
		},
		{
			name:       "TraceBasedAsyncGauge",
			kind:       InstrumentKindObservableGauge,
			filter:     exemplar.TraceBasedFilter,
			expectDrop: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var invoked bool
			provider := func(attribute.Set) exemplar.Reservoir {
				invoked = true
				return nil
			}

			f := reservoirFunc[int64](tc.kind, provider, tc.filter)
			_ = f(*attribute.EmptySet())

			if tc.expectDrop {
				require.False(t, invoked, "ReservoirProvider should not be invoked")
			} else {
				require.True(t, invoked, "ReservoirProvider should be invoked")
			}
		})
	}
}
