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

	"go.opentelemetry.io/otel/sdk/export"
)

type testSpanProcesor struct {
	spansStarted  []*export.SpanData
	spansEnded    []*export.SpanData
	shutdownCount int
}

func (t *testSpanProcesor) OnStart(s *export.SpanData) {
	t.spansStarted = append(t.spansStarted, s)
}

func (t *testSpanProcesor) OnEnd(s *export.SpanData) {
	t.spansEnded = append(t.spansEnded, s)
}

func (t *testSpanProcesor) Shutdown() {
	t.shutdownCount++
}

func TestRegisterSpanProcessort(t *testing.T) {
	name := "Register span processor before span starts"
	tp := basicProvider(t)
	sp := NewTestSpanProcessor()
	tp.RegisterSpanProcessor(sp)

	tr := tp.GetTracer("SpanProcessor")
	_, span := tr.Start(context.Background(), "OnStart")
	span.End()
	wantCount := 1
	gotCount := len(sp.spansStarted)
	if gotCount != wantCount {
		t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
	}
	gotCount = len(sp.spansEnded)
	if gotCount != wantCount {
		t.Errorf("%s: ended count: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestUnregisterSpanProcessor(t *testing.T) {
	name := "Start span after unregistering span processor"
	tp := basicProvider(t)
	sp := NewTestSpanProcessor()
	tp.RegisterSpanProcessor(sp)

	tr := tp.GetTracer("SpanProcessor")
	_, span := tr.Start(context.Background(), "OnStart")
	span.End()
	tp.UnregisterSpanProcessor(sp)

	// start another span after unregistering span processor.
	_, span = tr.Start(context.Background(), "Start span after unregister")
	span.End()

	wantCount := 1
	gotCount := len(sp.spansStarted)
	if gotCount != wantCount {
		t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
	}

	gotCount = len(sp.spansEnded)
	if gotCount != wantCount {
		t.Errorf("%s: ended count: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestUnregisterSpanProcessorWhileSpanIsActive(t *testing.T) {
	name := "Unregister span processor while span is active"
	tp := basicProvider(t)
	sp := NewTestSpanProcessor()
	tp.RegisterSpanProcessor(sp)

	tr := tp.GetTracer("SpanProcessor")
	_, span := tr.Start(context.Background(), "OnStart")
	tp.UnregisterSpanProcessor(sp)

	span.End()

	wantCount := 1
	gotCount := len(sp.spansStarted)
	if gotCount != wantCount {
		t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
	}

	wantCount = 0
	gotCount = len(sp.spansEnded)
	if gotCount != wantCount {
		t.Errorf("%s: ended count: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestSpanProcessorShutdown(t *testing.T) {
	name := "Increment shutdown counter of a span processor"
	tp := basicProvider(t)
	sp := NewTestSpanProcessor()
	if sp == nil {
		t.Fatalf("Error creating new instance of TestSpanProcessor\n")
	}
	tp.RegisterSpanProcessor(sp)

	wantCount := 1
	sp.Shutdown()

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestMultipleUnregisterSpanProcessorCalls(t *testing.T) {
	name := "Increment shutdown counter after first UnregisterSpanProcessor call"
	tp := basicProvider(t)
	sp := NewTestSpanProcessor()
	if sp == nil {
		t.Fatalf("Error creating new instance of TestSpanProcessor\n")
	}

	wantCount := 1

	tp.RegisterSpanProcessor(sp)
	tp.UnregisterSpanProcessor(sp)

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}

	// Multiple UnregisterSpanProcessor should not trigger multiple Shutdown calls.
	tp.UnregisterSpanProcessor(sp)

	gotCount = sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func NewTestSpanProcessor() *testSpanProcesor {
	return &testSpanProcesor{}
}
