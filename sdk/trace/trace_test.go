// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const envVarResourceAttributes = "OTEL_RESOURCE_ATTRIBUTES"

type storingHandler struct {
	errs []error
}

func (s *storingHandler) Handle(err error) {
	s.errs = append(s.errs, err)
}

func (s *storingHandler) Reset() {
	s.errs = nil
}

var (
	tid trace.TraceID
	sid trace.SpanID
	sc  trace.SpanContext
	ts  trace.TraceState

	handler = &storingHandler{}
)

func init() {
	tid, _ = trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ = trace.SpanIDFromHex("0102040810203040")
	sc = trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	})
	ts, _ = trace.ParseTraceState("k=v")

	otel.SetErrorHandler(handler)
}

func TestTracerFollowsExpectedAPIBehaviour(t *testing.T) {
	harness := newHarness(t)

	harness.testTracerProvider(func() trace.TracerProvider {
		return NewTracerProvider(WithSampler(TraceIDRatioBased(0)))
	})

	tp := NewTracerProvider(WithSampler(TraceIDRatioBased(0)))
	harness.testTracer(func() trace.Tracer {
		return tp.Tracer("")
	})
}

type testExporter struct {
	mu    sync.RWMutex
	idx   map[string]int
	spans []*snapshot
}

func NewTestExporter() *testExporter {
	return &testExporter{idx: make(map[string]int)}
}

func (te *testExporter) ExportSpans(_ context.Context, spans []ReadOnlySpan) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	i := len(te.spans)
	for _, s := range spans {
		te.idx[s.Name()] = i
		te.spans = append(te.spans, s.(*snapshot))
		i++
	}
	return nil
}

func (te *testExporter) Spans() []*snapshot {
	te.mu.RLock()
	defer te.mu.RUnlock()

	cp := make([]*snapshot, len(te.spans))
	copy(cp, te.spans)
	return cp
}

func (te *testExporter) GetSpan(name string) (*snapshot, bool) {
	te.mu.RLock()
	defer te.mu.RUnlock()
	i, ok := te.idx[name]
	if !ok {
		return nil, false
	}
	return te.spans[i], true
}

func (te *testExporter) Len() int {
	te.mu.RLock()
	defer te.mu.RUnlock()
	return len(te.spans)
}

func (te *testExporter) Shutdown(context.Context) error {
	te.Reset()
	return nil
}

func (te *testExporter) Reset() {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.idx = make(map[string]int)
	te.spans = te.spans[:0]
}

type testSampler struct {
	callCount int
	prefix    string
	t         *testing.T
}

func (ts *testSampler) ShouldSample(p SamplingParameters) SamplingResult {
	ts.callCount++
	ts.t.Logf("called sampler for name %q", p.Name)
	decision := Drop
	if strings.HasPrefix(p.Name, ts.prefix) {
		decision = RecordAndSample
	}
	return SamplingResult{
		Decision:   decision,
		Attributes: []attribute.KeyValue{attribute.Int("callCount", ts.callCount)},
	}
}

func (ts testSampler) Description() string {
	return "testSampler"
}

func TestSetName(t *testing.T) {
	tp := NewTracerProvider()

	type testCase struct {
		name    string
		newName string
	}
	for idx, tt := range []testCase{
		{ // 0
			name:    "foobar",
			newName: "foobaz",
		},
		{ // 1
			name:    "foobar",
			newName: "barbaz",
		},
		{ // 2
			name:    "barbar",
			newName: "barbaz",
		},
		{ // 3
			name:    "barbar",
			newName: "foobar",
		},
	} {
		sp := startNamedSpan(tp, "SetName", tt.name)
		if sdkspan, ok := sp.(*recordingSpan); ok {
			if sdkspan.Name() != tt.name {
				t.Errorf("%d: invalid name at span creation, expected %v, got %v", idx, tt.name, sdkspan.Name())
			}
		} else {
			t.Errorf("%d: unable to coerce span to SDK span, is type %T", idx, sp)
		}
		sp.SetName(tt.newName)
		if sdkspan, ok := sp.(*recordingSpan); ok {
			if sdkspan.Name() != tt.newName {
				t.Errorf("%d: span name not changed, expected %v, got %v", idx, tt.newName, sdkspan.Name())
			}
		} else {
			t.Errorf("%d: unable to coerce span to SDK span, is type %T", idx, sp)
		}
		sp.End()
	}
}

func TestSpanIsRecording(t *testing.T) {
	t.Run("while Span active", func(t *testing.T) {
		for name, tc := range map[string]struct {
			sampler Sampler
			want    bool
		}{
			"Always sample, recording on": {sampler: AlwaysSample(), want: true},
			"Never sample recording off":  {sampler: NeverSample(), want: false},
		} {
			tp := NewTracerProvider(WithSampler(tc.sampler))
			_, span := tp.Tracer(name).Start(context.Background(), "StartSpan")
			got := span.IsRecording()
			span.End()
			assert.Equal(t, tc.want, got, name)
		}
	})

	t.Run("after Span end", func(t *testing.T) {
		for name, tc := range map[string]Sampler{
			"Always Sample": AlwaysSample(),
			"Never Sample":  NeverSample(),
		} {
			tp := NewTracerProvider(WithSampler(tc))
			_, span := tp.Tracer(name).Start(context.Background(), "StartSpan")
			span.End()
			got := span.IsRecording()
			assert.False(t, got, name)
		}
	})
}

func TestSampling(t *testing.T) {
	idg := defaultIDGenerator()
	const total = 10000
	for name, tc := range map[string]struct {
		sampler       Sampler
		expect        float64
		parent        bool
		sampledParent bool
	}{
		// Span w/o a parent
		"NeverSample":           {sampler: NeverSample(), expect: 0},
		"AlwaysSample":          {sampler: AlwaysSample(), expect: 1.0},
		"TraceIdRatioBased_-1":  {sampler: TraceIDRatioBased(-1.0), expect: 0},
		"TraceIdRatioBased_.25": {sampler: TraceIDRatioBased(0.25), expect: .25},
		"TraceIdRatioBased_.50": {sampler: TraceIDRatioBased(0.50), expect: .5},
		"TraceIdRatioBased_.75": {sampler: TraceIDRatioBased(0.75), expect: .75},
		"TraceIdRatioBased_2.0": {sampler: TraceIDRatioBased(2.0), expect: 1},

		// Spans w/o a parent and using ParentBased(DelegateSampler()) Sampler, receive DelegateSampler's sampling decision
		"ParentNeverSample":           {sampler: ParentBased(NeverSample()), expect: 0},
		"ParentAlwaysSample":          {sampler: ParentBased(AlwaysSample()), expect: 1},
		"ParentTraceIdRatioBased_.50": {sampler: ParentBased(TraceIDRatioBased(0.50)), expect: .5},

		// An unadorned TraceIDRatioBased sampler ignores parent spans
		"UnsampledParentSpanWithTraceIdRatioBased_.25": {sampler: TraceIDRatioBased(0.25), expect: .25, parent: true},
		"SampledParentSpanWithTraceIdRatioBased_.25":   {sampler: TraceIDRatioBased(0.25), expect: .25, parent: true, sampledParent: true},
		"UnsampledParentSpanWithTraceIdRatioBased_.50": {sampler: TraceIDRatioBased(0.50), expect: .5, parent: true},
		"SampledParentSpanWithTraceIdRatioBased_.50":   {sampler: TraceIDRatioBased(0.50), expect: .5, parent: true, sampledParent: true},
		"UnsampledParentSpanWithTraceIdRatioBased_.75": {sampler: TraceIDRatioBased(0.75), expect: .75, parent: true},
		"SampledParentSpanWithTraceIdRatioBased_.75":   {sampler: TraceIDRatioBased(0.75), expect: .75, parent: true, sampledParent: true},

		// Spans with a sampled parent but using NeverSample Sampler, are not sampled
		"SampledParentSpanWithNeverSample": {sampler: NeverSample(), expect: 0, parent: true, sampledParent: true},

		// Spans with a sampled parent and using ParentBased(DelegateSampler()) Sampler, inherit the parent span's sampling status
		"SampledParentSpanWithParentNeverSample":             {sampler: ParentBased(NeverSample()), expect: 1, parent: true, sampledParent: true},
		"UnsampledParentSpanWithParentNeverSampler":          {sampler: ParentBased(NeverSample()), expect: 0, parent: true, sampledParent: false},
		"SampledParentSpanWithParentAlwaysSampler":           {sampler: ParentBased(AlwaysSample()), expect: 1, parent: true, sampledParent: true},
		"UnsampledParentSpanWithParentAlwaysSampler":         {sampler: ParentBased(AlwaysSample()), expect: 0, parent: true, sampledParent: false},
		"SampledParentSpanWithParentTraceIdRatioBased_.50":   {sampler: ParentBased(TraceIDRatioBased(0.50)), expect: 1, parent: true, sampledParent: true},
		"UnsampledParentSpanWithParentTraceIdRatioBased_.50": {sampler: ParentBased(TraceIDRatioBased(0.50)), expect: 0, parent: true, sampledParent: false},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := NewTracerProvider(WithSampler(tc.sampler))
			tr := p.Tracer("test")
			var sampled int
			for i := 0; i < total; i++ {
				ctx := context.Background()
				if tc.parent {
					tid, sid := idg.NewIDs(ctx)
					psc := trace.NewSpanContext(trace.SpanContextConfig{
						TraceID: tid,
						SpanID:  sid,
					})
					if tc.sampledParent {
						psc = psc.WithTraceFlags(trace.FlagsSampled)
					}
					ctx = trace.ContextWithRemoteSpanContext(ctx, psc)
				}
				_, span := tr.Start(ctx, "test")
				if span.SpanContext().IsSampled() {
					sampled++
				}
			}
			tolerance := 0.0
			got := float64(sampled) / float64(total)

			if tc.expect > 0 && tc.expect < 1 {
				// See https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval
				const z = 4.75342 // This should succeed 99.9999% of the time
				tolerance = z * math.Sqrt(got*(1-got)/total)
			}

			diff := math.Abs(got - tc.expect)
			if diff > tolerance {
				t.Errorf("got %f (diff: %f), expected %f (w/tolerance: %f)", got, diff, tc.expect, tolerance)
			}
		})
	}
}

