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
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testBatchExporter struct {
	mu         sync.Mutex
	spans      []*export.SpanData
	sizes      []int
	batchCount int
}

func (t *testBatchExporter) ExportSpans(ctx context.Context, sds []*export.SpanData) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.spans = append(t.spans, sds...)
	t.sizes = append(t.sizes, len(sds))
	t.batchCount++
	return nil
}

func (t *testBatchExporter) Shutdown(context.Context) error { return nil }

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

var _ export.SpanExporter = (*testBatchExporter)(nil)

func TestNewBatchSpanProcessorWithNilExporter(t *testing.T) {
	bsp := sdktrace.NewBatchSpanProcessor(nil)
	// These should not panic.
	bsp.OnStart(context.Background(), &export.SpanData{})
	bsp.OnEnd(&export.SpanData{})
	bsp.ForceFlush()
	err := bsp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the BatchSpanProcessor down\n")
	}
}

type testOption struct {
	name           string
	o              []sdktrace.BatchSpanProcessorOption
	wantNumSpans   int
	wantBatchCount int
	genNumSpans    int
	parallel       bool
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
			if option.wantNumSpans != gotNumOfSpans {
				t.Errorf("number of exported span: got %+v, want %+v\n",
					gotNumOfSpans, option.wantNumSpans)
			}

			gotBatchCount := te.getBatchCount()
			if gotBatchCount < option.wantBatchCount {
				t.Errorf("number batches: got %+v, want >= %+v\n",
					gotBatchCount, option.wantBatchCount)
				t.Errorf("Batches %v\n", te.sizes)
			}
		})
	}
}

func createAndRegisterBatchSP(option testOption, te *testBatchExporter) *sdktrace.BatchSpanProcessor {
	// Always use blocking queue to avoid flaky tests.
	options := append(option.o, sdktrace.WithBlocking())
	return sdktrace.NewBatchSpanProcessor(te, options...)
}

func generateSpan(t *testing.T, parallel bool, tr trace.Tracer, option testOption) {
	sc := getSpanContext()

	wg := &sync.WaitGroup{}
	for i := 0; i < option.genNumSpans; i++ {
		binary.BigEndian.PutUint64(sc.TraceID[0:8], uint64(i+1))
		wg.Add(1)
		f := func(sc trace.SpanContext) {
			ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
			_, span := tr.Start(ctx, option.name)
			span.End()
			wg.Done()
		}
		if parallel {
			go f(sc)
		} else {
			f(sc)
		}
	}
	wg.Wait()
}

func getSpanContext() trace.SpanContext {
	tid, _ := trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ := trace.SpanIDFromHex("0102040810203040")
	return trace.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	}
}

func TestBatchSpanProcessorShutdown(t *testing.T) {
	bsp := sdktrace.NewBatchSpanProcessor(&testBatchExporter{})

	err := bsp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the BatchSpanProcessor down\n")
	}

	// Multiple call to Shutdown() should not panic.
	err = bsp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the BatchSpanProcessor down\n")
	}
}
