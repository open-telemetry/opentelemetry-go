// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package processcontext

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
)

// Publisher manages the OTEL_CTX memory-mapped region for this process.
// External readers (e.g. OBI or eBPF profiler) discover and read resource
// attributes via /proc/<pid>/maps.
//
// At most one Publisher should be active per process at a time. Publisher is
// safe for concurrent use.
type Publisher struct {
	impl *publisher
}

// NewPublisher creates a Publisher and publishes r as the initial process
// context. It returns an error if the mapping cannot be created or if the
// serialized payload exceeds [MaxPayloadSize].
//
// Call [Publisher.Shutdown] to release the mapping when done.
func NewPublisher(r *resource.Resource, opts ...Option) (*Publisher, error) {
	cfg := newConfig(opts)
	impl, err := newPublisher(r, cfg)
	if err != nil {
		return nil, err
	}
	return &Publisher{impl: impl}, nil
}

// Update republishes updated resource attributes. It returns an error if the
// serialized payload exceeds [MaxPayloadSize] or if the Publisher has been
// shut down.
func (p *Publisher) Update(r *resource.Resource) error {
	return p.impl.update(r)
}

// Shutdown zeros the timestamp (signaling the context is unavailable) and
// unmaps the region. After Shutdown, calls to Update return an error.
func (p *Publisher) Shutdown(_ context.Context) error {
	p.impl.shutdown()
	return nil
}

// MaxPayloadSize is the maximum size in bytes of the serialized
// [ProcessContext] payload.
const MaxPayloadSize = 65536

// Option configures a [Publisher].
type Option interface {
	apply(*config)
}

type config struct{}

func newConfig(opts []Option) config {
	var cfg config
	for _, o := range opts {
		o.apply(&cfg)
	}
	return cfg
}
