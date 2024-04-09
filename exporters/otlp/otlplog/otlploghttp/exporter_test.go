// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/log"
)

func TestExporterShutdown(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown Exporter")

	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	r := make([]log.Record, 1)
	assert.NoError(t, e.Export(ctx, r), "Export on Shutdown Exporter")
	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush on Shutdown Exporter")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown on Shutdown Exporter")
}

func TestExporterForceFlush(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")

	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush")
}

func TestExporterConcurrentSafe(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "newExporter")

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	runs := new(uint64)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := make([]log.Record, 1)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					atomic.AddUint64(runs, 1)
					_ = e.Export(ctx, r)
					_ = e.ForceFlush(ctx)
				}
			}
		}()
	}

	for atomic.LoadUint64(runs) == 0 {
		runtime.Gosched()
	}

	_ = e.Shutdown(ctx)
	cancel()
	wg.Wait()
}
