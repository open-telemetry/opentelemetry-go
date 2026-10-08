// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync/atomic"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/counter"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/observ"

	"go.opentelemetry.io/otel/sdk/log"
)

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	writer     atomic.Pointer[exportWriter]
	stopped    atomic.Bool
	timestamps bool
	inst       *observ.Instrumentation
}

// New creates an [Exporter].
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	w := &exportWriter{writer: cfg.Writer, available: make(chan struct{}, 1)}
	w.available <- struct{}{}
	w.encoder = json.NewEncoder(&w.buffer)
	if cfg.PrettyPrint {
		w.encoder.SetIndent("", "\t")
	}

	e := &Exporter{
		timestamps: cfg.Timestamps,
	}
	e.writer.Store(w)

	var err error
	e.inst, err = observ.NewInstrumentation(counter.NextExporterID())
	return e, err
}

// Export exports log records to the writer. It returns [log.ErrExporterShutdown]
// if called after Shutdown.
//
// Cancellation stops waiting for the writer, but cannot interrupt an in-progress
// [io.Writer.Write]. At most one write can remain in progress, and subsequent
// exports wait for it with their own contexts. Only encoded JSON is retained
// after Export returns; the supplied records are not retained.
func (e *Exporter) Export(ctx context.Context, records []log.Record) (err error) {
	w := e.writer.Load()
	if e.stopped.Load() {
		return log.ErrExporterShutdown
	}

	if w == nil {
		return nil
	}

	var success int64
	if e.inst != nil {
		op := e.inst.ExportLogs(ctx, int64(len(records)))
		defer func() {
			op.End(success, err)
		}()
	}

	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := w.acquire(ctx); err != nil {
			return err
		}
		if e.stopped.Load() {
			w.release()
			return log.ErrExporterShutdown
		}

		// Encode before starting the write so it cannot outlive record ownership.
		w.buffer.Reset()
		if err := w.encoder.Encode(e.newRecordJSON(record)); err != nil {
			w.release()
			return err
		}
		if err := ctx.Err(); err != nil {
			w.release()
			return err
		}
		if err := w.write(ctx); err != nil {
			return err
		}
		success++
	}
	return nil
}

// Shutdown shuts down the Exporter. Calls to Export after Shutdown return
// [log.ErrExporterShutdown].
// The first call waits for an in-progress write to complete or ctx to end,
// without closing the writer. A write may complete after Shutdown returns if
// ctx ends first. Subsequent calls to Shutdown or ForceFlush do nothing.
func (e *Exporter) Shutdown(ctx context.Context) error {
	// Store stopped first so Export cannot observe a cleared writer while the
	// Exporter still appears active.
	e.stopped.Store(true)
	if w := e.writer.Swap(nil); w != nil {
		return w.wait(ctx)
	}
	return nil
}

// ForceFlush waits for an in-progress write to complete or ctx to end.
// It does not call Flush on the writer and does nothing after Shutdown.
func (e *Exporter) ForceFlush(ctx context.Context) error {
	if w := e.writer.Load(); w != nil && !e.stopped.Load() {
		return w.wait(ctx)
	}
	return nil
}

// exportWriter keeps its buffer unavailable until Write actually returns, even
// if Export has already returned. This bounds background work to one write and
// one encoded record, and prevents a subsequent export from reusing its bytes.
type exportWriter struct {
	writer    io.Writer
	buffer    bytes.Buffer
	encoder   *json.Encoder
	available chan struct{}
}

func (w *exportWriter) acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.available:
		if err := ctx.Err(); err != nil {
			w.release()
			return err
		}
		return nil
	}
}

func (w *exportWriter) release() {
	w.available <- struct{}{}
}

func (w *exportWriter) wait(ctx context.Context) error {
	if err := w.acquire(ctx); err != nil {
		return err
	}
	w.release()
	return nil
}

// write transfers the acquired buffer to the writing goroutine until it returns.
func (w *exportWriter) write(ctx context.Context) error {
	result := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			// A writer callback now runs outside the caller's recovery boundary.
			// Contain its panic and always make the buffer available again.
			if p := recover(); p != nil {
				err = fmt.Errorf("stdoutlog: writer panicked: %v", p)
			}
			result <- err
			w.release()
		}()
		_, err = w.writer.Write(w.buffer.Bytes())
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-result:
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return err
	}
}
