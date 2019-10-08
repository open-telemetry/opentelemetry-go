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

	apitrace "go.opentelemetry.io/api/trace"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

type testSpanProcesor struct {
	spansStarted  []*sdktrace.SpanData
	spansEnded    []*sdktrace.SpanData
	shutdownCount int
}

func init() {
	sdktrace.Register()
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
}

func (t *testSpanProcesor) OnStart(s *sdktrace.SpanData) {
	t.spansStarted = append(t.spansStarted, s)
}

func (t *testSpanProcesor) OnEnd(s *sdktrace.SpanData) {
	t.spansEnded = append(t.spansEnded, s)
}

func (t *testSpanProcesor) Shutdown() {
	t.shutdownCount++
}

func TestRegisterSpanProcessort(t *testing.T) {
	name := "Register span processor before span starts"
	sp := NewTestSpanProcessor()
	sdktrace.RegisterSpanProcessor(sp)
	defer sdktrace.UnregisterSpanProcessor(sp)
	_, span := apitrace.GlobalTracer().Start(context.Background(), "OnStart")
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
	sp := NewTestSpanProcessor()
	sdktrace.RegisterSpanProcessor(sp)
	_, span := apitrace.GlobalTracer().Start(context.Background(), "OnStart")
	span.End()
	sdktrace.UnregisterSpanProcessor(sp)

	// start another span after unregistering span processor.
	_, span = apitrace.GlobalTracer().Start(context.Background(), "Start span after unregister")
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
	sp := NewTestSpanProcessor()
	sdktrace.RegisterSpanProcessor(sp)
	_, span := apitrace.GlobalTracer().Start(context.Background(), "OnStart")
	sdktrace.UnregisterSpanProcessor(sp)

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
	sp := NewTestSpanProcessor()
	if sp == nil {
		t.Fatalf("Error creating new instance of TestSpanProcessor\n")
	}

	wantCount := sp.shutdownCount + 1
	sp.Shutdown()

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestMultipleUnregisterSpanProcessorCalls(t *testing.T) {
	name := "Increment shutdown counter after each UnregisterSpanProcessor call"
	sp := NewTestSpanProcessor()
	if sp == nil {
		t.Fatalf("Error creating new instance of TestSpanProcessor\n")
	}

	wantCount := sp.shutdownCount + 1

	sdktrace.RegisterSpanProcessor(sp)
	sdktrace.UnregisterSpanProcessor(sp)

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}

	// Multiple UnregisterSpanProcessor triggers multiple Shutdown calls.
	wantCount = wantCount + 1
	sdktrace.UnregisterSpanProcessor(sp)

	gotCount = sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func NewTestSpanProcessor() *testSpanProcesor {
	return &testSpanProcesor{}
}
