// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
)

// Processor handles the processing of log records.
//
// Any of the Processor's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Processor to manage
// this concurrency.
type Processor interface {
	// OnEmit is called when a Record is emitted.
	//
	// OnEmit will be called independent of Enabled. Implementations need to
	// validate the arguments themselves before processing.
	//
	// Implementation should not interrupt the record processing
	// if the context is canceled.
	//
	// All retry logic must be contained in this function. The SDK does not
	// implement any retry logic. All errors returned by this function are
	// considered unrecoverable and will be reported to a configured error
	// Handler.
	//
	// Before modifying a Record, the implementation must use Record.Clone
	// to create a copy that shares no state with the original.
	OnEmit(ctx context.Context, record Record) error
	// Enabled returns whether the Processor will process for the given context
	// and record.
	//
	// The passed record is likely to be a partial record with only the
	// bridge-relevant information being provided (e.g a record with only the
	// Severity set). If a Logger needs more information than is provided, it
	// is said to be in an indeterminate state (see below).
	//
	// The returned value will be true when the Processor will process for the
	// provided context and record, and will be false if the Processor will not
	// process. The returned value may be true or false in an indeterminate
	// state. An implementation should default to returning true for an
	// indeterminate state, but may return false if valid reasons in particular
	// circumstances exist (e.g. performance, correctness).
	//
	// Before modifying a Record, the implementation must use Record.Clone
	// to create a copy that shares no state with the original.
	Enabled(ctx context.Context, record Record) bool
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
