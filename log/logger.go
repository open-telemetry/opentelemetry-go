// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"context"

	"go.opentelemetry.io/otel/log/embedded"
)

// Logger emits log records.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type Logger interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.Logger

	// Emit emits a log record.
	//
	// The record may be held by the implementation. Callers should not mutate
	// the record after passed.
	//
	// Implementations of this method need to be safe for a user to call
	// concurrently.
	Emit(ctx context.Context, record Record)

	// Enabled returns whether the Logger emits for the given context and
	// record.
	//
	// The passed record is likely to be a partial record with only the
	// bridge-relevant information being provided (e.g a record with only the
	// Severity set). If a Logger needs more information than is provided, it
	// is said to be in an indeterminate state (see below).
	//
	// The returned value will be true when the Logger will emit for the
	// provided context and record, and will be false if the Logger will not
	// emit. The returned value may be true or false in an indeterminate state.
	// An implementation should default to returning true for an indeterminate
	// state, but may return false if valid reasons in particular circumstances
	// exist (e.g. performance, correctness).
	//
	// The record should not be held by the implementation. A copy should be
	// made if the record needs to be held after the call returns.
	//
	// Implementations of this method need to be safe for a user to call
	// concurrently.
	Enabled(ctx context.Context, record Record) bool
}
