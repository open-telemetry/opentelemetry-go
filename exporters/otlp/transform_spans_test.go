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

package otlp_test

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
	"go.opentelemetry.io/otel/exporters/otlp"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type testCases struct {
	otSpan   *export.SpanData
	otlpSpan *tracepb.Span
}

func TestOtSpanToOtlpSpan_Basic(t *testing.T) {
	// The goal of this test is to ensure that each
	// spanData is transformed and exported correctly!
	testAndVerify("Basic End-2-End", t, func(t *testing.T) []testCases {

		startTime := time.Now()
		endTime := startTime.Add(10 * time.Second)

		tcs := []testCases{
			{
				otSpan: &export.SpanData{
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
							},
						},
						{Time: endTime,
							Attributes: []core.KeyValue{
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
					},
					DroppedAttributeCount:    1,
					DroppedMessageEventCount: 2,
					DroppedLinkCount:         3,
				},
				otlpSpan: &tracepb.Span{
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
							},
						},
						{
							TimeUnixnano: uint64(endTime.Nanosecond()),
							Attributes: []*commonpb.AttributeKeyValue{
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
					},
					DroppedAttributesCount: 1,
					DroppedEventsCount:     2,
					DroppedLinksCount:      3,
				},
			},
		}
		return tcs
	})
}

func TestOtSpanToOtlpSpan_SpanKind(t *testing.T) {
	testAndVerify("Test SpanKind", t, func(t *testing.T) []testCases {
		kinds := []struct {
			in  apitrace.SpanKind
			out tracepb.Span_SpanKind
		}{
			{
				in:  apitrace.SpanKindClient,
				out: tracepb.Span_CLIENT,
			},
			{
				in:  apitrace.SpanKindServer,
				out: tracepb.Span_SERVER,
			},
			{
				in:  apitrace.SpanKindProducer,
				out: tracepb.Span_PRODUCER,
			},
			{
				in:  apitrace.SpanKindConsumer,
				out: tracepb.Span_CONSUMER,
			},
			{
				in:  apitrace.SpanKindInternal,
				out: tracepb.Span_INTERNAL,
			},
			{
				in:  apitrace.SpanKindUnspecified,
				out: tracepb.Span_SPAN_KIND_UNSPECIFIED,
			},
		}

		tcs := make([]testCases, 0, len(kinds))
		for _, kind := range kinds {
			otSpan, otlpSpan := getSpan()
			otSpan.SpanKind = kind.in
			otlpSpan.Kind = kind.out
			tc := testCases{
				otSpan:   otSpan,
				otlpSpan: otlpSpan,
			}
			tcs = append(tcs, tc)
		}
		return tcs
	})
}

func TestOtSpanToOtlpSpan_Attribute(t *testing.T) {
	testAndVerify("Test SpanAttribute", t, func(t *testing.T) []testCases {
		attrInt := &commonpb.AttributeKeyValue{
			Key:      "commonInt",
			Type:     commonpb.AttributeKeyValue_INT,
			IntValue: 25,
		}
		attrInt64 := &commonpb.AttributeKeyValue{
			Key:      "commonInt64",
			Type:     commonpb.AttributeKeyValue_INT,
			IntValue: 12e9,
		}
		kinds := []struct {
			in  core.KeyValue
			out *commonpb.AttributeKeyValue
		}{
			{
				in:  core.Key("commonInt").Int(25),
				out: attrInt,
			},
			{
				in:  core.Key("commonInt").Uint(25),
				out: attrInt,
			},
			{
				in:  core.Key("commonInt").Int32(25),
				out: attrInt,
			},
			{
				in:  core.Key("commonInt").Uint32(25),
				out: attrInt,
			},
			{
				in:  core.Key("commonInt64").Int64(12e9),
				out: attrInt64,
			},
			{
				in:  core.Key("commonInt64").Uint64(12e9),
				out: attrInt64,
			},
			{
				in: core.Key("float32").Float32(3.598549),
				out: &commonpb.AttributeKeyValue{

					Key:         "float32",
					Type:        commonpb.AttributeKeyValue_DOUBLE,
					DoubleValue: 3.5985488891601562,
				},
			},
			{
				in: core.Key("float64").Float64(14.598549),
				out: &commonpb.AttributeKeyValue{

					Key:         "float64",
					Type:        commonpb.AttributeKeyValue_DOUBLE,
					DoubleValue: 14.598549,
				},
			},
			{
				in: core.Key("string").String("string"),
				out: &commonpb.AttributeKeyValue{

					Key:         "string",
					Type:        commonpb.AttributeKeyValue_STRING,
					StringValue: "string",
				},
			},
			{
				in: core.Key("bool").Bool(true),
				out: &commonpb.AttributeKeyValue{

					Key:       "bool",
					Type:      commonpb.AttributeKeyValue_BOOL,
					BoolValue: true,
				},
			},
		}

		tcs := make([]testCases, 0, len(kinds))
		for _, kind := range kinds {
			otSpan, otlpSpan := getSpan()
			otSpan.Attributes = []core.KeyValue{kind.in}
			otlpSpan.Attributes = []*commonpb.AttributeKeyValue{kind.out}
			tc := testCases{
				otSpan:   otSpan,
				otlpSpan: otlpSpan,
			}
			tcs = append(tcs, tc)
		}
		return tcs
	})
}

func getSpan() (*export.SpanData, *tracepb.Span) {
	startTime := time.Now()
	endTime := startTime.Add(10 * time.Second)

	otSpan := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID: core.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:  core.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		},
		SpanKind:     apitrace.SpanKindServer,
		ParentSpanID: core.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		Name:         "Test Span",
		StartTime:    startTime,
		EndTime:      endTime,
	}
	otlpSpan := &tracepb.Span{
		TraceId:           []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		SpanId:            []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		ParentSpanId:      []byte{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		Name:              "Test Span",
		Kind:              tracepb.Span_SERVER,
		StartTimeUnixnano: uint64(startTime.Nanosecond()),
		EndTimeUnixnano:   uint64(endTime.Nanosecond()),
		Status: &tracepb.Status{
			Code: 0,
		},
	}

	return otSpan, otlpSpan
}

func testAndVerify(name string, t *testing.T, f func(t *testing.T) []testCases) {
	// The goal of this test is to ensure that each
	// spanData is transformed and exported correctly!

	collector := runMockCol(t)
	defer func() {
		_ = collector.stop()
	}()

	exp, err := otlp.NewExporter(otlp.WithInsecure(),
		otlp.WithAddress(collector.address),
		otlp.WithReconnectionPeriod(50*time.Millisecond))
	if err != nil {
		t.Fatalf("Failed to create a new collector exporter: %v", err)
	}
	defer func() {
		_ = exp.Stop()
	}()

	// Give the background collector connection sometime to setup.
	<-time.After(20 * time.Millisecond)

	tcs := f(t)

	for _, tc := range tcs {
		exp.ExportSpans(context.Background(), []*export.SpanData{tc.otSpan})
	}

	_ = exp.Stop()
	_ = collector.stop()

	spans := collector.getSpans()
	gotCount := len(spans)
	wantCount := len(tcs)
	if gotCount != wantCount {
		t.Fatalf("%s: got %d spans, want %d spans", name, gotCount, wantCount)
	}
	for i, tc := range tcs {
		exp.ExportSpans(context.Background(), []*export.SpanData{tc.otSpan})
		if diff := cmp.Diff(spans[i], tc.otlpSpan, cmp.Comparer(proto.Equal)); diff != "" {
			t.Fatalf("%s transformed span differs %v\n", name, diff)
		}
	}
}
