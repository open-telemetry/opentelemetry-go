// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package errorhandler provides the global error handler for OpenTelemetry.
//
// This package has no OTel dependencies, allowing it to be imported by any
// package in the module without creating import cycles.
package errorhandler // import "go.opentelemetry.io/otel/internal/errorhandler"

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

// ErrorHandler handles irremediable events.
type ErrorHandler interface {
	// Handle handles any error deemed irremediable by an OpenTelemetry
	// component.
	Handle(error)
}

// ErrDelegator delegates error handling to a wrapped ErrorHandler.
type ErrDelegator struct {
	Delegate atomic.Pointer[ErrorHandler]
}

// Compile-time check that ErrDelegator implements ErrorHandler.
var _ ErrorHandler = (*ErrDelegator)(nil)

func (d *ErrDelegator) Handle(err error) {
	if eh := d.Delegate.Load(); eh != nil {
		(*eh).Handle(err)
		return
	}
	log.Print(err)
}

// SetDelegate sets the ErrorHandler delegate.
func (d *ErrDelegator) SetDelegate(eh ErrorHandler) {
	d.Delegate.Store(&eh)
}

type errorHandlerHolder struct {
	eh ErrorHandler
}

var (
	globalErrorHandler       = defaultErrorHandler()
	delegateErrorHandlerOnce sync.Once
)

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
	return globalErrorHandler.Load().(errorHandlerHolder).eh
}

// SetErrorHandler sets the global ErrorHandler to h.
//
// The first time this is called all ErrorHandler previously returned from
// GetErrorHandler will send errors to h instead of the default logging
// ErrorHandler. Subsequent calls will set the global ErrorHandler, but not
// delegate errors to h.
func SetErrorHandler(h ErrorHandler) {
	current := GetErrorHandler()

	if _, cOk := current.(*ErrDelegator); cOk {
		if _, ehOk := h.(*ErrDelegator); ehOk && current == h {
			// Do not assign to the delegate of the default ErrDelegator to be
			// itself.
			log.Print(errors.New("no ErrorHandler delegate configured"), " ErrorHandler remains its current value.")
			return
		}
	}

	delegateErrorHandlerOnce.Do(func() {
		if def, ok := current.(*ErrDelegator); ok {
			def.SetDelegate(h)
		}
	})
	globalErrorHandler.Store(errorHandlerHolder{eh: h})
}

func defaultErrorHandler() *atomic.Value {
	v := &atomic.Value{}
	v.Store(errorHandlerHolder{eh: &ErrDelegator{}})
	return v
}
