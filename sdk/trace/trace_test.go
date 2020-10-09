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

package trace

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/apitest"
	ottest "go.opentelemetry.io/otel/internal/testing"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	tid otel.TraceID
	sid otel.SpanID
)

type discardHandler struct{}

func (*discardHandler) Handle(_ error) {}

func init() {
	tid, _ = otel.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ = otel.SpanIDFromHex("0102040810203040")

	global.SetErrorHandler(new(discardHandler))
}

func TestTracerFollowsExpectedAPIBehaviour(t *testing.T) {
	tp := NewTracerProvider(WithConfig(Config{DefaultSampler: TraceIDRatioBased(0)}))
	harness := apitest.NewHarness(t)
	subjectFactory := func() otel.Tracer {
		return tp.Tracer("")
	}

	harness.TestTracer(subjectFactory)
}

type testExporter struct {
	mu    sync.RWMutex
	idx   map[string]int
	spans []*export.SpanData
}

func NewTestExporter() *testExporter {
	return &testExporter{idx: make(map[string]int)}
}

func (te *testExporter) ExportSpans(_ context.Context, spans []*export.SpanData) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	i := len(te.spans)
	for _, s := range spans {
		te.idx[s.Name] = i
		te.spans = append(te.spans, s)
		i++
	}
	return nil
}

func (te *testExporter) Spans() []*export.SpanData {
	te.mu.RLock()
	defer te.mu.RUnlock()

	cp := make([]*export.SpanData, len(te.spans))
	copy(cp, te.spans)
	return cp
}

func (te *testExporter) GetSpan(name string) (*export.SpanData, bool) {
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
	return SamplingResult{Decision: decision, Attributes: []label.KeyValue{label.Int("callCount", ts.callCount)}}
}

func (ts testSampler) Description() string {
	return "testSampler"
}

func TestSetName(t *testing.T) {
	fooSampler := &testSampler{prefix: "foo", t: t}
	tp := NewTracerProvider(WithConfig(Config{DefaultSampler: fooSampler}))

	type testCase struct {
		name          string
		newName       string
		sampledBefore bool
		sampledAfter  bool
	}
	for idx, tt := range []testCase{
		{ // 0
			name:          "foobar",
			newName:       "foobaz",
			sampledBefore: true,
			sampledAfter:  true,
		},
		{ // 1
			name:          "foobar",
			newName:       "barbaz",
			sampledBefore: true,
			sampledAfter:  false,
		},
		{ // 2
			name:          "barbar",
			newName:       "barbaz",
			sampledBefore: false,
			sampledAfter:  false,
		},
		{ // 3
			name:          "barbar",
			newName:       "foobar",
			sampledBefore: false,
			sampledAfter:  true,
		},
	} {
		span := startNamedSpan(tp, "SetName", tt.name)
		if fooSampler.callCount == 0 {
			t.Errorf("%d: the sampler was not even called during span creation", idx)
		}
		fooSampler.callCount = 0
		if gotSampledBefore := span.SpanContext().IsSampled(); tt.sampledBefore != gotSampledBefore {
			t.Errorf("%d: invalid sampling decision before rename, expected %v, got %v", idx, tt.sampledBefore, gotSampledBefore)
		}
		span.SetName(tt.newName)
		if fooSampler.callCount == 0 {
			t.Errorf("%d: the sampler was not even called during span rename", idx)
		}
		fooSampler.callCount = 0
		if gotSampledAfter := span.SpanContext().IsSampled(); tt.sampledAfter != gotSampledAfter {
			t.Errorf("%d: invalid sampling decision after rename, expected %v, got %v", idx, tt.sampledAfter, gotSampledAfter)
		}
		span.End()
	}
}

func TestRecordingIsOn(t *testing.T) {
	tp := NewTracerProvider()
	_, span := tp.Tracer("Recording on").Start(context.Background(), "StartSpan")
	defer span.End()
	if span.IsRecording() == false {
		t.Error("new span is not recording events")
	}
}

