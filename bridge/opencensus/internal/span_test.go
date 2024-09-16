// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
	links     []trace.Link
}

func (s *span) IsRecording() bool                         { return s.recording }
func (s *span) End(...trace.SpanEndOption)                { s.ended = true }
func (s *span) SpanContext() trace.SpanContext            { return s.sc }
func (s *span) SetName(n string)                          { s.name = n }
func (s *span) SetStatus(c codes.Code, d string)          { s.sCode, s.sMsg = c, d }
func (s *span) SetAttributes(a ...attribute.KeyValue)     { s.attrs = a }
func (s *span) AddEvent(n string, o ...trace.EventOption) { s.eName, s.eOpts = n, o }
func (s *span) AddLink(l trace.Link)                      { s.links = append(s.links, l) }

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

	for _, tt := range []struct {
		name string

		code    int32
		message string

		wantCode codes.Code
	}{
		{
			name:    "with an error code",
			code:    int32(codes.Error),
			message: "error",

			wantCode: codes.Error,
		},
		{
			name:    "with a negative/invalid code",
			code:    -42,
			message: "error",

			wantCode: codes.Unset,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			status := octrace.Status{Code: tt.code, Message: tt.message}
			ocS.SetStatus(status)

			if s.sCode != tt.wantCode {
				t.Errorf("span.SetStatus failed to set OpenTelemetry status code. Expected %d, got %d", tt.wantCode, s.sCode)
			}
			if s.sMsg != tt.message {
				t.Errorf("span.SetStatus failed to set OpenTelemetry status description. Expected %s, got %s", tt.message, s.sMsg)
			}
		})
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
	// OpenCensus does not try to set links if not recording.
	s := &span{recording: true}
	ocS := internal.NewSpan(s)
	ocS.AddLink(octrace.Link{})
	ocS.AddLink(octrace.Link{
		TraceID: octrace.TraceID([16]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		SpanID:  octrace.SpanID([8]byte{2, 0, 0, 0, 0, 0, 0, 0}),
		Attributes: map[string]interface{}{
			"foo":    "bar",
			"number": int64(3),
		},
	})

	wantLinks := []trace.Link{
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceFlags: trace.FlagsSampled,
			}),
		},
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
				SpanID:     trace.SpanID([]byte{2, 0, 0, 0, 0, 0, 0, 0}),
				TraceFlags: trace.FlagsSampled,
			}),
			Attributes: []attribute.KeyValue{
				attribute.String("foo", "bar"),
				attribute.Int64("number", 3),
			},
		},
	}

	if len(s.links) != len(wantLinks) {
		t.Fatalf("got wrong number of links; want %v, got %v", len(wantLinks), len(s.links))
	}

	for i, l := range s.links {
		if !l.SpanContext.Equal(wantLinks[i].SpanContext) {
			t.Errorf("link[%v] has the wrong span context; want %+v, got %+v", i, wantLinks[i].SpanContext, l.SpanContext)
		}
		gotAttributeSet := attribute.NewSet(l.Attributes...)
		wantAttributeSet := attribute.NewSet(wantLinks[i].Attributes...)
		if !gotAttributeSet.Equals(&wantAttributeSet) {
			t.Errorf("link[%v] has the wrong attributes; want %v, got %v", i, wantAttributeSet.Encoded(attribute.DefaultEncoder()), gotAttributeSet.Encoded(attribute.DefaultEncoder()))
		}
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
