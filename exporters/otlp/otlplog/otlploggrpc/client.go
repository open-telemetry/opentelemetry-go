// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

// The methods of this type are not expected to be called concurrently.
type client struct {
	metadata      metadata.MD
	exportTimeout time.Duration
	requestFunc   retry.RequestFunc

	// ourConn keeps track of where conn was created: true if created here in
	// NewClient, or false if passed with an option. This is important on
	// Shutdown as conn should only be closed if we created it. Otherwise,
	// it is up to the processes that passed conn to close it.
	ourConn bool
	conn    *grpc.ClientConn
	lsc     collogpb.LogsServiceClient
}

// Used for testing.
var newGRPCClient = grpc.NewClient

// newClient creates a new gRPC log client.
func newClient(cfg config) (*client, error) {
	c := &client{
		exportTimeout: cfg.timeout.Value,
		requestFunc:   cfg.retryCfg.Value.RequestFunc(retryable),
		conn:          cfg.gRPCConn.Value,
	}

	if len(cfg.headers.Value) > 0 {
		c.metadata = metadata.New(cfg.headers.Value)
	}

	if c.conn == nil {
		// If the caller did not provide a ClientConn when the client was
		// created, create one using the configuration they did provide.
		dialOpts := newGRPCDialOptions(cfg)

		conn, err := newGRPCClient(cfg.endpoint.Value, dialOpts...)
		if err != nil {
			return nil, err
		}
		// Keep track that we own the lifecycle of this conn and need to close
		// it on Shutdown.
		c.ourConn = true
		c.conn = conn
	}

	c.lsc = collogpb.NewLogsServiceClient(c.conn)

	return c, nil
}

func newGRPCDialOptions(cfg config) []grpc.DialOption {
	userAgent := "OTel Go OTLP over gRPC logs exporter/" + Version()
	dialOpts := []grpc.DialOption{grpc.WithUserAgent(userAgent)}
	dialOpts = append(dialOpts, cfg.dialOptions.Value...)

	// Convert other grpc configs to the dial options.
	// Service config
	if cfg.serviceConfig.Value != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(cfg.serviceConfig.Value))
	}
	// Prioritize GRPCCredentials over Insecure (passing both is an error).
	if cfg.gRPCCredentials.Value != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(cfg.gRPCCredentials.Value))
	} else if cfg.insecure.Value {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Default to using the host's root CA.
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(nil),
		))
	}
	// Compression
	if cfg.compression.Value == GzipCompression {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}
	// Reconnection period
	if cfg.reconnectionPeriod.Value != 0 {
		p := grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: cfg.reconnectionPeriod.Value,
		}
		dialOpts = append(dialOpts, grpc.WithConnectParams(p))
	}

	return dialOpts
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
		// Additionally, handle RetryInfo.
		_, d := throttleDelay(s)
		return true, d
	case codes.ResourceExhausted:
		// Retry only if the server signals that the recovery from resource exhaustion is possible.
		return throttleDelay(s)
	}

	// Not a retry-able error.
	return false, 0
}

// throttleDelay returns if the status is RetryInfo
// and the duration to wait for if an explicit throttle time is included.
func throttleDelay(s *status.Status) (bool, time.Duration) {
	for _, detail := range s.Details() {
		if t, ok := detail.(*errdetails.RetryInfo); ok {
			return true, t.RetryDelay.AsDuration()
		}
	}
	return false, 0
}
