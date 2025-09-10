// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for OTLP log exporters.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/selfobservability"

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

// Instrumentation is experimental instrumentation for the exporter.
type Instrumentation struct {
	logInflightMetric         otelconv.SDKExporterLogInflight
	logExportedMetric         otelconv.SDKExporterLogExported
	logExportedDurationMetric otelconv.SDKExporterOperationDuration
	presetAttrs               []attribute.KeyValue
	setOpt                    metric.MeasurementOption
}

// NewInstrumentation returns instrumentation for otlplog grpc exporter.
func NewInstrumentation(
	name, componentName string,
	componentType otelconv.ComponentTypeAttr,
	target string,
) (*Instrumentation, error) {
	em := &Instrumentation{}

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		name,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err, e error
	if em.logInflightMetric, e = otelconv.NewSDKExporterLogInflight(m); e != nil {
		e = fmt.Errorf("failed to create span inflight metric: %w", e)
		otel.Handle(e)
		err = errors.Join(err, e)
	}
	if em.logExportedMetric, e = otelconv.NewSDKExporterLogExported(m); e != nil {
		e = fmt.Errorf("failed to create span exported metric: %w", e)
		otel.Handle(e)
		err = errors.Join(err, e)
	}
	if em.logExportedDurationMetric, e = otelconv.NewSDKExporterOperationDuration(m); e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		otel.Handle(e)
		err = errors.Join(err, e)
	}

	if err != nil {
		return nil, err
	}

	em.presetAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(string(componentType)),
	}
	em.presetAttrs = append(em.presetAttrs, ServerAddrAttrs(target)...)
	s := attribute.NewSet(em.presetAttrs...)
	em.setOpt = metric.WithAttributeSet(s)

	return em, nil
}

// ExportSpanDone is a function that is called when a call to an Exporter's
// ExportSpans method completes
//
// The number of successful exports is provided as success. Any error that is encountered is provided as error
// The code of last gRPC requests performed in scope of this export call.
type ExportSpanDone func(err error, success int64, code codes.Code)

// ExportSpans instruments the ExportSpans method of the exporter. It returns a
// function that needs to be deferred so it is called when the method returns.
func (i *Instrumentation) ExportSpans(ctx context.Context, count int64) ExportSpanDone {
	addOpt := get[metric.AddOption](addOpPool)
	defer put(addOpPool, addOpt)

	*addOpt = append(*addOpt, i.setOpt)

	start := time.Now()
	i.logInflightMetric.Inst().Add(ctx, count, *addOpt...)

	return i.end(ctx, start, count)
}

func (i *Instrumentation) end(ctx context.Context, start time.Time, count int64) ExportSpanDone {
	return func(err error, success int64, code codes.Code) {
		addOpt := get[metric.AddOption](addOpPool)
		defer put(addOpPool, addOpt)

		*addOpt = append(*addOpt, i.setOpt)

		duration := time.Since(start).Seconds()
		i.logInflightMetric.Inst().Add(ctx, -count, *addOpt...)
		i.logExportedMetric.Inst().Add(ctx, success, *addOpt...)

		mOpt := i.setOpt
		if err != nil {
			attrs := get[attribute.KeyValue](attrsPool)
			defer put(attrsPool, attrs)
			*attrs = append(*attrs, i.presetAttrs...)
			*attrs = append(*attrs, semconv.ErrorType(err))

			set := attribute.NewSet(*attrs...)
			mOpt = metric.WithAttributeSet(set)

			*addOpt = append((*addOpt)[:0], mOpt)

			i.logExportedMetric.Inst().Add(ctx, count-success, *addOpt...)
		}

		recordOpt := get[metric.RecordOption](recordOptPool)
		defer put(recordOptPool, recordOpt)
		*recordOpt = append(
			*recordOpt,
			mOpt,
			metric.WithAttributes(
				i.logExportedDurationMetric.AttrRPCGRPCStatusCode(otelconv.RPCGRPCStatusCodeAttr(code)),
			),
		)
		i.logExportedDurationMetric.Inst().Record(ctx, duration, *recordOpt...)
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
