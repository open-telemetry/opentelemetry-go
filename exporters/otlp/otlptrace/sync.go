// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptrace

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// SyncExporter exports trace data using a SyncClient with an explicit
// synchronous ownership contract. It reuses arena-backed allocations across
// exports safely because SyncClient.UploadTracesSync guarantees the data is
// fully consumed before returning.
type SyncExporter struct {
	client SyncClient

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once

	cfg  syncConfig
	pool sync.Pool
}

// syncConfig holds configuration for SyncExporter.
type syncConfig struct {
	initialSize     int
	maxRetainedSize int
}

// Option applies an option to a SyncExporter.
type Option interface {
	apply(syncConfig) syncConfig
}

type optionFunc func(syncConfig) syncConfig

func (fn optionFunc) apply(c syncConfig) syncConfig { return fn(c) }

func newSyncConfig(opts ...Option) syncConfig {
	cfg := syncConfig{
		initialSize:     32,
		maxRetainedSize: 512,
	}
	for _, o := range opts {
		cfg = o.apply(cfg)
	}
	if cfg.initialSize <= 0 {
		cfg.initialSize = 32
	}
	if cfg.maxRetainedSize <= 0 {
		cfg.maxRetainedSize = 512
	}
	return cfg
}

// WithInitialBatchSize sets the initial arena batch size hint for the
// SyncExporter. Larger values reduce early growth for predictable workloads.
func WithInitialBatchSize(n int) Option {
	return optionFunc(func(c syncConfig) syncConfig {
		c.initialSize = n
		return c
	})
}

// WithMaxRetainedBatchSize sets the maximum batch size whose arena will be
// retained in the pool. Arenas that grow beyond this size are discarded
// after use to bound memory retention.
func WithMaxRetainedBatchSize(n int) Option {
	return optionFunc(func(c syncConfig) syncConfig {
		c.maxRetainedSize = n
		return c
	})
}

// NewSync constructs a new SyncExporter and starts it.
func NewSync(ctx context.Context, client SyncClient, opts ...Option) (*SyncExporter, error) {
	exp := NewSyncUnstarted(client, opts...)
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}

// NewSyncUnstarted constructs a new SyncExporter and does not start it.
func NewSyncUnstarted(client SyncClient, opts ...Option) *SyncExporter {
	cfg := newSyncConfig(opts...)
	exp := &SyncExporter{
		client: client,
		cfg:    cfg,
	}
	exp.pool.New = func() any {
		return tracetransform.NewArena(cfg.initialSize)
	}
	return exp
}

// ExportSpans exports a batch of spans synchronously.
func (e *SyncExporter) ExportSpans(ctx context.Context, ss []tracesdk.ReadOnlySpan) error {
	if len(ss) == 0 {
		return nil
	}

	// Acquire arena from pool or create new.
	var arena *tracetransform.Arena
	if v := e.pool.Get(); v != nil {
		arena = v.(*tracetransform.Arena)
	} else {
		arena = tracetransform.NewArena(e.cfg.initialSize)
	}

	protoSpans := tracetransform.SpansWithArena(ss, arena)
	if len(protoSpans) == 0 {
		arena.Reset()
		if !arena.Exceeds(e.cfg.maxRetainedSize) {
			e.pool.Put(arena)
		}
		return nil
	}

	err := e.client.UploadTracesSync(ctx, protoSpans)

	// Arena can be safely reset and reused only after UploadTracesSync
	// returns, per the SyncClient contract.
	arena.Reset()
	if !arena.Exceeds(e.cfg.maxRetainedSize) {
		e.pool.Put(arena)
	}

	if err != nil {
		return fmt.Errorf("traces export: %w", err)
	}
	return nil
}

// Start establishes a connection to the receiving endpoint.
func (e *SyncExporter) Start(ctx context.Context) error {
	err := errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()
		err = e.client.Start(ctx)
	})
	return err
}

// Shutdown flushes all exports and closes all connections to the receiving endpoint.
func (e *SyncExporter) Shutdown(ctx context.Context) error {
	e.mu.RLock()
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	var err error
	e.stopOnce.Do(func() {
		err = e.client.Stop(ctx)
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
	})

	return err
}

var _ tracesdk.SpanExporter = (*SyncExporter)(nil)

// MarshalLog is the marshaling function used by the logging system to represent this SyncExporter.
func (e *SyncExporter) MarshalLog() any {
	return struct {
		Type   string
		Client string
	}{
		Type:   "otlptrace-sync",
		Client: fmt.Sprintf("%T", e.client),
	}
}
