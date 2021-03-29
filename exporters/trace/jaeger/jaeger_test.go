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
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	collectorEndpoint = "http://localhost:14268/api/traces"
	agentEndpoint     = "localhost:6831"
)

func TestInstallNewPipeline(t *testing.T) {
	testCases := []struct {
		name             string
		endpoint         EndpointOption
		options          []Option
		expectedProvider trace.TracerProvider
	}{
		{
			name:             "simple pipeline",
			endpoint:         WithCollectorEndpoint(collectorEndpoint),
			expectedProvider: &sdktrace.TracerProvider{},
		},
		{
			name:             "with agent endpoint",
			endpoint:         WithAgentEndpoint(agentEndpoint),
			expectedProvider: &sdktrace.TracerProvider{},
		},
		{
			name:     "with disabled",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithDisabled(true),
			},
			expectedProvider: trace.NewNoopTracerProvider(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn, err := InstallNewPipeline(
				tc.endpoint,
				tc.options...,
			)
			defer fn()

			assert.NoError(t, err)
			assert.IsType(t, tc.expectedProvider, otel.GetTracerProvider())

			otel.SetTracerProvider(nil)
		})
	}
}

func TestNewExportPipeline(t *testing.T) {
	testCases := []struct {
		name                                  string
		endpoint                              EndpointOption
		options                               []Option
		expectedProviderType                  trace.TracerProvider
		testSpanSampling, spanShouldBeSampled bool
	}{
		{
			name:                 "simple pipeline",
			endpoint:             WithCollectorEndpoint(collectorEndpoint),
			expectedProviderType: &sdktrace.TracerProvider{},
		},
		{
			name:     "with disabled",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithDisabled(true),
			},
			expectedProviderType: trace.NewNoopTracerProvider(),
		},
		{
			name:     "always on",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithSDKOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
			},
			expectedProviderType: &sdktrace.TracerProvider{},
			testSpanSampling:     true,
			spanShouldBeSampled:  true,
		},
		{
			name:     "never",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithSDKOptions(sdktrace.WithSampler(sdktrace.NeverSample())),
			},
			expectedProviderType: &sdktrace.TracerProvider{},
			testSpanSampling:     true,
			spanShouldBeSampled:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp, fn, err := NewExportPipeline(
				tc.endpoint,
				tc.options...,
			)
			defer fn()

			assert.NoError(t, err)
			assert.NotEqual(t, tp, otel.GetTracerProvider())
			assert.IsType(t, tc.expectedProviderType, tp)

			if tc.testSpanSampling {
				_, span := tp.Tracer("jaeger test").Start(context.Background(), tc.name)
				spanCtx := span.SpanContext()
				assert.Equal(t, tc.spanShouldBeSampled, spanCtx.IsSampled())
				span.End()
			}
		})
	}
}

func TestNewExportPipelineWithDisabledFromEnv(t *testing.T) {
	envStore, err := ottest.SetEnvVariables(map[string]string{
		envDisabled: "true",
	})
	require.NoError(t, err)
	envStore.Record(envDisabled)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	tp, fn, err := NewExportPipeline(
		WithCollectorEndpoint(collectorEndpoint),
	)
	defer fn()
	assert.NoError(t, err)
	assert.IsType(t, trace.NewNoopTracerProvider(), tp)
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
			endpoint:               WithCollectorEndpoint(collectorEndpoint),
			expectedServiceName:    "unknown_service",
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:                   "default exporter with agent endpoint",
			endpoint:               WithAgentEndpoint(agentEndpoint),
			expectedServiceName:    "unknown_service",
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:     "with buffer and batch max count",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
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

func TestNewRawExporterShouldFail(t *testing.T) {
	testCases := []struct {
		name           string
		endpoint       EndpointOption
		expectedErrMsg string
	}{
		{
			name:           "with empty collector endpoint",
			endpoint:       WithCollectorEndpoint(""),
			expectedErrMsg: "collectorEndpoint must not be empty",
		},
		{
			name:           "with empty agent endpoint",
			endpoint:       WithAgentEndpoint(""),
			expectedErrMsg: "agentEndpoint must not be empty",
		},
		{
			name:           "with invalid agent endpoint",
			endpoint:       WithAgentEndpoint("localhost"),
			expectedErrMsg: "address localhost: missing port in address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRawExporter(
				tc.endpoint,
			)

			assert.Error(t, err)
			assert.EqualError(t, err, tc.expectedErrMsg)
		})
	}
}

