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

package jaeger

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/support/bundler"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	collectorEndpoint = "http://localhost:14268/api/traces"
)

func TestInstallNewPipeline(t *testing.T) {
	tp, err := InstallNewPipeline(WithCollectorEndpoint(WithEndpoint(collectorEndpoint)))
	require.NoError(t, err)
	// Ensure InstallNewPipeline sets the global TracerProvider. By default
	// the global tracer provider will be a NoOp implementation, this checks
	// if that has been overwritten.
	assert.IsType(t, tp, otel.GetTracerProvider())
}

func TestNewExportPipelinePassthroughError(t *testing.T) {
	for _, testcase := range []struct {
		name    string
		failing bool
		epo     EndpointOption
	}{
		{
			name:    "failing underlying NewRawExporter",
			failing: true,
			epo: func() (batchUploader, error) {
				return nil, errors.New("error")
			},
		},
		{
			name: "with default agent endpoint",
			epo:  WithAgentEndpoint(),
		},
		{
			name: "with collector endpoint",
			epo:  WithCollectorEndpoint(WithEndpoint(collectorEndpoint)),
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			_, err := NewExportPipeline(testcase.epo)
			if testcase.failing {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewRawExporter(t *testing.T) {
	testCases := []struct {
		name                                                           string
		endpoint                                                       EndpointOption
		options                                                        []Option
		expectedServiceName                                            string
		expectedTagsLen, expectedBufferMaxCount, expectedBatchMaxCount int
	}{
		{
			name:                   "default exporter",
			endpoint:               WithCollectorEndpoint(),
			expectedServiceName:    "unknown_service",
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:                   "default exporter with agent endpoint",
			endpoint:               WithAgentEndpoint(),
			expectedServiceName:    "unknown_service",
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:     "with buffer and batch max count",
			endpoint: WithCollectorEndpoint(WithEndpoint(collectorEndpoint)),
			options: []Option{
				WithBufferMaxCount(99),
				WithBatchMaxCount(99),
			},
			expectedServiceName:    "unknown_service",
			expectedBufferMaxCount: 99,
			expectedBatchMaxCount:  99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exp, err := NewRawExporter(
				tc.endpoint,
				tc.options...,
			)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBufferMaxCount, exp.bundler.BufferedByteLimit)
			assert.Equal(t, tc.expectedBatchMaxCount, exp.bundler.BundleCountThreshold)
			assert.NotEmpty(t, exp.defaultServiceName)
			assert.True(t, strings.HasPrefix(exp.defaultServiceName, tc.expectedServiceName))
		})
	}
}

func TestNewRawExporterUseEnvVarIfOptionUnset(t *testing.T) {
	// Record and restore env
	envStore := ottest.NewEnvStore()
	envStore.Record(envEndpoint)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// If the user sets the environment variable OTEL_EXPORTER_JAEGER_ENDPOINT, endpoint will always get a value.
	require.NoError(t, os.Unsetenv(envEndpoint))
	_, err := NewRawExporter(
		WithCollectorEndpoint(),
	)

	assert.NoError(t, err)
}

type testCollectorEndpoint struct {
	batchesUploaded []*gen.Batch
}

func (c *testCollectorEndpoint) upload(_ context.Context, batch *gen.Batch) error {
	c.batchesUploaded = append(c.batchesUploaded, batch)
	return nil
}

var _ batchUploader = (*testCollectorEndpoint)(nil)

func withTestCollectorEndpoint() func() (batchUploader, error) {
	return func() (batchUploader, error) {
		return &testCollectorEndpoint{}, nil
	}
}

func withTestCollectorEndpointInjected(ce *testCollectorEndpoint) func() (batchUploader, error) {
	return func() (batchUploader, error) {
		return ce, nil
	}
}

func TestExporterExportSpan(t *testing.T) {
	const (
		serviceName = "test-service"
		tagKey      = "key"
		tagVal      = "val"
	)

	testCollector := &testCollectorEndpoint{}
	exp, err := NewRawExporter(withTestCollectorEndpointInjected(testCollector))
	require.NoError(t, err)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String(tagKey, tagVal),
		)),
	)

	tracer := tp.Tracer("test-tracer")

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		_, span := tracer.Start(ctx, fmt.Sprintf("test-span-%d", i))
		span.End()
		assert.True(t, span.SpanContext().IsValid())
	}

	require.NoError(t, tp.Shutdown(ctx))

	batchesUploaded := testCollector.batchesUploaded
	require.Len(t, batchesUploaded, 1)
	uploadedBatch := batchesUploaded[0]
	assert.Equal(t, serviceName, uploadedBatch.GetProcess().GetServiceName())
	assert.Len(t, uploadedBatch.GetSpans(), 3)

	require.Len(t, uploadedBatch.GetProcess().GetTags(), 1)
	assert.Equal(t, tagKey, uploadedBatch.GetProcess().GetTags()[0].GetKey())
	assert.Equal(t, tagVal, uploadedBatch.GetProcess().GetTags()[0].GetVStr())
}

