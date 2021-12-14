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

package internal_test

import (
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus/internal"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type span struct {
	trace.Span

	recording bool
	ended     bool
	sc        trace.SpanContext
	name      string
	sCode     codes.Code
	sMsg      string
	attrs     []attribute.KeyValue
	eName     string
	eOpts     []trace.EventOption
}

func (s *span) IsRecording() bool                         { return s.recording }
func (s *span) End(...trace.SpanEndOption)                { s.ended = true }
func (s *span) SpanContext() trace.SpanContext            { return s.sc }
func (s *span) SetName(n string)                          { s.name = n }
func (s *span) SetStatus(c codes.Code, d string)          { s.sCode, s.sMsg = c, d }
func (s *span) SetAttributes(a ...attribute.KeyValue)     { s.attrs = a }
func (s *span) AddEvent(n string, o ...trace.EventOption) { s.eName, s.eOpts = n, o }

func TestSpanIsRecordingEvents(t *testing.T) {
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	if !ocS.IsRecordingEvents() {
		t.Errorf("span.IsRecordingEvents() = false, want true")
	}
	s.recording = false
	if ocS.IsRecordingEvents() {
		t.Errorf("span.IsRecordingEvents() = true, want false")
	}
}

func TestSpanEnd(t *testing.T) {
	s := new(span)
	ocS := internal.NewSpan(s)
	if s.ended {
		t.Fatal("new span already ended")
	}

	ocS.End()
	if !s.ended {
		t.Error("span.End() did not end OpenTelemetry span")
	}
}

func TestSpanSpanContext(t *testing.T) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})
	// Do not test the conversion, only that the method is called.
	converted := otel2oc.SpanContext(sc)

	s := &span{sc: sc}
	ocS := internal.NewSpan(s)
	if ocS.SpanContext() != converted {
		t.Error("span.SpanContext did not use OpenTelemetry SpanContext")
	}
}

func TestSpanSetName(t *testing.T) {
	// OpenCensus does not set a name if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	name := "test name"
	ocS.SetName(name)
	if s.name != name {
		t.Error("span.SetName did not set OpenTelemetry span name")
	}
}

func TestSpanSetStatus(t *testing.T) {
	// OpenCensus does not set a status if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)

	c, d := codes.Error, "error"
	status := octrace.Status{Code: int32(c), Message: d}
	ocS.SetStatus(status)

	if s.sCode != c {
		t.Error("span.SetStatus failed to set OpenTelemetry status code")
	}
	if s.sMsg != d {
		t.Error("span.SetStatus failed to set OpenTelemetry status description")
	}
}

func TestSpanAddAttributes(t *testing.T) {
	attrs := []octrace.Attribute{
		octrace.BoolAttribute("a", true),
	}
	// Do not test the conversion, only that the method is called.
	converted := oc2otel.Attributes(attrs)

	// OpenCensus does not set attributes if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.AddAttributes(attrs...)

	if len(s.attrs) != len(converted) || s.attrs[0] != converted[0] {
		t.Error("span.AddAttributes failed to set OpenTelemetry attributes")
	}
}

func TestSpanAnnotate(t *testing.T) {
	name := "annotation"
	attrs := []octrace.Attribute{
		octrace.BoolAttribute("a", true),
	}
	// Do not test the conversion, only that the method is called.
	want := oc2otel.Attributes(attrs)

	// OpenCensus does not set events if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.Annotate(attrs, name)

	if s.eName != name {
		t.Error("span.Annotate did not set event name")
	}

	config := trace.NewEventConfig(s.eOpts...)
	got := config.Attributes()
	if len(want) != len(got) || want[0] != got[0] {
		t.Error("span.Annotate did not set event options")
	}
}

func TestSpanAnnotatef(t *testing.T) {
	format := "annotation %s"
	attrs := []octrace.Attribute{
		octrace.BoolAttribute("a", true),
	}
	// Do not test the conversion, only that the method is called.
	want := oc2otel.Attributes(attrs)

	// OpenCensus does not set events if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.Annotatef(attrs, format, "a")

	if s.eName != "annotation a" {
		t.Error("span.Annotatef did not set event name")
	}

	config := trace.NewEventConfig(s.eOpts...)
	got := config.Attributes()
	if len(want) != len(got) || want[0] != got[0] {
		t.Error("span.Annotatef did not set event options")
	}
}

func TestSpanAddMessageSendEvent(t *testing.T) {
	var u, c int64 = 1, 2

	// OpenCensus does not set events if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.AddMessageSendEvent(0, u, c)

	if s.eName != internal.MessageSendEvent {
		t.Error("span.AddMessageSendEvent did not set event name")
	}

	config := trace.NewEventConfig(s.eOpts...)
	got := config.Attributes()
	if len(got) != 2 {
		t.Fatalf("span.AddMessageSendEvent set %d attributes, want 2", len(got))
	}

	want := attribute.KeyValue{Key: internal.UncompressedKey, Value: attribute.Int64Value(u)}
	if got[0] != want {
		t.Errorf("span.AddMessageSendEvent wrong uncompressed attribute: %v", got[0])
	}

	want = attribute.KeyValue{Key: internal.CompressedKey, Value: attribute.Int64Value(c)}
	if got[1] != want {
		t.Errorf("span.AddMessageSendEvent wrong compressed attribute: %v", got[1])
	}
}

func TestSpanAddMessageReceiveEvent(t *testing.T) {
	var u, c int64 = 3, 4

	// OpenCensus does not set events if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.AddMessageReceiveEvent(0, u, c)

	if s.eName != internal.MessageReceiveEvent {
		t.Error("span.AddMessageReceiveEvent did not set event name")
	}

	config := trace.NewEventConfig(s.eOpts...)
	got := config.Attributes()
	if len(got) != 2 {
		t.Fatalf("span.AddMessageReceiveEvent set %d attributes, want 2", len(got))
	}

	want := attribute.KeyValue{Key: internal.UncompressedKey, Value: attribute.Int64Value(u)}
	if got[0] != want {
		t.Errorf("span.AddMessageReceiveEvent wrong uncompressed attribute: %v", got[0])
	}

	want = attribute.KeyValue{Key: internal.CompressedKey, Value: attribute.Int64Value(c)}
	if got[1] != want {
		t.Errorf("span.AddMessageReceiveEvent wrong compressed attribute: %v", got[1])
	}
}

func TestSpanAddLinkFails(t *testing.T) {
	h, restore := withHandler()
	defer restore()

	// OpenCensus does not try to set links if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.AddLink(octrace.Link{})

	if h.err == nil {
		t.Error("span.AddLink failed to raise an error")
	}
}

func TestSpanString(t *testing.T) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})

	s := &span{sc: sc}
	ocS := internal.NewSpan(s)
	if expected := "span 0100000000000000"; ocS.String() != expected {
		t.Errorf("span.String = %q, not %q", ocS.String(), expected)
	}
}
