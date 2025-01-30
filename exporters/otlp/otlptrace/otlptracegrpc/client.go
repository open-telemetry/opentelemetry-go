// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlpconfig"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/retry"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

const selfObsScopeName = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

type client struct {
	endpoint      string
	dialOpts      []grpc.DialOption
	metadata      metadata.MD
	exportTimeout time.Duration
	requestFunc   retry.RequestFunc

	// stopCtx is used as a parent context for all exports. Therefore, when it
	// is canceled with the stopFunc all exports are canceled.
	stopCtx context.Context
	// stopFunc cancels stopCtx, stopping any active exports.
	stopFunc context.CancelFunc

	// ourConn keeps track of where conn was created: true if created here on
	// Start, or false if passed with an option. This is important on Shutdown
	// as the conn should only be closed if created here on start. Otherwise,
	// it is up to the processes that passed the conn to close it.
	ourConn bool
	conn    *grpc.ClientConn
	tscMu   sync.RWMutex
	tsc     coltracepb.TraceServiceClient

	spansInflightUpDownCounter metric.Int64UpDownCounter
	spansExportedCounter       metric.Int64Counter
	baseAttributes             metric.MeasurementOption
	successAttributes          metric.MeasurementOption
	exportFailedAttributes     metric.MeasurementOption
}

// Compile time check *client implements otlptrace.Client.
var _ otlptrace.Client = (*client)(nil)

// NewClient creates a new gRPC trace client.
func NewClient(opts ...Option) otlptrace.Client {
	return newClient(opts...)
}

func newClient(opts ...Option) *client {
	cfg := otlpconfig.NewGRPCConfig(asGRPCOptions(opts)...)

	ctx, cancel := context.WithCancel(context.Background())

	c := &client{
		endpoint:      cfg.Traces.Endpoint,
		exportTimeout: cfg.Traces.Timeout,
		requestFunc:   cfg.RetryConfig.RequestFunc(retryable),
		dialOpts:      cfg.DialOptions,
		stopCtx:       ctx,
		stopFunc:      cancel,
		conn:          cfg.GRPCConn,
	}

	if len(cfg.Traces.Headers) > 0 {
		c.metadata = metadata.New(cfg.Traces.Headers)
	}

	c.configureSelfObservability()

	return c
}

var exporterID atomic.Int64

// nextExporterID returns an identifier for this otlp grpc trace exporter,
// starting with 0 and incrementing by 1 each time it is called.
func nextExporterID() int64 {
	return exporterID.Add(1) - 1
}

// configureSelfObservability configures metrics for the batch span processor.
func (c *client) configureSelfObservability() {
	mp := otel.GetMeterProvider()
	if !x.SelfObservability.Enabled() {
		mp = metric.MeterProvider(noop.NewMeterProvider())
	}
	meter := mp.Meter(
		selfObsScopeName,
		metric.WithInstrumentationVersion(otlptrace.Version()),
	)
	var err error
	c.spansInflightUpDownCounter, err = meter.Int64UpDownCounter("otel.sdk.span.exporter.spans_inflight",
		metric.WithUnit("{span}"),
		metric.WithDescription("The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)."),
	)
	if err != nil {
		otel.Handle(err)
	}
	c.spansExportedCounter, err = meter.Int64Counter("otel.sdk.span.exporter.spans_exported",
		metric.WithUnit("{span}"),
		metric.WithDescription("The number of spans for which the export has finished, either successful or failed."),
	)
	if err != nil {
		otel.Handle(err)
	}

	componentTypeAttr := attribute.String("otel.sdk.component.type", "otlp_grpc_span_exporter")
	componentNameAttr := attribute.String("otel.sdk.component.name", fmt.Sprintf("otlp_grpc_span_exporter/%d", nextExporterID()))
	c.baseAttributes = metric.WithAttributes(componentNameAttr, componentTypeAttr)
	c.successAttributes = metric.WithAttributes(componentNameAttr, componentTypeAttr, attribute.String("error.type", ""))
	c.exportFailedAttributes = metric.WithAttributes(componentNameAttr, componentTypeAttr, attribute.String("error.type", "export_failed"))
}

// Start establishes a gRPC connection to the collector.
func (c *client) Start(context.Context) error {
	if c.conn == nil {
		// If the caller did not provide a ClientConn when the client was
		// created, create one using the configuration they did provide.
		conn, err := grpc.NewClient(c.endpoint, c.dialOpts...)
		if err != nil {
			return err
		}
		// Keep track that we own the lifecycle of this conn and need to close
		// it on Shutdown.
		c.ourConn = true
		c.conn = conn
	}

	// The otlptrace.Client interface states this method is called just once,
	// so no need to check if already started.
	c.tscMu.Lock()
	c.tsc = coltracepb.NewTraceServiceClient(c.conn)
	c.tscMu.Unlock()

	return nil
}

var errAlreadyStopped = errors.New("the client is already stopped")

