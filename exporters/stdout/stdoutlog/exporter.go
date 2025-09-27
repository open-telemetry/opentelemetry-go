// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/counter"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/observ"
	"go.opentelemetry.io/otel/sdk/log"
)

// otelComponentType is a name identifying the type of the OpenTelemetry component.
const (
	otelComponentType = "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

	// Version is the current version of this instrumentation.
	//
	// This matches the version of the exporter.
	Version = internal.Version
)

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	encoder         atomic.Pointer[json.Encoder]
	timestamps      bool
	instrumentation *observ.Instrumentation
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

	exporterID := counter.NextExporterID()
	inst, err := observ.NewInstrumentation(otelComponentType, exporterID)
	if err != nil {
		return nil, err
	}
	e.instrumentation = inst

	return &e, nil
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if inst := e.instrumentation; inst != nil {
		done := inst.ExportLogs(ctx, len(records))
		exported, err := e.exportRecords(ctx, records)
		done(exported, err)
		return err
	}

	_, err := e.exportRecords(ctx, records)
	return err
}

func (e *Exporter) exportRecords(ctx context.Context, records []log.Record) (int64, error) {
	enc := e.encoder.Load()
	if enc == nil {
		return 0, nil
	}

	var exported int64
	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return exported, err
		}

		recordJSON := e.newRecordJSON(record)
		if err := enc.Encode(recordJSON); err != nil {
			return exported, err
		}
		exported++
	}

	return exported, nil
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
