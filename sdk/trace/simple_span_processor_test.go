// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestSimpleSpanProcessorShutdownOnEndConcurrency(t *testing.T) {
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