// Stop shuts down the client.
//
// Any active connections to a remote endpoint are closed if they were created
// by the client. Any gRPC connection passed during creation using
// WithGRPCConn will not be closed. It is the caller's responsibility to
// handle cleanup of that resource.
//
// This method synchronizes with the UploadTraces method of the client. It
// will wait for any active calls to that method to complete unimpeded, or it
// will cancel any active calls if ctx expires. If ctx expires, the context
// error will be forwarded as the returned error. All client held resources
// will still be released in this situation.
//
// If the client has already stopped, an error will be returned describing
// this.
func (c *client) Stop(ctx context.Context) error {
	// Make sure to return context error if the context is done when calling this method.
	err := ctx.Err()

	// Acquire the c.tscMu lock within the ctx lifetime.
	acquired := make(chan struct{})
	go func() {
		c.tscMu.Lock()
		close(acquired)
	}()

	select {
	case <-ctx.Done():
		// The Stop timeout is reached. Kill any remaining exports to force
		// the clear of the lock and save the timeout error to return and
		// signal the shutdown timed out before cleanly stopping.
		c.stopFunc()
		err = ctx.Err()

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
		return errAlreadyStopped
	}

	// Clear c.tsc to signal the client is stopped.
	c.tsc = nil

	if c.ourConn {
		closeErr := c.conn.Close()
		// A context timeout error takes precedence over this error.
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}
	return err
}

var errShutdown = errors.New("the client is shutdown")

// UploadTraces sends a batch of spans.
//
// Retryable errors from the server will be handled according to any
// RetryConfig the client was created with.
func (c *client) UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	var numSpans int64
	for _, rs := range protoSpans {
		for _, ss := range rs.GetScopeSpans() {
			numSpans += int64(len(ss.GetSpans()))
		}
	}
	c.spansInflightUpDownCounter.Add(ctx, numSpans, c.baseAttributes)
	defer func() {
		c.spansInflightUpDownCounter.Add(ctx, -numSpans, c.baseAttributes)
	}()
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

	var partialRejected int64

	err := c.requestFunc(ctx, func(iCtx context.Context) error {
		resp, err := c.tsc.Export(iCtx, &coltracepb.ExportTraceServiceRequest{
			ResourceSpans: protoSpans,
		})
		if resp != nil && resp.PartialSuccess != nil {
			msg := resp.PartialSuccess.GetErrorMessage()
			partialRejected = resp.PartialSuccess.GetRejectedSpans()
			if partialRejected != 0 || msg != "" {
				err := internal.TracePartialSuccessError(partialRejected, msg)
				otel.Handle(err)
			}
		}
		// nil is converted to OK.
		if status.Code(err) == codes.OK {
			// Success.
			return nil
		}
		return err
	})
	if err == nil {
		c.spansExportedCounter.Add(ctx, numSpans, c.successAttributes)
	} else if partialRejected == 0 {
		c.spansExportedCounter.Add(ctx, numSpans, c.exportFailedAttributes)
	} else {
		// partial success
		c.spansExportedCounter.Add(ctx, partialRejected, c.exportFailedAttributes)
		c.spansExportedCounter.Add(ctx, numSpans-partialRejected, c.successAttributes)
	}
	return err
}

// exportContext returns a copy of parent with an appropriate deadline and
// cancellation function.
//
// It is the callers responsibility to cancel the returned context once its
// use is complete, via the parent or directly with the returned CancelFunc, to
// ensure all resources are correctly released.
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
		md := c.metadata
		if outMD, ok := metadata.FromOutgoingContext(ctx); ok {
			md = metadata.Join(md, outMD)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
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

// retryable returns if err identifies a request that can be retried and a
// duration to wait for if an explicit throttle time is included in err.
func retryable(err error) (bool, time.Duration) {
	s := status.Convert(err)
	return retryableGRPCStatus(s)
}

func retryableGRPCStatus(s *status.Status) (bool, time.Duration) {
	switch s.Code() {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable,
		codes.DataLoss:
		// Additionally handle RetryInfo.
		_, d := throttleDelay(s)
		return true, d
	case codes.ResourceExhausted:
		// Retry only if the server signals that the recovery from resource exhaustion is possible.
		return throttleDelay(s)
	}

	// Not a retry-able error.
	return false, 0
}

// throttleDelay returns of the status is RetryInfo
// and the its duration to wait for if an explicit throttle time.
func throttleDelay(s *status.Status) (bool, time.Duration) {
	for _, detail := range s.Details() {
		if t, ok := detail.(*errdetails.RetryInfo); ok {
			return true, t.RetryDelay.AsDuration()
		}
	}
	return false, 0
}

// MarshalLog is the marshaling function used by the logging system to represent this Client.
func (c *client) MarshalLog() interface{} {
	return struct {
		Type     string
		Endpoint string
	}{
		Type:     "otlptracegrpc",
		Endpoint: c.endpoint,
	}
}