func TestStartSpanWithParent(t *testing.T) {
	tp := NewTracerProvider()
	tr := tp.Tracer("SpanWithParent")
	ctx := context.Background()

	_, s1 := tr.Start(trace.ContextWithRemoteSpanContext(ctx, sc), "span1-unsampled-parent1")
	if err := checkChild(t, sc, s1); err != nil {
		t.Error(err)
	}

	_, s2 := tr.Start(trace.ContextWithRemoteSpanContext(ctx, sc), "span2-unsampled-parent1")
	if err := checkChild(t, sc, s2); err != nil {
		t.Error(err)
	}

	sc2 := sc.WithTraceState(ts)
	_, s3 := tr.Start(trace.ContextWithRemoteSpanContext(ctx, sc2), "span3-sampled-parent2")
	if err := checkChild(t, sc2, s3); err != nil {
		t.Error(err)
	}

	ctx2, s4 := tr.Start(trace.ContextWithRemoteSpanContext(ctx, sc2), "span4-sampled-parent2")
	if err := checkChild(t, sc2, s4); err != nil {
		t.Error(err)
	}

	s4Sc := s4.SpanContext()
	_, s5 := tr.Start(ctx2, "span5-implicit-childof-span4")
	if err := checkChild(t, s4Sc, s5); err != nil {
		t.Error(err)
	}
}

// Test we get a successful span as a new root if a nil context is sent in, as opposed to a panic.
// See https://github.com/open-telemetry/opentelemetry-go/issues/3109
func TestStartSpanWithNilContext(t *testing.T) {
	tp := NewTracerProvider()
	tr := tp.Tracer("NoPanic")

	// nolint:staticcheck // no nil context, but that's the point of the test.
	assert.NotPanics(t, func() { tr.Start(nil, "should-not-panic") })
}

func TestStartSpanNewRootNotSampled(t *testing.T) {
	alwaysSampleTp := NewTracerProvider()
	sampledTr := alwaysSampleTp.Tracer("AlwaysSampled")
	neverSampleTp := NewTracerProvider(WithSampler(ParentBased(NeverSample())))
	neverSampledTr := neverSampleTp.Tracer("ParentBasedNeverSample")
	ctx := context.Background()

	ctx, s1 := sampledTr.Start(trace.ContextWithRemoteSpanContext(ctx, sc), "span1-sampled")
	if err := checkChild(t, sc, s1); err != nil {
		t.Error(err)
	}

	_, s2 := neverSampledTr.Start(ctx, "span2-no-newroot")
	if !s2.SpanContext().IsSampled() {
		t.Error(
			fmt.Errorf(
				"got child span is not sampled, want child span with sampler: ParentBased(NeverSample()) to be sampled",
			),
		)
	}

	// Adding WithNewRoot causes child spans to not sample based on parent context
	_, s3 := neverSampledTr.Start(ctx, "span3-newroot", trace.WithNewRoot())
	if s3.SpanContext().IsSampled() {
		t.Error(
			fmt.Errorf(
				"got child span is sampled, want child span WithNewRoot() and with sampler: ParentBased(NeverSample()) to not be sampled",
			),
		)
	}
}

func TestSetSpanAttributesOnStart(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	span := startSpan(tp,
		"StartSpanAttribute",
		trace.WithAttributes(attribute.String("key1", "value1")),
		trace.WithAttributes(attribute.String("key2", "value2")),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent: sc.WithRemote(true),
		name:   "span0",
		attributes: []attribute.KeyValue{
			attribute.String("key1", "value1"),
			attribute.String("key2", "value2"),
		},
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "StartSpanAttribute"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOnStart: -got +want %s", diff)
	}
}

func TestSamplerAttributesLocalChildSpan(t *testing.T) {
	sampler := &testSampler{prefix: "span", t: t}
	te := NewTestExporter()
	tp := NewTracerProvider(WithSampler(sampler), WithSyncer(te), WithResource(resource.Empty()))

	ctx := context.Background()
	ctx, span := startLocalSpan(ctx, tp, "SpanOne", "span0")
	_, spanTwo := startLocalSpan(ctx, tp, "SpanTwo", "span1")

	spanTwo.End()
	span.End()

	got := te.Spans()
	require.Len(t, got, 2)
	// FILO order above means spanTwo <-> gotSpan0 and span <-> gotSpan1.
	gotSpan0, gotSpan1 := got[0], got[1]
	// Ensure sampler is called for local child spans by verifying the
	// attributes set by the sampler are set on the child span.
	assert.Equal(t, []attribute.KeyValue{attribute.Int("callCount", 2)}, gotSpan0.Attributes())
	assert.Equal(t, []attribute.KeyValue{attribute.Int("callCount", 1)}, gotSpan1.Attributes())
}

