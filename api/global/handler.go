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

package global

import (
	"log"
	"os"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
)

var (
	// globalHandler provides a Handler that can be used throughout
	// an OpenTelemetry instrumented project. When a user specified Handler
	// is registered (`SetHandler`) all calls to `Handle` will be delegated
	// to the registered Handler.
	globalHandler = &handler{
		l: log.New(os.Stderr, "", log.LstdFlags),
	}

	// delegateHanderOnce ensures that a user provided Handler is only ever
	// registered once.
	delegateHanderOnce sync.Once

	// Ensure the handler implements Handle at build time.
	_ otel.ErrorHandler = (*handler)(nil)
)

// handler logs all errors to STDERR.
type handler struct {
	delegate atomic.Value

	l *log.Logger
}

// setDelegate sets the handler delegate if one is not already set.
func (h *handler) setDelegate(d otel.ErrorHandler) {
	if h.delegate.Load() != nil {
		// Delegate already registered
		return
	}
	h.delegate.Store(d)
}

// Handle implements otle.ErrorHandler.
func (h *handler) Handle(err error) {
	if d := h.delegate.Load(); d != nil {
		d.(otel.ErrorHandler).Handle(err)
		return
	}
	h.l.Print(err)
}

// ErrorHandler returns the global ErrorHandler instance. If no ErrorHandler
// instance has be explicitly set yet, a default ErrorHandler is returned
// that logs to STDERR until an ErrorHandler is set (all functionality is
// delegated to the set ErrorHandler once it is set).
func ErrorHandler() otel.ErrorHandler {
	return globalHandler
}

// SetErrorHandler sets the global ErrorHandler to be h.
func SetErrorHandler(h otel.ErrorHandler) {
	delegateHanderOnce.Do(func() {
		current := ErrorHandler()
		if current == h {
			return
		}
		if internalHandler, ok := current.(*handler); ok {
			internalHandler.setDelegate(h)
		}
	})
}

// Handle is a convience function for ErrorHandler().Handle(err)
func Handle(err error) {
	ErrorHandler().Handle(err)
}
