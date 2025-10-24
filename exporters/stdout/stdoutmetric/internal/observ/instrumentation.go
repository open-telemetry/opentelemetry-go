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
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/x"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

const (
	scope = "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

	// componentType is a name identifying the type of the OpenTelemetry
	// component. It is not a standardized OTel component type, so it uses the
	// Go package prefixed type name to ensure uniqueness and identity.
	componentType = "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric.exporter"

	// Version is the current version of this instrumentation.
	//
	// This matches the version of the exporter.
	Version = internal.Version
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
		metric.WithInstrumentationVersion(Version),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	inflightMetric, e := otelconv.NewSDKExporterMetricDataPointInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create metric_data_point inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	em.inflight = inflightMetric.Inst()
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

// ExportMetrics instruments the Export method of the exporter. It returns a
// function that needs to be deferred so it is called when the method returns.
func (i *Instrumentation) ExportMetrics(ctx context.Context, count int64) ExportOp {
	start := time.Now()
	i.inflight.Add(ctx, count, i.addOpts...)
	return ExportOp{
		ctx:     ctx,
		start:   start,
		nTraces: count,
		inst:    i,
	}
}

// ExportOp is an in-progress ExportMetrics operation.
type ExportOp struct {
	ctx     context.Context
	start   time.Time
	nTraces int64
	inst    *Instrumentation
}

// End ends the ExportMetrics operation, recording its duration.
//
// The err parameter indicates whether the operation failed. If err is not nil,
// it is added to metrics as attribute.
func (e ExportOp) End(err error) {
	durationSeconds := time.Since(e.start).Seconds()
	e.inst.inflight.Add(e.ctx, -e.nTraces, e.inst.addOpts...)
	if err == nil { // short circuit in case of success to avoid allocations
		e.inst.exported.Inst().Add(e.ctx, e.nTraces, e.inst.addOpts...)
		e.inst.duration.Inst().Record(e.ctx, durationSeconds, e.inst.recordOpts...)
		return
	}

	attrs := measureAttrsPool.Get().(*[]attribute.KeyValue)
	defer func() {
		*attrs = (*attrs)[:0] // reset the slice for reuse
		measureAttrsPool.Put(attrs)
	}()
	*attrs = append(*attrs, e.inst.attrs...)
	*attrs = append(*attrs, semconv.ErrorType(err))

	set := attribute.NewSet(*attrs...)
	e.inst.exported.AddSet(e.ctx, e.nTraces, set)
	e.inst.duration.RecordSet(e.ctx, durationSeconds, set)
}
