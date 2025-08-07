// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	encoder           atomic.Pointer[json.Encoder]
	timestamps        bool
	selfObservability *selfObservability
}

type selfObservability struct {
	inflight otelconv.SDKExporterLogInflight
	exported otelconv.SDKExporterLogExported
	duration otelconv.SDKExporterOperationDuration
}

// New creates an [Exporter].
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	var selfObs *selfObservability
	if cfg.SelfObservability {
		selfObs = newSelfObservability()
	}

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	e := Exporter{
		timestamps:        cfg.Timestamps,
		selfObservability: selfObs,
	}
	e.encoder.Store(enc)

	return &e, nil
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if len(records) == 0 {
		return nil
	}
	enc := e.encoder.Load()
	if enc == nil {
		return nil
	}
	e.initSelfObservability(ctx, &records)

	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		// Encode record, one by one.
		recordJSON := e.newRecordJSON(record)
		if err := enc.Encode(recordJSON); err != nil {
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
