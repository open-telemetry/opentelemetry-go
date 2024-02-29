// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/otest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExporterClientConcurrentSafe(t *testing.T) {
	const goroutines = 5

	coll, err := otest.NewHTTPCollector("", nil)
	require.NoError(t, err)

	ctx := context.Background()
	addr := coll.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)
	client, err := newClient(cfg)
	require.NoError(t, err)

	exp, err := newExporter(client, oconf.Config{})
	require.NoError(t, err)
	rm := new(metricdata.ResourceMetrics)

	done := make(chan struct{})
	var wg, someWork sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		someWork.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, exp.Export(ctx, rm))
			assert.NoError(t, exp.ForceFlush(ctx))

			// Ensure some work is done before shutting down.
			someWork.Done()

			for {
				_ = exp.Export(ctx, rm)
				_ = exp.ForceFlush(ctx)

				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	someWork.Wait()
	assert.NoError(t, exp.Shutdown(ctx))
	assert.ErrorIs(t, exp.Shutdown(ctx), errShutdown)

	close(done)
	wg.Wait()
}

func TestExporterDoesNotBlockTemporalityAndAggregation(t *testing.T) {
	rCh := make(chan otest.ExportResult, 1)
	coll, err := otest.NewHTTPCollector("", rCh)
	require.NoError(t, err)

	ctx := context.Background()
	addr := coll.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)
	client, err := newClient(cfg)
	require.NoError(t, err)

	exp, err := newExporter(client, oconf.Config{})
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rm := new(metricdata.ResourceMetrics)
		t.Log("starting export")
		require.NoError(t, exp.Export(ctx, rm))
		t.Log("export complete")
	}()

	assert.Eventually(t, func() bool {
		const inst = metric.InstrumentKindCounter
		// These should not be blocked.
		t.Log("getting temporality")
		_ = exp.Temporality(inst)
		t.Log("getting aggregation")
		_ = exp.Aggregation(inst)
		return true
	}, time.Second, 10*time.Millisecond)

	// Clear the export.
	rCh <- otest.ExportResult{}
	close(rCh)
	wg.Wait()
}
