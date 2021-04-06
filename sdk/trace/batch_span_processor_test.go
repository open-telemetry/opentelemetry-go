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
	"encoding/binary"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testBatchExporter struct {
	mu            sync.Mutex
	spans         []*sdktrace.SpanSnapshot
	sizes         []int
	batchCount    int
	shutdownCount int
	delay         time.Duration
	err           error
}

func (t *testBatchExporter) ExportSpans(ctx context.Context, ss []*sdktrace.SpanSnapshot) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	time.Sleep(t.delay)

	select {
	case <-ctx.Done():
		t.err = ctx.Err()
		return ctx.Err()
	default:
	}

	t.spans = append(t.spans, ss...)
	t.sizes = append(t.sizes, len(ss))
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

func TestNewBatchSpanProcessorWithNilExporter(t *testing.T) {
	tp := basicTracerProvider(t)
	bsp := sdktrace.NewBatchSpanProcessor(nil)
	tp.RegisterSpanProcessor(bsp)
	tr := tp.Tracer("NilExporter")

	_, span := tr.Start(context.Background(), "foo")
	span.End()

	// These should not panic.
	bsp.OnStart(context.Background(), span.(sdktrace.ReadWriteSpan))
	bsp.OnEnd(span.(sdktrace.ReadOnlySpan))
	if err := bsp.ForceFlush(context.Background()); err != nil {
		t.Errorf("failed to ForceFlush the BatchSpanProcessor: %v", err)
	}
	if err := bsp.Shutdown(context.Background()); err != nil {
		t.Errorf("failed to Shutdown the BatchSpanProcessor: %v", err)
	}
}

type testOption struct {
	name              string
	o                 []sdktrace.BatchSpanProcessorOption
	wantNumSpans      int
	wantBatchCount    int
	wantExportTimeout bool
	genNumSpans       int
	delayExportBy     time.Duration
	parallel          bool
}

func TestNewBatchSpanProcessorWithOptions(t *testing.T) {
	schDelay := 200 * time.Millisecond
	exportTimeout := time.Millisecond
	options := []testOption{
		{
			name:           "default BatchSpanProcessorOptions",
			wantNumSpans:   2053,
			wantBatchCount: 4,
			genNumSpans:    2053,
		},
		{
			name: "non-default ExportTimeout",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithExportTimeout(exportTimeout),
			},
			wantExportTimeout: true,
			genNumSpans:       2053,
			delayExportBy:     2 * exportTimeout, // to ensure export timeout
		},
		{
			name: "non-default BatchTimeout",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
			},
			wantNumSpans:   2053,
			wantBatchCount: 4,
			genNumSpans:    2053,
		},
		{
			name: "non-default MaxQueueSize and BatchTimeout",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
				sdktrace.WithMaxQueueSize(200),
			},
			wantNumSpans:   205,
			wantBatchCount: 1,
			genNumSpans:    205,
		},
		{
			name: "non-default MaxQueueSize, BatchTimeout and MaxExportBatchSize",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
				sdktrace.WithMaxQueueSize(205),
				sdktrace.WithMaxExportBatchSize(20),
			},
			wantNumSpans:   210,
			wantBatchCount: 11,
			genNumSpans:    210,
		},
		{
			name: "blocking option",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
				sdktrace.WithMaxQueueSize(200),
				sdktrace.WithMaxExportBatchSize(20),
			},
			wantNumSpans:   205,
			wantBatchCount: 11,
			genNumSpans:    205,
		},
		{
			name: "parallel span generation",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
				sdktrace.WithMaxQueueSize(200),
			},
			wantNumSpans:   205,
			wantBatchCount: 1,
			genNumSpans:    205,
			parallel:       true,
		},
		{
			name: "parallel span blocking",
			o: []sdktrace.BatchSpanProcessorOption{
				sdktrace.WithBatchTimeout(schDelay),
				sdktrace.WithMaxExportBatchSize(200),
			},
			wantNumSpans:   2000,
			wantBatchCount: 10,
			genNumSpans:    2000,
			parallel:       true,
		},
	}
	for _, option := range options {
		t.Run(option.name, func(t *testing.T) {
			te := testBatchExporter{
				delay: option.delayExportBy,
			}
			tp := basicTracerProvider(t)
			ssp := createAndRegisterBatchSP(option, &te)
			if ssp == nil {
				t.Fatalf("%s: Error creating new instance of BatchSpanProcessor\n", option.name)
			}
			tp.RegisterSpanProcessor(ssp)
			tr := tp.Tracer("BatchSpanProcessorWithOptions")

			generateSpan(t, option.parallel, tr, option)

			tp.UnregisterSpanProcessor(ssp)

			gotNumOfSpans := te.len()
			if option.wantNumSpans > 0 && option.wantNumSpans != gotNumOfSpans {
				t.Errorf("number of exported span: got %+v, want %+v\n",
					gotNumOfSpans, option.wantNumSpans)
			}

			gotBatchCount := te.getBatchCount()
			if option.wantBatchCount > 0 && gotBatchCount < option.wantBatchCount {
				t.Errorf("number batches: got %+v, want >= %+v\n",
					gotBatchCount, option.wantBatchCount)
				t.Errorf("Batches %v\n", te.sizes)
			}

			if option.wantExportTimeout && te.err != context.DeadlineExceeded {
				t.Errorf("context deadline: got err %+v, want %+v\n",
					te.err, context.DeadlineExceeded)
			}
		})
	}
}