func Test_spanSnapshotToThrift(t *testing.T) {
	now := time.Now()
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")
	parentSpanID, _ := trace.SpanIDFromHex("0807060504030201")

	linkTraceID, _ := trace.TraceIDFromHex("0102030405060709090a0b0c0d0e0f11")
	linkSpanID, _ := trace.SpanIDFromHex("0102030405060709")

	eventNameValue := "event-test"
	eventDropped := int64(10)
	keyValue := "value"
	statusCodeValue := int64(1)
	doubleValue := 123.456
	intValue := int64(123)
	boolTrue := true
	arrValue := "[0,1,2,3]"
	statusMessage := "this is a problem"
	spanKind := "client"
	rv1 := "rv11"
	rv2 := int64(5)
	instrLibName := "instrumentation-library"
	instrLibVersion := "semver:1.0.0"

	tests := []struct {
		name string
		data *sdktrace.SpanSnapshot
		want *gen.Span
	}{
		{
			name: "no status description",
			data: &sdktrace.SpanSnapshot{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
				}),
				Name:       "/foo",
				StartTime:  now,
				EndTime:    now,
				StatusCode: codes.Error,
				SpanKind:   trace.SpanKindClient,
				InstrumentationLibrary: instrumentation.Library{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: &gen.Span{
				TraceIdLow:    651345242494996240,
				TraceIdHigh:   72623859790382856,
				SpanId:        72623859790382856,
				OperationName: "/foo",
				StartTime:     now.UnixNano() / 1000,
				Duration:      0,
				Tags: []*gen.Tag{
					{Key: keyError, VType: gen.TagType_BOOL, VBool: &boolTrue},
					{Key: keyInstrumentationLibraryName, VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: keyInstrumentationLibraryVersion, VType: gen.TagType_STRING, VStr: &instrLibVersion},
					{Key: keyStatusCode, VType: gen.TagType_LONG, VLong: &statusCodeValue},
					// Should not have a status message because it was unset
					{Key: keySpanKind, VType: gen.TagType_STRING, VStr: &spanKind},
				},
			},
		},
		{
			name: "no parent",
			data: &sdktrace.SpanSnapshot{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
				}),
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Links: []trace.Link{
					{
						SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
							TraceID: linkTraceID,
							SpanID:  linkSpanID,
						}),
					},
				},
				Attributes: []attribute.KeyValue{
					attribute.String("key", keyValue),
					attribute.Float64("double", doubleValue),
					attribute.Int64("int", intValue),
				},
				MessageEvents: []trace.Event{
					{
						Name:                  eventNameValue,
						Attributes:            []attribute.KeyValue{attribute.String("k1", keyValue)},
						DroppedAttributeCount: int(eventDropped),
						Time:                  now,
					},
				},
				StatusCode:    codes.Error,
				StatusMessage: statusMessage,
				SpanKind:      trace.SpanKindClient,
				InstrumentationLibrary: instrumentation.Library{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: &gen.Span{
				TraceIdLow:    651345242494996240,
				TraceIdHigh:   72623859790382856,
				SpanId:        72623859790382856,
				OperationName: "/foo",
				StartTime:     now.UnixNano() / 1000,
				Duration:      0,
				Tags: []*gen.Tag{
					{Key: "double", VType: gen.TagType_DOUBLE, VDouble: &doubleValue},
					{Key: "key", VType: gen.TagType_STRING, VStr: &keyValue},
					{Key: "int", VType: gen.TagType_LONG, VLong: &intValue},
					{Key: keyError, VType: gen.TagType_BOOL, VBool: &boolTrue},
					{Key: keyInstrumentationLibraryName, VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: keyInstrumentationLibraryVersion, VType: gen.TagType_STRING, VStr: &instrLibVersion},
					{Key: keyStatusCode, VType: gen.TagType_LONG, VLong: &statusCodeValue},
					{Key: keyStatusMessage, VType: gen.TagType_STRING, VStr: &statusMessage},
					{Key: keySpanKind, VType: gen.TagType_STRING, VStr: &spanKind},
				},
				References: []*gen.SpanRef{
					{
						RefType:     gen.SpanRefType_FOLLOWS_FROM,
						TraceIdHigh: int64(binary.BigEndian.Uint64(linkTraceID[0:8])),
						TraceIdLow:  int64(binary.BigEndian.Uint64(linkTraceID[8:16])),
						SpanId:      int64(binary.BigEndian.Uint64(linkSpanID[:])),
					},
				},
				Logs: []*gen.Log{
					{
						Timestamp: now.UnixNano() / 1000,
						Fields: []*gen.Tag{
							{
								Key:   keyEventName,
								VStr:  &eventNameValue,
								VType: gen.TagType_STRING,
							},
							{
								Key:   "k1",
								VStr:  &keyValue,
								VType: gen.TagType_STRING,
							},
							{
								Key:   keyDroppedAttributeCount,
								VLong: &eventDropped,
								VType: gen.TagType_LONG,
							},
						},
					},
				},
			},
		},
		{
			name: "with parent",
			data: &sdktrace.SpanSnapshot{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
				}),
				Parent: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  parentSpanID,
				}),
				Links: []trace.Link{
					{
						SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
							TraceID: linkTraceID,
							SpanID:  linkSpanID,
						}),
					},
				},
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Attributes: []attribute.KeyValue{
					attribute.Array("arr", []int{0, 1, 2, 3}),
				},
				StatusCode:    codes.Unset,
				StatusMessage: statusMessage,
				SpanKind:      trace.SpanKindInternal,
				InstrumentationLibrary: instrumentation.Library{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: &gen.Span{
				TraceIdLow:    651345242494996240,
				TraceIdHigh:   72623859790382856,
				SpanId:        72623859790382856,
				ParentSpanId:  578437695752307201,
				OperationName: "/foo",
				StartTime:     now.UnixNano() / 1000,
				Duration:      0,
				Tags: []*gen.Tag{
					// status code, message and span kind should NOT be populated
					{Key: "arr", VType: gen.TagType_STRING, VStr: &arrValue},
					{Key: keyInstrumentationLibraryName, VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: keyInstrumentationLibraryVersion, VType: gen.TagType_STRING, VStr: &instrLibVersion},
				},
				References: []*gen.SpanRef{
					{
						RefType:     gen.SpanRefType_FOLLOWS_FROM,
						TraceIdHigh: int64(binary.BigEndian.Uint64(linkTraceID[0:8])),
						TraceIdLow:  int64(binary.BigEndian.Uint64(linkTraceID[8:16])),
						SpanId:      int64(binary.BigEndian.Uint64(linkSpanID[:])),
					},
				},
			},
		},
		{
			name: "resources do not affect the tags",
			data: &sdktrace.SpanSnapshot{
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
				}),
				Parent: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  parentSpanID,
				}),
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Resource: resource.NewWithAttributes(
					attribute.String("rk1", rv1),
					attribute.Int64("rk2", rv2),
					semconv.ServiceNameKey.String("service name"),
				),
				StatusCode:    codes.Unset,
				StatusMessage: statusMessage,
				SpanKind:      trace.SpanKindInternal,
				InstrumentationLibrary: instrumentation.Library{
					Name:    instrLibName,
					Version: instrLibVersion,
				},
			},
			want: &gen.Span{
				TraceIdLow:    651345242494996240,
				TraceIdHigh:   72623859790382856,
				SpanId:        72623859790382856,
				ParentSpanId:  578437695752307201,
				OperationName: "/foo",
				StartTime:     now.UnixNano() / 1000,
				Duration:      0,
				Tags: []*gen.Tag{
					{Key: keyInstrumentationLibraryName, VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: keyInstrumentationLibraryVersion, VType: gen.TagType_STRING, VStr: &instrLibVersion},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spanSnapshotToThrift(tt.data)
			sort.Slice(got.Tags, func(i, j int) bool {
				return got.Tags[i].Key < got.Tags[j].Key
			})
			sort.Slice(tt.want.Tags, func(i, j int) bool {
				return tt.want.Tags[i].Key < tt.want.Tags[j].Key
			})
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff%v", diff)
			}
		})
	}
}

