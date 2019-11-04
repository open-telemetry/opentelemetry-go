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

package trace

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/testharness"
	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/export"
)

var (
	tid core.TraceID
	sid core.SpanID
)

func init() {
	tid, _ = core.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ = core.SpanIDFromHex("0102040810203040")
}

func TestTracerFollowsExpectedAPIBehaviour(t *testing.T) {
	tp, err := NewProvider(WithConfig(Config{DefaultSampler: ProbabilitySampler(0)}))
	if err != nil {
		t.Fatalf("failed to create provider, err: %v\n", err)
	}
	harness := testharness.NewHarness(t)
	subjectFactory := func() trace.Tracer {
		return tp.GetTracer("")
	}

	harness.TestTracer(subjectFactory)
}

type testExporter struct {
	spans []*export.SpanData
}

func (t *testExporter) ExportSpan(ctx context.Context, d *export.SpanData) {
	t.spans = append(t.spans, d)
}

func TestSetName(t *testing.T) {
	samplerIsCalled := false
	fooSampler := Sampler(func(p SamplingParameters) SamplingDecision {
		samplerIsCalled = true
		t.Logf("called sampler for name %q", p.Name)
		return SamplingDecision{Sample: strings.HasPrefix(p.Name, "SetName/foo")}
	})
	tp, _ := NewProvider(WithConfig(Config{DefaultSampler: fooSampler}))

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
		if !samplerIsCalled {
			t.Errorf("%d: the sampler was not even called during span creation", idx)
		}
		samplerIsCalled = false
		if gotSampledBefore := span.SpanContext().IsSampled(); tt.sampledBefore != gotSampledBefore {
			t.Errorf("%d: invalid sampling decision before rename, expected %v, got %v", idx, tt.sampledBefore, gotSampledBefore)
		}
		span.SetName(tt.newName)
		if !samplerIsCalled {
			t.Errorf("%d: the sampler was not even called during span rename", idx)
		}
		samplerIsCalled = false
		if gotSampledAfter := span.SpanContext().IsSampled(); tt.sampledAfter != gotSampledAfter {
			t.Errorf("%d: invalid sampling decision after rename, expected %v, got %v", idx, tt.sampledAfter, gotSampledAfter)
		}
		span.End()
	}
}

func TestRecordingIsOff(t *testing.T) {
	tp, _ := NewProvider()
	_, span := tp.GetTracer("Recording off").Start(context.Background(), "StartSpan")
	defer span.End()
	if span.IsRecording() == true {
		t.Error("new span is recording events")
	}
}

