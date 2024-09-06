// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlplogfile // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"

import (
	"context"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile/internal/transform"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile/internal/writer"
	"go.opentelemetry.io/otel/sdk/log"
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

// Exporter is an OpenTelemetry log exporter that outputs log records
// into files, as JSON. The implementation is based on the specification
// defined here: https://github.com/open-telemetry/opentelemetry-specification/blob/v1.36.0/specification/protocol/file-exporter.md
type Exporter struct {
	mu      sync.Mutex
	w       *writer.Writer
	stopped bool
}

// Compile-time check that the implementation satisfies the interface.
var _ log.Exporter = &Exporter{}

// New returns a new [Exporter].
func New(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options)
	if err != nil {
		return nil, err
	}

	w, err := writer.New(cfg.out, cfg.flushInterval)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		w:       w,
		stopped: false,
	}, nil
}

// Export exports logs records to the file.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	// Honor context cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopped {
		return nil
	}

	data := &lpb.LogsData{
		ResourceLogs: transform.ResourceLogs(records),
	}

	by, err := protojson.Marshal(data)
	if err != nil {
		return err
	}

	return e.w.Export(by)
}

// ForceFlush flushes data to the file.
func (e *Exporter) ForceFlush(_ context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopped {
		return nil
	}

	return e.w.Flush()
}

// Shutdown shuts down the exporter. Buffered data is written to disk,
// and opened resources such as file will be closed.
func (e *Exporter) Shutdown(_ context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopped {
		return nil
	}

	e.stopped = true
	return e.w.Shutdown()
}
