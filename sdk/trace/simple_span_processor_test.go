// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type simpleTestExporter struct {
	spans    []ReadOnlySpan
	shutdown bool
}

func (t *simpleTestExporter) ExportSpans(ctx context.Context, spans []ReadOnlySpan) error {
	t.spans = append(t.spans, spans...)
	return nil
}

func (t *simpleTestExporter) Shutdown(ctx context.Context) error {
	t.shutdown = true
	select {
	case <-ctx.Done():
		// Ensure context deadline tests receive the expected error.
		return ctx.Err()
	default:
		return nil
	}
}

var _ SpanExporter = (*simpleTestExporter)(nil)

func TestNewSimpleSpanProcessor(t *testing.T) {
	if ssp := NewSimpleSpanProcessor(&simpleTestExporter{}); ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor")
	}
}

func TestNewSimpleSpanProcessorWithNilExporter(t *testing.T) {
	if ssp := NewSimpleSpanProcessor(nil); ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor with nil exporter")
	}
}

func TestSimpleSpanProcessorOnEnd(t *testing.T) {
	tp := basicTracerProvider(t)
	te := simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(&te)

	tp.RegisterSpanProcessor(ssp)
	startSpan(tp, "TestSimpleSpanProcessorOnEnd").End()

	wantTraceID := tid
	gotTraceID := te.spans[0].SpanContext().TraceID()
	if wantTraceID != gotTraceID {
		t.Errorf("SimplerSpanProcessor OnEnd() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}

func TestSimpleSpanProcessorShutdown(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)

	// Ensure we can export a span before we test we cannot after shutdown.
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)
	startSpan(tp, "TestSimpleSpanProcessorShutdown").End()
	nExported := len(exporter.spans)
	if nExported != 1 {
		t.Error("failed to verify span export")
	}

	if err := ssp.Shutdown(context.Background()); err != nil {
		t.Errorf("shutting the SimpleSpanProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleSpanProcessor.Shutdown did not shut down exporter")
	}

	startSpan(tp, "TestSimpleSpanProcessorShutdown").End()
	if len(exporter.spans) > nExported {
		t.Error("exported span to shutdown exporter")
	}
}

func TestSimpleSpanProcessorShutdownOnEndConcurrentSafe(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for {
			select {
			case <-stop:
				return
			default:
				startSpan(tp, "TestSimpleSpanProcessorShutdownOnEndConcurrentSafe").End()
			}
		}
	}()

	if err := ssp.Shutdown(context.Background()); err != nil {
		t.Errorf("shutting the SimpleSpanProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleSpanProcessor.Shutdown did not shut down exporter")
	}

	stop <- struct{}{}
	<-done
}

func TestSimpleSpanProcessorShutdownOnEndConcurrentSafe2(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)

	var wg sync.WaitGroup
	wg.Add(2)

	span := func(spanName string) {
		assert.NotPanics(t, func() {
			defer wg.Done()
			_, span := tp.Tracer("test").Start(context.Background(), spanName)
			span.End()
		})
	}

	go span("test-span-1")
	go span("test-span-2")

	wg.Wait()

	assert.NoError(t, ssp.Shutdown(context.Background()))
	assert.True(t, exporter.shutdown, "exporter shutdown")
}

func TestSimpleSpanProcessorShutdownHonorsContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	ssp := NewSimpleSpanProcessor(&simpleTestExporter{})
	if got, want := ssp.Shutdown(ctx), context.DeadlineExceeded; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}

func TestSimpleSpanProcessorShutdownHonorsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ssp := NewSimpleSpanProcessor(&simpleTestExporter{})
	if got, want := ssp.Shutdown(ctx), context.Canceled; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}