func TestSampling(t *testing.T) {
	idg := defIDGenerator()
	total := 10000
	for name, tc := range map[string]struct {
		sampler       Sampler
		expect        float64
		tolerance     float64
		parent        bool
		sampledParent bool
	}{
		// Span w/o a parent
		"NeverSample":            {sampler: NeverSample(), expect: 0, tolerance: 0},
		"AlwaysSample":           {sampler: AlwaysSample(), expect: 1.0, tolerance: 0},
		"ProbabilitySampler_-1":  {sampler: ProbabilitySampler(-1.0), expect: 0, tolerance: 0},
		"ProbabilitySampler_.25": {sampler: ProbabilitySampler(0.25), expect: .25, tolerance: 0.015},
		"ProbabilitySampler_.50": {sampler: ProbabilitySampler(0.50), expect: .5, tolerance: 0.015},
		"ProbabilitySampler_.75": {sampler: ProbabilitySampler(0.75), expect: .75, tolerance: 0.015},
		"ProbabilitySampler_2.0": {sampler: ProbabilitySampler(2.0), expect: 1, tolerance: 0},
		// Spans with a parent that is *not* sampled act like spans w/o a parent
		"UnsampledParentSpanWithProbabilitySampler_-1":  {sampler: ProbabilitySampler(-1.0), expect: 0, tolerance: 0, parent: true},
		"UnsampledParentSpanWithProbabilitySampler_.25": {sampler: ProbabilitySampler(.25), expect: .25, tolerance: 0.015, parent: true},
		"UnsampledParentSpanWithProbabilitySampler_.50": {sampler: ProbabilitySampler(0.50), expect: .5, tolerance: 0.015, parent: true},
		"UnsampledParentSpanWithProbabilitySampler_.75": {sampler: ProbabilitySampler(0.75), expect: .75, tolerance: 0.015, parent: true},
		"UnsampledParentSpanWithProbabilitySampler_2.0": {sampler: ProbabilitySampler(2.0), expect: 1, tolerance: 0, parent: true},
		// Spans with a parent that is sampled, will always sample, regardless of the probability
		"SampledParentSpanWithProbabilitySampler_-1":  {sampler: ProbabilitySampler(-1.0), expect: 1, tolerance: 0, parent: true, sampledParent: true},
		"SampledParentSpanWithProbabilitySampler_.25": {sampler: ProbabilitySampler(.25), expect: 1, tolerance: 0, parent: true, sampledParent: true},
		"SampledParentSpanWithProbabilitySampler_2.0": {sampler: ProbabilitySampler(2.0), expect: 1, tolerance: 0, parent: true, sampledParent: true},
		// Spans with a sampled parent, but when using the NeverSample Sampler, aren't sampled
		"SampledParentSpanWithNeverSample": {sampler: NeverSample(), expect: 0, tolerance: 0, parent: true, sampledParent: true},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p, err := NewProvider(WithConfig(Config{DefaultSampler: tc.sampler}))
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			tr := p.GetTracer("test")
			var sampled int
			for i := 0; i < total; i++ {
				var opts []apitrace.SpanOption
				if tc.parent {
					psc := core.SpanContext{
						TraceID: idg.NewTraceID(),
						SpanID:  idg.NewSpanID(),
					}
					if tc.sampledParent {
						psc.TraceFlags = core.TraceFlagsSampled
					}
					opts = append(opts, apitrace.ChildOf(psc))
				}
				_, span := tr.Start(context.Background(), "test", opts...)
				if span.SpanContext().IsSampled() {
					sampled++
				}
			}
			got := float64(sampled) / float64(total)
			diff := math.Abs(got - tc.expect)
			if diff > tc.tolerance {
				t.Errorf("got %f (diff: %f), expected %f (w/tolerance: %f)", got, diff, tc.expect, tc.tolerance)
			}
		})
	}
}

func TestStartSpanWithChildOf(t *testing.T) {
	tp, _ := NewProvider()
	tr := tp.GetTracer("SpanWith ChildOf")

	sc1 := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x0,
	}
	_, s1 := tr.Start(context.Background(), "span1-unsampled-parent1", apitrace.ChildOf(sc1))
	if err := checkChild(sc1, s1); err != nil {
		t.Error(err)
	}

	_, s2 := tr.Start(context.Background(), "span2-unsampled-parent1", apitrace.ChildOf(sc1))
	if err := checkChild(sc1, s2); err != nil {
		t.Error(err)
	}

	sc2 := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
		//Tracestate:   testTracestate,
	}
	_, s3 := tr.Start(context.Background(), "span3-sampled-parent2", apitrace.ChildOf(sc2))
	if err := checkChild(sc2, s3); err != nil {
		t.Error(err)
	}

	ctx, s4 := tr.Start(context.Background(), "span4-sampled-parent2", apitrace.ChildOf(sc2))
	if err := checkChild(sc2, s4); err != nil {
		t.Error(err)
	}

	s4Sc := s4.SpanContext()
	_, s5 := tr.Start(ctx, "span5-implicit-childof-span4")
	if err := checkChild(s4Sc, s5); err != nil {
		t.Error(err)
	}
}

func TestStartSpanWithFollowsFrom(t *testing.T) {
	tp, _ := NewProvider()
	tr := tp.GetTracer("SpanWith FollowsFrom")

	sc1 := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x0,
	}
	_, s1 := tr.Start(context.Background(), "span1-unsampled-parent1", apitrace.FollowsFrom(sc1))
	if err := checkChild(sc1, s1); err != nil {
		t.Error(err)
	}

	_, s2 := tr.Start(context.Background(), "span2-unsampled-parent1", apitrace.FollowsFrom(sc1))
	if err := checkChild(sc1, s2); err != nil {
		t.Error(err)
	}

	sc2 := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 0x1,
		//Tracestate:   testTracestate,
	}
	_, s3 := tr.Start(context.Background(), "span3-sampled-parent2", apitrace.FollowsFrom(sc2))
	if err := checkChild(sc2, s3); err != nil {
		t.Error(err)
	}

	ctx, s4 := tr.Start(context.Background(), "span4-sampled-parent2", apitrace.FollowsFrom(sc2))
	if err := checkChild(sc2, s4); err != nil {
		t.Error(err)
	}

	s4Sc := s4.SpanContext()
	_, s5 := tr.Start(ctx, "span5-implicit-childof-span4")
	if err := checkChild(s4Sc, s5); err != nil {
		t.Error(err)
	}
}