func TestSampling(t *testing.T) {
	idg := defIDGenerator()
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
			p := NewTracerProvider(WithConfig(Config{DefaultSampler: tc.sampler}))
			tr := p.Tracer("test")
			var sampled int
			for i := 0; i < total; i++ {
				ctx := context.Background()
				if tc.parent {
					psc := otel.SpanContext{
						TraceID: idg.NewTraceID(),
						SpanID:  idg.NewSpanID(),
					}
					if tc.sampledParent {
						psc.TraceFlags = otel.FlagsSampled
					}
					ctx = otel.ContextWithRemoteSpanContext(ctx, psc)
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

	sc1 := otel.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
	}
	_, s1 := tr.Start(otel.ContextWithRemoteSpanContext(ctx, sc1), "span1-unsampled-parent1")
	if err := checkChild(sc1, s1); err != nil {
		t.Error(err)
	}

	_, s2 := tr.Start(otel.ContextWithRemoteSpanContext(ctx, sc1), "span2-unsampled-parent1")
	if err := checkChild(sc1, s2); err != nil {
		t.Error(err)
	}

	sc2 := otel.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
		//Tracestate:   testTracestate,
	}
	_, s3 := tr.Start(otel.ContextWithRemoteSpanContext(ctx, sc2), "span3-sampled-parent2")
	if err := checkChild(sc2, s3); err != nil {
		t.Error(err)
	}

	ctx2, s4 := tr.Start(otel.ContextWithRemoteSpanContext(ctx, sc2), "span4-sampled-parent2")
	if err := checkChild(sc2, s4); err != nil {
		t.Error(err)
	}

	s4Sc := s4.SpanContext()
	_, s5 := tr.Start(ctx2, "span5-implicit-childof-span4")
	if err := checkChild(s4Sc, s5); err != nil {
		t.Error(err)
	}
}

