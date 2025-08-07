// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for OTLP metric exporters.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/selfobservability"

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// exporterIDCounter is used to generate unique component names for exporters.
var exporterIDCounter atomic.Uint64

// nextExporterID returns the next unique exporter ID.
func nextExporterID() uint64 {
	return exporterIDCounter.Add(1) - 1
}

// ExporterMetrics holds the self-observability metric instruments for an OTLP metric exporter.
type ExporterMetrics struct {
	exported otelconv.SDKExporterMetricDataPointExported
	inflight otelconv.SDKExporterMetricDataPointInflight
	duration otelconv.SDKExporterOperationDuration
	attrs    []attribute.KeyValue
	enabled  bool
}

// NewExporterMetrics creates a new ExporterMetrics instance.
// If self-observability is disabled, returns a no-op instance.
func NewExporterMetrics(componentType, serverAddress string, serverPort int) *ExporterMetrics {
	em := &ExporterMetrics{
		enabled: isSelfObservabilityEnabled(),
	}

	if !em.enabled {
		return em
	}

	meter := otel.GetMeterProvider().Meter(
		"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error
	em.exported, err = otelconv.NewSDKExporterMetricDataPointExported(meter)
	if err != nil {
		em.enabled = false
		return em
	}

	em.inflight, err = otelconv.NewSDKExporterMetricDataPointInflight(meter)
	if err != nil {
		em.enabled = false
		return em
	}

	em.duration, err = otelconv.NewSDKExporterOperationDuration(meter)
	if err != nil {
		em.enabled = false
		return em
	}

	// Set up common attributes
	componentName := fmt.Sprintf("%s/%d", componentType, nextExporterID())
	em.attrs = []attribute.KeyValue{
		semconv.OTelComponentTypeKey.String(componentType),
		semconv.OTelComponentName(componentName),
		semconv.ServerAddress(serverAddress),
		semconv.ServerPort(serverPort),
	}

	return em
}

// TrackExport tracks an export operation and returns a function to complete the tracking.
// The returned function should be called when the export operation completes.
func (em *ExporterMetrics) TrackExport(ctx context.Context, rm *metricdata.ResourceMetrics) func(error) {
	if !em.enabled {
		return func(error) {}
	}

	dataPointCount := countDataPoints(rm)
	startTime := time.Now()

	em.inflight.Add(ctx, dataPointCount, em.attrs...)

	return func(err error) {
		em.inflight.Add(ctx, -dataPointCount, em.attrs...)

		duration := time.Since(startTime).Seconds()
		attrs := em.attrs
		if err != nil {
			attrs = append(attrs, semconv.ErrorTypeOther)
		}
		em.duration.Inst().Record(ctx, duration, metric.WithAttributes(attrs...))

		if err == nil {
			em.exported.Add(ctx, dataPointCount, em.attrs...)
		}
	}
}

// countDataPoints counts the total number of data points in a ResourceMetrics.
func countDataPoints(rm *metricdata.ResourceMetrics) int64 {
	if rm == nil {
		return 0
	}

	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Gauge[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Summary:
				total += int64(len(data.DataPoints))
			}
		}
	}
	return total
}

// ParseEndpoint extracts server address and port from an endpoint URL.
// Returns defaults if parsing fails.
func ParseEndpoint(endpoint string, defaultPort int) (address string, port int) {
	address = "localhost"
	port = defaultPort

	if endpoint == "" {
		return
	}

	// Handle endpoint without scheme
	if !strings.Contains(endpoint, "://") {
		endpoint = "http://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return
	}

	if u.Hostname() != "" {
		address = u.Hostname()
	}

	if u.Port() != "" {
		if p, err := strconv.Atoi(u.Port()); err == nil {
			port = p
		}
	}

	return
}

// isSelfObservabilityEnabled checks if self-observability is enabled via environment variable.
// It follows OpenTelemetry specification for boolean environment variable parsing.
func isSelfObservabilityEnabled() bool {
	value := os.Getenv("OTEL_GO_X_SELF_OBSERVABILITY")
	// Only "true" (case-insensitive) is considered true, all other values are false
	return strings.EqualFold(value, "true")
}