// TODO: [rghetia] Equivalent of SpanKind Test.

func TestSetSpanAttributesOnStart(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))
	span := startSpan(tp, "StartSpanAttribute", apitrace.WithAttributes(key.String("key1", "value1")))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "StartSpanAttribute/span0",
		Attributes: []core.KeyValue{
			key.String("key1", "value1"),
		},
		SpanKind:        "internal",
		HasRemoteParent: true,
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOnStart: -got +want %s", diff)
	}
}

func TestSetSpanAttributes(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))
	span := startSpan(tp, "SpanAttribute")
	span.SetAttribute(key.New("key1").String("value1"))
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "SpanAttribute/span0",
		Attributes: []core.KeyValue{
			key.String("key1", "value1"),
		},
		SpanKind:        "internal",
		HasRemoteParent: true,
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributes: -got +want %s", diff)
	}
}

func TestSetSpanAttributesOverLimit(t *testing.T) {
	te := &testExporter{}
	cfg := Config{MaxAttributesPerSpan: 2}
	tp, _ := NewProvider(WithConfig(cfg), WithSyncer(te))

	span := startSpan(tp, "SpanAttributesOverLimit")
	span.SetAttribute(key.Bool("key1", true))
	span.SetAttribute(key.String("key2", "value2"))
	span.SetAttribute(key.Bool("key1", false)) // Replace key1.
	span.SetAttribute(key.Int64("key4", 4))    // Remove key2 and add key4
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "SpanAttributesOverLimit/span0",
		Attributes: []core.KeyValue{
			key.Bool("key1", false),
			key.Int64("key4", 4),
		},
		SpanKind:              "internal",
		HasRemoteParent:       true,
		DroppedAttributeCount: 1,
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOverLimit: -got +want %s", diff)
	}
}

func TestEvents(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))

	span := startSpan(tp, "Events")
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.Bool("key2", true)
	k3v3 := key.Int64("key3", 3)

	span.AddEvent(context.Background(), "foo", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "bar",
		key.Bool("key2", true),
		key.Int64("key3", 3),
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
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "Events/span0",
		HasRemoteParent: true,
		MessageEvents: []export.Event{
			{Message: "foo", Attributes: []core.KeyValue{k1v1}},
			{Message: "bar", Attributes: []core.KeyValue{k2v2, k3v3}},
		},
		SpanKind: "internal",
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Events: -got +want %s", diff)
	}
}

func TestEventsOverLimit(t *testing.T) {
	te := &testExporter{}
	cfg := Config{MaxEventsPerSpan: 2}
	tp, _ := NewProvider(WithConfig(cfg), WithSyncer(te))

	span := startSpan(tp, "EventsOverLimit")
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.Bool("key2", false)
	k3v3 := key.New("key3").String("value3")

	span.AddEvent(context.Background(), "fooDrop", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "barDrop",
		key.Bool("key2", true),
		key.New("key3").String("value3"),
	)
	span.AddEvent(context.Background(), "foo", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "bar",
		key.Bool("key2", false),
		key.New("key3").String("value3"),
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
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "EventsOverLimit/span0",
		MessageEvents: []export.Event{
			{Message: "foo", Attributes: []core.KeyValue{k1v1}},
			{Message: "bar", Attributes: []core.KeyValue{k2v2, k3v3}},
		},
		DroppedMessageEventCount: 2,
		HasRemoteParent:          true,
		SpanKind:                 "internal",
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Message Event over limit: -got +want %s", diff)
	}
}

