// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"reflect"
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

func TestReservoirFuncAlwaysOff(t *testing.T) {
	f := reservoirFunc[int64](
		exemplar.FixedSizeReservoirProvider(1),
		exemplar.AlwaysOffFilter,
	)

	res := f(*attribute.EmptySet())
	require.Contains(
		t,
		reflect.TypeOf(res).String(),
		"dropRes",
		"reservoirFunc should return DropReservoir when AlwaysOffFilter is used",
	)
}
