// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/log"
)

var zeroTime time.Time

var _ log.Exporter = &Exporter{}

// Exporter is an implementation of  that writes spans to stdout.
type Exporter struct {
	encoder    *json.Encoder
	timestamps bool

	stoppedMu sync.RWMutex
	stopped   bool
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
// The writer is os.Stdout by default.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(records) == 0 {
		return nil
	}

	for i := range records {
		// Honor context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		record := records[i]
		// Remove timestamps
		if !e.timestamps {
			// Clone before make changes
			record = records[i].Clone()

			record.SetTimestamp(zeroTime)
			record.SetObservedTimestamp(zeroTime)
		}

		// Encode record, one by one
		if err := e.encoder.Encode(&record); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown stops the exporter.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// ForceFlush performs no action.
func (e *Exporter) ForceFlush(ctx context.Context) error {
	return nil
}
