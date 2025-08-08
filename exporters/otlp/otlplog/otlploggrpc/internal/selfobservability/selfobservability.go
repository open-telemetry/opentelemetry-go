// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package selfobservability provides self-observability metrics for OTLP log exporters.
// This is an experimental feature controlled by the x.SelfObservability feature flag.
package selfobservability // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/selfobservability"

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

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
) *ExporterMetrics {
	em := &ExporterMetrics{}
	mp := otel.GetMeterProvider()
	m := mp.Meter(
		name,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	if em.logInflightMetric, err = otelconv.NewSDKExporterLogInflight(m); err != nil {
		otel.Handle(err)
	}
	if em.logExportedMetric, err = otelconv.NewSDKExporterLogExported(m); err != nil {
		otel.Handle(err)
	}
	if em.logExportedDurationMetric, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}

	em.presetAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(componentName),
		semconv.OTelComponentTypeKey.String(string(componentType)),
	}
	em.presetAttrs = append(em.presetAttrs, ServerAddrAttrs(target)...)
	return em
}

func (em *ExporterMetrics) TrackExport(ctx context.Context, count int64) func(err error, code codes.Code) {
	begin := time.Now()
	em.logInflightMetric.Add(ctx, count, em.presetAttrs...)
	return func(err error, code codes.Code) {
		duration := time.Since(begin).Seconds()
		em.logInflightMetric.Add(ctx, -count, em.presetAttrs...)
		if err != nil {
			em.presetAttrs = append(em.presetAttrs, semconv.ErrorType(err))
		}
		em.logExportedMetric.Add(ctx, count, em.presetAttrs...)
		em.presetAttrs = append(
			em.presetAttrs,
			em.logExportedDurationMetric.AttrRPCGRPCStatusCode(otelconv.RPCGRPCStatusCodeAttr(code)),
		)
		em.logExportedDurationMetric.Record(ctx, duration, em.presetAttrs...)
	}
}

func ServerAddrAttrs(target string) []attribute.KeyValue {
	if strings.HasPrefix(target, "unix://") {
		path := strings.TrimPrefix(target, "unix://")
		return []attribute.KeyValue{semconv.ServerAddress(path)}
	}

	if idx := strings.Index(target, "://"); idx != -1 {
		target = target[idx+4:]
	}

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