func TestSpanSetAttributes(t *testing.T) {
	attrs := [...]attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key3", "value3"),
		attribute.String("key4", "value4"),
		attribute.String("key1", "value5"),
		attribute.String("key2", "value6"),
		attribute.String("key3", "value7"),
	}
	invalid := attribute.KeyValue{}

	tests := []struct {
		name        string
		input       [][]attribute.KeyValue
		wantAttrs   []attribute.KeyValue
		wantDropped int
	}{
		{
			name:      "array",
			input:     [][]attribute.KeyValue{attrs[:3]},
			wantAttrs: attrs[:3],
		},
		{
			name:      "single_value:array",
			input:     [][]attribute.KeyValue{attrs[:1], attrs[1:3]},
			wantAttrs: attrs[:3],
		},
		{
			name:      "array:single_value",
			input:     [][]attribute.KeyValue{attrs[:2], attrs[2:3]},
			wantAttrs: attrs[:3],
		},
		{
			name:      "single_values",
			input:     [][]attribute.KeyValue{attrs[:1], attrs[1:2], attrs[2:3]},
			wantAttrs: attrs[:3],
		},

		// The tracing specification states:
		//
		//   For each unique attribute key, addition of which would result in
		//   exceeding the limit, SDK MUST discard that key/value pair
		//
		// Therefore, adding attributes after the capacity is reached should
		// result in those attributes being dropped.

		{
			name:        "drop_last_added",
			input:       [][]attribute.KeyValue{attrs[:3], attrs[3:4], attrs[3:4]},
			wantAttrs:   attrs[:3],
			wantDropped: 2,
		},

		// The tracing specification states:
		//
		//   Setting an attribute with the same key as an existing attribute
		//   SHOULD overwrite the existing attribute's value.
		//
		// Therefore, attributes are updated regardless of capacity state.

		{
			name:      "single_value_update",
			input:     [][]attribute.KeyValue{attrs[:1], attrs[:3]},
			wantAttrs: attrs[:3],
		},
		{
			name:      "all_update",
			input:     [][]attribute.KeyValue{attrs[:3], attrs[4:7]},
			wantAttrs: attrs[4:7],
		},
		{
			name:      "all_update/multi",
			input:     [][]attribute.KeyValue{attrs[:3], attrs[4:7], attrs[:3]},
			wantAttrs: attrs[:3],
		},
		{
			name:      "deduplicate/under_capacity",
			input:     [][]attribute.KeyValue{attrs[:1], attrs[:1], attrs[:1]},
			wantAttrs: attrs[:1],
		},
		{
			name:      "deduplicate/over_capacity",
			input:     [][]attribute.KeyValue{attrs[:1], attrs[:1], attrs[:1], attrs[:3]},
			wantAttrs: attrs[:3],
		},
		{
			name: "deduplicate/added",
			input: [][]attribute.KeyValue{
				attrs[:2],
				{attrs[2], attrs[2], attrs[2]},
			},
			wantAttrs: attrs[:3],
		},
		{
			name: "deduplicate/added_at_cappacity",
			input: [][]attribute.KeyValue{
				attrs[:3],
				{attrs[2], attrs[2], attrs[2]},
			},
			wantAttrs: attrs[:3],
		},
		{
			name: "invalid",
			input: [][]attribute.KeyValue{
				{invalid},
			},
			wantDropped: 1,
		},
		{
			name: "invalid_with_valid",
			input: [][]attribute.KeyValue{
				{invalid, attrs[0]},
			},
			wantAttrs:   attrs[:1],
			wantDropped: 1,
		},
		{
			name: "invalid_over_capacity",
			input: [][]attribute.KeyValue{
				{invalid, invalid, invalid, invalid, attrs[0]},
			},
			wantAttrs:   attrs[:1],
			wantDropped: 4,
		},
		{
			name: "valid:invalid/under_capacity",
			input: [][]attribute.KeyValue{
				attrs[:1],
				{invalid},
			},
			wantAttrs:   attrs[:1],
			wantDropped: 1,
		},
		{
			name: "valid:invalid/over_capacity",
			input: [][]attribute.KeyValue{
				attrs[:1],
				{invalid, invalid, invalid, invalid},
			},
			wantAttrs:   attrs[:1],
			wantDropped: 4,
		},
		{
			name: "valid_at_capacity:invalid",
			input: [][]attribute.KeyValue{
				attrs[:3],
				{invalid, invalid, invalid, invalid},
			},
			wantAttrs:   attrs[:3],
			wantDropped: 4,
		},
	}

	const (
		capacity = 3
		instName = "TestSpanAttributeCapacity"
		spanName = "test span"
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			te := NewTestExporter()
			sl := NewSpanLimits()
			sl.AttributeCountLimit = capacity
			tp := NewTracerProvider(WithSyncer(te), WithSpanLimits(sl))
			_, span := tp.Tracer(instName).Start(context.Background(), spanName)
			for _, a := range test.input {
				span.SetAttributes(a...)
			}
			span.End()

			require.Implements(t, (*ReadOnlySpan)(nil), span)
			roSpan := span.(ReadOnlySpan)

			// Ensure the span itself is valid.
			assert.ElementsMatch(t, test.wantAttrs, roSpan.Attributes(), "expected attributes")
			assert.Equal(t, test.wantDropped, roSpan.DroppedAttributes(), "dropped attributes")

			snap, ok := te.GetSpan(spanName)
			require.Truef(t, ok, "span %s not exported", spanName)

			// Ensure the exported span snapshot is valid.
			assert.ElementsMatch(t, test.wantAttrs, snap.Attributes(), "expected attributes")
			assert.Equal(t, test.wantDropped, snap.DroppedAttributes(), "dropped attributes")
		})
	}
}

func TestEvents(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))

	span := startSpan(tp, "Events")
	k1v1 := attribute.String("key1", "value1")
	k2v2 := attribute.Bool("key2", true)
	k3v3 := attribute.Int64("key3", 3)

	span.AddEvent("foo", trace.WithAttributes(attribute.String("key1", "value1")))
	span.AddEvent("bar", trace.WithAttributes(
		attribute.Bool("key2", true),
		attribute.Int64("key3", 3),
	))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.Events() {
		if !checkTime(&got.Events()[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent: sc.WithRemote(true),
		name:   "span0",
		events: []Event{
			{Name: "foo", Attributes: []attribute.KeyValue{k1v1}},
			{Name: "bar", Attributes: []attribute.KeyValue{k2v2, k3v3}},
		},
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "Events"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Events: -got +want %s", diff)
	}
}

func TestEventsOverLimit(t *testing.T) {
	te := NewTestExporter()
	sl := NewSpanLimits()
	sl.EventCountLimit = 2
	tp := NewTracerProvider(WithSpanLimits(sl), WithSyncer(te), WithResource(resource.Empty()))

	span := startSpan(tp, "EventsOverLimit")
	k1v1 := attribute.String("key1", "value1")
	k2v2 := attribute.Bool("key2", false)
	k3v3 := attribute.String("key3", "value3")

	span.AddEvent("fooDrop", trace.WithAttributes(attribute.String("key1", "value1")))
	span.AddEvent("barDrop", trace.WithAttributes(
		attribute.Bool("key2", true),
		attribute.String("key3", "value3"),
	))
	span.AddEvent("foo", trace.WithAttributes(attribute.String("key1", "value1")))
	span.AddEvent("bar", trace.WithAttributes(
		attribute.Bool("key2", false),
		attribute.String("key3", "value3"),
	))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.Events() {
		if !checkTime(&got.Events()[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent: sc.WithRemote(true),
		name:   "span0",
		events: []Event{
			{Name: "foo", Attributes: []attribute.KeyValue{k1v1}},
			{Name: "bar", Attributes: []attribute.KeyValue{k2v2, k3v3}},
		},
		droppedEventCount:    2,
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "EventsOverLimit"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Event over limit: -got +want %s", diff)
	}
}

func TestLinks(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))

	k1v1 := attribute.String("key1", "value1")
	k2v2 := attribute.String("key2", "value2")
	k3v3 := attribute.String("key3", "value3")

	sc1 := trace.NewSpanContext(
		trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}},
	)
	sc2 := trace.NewSpanContext(
		trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}},
	)

	l1 := trace.Link{SpanContext: sc1, Attributes: []attribute.KeyValue{k1v1}}
	l2 := trace.Link{SpanContext: sc2, Attributes: []attribute.KeyValue{k2v2, k3v3}}

	links := []trace.Link{l1, l2}
	span := startSpan(tp, "Links", trace.WithLinks(links...))

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:               sc.WithRemote(true),
		name:                 "span0",
		links:                []Link{{l1.SpanContext, l1.Attributes, 0}, {l2.SpanContext, l2.Attributes, 0}},
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "Links"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link: -got +want %s", diff)
	}
	sc1 = trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}})

	span1 := startSpan(tp, "name", trace.WithLinks([]trace.Link{
		{SpanContext: trace.SpanContext{}},
		{SpanContext: sc1},
	}...))

	sdkspan, _ := span1.(*recordingSpan)
	require.Len(t, sdkspan.Links(), 1)
}

func TestLinksOverLimit(t *testing.T) {
	te := NewTestExporter()

	sc1 := trace.NewSpanContext(
		trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}},
	)
	sc2 := trace.NewSpanContext(
		trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}},
	)
	sc3 := trace.NewSpanContext(
		trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{1, 1}), SpanID: trace.SpanID{3}},
	)

	sl := NewSpanLimits()
	sl.LinkCountLimit = 2
	tp := NewTracerProvider(WithSpanLimits(sl), WithSyncer(te), WithResource(resource.Empty()))

	span := startSpan(tp, "LinksOverLimit",
		trace.WithLinks(
			trace.Link{SpanContext: sc1, Attributes: []attribute.KeyValue{attribute.String("key1", "value1")}},
			trace.Link{SpanContext: sc2, Attributes: []attribute.KeyValue{attribute.String("key2", "value2")}},
			trace.Link{SpanContext: sc3, Attributes: []attribute.KeyValue{attribute.String("key3", "value3")}},
		),
	)

	k2v2 := attribute.String("key2", "value2")
	k3v3 := attribute.String("key3", "value3")

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent: sc.WithRemote(true),
		name:   "span0",
		links: []Link{
			{SpanContext: sc2, Attributes: []attribute.KeyValue{k2v2}, DroppedAttributeCount: 0},
			{SpanContext: sc3, Attributes: []attribute.KeyValue{k3v3}, DroppedAttributeCount: 0},
		},
		droppedLinkCount:     1,
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "LinksOverLimit"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link over limit: -got +want %s", diff)
	}
}

