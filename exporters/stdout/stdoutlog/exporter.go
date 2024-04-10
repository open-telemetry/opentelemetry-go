// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/log"
)

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	encoder    *json.Encoder
	timestamps bool

	running atomic.Bool
}

// New creates an [Exporter] with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	e := Exporter{
		encoder:    enc,
		timestamps: cfg.Timestamps,
	}
	e.running.Store(true)

	return &e, nil
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if !e.running.Load() {
		// Free the encoder resources.
		e.encoder = nil
		return nil
	}

	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		// Encode record, one by one.
		recordJSON := e.newRecordJSON(record)
		if err := e.encoder.Encode(recordJSON); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown stops the exporter.
func (e *Exporter) Shutdown(context.Context) error {
	e.running.Store(false)

	return nil
}

// ForceFlush performs no action.
func (e *Exporter) ForceFlush(context.Context) error {
	return nil
}
