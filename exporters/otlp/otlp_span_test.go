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

package otlp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	commonpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1"
	resourcepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/resource/v1"
	tracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/trace/v1"
	"go.opentelemetry.io/otel/label"

	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

type traceServiceClientStub struct {
	rs []tracepb.ResourceSpans
}

func (t *traceServiceClientStub) Export(ctx context.Context, in *coltracepb.ExportTraceServiceRequest, opts ...grpc.CallOption) (*coltracepb.ExportTraceServiceResponse, error) {
	for _, rs := range in.GetResourceSpans() {
		if rs == nil {
			continue
		}
		t.rs = append(t.rs, *rs)
	}
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

func (t *traceServiceClientStub) ResourceSpans() []tracepb.ResourceSpans {
	return t.rs
}

func (t *traceServiceClientStub) Reset() {
	t.rs = nil
}

func TestExportSpans(t *testing.T) {
	tsc := &traceServiceClientStub{}
	exp := NewUnstartedExporter()
	exp.traceExporter = tsc
	exp.started = true

	// March 31, 2020 5:01:26 1234nanos (UTC)
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)

	for _, test := range []struct {
		sd   []*tracesdk.SpanData
		want []tracepb.ResourceSpans
	}{
		{
			[]*tracesdk.SpanData(nil),
			[]tracepb.ResourceSpans(nil),
		},
		{
			[]*tracesdk.SpanData{},
			[]tracepb.ResourceSpans(nil),
		},
		{
			[]*tracesdk.SpanData{
				{
					SpanContext: otel.SpanContext{
						TraceID:    otel.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
						SpanID:     otel.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: byte(1),
					},
					SpanKind:  otel.SpanKindServer,
					Name:      "parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []label.KeyValue{
						label.String("user", "alice"),
						label.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.New(label.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: otel.SpanContext{
						TraceID:    otel.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
						SpanID:     otel.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: byte(1),
					},
					SpanKind:  otel.SpanKindServer,
					Name:      "secondary parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []label.KeyValue{
						label.String("user", "alice"),
						label.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.New(label.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-b",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: otel.SpanContext{
						TraceID:    otel.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
						SpanID:     otel.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 2}),
						TraceFlags: byte(1),
					},
					ParentSpanID: otel.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
					SpanKind:     otel.SpanKindInternal,
					Name:         "internal process",
					StartTime:    startTime,
					EndTime:      endTime,
					Attributes: []label.KeyValue{
						label.String("user", "alice"),
						label.Bool("authenticated", true),
					},
					StatusCode:    codes.Ok,
					StatusMessage: "Ok",
					Resource:      resource.New(label.String("instance", "tester-a")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v0.1.0",
					},
				},
				{
					SpanContext: otel.SpanContext{
						TraceID:    otel.TraceID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
						SpanID:     otel.SpanID([8]byte{0, 0, 0, 0, 0, 0, 0, 1}),
						TraceFlags: byte(1),
					},
					SpanKind:  otel.SpanKindServer,
					Name:      "parent process",
					StartTime: startTime,
					EndTime:   endTime,
					Attributes: []label.KeyValue{
						label.String("user", "bob"),
						label.Bool("authenticated", false),
					},
					StatusCode:    codes.Error,
					StatusMessage: "Unauthenticated",
					Resource:      resource.New(label.String("instance", "tester-b")),
					InstrumentationLibrary: instrumentation.Library{
						Name:    "lib-a",
						Version: "v1.1.0",
					},
				},
			},
			[]tracepb.ResourceSpans{
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
									Kind:              tracepb.Span_SERVER,
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
										Code:    tracepb.Status_Ok,
										Message: "Ok",
									},
								},
								{
									TraceId:           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
									SpanId:            []byte{0, 0, 0, 0, 0, 0, 0, 2},
									ParentSpanId:      []byte{0, 0, 0, 0, 0, 0, 0, 1},
									Name:              "internal process",
									Kind:              tracepb.Span_INTERNAL,
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
										Code:    tracepb.Status_Ok,
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
									Kind:              tracepb.Span_SERVER,
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
										Code:    tracepb.Status_Ok,
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
									Kind:              tracepb.Span_SERVER,
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
										Code:    tracepb.Status_UnknownError,
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
		tsc.Reset()
		assert.NoError(t, exp.ExportSpans(context.Background(), test.sd))
		assert.ElementsMatch(t, test.want, tsc.ResourceSpans())
	}
}