func TestSetSpanName(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	ctx := context.Background()

	want := "SpanName-1"
	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	_, span := tp.Tracer("SetSpanName").Start(ctx, "SpanName-1")
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name() != want {
		t.Errorf("span.Name: got %q; want %q", got.Name(), want)
	}
}

func TestSetSpanStatus(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))

	span := startSpan(tp, "SpanStatus")
	span.SetStatus(codes.Error, "Error")
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:   sc.WithRemote(true),
		name:     "span0",
		spanKind: trace.SpanKindInternal,
		status: Status{
			Code:        codes.Error,
			Description: "Error",
		},
		instrumentationScope: instrumentation.Scope{Name: "SpanStatus"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanStatus: -got +want %s", diff)
	}
}

func TestSetSpanStatusWithoutMessageWhenStatusIsNotError(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))

	span := startSpan(tp, "SpanStatus")
	span.SetStatus(codes.Ok, "This message will be ignored")
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:   sc.WithRemote(true),
		name:     "span0",
		spanKind: trace.SpanKindInternal,
		status: Status{
			Code:        codes.Ok,
			Description: "",
		},
		instrumentationScope: instrumentation.Scope{Name: "SpanStatus"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanStatus: -got +want %s", diff)
	}
}

func cmpDiff(x, y interface{}) string {
	return cmp.Diff(x, y,
		cmp.AllowUnexported(snapshot{}),
		cmp.AllowUnexported(attribute.Value{}),
		cmp.AllowUnexported(Event{}),
		cmp.AllowUnexported(trace.TraceState{}),
		cmp.Comparer(func(x, y attribute.Set) bool {
			return x.Equals(&y)
		}),
	)
}

// checkChild is test utility function that tests that c has fields set appropriately,
// given that it is a child span of p.
func checkChild(t *testing.T, p trace.SpanContext, apiSpan trace.Span) error {
	s := apiSpan.(*recordingSpan)
	if s == nil {
		return fmt.Errorf("got nil child span, want non-nil")
	}
	if got, want := s.spanContext.TraceID().String(), p.TraceID().String(); got != want {
		return fmt.Errorf("got child trace ID %s, want %s", got, want)
	}
	if childID, parentID := s.spanContext.SpanID().String(), p.SpanID().String(); childID == parentID {
		return fmt.Errorf("got child span ID %s, parent span ID %s; want unequal IDs", childID, parentID)
	}
	if got, want := s.spanContext.TraceFlags(), p.TraceFlags(); got != want {
		return fmt.Errorf("got child trace options %d, want %d", got, want)
	}
	got, want := s.spanContext.TraceState(), p.TraceState()
	assert.Equal(t, want, got)
	return nil
}

// startSpan starts a span with a name "span0". See startNamedSpan for
// details.
func startSpan(tp *TracerProvider, trName string, args ...trace.SpanStartOption) trace.Span {
	return startNamedSpan(tp, trName, "span0", args...)
}

// startNamed Span is a test utility func that starts a span with a
// passed name and with remote span context as parent. The remote span
// context contains TraceFlags with sampled bit set. This allows the
// span to be automatically sampled.
func startNamedSpan(tp *TracerProvider, trName, name string, args ...trace.SpanStartOption) trace.Span {
	_, span := tp.Tracer(trName).Start(
		trace.ContextWithRemoteSpanContext(context.Background(), sc),
		name,
		args...,
	)
	return span
}

// startLocalSpan is a test utility func that starts a span with a
// passed name and with the passed context. The context is returned
// along with the span so this parent can be used to create child
// spans.
func startLocalSpan(
	ctx context.Context,
	tp *TracerProvider,
	trName, name string,
	args ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	ctx, span := tp.Tracer(trName).Start(
		ctx,
		name,
		args...,
	)
	return ctx, span
}

// endSpan is a test utility function that ends the span in the context and
// returns the exported span.
// It requires that span be sampled using one of these methods
//  1. Passing parent span context in context
//  2. Use WithSampler(AlwaysSample())
//  3. Configuring AlwaysSample() as default sampler
//
// It also does some basic tests on the span.
// It also clears spanID in the to make the comparison easier.
func endSpan(te *testExporter, span trace.Span) (*snapshot, error) {
	if !span.IsRecording() {
		return nil, fmt.Errorf("method IsRecording: got false, want true")
	}
	if !span.SpanContext().IsSampled() {
		return nil, fmt.Errorf("method IsSampled: got false, want true")
	}
	span.End()
	if te.Len() != 1 {
		return nil, fmt.Errorf("got %d exported spans, want one span", te.Len())
	}
	got := te.Spans()[0]
	if !got.SpanContext().SpanID().IsValid() {
		return nil, fmt.Errorf("exporting span: expected nonzero SpanID")
	}
	got.spanContext = got.SpanContext().WithSpanID(trace.SpanID{})
	if !checkTime(&got.startTime) {
		return nil, fmt.Errorf("exporting span: expected nonzero StartTime")
	}
	if !checkTime(&got.endTime) {
		return nil, fmt.Errorf("exporting span: expected nonzero EndTime")
	}
	return got, nil
}

// checkTime checks that a nonzero time was set in x, then clears it.
func checkTime(x *time.Time) bool {
	if x.IsZero() {
		return false
	}
	*x = time.Time{}
	return true
}

func TestEndSpanTwice(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	st := time.Now()
	et1 := st.Add(100 * time.Millisecond)
	et2 := st.Add(200 * time.Millisecond)

	span := startSpan(tp, "EndSpanTwice", trace.WithTimestamp(st))
	span.End(trace.WithTimestamp(et1))
	span.End(trace.WithTimestamp(et2))

	if te.Len() != 1 {
		t.Fatalf("expected only a single span, got %#v", te.Spans())
	}

	ro := span.(ReadOnlySpan)
	if ro.EndTime() != et1 {
		t.Fatalf("2nd call to End() should not modify end time")
	}
}

func TestStartSpanAfterEnd(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSampler(AlwaysSample()), WithSyncer(te))
	ctx := context.Background()

	tr := tp.Tracer("SpanAfterEnd")
	ctx, span0 := tr.Start(trace.ContextWithRemoteSpanContext(ctx, sc), "parent")
	ctx1, span1 := tr.Start(ctx, "span-1")
	span1.End()
	// Start a new span with the context containing span-1
	// even though span-1 is ended, we still add this as a new child of span-1
	_, span2 := tr.Start(ctx1, "span-2")
	span2.End()
	span0.End()
	if got, want := te.Len(), 3; got != want {
		t.Fatalf("len(%#v) = %d; want %d", te.Spans(), got, want)
	}

	gotParent, ok := te.GetSpan("parent")
	if !ok {
		t.Fatal("parent not recorded")
	}
	gotSpan1, ok := te.GetSpan("span-1")
	if !ok {
		t.Fatal("span-1 not recorded")
	}
	gotSpan2, ok := te.GetSpan("span-2")
	if !ok {
		t.Fatal("span-2 not recorded")
	}

	if got, want := gotSpan1.SpanContext().TraceID(), gotParent.SpanContext().TraceID(); got != want {
		t.Errorf("span-1.TraceID=%q; want %q", got, want)
	}
	if got, want := gotSpan2.SpanContext().TraceID(), gotParent.SpanContext().TraceID(); got != want {
		t.Errorf("span-2.TraceID=%q; want %q", got, want)
	}
	if got, want := gotSpan1.Parent().SpanID(), gotParent.SpanContext().SpanID(); got != want {
		t.Errorf("span-1.ParentSpanID=%q; want %q (parent.SpanID)", got, want)
	}
	if got, want := gotSpan2.Parent().SpanID(), gotSpan1.SpanContext().SpanID(); got != want {
		t.Errorf("span-2.ParentSpanID=%q; want %q (span1.SpanID)", got, want)
	}
}

