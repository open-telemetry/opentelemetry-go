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
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

var attrsPool = sync.Pool{
	New: func() any {
		// "component.name" + "component.type" + "error.type" + "server.address" + "server.port" + "rpc.grpc.status_code"
		const n = 1 + 1 + 1 + 1 + 1 + 1
		s := make([]attribute.KeyValue, 0, n)
		return &s
	},
}

type ExporterMetrics struct {
	logInflightMetric         otelconv.SDKExporterLogInflight
	logExportedMetric         otelconv.SDKExporterLogExported
	logExportedDurationMetric otelconv.SDKExporterOperationDuration
	presetAttrs               []attribute.KeyValue
}

func NewExporterMetrics(
	name, componentName string,
	componentType otelconv.ComponentTypeAttr,
	target string,
) (*ExporterMetrics, error) {
	em := &ExporterMetrics{}
	em.presetAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(string(componentType)),
	}
	em.presetAttrs = append(em.presetAttrs, ServerAddrAttrs(target)...)

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
		otel.Handle(err)
		err = errors.Join(err, e)
	}
	return em, err
}

func (em *ExporterMetrics) TrackExport(
	ctx context.Context,
	count int64,
) func(err error, successCount int64, code codes.Code) {
	attrs := attrsPool.Get().(*[]attribute.KeyValue)
	*attrs = append([]attribute.KeyValue{}, em.presetAttrs...)

	begin := time.Now()
	em.logInflightMetric.Add(ctx, count, *attrs...)
	return func(err error, successCount int64, code codes.Code) {
		defer func() {
			*attrs = (*attrs)[:0]
			attrsPool.Put(attrs)
		}()

		duration := time.Since(begin).Seconds()
		em.logInflightMetric.Add(ctx, -count, *attrs...)
		em.logExportedMetric.Add(ctx, successCount, *attrs...)
		if err != nil {
			*attrs = append(*attrs, semconv.ErrorType(err))
			em.logExportedMetric.Add(ctx, count-successCount, *attrs...)
		}
		*attrs = append(
			*attrs,
			em.logExportedDurationMetric.AttrRPCGRPCStatusCode(otelconv.RPCGRPCStatusCodeAttr(code)),
		)
		em.logExportedDurationMetric.Record(ctx, duration, *attrs...)
	}
}

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
