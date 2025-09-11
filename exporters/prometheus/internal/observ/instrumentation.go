// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides experimental observability instrumentation
// for the prometheus exporter.
package observ // import "go.opentelemetry.io/otel/exporters/prometheus/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus/internal"
	"go.opentelemetry.io/otel/exporters/prometheus/internal/x"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

const (
	// ComponentType uniquely identifies the OpenTelemetry Exporter component
	// being instrumented.
	ComponentType = "go.opentelemetry.io/otel/exporters/prometheus/prometheus.Exporter"

	// ScopeName is the unique name of the meter used for instrumentation.
	ScopeName = "go.opentelemetry.io/otel/exporters/prometheus/internal/x"

	// SchemaURL is the schema URL of the metrics produced by this
	// instrumentation.
	SchemaURL = semconv.SchemaURL

	// Version is the current version of this instrumentation.
	//
	// This matches the version of the exporter.
	Version = internal.Version
)

var (
	measureAttrsPool = &sync.Pool{
		New: func() any {
			// "component.name" + "component.type" + "error.type"
			const n = 1 + 1 + 1
			s := make([]attribute.KeyValue, 0, n)
			// Return a pointer to a slice instead of a slice itself
			// to avoid allocations on every call.
			return &s
		},
	}

	addOptPool = &sync.Pool{
		New: func() any {
			const n = 1 // WithAttributeSet
			o := make([]metric.AddOption, 0, n)
			return &o
		},
	}

	recordOptPool = &sync.Pool{
		New: func() any {
			const n = 1 // WithAttributeSet
			o := make([]metric.RecordOption, 0, n)
			return &o
		},
	}
)

func get[T any](p *sync.Pool) *[]T { return p.Get().(*[]T) }

func put[T any](p *sync.Pool, s *[]T) {
	*s = (*s)[:0] // Reset.
	p.Put(s)
}

func ComponentName(id int64) string {
	return fmt.Sprintf("%s/%d", ComponentType, id)
}

type Instrumentation struct {
	inflightMetric     metric.Int64UpDownCounter
	exportedMetric     metric.Int64Counter
	operationDuration  metric.Float64Histogram
	collectionDuration metric.Float64Histogram

	attrs  []attribute.KeyValue
	setOpt metric.MeasurementOption
}

func NewInstrumentation(id int64) (*Instrumentation, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	i := &Instrumentation{
		attrs: []attribute.KeyValue{
			semconv.OTelComponentName(ComponentName(id)),
			semconv.OTelComponentTypeKey.String(ComponentType),
		},
	}

	s := attribute.NewSet(i.attrs...)
	i.setOpt = metric.WithAttributeSet(s)

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		ScopeName,
		metric.WithInstrumentationVersion(Version),
		metric.WithSchemaURL(SchemaURL),
	)

	var err, e error

	inflightMetric, e := otelconv.NewSDKExporterMetricDataPointInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	i.inflightMetric = inflightMetric.Inst()

	exportedMetric, e := otelconv.NewSDKExporterMetricDataPointExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create exported metric: %w", e)
		err = errors.Join(err, e)
	}
	i.exportedMetric = exportedMetric.Inst()

	operationDuration, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	i.operationDuration = operationDuration.Inst()

	collectionDuration, e := otelconv.NewSDKMetricReaderCollectionDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create collection duration metric: %w", e)
		err = errors.Join(err, e)
	}
	i.collectionDuration = collectionDuration.Inst()

	return i, err
}

func (i *Instrumentation) RecordOperationDuration(
	ctx context.Context,
) func(err error) {
	start := time.Now()

	return func(err error) {
		recordOpt := get[metric.RecordOption](recordOptPool)
		defer put(recordOptPool, recordOpt)
		*recordOpt = append(*recordOpt, i.setOpt)

		if err != nil {
			attrs := get[attribute.KeyValue](measureAttrsPool)
			defer put(measureAttrsPool, attrs)
			*attrs = append(*attrs, i.attrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			set := attribute.NewSet(*attrs...)
			*recordOpt = append((*recordOpt)[:0], metric.WithAttributeSet(set))
		}

		i.operationDuration.Record(ctx, time.Since(start).Seconds(), *recordOpt...)
	}
}

func (i *Instrumentation) RecordCollectionDuration(
	ctx context.Context,
	operation func() error,
) error {
	start := time.Now()

	recordOpt := get[metric.RecordOption](recordOptPool)
	defer put(recordOptPool, recordOpt)
	*recordOpt = append(*recordOpt, i.setOpt)

	err := operation()
	if err != nil {
		attrs := get[attribute.KeyValue](measureAttrsPool)
		defer put(measureAttrsPool, attrs)
		*attrs = append(*attrs, i.attrs...)
		*attrs = append(*attrs, semconv.ErrorType(err))

		set := attribute.NewSet(*attrs...)
		*recordOpt = append((*recordOpt)[:0], metric.WithAttributeSet(set))
	}

	i.collectionDuration.Record(ctx, time.Since(start).Seconds(), *recordOpt...)

	return err
}

type ScrapeDone func(success int64, err error)

func (i *Instrumentation) ExportMetrics(ctx context.Context, n int64) ExportMetricsDone {
	addOpt := get[metric.AddOption](addOptPool)
	defer put(addOptPool, addOpt)
	*addOpt = append(*addOpt, i.setOpt)

	i.inflightMetric.Add(ctx, n, *addOpt...)

	return i.end(ctx, n)
}

func (i *Instrumentation) end(ctx context.Context, n int64) ScrapeDone {
	return func(success int64, err error) {
		addOpt := get[metric.AddOption](addOptPool)
		defer put(addOptPool, addOpt)
		*addOpt = append(*addOpt, i.setOpt)

		i.inflightMetric.Add(ctx, -n, *addOpt...)
		i.exportedMetric.Add(ctx, success, *addOpt...)

		if err != nil {
			attrs := get[attribute.KeyValue](measureAttrsPool)
			defer put(measureAttrsPool, attrs)
			*attrs = append(*attrs, i.attrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			set := attribute.NewSet(*attrs...)

			*addOpt = append((*addOpt)[:0], metric.WithAttributeSet(set))
			i.exportedMetric.Add(ctx, n-success, *addOpt...)
		}
	}
}
