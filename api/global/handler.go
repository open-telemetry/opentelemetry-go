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
	defaultLogger  = log.New(os.Stderr, "", log.LstdFlags)
	defaultHandler = handler{l: defaultLogger}
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

// Error implements oterror.Handler.
func (h *handler) Error(err error) {
	h.Lock()
	defer h.Unlock()
	if h.delegate != nil {
		h.delegate.Error(err)
		return
	}
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
