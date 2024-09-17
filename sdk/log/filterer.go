// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains support for Logs SDK experimental features.
package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"

	"go.opentelemetry.io/otel/log"
)

// Filterer handles filtering of log records.
//
// Any of the Filterer's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Filterer to manage
// this concurrency.
type Filterer interface {
	// Filter returns whether the SDK will process for the given context
	// and param.
	//
	// The passed param is may be a partial record with only the
	// bridge-relevant information being provided (e.g a record with only the
	// Severity set). If a Logger needs more information than is provided, it
	// is said to be in an indeterminate state (see below).
	//
	// The returned value will be true when the Processor will process for the
	// provided context and param, and will be false if the Processor will not
	// process. An implementation should default to returning true for an
	// indeterminate state.
	Filter(ctx context.Context, param log.EnabledParameters) bool
}
