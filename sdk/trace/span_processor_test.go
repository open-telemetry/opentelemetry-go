// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type testSpanProcessor struct {
	name          string
	spansStarted  []ReadWriteSpan
	spansEnded    []ReadOnlySpan
	shutdownCount int
}

func (t *testSpanProcessor) OnStart(parent context.Context, s ReadWriteSpan) {
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

func (t *testSpanProcessor) OnEnd(s ReadOnlySpan) {
	if t == nil {
		return
	}
	t.spansEnded = append(t.spansEnded, s)
}

func (t *testSpanProcessor) Shutdown(context.Context) error {
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
	ctx := trace.ContextWithRemoteSpanContext(t.Context(), parent)

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
	_, span := tr.Start(t.Context(), "OnStart")
	span.End()
	for _, sp := range sps {
		tp.UnregisterSpanProcessor(sp)
	}

	// start another span after unregistering span processor.
	_, span = tr.Start(t.Context(), "Start span after unregister")
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
	_, span := tr.Start(t.Context(), "OnStart")
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
	err := sp.Shutdown(t.Context())
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

// onEndingProcessor is a SpanProcessor that also implements OnEnding.
type onEndingProcessor struct {
	testSpanProcessor
	onEndingCalled bool
	endTimeSet     bool
	setAttr        attribute.KeyValue
	order          *[]string
}

func (p *onEndingProcessor) OnEnding(s ReadWriteSpan, _ trace.Span) {
	p.onEndingCalled = true
	p.endTimeSet = !s.EndTime().IsZero()
	if p.order != nil {
		*p.order = append(*p.order, p.name)
	}
	s.SetAttributes(p.setAttr)
}

func TestOnEndingCalledBeforeOnEnd(t *testing.T) {
	tp := basicTracerProvider(t)
	p := &onEndingProcessor{
		testSpanProcessor: testSpanProcessor{name: "p1"},
		setAttr:           attribute.Bool("from-onending", true),
	}
	tp.RegisterSpanProcessor(p)

	_, span := tp.Tracer("t").Start(t.Context(), "s")
	span.End()

	assert.True(t, p.onEndingCalled, "OnEnding should be called")
	assert.True(t, p.endTimeSet, "EndTime should be set when OnEnding is called")

	require.Len(t, p.spansEnded, 1)
	var found bool
	for _, kv := range p.spansEnded[0].Attributes() {
		if kv.Key == "from-onending" && kv.Value.AsBool() {
			found = true
		}
	}
	assert.True(t, found, "attribute set in OnEnding should appear in OnEnd snapshot")
}

func TestOnEndingOrder(t *testing.T) {
	tp := basicTracerProvider(t)
	order := []string{}
	p1 := &onEndingProcessor{testSpanProcessor: testSpanProcessor{name: "p1"}, order: &order}
	p2 := &onEndingProcessor{testSpanProcessor: testSpanProcessor{name: "p2"}, order: &order}
	tp.RegisterSpanProcessor(p1)
	tp.RegisterSpanProcessor(p2)

	_, span := tp.Tracer("t").Start(t.Context(), "s")
	span.End()

	assert.Equal(t, []string{"p1", "p2"}, order, "OnEnding should run in registration order")
}

func TestOnEndingSkippedForProcessorsWithoutIt(t *testing.T) {
	tp := basicTracerProvider(t)
	regular := NewTestSpanProcessor("regular")
	withHook := &onEndingProcessor{
		testSpanProcessor: testSpanProcessor{name: "hook"},
		setAttr:           attribute.String("marker", "set"),
	}
	tp.RegisterSpanProcessor(regular)
	tp.RegisterSpanProcessor(withHook)

	_, span := tp.Tracer("t").Start(t.Context(), "s")
	span.End()

	assert.True(t, withHook.onEndingCalled)
	require.Len(t, regular.spansEnded, 1)
	var found bool
	for _, kv := range regular.spansEnded[0].Attributes() {
		if kv.Key == "marker" {
			found = true
		}
	}
	assert.True(t, found, "regular OnEnd should see attribute set in OnEnding of another processor")
}

// blockingOnEndingProcessor lets the test synchronize with OnEnding so that a
// concurrent mutation attempt through the original span reference can be
// verified to have no effect.
type blockingOnEndingProcessor struct {
	testSpanProcessor
	started chan struct{}
	release chan struct{}
}

func (p *blockingOnEndingProcessor) OnEnding(ending ReadWriteSpan, _ trace.Span) {
	close(p.started)
	<-p.release
	ending.SetAttributes(attribute.String("from-onending", "yes"))
}

// TestOnEndingConcurrencyGuarantee verifies the spec requirement that no other
// goroutine may modify the span while OnEnding is executing. A concurrent
// SetAttributes call on the original span reference must be dropped, while an
// attribute set via the ending parameter must appear in the OnEnd snapshot.
func TestOnEndingConcurrencyGuarantee(t *testing.T) {
	tp := basicTracerProvider(t)
	p := &blockingOnEndingProcessor{
		testSpanProcessor: testSpanProcessor{name: "p"},
		started:           make(chan struct{}),
		release:           make(chan struct{}),
	}
	tp.RegisterSpanProcessor(p)

	_, span := tp.Tracer("t").Start(t.Context(), "s")

	endDone := make(chan struct{})
	go func() {
		span.End()
		close(endDone)
	}()

	// Wait until OnEnding has started, then try to mutate through the original
	// span reference. The SDK must ignore this because the span is in the
	// ending-only window.
	<-p.started
	span.SetAttributes(attribute.String("race", "should-not-appear"))
	close(p.release)
	<-endDone

	require.Len(t, p.spansEnded, 1)
	snap := p.spansEnded[0].Attributes()
	for _, kv := range snap {
		assert.NotEqual(t, attribute.Key("race"), kv.Key, "concurrent mutation must be dropped")
	}
	var found bool
	for _, kv := range snap {
		if kv.Key == "from-onending" && kv.Value.AsString() == "yes" {
			found = true
		}
	}
	assert.True(t, found, "attribute set via ending parameter must appear in OnEnd snapshot")
}

// correlatingProcessor records the span from OnStart and checks that the
// original parameter in OnEnding is the same Go value.
type correlatingProcessor struct {
	testSpanProcessor
	started trace.Span
}

func (p *correlatingProcessor) OnStart(_ context.Context, s ReadWriteSpan) {
	p.started = s
}

func (p *correlatingProcessor) OnEnding(_ ReadWriteSpan, original trace.Span) {
	if original != p.started {
		panic("original span in OnEnding is not the same value as the one from OnStart")
	}
}

func TestOnEndingOriginalMatchesOnStart(t *testing.T) {
	tp := basicTracerProvider(t)
	p := &correlatingProcessor{}
	tp.RegisterSpanProcessor(p)

	_, span := tp.Tracer("t").Start(t.Context(), "s")
	// Confirm OnStart saw the same span the caller holds.
	assert.Equal(t, span, p.started)
	// End calls OnEnding; the correlatingProcessor panics if original != started.
	assert.NotPanics(t, func() { span.End() })
}

// reentrantRecordErrorOnEndingProcessor calls RecordError from OnEnding using
// an error whose Error() method calls back into the same span. This exercises
// the same reentrancy class fixed for recordingSpan.RecordError by #8815.
type reentrantRecordErrorOnEndingProcessor struct {
	testSpanProcessor
}

func (*reentrantRecordErrorOnEndingProcessor) OnEnding(s ReadWriteSpan, _ trace.Span) {
	s.RecordError(reentrantRecordError{span: s})
}

func TestOnEndingRecordErrorAllowsReentrantErrorFormatting(t *testing.T) {
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(&reentrantRecordErrorOnEndingProcessor{})

	_, span := tp.Tracer(t.Name()).Start(t.Context(), "span")
	done := make(chan struct{})
	go func() {
		span.End()
		close(done)
	}()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("Span.End deadlocked while OnEnding formatted a reentrant error")
	}
}
