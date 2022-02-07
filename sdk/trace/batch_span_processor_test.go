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
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	ottest "go.opentelemetry.io/otel/internal/internaltest"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/internal/env"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
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
	name           string
	o              []sdktrace.BatchSpanProcessorOption
	wantNumSpans   int
	wantBatchCount int
	genNumSpans    int
	parallel       bool
	envs           map[string]string
}

func TestNewBatchSpanProcessorWithOptions(t *testing.T) {
	schDelay := 200 * time.Millisecond
	options := []testOption{
		{
			name:           "default BatchSpanProcessorOptions",
			wantNumSpans:   2053,
			wantBatchCount: 4,
			genNumSpans:    2053,
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
			te := testBatchExporter{}
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
		})
	}
}

func TestNewBatchSpanProcessorWithEnvOptions(t *testing.T) {
	options := []testOption{
		{
			name:           "BatchSpanProcessorEnvOptions - Basic",
			wantNumSpans:   2053,
			wantBatchCount: 1,
			genNumSpans:    2053,
			envs: map[string]string{
				env.BatchSpanProcessorMaxQueueSizeKey:       "5000",
				env.BatchSpanProcessorMaxExportBatchSizeKey: "5000",
			},
		},
		{
			name:           "BatchSpanProcessorEnvOptions - A lager max export batch size than queue size",
			wantNumSpans:   2053,
			wantBatchCount: 4,
			genNumSpans:    2053,
			envs: map[string]string{
				env.BatchSpanProcessorMaxQueueSizeKey:       "5000",
				env.BatchSpanProcessorMaxExportBatchSizeKey: "10000",
			},
		},
		{
			name:           "BatchSpanProcessorEnvOptions - A lage max export batch size with a small queue size",
			wantNumSpans:   2053,
			wantBatchCount: 42,
			genNumSpans:    2053,
			envs: map[string]string{
				env.BatchSpanProcessorMaxQueueSizeKey:       "50",
				env.BatchSpanProcessorMaxExportBatchSizeKey: "10000",
			},
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(env.BatchSpanProcessorScheduleDelayKey)
	envStore.Record(env.BatchSpanProcessorExportTimeoutKey)
	envStore.Record(env.BatchSpanProcessorMaxQueueSizeKey)
	envStore.Record(env.BatchSpanProcessorMaxExportBatchSizeKey)

	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	for _, option := range options {
		t.Run(option.name, func(t *testing.T) {
			for k, v := range option.envs {
				require.NoError(t, os.Setenv(k, v))
			}

			te := testBatchExporter{}
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
		})
	}
}

type stuckExporter struct {
	testBatchExporter
}

// ExportSpans waits for ctx to expire and returns that error.
func (e *stuckExporter) ExportSpans(ctx context.Context, _ []sdktrace.ReadOnlySpan) error {
	<-ctx.Done()
	e.err = ctx.Err()
	return ctx.Err()
}

func TestBatchSpanProcessorExportTimeout(t *testing.T) {
	exp := new(stuckExporter)
	bsp := sdktrace.NewBatchSpanProcessor(
		exp,
		// Set a non-zero export timeout so a deadline is set.
		sdktrace.WithExportTimeout(1*time.Microsecond),
		sdktrace.WithBlocking(),
	)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(bsp)

	tr := tp.Tracer("BatchSpanProcessorExportTimeout")
	generateSpan(t, false, tr, testOption{genNumSpans: 1})
	tp.UnregisterSpanProcessor(bsp)

	if exp.err != context.DeadlineExceeded {
		t.Errorf("context deadline error not returned: got %+v", exp.err)
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

	assertMaxSpanDiff(t, te.len(), option.wantNumSpans, 10)

	gotBatchCount := te.getBatchCount()
	if gotBatchCount < option.wantBatchCount {
		t.Errorf("number batches: got %+v, want >= %+v\n",
			gotBatchCount, option.wantBatchCount)
		t.Errorf("Batches %v\n", te.sizes)
	}
	assert.NoError(t, err)
}

func TestBatchSpanProcessorDropBatchIfFailed(t *testing.T) {
	te := testBatchExporter{
		errors: []error{errors.New("fail to export")},
	}
	tp := basicTracerProvider(t)
	option := testOption{
		o: []sdktrace.BatchSpanProcessorOption{
			sdktrace.WithMaxQueueSize(0),
			sdktrace.WithMaxExportBatchSize(2000),
		},
		wantNumSpans:   1000,
		wantBatchCount: 1,
		genNumSpans:    1000,
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
	assert.Error(t, err)
	assert.EqualError(t, err, "fail to export")

	// First flush will fail, nothing should be exported.
	assertMaxSpanDiff(t, te.droppedCount, option.wantNumSpans, 10)
	assert.Equal(t, 0, te.len())
	assert.Equal(t, 0, te.getBatchCount())

	// Generate a new batch, this will succeed
	generateSpan(t, option.parallel, tr, option)

	// Force flush any held span batches
	err = ssp.ForceFlush(context.Background())
	assert.NoError(t, err)

	assertMaxSpanDiff(t, te.len(), option.wantNumSpans, 10)
	gotBatchCount := te.getBatchCount()
	if gotBatchCount < option.wantBatchCount {
		t.Errorf("number batches: got %+v, want >= %+v\n",
			gotBatchCount, option.wantBatchCount)
		t.Errorf("Batches %v\n", te.sizes)
	}
}

func assertMaxSpanDiff(t *testing.T, want, got, maxDif int) {
	spanDifference := want - got
	if spanDifference < 0 {
		spanDifference = spanDifference * -1
	}
	if spanDifference > maxDif {
		t.Errorf("number of exported span not equal to or within %d less than: got %+v, want %+v\n",
			maxDif, got, want)
	}
}

type indefiniteExporter struct{}

func (indefiniteExporter) Shutdown(context.Context) error { return nil }
func (indefiniteExporter) ExportSpans(ctx context.Context, _ []sdktrace.ReadOnlySpan) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestBatchSpanProcessorForceFlushTimeout(t *testing.T) {
	// Add timeout to context to test deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	bsp := sdktrace.NewBatchSpanProcessor(indefiniteExporter{})
	if got, want := bsp.ForceFlush(ctx), context.DeadlineExceeded; !errors.Is(got, want) {
		t.Errorf("expected %q error, got %v", want, got)
	}
}

func TestBatchSpanProcessorForceFlushCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context
	cancel()

	bsp := sdktrace.NewBatchSpanProcessor(indefiniteExporter{})
	if got, want := bsp.ForceFlush(ctx), context.Canceled; !errors.Is(got, want) {
		t.Errorf("expected %q error, got %v", want, got)
	}
}

func TestBatchSpanProcessorForceFlushQueuedSpans(t *testing.T) {
	ctx := context.Background()

	exp := tracetest.NewInMemoryExporter()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)

	tracer := tp.Tracer("tracer")

	for i := 0; i < 10; i++ {
		_, span := tracer.Start(ctx, fmt.Sprintf("span%d", i))
		span.End()

		err := tp.ForceFlush(ctx)
		assert.NoError(t, err)

		assert.Len(t, exp.GetSpans(), i+1)
	}
}

func BenchmarkSpanProcessor(b *testing.B) {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			tracetest.NewNoopExporter(),
			sdktrace.WithMaxExportBatchSize(10),
		))
	tracer := tp.Tracer("bench")
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_, span := tracer.Start(ctx, "bench")
			span.End()
		}
	}
}

func BenchmarkSpanProcessorVerboseLogging(b *testing.B) {
	global.SetLogger(funcr.New(func(prefix, args string) {}, funcr.Options{Verbosity: 5}))
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			tracetest.NewNoopExporter(),
			sdktrace.WithMaxExportBatchSize(10),
		))
	tracer := tp.Tracer("bench")
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_, span := tracer.Start(ctx, "bench")
			span.End()
		}
	}
}
