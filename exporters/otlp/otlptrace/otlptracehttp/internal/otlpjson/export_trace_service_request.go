// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package otlpjson implements OTLP JSON Protobuf encoding for trace data.
//
// The encoding conforms to the OTLP specs
// (https://opentelemetry.io/docs/specs/otlp/#json-protobuf-encoding):
//   - trace ID and span ID byte arrays are encoded as case-insensitive hex-encoded strings
//   - enum values encoded as integers
//   - field names in lowerCamelCase
//   - 64-bit integers encoded as quoted decimal strings (ProtoJSON specs)
//
// Common OTLP/JSON types and scalar codecs are generated from
// internal/shared/otlp/otlpjson.
package otlpjson

import (
	"encoding/json"

	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

// ExportTraceServiceRequest corresponds to coltracepb.ExportTraceServiceRequest.
type ExportTraceServiceRequest struct {
	ResourceSpans []*ResourceSpans `json:"resourceSpans,omitempty"`
}

// ResourceSpans corresponds to tracepb.ResourceSpans.
type ResourceSpans struct {
	Resource   *Resource     `json:"resource,omitempty"`
	ScopeSpans []*ScopeSpans `json:"scopeSpans,omitempty"`
	SchemaURL  string        `json:"schemaUrl,omitempty"`
}

// ScopeSpans corresponds to tracepb.ScopeSpans.
type ScopeSpans struct {
	Scope     *InstrumentationScope `json:"scope,omitempty"`
	Spans     []*Span               `json:"spans,omitempty"`
	SchemaURL string                `json:"schemaUrl,omitempty"`
}

// Span corresponds to tracepb.Span.
type Span struct {
	TraceID                TraceID     `json:"traceId"`
	SpanID                 SpanID      `json:"spanId"`
	TraceState             string      `json:"traceState,omitempty"`
	ParentSpanID           *SpanID     `json:"parentSpanId,omitempty"`
	Flags                  uint32      `json:"flags,omitempty"`
	Name                   string      `json:"name,omitempty"`
	Kind                   int32       `json:"kind,omitempty"`
	StartTimeUnixNano      Uint64      `json:"startTimeUnixNano,omitempty"`
	EndTimeUnixNano        Uint64      `json:"endTimeUnixNano,omitempty"`
	Attributes             []*KeyValue `json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `json:"droppedAttributesCount,omitempty"`
	Events                 []*Event    `json:"events,omitempty"`
	DroppedEventsCount     uint32      `json:"droppedEventsCount,omitempty"`
	Links                  []*Link     `json:"links,omitempty"`
	DroppedLinksCount      uint32      `json:"droppedLinksCount,omitempty"`
	Status                 *Status     `json:"status,omitempty"`
}

// Event corresponds to tracepb.Span_Event.
type Event struct {
	TimeUnixNano           Uint64      `json:"timeUnixNano,omitempty"`
	Name                   string      `json:"name,omitempty"`
	Attributes             []*KeyValue `json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `json:"droppedAttributesCount,omitempty"`
}

// Link corresponds to tracepb.Span_Link.
type Link struct {
	TraceID                TraceID     `json:"traceId"`
	SpanID                 SpanID      `json:"spanId"`
	TraceState             string      `json:"traceState,omitempty"`
	Attributes             []*KeyValue `json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `json:"droppedAttributesCount,omitempty"`
	Flags                  uint32      `json:"flags,omitempty"`
}

// Status corresponds to tracepb.Status.
type Status struct {
	Message string `json:"message,omitempty"`
	Code    int32  `json:"code,omitempty"`
}

// MarshalExportTraceServiceRequest encodes an ExportTraceServiceRequest as JSON Protobuf encoded bytes.
func MarshalExportTraceServiceRequest(req *coltracepb.ExportTraceServiceRequest) ([]byte, error) {
	if req == nil {
		return []byte("{}"), nil
	}
	r := &ExportTraceServiceRequest{}
	for _, rs := range req.ResourceSpans {
		r.ResourceSpans = append(r.ResourceSpans, encodeResourceSpans(rs))
	}
	return json.Marshal(r)
}

func encodeResourceSpans(rs *tracepb.ResourceSpans) *ResourceSpans {
	if rs == nil {
		return nil
	}
	out := &ResourceSpans{SchemaURL: rs.SchemaUrl}
	if rs.Resource != nil {
		out.Resource = &Resource{
			Attributes:             encodeKeyValues(rs.Resource.Attributes),
			DroppedAttributesCount: rs.Resource.DroppedAttributesCount,
			EntityRefs:             encodeEntityRefs(rs.Resource.EntityRefs),
		}
	}
	for _, ss := range rs.ScopeSpans {
		out.ScopeSpans = append(out.ScopeSpans, encodeScopeSpans(ss))
	}
	return out
}

func encodeScopeSpans(ss *tracepb.ScopeSpans) *ScopeSpans {
	if ss == nil {
		return nil
	}
	out := &ScopeSpans{SchemaURL: ss.SchemaUrl}
	if ss.Scope != nil {
		out.Scope = &InstrumentationScope{
			Name:                   ss.Scope.Name,
			Version:                ss.Scope.Version,
			Attributes:             encodeKeyValues(ss.Scope.Attributes),
			DroppedAttributesCount: ss.Scope.DroppedAttributesCount,
		}
	}
	for _, s := range ss.Spans {
		out.Spans = append(out.Spans, encodeSpan(s))
	}
	return out
}

func encodeSpan(s *tracepb.Span) *Span {
	if s == nil {
		return nil
	}
	out := &Span{
		TraceState:             s.TraceState,
		Flags:                  s.Flags,
		Name:                   s.Name,
		Kind:                   int32(s.Kind),
		StartTimeUnixNano:      Uint64(s.StartTimeUnixNano),
		EndTimeUnixNano:        Uint64(s.EndTimeUnixNano),
		Attributes:             encodeKeyValues(s.Attributes),
		DroppedAttributesCount: s.DroppedAttributesCount,
		DroppedEventsCount:     s.DroppedEventsCount,
		DroppedLinksCount:      s.DroppedLinksCount,
	}

	copy(out.TraceID[:], s.TraceId)
	copy(out.SpanID[:], s.SpanId)

	if len(s.ParentSpanId) > 0 {
		var psid SpanID
		copy(psid[:], s.ParentSpanId)
		if psid != (SpanID{}) {
			out.ParentSpanID = &psid
		}
	}

	for _, e := range s.Events {
		if e == nil {
			continue
		}
		out.Events = append(out.Events, &Event{
			TimeUnixNano:           Uint64(e.TimeUnixNano),
			Name:                   e.Name,
			Attributes:             encodeKeyValues(e.Attributes),
			DroppedAttributesCount: e.DroppedAttributesCount,
		})
	}
	for _, l := range s.Links {
		out.Links = append(out.Links, encodeLink(l))
	}
	if s.Status != nil {
		out.Status = &Status{
			Message: s.Status.Message,
			Code:    int32(s.Status.Code),
		}
	}
	return out
}

func encodeLink(l *tracepb.Span_Link) *Link {
	if l == nil {
		return nil
	}
	out := &Link{
		TraceState:             l.TraceState,
		Attributes:             encodeKeyValues(l.Attributes),
		DroppedAttributesCount: l.DroppedAttributesCount,
		Flags:                  l.Flags,
	}
	copy(out.TraceID[:], l.TraceId)
	copy(out.SpanID[:], l.SpanId)
	return out
}

// UnmarshalExportTraceServiceRequest decodes JSON Protobuf encoded payload into an ExportTraceServiceRequest.
func UnmarshalExportTraceServiceRequest(data []byte, req *coltracepb.ExportTraceServiceRequest) error {
	var jr ExportTraceServiceRequest
	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}

	for _, rs := range jr.ResourceSpans {
		req.ResourceSpans = append(req.ResourceSpans, decodeResourceSpans(rs))
	}
	return nil
}

func decodeResourceSpans(jrs *ResourceSpans) *tracepb.ResourceSpans {
	rs := &tracepb.ResourceSpans{SchemaUrl: jrs.SchemaURL}
	if jrs.Resource != nil {
		rs.Resource = &resourcepb.Resource{
			Attributes:             decodeKeyValues(jrs.Resource.Attributes),
			DroppedAttributesCount: jrs.Resource.DroppedAttributesCount,
			EntityRefs:             decodeEntityRefs(jrs.Resource.EntityRefs),
		}
	}
	for _, ss := range jrs.ScopeSpans {
		rs.ScopeSpans = append(rs.ScopeSpans, decodeScopeSpans(ss))
	}
	return rs
}

func decodeScopeSpans(jss *ScopeSpans) *tracepb.ScopeSpans {
	ss := &tracepb.ScopeSpans{SchemaUrl: jss.SchemaURL}
	if jss.Scope != nil {
		ss.Scope = &commonpb.InstrumentationScope{
			Name:                   jss.Scope.Name,
			Version:                jss.Scope.Version,
			Attributes:             decodeKeyValues(jss.Scope.Attributes),
			DroppedAttributesCount: jss.Scope.DroppedAttributesCount,
		}
	}
	for _, s := range jss.Spans {
		ss.Spans = append(ss.Spans, decodeSpan(s))
	}
	return ss
}

func decodeSpan(js *Span) *tracepb.Span {
	s := &tracepb.Span{
		TraceId:                js.TraceID[:],
		SpanId:                 js.SpanID[:],
		TraceState:             js.TraceState,
		Flags:                  js.Flags,
		Name:                   js.Name,
		Kind:                   tracepb.Span_SpanKind(js.Kind),
		StartTimeUnixNano:      uint64(js.StartTimeUnixNano),
		EndTimeUnixNano:        uint64(js.EndTimeUnixNano),
		Attributes:             decodeKeyValues(js.Attributes),
		DroppedAttributesCount: js.DroppedAttributesCount,
		DroppedEventsCount:     js.DroppedEventsCount,
		DroppedLinksCount:      js.DroppedLinksCount,
	}
	if js.ParentSpanID != nil {
		s.ParentSpanId = js.ParentSpanID[:]
	}
	for _, e := range js.Events {
		s.Events = append(s.Events, &tracepb.Span_Event{
			TimeUnixNano:           uint64(e.TimeUnixNano),
			Name:                   e.Name,
			Attributes:             decodeKeyValues(e.Attributes),
			DroppedAttributesCount: e.DroppedAttributesCount,
		})
	}
	for _, l := range js.Links {
		s.Links = append(s.Links, decodeLink(l))
	}
	if js.Status != nil {
		s.Status = &tracepb.Status{
			Message: js.Status.Message,
			Code:    tracepb.Status_StatusCode(js.Status.Code),
		}
	}
	return s
}

func decodeLink(jl *Link) *tracepb.Span_Link {
	return &tracepb.Span_Link{
		TraceId:                jl.TraceID[:],
		SpanId:                 jl.SpanID[:],
		TraceState:             jl.TraceState,
		Attributes:             decodeKeyValues(jl.Attributes),
		DroppedAttributesCount: jl.DroppedAttributesCount,
		Flags:                  jl.Flags,
	}
}
