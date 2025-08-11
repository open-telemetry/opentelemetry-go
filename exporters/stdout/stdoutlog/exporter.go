// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// otelComponentType is a name identifying the type of the OpenTelemetry component.
const otelComponentType = "stdout_log_exporter"

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	encoder           atomic.Pointer[json.Encoder]
	timestamps        bool
	selfObservability *selfObservability
}

type selfObservability struct {
	enabled                 bool
	attrs                   []attribute.KeyValue
	inflightMetric          otelconv.SDKExporterLogInflight
	exportedMetric          otelconv.SDKExporterLogExported
	operationDurationMetric otelconv.SDKExporterOperationDuration
}

// New creates an [Exporter].
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	e := Exporter{
		timestamps: cfg.Timestamps,
	}
	e.encoder.Store(enc)
	e.initSelfObservability()

	return &e, nil
}

// initSelfObservability initializes self-observability for the exporter if enabled.
func (e *Exporter) initSelfObservability() {
	if !x.SelfObservability.Enabled() {
		return
	}

	e.selfObservability = &selfObservability{
		enabled: true,
		attrs: []attribute.KeyValue{
			semconv.OTelComponentName(fmt.Sprintf("%s/%d", otelComponentType, nextExporterID())),
			semconv.OTelComponentTypeKey.String(otelComponentType),
		},
	}

	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/exporters/stdout/stdoutlog",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error
	if e.selfObservability.inflightMetric, err = otelconv.NewSDKExporterLogInflight(m); err != nil {
		otel.Handle(err)
	}
	if e.selfObservability.exportedMetric, err = otelconv.NewSDKExporterLogExported(m); err != nil {
		otel.Handle(err)
	}
	if e.selfObservability.operationDurationMetric, err = otelconv.NewSDKExporterOperationDuration(m); err != nil {
		otel.Handle(err)
	}
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	enc := e.encoder.Load()
	if enc == nil {
		return nil
	}

	var err error
	if e.selfObservability != nil && e.selfObservability.enabled {
		err = e.exportWithSelfObservability(ctx, records)
	} else {
		err = e.exportWithoutSelfObservability(ctx, records)
	}
	return err
}

// exportWithSelfObservability exports logs with self-observability metrics.
func (e *Exporter) exportWithSelfObservability(ctx context.Context, records []log.Record) error {
	if len(records) == 0 {
		return nil
	}

	count := int64(len(records))
	start := time.Now()

	e.selfObservability.inflightMetric.Add(context.Background(), count, e.selfObservability.attrs...)
	defer func() {
		addAttrs := make([]attribute.KeyValue, len(e.selfObservability.attrs), len(e.selfObservability.attrs)+1)
		copy(addAttrs, e.selfObservability.attrs)
		if err := ctx.Err(); err != nil {
			addAttrs = append(addAttrs, semconv.ErrorType(err))
		}

		e.selfObservability.inflightMetric.Add(context.Background(), -count, e.selfObservability.attrs...)
		e.selfObservability.exportedMetric.Add(context.Background(), count, addAttrs...)
		e.selfObservability.operationDurationMetric.Record(context.Background(), time.Since(start).Seconds(), addAttrs...)
	}()

	return e.exportWithoutSelfObservability(ctx, records)
}

// exportWithoutSelfObservability exports logs without self-observability metrics.
func (e *Exporter) exportWithoutSelfObservability(ctx context.Context, records []log.Record) error {
	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		recordJSON := e.newRecordJSON(record)
		if err := e.encoder.Load().Encode(recordJSON); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shuts down the Exporter.
// Calls to Export will perform no operation after this is called.
func (e *Exporter) Shutdown(context.Context) error {
	e.encoder.Store(nil)
	return nil
}

// ForceFlush performs no action.
func (*Exporter) ForceFlush(context.Context) error {
	return nil
}

var exporterIDCounter atomic.Int64

// nextExporterID returns a new unique ID for an exporter.
// the starting value is 0, and it increments by 1 for each call.
func nextExporterID() int64 {
	return exporterIDCounter.Add(1) - 1
}
