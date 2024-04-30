// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace // import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var zeroTime time.Time

var _ trace.SpanExporter = &Exporter{}

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}
	exp := &Exporter{
		encoder: enc,
	}
	exp.timestamps.Store(cfg.Timestamps)
	return exp, nil
}

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type Exporter struct {
	encoder    *json.Encoder
	timestamps atomic.Bool
	stopped    atomic.Bool
}

// ExportSpans writes spans in json format to stdout.
func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if e.stopped.Load() {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}

	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)

	for i := range stubs {
		stub := &stubs[i]
		// Remove timestamps
		if !(e.timestamps.Load()) {
			stub.StartTime = zeroTime
			stub.EndTime = zeroTime
			for j := range stub.Events {
				ev := &stub.Events[j]
				ev.Time = zeroTime
			}
		}

		// Encode span stubs, one by one
		if err := e.encoder.Encode(stub); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown is called to stop the exporter, it performs no action.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stopped.Store(true)
	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent this Exporter.
func (e *Exporter) MarshalLog() interface{} {
	return struct {
		Type           string
		WithTimestamps bool
	}{
		Type:           "stdout",
		WithTimestamps: e.timestamps.Load(),
	}
}
