// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	ScopeName = "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal/observ"
	Version   = internal.Version
)

type Instrumentation struct {
	inflight  metric.Int64UpDownCounter
	exported  metric.Int64Counter
	operation metric.Float64Histogram

	presetAttrs []attribute.KeyValue
	addOpt      metric.AddOption
	recordOpt   metric.RecordOption
}

func NewInstrumentation(name string, id int64, target string) (*Instrumentation, error) {
	inst := &Instrumentation{}

	provider := otel.GetMeterProvider()
	m := provider.Meter(
		ScopeName,
		metric.WithSchemaURL(semconv.SchemaURL),
		metric.WithInstrumentationVersion(Version),
	)

	var e, err error
	logInflight, e := otelconv.NewSDKExporterLogInflight(m)
	if e != nil {
		e = fmt.Errorf("failed to create the inflight metric %w", e)
		err = errors.Join(err, e)
	}
	inst.inflight = logInflight.Inst()

	exported, e := otelconv.NewSDKExporterLogExported(m)
	if e != nil {
		e = fmt.Errorf("failed to create the exported metric %w", e)
		err = errors.Join(err, e)
	}
	inst.exported = exported.Inst()

	operation, e := otelconv.NewSDKExporterOperationDuration(m)
	if e != nil {
		e = fmt.Errorf("failed to create the operation metric %w", e)
		err = errors.Join(err, e)
	}
	inst.operation = operation.Inst()

	if err != nil {
		return nil, err
	}
	return inst, err
}
