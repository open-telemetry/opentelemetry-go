// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/internal/global"

import (
	"log"
	"os"
	"sync/atomic"
)

var (
	// GlobalErrorHandler provides an ErrorHandler that can be used
	// throughout an OpenTelemetry instrumented project. When a user
	// specified ErrorHandler is registered (`SetErrorHandler`) all calls to
	// `Handle` and will be delegated to the registered ErrorHandler.
	GlobalErrorHandler = defaultErrorHandler()

	// Compile-time check that delegator implements ErrorHandler.
	_ ErrorHandler = (*ErrDelegator)(nil)
	// Compile-time check that errLogger implements ErrorHandler.
	_ ErrorHandler = (*ErrLogger)(nil)
)

// ErrorHandler handles irremediable events.
type ErrorHandler interface {
	// Handle handles any error deemed irremediable by an OpenTelemetry
	// component.
	Handle(error)
}

type ErrDelegator struct {
	delegate atomic.Pointer[ErrorHandler]
}

func (d *ErrDelegator) Handle(err error) {
	d.getDelegate().Handle(err)
}

func (d *ErrDelegator) getDelegate() ErrorHandler {
	return *d.delegate.Load()
}

// setDelegate sets the ErrorHandler delegate.
func (d *ErrDelegator) setDelegate(eh ErrorHandler) {
	d.delegate.Store(&eh)
}

func defaultErrorHandler() *ErrDelegator {
	d := &ErrDelegator{}
	d.setDelegate(&ErrLogger{l: log.New(os.Stderr, "", log.LstdFlags)})
	return d
}

// ErrLogger logs errors if no delegate is set, otherwise they are delegated.
type ErrLogger struct {
	l *log.Logger
}

// Handle logs err if no delegate is set, otherwise it is delegated.
func (h *ErrLogger) Handle(err error) {
	h.l.Print(err)
}

// GetErrorHandler returns the global ErrorHandler instance.
//
// The default ErrorHandler instance returned will log all errors to STDERR
// until an override ErrorHandler is set with SetErrorHandler. All
// ErrorHandler returned prior to this will automatically forward errors to
// the set instance instead of logging.
//
// Subsequent calls to SetErrorHandler after the first will not forward errors
// to the new ErrorHandler for prior returned instances.
func GetErrorHandler() ErrorHandler {
	return GlobalErrorHandler
}

// SetErrorHandler sets the global ErrorHandler to h.
//
// The first time this is called all ErrorHandler previously returned from
// GetErrorHandler will send errors to h instead of the default logging
// ErrorHandler. Subsequent calls will set the global ErrorHandler, but not
// delegate errors to h.
func SetErrorHandler(h ErrorHandler) {
	GlobalErrorHandler.setDelegate(h)
}

// Handle is a convenience function for ErrorHandler().Handle(err).
func Handle(err error) {
	GetErrorHandler().Handle(err)
}
