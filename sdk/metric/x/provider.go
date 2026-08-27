// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// MeterProvider is the experimental metric SDK provider. It owns one local
// pipeline for each configured experimental Reader.
type MeterProvider struct {
	noop.MeterProvider

	mu      sync.Mutex
	meters  map[instrumentation.Scope]*meter
	pipes   []*pipeline
	readers []Reader
	stopped atomic.Bool
}

var _ metric.MeterProvider = (*MeterProvider)(nil)

// NewMeterProvider returns a configured experimental MeterProvider.
func NewMeterProvider(options ...Option) *MeterProvider {
	cfg := newConfig(options)
	provider := &MeterProvider{
		meters:  make(map[instrumentation.Scope]*meter),
		readers: cfg.readers,
	}
	for _, reader := range cfg.readers {
		limit, fallback := reader.cardinalityLimit(InstrumentKindCounter)
		if fallback {
			limit = cfg.cardinalityLimit
		}
		pipe := newPipeline(cfg.resource, reader, cfg.views, cfg.exemplarFilter, limit)
		provider.pipes = append(provider.pipes, pipe)
		reader.register(pipe)
	}
	return provider
}

// Meter returns a Meter for name and options.
func (p *MeterProvider) Meter(name string, options ...metric.MeterOption) metric.Meter {
	if p.stopped.Load() {
		return noop.Meter{}
	}
	cfg := metric.NewMeterConfig(options...)
	scope := instrumentation.Scope{
		Name:       name,
		Version:    cfg.InstrumentationVersion(),
		SchemaURL:  cfg.SchemaURL(),
		Attributes: cfg.InstrumentationAttributes(),
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if existing, ok := p.meters[scope]; ok {
		return existing
	}
	meter := &meter{scope: scope, pipes: p.pipes, counters: make(map[counterID]*int64Counter)}
	p.meters[scope] = meter
	return meter
}

// ForceFlush asks all Readers that support ForceFlush to flush.
func (p *MeterProvider) ForceFlush(ctx context.Context) error {
	for _, reader := range p.readers {
		if flusher, ok := reader.(interface{ ForceFlush(context.Context) error }); ok {
			if err := flusher.ForceFlush(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// Shutdown shuts down all Readers. It is idempotent.
func (p *MeterProvider) Shutdown(ctx context.Context) error {
	if p.stopped.Swap(true) {
		return nil
	}
	var first error
	for _, reader := range p.readers {
		if err := reader.Shutdown(ctx); err != nil && first == nil {
			first = err
		}
	}
	return first
}
