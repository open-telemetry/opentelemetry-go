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
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	apitrace "go.opentelemetry.io/api/trace"
)

var (
	tid = core.TraceID{High: 0x0102030405060708, Low: 0x0102040810203040}
	sid = uint64(0x0102040810203040)
)

func init() {
	Register()
	setupDefaultSamplerConfig()
}

func setupDefaultSamplerConfig() {
	// no random sampling, but sample children of sampled spans.
	ApplyConfig(Config{DefaultSampler: ProbabilitySampler(0)})
}

type testExporter struct {
	spans []*SpanData
}

func (t *testExporter) ExportSpan(s *SpanData) {
	t.spans = append(t.spans, s)
}

func TestStartSpan(t *testing.T) {
	_, span := apitrace.GlobalTracer().Start(context.Background(), "StartSpan")
	defer span.Finish()
	if span == nil {
		t.Errorf("span not started")
	}
}

func TestSetName(t *testing.T) {
	samplerIsCalled := false
	fooSampler := Sampler(func(p SamplingParameters) SamplingDecision {
		samplerIsCalled = true
		t.Logf("called sampler for name %q", p.Name)
		return SamplingDecision{Sample: strings.HasPrefix(p.Name, "foo")}
	})
	ApplyConfig(Config{DefaultSampler: fooSampler})
	defer setupDefaultSamplerConfig()
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
		span := startNamedSpan(tt.name)
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
		span.Finish()
	}
}

func TestRecordingIsOff(t *testing.T) {
	_, span := apitrace.GlobalTracer().Start(context.Background(), "StartSpan")
	defer span.Finish()
	if span.IsRecordingEvents() == true {
		t.Error("new span is recording events")
	}
}

// TODO: [rghetia] enable sampling test when Sampling is working.

func TestStartSpanWithChildOf(t *testing.T) {
	sc1 := core.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: 0x0,
	}
	_, s1 := apitrace.GlobalTracer().Start(context.Background(), "span1-unsampled-parent1", apitrace.ChildOf(sc1))
	if err := checkChild(sc1, s1); err != nil {
		t.Error(err)
	}

	_, s2 := apitrace.GlobalTracer().Start(context.Background(), "span2-unsampled-parent1", apitrace.ChildOf(sc1))
	if err := checkChild(sc1, s2); err != nil {
		t.Error(err)
	}

	sc2 := core.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: 0x1,
		//Tracestate:   testTracestate,
	}
	_, s3 := apitrace.GlobalTracer().Start(context.Background(), "span3-sampled-parent2", apitrace.ChildOf(sc2))
	if err := checkChild(sc2, s3); err != nil {
		t.Error(err)
	}

	ctx, s4 := apitrace.GlobalTracer().Start(context.Background(), "span4-sampled-parent2", apitrace.ChildOf(sc2))
	if err := checkChild(sc2, s4); err != nil {
		t.Error(err)
	}

	s4Sc := s4.SpanContext()
	_, s5 := apitrace.GlobalTracer().Start(ctx, "span5-implicit-childof-span4")
	if err := checkChild(s4Sc, s5); err != nil {
		t.Error(err)
	}
}

// TODO: [rghetia] Equivalent of SpanKind Test.

func TestSetSpanAttributes(t *testing.T) {
	span := startSpan()
	span.SetAttribute(key.New("key1").String("value1"))
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []core.KeyValue{{
			Key:   core.Key{Name: "key1"},
			Value: core.Value{Type: core.STRING, String: "value1"},
		}},
		HasRemoteParent: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("SetSpanAttributes: -got +want %s", diff)
	}
}

func TestSetSpanAttributesOverLimit(t *testing.T) {
	cfg := Config{MaxAttributesPerSpan: 2}
	ApplyConfig(cfg)

	span := startSpan()
	span.SetAttribute(key.New("key1").String("value1"))
	span.SetAttribute(key.New("key2").String("value2"))
	span.SetAttribute(key.New("key1").String("value3")) // Replace key1.
	span.SetAttribute(key.New("key4").String("value4")) // Remove key2 and add key4
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Attributes: []core.KeyValue{
			{
				Key:   core.Key{Name: "key1"},
				Value: core.Value{Type: core.STRING, String: "value3"},
			},
			{
				Key:   core.Key{Name: "key4"},
				Value: core.Value{Type: core.STRING, String: "value4"},
			},
		},
		HasRemoteParent:       true,
		DroppedAttributeCount: 1,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("SetSpanAttributesOverLimit: -got +want %s", diff)
	}
}