func TestChildSpanCount(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSampler(AlwaysSample()), WithSyncer(te))

	tr := tp.Tracer("ChidSpanCount")
	ctx, span0 := tr.Start(context.Background(), "parent")
	ctx1, span1 := tr.Start(ctx, "span-1")
	_, span2 := tr.Start(ctx1, "span-2")
	span2.End()
	span1.End()

	_, span3 := tr.Start(ctx, "span-3")
	span3.End()
	span0.End()
	if got, want := te.Len(), 4; got != want {
		t.Fatalf("len(%#v) = %d; want %d", te.Spans(), got, want)
	}

	gotParent, ok := te.GetSpan("parent")
	if !ok {
		t.Fatal("parent not recorded")
	}
	gotSpan1, ok := te.GetSpan("span-1")
	if !ok {
		t.Fatal("span-1 not recorded")
	}
	gotSpan2, ok := te.GetSpan("span-2")
	if !ok {
		t.Fatal("span-2 not recorded")
	}
	gotSpan3, ok := te.GetSpan("span-3")
	if !ok {
		t.Fatal("span-3 not recorded")
	}

	if got, want := gotSpan3.ChildSpanCount(), 0; got != want {
		t.Errorf("span-3.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotSpan2.ChildSpanCount(), 0; got != want {
		t.Errorf("span-2.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotSpan1.ChildSpanCount(), 1; got != want {
		t.Errorf("span-1.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotParent.ChildSpanCount(), 2; got != want {
		t.Errorf("parent.ChildSpanCount=%d; want %d", got, want)
	}
}

func TestNilSpanEnd(t *testing.T) {
	var span *recordingSpan
	span.End()
}

func TestSpanWithCanceledContext(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, span := tp.Tracer(t.Name()).Start(ctx, "span")
	span.End()

	assert.Equal(t, 1, te.Len(), "span recording must ignore context cancellation")
}

func TestNonRecordingSpanDoesNotTrackRuntimeTracerTask(t *testing.T) {
	tp := NewTracerProvider(WithSampler(NeverSample()))
	tr := tp.Tracer("TestNonRecordingSpanDoesNotTrackRuntimeTracerTask")

	_, apiSpan := tr.Start(context.Background(), "foo")
	if _, ok := apiSpan.(runtimeTracer); ok {
		t.Fatalf("non recording span implements runtime trace task tracking")
	}
}

func TestRecordingSpanRuntimeTracerTaskEnd(t *testing.T) {
	tp := NewTracerProvider(WithSampler(AlwaysSample()))
	tr := tp.Tracer("TestRecordingSpanRuntimeTracerTaskEnd")

	var n uint64
	executionTracerTaskEnd := func() {
		atomic.AddUint64(&n, 1)
	}
	_, apiSpan := tr.Start(context.Background(), "foo")
	s, ok := apiSpan.(*recordingSpan)
	if !ok {
		t.Fatal("recording span not returned from always sampled Tracer")
	}

	s.executionTracerTaskEnd = executionTracerTaskEnd
	s.End()

	if n != 1 {
		t.Error("recording span did not end runtime trace task")
	}
}

func TestCustomStartEndTime(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithSampler(AlwaysSample()))

	startTime := time.Date(2019, time.August, 27, 14, 42, 0, 0, time.UTC)
	endTime := startTime.Add(time.Second * 20)
	_, span := tp.Tracer("Custom Start and End time").Start(
		context.Background(),
		"testspan",
		trace.WithTimestamp(startTime),
	)
	span.End(trace.WithTimestamp(endTime))

	if te.Len() != 1 {
		t.Fatalf("got %d exported spans, want one span", te.Len())
	}
	got := te.Spans()[0]
	if !got.StartTime().Equal(startTime) {
		t.Errorf("expected start time to be %s, got %s", startTime, got.StartTime())
	}
	if !got.EndTime().Equal(endTime) {
		t.Errorf("expected end time to be %s, got %s", endTime, got.EndTime())
	}
}

func TestRecordError(t *testing.T) {
	scenarios := []struct {
		err error
		typ string
		msg string
	}{
		{
			err: newTestError("test error"),
			typ: "go.opentelemetry.io/otel/sdk/trace.testError",
			msg: "test error",
		},
		{
			err: errors.New("test error 2"),
			typ: "*errors.errorString",
			msg: "test error 2",
		},
	}

	for _, s := range scenarios {
		te := NewTestExporter()
		tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
		span := startSpan(tp, "RecordError")

		errTime := time.Now()
		span.RecordError(s.err, trace.WithTimestamp(errTime))

		got, err := endSpan(te, span)
		if err != nil {
			t.Fatal(err)
		}

		want := &snapshot{
			spanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    tid,
				TraceFlags: 0x1,
			}),
			parent:   sc.WithRemote(true),
			name:     "span0",
			status:   Status{Code: codes.Unset},
			spanKind: trace.SpanKindInternal,
			events: []Event{
				{
					Name: semconv.ExceptionEventName,
					Time: errTime,
					Attributes: []attribute.KeyValue{
						semconv.ExceptionType(s.typ),
						semconv.ExceptionMessage(s.msg),
					},
				},
			},
			instrumentationScope: instrumentation.Scope{Name: "RecordError"},
		}
		if diff := cmpDiff(got, want); diff != "" {
			t.Errorf("SpanErrorOptions: -got +want %s", diff)
		}
	}
}

func TestRecordErrorWithStackTrace(t *testing.T) {
	err := newTestError("test error")
	typ := "go.opentelemetry.io/otel/sdk/trace.testError"
	msg := "test error"

	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	span := startSpan(tp, "RecordError")

	errTime := time.Now()
	span.RecordError(err, trace.WithTimestamp(errTime), trace.WithStackTrace(true))

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:   sc.WithRemote(true),
		name:     "span0",
		status:   Status{Code: codes.Unset},
		spanKind: trace.SpanKindInternal,
		events: []Event{
			{
				Name: semconv.ExceptionEventName,
				Time: errTime,
				Attributes: []attribute.KeyValue{
					semconv.ExceptionType(typ),
					semconv.ExceptionMessage(msg),
				},
			},
		},
		instrumentationScope: instrumentation.Scope{Name: "RecordError"},
	}

	assert.Equal(t, want.spanContext, got.spanContext)
	assert.Equal(t, want.parent, got.parent)
	assert.Equal(t, want.name, got.name)
	assert.Equal(t, want.status, got.status)
	assert.Equal(t, want.spanKind, got.spanKind)
	assert.Equal(t, got.events[0].Attributes[0].Value.AsString(), want.events[0].Attributes[0].Value.AsString())
	assert.Equal(t, got.events[0].Attributes[1].Value.AsString(), want.events[0].Attributes[1].Value.AsString())
	gotStackTraceFunctionName := strings.Split(got.events[0].Attributes[2].Value.AsString(), "\n")

	assert.Truef(
		t,
		strings.HasPrefix(gotStackTraceFunctionName[1], "go.opentelemetry.io/otel/sdk/trace.recordStackTrace"),
		"%q not prefixed with go.opentelemetry.io/otel/sdk/trace.recordStackTrace",
		gotStackTraceFunctionName[1],
	)
	assert.Truef(
		t,
		strings.HasPrefix(
			gotStackTraceFunctionName[3],
			"go.opentelemetry.io/otel/sdk/trace.(*recordingSpan).RecordError",
		),
		"%q not prefixed with go.opentelemetry.io/otel/sdk/trace.(*recordingSpan).RecordError",
		gotStackTraceFunctionName[3],
	)
}

func TestRecordErrorNil(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	span := startSpan(tp, "RecordErrorNil")

	span.RecordError(nil)

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:   sc.WithRemote(true),
		name:     "span0",
		spanKind: trace.SpanKindInternal,
		status: Status{
			Code:        codes.Unset,
			Description: "",
		},
		instrumentationScope: instrumentation.Scope{Name: "RecordErrorNil"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SpanErrorOptions: -got +want %s", diff)
	}
}

