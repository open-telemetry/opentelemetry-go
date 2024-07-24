// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package zipkin

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	zkmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv125 "go.opentelemetry.io/otel/semconv/v1.25.0"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func TestModelConversion(t *testing.T) {
	res := resource.NewSchemaless(
		semconv.ServiceName("model-test"),
		semconv.ServiceVersion("0.1.0"),
		attribute.Int64("resource-attr1", 42),
		attribute.IntSlice("resource-attr2", []int{0, 1, 2}),
	)

	inputBatch := tracetest.SpanStubs{
		// typical span data with UNSET status
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
				attribute.IntSlice("attr3", []int{0, 1, 2}),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Unset,
				Description: "",
			},
			Resource: res,
		},
		// typical span data with OK status
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
				attribute.IntSlice("attr3", []int{0, 1, 2}),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Ok,
				Description: "",
			},
			Resource: res,
		},
		// typical span data with ERROR status
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
				attribute.IntSlice("attr3", []int{0, 1, 2}),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data with no parent (same as typical, but has
		// invalid parent)
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data of unspecified kind
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindUnspecified,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data of internal kind
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindInternal,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data of client kind
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindClient,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
				attribute.String("peer.hostname", "test-peer-hostname"),
				attribute.String("network.peer.address", "1.2.3.4"),
				attribute.Int64("network.peer.port", 9876),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data of producer kind
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindProducer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data of consumer kind
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindConsumer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data with no events
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.Int64("attr1", 42),
				attribute.String("attr2", "bar"),
			},
			Events: nil,
			Status: tracesdk.Status{
				Code:        codes.Error,
				Description: "404, file not found",
			},
			Resource: res,
		},
		// span data with an "error" attribute set to "false"
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38},
			}),
			SpanKind:  trace.SpanKindServer,
			Name:      "foo",
			StartTime: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:   time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes: []attribute.KeyValue{
				attribute.String("error", "false"),
			},
			Events: []tracesdk.Event{
				{
					Time: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Name: "ev1",
					Attributes: []attribute.KeyValue{
						attribute.Int64("eventattr1", 123),
					},
				},
				{
					Time:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Name:       "ev2",
					Attributes: nil,
				},
			},
			Resource: res,
		},
	}.Snapshots()

	expectedOutputBatch := []zkmodel.SpanModel{
		// model for typical span data with UNSET status
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":           "42",
				"attr2":           "bar",
				"attr3":           "[0,1,2]",
				"service.name":    "model-test",
				"service.version": "0.1.0",
				"resource-attr1":  "42",
				"resource-attr2":  "[0,1,2]",
			},
		},
		// model for typical span data with OK status
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"attr3":            "[0,1,2]",
				"otel.status_code": "OK",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for typical span data with ERROR status
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"attr3":            "[0,1,2]",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data with no parent
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: nil,
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data of unspecified kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data of internal kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data of client kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "CLIENT",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: &zkmodel.Endpoint{
				IPv4: net.ParseIP("1.2.3.4"),
				Port: 9876,
			},
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":                "42",
				"attr2":                "bar",
				"network.peer.address": "1.2.3.4",
				"network.peer.port":    "9876",
				"peer.hostname":        "test-peer-hostname",
				"otel.status_code":     "ERROR",
				"error":                "404, file not found",
				"service.name":         "model-test",
				"service.version":      "0.1.0",
				"resource-attr1":       "42",
				"resource-attr2":       "[0,1,2]",
			},
		},
		// model for span data of producer kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "PRODUCER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data of consumer kind
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "CONSUMER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data with no events
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations:    nil,
			Tags: map[string]string{
				"attr1":            "42",
				"attr2":            "bar",
				"otel.status_code": "ERROR",
				"error":            "404, file not found",
				"service.name":     "model-test",
				"service.version":  "0.1.0",
				"resource-attr1":   "42",
				"resource-attr2":   "[0,1,2]",
			},
		},
		// model for span data with an "error" attribute set to "false"
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: zkmodelIDPtr(0x3f3e3d3c3b3a3938),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "model-test",
			},
			RemoteEndpoint: nil,
			Annotations: []zkmodel.Annotation{
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 30, 0, time.UTC),
					Value:     `ev1: {"eventattr1":123}`,
				},
				{
					Timestamp: time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
					Value:     "ev2",
				},
			},
			Tags: map[string]string{
				"service.name":    "model-test",
				"service.version": "0.1.0",
				"resource-attr1":  "42",
				"resource-attr2":  "[0,1,2]",
			}, // only resource tags should be included
		},
	}
	gottenOutputBatch := SpanModels(inputBatch)
	require.Equal(t, expectedOutputBatch, gottenOutputBatch)
}

