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
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	apitrace "go.opentelemetry.io/otel/api/trace"
	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewExporterPipelineWithRegistration(t *testing.T) {
	tp, fn, err := NewExportPipeline(
		WithCollectorEndpoint("http://localhost:14268/api/traces"),
		RegisterAsGlobal(),
	)
	defer fn()
	assert.NoError(t, err)
	assert.Same(t, tp, global.TraceProvider())
}

func TestNewExporterPipelineWithoutRegistration(t *testing.T) {
	tp, fn, err := NewExportPipeline(
		WithCollectorEndpoint("http://localhost:14268/api/traces"),
	)
	defer fn()
	assert.NoError(t, err)
	assert.NotEqual(t, tp, global.TraceProvider())
}

func TestNewExporterPipelineWithSDK(t *testing.T) {
	tp, fn, err := NewExportPipeline(
		WithCollectorEndpoint("http://localhost:14268/api/traces"),
		WithSDK(&sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
		}),
	)
	defer fn()
	assert.NoError(t, err)
	_, span := tp.Tracer("jaeger test").Start(context.Background(), "always-on")
	spanCtx := span.SpanContext()
	assert.True(t, spanCtx.IsSampled())
	span.End()

	tp2, fn, err := NewExportPipeline(
		WithCollectorEndpoint("http://localhost:14268/api/traces"),
		WithSDK(&sdktrace.Config{
			DefaultSampler: sdktrace.NeverSample(),
		}),
	)
	defer fn()
	assert.NoError(t, err)
	_, span2 := tp2.Tracer("jaeger test").Start(context.Background(), "never")
	span2Ctx := span2.SpanContext()
	assert.False(t, span2Ctx.IsSampled())
	span2.End()
}

func TestNewRawExporter(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
		serviceName       = "test-service"
		tagKey            = "key"
		tagVal            = "val"
	)
	// Create Jaeger Exporter
	exp, err := NewRawExporter(
		WithCollectorEndpoint(collectorEndpoint),
		WithProcess(Process{
			ServiceName: serviceName,
			Tags: []kv.KeyValue{
				kv.String(tagKey, tagVal),
			},
		}),
	)

	assert.NoError(t, err)
	assert.EqualValues(t, serviceName, exp.process.ServiceName)
	assert.Len(t, exp.process.Tags, 1)
}

func TestNewRawExporterShouldFailIfCollectorEndpointEmpty(t *testing.T) {
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

func TestNewRawExporterWithAgentEndpoint(t *testing.T) {
	const agentEndpoint = "localhost:6831"
	// Create Jaeger Exporter
	_, err := NewRawExporter(
		WithAgentEndpoint(agentEndpoint),
	)
	assert.NoError(t, err)
}

func TestNewRawExporterWithAgentShouldFailIfEndpointInvalid(t *testing.T) {
	//empty
	_, err := NewRawExporter(
		WithAgentEndpoint(""),
	)
	assert.Error(t, err)

	//invalid endpoint addr
	_, err = NewRawExporter(
		WithAgentEndpoint("http://localhost"),
	)
	assert.Error(t, err)
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
					// Jaeger doesn't handle Uint tags, this should be ignored.
					kv.Uint64("ignored", 123),
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