func TestWithSpanKind(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithSampler(AlwaysSample()), WithResource(resource.Empty()))
	tr := tp.Tracer("withSpanKind")

	_, span := tr.Start(context.Background(), "WithoutSpanKind")
	spanData, err := endSpan(te, span)
	if err != nil {
		t.Error(err.Error())
	}

	if spanData.SpanKind() != trace.SpanKindInternal {
		t.Errorf(
			"Default value of Spankind should be Internal: got %+v, want %+v\n",
			spanData.SpanKind(),
			trace.SpanKindInternal,
		)
	}

	sks := []trace.SpanKind{
		trace.SpanKindInternal,
		trace.SpanKindServer,
		trace.SpanKindClient,
		trace.SpanKindProducer,
		trace.SpanKindConsumer,
	}

	for _, sk := range sks {
		te.Reset()

		_, span := tr.Start(context.Background(), fmt.Sprintf("SpanKind-%v", sk), trace.WithSpanKind(sk))
		spanData, err := endSpan(te, span)
		if err != nil {
			t.Error(err.Error())
		}

		if spanData.SpanKind() != sk {
			t.Errorf("WithSpanKind check: got %+v, want %+v\n", spanData.SpanKind(), sks)
		}
	}
}

func mergeResource(t *testing.T, r1, r2 *resource.Resource) *resource.Resource {
	r, err := resource.Merge(r1, r2)
	assert.NoError(t, err)
	return r
}

func TestWithResource(t *testing.T) {
	t.Setenv(envVarResourceAttributes, "key=value,rk5=7")
	cases := []struct {
		name    string
		options []TracerProviderOption
		want    *resource.Resource
		msg     string
	}{
		{
			name:    "explicitly empty resource",
			options: []TracerProviderOption{WithResource(resource.Empty())},
			want:    resource.Environment(),
		},
		{
			name:    "uses default if no resource option",
			options: []TracerProviderOption{},
			want:    resource.Default(),
		},
		{
			name: "explicit resource",
			options: []TracerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5)),
			),
		},
		{
			name: "last resource wins",
			options: []TracerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "vk1"), attribute.Int64("rk2", 5))),
				WithResource(resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10)),
			),
		},
		{
			name: "overlapping attributes with environment resource",
			options: []TracerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10)),
			),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			te := NewTestExporter()
			defaultOptions := []TracerProviderOption{WithSyncer(te), WithSampler(AlwaysSample())}
			tp := NewTracerProvider(append(defaultOptions, tc.options...)...)
			span := startSpan(tp, "WithResource")
			span.SetAttributes(attribute.String("key1", "value1"))
			got, err := endSpan(te, span)
			if err != nil {
				t.Error(err.Error())
			}
			want := &snapshot{
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent: sc.WithRemote(true),
				name:   "span0",
				attributes: []attribute.KeyValue{
					attribute.String("key1", "value1"),
				},
				spanKind:             trace.SpanKindInternal,
				resource:             tc.want,
				instrumentationScope: instrumentation.Scope{Name: "WithResource"},
			}
			if diff := cmpDiff(got, want); diff != "" {
				t.Errorf("WithResource:\n  -got +want %s", diff)
			}
		})
	}
}

func TestWithInstrumentationVersionAndSchema(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))

	ctx := context.Background()
	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	_, span := tp.Tracer(
		"WithInstrumentationVersion",
		trace.WithInstrumentationVersion("v0.1.0"),
		trace.WithSchemaURL("https://opentelemetry.io/schemas/1.2.0"),
	).Start(ctx, "span0")
	got, err := endSpan(te, span)
	if err != nil {
		t.Error(err.Error())
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:   sc.WithRemote(true),
		name:     "span0",
		spanKind: trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{
			Name:      "WithInstrumentationVersion",
			Version:   "v0.1.0",
			SchemaURL: "https://opentelemetry.io/schemas/1.2.0",
		},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("WithResource:\n  -got +want %s", diff)
	}
}

func TestSpanCapturesPanic(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	_, span := tp.Tracer("CatchPanic").Start(
		context.Background(),
		"span",
	)

	f := func() {
		defer span.End()
		panic(errors.New("error message"))
	}
	require.PanicsWithError(t, "error message", f)
	spans := te.Spans()
	require.Len(t, spans, 1)
	require.Len(t, spans[0].Events(), 1)
	assert.Equal(t, semconv.ExceptionEventName, spans[0].Events()[0].Name)
	assert.Equal(t, []attribute.KeyValue{
		semconv.ExceptionType("*errors.errorString"),
		semconv.ExceptionMessage("error message"),
	}, spans[0].Events()[0].Attributes)
}

func TestSpanCapturesPanicWithStackTrace(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	_, span := tp.Tracer("CatchPanic").Start(
		context.Background(),
		"span",
	)

	f := func() {
		defer span.End(trace.WithStackTrace(true))
		panic(errors.New("error message"))
	}
	require.PanicsWithError(t, "error message", f)
	spans := te.Spans()
	require.Len(t, spans, 1)
	require.Len(t, spans[0].Events(), 1)
	assert.Equal(t, semconv.ExceptionEventName, spans[0].Events()[0].Name)
	assert.Equal(t, "*errors.errorString", spans[0].Events()[0].Attributes[0].Value.AsString())
	assert.Equal(t, "error message", spans[0].Events()[0].Attributes[1].Value.AsString())

	gotStackTraceFunctionName := strings.Split(spans[0].Events()[0].Attributes[2].Value.AsString(), "\n")
	assert.Truef(
		t,
		strings.HasPrefix(gotStackTraceFunctionName[1], "go.opentelemetry.io/otel/sdk/trace.recordStackTrace"),
		"%q not prefixed with go.opentelemetry.io/otel/sdk/trace.recordStackTrace",
		gotStackTraceFunctionName[1],
	)
	assert.Truef(
		t,
		strings.HasPrefix(gotStackTraceFunctionName[3], "go.opentelemetry.io/otel/sdk/trace.(*recordingSpan).End"),
		"%q not prefixed with go.opentelemetry.io/otel/sdk/trace.(*recordingSpan).End",
		gotStackTraceFunctionName[3],
	)
}

func TestReadOnlySpan(t *testing.T) {
	kv := attribute.String("foo", "bar")

	tp := NewTracerProvider(WithResource(resource.NewSchemaless(kv)))
	tr := tp.Tracer("ReadOnlySpan", trace.WithInstrumentationVersion("3"))

	// Initialize parent context.
	tID, sID := tp.idGenerator.NewIDs(context.Background())
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: 0x1,
		Remote:     true,
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)

	// Initialize linked context.
	tID, sID = tp.idGenerator.NewIDs(context.Background())
	linked := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: 0x1,
	})

	st := time.Now()
	ctx, s := tr.Start(ctx, "foo", trace.WithTimestamp(st),
		trace.WithLinks(trace.Link{SpanContext: linked}))
	s.SetAttributes(kv)
	s.AddEvent("foo", trace.WithAttributes(kv))
	s.SetStatus(codes.Ok, "foo")

	// Verify span implements ReadOnlySpan.
	ro, ok := s.(ReadOnlySpan)
	require.True(t, ok)

	assert.Equal(t, "foo", ro.Name())
	assert.Equal(t, trace.SpanContextFromContext(ctx), ro.SpanContext())
	assert.Equal(t, parent, ro.Parent())
	assert.Equal(t, trace.SpanKindInternal, ro.SpanKind())
	assert.Equal(t, st, ro.StartTime())
	assert.True(t, ro.EndTime().IsZero())
	assert.Equal(t, kv.Key, ro.Attributes()[0].Key)
	assert.Equal(t, kv.Value, ro.Attributes()[0].Value)
	assert.Equal(t, linked, ro.Links()[0].SpanContext)
	assert.Equal(t, kv.Key, ro.Events()[0].Attributes[0].Key)
	assert.Equal(t, kv.Value, ro.Events()[0].Attributes[0].Value)
	assert.Equal(t, codes.Ok, ro.Status().Code)
	assert.Empty(t, ro.Status().Description)
	assert.Equal(t, "ReadOnlySpan", ro.InstrumentationLibrary().Name)
	assert.Equal(t, "3", ro.InstrumentationLibrary().Version)
	assert.Equal(t, "ReadOnlySpan", ro.InstrumentationScope().Name)
	assert.Equal(t, "3", ro.InstrumentationScope().Version)
	assert.Equal(t, kv.Key, ro.Resource().Attributes()[0].Key)
	assert.Equal(t, kv.Value, ro.Resource().Attributes()[0].Value)

	// Verify changes to the original span are reflected in the ReadOnlySpan.
	s.SetName("bar")
	assert.Equal(t, "bar", ro.Name())

	// Verify snapshot() returns snapshots that are independent from the
	// original span and from one another.
	d1 := s.(*recordingSpan).snapshot()
	s.AddEvent("baz")
	d2 := s.(*recordingSpan).snapshot()
	for _, e := range d1.Events() {
		if e.Name == "baz" {
			t.Errorf("Didn't expect to find 'baz' event")
		}
	}
	var exists bool
	for _, e := range d2.Events() {
		if e.Name == "baz" {
			exists = true
		}
	}
	if !exists {
		t.Errorf("Expected to find 'baz' event")
	}

	et := st.Add(time.Millisecond)
	s.End(trace.WithTimestamp(et))
	assert.Equal(t, et, ro.EndTime())
}

