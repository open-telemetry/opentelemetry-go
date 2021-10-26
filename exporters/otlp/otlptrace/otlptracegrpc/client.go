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

package otlptracegrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/otlpconfig"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type client struct {
	endpoint      string
	dialOpts      []grpc.DialOption
	metadata      metadata.MD
	exportTimeout time.Duration

	// stopCtx is used as a parent context for all exports therefore ensuring
	// that when it is canceled with the stopFunc, all exports are canceled.
	stopCtx context.Context
	// stopFunc cancels stopCtx, stopping any active exports.
	stopFunc context.CancelFunc

	conn  *grpc.ClientConn
	tscMu sync.RWMutex
	tsc   coltracepb.TraceServiceClient
}

// Compile time check *client implements otlptrace.Client.
var _ otlptrace.Client = (*client)(nil)

// NewClient creates a new gRPC trace client.
func NewClient(opts ...Option) otlptrace.Client {
	cfg := otlpconfig.NewGRPCConfig(asGRPCOptions(opts)...)

	ctx, cancel := context.WithCancel(context.Background())

	c := &client{
		endpoint:      cfg.Traces.Endpoint,
		exportTimeout: cfg.Traces.Timeout,
		dialOpts:      cfg.DialOptions,
		stopCtx:       ctx,
		stopFunc:      cancel,
		conn:          cfg.GRPCConn,
	}

	if len(cfg.Traces.Headers) > 0 {
		c.metadata = metadata.New(cfg.Traces.Headers)
	}

	return c
}

// Start establishes a gRPC connection to the collector.
func (c *client) Start(ctx context.Context) error {
	if c.conn == nil {
		// If the caller did not provide a ClientConn when the clinet was
		// created, create one using the configuration they did provide.
		conn, err := grpc.DialContext(ctx, c.endpoint, c.dialOpts...)
		if err != nil {
			return err
		}
		c.conn = conn
	}

	// The otlptrace.Client interface states this method is called just once,
	// so no need to check if already started.
	c.tscMu.Lock()
	c.tsc = coltracepb.NewTraceServiceClient(c.conn)
	c.tscMu.Unlock()

	return nil
}

var errNotStarted = errors.New("client not started")

// Stop shuts down the gRPC connection to the collector.
func (c *client) Stop(ctx context.Context) error {
	// Acquire the c.tscMu lock within the ctx lifetime.
	acquired := make(chan struct{})
	go func() {
		c.tscMu.Lock()
		close(acquired)
	}()
	var ctxErr error
	select {
	case <-ctx.Done():
		// The Stop timeout is reached. Kill any remaining exports to force
		// the clear of the lock and save the timeout error to return and
		// signal the shutdown timed out before cleaning stopping.
		c.stopFunc()
		ctxErr = ctx.Err()

		// To ensure the client is not left in a dirty state c.tsc needs to be
		// set to nil. To avoid the race condition when doing this, ensure
		// that all the exports are killed (initiated by c.stopFunc).
		<-acquired
	case <-acquired:
	}
	// Hold the tscMu lock for the rest of the function to ensure no new
	// exports are started.
	defer c.tscMu.Unlock()

	// The otlptrace.Client interface states this method is called only
	// once, but there is no guarantee it is called after Start. Ensure the
	// client is started before doing anything and let the called know if they
	// made a mistake.
	if c.tsc == nil {
		return errNotStarted
	}

	// Clear c.tsc to signal the client is stopped.
	c.tsc = nil

	connErr := c.conn.Close()
	if ctxErr != nil {
		return ctxErr
	}
	return connErr
}

var errShutdown = errors.New("exporter is shutdown")

// UploadTraces sends a batch of spans to the collector.
func (c *client) UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	// Hold a read lock to ensure a shut down initiated after this starts does
	// not abandon the export. This read lock acquire has less priority than a
	// write lock acquire (i.e. Stop), meaning if the client is shutting down
	// this will come after the shut down.
	c.tscMu.RLock()
	defer c.tscMu.RUnlock()

	if c.tsc == nil {
		return errShutdown
	}

	ctx, cancel := c.exportContext(ctx)
	defer cancel()

	_, err := c.tsc.Export(ctx, &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: protoSpans,
	})
	return err
}

func (c *client) exportContext(parent context.Context) (context.Context, context.CancelFunc) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if c.exportTimeout > 0 {
		ctx, cancel = context.WithTimeout(parent, c.exportTimeout)
	} else {
		ctx, cancel = context.WithCancel(parent)
	}

	if c.metadata.Len() > 0 {
		ctx = metadata.NewOutgoingContext(ctx, c.metadata)
	}

	// Unify the client stopCtx with the parent.
	go func() {
		select {
		case <-ctx.Done():
		case <-c.stopCtx.Done():
			// Cancel the export as the shutdown has timed out.
			cancel()
		}
	}()

	return ctx, cancel
}
