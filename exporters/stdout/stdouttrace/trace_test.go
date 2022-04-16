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

package stdouttrace_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestExporterExportSpan(t *testing.T) {
	// setup test span
	now := time.Now()
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")
	traceState, _ := trace.ParseTraceState("key=val")
	keyValue := "value"
	doubleValue := 123.456
	resource := resource.NewSchemaless(attribute.String("rk1", "rv11"))

	ss := tracetest.SpanStub{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceState: traceState,
		}),
		Name:      "/foo",
		StartTime: now,
		EndTime:   now,
		Attributes: []attribute.KeyValue{
			attribute.String("key", keyValue),
			attribute.Float64("double", doubleValue),
		},
		Events: []tracesdk.Event{
			{Name: "foo", Attributes: []attribute.KeyValue{attribute.String("key", keyValue)}, Time: now},
			{Name: "bar", Attributes: []attribute.KeyValue{attribute.Float64("double", doubleValue)}, Time: now},
		},
		SpanKind: trace.SpanKindInternal,
		Status: tracesdk.Status{
			Code:        codes.Error,
			Description: "interesting",
		},
		Resource: resource,
	}

	tests := []struct {
		opts      []stdouttrace.Option
		expectNow time.Time
	}{
		{
			opts:      []stdouttrace.Option{stdouttrace.WithPrettyPrint()},
			expectNow: now,
		},
		{
			opts: []stdouttrace.Option{stdouttrace.WithPrettyPrint(), stdouttrace.WithoutTimestamps()},
			// expectNow is an empty time.Time
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		// write to buffer for testing
		var b bytes.Buffer
		ex, err := stdouttrace.New(append(tt.opts, stdouttrace.WithWriter(&b))...)
		require.Nil(t, err)

		err = ex.ExportSpans(ctx, tracetest.SpanStubs{ss, ss}.Snapshots())
		require.Nil(t, err)

		got := b.String()
		wantone := expectedJSON(tt.expectNow)
		assert.Equal(t, wantone+wantone, got)
	}
}

func expectedJSON(now time.Time) string {
	serializedNow, _ := json.Marshal(now)
	return `{
	"Name": "/foo",
	"SpanContext": {
		"TraceID": "0102030405060708090a0b0c0d0e0f10",
		"SpanID": "0102030405060708",
		"TraceFlags": "00",
		"TraceState": "key=val",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": ` + string(serializedNow) + `,
	"EndTime": ` + string(serializedNow) + `,
	"Attributes": [
		{
			"Key": "key",
			"Value": {
				"Type": "STRING",
				"Value": "value"
			}
		},
		{
			"Key": "double",
			"Value": {
				"Type": "FLOAT64",
				"Value": 123.456
			}
		}
	],
	"Events": [
		{
			"Name": "foo",
			"Attributes": [
				{
					"Key": "key",
					"Value": {
						"Type": "STRING",
						"Value": "value"
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": ` + string(serializedNow) + `
		},
		{
			"Name": "bar",
			"Attributes": [
				{
					"Key": "double",
					"Value": {
						"Type": "FLOAT64",
						"Value": 123.456
					}
				}
			],
			"DroppedAttributeCount": 0,
			"Time": ` + string(serializedNow) + `
		}
	],
	"Links": null,
	"Status": {
		"Code": "Error",
		"Description": "interesting"
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "rk1",
			"Value": {
				"Type": "STRING",
				"Value": "rv11"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "",
		"Version": "",
		"SchemaURL": ""
	}
}
`
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e, err := stdouttrace.New()
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
	defer innerCancel()
	<-innerCtx.Done()
	err = e.Shutdown(innerCtx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestExporterShutdownHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e, err := stdouttrace.New()
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	innerCtx, innerCancel := context.WithCancel(ctx)
	innerCancel()
	err = e.Shutdown(innerCtx)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestExporterShutdownNoError(t *testing.T) {
	e, err := stdouttrace.New()
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	if err := e.Shutdown(context.Background()); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}
