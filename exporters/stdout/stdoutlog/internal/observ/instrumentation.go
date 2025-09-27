// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides experimental observability instrumentation for the
// stdout log exporter.
package observ // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/x"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

// InstrumentationVersion matches the stdout log exporter version.
const InstrumentationVersion = internal.Version

var (
	attrsPool = &sync.Pool{
		New: func() any {
			// component.name + component.type + error.type
			const n = 1 + 1 + 1
			s := make([]attribute.KeyValue, 0, n)
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
	*s = (*s)[:0]
	p.Put(s)
}

// Instrumentation instruments the stdout log exporter.
type Instrumentation struct {
	inflight metric.Int64UpDownCounter
	exported metric.Int64Counter
	duration metric.Float64Histogram

	attrs  []attribute.KeyValue
	setOpt metric.MeasurementOption
}

// ExportLogsDone completes an export observation.
type ExportLogsDone func(success int64, err error)

// NewInstrumentation returns instrumentation for the stdout log exporter with
// the provided component type and exporter identifier using the global
// MeterProvider.
//
// If the experimental observability feature is disabled, nil is returned.
func NewInstrumentation(componentType string, exporterID int64) (*Instrumentation, error) {
	if !x.SelfObservability.Enabled() {
		return nil, nil
	}

	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(fmt.Sprintf("%s/%d", componentType, exporterID)),
		semconv.OTelComponentTypeKey.String(componentType),
	}

	inst := &Instrumentation{
		attrs: attrs,
	}

	set := attribute.NewSet(attrs...)
	inst.setOpt = metric.WithAttributeSet(set)

	mp := otel.GetMeterProvider()
	meter := mp.Meter(
		componentType,
		metric.WithInstrumentationVersion(InstrumentationVersion),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error

	inflight, e := otelconv.NewSDKExporterLogInflight(meter)
	if e != nil {
		err = errors.Join(err, fmt.Errorf("failed to create log inflight metric: %w", e))
	}
	inst.inflight = inflight.Inst()

	exported, e := otelconv.NewSDKExporterLogExported(meter)
	if e != nil {
		err = errors.Join(err, fmt.Errorf("failed to create log exported metric: %w", e))
	}
	inst.exported = exported.Inst()

	duration, e := otelconv.NewSDKExporterOperationDuration(meter)
	if e != nil {
		err = errors.Join(err, fmt.Errorf("failed to create export duration metric: %w", e))
	}
	inst.duration = duration.Inst()

	return inst, err
}

// ExportLogs instruments the exporter Export method. It returns a callback that
// MUST be invoked when the export completes with the number of successfully
// exported records and the resulting error.
func (i *Instrumentation) ExportLogs(ctx context.Context, total int) ExportLogsDone {
	start := time.Now()

	addOpt := get[metric.AddOption](addOptPool)
	*addOpt = append(*addOpt, i.setOpt)
	i.inflight.Add(ctx, int64(total), *addOpt...)
	put(addOptPool, addOpt)

	return func(success int64, err error) {
		addOpt := get[metric.AddOption](addOptPool)
		defer put(addOptPool, addOpt)
		*addOpt = append(*addOpt, i.setOpt)

		n := int64(total)
		i.inflight.Add(ctx, -n, *addOpt...)
		if success > 0 || (n == 0 && err == nil) {
			i.exported.Add(ctx, success, *addOpt...)
		}

		measurementOpt := i.setOpt

		if err != nil {
			attrs := get[attribute.KeyValue](attrsPool)
			defer put(attrsPool, attrs)
			*attrs = append(*attrs, i.attrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			set := attribute.NewSet(*attrs...)
			measurementOpt = metric.WithAttributeSet(set)

			*addOpt = append((*addOpt)[:0], measurementOpt)
			failures := n - success
			if failures < 0 {
				failures = 0
			}
			if failures > 0 || (n == 0 && err != nil) {
				i.exported.Add(ctx, failures, *addOpt...)
			}
		}

		recordOpt := get[metric.RecordOption](recordOptPool)
		defer put(recordOptPool, recordOpt)
		*recordOpt = append(*recordOpt, measurementOpt)

		i.duration.Record(ctx, time.Since(start).Seconds(), *recordOpt...)
	}
}