func zkmodelIDPtr(n uint64) *zkmodel.ID {
	id := zkmodel.ID(n)
	return &id
}

func TestTagsTransformation(t *testing.T) {
	keyValue := "value"
	doubleValue := 123.456
	uintValue := int64(123)
	statusMessage := "this is a problem"
	instrLibName := "instrumentation-library"
	instrLibVersion := "semver:1.0.0"

	tests := []struct {
		name string
		data tracetest.SpanStub
		want map[string]string
	}{
		{
			name: "attributes",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{
					attribute.String("key", keyValue),
					attribute.Float64("double", doubleValue),
					attribute.Int64("uint", uintValue),
					attribute.Bool("ok", true),
				},
			},
			want: map[string]string{
				"double": fmt.Sprint(doubleValue),
				"key":    keyValue,
				"ok":     "true",
				"uint":   strconv.FormatInt(uintValue, 10),
			},
		},
		{
			name: "no attributes",
			data: tracetest.SpanStub{},
			want: nil,
		},
		{
			name: "omit-noerror",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{
					attribute.Bool("error", false),
				},
			},
			want: nil,
		},
		{
			name: "statusCode UNSET",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{
					attribute.String("key", keyValue),
				},
				Status: tracesdk.Status{
					Code:        codes.Unset,
					Description: "",
				},
			},
			want: map[string]string{
				"key": keyValue,
			},
		},
		{
			name: "statusCode OK",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{
					attribute.String("key", keyValue),
					attribute.Bool("ok", true),
				},
				Status: tracesdk.Status{
					Code:        codes.Ok,
					Description: "",
				},
			},
			want: map[string]string{
				"key":              keyValue,
				"ok":               "true",
				"otel.status_code": "OK",
			},
		},
		{
			name: "statusCode ERROR",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{
					attribute.String("key", keyValue),
					attribute.Bool("error", true),
				},
				Status: tracesdk.Status{
					Code:        codes.Error,
					Description: statusMessage,
				},
			},
			want: map[string]string{
				"error":            statusMessage,
				"key":              keyValue,
				"otel.status_code": "ERROR",
			},
		},
		{
			name: "instrLib-empty",
			data: tracetest.SpanStub{
				InstrumentationScope: instrumentation.Scope{},
			},
			want: nil,
		},
		{
			name: "instrLib-noversion",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{},
				InstrumentationScope: instrumentation.Scope{
					Name: instrLibName,
				},
			},
			want: map[string]string{
				"otel.scope.name": instrLibName,
			},
		},
		{
			name: "instrLib-with-version",
			data: tracetest.SpanStub{
				Attributes: []attribute.KeyValue{},
				InstrumentationScope: instrumentation.Scope{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: map[string]string{
				"otel.scope.name":    instrLibName,
				"otel.scope.version": instrLibVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toZipkinTags(tt.data.Snapshot())
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff%v", diff)
			}
		})
	}
}