func TestReadWriteSpan(t *testing.T) {
	tp := NewTracerProvider(WithResource(resource.Empty()))
	tr := tp.Tracer("ReadWriteSpan")

	// Initialize parent context.
	tID, sID := tp.idGenerator.NewIDs(context.Background())
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tID,
		SpanID:     sID,
		TraceFlags: 0x1,
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)

	_, span := tr.Start(ctx, "foo")
	defer span.End()

	// Verify span implements ReadOnlySpan.
	rw, ok := span.(ReadWriteSpan)
	require.True(t, ok)

	// Verify the span can be read from.
	assert.False(t, rw.StartTime().IsZero())

	// Verify the span can be written to.
	rw.SetName("bar")
	assert.Equal(t, "bar", rw.Name())
	// NOTE: This function tests ReadWriteSpan which is an interface which
	// embeds trace.Span and ReadOnlySpan. Since both of these interfaces have
	// their own tests, there is no point in testing all the possible methods
	// available via ReadWriteSpan as doing so would mean creating a lot of
	// duplication.
}

func TestAddEventsWithMoreAttributesThanLimit(t *testing.T) {
	te := NewTestExporter()
	sl := NewSpanLimits()
	sl.AttributePerEventCountLimit = 2
	tp := NewTracerProvider(
		WithSpanLimits(sl),
		WithSyncer(te),
		WithResource(resource.Empty()),
	)

	span := startSpan(tp, "AddSpanEventWithOverLimitedAttributes")
	span.AddEvent("test1", trace.WithAttributes(
		attribute.Bool("key1", true),
		attribute.String("key2", "value2"),
	))
	// Parts of the attribute should be discard
	span.AddEvent("test2", trace.WithAttributes(
		attribute.Bool("key1", true),
		attribute.String("key2", "value2"),
		attribute.String("key3", "value3"),
		attribute.String("key4", "value4"),
	))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.Events() {
		if !checkTime(&got.Events()[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &snapshot{
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:     sc.WithRemote(true),
		name:       "span0",
		attributes: nil,
		events: []Event{
			{
				Name: "test1",
				Attributes: []attribute.KeyValue{
					attribute.Bool("key1", true),
					attribute.String("key2", "value2"),
				},
			},
			{
				Name: "test2",
				Attributes: []attribute.KeyValue{
					attribute.Bool("key1", true),
					attribute.String("key2", "value2"),
				},
				DroppedAttributeCount: 2,
			},
		},
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "AddSpanEventWithOverLimitedAttributes"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOverLimit: -got +want %s", diff)
	}
}

type stateSampler struct {
	prefix string
	f      func(trace.TraceState) trace.TraceState
}

func (s *stateSampler) ShouldSample(p SamplingParameters) SamplingResult {
	decision := Drop
	if strings.HasPrefix(p.Name, s.prefix) {
		decision = RecordAndSample
	}
	ts := s.f(trace.SpanContextFromContext(p.ParentContext).TraceState())
	return SamplingResult{Decision: decision, Tracestate: ts}
}

func (s stateSampler) Description() string {
	return "stateSampler"
}

// Check that a new span propagates the SamplerResult.TraceState.
func TestSamplerTraceState(t *testing.T) {
	mustTS := func(ts trace.TraceState, err error) trace.TraceState {
		require.NoError(t, err)
		return ts
	}
	makeInserter := func(k, v, prefix string) Sampler {
		return &stateSampler{
			prefix: prefix,
			f:      func(t trace.TraceState) trace.TraceState { return mustTS(t.Insert(k, v)) },
		}
	}
	makeDeleter := func(k, prefix string) Sampler {
		return &stateSampler{
			prefix: prefix,
			f:      func(t trace.TraceState) trace.TraceState { return t.Delete(k) },
		}
	}
	clearer := func(prefix string) Sampler {
		return &stateSampler{
			prefix: prefix,
			f:      func(t trace.TraceState) trace.TraceState { return trace.TraceState{} },
		}
	}

	tests := []struct {
		name       string
		sampler    Sampler
		spanName   string
		input      trace.TraceState
		want       trace.TraceState
		exportSpan bool
	}{
		{
			name:       "alwaysOn",
			sampler:    AlwaysSample(),
			input:      mustTS(trace.ParseTraceState("k1=v1")),
			want:       mustTS(trace.ParseTraceState("k1=v1")),
			exportSpan: true,
		},
		{
			name:       "alwaysOff",
			sampler:    NeverSample(),
			input:      mustTS(trace.ParseTraceState("k1=v1")),
			want:       mustTS(trace.ParseTraceState("k1=v1")),
			exportSpan: false,
		},
		{
			name:       "insertKeySampled",
			sampler:    makeInserter("k2", "v2", "span"),
			spanName:   "span0",
			input:      mustTS(trace.ParseTraceState("k1=v1")),
			want:       mustTS(trace.ParseTraceState("k2=v2,k1=v1")),
			exportSpan: true,
		},
		{
			name:       "insertKeyDropped",
			sampler:    makeInserter("k2", "v2", "span"),
			spanName:   "nospan0",
			input:      mustTS(trace.ParseTraceState("k1=v1")),
			want:       mustTS(trace.ParseTraceState("k2=v2,k1=v1")),
			exportSpan: false,
		},
		{
			name:       "deleteKeySampled",
			sampler:    makeDeleter("k1", "span"),
			spanName:   "span0",
			input:      mustTS(trace.ParseTraceState("k1=v1,k2=v2")),
			want:       mustTS(trace.ParseTraceState("k2=v2")),
			exportSpan: true,
		},
		{
			name:       "deleteKeyDropped",
			sampler:    makeDeleter("k1", "span"),
			spanName:   "nospan0",
			input:      mustTS(trace.ParseTraceState("k1=v1,k2=v2,k3=v3")),
			want:       mustTS(trace.ParseTraceState("k2=v2,k3=v3")),
			exportSpan: false,
		},
		{
			name:       "clearer",
			sampler:    clearer("span"),
			spanName:   "span0",
			input:      mustTS(trace.ParseTraceState("k1=v1,k3=v3")),
			want:       mustTS(trace.ParseTraceState("")),
			exportSpan: true,
		},
	}

	for _, ts := range tests {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			te := NewTestExporter()
			tp := NewTracerProvider(WithSampler(ts.sampler), WithSyncer(te), WithResource(resource.Empty()))
			tr := tp.Tracer("TraceState")

			sc1 := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    tid,
				SpanID:     sid,
				TraceFlags: trace.FlagsSampled,
				TraceState: ts.input,
			})
			ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc1)
			_, span := tr.Start(ctx, ts.spanName)

			// span's TraceState should be set regardless of Sampled/NonSampled state.
			require.Equal(t, ts.want, span.SpanContext().TraceState())

			span.End()

			got := te.Spans()
			if len(got) > 0 != ts.exportSpan {
				t.Errorf("unexpected number of exported spans %d", len(got))
			}
			if len(got) == 0 {
				return
			}

			receivedState := got[0].SpanContext().TraceState()

			if diff := cmpDiff(receivedState, ts.want); diff != "" {
				t.Errorf("TraceState not propagated: -got +want %s", diff)
			}
		})
	}
}

type testIDGenerator struct {
	traceID int
	spanID  int
}

