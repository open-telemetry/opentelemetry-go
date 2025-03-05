// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testBatchExporter struct {
	mu            sync.Mutex
	spans         []sdktrace.ReadOnlySpan
	sizes         []int
	batchCount    int
	shutdownCount int
	errors        []error
	droppedCount  int
	idx           int
	err           error
}

func (t *testBatchExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.idx < len(t.errors) {
		t.droppedCount += len(spans)
		err := t.errors[t.idx]
		t.idx++
		return err
	}

	select {
	case <-ctx.Done():
		t.err = ctx.Err()
		return ctx.Err()
	default:
	}

	t.spans = append(t.spans, spans...)
	t.sizes = append(t.sizes, len(spans))
	t.batchCount++
	return nil
}

func (t *testBatchExporter) Shutdown(context.Context) error {
	t.shutdownCount++
	return nil
}

func (t *testBatchExporter) len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.spans)
}

func (t *testBatchExporter) getBatchCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.batchCount
}

var _ sdktrace.SpanExporter = (*testBatchExporter)(nil)

func TestBatchSpanProcessorShutdownSelfObs(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	// TODO: setup meterprovider and verify metrics produced
	var bp testBatchExporter
	bsp := sdktrace.NewBatchSpanProcessor(&bp)

	err := bsp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the BatchSpanProcessor down\n")
	}
	assert.Equal(t, 1, bp.shutdownCount, "shutdown from span exporter not called")

	// Multiple call to Shutdown() should not panic.
	err = bsp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the BatchSpanProcessor down\n")
	}
	assert.Equal(t, 1, bp.shutdownCount)
}