func TestEvents(t *testing.T) {
	span := startSpan()
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	span.AddEvent(context.Background(), "foo", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "bar",
		key.New("key2").String("value2"),
		key.New("key3").String("value3"),
	)
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.MessageEvents {
		if !checkTime(&got.MessageEvents[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		HasRemoteParent: true,
		MessageEvents: []Event{
			{Message: "foo", Attributes: []core.KeyValue{k1v1}},
			{Message: "bar", Attributes: []core.KeyValue{k2v2, k3v3}},
		},
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Event{})); diff != "" {
		t.Errorf("Message Events: -got +want %s", diff)
	}
}

func TestEventsOverLimit(t *testing.T) {
	cfg := Config{MaxEventsPerSpan: 2}
	ApplyConfig(cfg)
	span := startSpan()
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	span.AddEvent(context.Background(), "fooDrop", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "barDrop",
		key.New("key2").String("value2"),
		key.New("key3").String("value3"),
	)
	span.AddEvent(context.Background(), "foo", key.New("key1").String("value1"))
	span.AddEvent(context.Background(), "bar",
		key.New("key2").String("value2"),
		key.New("key3").String("value3"),
	)
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	for i := range got.MessageEvents {
		if !checkTime(&got.MessageEvents[i].Time) {
			t.Error("exporting span: expected nonzero Event Time")
		}
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		MessageEvents: []Event{
			{Message: "foo", Attributes: []core.KeyValue{k1v1}},
			{Message: "bar", Attributes: []core.KeyValue{k2v2, k3v3}},
		},
		DroppedMessageEventCount: 2,
		HasRemoteParent:          true,
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Event{})); diff != "" {
		t.Errorf("Message Event over limit: -got +want %s", diff)
	}
}

func TestAddLinks(t *testing.T) {
	span := startSpan()
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")

	sc1 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x1}, SpanID: 0x3}
	sc2 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x2}, SpanID: 0x3}

	link1 := apitrace.Link{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}}
	link2 := apitrace.Link{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}}
	span.AddLink(link1)
	span.AddLink(link2)

	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		HasRemoteParent: true,
		Links: []apitrace.Link{
			{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}},
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}},
		},
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Event{})); diff != "" {
		t.Errorf("AddLink: -got +want %s", diff)
	}
}

func TestLinks(t *testing.T) {
	span := startSpan()
	k1v1 := key.New("key1").String("value1")
	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	sc1 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x1}, SpanID: 0x3}
	sc2 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x2}, SpanID: 0x3}

	span.Link(sc1, key.New("key1").String("value1"))
	span.Link(sc2,
		key.New("key2").String("value2"),
		key.New("key3").String("value3"),
	)
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		HasRemoteParent: true,
		Links: []apitrace.Link{
			{SpanContext: sc1, Attributes: []core.KeyValue{k1v1}},
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2, k3v3}},
		},
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Event{})); diff != "" {
		t.Errorf("Link: -got +want %s", diff)
	}
}

func TestLinksOverLimit(t *testing.T) {
	cfg := Config{MaxLinksPerSpan: 2}
	ApplyConfig(cfg)
	sc1 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x1}, SpanID: 0x3}
	sc2 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x2}, SpanID: 0x3}
	sc3 := core.SpanContext{TraceID: core.TraceID{High: 0x1, Low: 0x3}, SpanID: 0x3}

	span := startSpan()
	k2v2 := key.New("key2").String("value2")
	k3v3 := key.New("key3").String("value3")

	span.Link(sc1, key.New("key1").String("value1"))
	span.Link(sc2, key.New("key2").String("value2"))
	span.Link(sc3, key.New("key3").String("value3"))

	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID: sid,
		Name:         "span0",
		Links: []apitrace.Link{
			{SpanContext: sc2, Attributes: []core.KeyValue{k2v2}},
			{SpanContext: sc3, Attributes: []core.KeyValue{k3v3}},
		},
		DroppedLinkCount: 1,
		HasRemoteParent:  true,
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(Event{})); diff != "" {
		t.Errorf("Link over limit: -got +want %s", diff)
	}
}

