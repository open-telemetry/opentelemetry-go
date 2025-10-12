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
	inflight        otelconv.SDKExporterMetricDataPointInflight
	inflightCounter metric.Int64UpDownCounter
	addOpts         []metric.AddOption
	exported        otelconv.SDKExporterMetricDataPointExported
	duration        otelconv.SDKExporterOperationDuration
	recordOpts      []metric.RecordOption
	attrs           []attribute.KeyValue
	set             attribute.Set
}

func NewExporterMetrics(
	name string,
	componentName, componentType attribute.KeyValue,
) (*ExporterMetrics, error) {
	attrs := []attribute.KeyValue{componentName, componentType}
	attrSet := attribute.NewSet(attrs...)
	attrOpts := metric.WithAttributeSet(attrSet)
	addOpts := []metric.AddOption{attrOpts}
	recordOpts := []metric.RecordOption{attrOpts}
	em := &ExporterMetrics{
		attrs:      attrs,
		addOpts:    addOpts,
		set:        attrSet,
		recordOpts: recordOpts,
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
	em.inflightCounter = em.inflight.Int64UpDownCounter
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
	em.inflightCounter.Add(ctx, count, em.addOpts...)
	return func(err error) {
		durationSeconds := time.Since(begin).Seconds()
		em.inflightCounter.Add(ctx, -count, em.addOpts...)
		if err == nil {
			em.exported.Int64Counter.Add(ctx, count, em.addOpts...)
			em.duration.Float64Histogram.Record(ctx, durationSeconds, em.recordOpts...)
			return
		}

		attrs := measureAttrsPool.Get().(*[]attribute.KeyValue)
		defer func() {
			*attrs = (*attrs)[:0] // reset the slice for reuse
			measureAttrsPool.Put(attrs)
		}()
		*attrs = append(*attrs, em.attrs...)
		*attrs = append(*attrs, semconv.ErrorType(err))

		set := attribute.NewSet(*attrs...)
		em.exported.AddSet(ctx, count, set)
		em.duration.RecordSet(ctx, durationSeconds, set)
	}
}