func TestSetSpanAttributesOnStart(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	span := startSpan(tp,
		"StartSpanAttribute",
		otel.WithAttributes(label.String("key1", "value1")),
		otel.WithAttributes(label.String("key2", "value2")),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []label.KeyValue{
			label.String("key1", "value1"),
			label.String("key2", "value2"),
		},
		SpanKind:               otel.SpanKindInternal,
		HasRemoteParent:        true,
		InstrumentationLibrary: instrumentation.Library{Name: "StartSpanAttribute"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOnStart: -got +want %s", diff)
	}
}

func TestSetSpanAttributes(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	span := startSpan(tp, "SpanAttribute")
	span.SetAttributes(label.String("key1", "value1"))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []label.KeyValue{
			label.String("key1", "value1"),
		},
		SpanKind:               otel.SpanKindInternal,
		HasRemoteParent:        true,
		InstrumentationLibrary: instrumentation.Library{Name: "SpanAttribute"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributes: -got +want %s", diff)
	}
}

func TestSetSpanAttributesOverLimit(t *testing.T) {
	te := NewTestExporter()
	cfg := Config{MaxAttributesPerSpan: 2}
	tp := NewTracerProvider(WithConfig(cfg), WithSyncer(te))

	span := startSpan(tp, "SpanAttributesOverLimit")
	span.SetAttributes(
		label.Bool("key1", true),
		label.String("key2", "value2"),
		label.Bool("key1", false), // Replace key1.
		label.Int64("key4", 4),    // Remove key2 and add key4
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []label.KeyValue{
			label.Bool("key1", false),
			label.Int64("key4", 4),
		},
		SpanKind:               otel.SpanKindInternal,
		HasRemoteParent:        true,
		DroppedAttributeCount:  1,
		InstrumentationLibrary: instrumentation.Library{Name: "SpanAttributesOverLimit"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOverLimit: -got +want %s", diff)
	}
}

func TestEvents(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	span := startSpan(tp, "Events")
	k1v1 := label.String("key1", "value1")
	k2v2 := label.Bool("key2", true)
	k3v3 := label.Int64("key3", 3)

	span.AddEvent(context.Background(), "foo", label.String("key1", "value1"))
	span.AddEvent(context.Background(), "bar",
		label.Bool("key2", true),
		label.Int64("key3", 3),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.MessageEvents {
		if !checkTime(&got.MessageEvents[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		HasRemoteParent: true,
		MessageEvents: []export.Event{
			{Name: "foo", Attributes: []label.KeyValue{k1v1}},
			{Name: "bar", Attributes: []label.KeyValue{k2v2, k3v3}},
		},
		SpanKind:               otel.SpanKindInternal,
		InstrumentationLibrary: instrumentation.Library{Name: "Events"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Events: -got +want %s", diff)
	}
}

func TestEventsOverLimit(t *testing.T) {
	te := NewTestExporter()
	cfg := Config{MaxEventsPerSpan: 2}
	tp := NewTracerProvider(WithConfig(cfg), WithSyncer(te))

	span := startSpan(tp, "EventsOverLimit")
	k1v1 := label.String("key1", "value1")
	k2v2 := label.Bool("key2", false)
	k3v3 := label.String("key3", "value3")

	span.AddEvent(context.Background(), "fooDrop", label.String("key1", "value1"))
	span.AddEvent(context.Background(), "barDrop",
		label.Bool("key2", true),
		label.String("key3", "value3"),
	)
	span.AddEvent(context.Background(), "foo", label.String("key1", "value1"))
	span.AddEvent(context.Background(), "bar",
		label.Bool("key2", false),
		label.String("key3", "value3"),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.MessageEvents {
		if !checkTime(&got.MessageEvents[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		MessageEvents: []export.Event{
			{Name: "foo", Attributes: []label.KeyValue{k1v1}},
			{Name: "bar", Attributes: []label.KeyValue{k2v2, k3v3}},
		},
		DroppedMessageEventCount: 2,
		HasRemoteParent:          true,
		SpanKind:                 otel.SpanKindInternal,
		InstrumentationLibrary:   instrumentation.Library{Name: "EventsOverLimit"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Event over limit: -got +want %s", diff)
	}
}

func TestLinks(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	k1v1 := label.String("key1", "value1")
	k2v2 := label.String("key2", "value2")
	k3v3 := label.String("key3", "value3")

	sc1 := otel.SpanContext{TraceID: otel.TraceID([16]byte{1, 1}), SpanID: otel.SpanID{3}}
	sc2 := otel.SpanContext{TraceID: otel.TraceID([16]byte{1, 1}), SpanID: otel.SpanID{3}}

	links := []otel.Link{
		{SpanContext: sc1, Attributes: []label.KeyValue{k1v1}},
		{SpanContext: sc2, Attributes: []label.KeyValue{k2v2, k3v3}},
	}
	span := startSpan(tp, "Links", otel.WithLinks(links...))

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:           sid,
		Name:                   "span0",
		HasRemoteParent:        true,
		Links:                  links,
		SpanKind:               otel.SpanKindInternal,
		InstrumentationLibrary: instrumentation.Library{Name: "Links"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link: -got +want %s", diff)
	}
}

func TestLinksOverLimit(t *testing.T) {
	te := NewTestExporter()
	cfg := Config{MaxLinksPerSpan: 2}

	sc1 := otel.SpanContext{TraceID: otel.TraceID([16]byte{1, 1}), SpanID: otel.SpanID{3}}
	sc2 := otel.SpanContext{TraceID: otel.TraceID([16]byte{1, 1}), SpanID: otel.SpanID{3}}
	sc3 := otel.SpanContext{TraceID: otel.TraceID([16]byte{1, 1}), SpanID: otel.SpanID{3}}

	tp := NewTracerProvider(WithConfig(cfg), WithSyncer(te))

	span := startSpan(tp, "LinksOverLimit",
		otel.WithLinks(
			otel.Link{SpanContext: sc1, Attributes: []label.KeyValue{label.String("key1", "value1")}},
			otel.Link{SpanContext: sc2, Attributes: []label.KeyValue{label.String("key2", "value2")}},
			otel.Link{SpanContext: sc3, Attributes: []label.KeyValue{label.String("key3", "value3")}},
		),
	)

	k2v2 := label.String("key2", "value2")
	k3v3 := label.String("key3", "value3")

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Links: []otel.Link{
			{SpanContext: sc2, Attributes: []label.KeyValue{k2v2}},
			{SpanContext: sc3, Attributes: []label.KeyValue{k3v3}},
		},
		DroppedLinkCount:       1,
		HasRemoteParent:        true,
		SpanKind:               otel.SpanKindInternal,
		InstrumentationLibrary: instrumentation.Library{Name: "LinksOverLimit"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link over limit: -got +want %s", diff)
	}
}

func TestSetSpanName(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	ctx := context.Background()

	want := "SpanName-1"
	ctx = otel.ContextWithRemoteSpanContext(ctx, otel.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 1,
	})
	_, span := tp.Tracer("SetSpanName").Start(ctx, "SpanName-1")
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != want {
		t.Errorf("span.Name: got %q; want %q", got.Name, want)
	}
}

func TestSetSpanStatus(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	span := startSpan(tp, "SpanStatus")
	span.SetStatus(codes.Error, "Error")
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:           sid,
		Name:                   "span0",
		SpanKind:               otel.SpanKindInternal,
		StatusCode:             codes.Error,
		StatusMessage:          "Error",
		HasRemoteParent:        true,
		InstrumentationLibrary: instrumentation.Library{Name: "SpanStatus"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanStatus: -got +want %s", diff)
	}
}

func cmpDiff(x, y interface{}) string {
	return cmp.Diff(x, y,
		cmp.AllowUnexported(label.Value{}),
		cmp.AllowUnexported(export.Event{}))
}

func remoteSpanContext() otel.SpanContext {
	return otel.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 1,
	}
}

// checkChild is test utility function that tests that c has fields set appropriately,
// given that it is a child span of p.
func checkChild(p otel.SpanContext, apiSpan otel.Span) error {
	s := apiSpan.(*span)
	if s == nil {
		return fmt.Errorf("got nil child span, want non-nil")
	}
	if got, want := s.spanContext.TraceID.String(), p.TraceID.String(); got != want {
		return fmt.Errorf("got child trace ID %s, want %s", got, want)
	}
	if childID, parentID := s.spanContext.SpanID.String(), p.SpanID.String(); childID == parentID {
		return fmt.Errorf("got child span ID %s, parent span ID %s; want unequal IDs", childID, parentID)
	}
	if got, want := s.spanContext.TraceFlags, p.TraceFlags; got != want {
		return fmt.Errorf("got child trace options %d, want %d", got, want)
	}
	// TODO [rgheita] : Fix tracestate test
	//if got, want := c.spanContext.Tracestate, p.Tracestate; got != want {
	//	return fmt.Errorf("got child tracestate %v, want %v", got, want)
	//}
	return nil
}

// startSpan starts a span with a name "span0". See startNamedSpan for
// details.
func startSpan(tp *TracerProvider, trName string, args ...otel.SpanOption) otel.Span {
	return startNamedSpan(tp, trName, "span0", args...)
}

// startNamed Span is a test utility func that starts a span with a
// passed name and with remote span context as parent. The remote span
// context contains TraceFlags with sampled bit set. This allows the
// span to be automatically sampled.
func startNamedSpan(tp *TracerProvider, trName, name string, args ...otel.SpanOption) otel.Span {
	ctx := context.Background()
	ctx = otel.ContextWithRemoteSpanContext(ctx, remoteSpanContext())
	args = append(args, otel.WithRecord())
	_, span := tp.Tracer(trName).Start(
		ctx,
		name,
		args...,
	)
	return span
}

// endSpan is a test utility function that ends the span in the context and
// returns the exported export.SpanData.
// It requires that span be sampled using one of these methods
//  1. Passing parent span context in context
//  2. Use WithSampler(AlwaysSample())
//  3. Configuring AlwaysSample() as default sampler
//
// It also does some basic tests on the span.
// It also clears spanID in the export.SpanData to make the comparison easier.
func endSpan(te *testExporter, span otel.Span) (*export.SpanData, error) {
	if !span.IsRecording() {
		return nil, fmt.Errorf("IsRecording: got false, want true")
	}
	if !span.SpanContext().IsSampled() {
		return nil, fmt.Errorf("IsSampled: got false, want true")
	}
	span.End()
	if te.Len() != 1 {
		return nil, fmt.Errorf("got %d exported spans, want one span", te.Len())
	}
	got := te.Spans()[0]
	if !got.SpanContext.SpanID.IsValid() {
		return nil, fmt.Errorf("exporting span: expected nonzero SpanID")
	}
	got.SpanContext.SpanID = otel.SpanID{}
	if !checkTime(&got.StartTime) {
		return nil, fmt.Errorf("exporting span: expected nonzero StartTime")
	}
	if !checkTime(&got.EndTime) {
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

	span := startSpan(tp, "EndSpanTwice")
	span.End()
	span.End()
	if te.Len() != 1 {
		t.Fatalf("expected only a single span, got %#v", te.Spans())
	}
}

func TestStartSpanAfterEnd(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithConfig(Config{DefaultSampler: AlwaysSample()}), WithSyncer(te))
	ctx := context.Background()

	tr := tp.Tracer("SpanAfterEnd")
	ctx, span0 := tr.Start(otel.ContextWithRemoteSpanContext(ctx, remoteSpanContext()), "parent")
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

	if got, want := gotSpan1.SpanContext.TraceID, gotParent.SpanContext.TraceID; got != want {
		t.Errorf("span-1.TraceID=%q; want %q", got, want)
	}
	if got, want := gotSpan2.SpanContext.TraceID, gotParent.SpanContext.TraceID; got != want {
		t.Errorf("span-2.TraceID=%q; want %q", got, want)
	}
	if got, want := gotSpan1.ParentSpanID, gotParent.SpanContext.SpanID; got != want {
		t.Errorf("span-1.ParentSpanID=%q; want %q (parent.SpanID)", got, want)
	}
	if got, want := gotSpan2.ParentSpanID, gotSpan1.SpanContext.SpanID; got != want {
		t.Errorf("span-2.ParentSpanID=%q; want %q (span1.SpanID)", got, want)
	}
}

func TestChildSpanCount(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithConfig(Config{DefaultSampler: AlwaysSample()}), WithSyncer(te))

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

	if got, want := gotSpan3.ChildSpanCount, 0; got != want {
		t.Errorf("span-3.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotSpan2.ChildSpanCount, 0; got != want {
		t.Errorf("span-2.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotSpan1.ChildSpanCount, 1; got != want {
		t.Errorf("span-1.ChildSpanCount=%d; want %d", got, want)
	}
	if got, want := gotParent.ChildSpanCount, 2; got != want {
		t.Errorf("parent.ChildSpanCount=%d; want %d", got, want)
	}
}

func TestNilSpanEnd(t *testing.T) {
	var span *span
	span.End()
}

func TestExecutionTracerTaskEnd(t *testing.T) {
	var n uint64
	tp := NewTracerProvider(WithConfig(Config{DefaultSampler: NeverSample()}))
	tr := tp.Tracer("Execution Tracer Task End")

	executionTracerTaskEnd := func() {
		atomic.AddUint64(&n, 1)
	}

	var spans []*span
	_, apiSpan := tr.Start(context.Background(), "foo")
	s := apiSpan.(*span)

	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // never sample

	tID, _ := otel.TraceIDFromHex("0102030405060708090a0b0c0d0e0f")
	sID, _ := otel.SpanIDFromHex("0001020304050607")
	ctx := context.Background()

	ctx = otel.ContextWithRemoteSpanContext(ctx,
		otel.SpanContext{
			TraceID:    tID,
			SpanID:     sID,
			TraceFlags: 0,
		},
	)
	_, apiSpan = tr.Start(
		ctx,
		"foo",
	)
	s = apiSpan.(*span)
	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // parent not sampled

	//tp.ApplyConfig(Config{DefaultSampler: AlwaysSample()})
	_, apiSpan = tr.Start(context.Background(), "foo")
	s = apiSpan.(*span)
	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // always sample

	for _, span := range spans {
		span.End()
	}
	if got, want := n, uint64(len(spans)); got != want {
		t.Fatalf("Execution tracer task ended for %v spans; want %v", got, want)
	}
}

func TestCustomStartEndTime(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithConfig(Config{DefaultSampler: AlwaysSample()}))

	startTime := time.Date(2019, time.August, 27, 14, 42, 0, 0, time.UTC)
	endTime := startTime.Add(time.Second * 20)
	_, span := tp.Tracer("Custom Start and End time").Start(
		context.Background(),
		"testspan",
		otel.WithTimestamp(startTime),
	)
	span.End(otel.WithTimestamp(endTime))

	if te.Len() != 1 {
		t.Fatalf("got %d exported spans, want one span", te.Len())
	}
	got := te.Spans()[0]
	if got.StartTime != startTime {
		t.Errorf("expected start time to be %s, got %s", startTime, got.StartTime)
	}
	if got.EndTime != endTime {
		t.Errorf("expected end time to be %s, got %s", endTime, got.EndTime)
	}
}

func TestRecordError(t *testing.T) {
	scenarios := []struct {
		err error
		typ string
		msg string
	}{
		{
			err: ottest.NewTestError("test error"),
			typ: "go.opentelemetry.io/otel/internal/testing.TestError",
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
		tp := NewTracerProvider(WithSyncer(te))
		span := startSpan(tp, "RecordError")

		errTime := time.Now()
		span.RecordError(context.Background(), s.err,
			otel.WithErrorTime(errTime),
		)

		got, err := endSpan(te, span)
		if err != nil {
			t.Fatal(err)
		}

		want := &export.SpanData{
			SpanContext: otel.SpanContext{
				TraceID:    tid,
				TraceFlags: 0x1,
			},
			ParentSpanID:    sid,
			Name:            "span0",
			SpanKind:        otel.SpanKindInternal,
			HasRemoteParent: true,
			MessageEvents: []export.Event{
				{
					Name: errorEventName,
					Time: errTime,
					Attributes: []label.KeyValue{
						errorTypeKey.String(s.typ),
						errorMessageKey.String(s.msg),
					},
				},
			},
			InstrumentationLibrary: instrumentation.Library{Name: "RecordError"},
		}
		if diff := cmpDiff(got, want); diff != "" {
			t.Errorf("SpanErrorOptions: -got +want %s", diff)
		}
	}
}

func TestRecordErrorWithStatus(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	span := startSpan(tp, "RecordErrorWithStatus")

	testErr := ottest.NewTestError("test error")
	errTime := time.Now()
	testStatus := codes.Error
	span.RecordError(context.Background(), testErr,
		otel.WithErrorTime(errTime),
		otel.WithErrorStatus(testStatus),
	)

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		SpanKind:        otel.SpanKindInternal,
		StatusCode:      codes.Error,
		StatusMessage:   "",
		HasRemoteParent: true,
		MessageEvents: []export.Event{
			{
				Name: errorEventName,
				Time: errTime,
				Attributes: []label.KeyValue{
					errorTypeKey.String("go.opentelemetry.io/otel/internal/testing.TestError"),
					errorMessageKey.String("test error"),
				},
			},
		},
		InstrumentationLibrary: instrumentation.Library{Name: "RecordErrorWithStatus"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SpanErrorOptions: -got +want %s", diff)
	}
}

func TestRecordErrorNil(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	span := startSpan(tp, "RecordErrorNil")

	span.RecordError(context.Background(), nil)

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:           sid,
		Name:                   "span0",
		SpanKind:               otel.SpanKindInternal,
		HasRemoteParent:        true,
		StatusCode:             codes.Unset,
		StatusMessage:          "",
		InstrumentationLibrary: instrumentation.Library{Name: "RecordErrorNil"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SpanErrorOptions: -got +want %s", diff)
	}
}

func TestWithSpanKind(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithConfig(Config{DefaultSampler: AlwaysSample()}))
	tr := tp.Tracer("withSpanKind")

	_, span := tr.Start(context.Background(), "WithoutSpanKind")
	spanData, err := endSpan(te, span)
	if err != nil {
		t.Error(err.Error())
	}

	if spanData.SpanKind != otel.SpanKindInternal {
		t.Errorf("Default value of Spankind should be Internal: got %+v, want %+v\n", spanData.SpanKind, otel.SpanKindInternal)
	}

	sks := []otel.SpanKind{
		otel.SpanKindInternal,
		otel.SpanKindServer,
		otel.SpanKindClient,
		otel.SpanKindProducer,
		otel.SpanKindConsumer,
	}

	for _, sk := range sks {
		te.Reset()

		_, span := tr.Start(context.Background(), fmt.Sprintf("SpanKind-%v", sk), otel.WithSpanKind(sk))
		spanData, err := endSpan(te, span)
		if err != nil {
			t.Error(err.Error())
		}

		if spanData.SpanKind != sk {
			t.Errorf("WithSpanKind check: got %+v, want %+v\n", spanData.SpanKind, sks)
		}
	}
}

func TestWithResource(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te),
		WithConfig(Config{DefaultSampler: AlwaysSample()}),
		WithResource(resource.New(label.String("rk1", "rv1"), label.Int64("rk2", 5))))
	span := startSpan(tp, "WithResource")
	span.SetAttributes(label.String("key1", "value1"))
	got, err := endSpan(te, span)
	if err != nil {
		t.Error(err.Error())
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []label.KeyValue{
			label.String("key1", "value1"),
		},
		SpanKind:               otel.SpanKindInternal,
		HasRemoteParent:        true,
		Resource:               resource.New(label.String("rk1", "rv1"), label.Int64("rk2", 5)),
		InstrumentationLibrary: instrumentation.Library{Name: "WithResource"},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("WithResource:\n  -got +want %s", diff)
	}
}

func TestWithInstrumentationVersion(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))

	ctx := context.Background()
	ctx = otel.ContextWithRemoteSpanContext(ctx, remoteSpanContext())
	_, span := tp.Tracer(
		"WithInstrumentationVersion",
		otel.WithInstrumentationVersion("v0.1.0"),
	).Start(ctx, "span0", otel.WithRecord())
	got, err := endSpan(te, span)
	if err != nil {
		t.Error(err.Error())
	}

	want := &export.SpanData{
		SpanContext: otel.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		SpanKind:        otel.SpanKindInternal,
		HasRemoteParent: true,
		InstrumentationLibrary: instrumentation.Library{
			Name:    "WithInstrumentationVersion",
			Version: "v0.1.0",
		},
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("WithResource:\n  -got +want %s", diff)
	}
}

func TestSpanCapturesPanic(t *testing.T) {
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te))
	_, span := tp.Tracer("CatchPanic").Start(
		context.Background(),
		"span",
		otel.WithRecord(),
	)

	f := func() {
		defer span.End()
		panic(errors.New("error message"))
	}
	require.PanicsWithError(t, "error message", f)
	spans := te.Spans()
	require.Len(t, spans, 1)
	require.Len(t, spans[0].MessageEvents, 1)
	assert.Equal(t, spans[0].MessageEvents[0].Name, errorEventName)
	assert.Equal(t, spans[0].MessageEvents[0].Attributes, []label.KeyValue{
		errorTypeKey.String("*errors.errorString"),
		errorMessageKey.String("error message"),
	})
}
