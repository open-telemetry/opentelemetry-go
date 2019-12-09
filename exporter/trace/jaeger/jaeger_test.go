// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"

	apitrace "go.opentelemetry.io/otel/api/trace"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/api/core"
	gen "go.opentelemetry.io/otel/exporter/trace/jaeger/internal/gen-go/jaeger"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

func TestNewExporter(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
		serviceName       = "test-service"
		tagKey            = "key"
		tagVal            = "val"
	)
	// Create Jaeger Exporter
	exp, err := NewExporter(
		WithCollectorEndpoint(collectorEndpoint),
		WithProcess(Process{
			ServiceName: serviceName,
			Tags: []core.KeyValue{
				key.String(tagKey, tagVal),
			},
		}),
	)

	assert.NoError(t, err)
	assert.EqualValues(t, serviceName, exp.process.ServiceName)
	assert.Len(t, exp.process.Tags, 1)
}

func TestNewExporterShouldFailIfCollectorEndpointEmpty(t *testing.T) {
	_, err := NewExporter(
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
	exp, err := NewExporter(
		withTestCollectorEndpoint(),
		WithProcess(Process{
			ServiceName: serviceName,
			Tags: []core.KeyValue{
				key.String(tagKey, tagVal),
			},
		}),
	)

	assert.NoError(t, err)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exp))

	assert.NoError(t, err)

	global.SetTraceProvider(tp)
	_, span := global.TraceProvider().Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()

	assert.True(t, span.SpanContext().IsValid())

	exp.Flush()
	tc := exp.uploader.(*testCollectorEnpoint)
	assert.True(t, len(tc.spansUploaded) == 1)
}

func TestNewExporterWithAgentEndpoint(t *testing.T) {
	const agentEndpoint = "localhost:6831"
	// Create Jaeger Exporter
	_, err := NewExporter(
		WithAgentEndpoint(agentEndpoint),
	)
	assert.NoError(t, err)
}

func TestNewExporterWithAgentShouldFailIfEndpointInvalid(t *testing.T) {
	//empty
	_, err := NewExporter(
		WithAgentEndpoint(""),
	)
	assert.Error(t, err)

	//invalid endpoint addr
	_, err = NewExporter(
		WithAgentEndpoint("http://localhost"),
	)
	assert.Error(t, err)
}

func Test_spanDataToThrift(t *testing.T) {
	now := time.Now()
	traceID, _ := core.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := core.SpanIDFromHex("0102030405060708")

	linkTraceID, _ := core.TraceIDFromHex("0102030405060709090a0b0c0d0e0f11")
	linkSpanID, _ := core.SpanIDFromHex("0102030405060709")

	messageEventValue := "event-test"
	keyValue := "value"
	statusCodeValue := int64(2)
	doubleValue := 123.456
	boolTrue := true
	statusMessage := "Unknown"

	tests := []struct {
		name string
		data *export.SpanData
		want *gen.Span
	}{
		{
			name: "no parent",
			data: &export.SpanData{
				SpanContext: core.SpanContext{
					TraceID: traceID,
					SpanID:  spanID,
				},
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Links: []apitrace.Link{
					{
						SpanContext: core.SpanContext{
							TraceID: linkTraceID,
							SpanID:  linkSpanID,
						},
					},
				},
				Attributes: []core.KeyValue{
					key.String("key", keyValue),
					key.Float64("double", doubleValue),
					// Jaeger doesn't handle Uint tags, this should be ignored.
					key.Uint64("ignored", 123),
				},
				MessageEvents: []export.Event{
					{Message: messageEventValue, Attributes: []core.KeyValue{key.String("k1", keyValue)}, Time: now},
				},
				Status: codes.Unknown,
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
								Key:   "message",
								VStr:  &messageEventValue,
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
