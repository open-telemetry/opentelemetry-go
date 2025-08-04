// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"go.opentelemetry.io/otel/trace/internal/telemetry"
)

var (
	y2k = time.Unix(0, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()) // No location.

	attrsA = []telemetry.Attr{
		telemetry.String("user", "Alice"),
		telemetry.Bool("admin", true),
		telemetry.Int64("floor", -2),
		telemetry.Float64("impact", 0.21362),
		telemetry.Slice(
			"reports",
			telemetry.StringValue("Bob"),
			telemetry.StringValue("Dave"),
		),
		telemetry.Map(
			"favorites",
			telemetry.String("food", "hot dog"),
			telemetry.Int("number", 13),
		),
		telemetry.Bytes(
			"secret",
			[]byte("NUI4RUZGRjc5ODAzODEwM0QyNjlCNjMzODEzRkM2MEM="),
		),
	}
	pAttrsA = func() pcommon.Map {
		m := pcommon.NewMap()
		m.PutStr("user", "Alice")
		m.PutBool("admin", true)
		m.PutInt("floor", -2)
		m.PutDouble("impact", 0.21362)

		s := m.PutEmptySlice("reports")
		s.AppendEmpty().SetStr("Bob")
		s.AppendEmpty().SetStr("Dave")

		fav := m.PutEmptyMap("favorites")
		fav.PutStr("food", "hot dog")
		fav.PutInt("number", 13)

		sec := m.PutEmptyBytes("secret")
		sec.FromRaw([]byte("NUI4RUZGRjc5ODAzODEwM0QyNjlCNjMzODEzRkM2MEM="))

		return m
	}()

	link = &telemetry.SpanLink{
		TraceID:      telemetry.TraceID{0x2},
		SpanID:       telemetry.SpanID{0x1},
		TraceState:   "test=green",
		Attrs:        []telemetry.Attr{telemetry.Int("queue", 17)},
		DroppedAttrs: 8,
		Flags:        1,
	}
	pLink = func() ptrace.SpanLink {
		l := ptrace.NewSpanLink()
		l.SetTraceID(pcommon.TraceID{0x2})
		l.SetSpanID(pcommon.SpanID{0x1})
		l.TraceState().FromRaw("test=green")
		l.Attributes().PutInt("queue", 17)
		l.SetDroppedAttributesCount(8)
		l.SetFlags(1)
		return l
	}()

	event = &telemetry.SpanEvent{
		Time:         y2k.Add(10 * time.Microsecond),
		Name:         "span.event",
		Attrs:        []telemetry.Attr{telemetry.Float64("impact", 0.4372)},
		DroppedAttrs: 2,
	}
	prevent = func() ptrace.SpanEvent {
		e := ptrace.NewSpanEvent()
		e.SetTimestamp(pcommon.NewTimestampFromTime(y2k.Add(10 * time.Microsecond)))
		e.SetName("span.event")
		e.Attributes().PutDouble("impact", 0.4372)
		e.SetDroppedAttributesCount(2)
		return e
	}()

	spanA = &telemetry.Span{
		TraceID:       [16]byte{0x1},
		SpanID:        [8]byte{0x2},
		TraceState:    "test=a",
		ParentSpanID:  [8]byte{0x1},
		Flags:         1,
		Name:          "span.a",
		Kind:          telemetry.SpanKindClient,
		StartTime:     y2k,
		EndTime:       y2k.Add(time.Second),
		Attrs:         attrsA,
		DroppedAttrs:  2,
		Events:        []*telemetry.SpanEvent{event},
		DroppedEvents: 3,
		Links:         []*telemetry.SpanLink{link},
		DroppedLinks:  4,
		Status: &telemetry.Status{
			Message: "okay",
			Code:    telemetry.StatusCodeOK,
		},
	}
	pSpanA = func() ptrace.Span {
		s := ptrace.NewSpan()
		s.SetTraceID(pcommon.TraceID([16]byte{0x1}))
		s.SetSpanID(pcommon.SpanID([8]byte{0x2}))

		ts := s.TraceState()
		ts.FromRaw("test=a")

		s.SetParentSpanID(pcommon.SpanID([8]byte{0x1}))
		s.SetFlags(1)
		s.SetName("span.a")
		s.SetKind(ptrace.SpanKindClient)
		s.SetStartTimestamp(pcommon.NewTimestampFromTime(y2k))
		s.SetEndTimestamp(pcommon.NewTimestampFromTime(y2k.Add(time.Second)))
		pAttrsA.CopyTo(s.Attributes())
		s.SetDroppedAttributesCount(2)
		prevent.CopyTo(s.Events().AppendEmpty())
		s.SetDroppedEventsCount(3)
		pLink.CopyTo(s.Links().AppendEmpty())
		s.SetDroppedLinksCount(4)

		stat := s.Status()
		stat.SetMessage("okay")
		stat.SetCode(ptrace.StatusCodeOk)

		return s
	}()
	schema100 = "http://go.opentelemetry.io/schema/v1.0.0"

	scope = &telemetry.Scope{
		Name:         "go.opentelemetry.io/otel/trace/internal/telemetry/test",
		Version:      "v0.0.1",
		Attrs:        []telemetry.Attr{telemetry.String("department", "ops")},
		DroppedAttrs: 1,
	}
	pScope = func() pcommon.InstrumentationScope {
		s := pcommon.NewInstrumentationScope()
		s.SetName("go.opentelemetry.io/otel/trace/internal/telemetry/test")
		s.SetVersion("v0.0.1")
		s.Attributes().PutStr("department", "ops")
		s.SetDroppedAttributesCount(1)
		return s
	}()

	scopeSpans = &telemetry.ScopeSpans{
		Scope:     scope,
		Spans:     []*telemetry.Span{spanA},
		SchemaURL: schema100,
	}
	pScopeSpans = func() ptrace.ScopeSpans {
		s := ptrace.NewScopeSpans()
		pSpanA.CopyTo(s.Spans().AppendEmpty())
		pScope.CopyTo(s.Scope())
		s.SetSchemaUrl(schema100)
		return s
	}()

	res = telemetry.Resource{
		Attrs: []telemetry.Attr{
			telemetry.String("host", "hal"),
			telemetry.Int("id", 42),
		},
		DroppedAttrs: 100,
	}
	press = func() pcommon.Resource {
		r := pcommon.NewResource()
		r.Attributes().PutStr("host", "hal")
		r.Attributes().PutInt("id", 42)
		r.SetDroppedAttributesCount(100)
		return r
	}()

	resSpans = &telemetry.ResourceSpans{
		Resource:   res,
		SchemaURL:  schema100,
		ScopeSpans: []*telemetry.ScopeSpans{scopeSpans},
	}
	pResSpans = func() ptrace.ResourceSpans {
		rs := ptrace.NewResourceSpans()
		press.CopyTo(rs.Resource())
		pScopeSpans.CopyTo(rs.ScopeSpans().AppendEmpty())
		rs.SetSchemaUrl(schema100)
		return rs
	}()

	traces = telemetry.Traces{
		ResourceSpans: []*telemetry.ResourceSpans{
			resSpans,
		},
	}
	pTraces = func() ptrace.Traces {
		traces := ptrace.NewTraces()
		pResSpans.CopyTo(traces.ResourceSpans().AppendEmpty())
		return traces
	}()
)

func TestDecode(t *testing.T) {
	var enc ptrace.JSONMarshaler
	b, err := enc.MarshalTraces(pTraces)
	require.NoError(t, err)

	t.Log(string(b)) // This helps when test fails to understand what is being decoded.

	var got telemetry.Traces
	dec := json.NewDecoder(bytes.NewReader(b))
	require.NoError(t, dec.Decode(&got))

	assert.Equal(t, traces, got)
}

func TestEncode(t *testing.T) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	require.NoError(t, enc.Encode(traces))

	data := buf.Bytes()
	t.Log(string(data)) // This helps when test fails to understand what how the data has been encoded.

	var dec ptrace.JSONUnmarshaler
	got, err := dec.UnmarshalTraces(data)
	require.NoError(t, err)

	assert.Equal(t, pTraces, got)
}
