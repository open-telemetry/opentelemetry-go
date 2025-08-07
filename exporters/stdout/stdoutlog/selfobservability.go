// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

func newSelfObservability() *selfObservability {
	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/exporters/stdout/stdoutlog",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	so := selfObservability{}

	var err error
	if so.inflight, err = otelconv.NewSDKExporterLogInflight(m); err != nil {
		otel.Handle(err)
	}
	if so.exported, err = otelconv.NewSDKExporterLogExported(m); err != nil {
		otel.Handle(err)
	}
	if so.duration, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}
	return &so
}

func (e *Exporter) initSelfObservability(ctx context.Context, records *[]log.Record) {
	if records == nil || e.selfObservability == nil {
		return
	}

	e.selfObservability.inflight.Add(ctx, int64(len(*records)))

	start := time.Now()
	defer func() {
		dur := float64(time.Since(start).Nanoseconds())
		e.selfObservability.duration.Record(ctx, dur)
	}()

	e.selfObservability.exported.Add(ctx, int64(len(*records)))
}
