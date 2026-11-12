// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlphttp_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	collectorlog "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	collectormetric "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"
)

func TestMetricWireCompatibility(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL", "http/protobuf")
	server, requests := collector(t, "/v1/metrics", &collectormetric.ExportMetricsServiceResponse{
		PartialSuccess: &collectormetric.ExportMetricsPartialSuccess{},
	})
	exporter, err := otlpmetrichttp.New(t.Context(),
		otlpmetrichttp.WithEndpointURL(server.URL+"/v1/metrics"),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := exporter.Shutdown(ctx); err != nil {
			t.Error(err)
		}
	})

	now := time.Unix(1_700_000_000, 123)
	data := metricdata.ResourceMetrics{
		Resource: resource.NewWithAttributes("resource-schema", attribute.String("service.name", "wire-test")),
		ScopeMetrics: []metricdata.ScopeMetrics{{
			Scope: instrumentation.Scope{Name: "wire-meter", Version: "1.2.3", SchemaURL: "scope-schema"},
			Metrics: []metricdata.Metrics{{
				Name: "requests", Description: "request count", Unit: "{request}",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{{
						Attributes: attribute.NewSet(attribute.String("method", "GET")),
						StartTime:  now.Add(-time.Second), Time: now, Value: 42,
					}},
				},
			}},
		}},
	}
	if err := exporter.Export(t.Context(), &data); err != nil {
		t.Fatal(err)
	}
	var request collectormetric.ExportMetricsServiceRequest
	if err := proto.Unmarshal(receiveRequest(t, requests), &request); err != nil {
		t.Fatal(err)
	}
	if len(request.ResourceMetrics) != 1 || len(request.ResourceMetrics[0].ScopeMetrics) != 1 {
		t.Fatalf("unexpected resource/scope metrics: %v", &request)
	}
	rm := request.ResourceMetrics[0]
	sm := rm.ScopeMetrics[0]
	if rm.SchemaUrl != "resource-schema" || len(rm.GetResource().GetAttributes()) != 1 ||
		rm.Resource.Attributes[0].Key != "service.name" || rm.Resource.Attributes[0].Value.GetStringValue() != "wire-test" ||
		sm.GetScope().GetName() != "wire-meter" || sm.GetScope().GetVersion() != "1.2.3" || sm.SchemaUrl != "scope-schema" || len(sm.Metrics) != 1 {
		t.Fatalf("resource/scope metadata changed: %v", rm)
	}
	m := sm.Metrics[0]
	sum := m.GetSum()
	if m.Name != "requests" || m.Description != "request count" || m.Unit != "{request}" ||
		sum.GetAggregationTemporality() != metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE ||
		!sum.GetIsMonotonic() || len(sum.GetDataPoints()) != 1 {
		t.Fatalf("metric changed: %v", m)
	}
	dp := sum.DataPoints[0]
	if dp.GetAsInt() != 42 || dp.TimeUnixNano != uint64(now.UnixNano()) ||
		dp.StartTimeUnixNano != uint64(now.Add(-time.Second).UnixNano()) ||
		len(dp.Attributes) != 1 || dp.Attributes[0].Key != "method" || dp.Attributes[0].Value.GetStringValue() != "GET" {
		t.Fatalf("metric data point changed: %v", dp)
	}
}

func TestLogWireCompatibility(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_LOGS_PROTOCOL", "http/protobuf")
	server, requests := collector(t, "/v1/logs", &collectorlog.ExportLogsServiceResponse{
		PartialSuccess: &collectorlog.ExportLogsPartialSuccess{},
	})
	exporter, err := otlploghttp.New(t.Context(),
		otlploghttp.WithEndpointURL(server.URL+"/v1/logs"),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := exporter.Shutdown(ctx); err != nil {
			t.Error(err)
		}
	})

	now := time.Unix(1_700_000_000, 123)
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	record := logtest.RecordFactory{
		Timestamp: now, ObservedTimestamp: now.Add(time.Second),
		Severity: log.SeverityWarn, SeverityText: "WARN", Body: attribute.StringValue("request failed"),
		Attributes: []attribute.KeyValue{attribute.Int("status", 503)},
		TraceID:    traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled,
		Resource:             resource.NewWithAttributes("resource-schema", attribute.String("service.name", "wire-test")),
		InstrumentationScope: &instrumentation.Scope{Name: "wire-logger", Version: "1.2.3", SchemaURL: "scope-schema"},
	}.NewRecord()
	if err := exporter.Export(t.Context(), []sdklog.Record{record}); err != nil {
		t.Fatal(err)
	}
	var request collectorlog.ExportLogsServiceRequest
	if err := proto.Unmarshal(receiveRequest(t, requests), &request); err != nil {
		t.Fatal(err)
	}
	if len(request.ResourceLogs) != 1 || len(request.ResourceLogs[0].ScopeLogs) != 1 {
		t.Fatalf("unexpected resource/scope logs: %v", &request)
	}
	rl := request.ResourceLogs[0]
	sl := rl.ScopeLogs[0]
	if rl.SchemaUrl != "resource-schema" || len(rl.GetResource().GetAttributes()) != 1 ||
		rl.Resource.Attributes[0].Key != "service.name" || rl.Resource.Attributes[0].Value.GetStringValue() != "wire-test" ||
		sl.GetScope().GetName() != "wire-logger" || sl.GetScope().GetVersion() != "1.2.3" || sl.SchemaUrl != "scope-schema" || len(sl.LogRecords) != 1 {
		t.Fatalf("resource/scope metadata changed: %v", rl)
	}
	lr := sl.LogRecords[0]
	if lr.Body.GetStringValue() != "request failed" || lr.SeverityNumber != logpb.SeverityNumber_SEVERITY_NUMBER_WARN || lr.SeverityText != "WARN" ||
		lr.TimeUnixNano != uint64(now.UnixNano()) || lr.ObservedTimeUnixNano != uint64(now.Add(time.Second).UnixNano()) ||
		!bytes.Equal(lr.TraceId, traceID[:]) || !bytes.Equal(lr.SpanId, spanID[:]) || lr.Flags != uint32(trace.FlagsSampled) ||
		len(lr.Attributes) != 1 || lr.Attributes[0].Key != "status" || lr.Attributes[0].Value.GetIntValue() != 503 {
		t.Fatalf("log record changed: %v", lr)
	}
}

func receiveRequest(t *testing.T, requests <-chan []byte) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	select {
	case body := <-requests:
		return body
	case <-ctx.Done():
		t.Fatal("no export request received: ", ctx.Err())
		return nil
	}
}

func collector(t *testing.T, path string, response proto.Message) (*httptest.Server, <-chan []byte) {
	t.Helper()
	responseBody, err := proto.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	requests := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != path || r.Header.Get("Content-Type") != "application/x-protobuf" ||
			r.Header.Get("Content-Encoding") != "gzip" {
			t.Errorf("unexpected request: %s %s %v", r.Method, r.URL.Path, r.Header)
		}
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer reader.Close()
		body, err := io.ReadAll(reader)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		select {
		case requests <- body:
		default:
			t.Error("unexpected duplicate export request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/x-protobuf")
		if _, err := w.Write(responseBody); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)
	return server, requests
}
