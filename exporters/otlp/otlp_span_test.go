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

package otlp_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestExportSpans(t *testing.T) {
	exp, driver := newExporter(t)

	// March 31, 2020 5:01:26 1234nanos (UTC)
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)

	for _, test := range []struct {
		sd   []*tracesdk.SpanSnapshot
		want []*tracepb.ResourceSpans
	}{
		{
			[]*tracesdk.SpanSnapshot(nil),
			[]*tracepb.ResourceSpans(nil),
		},
		{
			[]*tracesdk.SpanSnapshot{},
			[]*tracepb.ResourceSpans(nil),
		},
		{
			[]*tracesdk.SpanSnapshot{
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    trace.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
						SpanID:     trace.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: trace.FlagsSampled,
					}),
					SpanKind:  trace.SpanKindServer,
					Name:      "parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []attribute.KeyValue{
						attribute.String("user", "alice"),
						attribute.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.NewWithAttributes(attribute.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    trace.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
						SpanID:     trace.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: trace.FlagsSampled,
					}),
					SpanKind:  trace.SpanKindServer,
					Name:      "secondary parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []attribute.KeyValue{
						attribute.String("user", "alice"),
						attribute.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.NewWithAttributes(attribute.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-b",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    trace.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
						SpanID:     trace.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 2}),
						TraceFlags: trace.FlagsSampled,
					}),
					Parent: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    trace.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
						SpanID:     trace.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: trace.FlagsSampled,
					}),
					SpanKind:  trace.SpanKindInternal,
					Name:      "internal process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []attribute.KeyValue{
						attribute.String("user", "alice"),
						attribute.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.NewWithAttributes(attribute.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    trace.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
						SpanID:     trace.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: trace.FlagsSampled,
					}),
					SpanKind:  trace.SpanKindServer,
					Name:      "parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []attribute.KeyValue{
						attribute.String("user", "bob"),
						attribute.Bool("authenticated", false),
					},
					StatusCode:    codes.Error,
					StatusMessage: "Unauthenticated",
					Resource:      resource.NewWithAttributes(attribute.String("instance", "tester-b")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v1.1.0",
					},
				},
			},
			[]*tracepb.ResourceSpans{
				{
					Resource: &resourcepb.Resource{
						Attributes: []*commonpb.KeyValue{
							{
								Key: "instance",
								Value: &commonpb.AnyValue{
									Value: &commonpb.AnyValue_StringValue{
										StringValue: "tester-a",
									},
								},
							},
						},
					},
					InstrumentationLibrarySpans: []*tracepb.InstrumentationLibrarySpans{
						{
							InstrumentationLibrary: &commonpb.InstrumentationLibrary{
								Name:    "lib-a",
								Version: "v0.1.0",
							},
							Spans: []*tracepb.Span{
								{
									TraceId:           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
									SpanId:            []byte{0, 0, 0, 0, 0, 0, 0, 1},
									Name:              "parent process",
									Kind:              tracepb.Span_SPAN_KIND_SERVER,
									StartTimeUnixNano: uint64(startTime.UnixNano()),
									EndTimeUnixNano:   uint64(endTime.UnixNano()),
									Attributes: []*commonpb.KeyValue{
										{
											Key: "user",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: "alice",
												},
											},
										},
										{
											Key: "authenticated",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BoolValue{
													BoolValue: true,
												},
											},
										},
									},
									Status: &tracepb.Status{
										Code:    tracepb.Status_STATUS_CODE_OK,
										Message: "Ok",
									},
								},
								{
									TraceId:           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
									SpanId:            []byte{0, 0, 0, 0, 0, 0, 0, 2},
									ParentSpanId:      []byte{0, 0, 0, 0, 0, 0, 0, 1},
									Name:              "internal process",
									Kind:              tracepb.Span_SPAN_KIND_INTERNAL,
									StartTimeUnixNano: uint64(startTime.UnixNano()),
									EndTimeUnixNano:   uint64(endTime.UnixNano()),
									Attributes: []*commonpb.KeyValue{
										{
											Key: "user",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: "alice",
												},
											},
										},
										{
											Key: "authenticated",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BoolValue{
													BoolValue: true,
												},
											},
										},
									},
									Status: &tracepb.Status{
										Code:    tracepb.Status_STATUS_CODE_OK,
										Message: "Ok",
									},
								},
							},
						},
						{
							InstrumentationLibrary: &commonpb.InstrumentationLibrary{
								Name:    "lib-b",
								Version: "v0.1.0",
							},
							Spans: []*tracepb.Span{
								{
									TraceId:           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
									SpanId:            []byte{0, 0, 0, 0, 0, 0, 0, 1},
									Name:              "secondary parent process",
									Kind:              tracepb.Span_SPAN_KIND_SERVER,
									StartTimeUnixNano: uint64(startTime.UnixNano()),
									EndTimeUnixNano:   uint64(endTime.UnixNano()),
									Attributes: []*commonpb.KeyValue{
										{
											Key: "user",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: "alice",
												},
											},
										},
										{
											Key: "authenticated",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BoolValue{
													BoolValue: true,
												},
											},
										},
									},
									Status: &tracepb.Status{
										Code:    tracepb.Status_STATUS_CODE_OK,
										Message: "Ok",
									},
								},
							},
						},
					},
				},
				{
					Resource: &resourcepb.Resource{
						Attributes: []*commonpb.KeyValue{
							{
								Key: "instance",
								Value: &commonpb.AnyValue{
									Value: &commonpb.AnyValue_StringValue{
										StringValue: "tester-b",
									},
								},
							},
						},
					},
					InstrumentationLibrarySpans: []*tracepb.InstrumentationLibrarySpans{
						{
							InstrumentationLibrary: &commonpb.InstrumentationLibrary{
								Name:    "lib-a",
								Version: "v1.1.0",
							},
							Spans: []*tracepb.Span{
								{
									TraceId:           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
									SpanId:            []byte{0, 0, 0, 0, 0, 0, 0, 1},
									Name:              "parent process",
									Kind:              tracepb.Span_SPAN_KIND_SERVER,
									StartTimeUnixNano: uint64(startTime.UnixNano()),
									EndTimeUnixNano:   uint64(endTime.UnixNano()),
									Attributes: []*commonpb.KeyValue{
										{
											Key: "user",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: "bob",
												},
											},
										},
										{
											Key: "authenticated",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BoolValue{
													BoolValue: false,
												},
											},
										},
									},
									Status: &tracepb.Status{
										Code:    tracepb.Status_STATUS_CODE_ERROR,
										Message: "Unauthenticated",
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		driver.Reset()
		assert.NoError(t, exp.ExportSpans(context.Background(), test.sd))
		assert.ElementsMatch(t, test.want, driver.rs)
	}
}
