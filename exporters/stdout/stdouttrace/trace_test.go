// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
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
	res := resource.NewSchemaless(attribute.String("rk1", "rv11"))

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
		Resource: res,
	}

	tests := []struct {
		opts      []stdouttrace.Option
		expectNow time.Time
		ctx       context.Context
		wantErr   error
	}{
		{
			opts:      []stdouttrace.Option{stdouttrace.WithPrettyPrint()},
			expectNow: now,
			ctx:       context.Background(),
		},
		{
			opts: []stdouttrace.Option{stdouttrace.WithPrettyPrint(), stdouttrace.WithoutTimestamps()},
			// expectNow is an empty time.Time
			ctx: context.Background(),
		},
		{
			opts: []stdouttrace.Option{},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		// write to buffer for testing
		var b bytes.Buffer
		ex, err := stdouttrace.New(append(tt.opts, stdouttrace.WithWriter(&b))...)
		require.NoError(t, err)

		err = ex.ExportSpans(tt.ctx, tracetest.SpanStubs{ss, ss}.Snapshots())
		assert.Equal(t, tt.wantErr, err)

		if tt.wantErr == nil {
			got := b.String()
			wantone := expectedJSON(tt.expectNow)
			assert.Equal(t, wantone+wantone, got)
		}
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
	"InstrumentationScope": {
		"Name": "",
		"Version": "",
		"SchemaURL": "",
		"Attributes": null
	},
	"InstrumentationLibrary": {
		"Name": "",
		"Version": "",
		"SchemaURL": "",
		"Attributes": null
	}
}
`
}

func TestExporterShutdownIgnoresContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e, err := stdouttrace.New()
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	innerCtx, innerCancel := context.WithCancel(ctx)
	innerCancel()
	err = e.Shutdown(innerCtx)
	assert.NoError(t, err)
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

func TestSelfObservability(t *testing.T) {
	defaultCallExportSpans := func(t *testing.T, exporter *stdouttrace.Exporter) {
		require.NoError(t, exporter.ExportSpans(context.Background(), tracetest.SpanStubs{
			{Name: "/foo"},
			{Name: "/bar"},
		}.Snapshots()))
	}

	tests := []struct {
		name            string
		enabled         bool
		callExportSpans func(t *testing.T, exporter *stdouttrace.Exporter)
		assertMetrics   func(t *testing.T, rm metricdata.ResourceMetrics)
	}{
		{
			name:            "Disabled",
			enabled:         false,
			callExportSpans: defaultCallExportSpans,
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Empty(t, rm.ScopeMetrics)
			},
		},
		{
			name:            "Enabled",
			enabled:         true,
			callExportSpans: defaultCallExportSpans,
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				t.Helper()
				require.Len(t, rm.ScopeMetrics, 1)

				sm := rm.ScopeMetrics[0]
				require.Len(t, sm.Metrics, 3)

				assert.Equal(t, instrumentation.Scope{
					Name:      "go.opentelemetry.io/otel/exporters/stdout/stdouttrace",
					Version:   sdk.Version(),
					SchemaURL: semconv.SchemaURL,
				}, sm.Scope)

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterSpanInflight{}.Name(),
					Description: otelconv.SDKExporterSpanInflight{}.Description(),
					Unit:        otelconv.SDKExporterSpanInflight{}.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/0"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
								),
								Value: 0,
							},
						},
					},
				}, sm.Metrics[0], metricdatatest.IgnoreTimestamp())

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterSpanExported{}.Name(),
					Description: otelconv.SDKExporterSpanExported{}.Description(),
					Unit:        otelconv.SDKExporterSpanExported{}.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/0"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
								),
								Value: 2,
							},
						},
					},
				}, sm.Metrics[1], metricdatatest.IgnoreTimestamp())

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterOperationDuration{}.Name(),
					Description: otelconv.SDKExporterOperationDuration{}.Description(),
					Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
					Data: metricdata.Histogram[float64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.HistogramDataPoint[float64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/0"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
								),
							},
						},
					},
				}, sm.Metrics[2], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
			},
		},
		{
			name:    "Enabled, but ExportSpans returns error",
			enabled: true,
			callExportSpans: func(t *testing.T, exporter *stdouttrace.Exporter) {
				t.Helper()
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				err := exporter.ExportSpans(ctx, tracetest.SpanStubs{
					{Name: "/foo"},
					{Name: "/bar"},
				}.Snapshots())
				require.Error(t, err)
			},
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				t.Helper()
				require.Len(t, rm.ScopeMetrics, 1)

				sm := rm.ScopeMetrics[0]
				require.Len(t, sm.Metrics, 3)

				assert.Equal(t, instrumentation.Scope{
					Name:      "go.opentelemetry.io/otel/exporters/stdout/stdouttrace",
					Version:   sdk.Version(),
					SchemaURL: semconv.SchemaURL,
				}, sm.Scope)

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterSpanInflight{}.Name(),
					Description: otelconv.SDKExporterSpanInflight{}.Description(),
					Unit:        otelconv.SDKExporterSpanInflight{}.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/1"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
								),
								Value: 0,
							},
						},
					},
				}, sm.Metrics[0], metricdatatest.IgnoreTimestamp())

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterSpanExported{}.Name(),
					Description: otelconv.SDKExporterSpanExported{}.Description(),
					Unit:        otelconv.SDKExporterSpanExported{}.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/1"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
									semconv.ErrorType(context.Canceled),
								),
								Value: 2,
							},
						},
					},
				}, sm.Metrics[1], metricdatatest.IgnoreTimestamp())

				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        otelconv.SDKExporterOperationDuration{}.Name(),
					Description: otelconv.SDKExporterOperationDuration{}.Description(),
					Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
					Data: metricdata.Histogram[float64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.HistogramDataPoint[float64]{
							{
								Attributes: attribute.NewSet(
									semconv.OTelComponentName("stdout_trace_exporter/1"),
									semconv.OTelComponentTypeKey.String("stdout_trace_exporter"),
									semconv.ErrorType(context.Canceled),
								),
							},
						},
					},
				}, sm.Metrics[2], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.enabled {
				t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
			}

			original := otel.GetMeterProvider()
			defer otel.SetMeterProvider(original)

			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(r))
			otel.SetMeterProvider(mp)

			exporter, err := stdouttrace.New(
				stdouttrace.WithWriter(io.Discard))
			require.NoError(t, err)

			tt.callExportSpans(t, exporter)

			var rm metricdata.ResourceMetrics
			require.NoError(t, r.Collect(context.Background(), &rm))

			tt.assertMetrics(t, rm)
		})
	}
}