func TestAddLinks(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))

	span := startSpan(tp, "AddLinks")
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")

	sc1 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}
	sc2 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}

	link1 := apitrace.Link{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}}
	link2 := apitrace.Link{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}}
	span.AddLink(link1)
	span.AddLink(link2)

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "AddLinks/span0",
		HasRemoteParent: true,
		Links: []apitrace.Link{
			{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}},
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}},
		},
		SpanKind: "internal",
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("AddLink: -got +want %s", diff)
	}
}

func TestLinks(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))

	span := startSpan(tp, "Links")
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	sc1 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}
	sc2 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}

	span.Link(sc1, key.New("key1").String("value1"))
	span.Link(sc2,
		key.New("key2").String("value2"),
		key.New("key3").String("value3"),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "Links/span0",
		HasRemoteParent: true,
		Links: []apitrace.Link{
			{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}},
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2, k3v3}},
		},
		SpanKind: "internal",
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link: -got +want %s", diff)
	}
}

func TestLinksOverLimit(t *testing.T) {
	te := &testExporter{}
	cfg := Config{MaxLinksPerSpan: 2}

	sc1 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}
	sc2 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}
	sc3 := core.SpanContext{TraceID: core.TraceID([16]byte{1, 1}), SpanID: core.SpanID{3}}

	tp, _ := NewProvider(WithConfig(cfg), WithSyncer(te))
	span := startSpan(tp, "LinksOverLimit")

	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	span.Link(sc1, key.New("key1").String("value1"))
	span.Link(sc2, key.New("key2").String("value2"))
	span.Link(sc3, key.New("key3").String("value3"))

	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID: sid,
		Name:         "LinksOverLimit/span0",
		Links: []apitrace.Link{
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}},
			{SpanContext: sc3, Attributes: []core.KeyValue{k3v3}},
		},
		DroppedLinkCount: 1,
		HasRemoteParent:  true,
		SpanKind:         "internal",
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("Link over limit: -got +want %s", diff)
	}
}

func TestSetSpanName(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))

	want := "SetSpanName/SpanName-1"
	_, span := tp.GetTracer("SetSpanName").Start(context.Background(), "SpanName-1",
		apitrace.ChildOf(core.SpanContext{
			TraceID:    tid,
			SpanID:     sid,
			TraceFlags: 1,
		}),
	)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != want {
		t.Errorf("span.Name: got %q; want %q", got.Name, want)
	}
}

func TestSetSpanStatus(t *testing.T) {
	te := &testExporter{}
	tp, _ := NewProvider(WithSyncer(te))

	span := startSpan(tp, "SpanStatus")
	span.SetStatus(codes.Canceled)
	got, err := endSpan(te, span)
	if err != nil {
		t.Fatal(err)
	}

	want := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID:    tid,
			TraceFlags: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "SpanStatus/span0",
		SpanKind:        "internal",
		Status:          codes.Canceled,
		HasRemoteParent: true,
	}
	if diff := cmpDiff(got, want); diff != "" {
		t.Errorf("SetSpanStatus: -got +want %s", diff)
	}
}

func cmpDiff(x, y interface{}) string {
	return cmp.Diff(x, y, cmp.AllowUnexported(core.Value{}), cmp.AllowUnexported(export.Event{}))
}

func remoteSpanContext() core.SpanContext {
	return core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: 1,
	}
}

