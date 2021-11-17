// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otel // import "go.opentelemetry.io/otel"

import (
	"sync"

	"go.opentelemetry.io/otel/internal/debug"
)

var (
	// globalErrorHandler provides an ErrorHandler that can be used
	// throughout an OpenTelemetry instrumented project. When a user
	// specified ErrorHandler is registered (`SetErrorHandler`) all calls to
	// `Handle` and will be delegated to the registered ErrorHandler.
	globalErrorHandler = &errorHandlerDelegate{
		delegate: &defaultErrorHandler{},
	}
	// delegateErrorHandlerOnce ensures that a user provided ErrorHandler is
	// only ever registered once.
	delegateErrorHandlerOnce sync.Once

	// Compile-time check that delegator implements ErrorHandler.
	_ ErrorHandler = (*errorHandlerDelegate)(nil)
	_ ErrorHandler = (*defaultErrorHandler)(nil)
)

// errorHandlerDelegate is a box type to enable updating of all handlers returned by GetErrorHandler()
type errorHandlerDelegate struct {
	delegate ErrorHandler
}

func (h *errorHandlerDelegate) setDelegate(d ErrorHandler) {
	h.delegate = d
}

// Handle handles any error deemed irremediable by an OpenTelemetry component.
func (h *errorHandlerDelegate) Handle(err error) {
	if h == nil {
		return
	}
	h.delegate.Handle(err)
}

// defaultErrorHandler utilizes the internal logger to manage the messages.
type defaultErrorHandler struct{}

func (h *defaultErrorHandler) Handle(err error) {
	debug.Error(err, "")
}

// DiscardErrorHandler drops all errors.
// Use `SetErrorHandler(DiscardErrorHandler{})` to disable error handling
type DiscardErrorHandler struct{}

func (DiscardErrorHandler) Handle(error) {}

// GetErrorHandler returns the global ErrorHandler instance.
//
// The default ErrorHandler instance returned will log all errors to STDERR
// until an override ErrorHandler is set with SetErrorHandler. All
// ErrorHandler returned prior to this will automatically forward errors to
// the set instance instead of logging.
func GetErrorHandler() ErrorHandler {
	return globalErrorHandler
}

// SetErrorHandler sets the global ErrorHandler to h.
func SetErrorHandler(h ErrorHandler) {
	globalErrorHandler.setDelegate(h)
}

// Handle is a convenience function for ErrorHandler().Handle(err)
func Handle(err error) {
	GetErrorHandler().Handle(err)
}
