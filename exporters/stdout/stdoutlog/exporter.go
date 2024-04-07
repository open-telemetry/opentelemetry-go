// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/sdk/log"
)

var zeroTime time.Time

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// The writer is os.Stdout by default.
type Exporter struct {
	encoder    *json.Encoder
	timestamps bool

	stopped atomic.Bool
}

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	return &Exporter{
		encoder:    enc,
		timestamps: cfg.Timestamps,
	}, nil
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if e.stopped.Load() {
		return nil
	}

	// Prevent panic if encoder is nil.
	if e.encoder == nil {
		e.encoder = json.NewEncoder(defaultWriter)
	}

	for i, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		// Remove timestamps.
		if !e.timestamps {
			// Clone before make changes.
			record = records[i].Clone()

			record.SetTimestamp(zeroTime)
			record.SetObservedTimestamp(zeroTime)
		}

		// Encode record, one by one.
		recordJSON := newRecordJSON(record)
		if err := e.encoder.Encode(recordJSON); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown stops the exporter.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stopped.Store(true)
	// Free the encoder resources.
	e.encoder = nil

	return nil
}

// ForceFlush performs no action.
func (e *Exporter) ForceFlush(ctx context.Context) error {
	return nil
}