func TestExporterShutdownHonorsCancel(t *testing.T) {
	orig := flush
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// Do this after the parent context is canceled to avoid a race.
	defer func() {
		<-ctx.Done()
		flush = orig
	}()
	defer cancel()
	flush = func(e *Exporter) {
		<-ctx.Done()
	}

	e, err := NewRawExporter(withTestCollectorEndpoint())
	require.NoError(t, err)
	innerCtx, innerCancel := context.WithCancel(ctx)
	go innerCancel()
	assert.Errorf(t, e.Shutdown(innerCtx), context.Canceled.Error())
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	orig := flush
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// Do this after the parent context is canceled to avoid a race.
	defer func() {
		<-ctx.Done()
		flush = orig
	}()
	defer cancel()
	flush = func(e *Exporter) {
		<-ctx.Done()
	}

	e, err := NewRawExporter(withTestCollectorEndpoint())
	require.NoError(t, err)
	innerCtx, innerCancel := context.WithTimeout(ctx, time.Microsecond*10)
	assert.Errorf(t, e.Shutdown(innerCtx), context.DeadlineExceeded.Error())
	innerCancel()
}

func TestErrorOnExportShutdownExporter(t *testing.T) {
	e, err := NewRawExporter(withTestCollectorEndpoint())
	require.NoError(t, err)
	assert.NoError(t, e.Shutdown(context.Background()))
	assert.NoError(t, e.ExportSpans(context.Background(), nil))
}

