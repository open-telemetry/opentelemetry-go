// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for stdout metric exporter.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/selfobservability"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

var measureAttrsPool = sync.Pool{
	New: func() any {
		// "component.name" + "component.type" + "error.type"
		const n = 1 + 1 + 1
		s := make([]attribute.KeyValue, 0, n)
		// Return a pointer to a slice instead of a slice itself
		// to avoid allocations on every call.
		return &s
	},
}

type ExporterMetrics struct {
	inflight otelconv.SDKExporterMetricDataPointInflight
	exported otelconv.SDKExporterMetricDataPointExported
	duration otelconv.SDKExporterOperationDuration
	attrs    []attribute.KeyValue
}

func NewExporterMetrics(
	name string,
	componentName, componentType attribute.KeyValue,
) (*ExporterMetrics, error) {
	em := &ExporterMetrics{
		attrs: []attribute.KeyValue{componentName, componentType},
	}
	mp := otel.GetMeterProvider()
	m := mp.Meter(
		name,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err, e error
	if em.inflight, e = otelconv.NewSDKExporterMetricDataPointInflight(m); e != nil {
		e = fmt.Errorf("failed to create metric_data_point inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	if em.exported, e = otelconv.NewSDKExporterMetricDataPointExported(m); e != nil {
		e = fmt.Errorf("failed to create metric_data_point exported metric: %w", e)
		err = errors.Join(err, e)
	}
	if em.duration, e = otelconv.NewSDKExporterOperationDuration(m); e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	return em, err
}

func (em *ExporterMetrics) TrackExport(ctx context.Context, count int64) func(err error) {
	begin := time.Now()
	em.inflight.Add(ctx, count, em.attrs...)
	return func(err error) {
		durationSeconds := time.Since(begin).Seconds()
		attrs := &em.attrs
		em.inflight.Add(ctx, -count, *attrs...)
		if err != nil {
			attrs = measureAttrsPool.Get().(*[]attribute.KeyValue)
			defer func() {
				*attrs = (*attrs)[:0] // reset the slice for reuse
				measureAttrsPool.Put(attrs)
			}()
			copy(*attrs, em.attrs)
			*attrs = append(*attrs, semconv.ErrorType(err))
		}
		em.exported.Add(ctx, count, *attrs...)
		em.duration.Record(ctx, durationSeconds, *attrs...)
	}
}
