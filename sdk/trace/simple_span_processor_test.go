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

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/export"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testExporter struct {
	spans []*export.SpanData
}

func (t *testExporter) ExportSpan(ctx context.Context, s *export.SpanData) {
	t.spans = append(t.spans, s)
}

var _ export.SpanSyncer = (*testExporter)(nil)

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
	tp := basicProvider(t)
	te := testExporter{}
	ssp := sdktrace.NewSimpleSpanProcessor(&te)
	if ssp == nil {
		t.Errorf("Error creating new instance of SimpleSpanProcessor with nil Exporter\n")
	}

	tp.RegisterSpanProcessor(ssp)
	tr := tp.GetTracer("SimpleSpanProcessor")
	tid, _ := core.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ := core.SpanIDFromHex("0102040810203040")
	sc := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	}
	_, span := tr.Start(context.Background(), "OnEnd", apitrace.ChildOf(sc))
	span.End()

	wantTraceID := tid
	gotTraceID := te.spans[0].SpanContext.TraceID
	if wantTraceID != gotTraceID {
		t.Errorf("SimplerSpanProcessor OnEnd() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}

func TestSimpleSpanProcessorShutdown(t *testing.T) {
	ssp := sdktrace.NewSimpleSpanProcessor(&testExporter{})
	if ssp == nil {
		t.Errorf("Error creating new instance of SimpleSpanProcessor\n")
	}

	ssp.Shutdown()
}