func (gen *testIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	traceIDHex := fmt.Sprintf("%032x", gen.traceID)
	traceID, _ := trace.TraceIDFromHex(traceIDHex)
	gen.traceID++

	spanID := gen.NewSpanID(ctx, traceID)
	return traceID, spanID
}

func (gen *testIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	spanIDHex := fmt.Sprintf("%016x", gen.spanID)
	spanID, _ := trace.SpanIDFromHex(spanIDHex)
	gen.spanID++
	return spanID
}

var _ IDGenerator = (*testIDGenerator)(nil)

func TestWithIDGenerator(t *testing.T) {
	const (
		startTraceID = 1
		startSpanID  = 10
		numSpan      = 5
	)

	gen := &testIDGenerator{traceID: startTraceID, spanID: startSpanID}
	te := NewTestExporter()
	tp := NewTracerProvider(
		WithSyncer(te),
		WithIDGenerator(gen),
	)
	for i := 0; i < numSpan; i++ {
		func() {
			_, span := tp.Tracer(t.Name()).Start(context.Background(), strconv.Itoa(i))
			defer span.End()

			gotSpanID, err := strconv.ParseUint(span.SpanContext().SpanID().String(), 16, 64)
			require.NoError(t, err)
			assert.Equal(t, uint64(startSpanID+i), gotSpanID)

			gotTraceID, err := strconv.ParseUint(span.SpanContext().TraceID().String(), 16, 64)
			require.NoError(t, err)
			assert.Equal(t, uint64(startTraceID+i), gotTraceID)
		}()
	}
}

func TestEmptyRecordingSpanAttributes(t *testing.T) {
	assert.Nil(t, (&recordingSpan{}).Attributes())
}

func TestEmptyRecordingSpanDroppedAttributes(t *testing.T) {
	assert.Equal(t, 0, (&recordingSpan{}).DroppedAttributes())
}

func TestSpanAddLink(t *testing.T) {
	tests := []struct {
		name               string
		attrLinkCountLimit int
		link               trace.Link
		want               *snapshot
	}{
		{
			name:               "AddLinkWithInvalidSpanContext",
			attrLinkCountLimit: 128,
			link: trace.Link{
				SpanContext: trace.NewSpanContext(
					trace.SpanContextConfig{TraceID: trace.TraceID([16]byte{}), SpanID: [8]byte{}},
				),
			},
			want: &snapshot{
				name: "span0",
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent:               sc.WithRemote(true),
				links:                nil,
				spanKind:             trace.SpanKindInternal,
				instrumentationScope: instrumentation.Scope{Name: "AddLinkWithInvalidSpanContext"},
			},
		},
		{
			name:               "AddLink",
			attrLinkCountLimit: 128,
			link: trace.Link{
				SpanContext: sc,
				Attributes:  []attribute.KeyValue{{Key: "k1", Value: attribute.StringValue("v1")}},
			},
			want: &snapshot{
				name: "span0",
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent: sc.WithRemote(true),
				links: []Link{
					{
						SpanContext: sc,
						Attributes:  []attribute.KeyValue{{Key: "k1", Value: attribute.StringValue("v1")}},
					},
				},
				spanKind:             trace.SpanKindInternal,
				instrumentationScope: instrumentation.Scope{Name: "AddLink"},
			},
		},
		{
			name:               "AddLinkWithMoreAttributesThanLimit",
			attrLinkCountLimit: 1,
			link: trace.Link{
				SpanContext: sc,
				Attributes: []attribute.KeyValue{
					{Key: "k1", Value: attribute.StringValue("v1")},
					{Key: "k2", Value: attribute.StringValue("v2")},
					{Key: "k3", Value: attribute.StringValue("v3")},
					{Key: "k4", Value: attribute.StringValue("v4")},
				},
			},
			want: &snapshot{
				name: "span0",
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent: sc.WithRemote(true),
				links: []Link{
					{
						SpanContext:           sc,
						Attributes:            []attribute.KeyValue{{Key: "k1", Value: attribute.StringValue("v1")}},
						DroppedAttributeCount: 3,
					},
				},
				spanKind:             trace.SpanKindInternal,
				instrumentationScope: instrumentation.Scope{Name: "AddLinkWithMoreAttributesThanLimit"},
			},
		},
		{
			name:               "AddLinkWithAttributesEmptySpanContext",
			attrLinkCountLimit: 128,
			link: trace.Link{
				Attributes: []attribute.KeyValue{{Key: "k1", Value: attribute.StringValue("v1")}},
			},
			want: &snapshot{
				name: "span0",
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent: sc.WithRemote(true),
				links: []Link{
					{
						Attributes: []attribute.KeyValue{{Key: "k1", Value: attribute.StringValue("v1")}},
					},
				},
				spanKind:             trace.SpanKindInternal,
				instrumentationScope: instrumentation.Scope{Name: "AddLinkWithAttributesEmptySpanContext"},
			},
		},
		{
			name:               "AddLinkWithTraceStateEmptySpanContext",
			attrLinkCountLimit: 128,
			link: trace.Link{
				SpanContext: trace.SpanContext{}.WithTraceState(ts),
			},
			want: &snapshot{
				name: "span0",
				spanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tid,
					TraceFlags: 0x1,
				}),
				parent: sc.WithRemote(true),
				links: []Link{
					{
						SpanContext: trace.SpanContext{}.WithTraceState(ts),
					},
				},
				spanKind:             trace.SpanKindInternal,
				instrumentationScope: instrumentation.Scope{Name: "AddLinkWithTraceStateEmptySpanContext"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			te := NewTestExporter()
			sl := NewSpanLimits()
			sl.AttributePerLinkCountLimit = tc.attrLinkCountLimit

			tp := NewTracerProvider(WithSpanLimits(sl), WithSyncer(te), WithResource(resource.Empty()))

			span := startSpan(tp, tc.name)
			span.AddLink(tc.link)

			got, err := endSpan(te, span)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmpDiff(got, tc.want); diff != "" {
				t.Errorf("-got +want %s", diff)
			}
		})
	}
}

func TestAddLinkToNonRecordingSpan(t *testing.T) {
	te := NewTestExporter()
	sl := NewSpanLimits()
	tp := NewTracerProvider(
		WithSpanLimits(sl),
		WithSyncer(te),
		WithResource(resource.Empty()),
	)

	attrs := []attribute.KeyValue{{Key: "k", Value: attribute.StringValue("v")}}

	span := startSpan(tp, "AddLinkToNonRecordingSpan")
	_, err := endSpan(te, span)
	require.NoError(t, err)

	// Add link to ended, non-recording, span. The link should be dropped.
	span.AddLink(trace.Link{
		SpanContext: sc,
		Attributes:  attrs,
	})

	require.Equal(t, 1, te.Len())
	got := te.Spans()[0]
	want := &snapshot{
		name: "span0",
		spanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    tid,
			TraceFlags: 0x1,
		}),
		parent:               sc.WithRemote(true),
		links:                nil,
		spanKind:             trace.SpanKindInternal,
		instrumentationScope: instrumentation.Scope{Name: "AddLinkToNonRecordingSpan"},
	}

	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("AddLinkToNonRecordingSpan: -got +want %s", diff)
	}
}

func BenchmarkTraceStart(b *testing.B) {
	tracer := NewTracerProvider().Tracer("")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.SpanContext{})

	l1 := trace.Link{SpanContext: trace.SpanContext{}, Attributes: []attribute.KeyValue{}}
	l2 := trace.Link{SpanContext: trace.SpanContext{}, Attributes: []attribute.KeyValue{}}

	links := []trace.Link{l1, l2}

	for _, tt := range []struct {
		name    string
		options []trace.SpanStartOption
	}{
		{
			name: "with a simple span",
		},
		{
			name: "with several links",
			options: []trace.SpanStartOption{
				trace.WithLinks(links...),
			},
		},
		{
			name: "with attributes",
			options: []trace.SpanStartOption{
				trace.WithAttributes(
					attribute.String("key1", "value1"),
					attribute.String("key2", "value2"),
				),
			},
		},
	} {
		b.Run(tt.name, func(b *testing.B) {
			spans := make([]trace.Span, b.N)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, span := tracer.Start(ctx, "", tt.options...)
				spans[i] = span
			}

			b.StopTimer()
			for i := 0; i < b.N; i++ {
				spans[i].End()
			}
		})
	}
}
