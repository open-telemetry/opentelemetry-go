// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import "sync"

// Registry coordinates shutdown of Finish-aware aggregations.
type Registry struct {
	mu        sync.Mutex
	stopped   bool
	shutdowns []func()
}

// Register adds f to the functions called by Shutdown. If Shutdown has
// already been called, Register calls f before returning.
func (r *Registry) Register(f func()) {
	r.mu.Lock()
	if r.stopped {
		r.mu.Unlock()
		f()
		return
	}
	r.shutdowns = append(r.shutdowns, f)
	r.mu.Unlock()
}

// Shutdown calls all registered shutdown functions once.
func (r *Registry) Shutdown() {
	r.mu.Lock()
	if r.stopped {
		r.mu.Unlock()
		return
	}
	r.stopped = true
	shutdowns := r.shutdowns
	r.shutdowns = nil
	r.mu.Unlock()

	for _, f := range shutdowns {
		f()
	}
}
