// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package selfobservability // import "go.opentelemetry.io/otel/exporters/prometheus/internal/selfobservability"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

type ExporterMetrics struct {
	selfObservabilityEnabled bool
	inflightMetric           otelconv.SDKExporterMetricDataPointInflight
	exportedMetric           otelconv.SDKExporterMetricDataPointExported
	operationDuration        otelconv.SDKExporterOperationDuration
	collectionDuration       otelconv.SDKMetricReaderCollectionDuration
	attrs                    []attribute.KeyValue
}

func NewExporterMetrics(componentName string) *ExporterMetrics {
	em := &ExporterMetrics{}
	em.selfObservabilityEnabled = true

	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/exporters/prometheus",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	if em.inflightMetric, err = otelconv.NewSDKExporterMetricDataPointInflight(m); err != nil {
		otel.Handle(err)
	}
	if em.exportedMetric, err = otelconv.NewSDKExporterMetricDataPointExported(m); err != nil {
		otel.Handle(err)
	}
	if em.operationDuration, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}
	if em.collectionDuration, err = otelconv.NewSDKMetricReaderCollectionDuration(m); err != nil {
		otel.Handle(err)
	}

	em.attrs = []attribute.KeyValue{
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypePrometheusHTTPTextMetricExporter)),
	}

	return em
}

// DisableSelfObservability disables self-observability for testing purposes.
func (em *ExporterMetrics) DisableSelfObservability() {
	em.selfObservabilityEnabled = false
}

// AddInflight adds the specified count to the inflight metric.
func (em *ExporterMetrics) AddInflight(ctx context.Context, count int64) {
	if !em.selfObservabilityEnabled {
		return
	}
	em.inflightMetric.Add(ctx, count, em.attrs...)
}

// AddExported adds the specified count to the exported metric.
func (em *ExporterMetrics) AddExported(ctx context.Context, count int64) {
	if !em.selfObservabilityEnabled {
		return
	}
	em.exportedMetric.Add(ctx, count, em.attrs...)
}

// AddExportedWithError adds the specified count to the exported metric with error attributes.
// This method is specifically for tracking failed exports, so an error must be provided.
func (em *ExporterMetrics) AddExportedWithError(ctx context.Context, count int64, err error) {
	if !em.selfObservabilityEnabled {
		return
	}
	attrs := append(em.attrs, semconv.ErrorType(err))
	em.exportedMetric.Add(ctx, count, attrs...)
}

// TrackCollectionDuration records the duration of a reader collection operation.
func (em *ExporterMetrics) TrackCollectionDuration(ctx context.Context) func(error) {
	if !em.selfObservabilityEnabled {
		return func(error) {}
	}
	start := time.Now()
	return func(err error) {
		duration := time.Since(start).Seconds()
		attrs := em.attrs
		if err != nil {
			attrs = append(attrs, semconv.ErrorType(err))
		}
		em.collectionDuration.Inst().Record(ctx, duration, metric.WithAttributes(attrs...))
	}
}

// TrackOperationDuration records the duration of an exporter operation (full scrape/export path).
func (em *ExporterMetrics) TrackOperationDuration(ctx context.Context) func(error) {
	if !em.selfObservabilityEnabled {
		return func(error) {}
	}
	start := time.Now()
	return func(err error) {
		duration := time.Since(start).Seconds()
		attrs := em.attrs
		if err != nil {
			attrs = append(attrs, semconv.ErrorType(err))
		}
		em.operationDuration.Inst().Record(ctx, duration, metric.WithAttributes(attrs...))
	}
}
