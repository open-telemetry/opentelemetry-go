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

	"go.opentelemetry.io/otel/api/oterror"
)

var (
	defaultHandler = &handler{
		l: log.New(os.Stderr, "", log.LstdFlags),
	}

	// Ensure the handler implements oterror.Handle at build time.
	_ oterror.Handler = (*handler)(nil)
)

// handler logs all errors to STDERR.
type handler struct {
	sync.Mutex
	delegate oterror.Handler

	l *log.Logger
}

func (h *handler) setDelegate(d oterror.Handler) {
	h.Lock()
	defer h.Unlock()
	if h.delegate != nil {
		// delegate already registered
		return
	}

	h.delegate = d
}

// Handle implements oterror.Handler.
func (h *handler) Handle(err error) {
	if h.delegate != nil {
		h.delegate.Handle(err)
		return
	}

	h.Lock()
	defer h.Unlock()
	h.l.Print(err)
}

// Handler returns the global Handler instance. If no Handler instance has
// be explicitly set yet, a default Handler is returned that logs to STDERR
// until an Handler is set (all functionality is delegated to the set
// Handler once it is set).
func Handler() oterror.Handler {
	return defaultHandler
}

// SetHandler sets the global Handler to be h.
func SetHandler(h oterror.Handler) {
	defaultHandler.setDelegate(h)
}

// Handle is a convience function for Handler().Handle(err)
func Handle(err error) {
	defaultHandler.Handle(err)
}
