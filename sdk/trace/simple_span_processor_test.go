// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	tid, _ = trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ = trace.SpanIDFromHex("0102040810203040")
)

type testExporter struct {
	spans    []sdktrace.ReadOnlySpan
	shutdown bool
}

func (t *testExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	t.spans = append(t.spans, spans...)
	return nil
}

func (t *testExporter) Shutdown(ctx context.Context) error {
	t.shutdown = true
	select {
	case <-ctx.Done():
		// Ensure context deadline tests receive the expected error.
		return ctx.Err()
	default:
		return nil
	}
}

var _ sdktrace.SpanExporter = (*testExporter)(nil)

func TestNewSimpleSpanProcessor(t *testing.T) {
	if ssp := sdktrace.NewSimpleSpanProcessor(&testExporter{}); ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor")
	}
}

func TestNewSimpleSpanProcessorWithNilExporter(t *testing.T) {
	if ssp := sdktrace.NewSimpleSpanProcessor(nil); ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor with nil exporter")
	}
}

func startSpan(tp trace.TracerProvider) trace.Span {
	tr := tp.Tracer("SimpleSpanProcessor")
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
	_, span := tr.Start(ctx, "OnEnd")
	return span
}

func TestSimpleSpanProcessorOnEnd(t *testing.T) {
	tp := basicTracerProvider(t)
	te := testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(&te)

	tp.RegisterSpanProcessor(ssp)
	startSpan(tp).End()

	wantTraceID := tid
	gotTraceID := te.spans[0].SpanContext().TraceID()
	if wantTraceID != gotTraceID {
		t.Errorf("SimplerSpanProcessor OnEnd() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}

func TestSimpleSpanProcessorShutdown(t *testing.T) {
	exporter := &testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(exporter)

	// Ensure we can export a span before we test we cannot after shutdown.
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)
	startSpan(tp).End()
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

	startSpan(tp).End()
	if len(exporter.spans) > nExported {
		t.Error("exported span to shutdown exporter")
	}
}

func TestSimpleSpanProcessorShutdownOnEndConcurrentSafe(t *testing.T) {
	exporter := &testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(exporter)
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
				startSpan(tp).End()
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
	exporter := &testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(exporter)
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

	ssp := sdktrace.NewSimpleSpanProcessor(&testExporter{})
	if got, want := ssp.Shutdown(ctx), context.DeadlineExceeded; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}

func TestSimpleSpanProcessorShutdownHonorsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ssp := sdktrace.NewSimpleSpanProcessor(&testExporter{})
	if got, want := ssp.Shutdown(ctx), context.Canceled; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}
