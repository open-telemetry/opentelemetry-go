// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides observability metrics for OTLP metric exporters.
// This is an experimental feature controlled by the x.Observability feature flag.
package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	// ScopeName is the unique name of the meter used for instrumentation.
	ScopeName = "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/observ"

	// Version is the current version of this instrumentation.
	//
	// This matches the version of the exporter.
	Version = internal.Version
)

var (
	attrsPool = &sync.Pool{
		New: func() any {
			const n = 1 /* component.name */ +
				1 /* component.type */ +
				1 /* server.addr */ +
				1 /* server.port */ +
				1 /* error.type */ +
				1 /* rpc.grpc.status_code */
			s := make([]attribute.KeyValue, 0, n)
			// Return a pointer to a slice instead of a slice itself
			// to avoid allocations on every call.
			return &s
		},
	}
	addOpPool = &sync.Pool{
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

// GetComponentName returns the constant name for the exporter with the
// provided id.
func GetComponentName(id int64) string {
	return fmt.Sprintf("%s/%d", otelconv.ComponentTypeOtlpGRPCMetricExporter, id)
}

// getPresetAttrs builds the preset attributes for instrumentation.
func getPresetAttrs(id int64, target string) []attribute.KeyValue {
	serverAttrs := ServerAddrAttrs(target)
	attrs := make([]attribute.KeyValue, 0, 2+len(serverAttrs))

	attrs = append(
		attrs,
		semconv.OTelComponentName(GetComponentName(id)),
		semconv.OTelComponentTypeOtlpGRPCMetricExporter,
	)
	attrs = append(attrs, serverAttrs...)

	return attrs
}

// Instrumentation is experimental instrumentation for the exporter.
type Instrumentation struct {
	metricDataPointInflightMetric metric.Int64UpDownCounter
	metricDataPointExportedMetric metric.Int64Counter
	operationDurationMetric       metric.Float64Histogram

	presetAttrs []attribute.KeyValue
	addOpt      metric.AddOption
	recOpt      metric.RecordOption
}

// NewInstrumentation returns instrumentation for otlpmetric grpc exporter.
func NewInstrumentation(id int64, target string) (*Instrumentation, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	i := &Instrumentation{}

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		ScopeName,
		metric.WithInstrumentationVersion(Version),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error

	metricDataPointInflightMetric, e := otelconv.NewSDKExporterMetricDataPointInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create metric data point inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	i.metricDataPointInflightMetric = metricDataPointInflightMetric.Inst()

	metricDataPointExportedMetric, e := otelconv.NewSDKExporterMetricDataPointExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create metric data point exported metric: %w", e)
		err = errors.Join(err, e)
	}
	i.metricDataPointExportedMetric = metricDataPointExportedMetric.Inst()

	operationDurationMetric, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	i.operationDurationMetric = operationDurationMetric.Inst()
	if err != nil {
		return nil, err
	}

	i.presetAttrs = getPresetAttrs(id, target)

	i.addOpt = metric.WithAttributeSet(attribute.NewSet(i.presetAttrs...))
	i.recOpt = metric.WithAttributeSet(attribute.NewSet(append(
		// Default to OK status code.
		[]attribute.KeyValue{semconv.RPCGRPCStatusCodeOk},
		i.presetAttrs...,
	)...))
	return i, nil
}

// ExportMetrics instruments the ExportMetrics method of the exporter. It returns
// an [ExportOp] that must have its [ExportOp.End] method called when the
// ExportMetrics method returns.
func (i *Instrumentation) ExportMetrics(ctx context.Context, rm *metricdata.ResourceMetrics) ExportOp {
	count := countDataPoints(rm)
	start := time.Now()
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)

	*addOpt = append(*addOpt, i.addOpt)

	i.metricDataPointInflightMetric.Add(ctx, count, *addOpt...)

	return ExportOp{
		nDataPoints: count,
		ctx:         ctx,
		start:       start,
		inst:        i,
	}
}

// ExportOp tracks the operation being observed by [Instrumentation.ExportMetrics].
type ExportOp struct {
	nDataPoints int64
	ctx         context.Context
	start       time.Time

	inst *Instrumentation
}

