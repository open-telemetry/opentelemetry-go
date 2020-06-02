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

	"go.opentelemetry.io/otel/api/oterror"
)

var (
	// globalHandler provides an oterror.Handler that can be used throughout
	// an OpenTelemetry instrumented project. When a user specified Handler
	// is registered (`SetHandler`) all calls to `Handle` will be delegated
	// to the registered Handler.
	globalHandler = &handler{
		l: log.New(os.Stderr, "", log.LstdFlags),
	}

	// delegateHanderOnce ensures that a user provided Handler is only ever
	// registered once.
	delegateHanderOnce sync.Once

	// Ensure the handler implements oterror.Handle at build time.
	_ oterror.Handler = (*handler)(nil)
)

// handler logs all errors to STDERR.
type handler struct {
	delegate atomic.Value

	l *log.Logger
}

// setDelegate sets the handler delegate if one is not already set.
func (h *handler) setDelegate(d oterror.Handler) {
	if h.delegate.Load() != nil {
		// Delegate already registered
		return
	}
	h.delegate.Store(d)
}

// Handle implements oterror.Handler.
func (h *handler) Handle(err error) {
	if d := h.delegate.Load(); d != nil {
		d.(oterror.Handler).Handle(err)
		return
	}
	h.l.Print(err)
}

// Handler returns the global Handler instance. If no Handler instance has
// be explicitly set yet, a default Handler is returned that logs to STDERR
// until an Handler is set (all functionality is delegated to the set
// Handler once it is set).
func Handler() oterror.Handler {
	return globalHandler
}

// SetHandler sets the global Handler to be h.
func SetHandler(h oterror.Handler) {
	delegateHanderOnce.Do(func() {
		current := Handler()
		if current == h {
			return
		}
		if internalHandler, ok := current.(*handler); ok {
			internalHandler.setDelegate(h)
		}
	})
}

// Handle is a convience function for Handler().Handle(err)
func Handle(err error) {
	globalHandler.Handle(err)
}
