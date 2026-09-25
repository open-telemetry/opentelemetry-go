// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"context"
	"errors"
	"sync"
)

type shutdown struct {
	stop func()
	wait func(context.Context) error
}

// Registry coordinates shutdown of Finish-aware aggregations.
type Registry struct {
	mu        sync.Mutex
	stopped   bool
	shutdowns []shutdown
	done      chan struct{}
	err       error
}

// Register adds stop and wait to the functions called by Shutdown. The caller
// must register the functions before exposing the resource they manage.
//
// If Shutdown has already been called, Register calls stop before returning.
// Because the resource has not yet been exposed, it has no work to wait for.
func (r *Registry) Register(stop func(), wait func(context.Context) error) {
	r.mu.Lock()
	if r.stopped {
		r.mu.Unlock()
		stop()
		return
	}
	r.shutdowns = append(r.shutdowns, shutdown{stop: stop, wait: wait})
	r.mu.Unlock()
}

// Shutdown stops all registered resources and waits for their admitted work to
// complete. It returns when all waits complete or ctx is canceled.
func (r *Registry) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	if r.stopped {
		done := r.done
		r.mu.Unlock()
		select {
		case <-done:
			r.mu.Lock()
			err := r.err
			r.mu.Unlock()
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	r.stopped = true
	r.done = make(chan struct{})
	shutdowns := r.shutdowns
	r.shutdowns = nil
	r.mu.Unlock()

	for _, shutdown := range shutdowns {
		shutdown.stop()
	}

	var err error
	for _, shutdown := range shutdowns {
		err = errors.Join(err, shutdown.wait(ctx))
	}
	r.mu.Lock()
	r.err = err
	close(r.done)
	r.mu.Unlock()
	return err
}
