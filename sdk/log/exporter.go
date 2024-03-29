// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
)

// Exporter handles the delivery of log records to external receivers.
type Exporter interface {
	// Export transmits log records to a receiver.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// All retry logic must be contained in this function. The SDK does not
	// implement any retry logic. All errors returned by this function are
	// considered unrecoverable and will be reported to a configured error
	// Handler.
	//
	// Implementations must not retain the records slice.
	//
	// Before modifying a Record, the implementation must use Record.Clone
	// to create a copy that shares no state with the original.
	Export(ctx context.Context, records []Record) error
	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the exporter should be done in this call.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	Shutdown(ctx context.Context) error
	// ForceFlush exports log records to the configured Exporter that have not yet
	// been exported.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	ForceFlush(ctx context.Context) error
}

var defaultNoopExporter = &noopExporter{}

type noopExporter struct{}

func (noopExporter) Export(context.Context, []Record) error { return nil }

func (noopExporter) Shutdown(context.Context) error { return nil }

func (noopExporter) ForceFlush(context.Context) error { return nil }

type timeoutExporter struct {
	Exporter

	// timeout is the maximum time an entire export (all batches) is attempted.
	//
	// If Timeout is less than or equal to 0 no timeout will be used.
	timeout time.Duration
}

func newTimeoutExporter(exp Exporter, timeout time.Duration) *timeoutExporter {
	return &timeoutExporter{Exporter: exp, timeout: timeout}
}

func (e *timeoutExporter) Export(ctx context.Context, records []Record) error {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()
	return e.Exporter.Export(ctx, records)
}

// chunker wraps an Exporter's Export method so it is called with
// appropriately sized export payloads and timeouts. Any payload larger than a
// defined size is chunked into smaller payloads and exported sequentially. The
// entire export (all chunks) needs to complete within the defined timeout,
// otherwise the export is canceled.
type chunker struct {
	Exporter

	// size is the maximum batch size exported.
	//
	// If size is less than or equal to 0 no chunking will be done.
	size int
}

func newChunkExporter(exp Exporter, size int) *chunker {
	return &chunker{Exporter: exp, size: size}
}

func (c chunker) Export(ctx context.Context, records []Record) error {
	if c.size <= 0 {
		return c.Exporter.Export(ctx, records)
	}

	n := len(records)
	for i, j := 0, min(c.size, n); i < n; i, j = i+c.size, min(j+c.size, n) {
		if err := c.Exporter.Export(ctx, records[i:j]); err != nil {
			return err
		}
	}
	return nil
}

// exportSync exports all data from input using exporter in a spawned
// goroutine. The returned chan will be closed when the spawned goroutine
// completes.
func exportSync(input <-chan exportData, exporter Exporter) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		defer close(done)
		for data := range input {
			data.DoExport(exporter.Export)
		}
	}()
	return done
}

// exportData is data related to an export.
type exportData struct {
	ctx     context.Context
	records []Record

	// respCh is the channel any error returned from the export will be sent
	// on. If this is nil, and the export error is non-nil, the error will
	// passed to the OTel error handler.
	respCh chan<- error
}

// DoExport calls exportFn with the data contained in e. The error response
// will be returned on e's respCh if not nil. The error will be handled by the
// default OTel error handle if it is not nil and respCh is nil or full.
func (e exportData) DoExport(exportFn func(context.Context, []Record) error) {
	if len(e.records) == 0 {
		e.respond(nil)
		return
	}

	e.respond(exportFn(e.ctx, e.records))
}

func (e exportData) respond(err error) {
	select {
	case e.respCh <- err:
	default:
		// e.respCh is nil or busy, default to otel.Handler.
		if err != nil {
			otel.Handle(err)
		}
	}
}

type bufferedExporter struct {
	Exporter

	input   chan exportData
	inputWG sync.WaitGroup

	done    chan struct{}
	stopped atomic.Bool
}

func newBufferedExporter(exporter Exporter, size int) *bufferedExporter {
	input := make(chan exportData, size)
	return &bufferedExporter{
		Exporter: exporter,

		input: input,
		done:  exportSync(input, exporter),
	}
}

var errStopped = errors.New("exporter stopped")

func (e *bufferedExporter) enqueue(ctx context.Context, records []Record, rCh chan<- error) error {
	data := exportData{ctx, records, rCh}

	e.inputWG.Add(1)
	defer e.inputWG.Done()

	// Check stopped before enqueueing now that e.inputWG is incremented to
	// prevent sends on a closed chan when Shutdown is called concurrently.
	if e.stopped.Load() {
		return errStopped
	}

	select {
	case e.input <- data:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (e *bufferedExporter) EnqueueExport(ctx context.Context, records []Record) bool {
	if len(records) == 0 {
		// Nothing to enqueue, do not waste input space.
		return true
	}
	return e.enqueue(ctx, records, nil) == nil
}

func (e *bufferedExporter) Export(ctx context.Context, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	resp := make(chan error, 1)
	err := e.enqueue(ctx, records, resp)
	if err != nil {
		if errors.Is(err, errStopped) {
			return nil
		}
		return fmt.Errorf("%w: dropping %d records", err, len(records))
	}

	select {
	case err := <-resp:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *bufferedExporter) ForceFlush(ctx context.Context) error {
	resp := make(chan error, 1)
	err := e.enqueue(ctx, nil, resp)
	if err != nil {
		if errors.Is(err, errStopped) {
			return nil
		}
		return err
	}

	select {
	case <-resp:
	case <-ctx.Done():
		return ctx.Err()
	}
	return e.Exporter.ForceFlush(ctx)
}

func (e *bufferedExporter) Shutdown(ctx context.Context) error {
	if e.stopped.Swap(true) {
		return nil
	}
	e.inputWG.Wait()

	// No more sends will be made.
	close(e.input)
	select {
	case <-e.done:
	case <-ctx.Done():
		return errors.Join(ctx.Err(), e.Shutdown(ctx))
	}
	return e.Exporter.Shutdown(ctx)
}