func TestExporterExportSpansHonorsCancel(t *testing.T) {
	e, err := NewRawExporter(withTestCollectorEndpoint())
	require.NoError(t, err)
	now := time.Now()
	ss := []*sdktrace.SpanSnapshot{
		{
			Name: "s1",
			Resource: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("name"),
				attribute.Key("r1").String("v1"),
			),
			StartTime: now,
			EndTime:   now,
		},
		{
			Name: "s2",
			Resource: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("name"),
				attribute.Key("r2").String("v2"),
			),
			StartTime: now,
			EndTime:   now,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.EqualError(t, e.ExportSpans(ctx, ss), context.Canceled.Error())
}

func TestExporterExportSpansHonorsTimeout(t *testing.T) {
	e, err := NewRawExporter(withTestCollectorEndpoint())
	require.NoError(t, err)
	now := time.Now()
	ss := []*sdktrace.SpanSnapshot{
		{
			Name: "s1",
			Resource: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("name"),
				attribute.Key("r1").String("v1"),
			),
			StartTime: now,
			EndTime:   now,
		},
		{
			Name: "s2",
			Resource: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("name"),
				attribute.Key("r2").String("v2"),
			),
			StartTime: now,
			EndTime:   now,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	assert.EqualError(t, e.ExportSpans(ctx, ss), context.DeadlineExceeded.Error())
}

