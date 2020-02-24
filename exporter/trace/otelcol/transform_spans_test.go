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

package otelcol_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"

	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/otelcol"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

func TestOCSpanToProtoSpan_endToEnd(t *testing.T) {
	// The goal of this test is to ensure that each
	// spanData is transformed and exported correctly!

	collector := runMockCol(t)
	defer func() {
		_ = collector.stop()
	}()

	serviceName := "spanTranslation"
	exp, err := otelcol.NewExporter(otelcol.WithInsecure(),
		otelcol.WithAddress(collector.address),
		otelcol.WithReconnectionPeriod(50*time.Millisecond),
		otelcol.WithServiceName(serviceName))
	if err != nil {
		t.Fatalf("Failed to create a new collector exporter: %v", err)
	}
	defer func() {
		_ = exp.Stop()
	}()

	// Give the background collector connection sometime to setup.
	<-time.After(20 * time.Millisecond)

	startTime := time.Now()
	endTime := startTime.Add(10 * time.Second)

	otSpanData := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID: core.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:  core.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		},
		SpanKind:     apitrace.SpanKindServer,
		ParentSpanID: core.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		Name:         "End-To-End Here",
		StartTime:    startTime,
		EndTime:      endTime,
		MessageEvents: []export.Event{
			{Time: startTime,
				Attributes: []core.KeyValue{
					core.Key("CompressedByteSize").Uint64(512),
					core.Key("UncompressedByteSize").Uint64(1024),
					core.Key("MessageEventType").String("Sent"),
				},
			},
			{Time: endTime,
				Attributes: []core.KeyValue{
					core.Key("CompressedByteSize").Uint64(1000),
					core.Key("UncompressedByteSize").Uint64(1024),
					core.Key("MessageEventType").String("Recv"),
				},
			},
		},
		Links: []apitrace.Link{
			{
				SpanContext: core.SpanContext{
					TraceID:    core.TraceID{0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF},
					SpanID:     core.SpanID{0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7},
					TraceFlags: 0,
				},
				Attributes: []core.KeyValue{
					core.Key("LinkType").String("Parent"),
				},
			},
			{
				SpanContext: core.SpanContext{
					TraceID:    core.TraceID{0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF},
					SpanID:     core.SpanID{0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7},
					TraceFlags: 0,
				},
				Attributes: []core.KeyValue{
					core.Key("LinkType").String("Child"),
				},
			},
		},
		Status:          codes.Internal,
		HasRemoteParent: true,
		Attributes: []core.KeyValue{
			core.Key("timeout_ns").Int64(12e9),
			core.Key("agent").String("otelcol"),
			core.Key("cache_hit").Bool(true),
			core.Key("ping_count").Int(25), // Should be transformed into int64
		},
		DroppedAttributeCount:    1,
		DroppedMessageEventCount: 2,
		DroppedLinkCount:         3,
	}

	exp.ExportSpans(context.Background(), []*export.SpanData{otSpanData})
	// Also try to export a nil span and it should never make it
	exp.ExportSpans(context.Background(), nil)

	_ = exp.Stop()
	_ = collector.stop()

	spans := collector.getSpans()
	if len(spans) == 0 || spans[0] == nil {
		t.Fatal("Expected the exported span")
	}

	wantProtoSpan := &tracepb.Span{
		TraceId:           []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		SpanId:            []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		ParentSpanId:      []byte{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		Name:              "End-To-End Here",
		Kind:              tracepb.Span_SERVER,
		StartTimeUnixnano: uint64(startTime.Nanosecond()),
		EndTimeUnixnano:   uint64(endTime.Nanosecond()),
		Status: &tracepb.Status{
			Code: 13,
		},
		Events: []*tracepb.Span_Event{
			{
				TimeUnixnano: uint64(startTime.Nanosecond()),
				Attributes: []*commonpb.AttributeKeyValue{
					{
						Key:         "CompressedByteSize",
						Type:        commonpb.AttributeKeyValue_INT,
						StringValue: "",
						IntValue:    512,
						DoubleValue: 0,
						BoolValue:   false,
					},
					{
						Key:         "UncompressedByteSize",
						Type:        commonpb.AttributeKeyValue_INT,
						StringValue: "",
						IntValue:    1024,
						DoubleValue: 0,
						BoolValue:   false,
					},
					{
						Key:         "MessageEventType",
						Type:        commonpb.AttributeKeyValue_STRING,
						StringValue: "Sent",
						IntValue:    0,
						DoubleValue: 0,
						BoolValue:   false,
					},
				},
			},
			{
				TimeUnixnano: uint64(endTime.Nanosecond()),
				Attributes: []*commonpb.AttributeKeyValue{
					{
						Key:         "CompressedByteSize",
						Type:        commonpb.AttributeKeyValue_INT,
						StringValue: "",
						IntValue:    1000,
						DoubleValue: 0,
						BoolValue:   false,
					},
					{
						Key:         "UncompressedByteSize",
						Type:        commonpb.AttributeKeyValue_INT,
						StringValue: "",
						IntValue:    1024,
						DoubleValue: 0,
						BoolValue:   false,
					},
					{
						Key:         "MessageEventType",
						Type:        commonpb.AttributeKeyValue_STRING,
						StringValue: "Recv",
						IntValue:    0,
						DoubleValue: 0,
						BoolValue:   false,
					},
				},
			},
		},
		Links: []*tracepb.Span_Link{
			{
				TraceId: []byte{0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF},
				SpanId:  []byte{0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7},
				Attributes: []*commonpb.AttributeKeyValue{
					{
						Key:         "LinkType",
						Type:        commonpb.AttributeKeyValue_STRING,
						StringValue: "Parent",
						IntValue:    0,
						DoubleValue: 0,
						BoolValue:   false,
					},
				},
			},
			{
				TraceId: []byte{0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF},
				SpanId:  []byte{0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7},
				Attributes: []*commonpb.AttributeKeyValue{
					{
						Key:         "LinkType",
						Type:        commonpb.AttributeKeyValue_STRING,
						StringValue: "Child",
						IntValue:    0,
						DoubleValue: 0,
						BoolValue:   false,
					},
				},
			},
		},
		Attributes: []*commonpb.AttributeKeyValue{
			{
				Key:         "timeout_ns",
				Type:        commonpb.AttributeKeyValue_INT,
				StringValue: "",
				IntValue:    12e9,
				DoubleValue: 0,
				BoolValue:   false,
			},
			{
				Key:         "agent",
				Type:        commonpb.AttributeKeyValue_STRING,
				StringValue: "otelcol",
				IntValue:    0,
				DoubleValue: 0,
				BoolValue:   false,
			},
			{
				Key:         "cache_hit",
				Type:        commonpb.AttributeKeyValue_BOOL,
				StringValue: "",
				IntValue:    0,
				DoubleValue: 0,
				BoolValue:   true,
			},
			{
				Key:         "ping_count",
				Type:        commonpb.AttributeKeyValue_INT,
				StringValue: "",
				IntValue:    25,
				DoubleValue: 0,
				BoolValue:   false,
			},
		},
		DroppedAttributesCount: 1,
		DroppedEventsCount:     2,
		DroppedLinksCount:      3,
	}

	if diff := cmp.Diff(spans[0], wantProtoSpan, cmp.Comparer(proto.Equal)); diff != "" {
		t.Fatalf("End-to-end transformed span differs %v\n", diff)
	}
}
