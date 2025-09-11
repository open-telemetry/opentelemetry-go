// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package observ provides self-observability metrics for OTLP log exporters.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
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

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

// ScopeName is the unique name of the meter used for instrumentation.
const ScopeName = "go.opentelemetry.io/otel/exporters/otlp/otlpgrpclog/internal/observ"

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
	i := &Instrumentation{}

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		ScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
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
		semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpGRPCLogExporter)),
	}
	i.presetAttrs = append(i.presetAttrs, ServerAddrAttrs(target)...)
	s := attribute.NewSet(i.presetAttrs...)
	i.setOpt = metric.WithAttributeSet(s)

	return i, nil
}

// ExportLogsDone is a function that is called when a call to an Exporter's
// ExportLogs method completes
//
// The number of successful exports is provided as success. Any error that is encountered is provided as error
// The code of last gRPC requests performed in scope of this export call.
type ExportLogsDone func(err error, success int64, code codes.Code)

// ExportLogs instruments the ExportLogs method of the exporter. It returns a
// function that needs to be deferred so it is called when the method returns.
func (i *Instrumentation) ExportLogs(ctx context.Context, count int64) ExportLogsDone {
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)

	*addOpt = append(*addOpt, i.setOpt)

	start := time.Now()
	i.logInflightMetric.Add(ctx, count, *addOpt...)

	return i.end(ctx, start, count)
}

func (i *Instrumentation) end(ctx context.Context, start time.Time, count int64) ExportLogsDone {
	return func(err error, success int64, code codes.Code) {
		addOpt := get[metric.AddOption](addOpPool)
		defer put(addOpPool, addOpt)

		*addOpt = append(*addOpt, i.setOpt)

		duration := time.Since(start).Seconds()
		i.logInflightMetric.Add(ctx, -count, *addOpt...)
		i.logExportedMetric.Add(ctx, success, *addOpt...)

		mOpt := i.setOpt
		if err != nil {
			attrs := get[attribute.KeyValue](attrsPool)
			defer put(attrsPool, attrs)
			*attrs = append(*attrs, i.presetAttrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			set := attribute.NewSet(*attrs...)
			mOpt = metric.WithAttributeSet(set)

			*addOpt = append((*addOpt)[:0], mOpt)

			i.logExportedMetric.Add(ctx, count-success, *addOpt...)
		}

		recordOpt := get[metric.RecordOption](recordOptPool)
		defer put(recordOptPool, recordOpt)
		*recordOpt = append(
			*recordOpt,
			mOpt,
			metric.WithAttributes(
				semconv.RPCGRPCStatusCodeKey.Int64(int64(code)),
			),
		)
		i.logExportedDurationMetric.Record(ctx, duration, *recordOpt...)
	}
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
