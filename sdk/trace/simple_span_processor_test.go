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
	"testing"

	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

type testExporter struct {
	spans []*sdktrace.SpanData
}

func (t *testExporter) ExportSpan(s *sdktrace.SpanData) {
	t.spans = append(t.spans, s)
}

var _ sdktrace.Exporter = (*testExporter)(nil)

func init() {
	sdktrace.Register()
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
}

func TestNewSimpleSpanProcessor(t *testing.T) {
	ssp := sdktrace.NewSimpleSpanProcessor(&testExporter{})
	if ssp == nil {
		t.Errorf("Error creating new instance of SimpleSpanProcessor\n")
	}
}

func TestNewSimpleSpanProcessorWithNilExporter(t *testing.T) {
	ssp := sdktrace.NewSimpleSpanProcessor(nil)
	if ssp == nil {
		t.Errorf("Error creating new instance of SimpleSpanProcessor with nil Exporter\n")
	}
}

func TestSimpleSpanProcessorOnEnd(t *testing.T) {
	te := testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(&te)
	if ssp == nil {
		t.Errorf("Error creating new instance of SimpleSpanProcessor with nil Exporter\n")
	}
	sdktrace.RegisterSpanProcessor(ssp)
	tid := core.TraceID{High: 0x0102030405060708, Low: 0x0102040810203040}
	sid := uint64(0x0102040810203040)
	sc := core.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: 0x1,
	}
	_, span := apitrace.GlobalTracer().Start(context.Background(), "OnEnd", apitrace.ChildOf(sc))
	span.Finish()

	wantTraceID := tid
	gotTraceID := te.spans[0].SpanContext.TraceID
	if wantTraceID != gotTraceID {
		t.Errorf("SimplerSpanProcessor OnEnd() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}
