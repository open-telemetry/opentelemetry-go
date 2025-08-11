// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for stdout metric exporter.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/selfobservability"

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
	inflight otelconv.SDKExporterMetricDataPointInflight
	exported otelconv.SDKExporterMetricDataPointExported
	duration otelconv.SDKExporterOperationDuration
	attrs    []attribute.KeyValue
}

func NewExporterMetrics(
	name string,
	componentName, componentType attribute.KeyValue,
) *ExporterMetrics {
	em := &ExporterMetrics{}
	mp := otel.GetMeterProvider()
	m := mp.Meter(
		name,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	if em.inflight, err = otelconv.NewSDKExporterMetricDataPointInflight(m); err != nil {
		otel.Handle(err)
	}
	if em.exported, err = otelconv.NewSDKExporterMetricDataPointExported(m); err != nil {
		otel.Handle(err)
	}
	if em.duration, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}

	em.attrs = []attribute.KeyValue{componentName, componentType}
	return em
}

func (em *ExporterMetrics) TrackExport(ctx context.Context, count int64) func(err error) {
	begin := time.Now()
	em.inflight.Add(ctx, count, em.attrs...)
	return func(err error) {
		durationSeconds := time.Since(begin).Seconds()
		attrs := em.attrs
		em.inflight.Add(ctx, -count, attrs...)
		if err != nil {
			attrs = make([]attribute.KeyValue, len(em.attrs)+1)
			copy(attrs, em.attrs)
			attrs = append(attrs, semconv.ErrorType(err))
		}
		em.exported.Add(ctx, count, attrs...)
		em.duration.Record(ctx, durationSeconds, attrs...)
	}
}
