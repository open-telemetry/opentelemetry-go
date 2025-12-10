// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
	"sync"
	"time"
)

const (
	// ScopeName is the unique name of the meter used for instrumentation.
	ScopeName = "go.opentelemetry.io/otel/exporters/stdoutlog/internal/observ"

	// ComponentType uniquely identifies the OpenTelemetry Exporter component
	// being instrumented.
	//
	// The STDOUT log exporter is not a standardized OTel component type, so
	// it uses the Go package prefixed type name to ensure uniqueness and
	// identity.
	ComponentType = "go.opentelemetry.io/otel/exporters/stdout/stdoutlog.Exporter"
)

var (
	addOptPool = &sync.Pool{
		New: func() any {
			const n = 1
			s := make([]metric.AddOption, 0, n)
			return &s
		},
	}
)

func get[T any](pool *sync.Pool) *[]T {
	return pool.Get().(*[]T)
}

func put[T any](pool *sync.Pool, value *[]T) {
	*value = (*value)[:0]
	pool.Put(value)
}

// Instrumentation is experimental instrumentation for the exporter.
type Instrumentation struct {
	inflight metric.Int64UpDownCounter
	exported metric.Int64Counter
	duration metric.Float64Histogram

	attrs  []attribute.KeyValue
	addOpt metric.AddOption
	recOpt metric.RecordOption
}

// GetComponentName returns the constant name for the exporter with the
// provided id.
func GetComponentName(id int64) string {
	return fmt.Sprintf("%s/%d", ComponentType, id)
}

// NewInstrumentation .
func NewInstrumentation(id int64) (*Instrumentation, error) {
	inst := &Instrumentation{}

	mp := otel.GetMeterProvider()
	m := mp.Meter(ScopeName)

	var err error

	inflight, e := otelconv.NewSDKExporterLogInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create the inflight metirc %w", e)
		err = errors.Join(err, e)
	}
	inst.inflight = inflight.Inst()

	exported, e := otelconv.NewSDKExporterLogExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create the exported metric %w", e)
		err = errors.Join(err, e)
	}
	inst.exported = exported.Inst()

	duration, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create the duration metric %w", e)
		err = errors.Join(err, e)
	}
	inst.duration = duration.Inst()

	if err != nil {
		return nil, err
	}

	return inst, nil
}

func (i *Instrumentation) Export(ctx context.Context, count int64) ExportOp {
	start := time.Now()

	addOpt := get[metric.AddOption](addOptPool)
	defer put(addOptPool, addOpt)
	*addOpt = append(*addOpt, i.addOpt)

	i.inflight.Add(ctx, count, *addOpt...)

	return ExportOp{
		count: count,
		ctx:   ctx,
		inst:  i,
		start: start,
	}
}

type ExportOp struct {
	count int64
	ctx   context.Context
	inst  *Instrumentation
	start time.Time
}

func (e *ExportOp) ExportLogs(err error) {
	addOpt := get[metric.AddOption](addOptPool)
	defer put(addOptPool, addOpt)
	*addOpt = append(*addOpt, e.inst.addOpt)

	e.inst.inflight.Add(e.ctx, e.count, *addOpt...)
}

func successful(err error, n int64) int64 {
	if err != nil {
		return n
	}
	return n - failed(err)
}

var errPool = sync.Pool{}

func failed(err error) int64 {
	return 0
}