func TestSetSpanName(t *testing.T) {
	want := "SpanName-1"
	_, span := apitrace.GlobalTracer().Start(context.Background(), want,
		apitrace.ChildOf(core.SpanContext{
			TraceID:      tid,
			SpanID:       sid,
			TraceOptions: 1,
		}),
	)
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != want {
		t.Errorf("span.Name: got %q; want %q", got.Name, want)
	}
}

func TestSetSpanStatus(t *testing.T) {
	span := startSpan()
	span.SetStatus(codes.Canceled)
	got, err := endSpan(span)
	if err != nil {
		t.Fatal(err)
	}

	want := &SpanData{
		SpanContext: core.SpanContext{
			TraceID:      tid,
			TraceOptions: 0x1,
		},
		ParentSpanID:    sid,
		Name:            "span0",
		Status:          codes.Canceled,
		HasRemoteParent: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("SetSpanStatus: -got +want %s", diff)
	}
}

func TestUnregisterExporter(t *testing.T) {
	var te testExporter
	RegisterExporter(&te)
	UnregisterExporter(&te)

	ctx := startSpan()
	_, _ = endSpan(ctx)
	if len(te.spans) != 0 {
		t.Error("unregistered Exporter was called")
	}
}

func remoteSpanContext() core.SpanContext {
	return core.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: 1,
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
	if got, want := s.spanContext.TraceOptions, p.TraceOptions; got != want {
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
func startSpan() apitrace.Span {
	return startNamedSpan("span0")
}

// startNamed Span is a test utility func that starts a span with a
// passed name and with ChildOf option.  remote span context contains
// traceoption with sampled bit set. This allows the span to be
// automatically sampled.
func startNamedSpan(name string) apitrace.Span {
	_, span := apitrace.GlobalTracer().Start(
		context.Background(),
		name,
		apitrace.ChildOf(remoteSpanContext()),
		apitrace.WithRecordEvents(),
	)
	return span
}

// endSpan is a test utility function that ends the span in the context and
// returns the exported SpanData.
// It requires that span be sampled using one of these methods
//  1. Passing parent span context using ChildOf option
//  2. Use WithSampler(AlwaysSample())
//  3. Configuring AlwaysSample() as default sampler
//
// It also does some basic tests on the span.
// It also clears spanID in the SpanData to make the comparison easier.
func endSpan(span apitrace.Span) (*SpanData, error) {

	if !span.IsRecordingEvents() {
		return nil, fmt.Errorf("IsRecordingEvents: got false, want true")
	}
	if !span.SpanContext().IsSampled() {
		return nil, fmt.Errorf("IsSampled: got false, want true")
	}
	var te testExporter
	RegisterExporter(&te)
	span.Finish()
	UnregisterExporter(&te)
	if len(te.spans) != 1 {
		return nil, fmt.Errorf("got exported spans %#v, want one span", te.spans)
	}
	got := te.spans[0]
	if got.SpanContext.SpanID == 0 {
		return nil, fmt.Errorf("exporting span: expected nonzero SpanID")
	}
	got.SpanContext.SpanID = 0
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

type exporter map[string]*SpanData

func (e exporter) ExportSpan(s *SpanData) {
	e[s.Name] = s
}

func TestEndSpanTwice(t *testing.T) {
	spans := make(exporter)
	RegisterExporter(&spans)
	defer UnregisterExporter(&spans)
	span := startSpan()
	span.Finish()
	span.Finish()
	UnregisterExporter(&spans)
	if len(spans) != 1 {
		t.Fatalf("expected only a single span, got %#v", spans)
	}
}

func TestStartSpanAfterEnd(t *testing.T) {
	spans := make(exporter)
	RegisterExporter(&spans)
	defer UnregisterExporter(&spans)
	ctx, span0 := apitrace.GlobalTracer().Start(context.Background(), "parent", apitrace.ChildOf(remoteSpanContext()))
	ctx1, span1 := apitrace.GlobalTracer().Start(ctx, "span-1")
	span1.Finish()
	// Start a new span with the context containing span-1
	// even though span-1 is ended, we still add this as a new child of span-1
	_, span2 := apitrace.GlobalTracer().Start(ctx1, "span-2")
	span2.Finish()
	span0.Finish()
	UnregisterExporter(&spans)
	if got, want := len(spans), 3; got != want {
		t.Fatalf("len(%#v) = %d; want %d", spans, got, want)
	}
	if got, want := spans["span-1"].SpanContext.TraceID, spans["parent"].SpanContext.TraceID; got != want {
		t.Errorf("span-1.TraceID=%q; want %q", got, want)
	}
	if got, want := spans["span-2"].SpanContext.TraceID, spans["parent"].SpanContext.TraceID; got != want {
		t.Errorf("span-2.TraceID=%q; want %q", got, want)
	}
	if got, want := spans["span-1"].ParentSpanID, spans["parent"].SpanContext.SpanID; got != want {
		t.Errorf("span-1.ParentSpanID=%q; want %q (parent.SpanID)", got, want)
	}
	if got, want := spans["span-2"].ParentSpanID, spans["span-1"].SpanContext.SpanID; got != want {
		t.Errorf("span-2.ParentSpanID=%q; want %q (span1.SpanID)", got, want)
	}
}

func TestChildSpanCount(t *testing.T) {
	ApplyConfig(Config{DefaultSampler: AlwaysSample()})
	spans := make(exporter)
	RegisterExporter(&spans)
	defer UnregisterExporter(&spans)
	ctx, span0 := apitrace.GlobalTracer().Start(context.Background(), "parent")
	ctx1, span1 := apitrace.GlobalTracer().Start(ctx, "span-1")
	_, span2 := apitrace.GlobalTracer().Start(ctx1, "span-2")
	span2.Finish()
	span1.Finish()

	_, span3 := apitrace.GlobalTracer().Start(ctx, "span-3")
	span3.Finish()
	span0.Finish()
	UnregisterExporter(&spans)
	if got, want := len(spans), 4; got != want {
		t.Fatalf("len(%#v) = %d; want %d", spans, got, want)
	}
	if got, want := spans["span-3"].ChildSpanCount, 0; got != want {
		t.Errorf("span-3.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["span-2"].ChildSpanCount, 0; got != want {
		t.Errorf("span-2.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["span-1"].ChildSpanCount, 1; got != want {
		t.Errorf("span-1.ChildSpanCount=%q; want %q", got, want)
	}
	if got, want := spans["parent"].ChildSpanCount, 2; got != want {
		t.Errorf("parent.ChildSpanCount=%q; want %q", got, want)
	}
}

func TestNilSpanFinish(t *testing.T) {
	var span *span
	span.Finish()
}

func TestExecutionTracerTaskEnd(t *testing.T) {
	var n uint64
	ApplyConfig(Config{DefaultSampler: NeverSample()})
	executionTracerTaskEnd := func() {
		atomic.AddUint64(&n, 1)
	}

	var spans []*span
	_, apiSpan := apitrace.GlobalTracer().Start(context.Background(), "foo")
	s := apiSpan.(*span)

	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // never sample

	_, apiSpan = apitrace.GlobalTracer().Start(
		context.Background(),
		"foo",
		apitrace.ChildOf(
			core.SpanContext{
				TraceID:      core.TraceID{High: 0x0102030405060708, Low: 0x090a0b0c0d0e0f},
				SpanID:       uint64(0x0001020304050607),
				TraceOptions: 0,
			},
		),
	)
	s = apiSpan.(*span)
	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // parent not sampled

	ApplyConfig(Config{DefaultSampler: AlwaysSample()})
	_, apiSpan = apitrace.GlobalTracer().Start(context.Background(), "foo")
	s = apiSpan.(*span)
	s.executionTracerTaskEnd = executionTracerTaskEnd
	spans = append(spans, s) // always sample

	for _, span := range spans {
		span.Finish()
	}
	if got, want := n, uint64(len(spans)); got != want {
		t.Fatalf("Execution tracer task ended for %v spans; want %v", got, want)
	}
}

func TestCustomStartEndTime(t *testing.T) {
	startTime := time.Date(2019, time.August, 27, 14, 42, 0, 0, time.UTC)
	endTime := startTime.Add(time.Second * 20)
	_, span := apitrace.Start(
		context.Background(),
		"testspan",
		apitrace.WithStartTime(startTime),
	)
	var te testExporter
	RegisterExporter(&te)
	span.Finish(apitrace.WithFinishTime(endTime))
	UnregisterExporter(&te)
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
