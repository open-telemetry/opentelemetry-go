// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"

import "fmt"

// PartialSuccess represents the underlying error for all handling
// stdoutlog partial success messages.  Use `errors.Is(err, PartialSuccess{})`
// to test whether an error passed to the stdoutlog error handler belongs to this category.
type PartialSuccess struct {
	ErrorMessage  string
	RejectedItems int64
	RejectedKind  string
}

var _ error = PartialSuccess{}

// Error implements the error interface.
func (ps PartialSuccess) Error() string {
	msg := ps.ErrorMessage
	if msg == "" {
		msg = "empty message"
	}
	return fmt.Sprintf("stdoutlog partial success: %s (%d %s failed)", msg, ps.RejectedItems, ps.RejectedKind)
}

// Is supports the errors.Is() interface.
func (PartialSuccess) Is(err error) bool {
	_, ok := err.(PartialSuccess)
	return ok
}

// LogPartialSuccessError returns an error describing a partial success
// response for the log signal.
func LogPartialSuccessError(itemsRejected int64, errorMessage string) error {
	return PartialSuccess{
		ErrorMessage:  errorMessage,
		RejectedItems: itemsRejected,
		RejectedKind:  "logs",
	}
}
