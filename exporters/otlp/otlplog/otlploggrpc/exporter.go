// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/transform"
	"go.opentelemetry.io/otel/sdk/log"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

// Exporter is a OpenTelemetry log Exporter. It transports log data encoded as
// OTLP protobufs using gRPC.
type Exporter struct {
	// Ensure synchronous access to the client across all functionality.
	clientMu sync.Mutex
	client   interface {
		UploadLogs(ctx context.Context, rl []*logpb.ResourceLogs) error
		Shutdown(context.Context) error
	}

	stopped atomic.Bool
}

// Compile-time check Exporter implements [log.Exporter].
var _ log.Exporter = (*Exporter)(nil)

// New returns a new [Exporter].
func New(_ context.Context, options ...Option) (*Exporter, error) {
	cfg := newConfig(options)
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return newExporter(c)
}

func newExporter(c *client) (*Exporter, error) {
	var e Exporter
	e.client = c
	return &e, nil
}

var transformResourceLogs = transform.ResourceLogs

// Export transforms and transmits log records to an OTLP receiver.
//
// This method returns nil error if called after Shutdown.
// This method returns an error if the method is canceled by the passed context.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if e.stopped.Load() {
		return nil
	}

	otlp := transformResourceLogs(records)
	if otlp == nil {
		return nil
	}

	e.clientMu.Lock()
	defer e.clientMu.Unlock()
	return e.client.UploadLogs(ctx, otlp)
}

// Shutdown shuts down the Exporter. Calls to Export or ForceFlush will perform
// no operation after this is called.
func (e *Exporter) Shutdown(ctx context.Context) error {
	if e.stopped.Swap(true) {
		return nil
	}

	e.clientMu.Lock()
	defer e.clientMu.Unlock()

	err := e.client.Shutdown(ctx)
	e.client = newNoopClient()
	return err
}

// ForceFlush does nothing. The Exporter holds no state.
func (e *Exporter) ForceFlush(ctx context.Context) error {
	return nil
}
