// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides experimental observability instrumentation
// for the stdout trace exporter.
package observ // import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	// ComponentType uniquely identifies the OpenTelemetry Exporter component
	// being instrumented.
	//
	// The STDOUT trace exporter is not a standardized OTel component type, so
	// it uses the Go package prefixed type name to ensure uniqueness and
	// identity.
	ComponentType = "go.opentelemetry.io/otel/exporters/stdout/stdouttrace.Exporter"

	// ScopeName is the unique name of the meter used for instrumentation.
	ScopeName = "go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/observ"

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

// Instrumentation is experimental instrumentation for the exporter.
type Instrumentation struct {
	inflightSpans metric.Int64UpDownCounter
	exportedSpans metric.Int64Counter
	opDuration    metric.Float64Histogram

	attrs  []attribute.KeyValue
	setOpt metric.MeasurementOption
}

// NewInstrumentation returns instrumentation for a STDOUT trace exporter with
// the provided ID using the global MeterProvider.
//
// If the experimental observability is disabled, nil is returned.
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

	var err error

	inflightSpans, e := otelconv.NewSDKExporterSpanInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create span inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	i.inflightSpans = inflightSpans.Inst()

	exportedSpans, e := otelconv.NewSDKExporterSpanExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create span exported metric: %w", e)
		err = errors.Join(err, e)
	}
	i.exportedSpans = exportedSpans.Inst()

	opDuration, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	i.opDuration = opDuration.Inst()

	return i, err
}

// ExportSpansDone is a function that is called when a call to an Exporter's
// ExportSpans method completes.
//
// The number of successful exports is provided as success. Any error that is
// encountered is provided as err.
type ExportSpansDone func(success int64, err error)

// ExportSpans instruments the ExportSpans method of the exporter. It returns a
// function that needs to be deferred so it is called when the method returns.
func (i *Instrumentation) ExportSpans(ctx context.Context, nSpans int) ExportSpansDone {
	start := time.Now()

	addOpt := get[metric.AddOption](addOptPool)
	defer put(addOptPool, addOpt)
	*addOpt = append(*addOpt, i.setOpt)
	i.inflightSpans.Add(ctx, int64(nSpans), *addOpt...)

	return i.end(ctx, start, int64(nSpans))
}

func (i *Instrumentation) end(ctx context.Context, start time.Time, n int64) ExportSpansDone {
	return func(success int64, err error) {
		addOpt := get[metric.AddOption](addOptPool)
		defer put(addOptPool, addOpt)
		*addOpt = append(*addOpt, i.setOpt)

		i.inflightSpans.Add(ctx, -n, *addOpt...)

		// Record the success and duration of the operation.
		//
		// Do not exclude 0 values, as they are valid and indicate no spans
		// were exported which is meaningful for certain aggregations.
		i.exportedSpans.Add(ctx, success, *addOpt...)

		mOpt := i.setOpt
		if err != nil {
			attrs := get[attribute.KeyValue](measureAttrsPool)
			defer put(measureAttrsPool, attrs)
			*attrs = append(*attrs, i.attrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			// Do not inefficiently make a copy of attrs by using
			// WithAttributes instead of WithAttributeSet.
			set := attribute.NewSet(*attrs...)
			mOpt = metric.WithAttributeSet(set)

			// Reset addOpt with new attribute set.
			*addOpt = append((*addOpt)[:0], mOpt)

			i.exportedSpans.Add(ctx, n-success, *addOpt...)
		}

		recordOpt := get[metric.RecordOption](recordOptPool)
		defer put(recordOptPool, recordOpt)
		*recordOpt = append(*recordOpt, mOpt)
		i.opDuration.Record(ctx, time.Since(start).Seconds(), *recordOpt...)
	}
}
