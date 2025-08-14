// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"io"
	"strings"
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

	r := new(log.Record)
	r.SetSeverityText("test")
	_ = s.OnEmit(context.Background(), r)

	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []log.Record{*r}, e.records)
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

type writerExporter struct {
	io.Writer
}

func (e *writerExporter) Export(_ context.Context, records []log.Record) error {
	for _, r := range records {
		_, _ = io.WriteString(e.Writer, r.Body().String())
	}
	return nil
}

func (*writerExporter) Shutdown(context.Context) error {
	return nil
}

func (*writerExporter) ForceFlush(context.Context) error {
	return nil
}

func TestSimpleProcessorEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var s log.SimpleProcessor
		ctx := context.Background()
		record := new(log.Record)
		assert.NoError(t, s.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, s.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, s.Shutdown(ctx), "Shutdown")
	})
}

func TestSimpleProcessorConcurrentSafe(*testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := context.Background()
	e := &writerExporter{new(strings.Builder)}
	s := log.NewSimpleProcessor(e)
	for range goRoutineN {
		go func() {
			defer wg.Done()

			_ = s.OnEmit(ctx, r)
			_ = s.Shutdown(ctx)
			_ = s.ForceFlush(ctx)
		}()
	}

	wg.Wait()
}

func BenchmarkSimpleProcessorOnEmit(b *testing.B) {
	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := context.Background()
	s := log.NewSimpleProcessor(nil)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var out error

		for pb.Next() {
			out = s.OnEmit(ctx, r)
		}

		_ = out
	})
}
