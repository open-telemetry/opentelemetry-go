// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/log"
)

type exporter struct {
	records []log.Record

	exportCalled     bool
	shutdownCalled   bool
	forceFlushCalled bool
}

func (e *exporter) Export(_ context.Context, r []log.Record) error {
	e.records = r
	e.exportCalled = true
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	e.shutdownCalled = true
	return nil
}

func (e *exporter) ForceFlush(context.Context) error {
	e.forceFlushCalled = true
	return nil
}

func TestSimpleProcessorOnEmit(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)

	var r log.Record
	r.SetSeverityText("test")
	_ = s.OnEmit(context.Background(), r)

	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []log.Record{r}, e.records)
}

func TestSimpleProcessorEnabled(t *testing.T) {
	s := log.NewSimpleProcessor(nil)
	assert.True(t, s.Enabled(context.Background(), log.Record{}))
}

func TestSimpleProcessorShutdown(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.Shutdown(context.Background())
	require.True(t, e.shutdownCalled, "exporter Shutdown not called")
}

func TestSimpleProcessorForceFlush(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.ForceFlush(context.Background())
	require.True(t, e.forceFlushCalled, "exporter ForceFlush not called")
}

func TestSimpleProcessorConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	var r log.Record
	r.SetSeverityText("test")
	ctx := context.Background()
	s := log.NewSimpleProcessor(nil)
	for i := 0; i < goRoutineN; i++ {
		go func() {
			defer wg.Done()

			_ = s.OnEmit(ctx, r)
			_ = s.Enabled(ctx, r)
			_ = s.Shutdown(ctx)
			_ = s.ForceFlush(ctx)
		}()
	}

	wg.Wait()
}

func BenchmarkSimpleProcessorOnEmit(b *testing.B) {
	var r log.Record
	r.SetSeverityText("test")
	ctx := context.Background()
	s := log.NewSimpleProcessor(nil)

	var out error

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out = s.OnEmit(ctx, r)
	}

	_ = out
}