// End completes the observation of the operation being observed by a call to
// [Instrumentation.ExportMetrics].
// Any error that is encountered is provided as err.
//
// If err is not nil, all metric data points will be recorded as failures unless error is of
// type [internal.PartialSuccess]. In the case of a PartialSuccess, the number
// of successfully exported metric data points will be determined by inspecting the
// RejectedDataPoints field of the PartialSuccess.
func (e ExportOp) End(err error) {
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)
	*addOpt = append(*addOpt, e.inst.addOpt)

	e.inst.metricDataPointInflightMetric.Add(e.ctx, -e.nDataPoints, *addOpt...)
	success := successful(e.nDataPoints, err)
	e.inst.metricDataPointExportedMetric.Add(e.ctx, success, *addOpt...)

	if err != nil {
		// Add the error.type attribute to the attribute set.
		attrs := get[attribute.KeyValue](attrsPool)
		defer put(attrsPool, attrs)
		*attrs = append(*attrs, e.inst.presetAttrs...)
		*attrs = append(*attrs, semconv.ErrorType(err))

		o := metric.WithAttributeSet(attribute.NewSet(*attrs...))

		// Reset addOpt with new attribute set
		*addOpt = append((*addOpt)[:0], o)

		e.inst.metricDataPointExportedMetric.Add(e.ctx, e.nDataPoints-success, *addOpt...)
	}

	recordOpt := get[metric.RecordOption](recordOptPool)
	defer put(recordOptPool, recordOpt)
	*recordOpt = append(*recordOpt, e.inst.recordOption(err))
	e.inst.operationDurationMetric.Record(e.ctx, time.Since(e.start).Seconds(), *recordOpt...)
}

func (i *Instrumentation) recordOption(err error) metric.RecordOption {
	if err == nil {
		return i.recOpt
	}
	attrs := get[attribute.KeyValue](attrsPool)
	defer put(attrsPool, attrs)

	*attrs = append(*attrs, i.presetAttrs...)
	code := int64(status.Code(err))
	*attrs = append(
		*attrs,
		semconv.RPCGRPCStatusCodeKey.Int64(code),
		semconv.ErrorType(err),
	)

	return metric.WithAttributeSet(attribute.NewSet(*attrs...))
}

// successful returns the number of successfully exported metric data points out of the n
// that were exported based on the provided error.
//
// If err is nil, n is returned. All metric data points were successfully exported.
//
// If err is not nil and not an [internal.PartialSuccess] error, 0 is returned.
// It is assumed all metric data points failed to be exported.
//
// If err is an [internal.PartialSuccess] error, the number of successfully
// exported metric data points is computed by subtracting the RejectedDataPoints field from n. If
// RejectedDataPoints is negative, n is returned. If RejectedDataPoints is greater than
// n, 0 is returned.
func successful(n int64, err error) int64 {
	if err == nil {
		return n // All metric data points successfully exported.
	}
	// Split rejection calculation so successful is inlineable.
	return n - rejectedCount(n, err)
}

var errPool = sync.Pool{
	New: func() any {
		return new(internal.PartialSuccess)
	},
}

// rejectedCount returns how many out of the n metric data points exported were rejected based on
// the provided non-nil err.
func rejectedCount(n int64, err error) int64 {
	ps := errPool.Get().(*internal.PartialSuccess)
	defer errPool.Put(ps)

	// check for partial success
	if errors.As(err, ps) {
		return min(max(ps.RejectedItems, 0), n)
	}
	// all metric data points rejected
	return n
}

// countDataPoints counts the total number of data points in a ResourceMetrics.
func countDataPoints(rm *metricdata.ResourceMetrics) int64 {
	if rm == nil {
		return 0
	}

	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Gauge[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Sum[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.Histogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[int64]:
				total += int64(len(data.DataPoints))
			case metricdata.ExponentialHistogram[float64]:
				total += int64(len(data.DataPoints))
			case metricdata.Summary:
				total += int64(len(data.DataPoints))
			}
		}
	}
	return total
}