func TestRemoteEndpointTransformation(t *testing.T) {
	tests := []struct {
		name string
		data tracetest.SpanStub
		want *zkmodel.Endpoint
	}{
		{
			name: "nil-not-applicable",
			data: tracetest.SpanStub{
				SpanKind:   trace.SpanKindClient,
				Attributes: []attribute.KeyValue{},
			},
			want: nil,
		},
		{
			name: "nil-not-found",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindConsumer,
				Attributes: []attribute.KeyValue{
					attribute.String("attr", "test"),
				},
			},
			want: nil,
		},
		{
			name: "peer.service-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					semconv.PeerService("peer-service-test"),
					semconv.ServerAddress("server-address-test"),
					semconv.NetworkPeerAddress("10.1.2.80"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "peer-service-test",
			},
		},
		{
			name: "server.address-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					semconv.ServerAddress("server-address-test"),
					attribute.String("net.peer.name", "net-peer-name-test"),
					semconv.NetworkPeerAddress("10.1.2.80"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "server-address-test",
			},
		},
		{
			name: "net.peer.name-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					attribute.String("net.peer.name", "net-peer-name-test"),
					semconv.NetworkPeerAddress("10.1.2.80"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "net-peer-name-test",
			},
		},
		{
			name: "network.peer.address-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					semconv.NetworkPeerAddress("10.1.2.80"),
					semconv125.DBName("db-name-test"),
					attribute.String("server.socket.domain", "server-socket-domain-test"),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("10.1.2.80"),
			},
		},
		{
			name: "server.socket.domain-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					semconv125.DBName("db-name-test"),
					attribute.String("server.socket.domain", "server-socket-domain-test"),
					attribute.String("server.socket.address", "10.2.3.4"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "server-socket-domain-test",
			},
		},
		{
			name: "server.socket.address-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					semconv125.DBName("db-name-test"),
					attribute.String("net.sock.peer.name", "server-socket-domain-test"),
					attribute.String("server.socket.address", "10.2.3.4"),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("10.2.3.4"),
			},
		},
		{
			name: "net.sock.peer.name-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					semconv125.DBName("db-name-test"),
					attribute.String("net.sock.peer.name", "net-sock-peer-name-test"),
					attribute.String("net.sock.peer.addr", "10.4.8.12"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "net-sock-peer-name-test",
			},
		},
		{
			name: "net.sock.peer.addr-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					semconv125.DBName("db-name-test"),
					attribute.String("net.sock.peer.addr", "10.4.8.12"),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("10.4.8.12"),
			},
		},
		{
			name: "peer.hostname-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerHostname.String("peer-hostname-test"),
					keyPeerAddress.String("peer-address-test"),
					semconv125.DBName("http-host-test"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "peer-hostname-test",
			},
		},
		{
			name: "peer.address-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					keyPeerAddress.String("peer-address-test"),
					semconv125.DBName("http-host-test"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "peer-address-test",
			},
		},
		{
			name: "db.name-rank",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					attribute.String("foo", "bar"),
					semconv125.DBName("db-name-test"),
				},
			},
			want: &zkmodel.Endpoint{
				ServiceName: "db-name-test",
			},
		},
		{
			name: "network.peer.address-invalid-ip",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					semconv.NetworkPeerAddress("INVALID"),
				},
			},
			want: nil,
		},
		{
			name: "network.peer.address-ipv6-no-port",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					semconv.NetworkPeerAddress("0:0:1:5ee:bad:c0de:0:0"),
				},
			},
			want: &zkmodel.Endpoint{
				IPv6: net.ParseIP("0:0:1:5ee:bad:c0de:0:0"),
			},
		},
		{
			name: "network.peer.address-ipv4-port",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					semconv.NetworkPeerAddress("1.2.3.4"),
					semconv.NetworkPeerPort(9876),
					attribute.Int("server.socket.port", 5432),
					attribute.Int("net.sock.peer.port", 2345),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("1.2.3.4"),
				Port: 9876,
			},
		},
		{
			name: "server.socket.address-ipv4-port",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					attribute.String("server.socket.address", "1.2.3.4"),
					semconv.NetworkPeerPort(9876),
					attribute.Int("server.socket.port", 5432),
					attribute.Int("net.sock.peer.port", 2345),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("1.2.3.4"),
				Port: 5432,
			},
		},
		{
			name: "net.sock.peer.addr-ipv4-port",
			data: tracetest.SpanStub{
				SpanKind: trace.SpanKindProducer,
				Attributes: []attribute.KeyValue{
					attribute.String("net.sock.peer.addr", "1.2.3.4"),
					semconv.NetworkPeerPort(9876),
					attribute.Int("server.socket.port", 5432),
					attribute.Int("net.sock.peer.port", 2345),
				},
			},
			want: &zkmodel.Endpoint{
				IPv4: net.ParseIP("1.2.3.4"),
				Port: 2345,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toZipkinRemoteEndpoint(tt.data.Snapshot())
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff%v", diff)
			}
		})
	}
}

func TestServiceName(t *testing.T) {
	attrs := []attribute.KeyValue{}
	assert.Equal(t, defaultServiceName, getServiceName(attrs))

	attrs = append(attrs, attribute.String("test_key", "test_value"))
	assert.Equal(t, defaultServiceName, getServiceName(attrs))

	attrs = append(attrs, semconv.ServiceName("my_service"))
	assert.Equal(t, "my_service", getServiceName(attrs))
}
