// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import (
	"context"

	logpb "go.opentelemetry.io/proto/otlp/logs/v1"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/transform"
	"go.opentelemetry.io/otel/sdk/log"
)

type logClient interface {
	UploadLogs(ctx context.Context, rl []*logpb.ResourceLogs) error
	Shutdown(context.Context) error
}

// Exporter is a OpenTelemetry log Exporter. It transports log data encoded as
// OTLP protobufs using gRPC.
// All Exporters must be created with [New].
// Exporter methods must not be called concurrently.
type Exporter struct {
	client  logClient
	stopped bool
}

// Compile-time check Exporter implements [log.Exporter].
var _ log.Exporter = (*Exporter)(nil)

// New returns a new [Exporter].
//
// It is recommended to use it with a [BatchProcessor]
// or other processor exporting records asynchronously.
func New(_ context.Context, options ...Option) (*Exporter, error) {
	cfg := newConfig(options)
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return newExporter(c), nil
}

func newExporter(c logClient) *Exporter {
	return &Exporter{client: c}
}

var transformResourceLogs = transform.ResourceLogs

// Export transforms and transmits log records to an OTLP receiver.
//
// This method returns nil and drops records if called after Shutdown.
// This method returns an error if the method is canceled by the passed context.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	if e.stopped {
		return nil
	}

	otlp := transformResourceLogs(records)
	if otlp == nil {
		return nil
	}

	return e.client.UploadLogs(ctx, otlp)
}

// Shutdown shuts down the Exporter. Calls to Export or ForceFlush will perform
// no operation after this is called.
func (e *Exporter) Shutdown(ctx context.Context) error {
	if e.stopped {
		return nil
	}
	e.stopped = true

	err := e.client.Shutdown(ctx)
	e.client = newNoopClient()
	return err
}

// ForceFlush does nothing. The Exporter holds no state.
func (*Exporter) ForceFlush(context.Context) error {
	return nil
}
