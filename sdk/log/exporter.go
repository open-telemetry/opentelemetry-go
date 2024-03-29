// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
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

// chunker wraps an Exporter's Export method so it is called with
// appropriately sized export payloads and timeouts. Any payload larger than a
// defined size is chunked into smaller payloads and exported sequentially. The
// entire export (all chunks) needs to complete within the defined timeout,
// otherwise the export is canceled.
type chunker struct {
	Exporter

	// Size is the maximum batch Size exported.
	//
	// If Size is less than or equal to 0 no chunking will be done.
	Size int
	// Timeout is the maximum time an entire export (all batches) is attempted.
	//
	// If Timeout is less than or equal to 0 no timeout will be used.
	Timeout time.Duration
}

func (c chunker) Export(ctx context.Context, records []Record) error {
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	if c.Size <= 0 {
		return c.Exporter.Export(ctx, records)
	}

	n := len(records)
	for i, j := 0, min(c.Size, n); i < n; i, j = i+c.Size, min(j+c.Size, n) {
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