func TestNewRawExporterShouldFailIfCollectorUnset(t *testing.T) {
	// Record and restore env
	envStore := ottest.NewEnvStore()
	envStore.Record(envEndpoint)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// If the user sets the environment variable JAEGER_ENDPOINT, endpoint will always get a value.
	require.NoError(t, os.Unsetenv(envEndpoint))

	_, err := NewRawExporter(
		WithCollectorEndpoint(""),
	)

	assert.Error(t, err)
}

type testCollectorEndpoint struct {
	batchesUploaded []*gen.Batch
}

func (c *testCollectorEndpoint) upload(batch *gen.Batch) error {
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

func TestExporter_ExportSpan(t *testing.T) {
	const (
		serviceName = "test-service"
		tagKey      = "key"
		tagVal      = "val"
	)
	// Create Jaeger Exporter
	exp, err := NewRawExporter(
		withTestCollectorEndpoint(),
		WithSDKOptions(
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.ServiceNameKey.String(serviceName),
				attribute.String(tagKey, tagVal),
			)),
		),
	)

	assert.NoError(t, err)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exp),
	)
	otel.SetTracerProvider(tp)
	_, span := otel.Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()

	assert.True(t, span.SpanContext().IsValid())

	exp.Flush()
	tc := exp.uploader.(*testCollectorEndpoint)
	assert.True(t, len(tc.batchesUploaded) == 1)
	assert.True(t, len(tc.batchesUploaded[0].GetSpans()) == 1)
}