func createAndRegisterBatchSP(option testOption, te *testBatchExporter) sdktrace.SpanProcessor {
	// Always use blocking queue to avoid flaky tests.
	options := append(option.o, sdktrace.WithBlocking())
	return sdktrace.NewBatchSpanProcessor(te, options...)
}

func generateSpan(t *testing.T, parallel bool, tr trace.Tracer, option testOption) {
	sc := getSpanContext()

	wg := &sync.WaitGroup{}
	for i := 0; i < option.genNumSpans; i++ {
		tid := sc.TraceID()
		binary.BigEndian.PutUint64(tid[0:8], uint64(i+1))
		newSc := sc.WithTraceID(tid)

		wg.Add(1)
		f := func(sc trace.SpanContext) {
			ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
			_, span := tr.Start(ctx, option.name)
			span.End()
			wg.Done()
		}
		if parallel {
			go f(newSc)
		} else {
			f(newSc)
		}
	}
	wg.Wait()
}

func getSpanContext() trace.SpanContext {
	tid, _ := trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ := trace.SpanIDFromHex("0102040810203040")
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	})
}

func TestBatchSpanProcessorShutdown(t *testing.T) {
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

func TestBatchSpanProcessorPostShutdown(t *testing.T) {
	tp := basicTracerProvider(t)
	be := testBatchExporter{}
	bsp := sdktrace.NewBatchSpanProcessor(&be)

	tp.RegisterSpanProcessor(bsp)
	tr := tp.Tracer("Normal")

	generateSpan(t, true, tr, testOption{
		o: []sdktrace.BatchSpanProcessorOption{
			sdktrace.WithMaxExportBatchSize(50),
		},
		genNumSpans: 60,
	})

	require.NoError(t, bsp.Shutdown(context.Background()), "shutting down BatchSpanProcessor")
	lenJustAfterShutdown := be.len()

	_, span := tr.Start(context.Background(), "foo")
	span.End()
	assert.NoError(t, bsp.ForceFlush(context.Background()), "force flushing BatchSpanProcessor")

	assert.Equal(t, lenJustAfterShutdown, be.len(), "OnEnd and ForceFlush should have no effect after Shutdown")
}

func TestBatchSpanProcessorForceFlushSucceeds(t *testing.T) {
	te := testBatchExporter{}
	tp := basicTracerProvider(t)
	option := testOption{
		name: "default BatchSpanProcessorOptions",
		o: []sdktrace.BatchSpanProcessorOption{
			sdktrace.WithMaxQueueSize(0),
			sdktrace.WithMaxExportBatchSize(3000),
		},
		wantNumSpans:   2053,
		wantBatchCount: 1,
		genNumSpans:    2053,
	}
	ssp := createAndRegisterBatchSP(option, &te)
	if ssp == nil {
		t.Fatalf("%s: Error creating new instance of BatchSpanProcessor\n", option.name)
	}
	tp.RegisterSpanProcessor(ssp)
	tr := tp.Tracer("BatchSpanProcessorWithOption")
	generateSpan(t, option.parallel, tr, option)

	// Force flush any held span batches
	err := ssp.ForceFlush(context.Background())

	gotNumOfSpans := te.len()
	spanDifference := option.wantNumSpans - gotNumOfSpans
	if spanDifference > 10 || spanDifference < 0 {
		t.Errorf("number of exported span not equal to or within 10 less than: got %+v, want %+v\n",
			gotNumOfSpans, option.wantNumSpans)
	}
	gotBatchCount := te.getBatchCount()
	if gotBatchCount < option.wantBatchCount {
		t.Errorf("number batches: got %+v, want >= %+v\n",
			gotBatchCount, option.wantBatchCount)
		t.Errorf("Batches %v\n", te.sizes)
	}
	assert.NoError(t, err)
}

func TestBatchSpanProcessorForceFlushTimeout(t *testing.T) {
	var bp testBatchExporter
	bsp := sdktrace.NewBatchSpanProcessor(&bp)
	// Add timeout to context to test deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	if err := bsp.ForceFlush(ctx); err == nil {
		t.Error("expected context DeadlineExceeded error, got nil")
	} else if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context DeadlineExceeded error, got %v", err)
	}
}

func TestBatchSpanProcessorForceFlushCancellation(t *testing.T) {
	var bp testBatchExporter
	bsp := sdktrace.NewBatchSpanProcessor(&bp)
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context
	cancel()

	if err := bsp.ForceFlush(ctx); err == nil {
		t.Error("expected context canceled error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error, got %v", err)
	}
}
