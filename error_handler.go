// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel // import "go.opentelemetry.io/otel"

import "go.opentelemetry.io/otel/internal/global"

// ErrorHandler handles irremediable events.
// It's an alias for [global.ErrorHandler].
type ErrorHandler = global.ErrorHandler

// ErrorHandlerFunc is a convenience adapter to allow the use of a function
// as an ErrorHandler.
type ErrorHandlerFunc func(error)

var _ ErrorHandler = ErrorHandlerFunc(nil)

// Handle handles the irremediable error by calling the ErrorHandlerFunc itself.
func (f ErrorHandlerFunc) Handle(err error) {
	f(err)
}