// checkChild is test utility function that tests that c has fields set appropriately,
// given that it is a child span of p.
func checkChild(p core.SpanContext, apiSpan apitrace.Span) error {
	s := apiSpan.(*span)
	if s == nil {
		return fmt.Errorf("got nil child span, want non-nil")
	}
	if got, want := s.spanContext.TraceIDString(), p.TraceIDString(); got != want {
		return fmt.Errorf("got child trace ID %s, want %s", got, want)
	}
	if childID, parentID := s.spanContext.SpanIDString(), p.SpanIDString(); childID == parentID {
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
func startSpan(tp *Provider, trName string, args ...apitrace.SpanOption) apitrace.Span {
	return startNamedSpan(tp, trName, "span0", args...)
}

// startNamed Span is a test utility func that starts a span with a
// passed name and with ChildOf option.  remote span context contains
// TraceFlags with sampled bit set. This allows the span to be
// automatically sampled.
func startNamedSpan(tp *Provider, trName, name string, args ...apitrace.SpanOption) apitrace.Span {
	args = append(args, apitrace.ChildOf(remoteSpanContext()), apitrace.WithRecord())
	_, span := tp.GetTracer(trName).Start(
		context.Background(),
		name,
		args...,
	)
	return span
}

// endSpan is a test utility function that ends the span in the context and
// returns the exported export.SpanData.
// It requires that span be sampled using one of these methods
//  1. Passing parent span context using ChildOf option
//  2. Use WithSampler(AlwaysSample())
//  3. Configuring AlwaysSample() as default sampler
//
// It also does some basic tests on the span.
// It also clears spanID in the export.SpanData to make the comparison easier.
func endSpan(te *testExporter, span apitrace.Span) (*export.SpanData, error) {

	if !span.IsRecording() {
		return nil, fmt.Errorf("IsRecording: got false, want true")
	}
	if !span.SpanContext().IsSampled() {
		return nil, fmt.Errorf("IsSampled: got false, want true")
	}
	span.End()
	if len(te.spans) != 1 {
		return nil, fmt.Errorf("got exported spans %#v, want one span", te.spans)
	}
	got := te.spans[0]
	if !got.SpanContext.SpanID.IsValid() {
		return nil, fmt.Errorf("exporting span: expected nonzero SpanID")
	}
	got.SpanContext.SpanID = core.SpanID{}
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

type fakeExporter map[string]*export.SpanData

func (f fakeExporter) ExportSpan(ctx context.Context, s *export.SpanData) {
	f[s.Name] = s
}

func TestEndSpanTwice(t *testing.T) {
	spans := make(fakeExporter)
	tp, _ := NewProvider(WithSyncer(spans))

	span := startSpan(tp, "EndSpanTwice")
	span.End()
	span.End()
	if len(spans) != 1 {
		t.Fatalf("expected only a single span, got %#v", spans)
	}
}

func TestStartSpanAfterEnd(t *testing.T) {
	spans := make(fakeExporter)
	tp, _ := NewProvider(WithConfig(Config{DefaultSampler: AlwaysSample()}), WithSyncer(spans))

	tr := tp.GetTracer("SpanAfterEnd")
	ctx, span0 := tr.Start(context.Background(), "parent", apitrace.ChildOf(remoteSpanContext()))
	ctx1, span1 := tr.Start(ctx, "span-1")
	span1.End()
	// Start a new span with the context containing span-1
	// even though span-1 is ended, we still add this as a new child of span-1
	_, span2 := tr.Start(ctx1, "span-2")
	span2.End()
	span0.End()
	if got, want := len(spans), 3; got != want {
		t.Fatalf("len(%#v) = %d; want %d", spans, got, want)
	}
	if got, want := spans["SpanAfterEnd/span-1"].SpanContext.TraceID, spans["SpanAfterEnd/parent"].SpanContext.TraceID; got != want {
		t.Errorf("span-1.TraceID=%q; want %q", got, want)
	}
	if got, want := spans["SpanAfterEnd/span-2"].SpanContext.TraceID, spans["SpanAfterEnd/parent"].SpanContext.TraceID; got != want {
		t.Errorf("span-2.TraceID=%q; want %q", got, want)
	}
	if got, want := spans["SpanAfterEnd/span-1"].ParentSpanID, spans["SpanAfterEnd/parent"].SpanContext.SpanID; got != want {
		t.Errorf("span-1.ParentSpanID=%q; want %q (parent.SpanID)", got, want)
	}
	if got, want := spans["SpanAfterEnd/span-2"].ParentSpanID, spans["SpanAfterEnd/span-1"].SpanContext.SpanID; got != want {
		t.Errorf("span-2.ParentSpanID=%q; want %q (span1.SpanID)", got, want)
	}
}

func TestChildSpanCount(t *testing.T) {
	spans := make(fakeExporter)
	tp, _ := NewProvider(WithConfig(Config{DefaultSampler: AlwaysSample()}), WithSyncer(spans))

	tr := tp.GetTracer("ChidSpanCount")
	ctx, span0 := tr.Start(context.Background(), "parent")
	ctx1, span1 := tr.Start(ctx, "span-1")
	_, span2 := tr.Start(ctx1, "span-2")
	span2.End()
	span1.End()

	_, span3 := tr.Start(ctx, "span-3")
	span3.End()
	span0.End()
	if got, want := len(spans), 4; got != want {
		t.Fatalf("len(%#v) = %d; want %d", spans, got, want)
	}
	if got, want := spans["ChidSpanCount/span-3"].ChildSpanCount, 0; got != want {
		t.Errorf("span-3.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["ChidSpanCount/span-2"].ChildSpanCount, 0; got != want {
		t.Errorf("span-2.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["ChidSpanCount/span-1"].ChildSpanCount, 1; got != want {
		t.Errorf("span-1.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["ChidSpanCount/parent"].ChildSpanCount, 2; got != want {
		t.Errorf("parent.ChildSpanCount=%q; want %q", got, want)
	}
}

func TestNilSpanEnd(t *testing.T) {
	var span *span
	span.End()
}

func TestExecutionTracerTaskEnd(t *testing.T) {
	var n uint64
	tp, _ := NewProvider(WithConfig(Config{DefaultSampler: NeverSample()}))
	tr := tp.GetTracer("Execution Tracer Task End")

	executionTracerTaskEnd := func() {
		atomic.AddUint64(&n, 1)
	}

	var spans []*span
	_, apiSpan := tr.Start(context.Background(), "foo")
	s := apiSpan.(*span)

	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // never sample

	tID, _ := core.TraceIDFromHex("0102030405060708090a0b0c0d0e0f")
	sID, _ := core.SpanIDFromHex("0001020304050607")

	_, apiSpan = tr.Start(
		context.Background(),
		"foo",
		apitrace.ChildOf(
			core.SpanContext{
				TraceID:    tID,
				SpanID:     sID,
				TraceFlags: 0,
			},
		),
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
	var te testExporter
	tp, _ := NewProvider(WithSyncer(&te), WithConfig(Config{DefaultSampler: AlwaysSample()}))

	startTime := time.Date(2019, time.August, 27, 14, 42, 0, 0, time.UTC)
	endTime := startTime.Add(time.Second * 20)
	_, span := tp.GetTracer("Custom Start and End time").Start(
		context.Background(),
		"testspan",
		apitrace.WithStartTime(startTime),
	)
	span.End(apitrace.WithEndTime(endTime))

	if len(te.spans) != 1 {
		t.Fatalf("got exported spans %#v, want one span", te.spans)
	}
	got := te.spans[0]
	if got.StartTime != startTime {
		t.Errorf("expected start time to be %s, got %s", startTime, got.StartTime)
	}
	if got.EndTime != endTime {
		t.Errorf("expected end time to be %s, got %s", endTime, got.EndTime)
	}
}

func TestWithSpanKind(t *testing.T) {
	var te testExporter
	tp, _ := NewProvider(WithSyncer(&te), WithConfig(Config{DefaultSampler: AlwaysSample()}))
	tr := tp.GetTracer("withSpanKind")

	_, span := tr.Start(context.Background(), "WithoutSpanKind")
	spanData, err := endSpan(&te, span)
	if err != nil {
		t.Error(err.Error())
	}

	if spanData.SpanKind != apitrace.SpanKindInternal {
		t.Errorf("Default value of Spankind should be Internal: got %+v, want %+v\n", spanData.SpanKind, apitrace.SpanKindInternal)
	}

	sks := []apitrace.SpanKind{
		apitrace.SpanKindInternal,
		apitrace.SpanKindServer,
		apitrace.SpanKindClient,
		apitrace.SpanKindProducer,
		apitrace.SpanKindConsumer,
	}

	for _, sk := range sks {
		te.spans = nil

		_, span := tr.Start(context.Background(), fmt.Sprintf("SpanKind-%v", sk), apitrace.WithSpanKind(sk))
		spanData, err := endSpan(&te, span)
		if err != nil {
			t.Error(err.Error())
		}

		if spanData.SpanKind != sk {
			t.Errorf("WithSpanKind check: got %+v, want %+v\n", spanData.SpanKind, sks)
		}
	}
}
