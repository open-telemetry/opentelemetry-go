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

package transform

import (
	"google.golang.org/grpc/codes"

	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"

	"go.opentelemetry.io/otel/api/label"
	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	maxMessageEventsPerSpan = 128
)

// SpanData transforms a slice of SpanData into a slice of OTLP ResourceSpans.
func SpanData(sdl []*export.SpanData) []*tracepb.ResourceSpans {
	if len(sdl) == 0 {
		return nil
	}
	// Group by the distinct representation of the Resource.
	rsm := make(map[label.Distinct]*tracepb.ResourceSpans)

	for _, sd := range sdl {
		if sd != nil {
			key := sd.Resource.Equivalent()

			rs, ok := rsm[key]
			if !ok {
				rs = &tracepb.ResourceSpans{
					Resource: Resource(sd.Resource),
					InstrumentationLibrarySpans: []*tracepb.InstrumentationLibrarySpans{
						{
							Spans: []*tracepb.Span{},
						},
					},
				}
				rsm[key] = rs
			}
			rs.InstrumentationLibrarySpans[0].Spans =
				append(rs.InstrumentationLibrarySpans[0].Spans, span(sd))
		}
	}
	rss := make([]*tracepb.ResourceSpans, 0, len(rsm))
	for _, rs := range rsm {
		rss = append(rss, rs)
	}
	return rss
}

// span transforms a Span into an OTLP span.
func span(sd *export.SpanData) *tracepb.Span {
	if sd == nil {
		return nil
	}

	s := &tracepb.Span{
		TraceId:           sd.SpanContext.TraceID[:],
		SpanId:            sd.SpanContext.SpanID[:],
		Status:            status(sd.StatusCode, sd.StatusMessage),
		StartTimeUnixNano: uint64(sd.StartTime.UnixNano()),
		EndTimeUnixNano:   uint64(sd.EndTime.UnixNano()),
		Links:             links(sd.Links),
		Kind:              spanKind(sd.SpanKind),
		Name:              sd.Name,
		Attributes:        Attributes(sd.Attributes),
		Events:            spanEvents(sd.MessageEvents),
		// TODO (rghetia): Add Tracestate: when supported.
		DroppedAttributesCount: uint32(sd.DroppedAttributeCount),
		DroppedEventsCount:     uint32(sd.DroppedMessageEventCount),
		DroppedLinksCount:      uint32(sd.DroppedLinkCount),
	}

	if sd.ParentSpanID.IsValid() {
		s.ParentSpanId = sd.ParentSpanID[:]
	}

	return s
}

// status transform a span code and message into an OTLP span status.
func status(status codes.Code, message string) *tracepb.Status {
	return &tracepb.Status{
		Code:    tracepb.Status_StatusCode(status),
		Message: message,
	}
}

// links transforms span Links to OTLP span links.
func links(links []apitrace.Link) []*tracepb.Span_Link {
	if len(links) == 0 {
		return nil
	}

	sl := make([]*tracepb.Span_Link, 0, len(links))
	for _, otLink := range links {
		// This redefinition is necessary to prevent otLink.*ID[:] copies
		// being reused -- in short we need a new otLink per iteration.
		otLink := otLink

		sl = append(sl, &tracepb.Span_Link{
			TraceId:    otLink.TraceID[:],
			SpanId:     otLink.SpanID[:],
			Attributes: Attributes(otLink.Attributes),
		})
	}
	return sl
}

// spanEvents transforms span Events to an OTLP span events.
func spanEvents(es []export.Event) []*tracepb.Span_Event {
	if len(es) == 0 {
		return nil
	}

	evCount := len(es)
	if evCount > maxMessageEventsPerSpan {
		evCount = maxMessageEventsPerSpan
	}
	events := make([]*tracepb.Span_Event, 0, evCount)
	messageEvents := 0

	// Transform message events
	for _, e := range es {
		if messageEvents >= maxMessageEventsPerSpan {
			break
		}
		messageEvents++
		events = append(events,
			&tracepb.Span_Event{
				Name:         e.Name,
				TimeUnixNano: uint64(e.Time.Nanosecond()),
				Attributes:   Attributes(e.Attributes),
				// TODO (rghetia) : Add Drop Counts when supported.
			},
		)
	}

	return events
}

// spanKind transforms a SpanKind to an OTLP span kind.
func spanKind(kind apitrace.SpanKind) tracepb.Span_SpanKind {
	switch kind {
	case apitrace.SpanKindInternal:
		return tracepb.Span_INTERNAL
	case apitrace.SpanKindClient:
		return tracepb.Span_CLIENT
	case apitrace.SpanKindServer:
		return tracepb.Span_SERVER
	case apitrace.SpanKindProducer:
		return tracepb.Span_PRODUCER
	case apitrace.SpanKindConsumer:
		return tracepb.Span_CONSUMER
	default:
		return tracepb.Span_SPAN_KIND_UNSPECIFIED
	}
}