func Test_spanSnapshotToThrift(t *testing.T) {
	now := time.Now()
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")
	parentSpanID, _ := trace.SpanIDFromHex("0807060504030201")

	linkTraceID, _ := trace.TraceIDFromHex("0102030405060709090a0b0c0d0e0f11")
	linkSpanID, _ := trace.SpanIDFromHex("0102030405060709")

	eventNameValue := "event-test"
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
		data *export.SpanSnapshot
		want *gen.Span
	}{
		{
			name: "no parent",
			data: &export.SpanSnapshot{
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
					{Name: eventNameValue, Attributes: []attribute.KeyValue{attribute.String("k1", keyValue)}, Time: now},
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
					{Key: "error", VType: gen.TagType_BOOL, VBool: &boolTrue},
					{Key: "otel.library.name", VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: "otel.library.version", VType: gen.TagType_STRING, VStr: &instrLibVersion},
					{Key: "status.code", VType: gen.TagType_LONG, VLong: &statusCodeValue},
					{Key: "status.message", VType: gen.TagType_STRING, VStr: &statusMessage},
					{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
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
								Key:   "k1",
								VStr:  &keyValue,
								VType: gen.TagType_STRING,
							},
							{
								Key:   "name",
								VStr:  &eventNameValue,
								VType: gen.TagType_STRING,
							},
						},
					},
				},
			},
		},
		{
			name: "with parent",
			data: &export.SpanSnapshot{
				ParentSpanID: parentSpanID,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
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
					{Key: "otel.library.name", VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: "otel.library.version", VType: gen.TagType_STRING, VStr: &instrLibVersion},
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
			data: &export.SpanSnapshot{
				ParentSpanID: parentSpanID,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					SpanID:  spanID,
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
					{Key: "otel.library.name", VType: gen.TagType_STRING, VStr: &instrLibName},
					{Key: "otel.library.version", VType: gen.TagType_STRING, VStr: &instrLibVersion},
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

func TestJaegerBatchList(t *testing.T) {
	newString := func(value string) *string {
		return &value
	}
	spanKind := "unspecified"
	now := time.Now()

	testCases := []struct {
		name                string
		spanSnapshotList    []*export.SpanSnapshot
		defaultServiceName  string
		resourceFromProcess *resource.Resource
		expectedBatchList   []*gen.Batch
	}{
		{
			name:              "no span shots",
			spanSnapshotList:  nil,
			expectedBatchList: nil,
		},
		{
			name: "span's snapshot contains nil span",
			spanSnapshotList: []*export.SpanSnapshot{
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
								{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
							},
							StartTime: now.UnixNano() / 1000,
						},
					},
				},
			},
		},
		{
			name: "merge spans that have the same resources",
			spanSnapshotList: []*export.SpanSnapshot{
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
			name: "merge resources that come from process",
			spanSnapshotList: []*export.SpanSnapshot{
				{
					Name: "s1",
					Resource: resource.NewWithAttributes(
						semconv.ServiceNameKey.String("name"),
						attribute.Key("r1").String("v1"),
						attribute.Key("r2").String("v2"),
					),
					StartTime: now,
					EndTime:   now,
				},
			},
			resourceFromProcess: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("new-name"),
				attribute.Key("r1").String("v2"),
			),
			expectedBatchList: []*gen.Batch{
				{
					Process: &gen.Process{
						ServiceName: "new-name",
						Tags: []*gen.Tag{
							{Key: "r1", VType: gen.TagType_STRING, VStr: newString("v2")},
							{Key: "r2", VType: gen.TagType_STRING, VStr: newString("v2")},
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
		{
			name: "span's snapshot contains no service name but resourceFromProcess does",
			spanSnapshotList: []*export.SpanSnapshot{
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
			resourceFromProcess: resource.NewWithAttributes(
				semconv.ServiceNameKey.String("new-name"),
			),
			expectedBatchList: []*gen.Batch{
				{
					Process: &gen.Process{
						ServiceName: "new-name",
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
		{
			name: "no service name in spans and resourceFromProcess",
			spanSnapshotList: []*export.SpanSnapshot{
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
			batchList := jaegerBatchList(tc.spanSnapshotList, tc.defaultServiceName, tc.resourceFromProcess)

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

func TestNewExporterPipelineWithOptions(t *testing.T) {
	envStore, err := ottest.SetEnvVariables(map[string]string{
		envDisabled: "false",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	const (
		serviceName     = "test-service"
		eventCountLimit = 10
		tagKey          = "key"
		tagVal          = "val"
	)

	testCollector := &testCollectorEndpoint{}
	tp, spanFlush, err := NewExportPipeline(
		withTestCollectorEndpointInjected(testCollector),
		WithSDKOptions(
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.ServiceNameKey.String(serviceName),
				attribute.String(tagKey, tagVal),
			)),
			sdktrace.WithSpanLimits(sdktrace.SpanLimits{
				EventCountLimit: eventCountLimit,
			}),
		),
	)
	assert.NoError(t, err)

	otel.SetTracerProvider(tp)
	_, span := otel.Tracer("test-tracer").Start(context.Background(), "test-span")
	for i := 0; i < eventCountLimit*2; i++ {
		span.AddEvent(fmt.Sprintf("event-%d", i))
	}
	span.End()
	spanFlush()

	assert.True(t, span.SpanContext().IsValid())

	batchesUploaded := testCollector.batchesUploaded
	assert.True(t, len(batchesUploaded) == 1)
	uploadedBatch := batchesUploaded[0]
	assert.Equal(t, serviceName, uploadedBatch.GetProcess().GetServiceName())
	assert.True(t, len(uploadedBatch.GetSpans()) == 1)
	uploadedSpan := uploadedBatch.GetSpans()[0]
	assert.Equal(t, eventCountLimit, len(uploadedSpan.GetLogs()))

	assert.Equal(t, 1, len(uploadedBatch.GetProcess().GetTags()))
	assert.Equal(t, tagKey, uploadedBatch.GetProcess().GetTags()[0].GetKey())
	assert.Equal(t, tagVal, uploadedBatch.GetProcess().GetTags()[0].GetVStr())
}
