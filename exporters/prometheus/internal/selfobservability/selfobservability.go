// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for prometheus exporter.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/prometheus/internal/selfobservability"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus/internal/counter"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

var otelComponentType = string(otelconv.ComponentTypePrometheusHTTPTextMetricExporter)

var attrsPool = sync.Pool{
	New: func() any {
		// "component.name" + "component.type" + "error.type"
		const n = 1 + 1 + 1
		s := make([]attribute.KeyValue, 0, n)
		return &s
	},
}

type SelfObservability struct {
	ctx                context.Context
	attrs              []attribute.KeyValue
	inflightMetric     otelconv.SDKExporterMetricDataPointInflight
	exportedMetric     otelconv.SDKExporterMetricDataPointExported
	operationDuration  otelconv.SDKExporterOperationDuration
	collectionDuration otelconv.SDKMetricReaderCollectionDuration
}

func NewSelfObservability() (*SelfObservability, error) {
	selfObs := &SelfObservability{}

	componentName := fmt.Sprintf("%s/%d", otelComponentType, counter.NextExporterID())
	selfObs.attrs = []attribute.KeyValue{
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(otelComponentType),
	}

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		"go.opentelemetry.io/otel/exporters/prometheus",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err, e error
	if selfObs.inflightMetric, e = otelconv.NewSDKExporterMetricDataPointInflight(m); e != nil {
		e = fmt.Errorf("failed to create inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	if selfObs.exportedMetric, e = otelconv.NewSDKExporterMetricDataPointExported(m); e != nil {
		e = fmt.Errorf("failed to create exported metric: %w", e)
		err = errors.Join(err, e)
	}
	if selfObs.operationDuration, e = otelconv.NewSDKExporterOperationDuration(m); e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	if selfObs.collectionDuration, e = otelconv.NewSDKMetricReaderCollectionDuration(m); e != nil {
		e = fmt.Errorf("failed to create collection duration metric: %w", e)
		err = errors.Join(err, e)
	}

	return selfObs, err
}

func (obs *SelfObservability) SetContext(ctx context.Context) {
	obs.ctx = ctx
}

func (obs *SelfObservability) GetContext() context.Context {
	return obs.ctx
}

func (obs *SelfObservability) RecordCollectionDuration(
	ctx context.Context,
	operation func() error,
) error {
	begin := time.Now()

	attrs := attrsPool.Get().(*[]attribute.KeyValue)
	*attrs = append([]attribute.KeyValue{}, obs.attrs...)
	defer func() {
		*attrs = (*attrs)[:0]
		attrsPool.Put(attrs)
	}()

	err := operation()
	if err != nil {
		*attrs = append(*attrs, semconv.ErrorType(err))
	}

	duration := time.Since(begin).Seconds()
	obs.collectionDuration.RecordSet(ctx, duration, attribute.NewSet(*attrs...))

	return err
}

func (obs *SelfObservability) RecordOperationDuration(
	ctx context.Context,
) func(err error) {
	begin := time.Now()

	attrs := attrsPool.Get().(*[]attribute.KeyValue)
	*attrs = append([]attribute.KeyValue{}, obs.attrs...)

	return func(err error) {
		defer func() {
			*attrs = (*attrs)[:0]
			attrsPool.Put(attrs)
		}()

		if err != nil {
			*attrs = append(*attrs, semconv.ErrorType(err))
		}
		duration := time.Since(begin).Seconds()
		obs.operationDuration.RecordSet(ctx, duration, attribute.NewSet(*attrs...))
	}
}

func (obs *SelfObservability) TrackExport(
	ctx context.Context,
	count int64,
) func(err error, successCount int64) {
	attrs := attrsPool.Get().(*[]attribute.KeyValue)
	*attrs = append([]attribute.KeyValue{}, obs.attrs...)

	obs.inflightMetric.AddSet(ctx, count, attribute.NewSet(*attrs...))
	return func(err error, successCount int64) {
		defer func() {
			*attrs = (*attrs)[:0]
			attrsPool.Put(attrs)
		}()

		obs.inflightMetric.AddSet(ctx, -count, attribute.NewSet(*attrs...))
		obs.exportedMetric.AddSet(ctx, successCount, attribute.NewSet(*attrs...))

		if err != nil {
			*attrs = append(*attrs, semconv.ErrorType(err))
			obs.exportedMetric.AddSet(ctx, count-successCount, attribute.NewSet(*attrs...))
		}
	}
}
