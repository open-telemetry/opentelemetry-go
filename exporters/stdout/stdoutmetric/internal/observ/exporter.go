// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides observability for stdout metric exporter.
// This is an experimental feature controlled by the x.Observability feature flag.
package observ // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

const (
	scope = "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

	// componentType is a name identifying the type of the OpenTelemetry
	// component. It is not a standardized OTel component type, so it uses the
	// Go package prefixed type name to ensure uniqueness and identity.
	componentType = "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric.exporter"
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

// Instrumentation is the instrumentation for stdout metric exporter.
type Instrumentation struct {
	inflight   metric.Int64UpDownCounter
	addOpts    []metric.AddOption
	exported   otelconv.SDKExporterMetricDataPointExported
	duration   otelconv.SDKExporterOperationDuration
	recordOpts []metric.RecordOption
	attrs      []attribute.KeyValue
}

func exporterComponentName(id int64) attribute.KeyValue {
	componentName := fmt.Sprintf("%s/%d", componentType, id)
	return semconv.OTelComponentName(componentName)
}

// NewInstrumentation returns a new Instrumentation for the stdout metric exporter
// with the provided ID.
//
// If the experimental observability is disabled, nil is returned.
func NewInstrumentation(id int64) (*Instrumentation, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}
	attrs := []attribute.KeyValue{
		exporterComponentName(id),
		semconv.OTelComponentTypeKey.String(componentType),
	}
	attrOpts := metric.WithAttributeSet(attribute.NewSet(attrs...))
	addOpts := []metric.AddOption{attrOpts}
	recordOpts := []metric.RecordOption{attrOpts}
	em := &Instrumentation{
		attrs:      attrs,
		addOpts:    addOpts,
		recordOpts: recordOpts,
	}
	mp := otel.GetMeterProvider()
	m := mp.Meter(
		scope,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	inflightMetric, e := otelconv.NewSDKExporterMetricDataPointInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create metric_data_point inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	em.inflight = inflightMetric.Int64UpDownCounter
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

func (em *Instrumentation) TrackExport(ctx context.Context, count int64) func(err error) {
	begin := time.Now()
	em.inflight.Add(ctx, count, em.addOpts...)
	return func(err error) {
		durationSeconds := time.Since(begin).Seconds()
		em.inflight.Add(ctx, -count, em.addOpts...)
		if err == nil { // short circuit in case of success to avoid allocations
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
