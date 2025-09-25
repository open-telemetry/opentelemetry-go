// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides observability metrics for OTLP log exporters.
// This is an experimental feature controlled by the x.Observability feature flag.
package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/observ"

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/x"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	// ScopeName is the unique name of the meter used for instrumentation.
	ScopeName = "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/observ"

	// Version is the current version of this instrumentation.
	//
	// This matches the version of the exporter.
	Version = internal.Version
)

var (
	attrsPool = &sync.Pool{
		New: func() any {
			// "component.name" + "component.type" + "error.type" + "server.address" + "server.port"
			const n = 1 + 1 + 1 + 1 + 1
			s := make([]attribute.KeyValue, 0, n)
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
			const n = 1 + 1 // WithAttributeSet + "rpc.grpc.status_code"
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

func GetComponentName(id int64) string {
	return fmt.Sprintf("%s/%d", otelconv.ComponentTypeOtlpGRPCLogExporter, id)
}

// Instrumentation is experimental instrumentation for the exporter.
type Instrumentation struct {
	logInflightMetric         metric.Int64UpDownCounter
	logExportedMetric         metric.Int64Counter
	logExportedDurationMetric metric.Float64Histogram
	presetAttrs               []attribute.KeyValue
	setOpt                    metric.MeasurementOption
}

// NewInstrumentation returns instrumentation for otlplog grpc exporter.
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

	logInflightMetric, e := otelconv.NewSDKExporterLogInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create log inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	i.logInflightMetric = logInflightMetric.Inst()

	logExportedMetric, e := otelconv.NewSDKExporterLogExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create log exported metric: %w", e)
		err = errors.Join(err, e)
	}
	i.logExportedMetric = logExportedMetric.Inst()

	logOpDurationMetric, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create log operation duration metric: %w", e)
		err = errors.Join(err, e)
	}
	i.logExportedDurationMetric = logOpDurationMetric.Inst()
	if err != nil {
		return nil, err
	}

	i.presetAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(id)),
		semconv.OTelComponentTypeOtlpGRPCLogExporter,
	}
	i.presetAttrs = append(i.presetAttrs, ServerAddrAttrs(target)...)
	s := attribute.NewSet(i.presetAttrs...)
	i.setOpt = metric.WithAttributeSet(s)

	return i, nil
}

// ExportLogs instruments the ExportLogs method of the exporter. It returns a
// function that needs to be deferred so it is called when the method returns.
func (i *Instrumentation) ExportLogs(ctx context.Context, count int64) ExportOp {
	start := time.Now()
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)

	*addOpt = append(*addOpt, i.setOpt)

	i.logInflightMetric.Add(ctx, count, *addOpt...)

	return ExportOp{
		nLogs: count,
		ctx:   ctx,
		start: start,
		inst:  i,
	}
}

// ExportOp tracks the operation being observed by [Instrumentation.ExportLogs].
type ExportOp struct {
	nLogs int64
	ctx   context.Context
	start time.Time

	inst *Instrumentation
}

// End completes the observation of the operation being observed by a call to
// [Instrumentation.ExportLogs].
// Any error that is encountered is provided as err.
//
// If err is not nil, all logs will be recorded as failures unless error is of
// type [internal.PartialSuccess]. In the case of a PartialSuccess, the number
// of successfully exported logs will be determined by inspecting the
// RejectedItems field of the PartialSuccess.
func (e ExportOp) End(err error) {
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)

	*addOpt = append(*addOpt, e.inst.setOpt)

	success := successful(e.nLogs, err)

	e.inst.logInflightMetric.Add(e.ctx, -e.nLogs, *addOpt...)
	e.inst.logExportedMetric.Add(e.ctx, success, *addOpt...)

	mOpt := e.inst.setOpt
	if err != nil {
		// Add the error.type attribute to the attribute set.
		attrs := get[attribute.KeyValue](attrsPool)
		defer put(attrsPool, attrs)
		*attrs = append(*attrs, e.inst.presetAttrs...)
		*attrs = append(*attrs, semconv.ErrorType(err))

		set := attribute.NewSet(*attrs...)
		mOpt = metric.WithAttributeSet(set)

		// Reset addOpt with new attribute set
		*addOpt = append((*addOpt)[:0], mOpt)

		e.inst.logExportedMetric.Add(e.ctx, e.nLogs-success, *addOpt...)
	}

	code := status.Code(err)

	recordOpt := get[metric.RecordOption](recordOptPool)
	defer put(recordOptPool, recordOpt)
	*recordOpt = append(
		*recordOpt,
		mOpt,
		metric.WithAttributes(
			semconv.RPCGRPCStatusCodeKey.Int64(int64(code)),
		),
	)
	duration := time.Since(e.start).Seconds()
	e.inst.logExportedDurationMetric.Record(e.ctx, duration, *recordOpt...)
}

// successful returns the number of successfully exported logs out of the n
// that were exported based on the provided error.
//
// If err is nil, n is returned. All logs were successfully exported.
//
// If err is not nil and not an [internal.PartialSuccess] error, 0 is returned.
// It is assumed all logs failed to be exported.
//
// If err is an [internal.PartialSuccess] error, the number of successfully
// exported logs is computed by subtracting the RejectedItems field from n. If
// RejectedItems is negative, n is returned. If RejectedItems is greater than
// n, 0 is returned.
func successful(n int64, err error) int64 {
	if err == nil {
		return n // All log successfully exported.
	}
	// Split rejection calculation so successful is inlineable.
	return n - rejectedCount(n, err)
}

var errPool = sync.Pool{
	New: func() any {
		return new(internal.PartialSuccess)
	},
}

// rejectedCount returns how many out of the n logs exporter were rejected based on
// the provided non-nil err.
func rejectedCount(n int64, err error) int64 {
	ps := errPool.Get().(*internal.PartialSuccess)
	defer errPool.Put(ps)

	// check for partial success
	if errors.As(err, ps) {
		return min(max(ps.RejectedItems, 0), n)
	}
	// all logs exporter
	return n
}

// ServerAddrAttrs is a function that extracts server address and port attributes
// from a target string.
func ServerAddrAttrs(target string) []attribute.KeyValue {
	if !strings.Contains(target, "://") {
		return splitHostPortAttrs(target)
	}

	u, err := url.Parse(target)
	if err != nil || u.Scheme == "" {
		return splitHostPortAttrs(target)
	}

	switch u.Scheme {
	case "unix":
		// unix:///path/to/socket
		return []attribute.KeyValue{semconv.ServerAddress(u.Path)}
	case "dns":
		// dns:///example.com:42 or dns://8.8.8.8/example.com:42
		addr := u.Opaque
		if addr == "" {
			addr = strings.TrimPrefix(u.Path, "/")
		}
		return splitHostPortAttrs(addr)
	default:
		return splitHostPortAttrs(u.Host)
	}
}

func splitHostPortAttrs(target string) []attribute.KeyValue {
	host, pStr, err := net.SplitHostPort(target)
	if err != nil {
		return []attribute.KeyValue{semconv.ServerAddress(target)}
	}
	port, err := strconv.Atoi(pStr)
	if err != nil {
		return []attribute.KeyValue{semconv.ServerAddress(host)}
	}
	return []attribute.KeyValue{
		semconv.ServerAddress(host),
		semconv.ServerPort(port),
	}
}
