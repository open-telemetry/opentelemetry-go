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

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testBatchExporter struct {
	mu         sync.Mutex
	spans      []*export.SpanData
	sizes      []int
	batchCount int
}

func (t *testBatchExporter) ExportSpans(ctx context.Context, sds []*export.SpanData) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.spans = append(t.spans, sds...)
	t.sizes = append(t.sizes, len(sds))
	t.batchCount++
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

var _ export.SpanBatcher = (*testBatchExporter)(nil)

func TestNewBatchSpanProcessorWithNilExporter(t *testing.T) {
	_, err := sdktrace.NewBatchSpanProcessor(nil)
	if err == nil {
		t.Errorf("Expected error while creating processor with nil exporter")
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
				sdktrace.WithBlocking(),
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
				sdktrace.WithBlocking(),
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
			tp := basicProvider(t)
			ssp := createAndRegisterBatchSP(t, option, &te)
			if ssp == nil {
				t.Fatalf("%s: Error creating new instance of BatchSpanProcessor\n", option.name)
			}
			tp.RegisterSpanProcessor(ssp)
			tr := tp.Tracer("BatchSpanProcessorWithOptions")

			generateSpan(t, option.parallel, tr, option)

			tp.UnregisterSpanProcessor(ssp)

			// TODO(https://github.com/open-telemetry/opentelemetry-go/issues/741)
			// Restore some sort of test here.
			_ = option.wantNumSpans
			_ = option.wantBatchCount
			_ = te.len()           // gotNumOfSpans
			_ = te.getBatchCount() // gotBatchCount
		})
	}
}

func createAndRegisterBatchSP(t *testing.T, option testOption, te *testBatchExporter) *sdktrace.BatchSpanProcessor {
	ssp, err := sdktrace.NewBatchSpanProcessor(te, option.o...)
	if ssp == nil {
		t.Errorf("%s: Error creating new instance of BatchSpanProcessor, error: %v\n", option.name, err)
	}
	return ssp
}

func generateSpan(t *testing.T, parallel bool, tr apitrace.Tracer, option testOption) {
	sc := getSpanContext()

	wg := &sync.WaitGroup{}
	for i := 0; i < option.genNumSpans; i++ {
		binary.BigEndian.PutUint64(sc.TraceID[0:8], uint64(i+1))
		wg.Add(1)
		f := func(sc apitrace.SpanContext) {
			ctx := apitrace.ContextWithRemoteSpanContext(context.Background(), sc)
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

func getSpanContext() apitrace.SpanContext {
	tid, _ := apitrace.IDFromHex("01020304050607080102040810203040")
	sid, _ := apitrace.SpanIDFromHex("0102040810203040")
	return apitrace.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	}
}

func TestBatchSpanProcessorShutdown(t *testing.T) {
	bsp, err := sdktrace.NewBatchSpanProcessor(&testBatchExporter{})
	if err != nil {
		t.Errorf("Unexpected error while creating processor\n")
	}

	if bsp == nil {
		t.Fatalf("Error creating new instance of BatchSpanProcessor\n")
	}

	bsp.Shutdown()

	// Multiple call to Shutdown() should not panic.
	bsp.Shutdown()
}
