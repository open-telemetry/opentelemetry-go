// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
	// The passed param may be a partial record (e.g a record with only the
	// Severity set). If a Filterer needs more information than is provided, it
	// is said to be in an indeterminate state. An implementation should
	// return true for an indeterminate state.
	//
	// The returned value will be true when the SDK should process for the
	// provided context and param, and will be false if the SDK should not
	// process.
	Filter(ctx context.Context, param FilterParameters) bool
}

// FilterParameters represent Filter parameters.
type FilterParameters struct {
	severity    log.Severity
	severitySet bool

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// Severity returns the [Severity] level value, or [SeverityUndefined] if no value was set.
// The ok result indicates whether the value was set.
func (r *FilterParameters) Severity() (value log.Severity, ok bool) {
	return r.severity, r.severitySet
}

// setSeverity sets the [Severity] level.
func (r *FilterParameters) setSeverity(level log.Severity) {
	r.severity = level
	r.severitySet = true
}