func TestJaegerBatchList(t *testing.T) {
	newString := func(value string) *string {
		return &value
	}
	spanKind := "unspecified"
	now := time.Now()

	testCases := []struct {
		name               string
		spanSnapshotList   []*sdktrace.SpanSnapshot
		defaultServiceName string
		expectedBatchList  []*gen.Batch
	}{
		{
			name:              "no span shots",
			spanSnapshotList:  nil,
			expectedBatchList: nil,
		},
		{
			name: "span's snapshot contains nil span",
			spanSnapshotList: []*sdktrace.SpanSnapshot{
				{
					Name: "s1",
					Resource: resource.NewWithAttributes(
						semconv.ServiceNameKey.String("name"),
						attribute.Key("r1").String("v1"),
					),
					StartTime: now,
					EndTime:   now,
				},
				nil,
			},
			expectedBatchList: []*gen.Batch{
				{
					Process: &gen.Process{
						ServiceName: "name",
						Tags: []*gen.Tag{
							{Key: "r1", VType: gen.TagType_STRING, VStr: newString("v1")},
						},
					},
					Spans: []*gen.Span{
						{
							OperationName: "s1",
							Tags: []*gen.Tag{
								{Key: keySpanKind, VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
					},
				},
			},
		},
		{
			name: "merge spans that have the same resources",
			spanSnapshotList: []*sdktrace.SpanSnapshot{
				{
					Name: "s1",
					Resource: resource.NewWithAttributes(
						semconv.ServiceNameKey.String("name"),
						attribute.Key("r1").String("v1"),
					),
					StartTime: now,
					EndTime:   now,
				},
				{
					Name: "s2",
					Resource: resource.NewWithAttributes(
						semconv.ServiceNameKey.String("name"),
						attribute.Key("r1").String("v1"),
					),
					StartTime: now,
					EndTime:   now,
				},
				{
					Name: "s3",
					Resource: resource.NewWithAttributes(
						semconv.ServiceNameKey.String("name"),
						attribute.Key("r2").String("v2"),
					),
					StartTime: now,
					EndTime:   now,
				},
			},
			expectedBatchList: []*gen.Batch{
				{
					Process: &gen.Process{
						ServiceName: "name",
						Tags: []*gen.Tag{
							{Key: "r1", VType: gen.TagType_STRING, VStr: newString("v1")},
						},
					},
					Spans: []*gen.Span{
						{
							OperationName: "s1",
							Tags: []*gen.Tag{
								{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
						{
							OperationName: "s2",
							Tags: []*gen.Tag{
								{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
					},
				},
				{
					Process: &gen.Process{
						ServiceName: "name",
						Tags: []*gen.Tag{
							{Key: "r2", VType: gen.TagType_STRING, VStr: newString("v2")},
						},
					},
					Spans: []*gen.Span{
						{
							OperationName: "s3",
							Tags: []*gen.Tag{
								{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
					},
				},
			},
		},
		{
			name: "no service name in spans",
			spanSnapshotList: []*sdktrace.SpanSnapshot{
				{
					Name: "s1",
					Resource: resource.NewWithAttributes(
						attribute.Key("r1").String("v1"),
					),
					StartTime: now,
					EndTime:   now,
				},
				nil,
			},
			defaultServiceName: "default service name",
			expectedBatchList: []*gen.Batch{
				{
					Process: &gen.Process{
						ServiceName: "default service name",
						Tags: []*gen.Tag{
							{Key: "r1", VType: gen.TagType_STRING, VStr: newString("v1")},
						},
					},
					Spans: []*gen.Span{
						{
							OperationName: "s1",
							Tags: []*gen.Tag{
								{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchList := jaegerBatchList(tc.spanSnapshotList, tc.defaultServiceName)

			assert.ElementsMatch(t, tc.expectedBatchList, batchList)
		})
	}
}

func TestProcess(t *testing.T) {
	v1 := "v1"

	testCases := []struct {
		name               string
		res                *resource.Resource
		defaultServiceName string
		expectedProcess    *gen.Process
	}{
		{
			name: "resources contain service name",
			res: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("service name"),
				attribute.Key("r1").String("v1"),
			),
			defaultServiceName: "default service name",
			expectedProcess: &gen.Process{
				ServiceName: "service name",
				Tags: []*gen.Tag{
					{Key: "r1", VType: gen.TagType_STRING, VStr: &v1},
				},
			},
		},
		{
			name:               "resources don't have service name",
			res:                resource.NewWithAttributes(attribute.Key("r1").String("v1")),
			defaultServiceName: "default service name",
			expectedProcess: &gen.Process{
				ServiceName: "default service name",
				Tags: []*gen.Tag{
					{Key: "r1", VType: gen.TagType_STRING, VStr: &v1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pro := process(tc.res, tc.defaultServiceName)

			assert.Equal(t, tc.expectedProcess, pro)
		})
	}
}
