// Copyright 2020, OpenTelemetry Authors
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

package otlp

import (
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/sdk/resource"

	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	maxMessageEventsPerSpan = 128
)

func otResourceToProtoResource(res *resource.Resource) *resourcepb.Resource {
	if res == nil {
		return nil
	}
	resProto := &resourcepb.Resource{
		Attributes: otAttributesToProtoAttributes(res.Attributes()),
	}
	return resProto
}

func otSpanToProtoSpan(sd *export.SpanData) *tracepb.Span {
	if sd == nil {
		return nil
	}
	return &tracepb.Span{
		TraceId:           sd.SpanContext.TraceID[:],
		SpanId:            sd.SpanContext.SpanID[:],
		ParentSpanId:      sd.ParentSpanID[:],
		Status:            otStatusToProtoStatus(sd.StatusCode, sd.StatusMessage),
		StartTimeUnixnano: uint64(sd.StartTime.Nanosecond()),
		EndTimeUnixnano:   uint64(sd.EndTime.Nanosecond()),
		Links:             otLinksToProtoLinks(sd.Links),
		Kind:              otSpanKindToProtoSpanKind(sd.SpanKind),
		Name:              sd.Name,
		Attributes:        otAttributesToProtoAttributes(sd.Attributes),
		Events:            otTimeEventsToProtoTimeEvents(sd.MessageEvents),
		// TODO (rghetia): Add Tracestate: when supported.
		DroppedAttributesCount: uint32(sd.DroppedAttributeCount),
		DroppedEventsCount:     uint32(sd.DroppedMessageEventCount),
		DroppedLinksCount:      uint32(sd.DroppedLinkCount),
	}
}

func otStatusToProtoStatus(status codes.Code, message string) *tracepb.Status {
	return &tracepb.Status{
		Code:    tracepb.Status_StatusCode(status),
		Message: message,
	}
}

func otLinksToProtoLinks(links []apitrace.Link) []*tracepb.Span_Link {
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
			Attributes: otAttributesToProtoAttributes(otLink.Attributes),
		})
	}
	return sl
}

func otAttributesToProtoAttributes(attrs []core.KeyValue) []*commonpb.AttributeKeyValue {
	if len(attrs) == 0 {
		return nil
	}
	out := make([]*commonpb.AttributeKeyValue, 0, len(attrs))
	for _, v := range attrs {
		switch v.Value.Type() {
		case core.BOOL:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:       string(v.Key),
				Type:      commonpb.AttributeKeyValue_BOOL,
				BoolValue: v.Value.AsBool(),
			})
		case core.INT64, core.INT32, core.UINT32, core.UINT64:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:      string(v.Key),
				Type:     commonpb.AttributeKeyValue_INT,
				IntValue: v.Value.AsInt64(),
			})
		case core.FLOAT32:
			f32 := v.Value.AsFloat32()
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_DOUBLE,
				DoubleValue: float64(f32),
			})
		case core.FLOAT64:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_DOUBLE,
				DoubleValue: v.Value.AsFloat64(),
			})
		case core.STRING:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_STRING,
				StringValue: v.Value.AsString(),
			})
		}
	}
	return out
}

func otTimeEventsToProtoTimeEvents(es []export.Event) []*tracepb.Span_Event {
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
				TimeUnixnano: uint64(e.Time.Nanosecond()),
				Attributes:   otAttributesToProtoAttributes(e.Attributes),
				// TODO (rghetia) : Add Drop Counts when supported.
			},
		)
	}

	return events
}

func otSpanKindToProtoSpanKind(kind apitrace.SpanKind) tracepb.Span_SpanKind {
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
