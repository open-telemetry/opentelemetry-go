// Copyright 2019, OpenTelemetry Authors
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
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

type testBatchExporter struct {
	mu sync.Mutex
	spans []*sdktrace.SpanData
	batchCount int
}

func (t *testBatchExporter) ExportSpans(sds []*sdktrace.SpanData) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = append(t.spans, sds...)
	t.batchCount++
}

func (t *testBatchExporter) len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.spans)
}

func (t *testBatchExporter) get(idx int) *sdktrace.SpanData {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.spans[idx]
}

var _ sdktrace.BatchExporter = (*testBatchExporter)(nil)

var defaultOpts = sdktrace.BatchSpanProcessorOption{}

func init() {
	sdktrace.Register()
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
}

func TestNewBatchSpanProcessor(t *testing.T) {
	ssp := sdktrace.NewBatchSpanProcessor(&testBatchExporter{}, defaultOpts)
	if ssp == nil {
		t.Errorf("Error creating new instance of BatchSpanProcessor\n")
	}
}

func TestNewBatchSpanProcessorWithNilExporter(t *testing.T) {
	ssp := sdktrace.NewBatchSpanProcessor(nil, defaultOpts)
	if ssp == nil {
		t.Errorf("Error creating new instance of BatchSpanProcessor with nil Exporter\n")
	}
}

func TestNewBatchSpanProcessorWithOptions(t *testing.T) {
	options := []struct {
		name string
		o sdktrace.BatchSpanProcessorOption
	} {
		{
			name: "default BatchSpanProcessorOption",
			o : sdktrace.BatchSpanProcessorOption{
			},
		},
		{
			name: "non-default ScheduledDelayMillis",
			o : sdktrace.BatchSpanProcessorOption{
				ScheduledDelayMillis: time.Duration(100 * time.Millisecond),
			},
		},
		{
			name: "non-default MaxQueueSize",
			o : sdktrace.BatchSpanProcessorOption{
				MaxQueueSize: 200,
				ScheduledDelayMillis: time.Duration(100 * time.Millisecond),
			},
		},
		{
			name: "non-default MaxExportBatchSize",
			o : sdktrace.BatchSpanProcessorOption{
				MaxQueueSize: 205,
				MaxExportBatchSize: 20,
				ScheduledDelayMillis: time.Duration(100 * time.Millisecond),
			},
		},
		{
			name: "all non-default batch option",
			o : sdktrace.BatchSpanProcessorOption{
				MaxExportBatchSize: 100,
				MaxQueueSize: 200,
				ScheduledDelayMillis: time.Duration(100 * time.Millisecond),
			},
		},
	}
	for _, option := range options {
		te := testBatchExporter{}
		expectO := getExpectedOptions(option.o)
		numOfSpan := expectO.MaxQueueSize + 5
		ssp := sdktrace.NewBatchSpanProcessor(&te, option.o)
		if ssp == nil {
			t.Errorf("Error creating new instance of BatchSpanProcessor with nil Exporter\n")
		}
		sdktrace.RegisterSpanProcessor(ssp)
		sc := getSpanContext()
		for i := 0; i < numOfSpan ; i++ {
			sc.TraceID.High = uint64(i+1)
			_, span := apitrace.GlobalTracer().Start(context.Background(), option.name, apitrace.ChildOf(sc))
			span.Finish()
		}

		time.Sleep(expectO.ScheduledDelayMillis + time.Duration(100 + time.Millisecond))

		gotNumOfSpans := te.len()
		wantNumOfSpans := expectO.MaxQueueSize
		if wantNumOfSpans != gotNumOfSpans {
			t.Errorf("BatchSpanProcessor number of exported span: got %+v, want %+v\n", gotNumOfSpans, wantNumOfSpans)
		}

		gotBatchCount := te.batchCount
		wantBatchCount := expectO.MaxQueueSize / expectO.MaxExportBatchSize
		if expectO.MaxQueueSize % expectO.MaxExportBatchSize > 0 {
			wantBatchCount += 1
		}
		if wantBatchCount != gotBatchCount {
			t.Errorf("BatchSpanProcessor number batches: got %+v, want %+v\n", gotBatchCount, wantBatchCount)
		}

		// Check first Span is reported. Most recent one is dropped.
		wantTraceID := sc.TraceID
		wantTraceID.High = 1;
		gotTraceID := te.get(0).SpanContext.TraceID
		if wantTraceID != gotTraceID {
			t.Errorf("BatchSpanProcessor first exported span: got %+v, want %+v\n", gotTraceID, wantTraceID)
		}
		sdktrace.UnregisterSpanProcessor(ssp)
	}
}

func getExpectedOptions(o sdktrace.BatchSpanProcessorOption) sdktrace.BatchSpanProcessorOption {
	expectedO := o
	if expectedO.ScheduledDelayMillis <= 0 {
		expectedO.ScheduledDelayMillis = time.Duration(5000 * time.Millisecond)
	}
	if expectedO.MaxQueueSize <= 0 {
		expectedO.MaxQueueSize = 2048
	}
	if expectedO.MaxExportBatchSize <= 0 {
		expectedO.MaxExportBatchSize = 512
	}
	return expectedO
}

func getSpanContext() core.SpanContext {
	tid := core.TraceID{High: 0x0102030405060708, Low: 0x0102040810203040}
	sid := uint64(0x0102040810203040)
	return core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	}
}