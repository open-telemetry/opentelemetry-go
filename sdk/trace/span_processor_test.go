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
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type testSpanProcessor struct {
	name          string
	spansStarted  []sdktrace.ReadWriteSpan
	spansEnded    []sdktrace.ReadOnlySpan
	shutdownCount int
}

func (t *testSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	if t == nil {
		return
	}
	psc := trace.SpanContextFromContext(parent)
	kv := []attribute.KeyValue{
		{
			Key:   "SpanProcessorName",
			Value: attribute.StringValue(t.name),
		},
		// Store parent trace ID and span ID as attributes to be read later in
		// tests so that we "do something" with the parent argument. Real
		// SpanProcessor implementations will likely use the parent argument in
		// a more meaningful way.
		{
			Key:   "ParentTraceID",
			Value: attribute.StringValue(psc.TraceID().String()),
		},
		{
			Key:   "ParentSpanID",
			Value: attribute.StringValue(psc.SpanID().String()),
		},
	}
	s.AddEvent("OnStart", trace.WithAttributes(kv...))
	t.spansStarted = append(t.spansStarted, s)
}

func (t *testSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	if t == nil {
		return
	}
	t.spansEnded = append(t.spansEnded, s)
}

func (t *testSpanProcessor) Shutdown(_ context.Context) error {
	if t == nil {
		return nil
	}
	t.shutdownCount++
	return nil
}

func (t *testSpanProcessor) ForceFlush(context.Context) error {
	if t == nil {
		return nil
	}
	return nil
}

func TestRegisterSpanProcessor(t *testing.T) {
	name := "Register span processor before span starts"
	tp := basicTracerProvider(t)
	spNames := []string{"sp1", "sp2", "sp3"}
	sps := NewNamedTestSpanProcessors(spNames)

	for _, sp := range sps {
		tp.RegisterSpanProcessor(sp)
	}

	tid, _ := trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ := trace.SpanIDFromHex("0102040810203040")
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: tid,
		SpanID:  sid,
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)

	tr := tp.Tracer("SpanProcessor")
	_, span := tr.Start(ctx, "OnStart")
	span.End()
	wantCount := 1

	for _, sp := range sps {
		gotCount := len(sp.spansStarted)
		if gotCount != wantCount {
			t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
		}
		gotCount = len(sp.spansEnded)
		if gotCount != wantCount {
			t.Errorf("%s: ended count: got %d, want %d\n", name, gotCount, wantCount)
		}

		c := 0
		tidOK := false
		sidOK := false
		for _, e := range sp.spansStarted[0].Events() {
			for _, kv := range e.Attributes {
				switch kv.Key {
				case "SpanProcessorName":
					gotValue := kv.Value.AsString()
					if gotValue != spNames[c] {
						t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, spNames[c])
					}
					c++
				case "ParentTraceID":
					gotValue := kv.Value.AsString()
					if gotValue != parent.TraceID().String() {
						t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, parent.TraceID())
					}
					tidOK = true
				case "ParentSpanID":
					gotValue := kv.Value.AsString()
					if gotValue != parent.SpanID().String() {
						t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, parent.SpanID())
					}
					sidOK = true
				default:
					continue
				}
			}
		}
		if c != len(spNames) {
			t.Errorf("%s: expected attributes(SpanProcessorName): got %d, want %d\n", name, c, len(spNames))
		}
		if !tidOK {
			t.Errorf("%s: expected attributes(ParentTraceID)\n", name)
		}
		if !sidOK {
			t.Errorf("%s: expected attributes(ParentSpanID)\n", name)
		}
	}
}

func TestUnregisterSpanProcessor(t *testing.T) {
	name := "Start span after unregistering span processor"
	tp := basicTracerProvider(t)
	spNames := []string{"sp1", "sp2", "sp3"}
	sps := NewNamedTestSpanProcessors(spNames)

	for _, sp := range sps {
		tp.RegisterSpanProcessor(sp)
	}

	tr := tp.Tracer("SpanProcessor")
	_, span := tr.Start(context.Background(), "OnStart")
	span.End()
	for _, sp := range sps {
		tp.UnregisterSpanProcessor(sp)
	}

	// start another span after unregistering span processor.
	_, span = tr.Start(context.Background(), "Start span after unregister")
	span.End()

	for _, sp := range sps {
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
}

func TestUnregisterSpanProcessorWhileSpanIsActive(t *testing.T) {
	name := "Unregister span processor while span is active"
	tp := basicTracerProvider(t)
	sp := NewTestSpanProcessor("sp")
	tp.RegisterSpanProcessor(sp)

	tr := tp.Tracer("SpanProcessor")
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
	tp := basicTracerProvider(t)
	sp := NewTestSpanProcessor("sp")
	tp.RegisterSpanProcessor(sp)

	wantCount := 1
	err := sp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the testSpanProcessor down\n")
	}

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestMultipleUnregisterSpanProcessorCalls(t *testing.T) {
	name := "Increment shutdown counter after first UnregisterSpanProcessor call"
	tp := basicTracerProvider(t)
	sp := NewTestSpanProcessor("sp")

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

func NewTestSpanProcessor(name string) *testSpanProcessor {
	return &testSpanProcessor{name: name}
}

func NewNamedTestSpanProcessors(names []string) []*testSpanProcessor {
	tsp := []*testSpanProcessor{}
	for _, n := range names {
		tsp = append(tsp, NewTestSpanProcessor(n))
	}
	return tsp
}
