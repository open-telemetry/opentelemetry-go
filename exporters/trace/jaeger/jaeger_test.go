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
	"math"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/support/bundler"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"
	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
	ottest "go.opentelemetry.io/otel/internal/testing"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
		expectedProvider trace.Provider
	}{
		{
			name:             "simple pipeline",
			endpoint:         WithCollectorEndpoint(collectorEndpoint),
			expectedProvider: &sdktrace.Provider{},
		},
		{
			name:             "with agent endpoint",
			endpoint:         WithAgentEndpoint(agentEndpoint),
			expectedProvider: &sdktrace.Provider{},
		},
		{
			name:     "with disabled",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithDisabled(true),
			},
			expectedProvider: &apitrace.NoopProvider{},
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
			assert.IsType(t, tc.expectedProvider, global.TraceProvider())

			global.SetTraceProvider(nil)
		})
	}
}

func TestNewExportPipeline(t *testing.T) {
	testCases := []struct {
		name                                  string
		endpoint                              EndpointOption
		options                               []Option
		expectedProviderType                  trace.Provider
		testSpanSampling, spanShouldBeSampled bool
	}{
		{
			name:                 "simple pipeline",
			endpoint:             WithCollectorEndpoint(collectorEndpoint),
			expectedProviderType: &sdktrace.Provider{},
		},
		{
			name:     "with disabled",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithDisabled(true),
			},
			expectedProviderType: &apitrace.NoopProvider{},
		},
		{
			name:     "always on",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithSDK(&sdktrace.Config{
					DefaultSampler: sdktrace.AlwaysSample(),
				}),
			},
			expectedProviderType: &sdktrace.Provider{},
			testSpanSampling:     true,
			spanShouldBeSampled:  true,
		},
		{
			name:     "never",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithSDK(&sdktrace.Config{
					DefaultSampler: sdktrace.NeverSample(),
				}),
			},
			expectedProviderType: &sdktrace.Provider{},
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
			assert.NotEqual(t, tp, global.TraceProvider())
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
	assert.IsType(t, &apitrace.NoopProvider{}, tp)
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
			expectedServiceName:    defaultServiceName,
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:                   "default exporter with agent endpoint",
			endpoint:               WithAgentEndpoint(agentEndpoint),
			expectedServiceName:    defaultServiceName,
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:     "with process",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithProcess(
					Process{
						ServiceName: "jaeger-test",
						Tags: []kv.KeyValue{
							kv.String("key", "val"),
						},
					},
				),
			},
			expectedServiceName:    "jaeger-test",
			expectedTagsLen:        1,
			expectedBufferMaxCount: bundler.DefaultBufferedByteLimit,
			expectedBatchMaxCount:  bundler.DefaultBundleCountThreshold,
		},
		{
			name:     "with buffer and batch max count",
			endpoint: WithCollectorEndpoint(collectorEndpoint),
			options: []Option{
				WithProcess(
					Process{
						ServiceName: "jaeger-test",
					},
				),
				WithBufferMaxCount(99),
				WithBatchMaxCount(99),
			},
			expectedServiceName:    "jaeger-test",
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
			assert.Equal(t, tc.expectedServiceName, exp.process.ServiceName)
			assert.Len(t, exp.process.Tags, tc.expectedTagsLen)
			assert.Equal(t, tc.expectedBufferMaxCount, exp.bundler.BufferedByteLimit)
			assert.Equal(t, tc.expectedBatchMaxCount, exp.bundler.BundleCountThreshold)
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

type testCollectorEnpoint struct {
	spansUploaded []*gen.Span
}

func (c *testCollectorEnpoint) upload(batch *gen.Batch) error {
	c.spansUploaded = append(c.spansUploaded, batch.Spans...)
	return nil
}

var _ batchUploader = (*testCollectorEnpoint)(nil)

func withTestCollectorEndpoint() func() (batchUploader, error) {
	return func() (batchUploader, error) {
		return &testCollectorEnpoint{}, nil
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
		WithProcess(Process{
			ServiceName: serviceName,
			Tags: []kv.KeyValue{
				kv.String(tagKey, tagVal),
			},
		}),
	)

	assert.NoError(t, err)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exp))

	assert.NoError(t, err)

	global.SetTraceProvider(tp)
	_, span := global.Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()

	assert.True(t, span.SpanContext().IsValid())

	exp.Flush()
	tc := exp.uploader.(*testCollectorEnpoint)
	assert.True(t, len(tc.spansUploaded) == 1)
}

func Test_spanDataToThrift(t *testing.T) {
	now := time.Now()
	traceID, _ := apitrace.IDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := apitrace.SpanIDFromHex("0102030405060708")

	linkTraceID, _ := apitrace.IDFromHex("0102030405060709090a0b0c0d0e0f11")
	linkSpanID, _ := apitrace.SpanIDFromHex("0102030405060709")

	eventNameValue := "event-test"
	keyValue := "value"
	statusCodeValue := int64(2)
	doubleValue := 123.456
	uintValue := int64(123)
	boolTrue := true
	statusMessage := "this is a problem"
	spanKind := "client"
	rv1 := "rv11"
	rv2 := int64(5)

	tests := []struct {
		name string
		data *export.SpanData
		want *gen.Span
	}{
		{
			name: "no parent",
			data: &export.SpanData{
				SpanContext: apitrace.SpanContext{
					TraceID: traceID,
					SpanID:  spanID,
				},
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Links: []apitrace.Link{
					{
						SpanContext: apitrace.SpanContext{
							TraceID: linkTraceID,
							SpanID:  linkSpanID,
						},
					},
				},
				Attributes: []kv.KeyValue{
					kv.String("key", keyValue),
					kv.Float64("double", doubleValue),
					kv.Uint64("uint", uint64(uintValue)),
					kv.Uint64("overflows", math.MaxUint64),
				},
				MessageEvents: []export.Event{
					{Name: eventNameValue, Attributes: []kv.KeyValue{kv.String("k1", keyValue)}, Time: now},
				},
				StatusCode:    codes.Unknown,
				StatusMessage: statusMessage,
				SpanKind:      apitrace.SpanKindClient,
				Resource:      resource.New(kv.String("rk1", rv1), kv.Int64("rk2", rv2)),
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
					{Key: "uint", VType: gen.TagType_LONG, VLong: &uintValue},
					{Key: "error", VType: gen.TagType_BOOL, VBool: &boolTrue},
					{Key: "status.code", VType: gen.TagType_LONG, VLong: &statusCodeValue},
					{Key: "status.message", VType: gen.TagType_STRING, VStr: &statusMessage},
					{Key: "span.kind", VType: gen.TagType_STRING, VStr: &spanKind},
					{Key: "rk1", VType: gen.TagType_STRING, VStr: &rv1},
					{Key: "rk2", VType: gen.TagType_LONG, VLong: &rv2},
				},
				References: []*gen.SpanRef{
					{
						RefType:     gen.SpanRefType_CHILD_OF,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spanDataToThrift(tt.data)
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
