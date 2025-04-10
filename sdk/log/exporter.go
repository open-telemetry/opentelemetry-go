// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"time"
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
	//
	// Export should never be called concurrently with other Export calls.
	// However, it may be called concurrently with other methods.
	Export(ctx context.Context, records []Record) error

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the exporter should be done in this call.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	//
	// Shutdown may be called concurrently with itself or with other methods.
	Shutdown(ctx context.Context) error

	// ForceFlush exports log records to the configured Exporter that have not yet
	// been exported.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// ForceFlush may be called concurrently with itself or with other methods.
	ForceFlush(ctx context.Context) error
}

var defaultNoopExporter = &noopExporter{}

type noopExporter struct{}

func (noopExporter) Export(context.Context, []Record) error { return nil }

func (noopExporter) Shutdown(context.Context) error { return nil }

func (noopExporter) ForceFlush(context.Context) error { return nil }

// timeoutExporter wraps an Exporter and ensures any call to Export will have a
// timeout for the context.
type timeoutExporter struct {
	Exporter

	// timeout is the maximum time an export is attempted.
	timeout time.Duration
}

// newTimeoutExporter wraps exporter with an Exporter that limits the context
// lifetime passed to Export to be timeout. If timeout is less than or equal to
// zero, exporter will be returned directly.
func newTimeoutExporter(exp Exporter, timeout time.Duration) Exporter {
	if timeout <= 0 {
		return exp
	}
	return &timeoutExporter{Exporter: exp, timeout: timeout}
}

// Export sets the timeout of ctx before calling the Exporter e wraps.
func (e *timeoutExporter) Export(ctx context.Context, records []Record) error {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()
	return e.Exporter.Export(ctx, records)
}
