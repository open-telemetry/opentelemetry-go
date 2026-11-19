// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	colprofilepb "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal/retry"
)

const (
	// profilesDevelopmentVersionKey identifies the OTLP Profiles v1development
	// schema version per opentelemetry-proto/docs/specification.md.
	profilesDevelopmentVersionKey = "otlp-profiles-development-version"
	profilesDevelopmentVersion    = "1"
)

type client struct {
	metadata       metadata.MD
	exportTimeout  time.Duration
	maxRequestSize int
	requestFunc    retry.RequestFunc

	ourConn bool
	conn    *grpc.ClientConn
	psc     colprofilepb.ProfilesServiceClient
}

var newGRPCClientFn = grpc.NewClient

func newClient(cfg config) (*client, error) {
	c := &client{
		exportTimeout:  cfg.timeout.Value,
		maxRequestSize: cfg.maxRequestSize.Value,
		requestFunc:    cfg.retryCfg.Value.RequestFunc(retryable),
		conn:           cfg.gRPCConn.Value,
	}

	md := metadata.New(cfg.headers.Value)
	md.Set(profilesDevelopmentVersionKey, profilesDevelopmentVersion)
	c.metadata = md

	if c.conn == nil {
		dialOpts := newGRPCDialOptions(cfg)
		conn, err := newGRPCClientFn(cfg.endpoint.Value, dialOpts...)
		if err != nil {
			return nil, err
		}
		c.ourConn = true
		c.conn = conn
	}

	c.psc = colprofilepb.NewProfilesServiceClient(c.conn)
	return c, nil
}

func newGRPCDialOptions(cfg config) []grpc.DialOption {
	userAgent := "OTel Go OTLP over gRPC profiles exporter/" + Version()
	dialOpts := []grpc.DialOption{grpc.WithUserAgent(userAgent)}

	if cfg.serviceConfig.Value != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(cfg.serviceConfig.Value))
	}
	switch {
	case cfg.gRPCCredentials.Value != nil:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(cfg.gRPCCredentials.Value))
	case cfg.tlsCfg.Value != nil:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(cfg.tlsCfg.Value)))
	case cfg.insecure.Value:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	default:
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}
	if cfg.compression.Value == GzipCompression {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}
	if cfg.reconnectionPeriod.Value != 0 {
		dialOpts = append(dialOpts, grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: cfg.reconnectionPeriod.Value,
		}))
	}

	return append(dialOpts, cfg.dialOptions.Value...)
}

func (c *client) UploadProfiles(
	ctx context.Context,
	pbRequest *colprofilepb.ExportProfilesServiceRequest,
) (uploadErr error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ctx, cancel := c.exportContext(ctx)
	defer cancel()

	if maxSize := c.maxRequestSize; maxSize > 0 && proto.Size(pbRequest) > maxSize {
		return fmt.Errorf("request message too large: exceeded %d bytes", maxSize)
	}

	err := c.requestFunc(ctx, func(ctx context.Context) error {
		resp, err := c.psc.Export(ctx, pbRequest)
		if resp != nil && resp.PartialSuccess != nil {
			msg := resp.PartialSuccess.GetErrorMessage()
			n := resp.PartialSuccess.GetRejectedProfiles()
			if n != 0 || msg != "" {
				uploadErr = errors.Join(uploadErr, internal.ProfilePartialSuccessError(n, msg))
			}
		}
		if status.Code(err) == codes.OK {
			return nil
		}
		return err
	})
	return errors.Join(uploadErr, err)
}

func (c *client) Shutdown(ctx context.Context) error {
	c.metadata = nil
	c.requestFunc = nil
	c.psc = nil

	err := ctx.Err()
	if c.ourConn {
		closeErr := c.conn.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}
	c.conn = nil
	return err
}

func (c *client) exportContext(parent context.Context) (context.Context, context.CancelFunc) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if c.exportTimeout > 0 {
		ctx, cancel = context.WithTimeoutCause(parent, c.exportTimeout, errors.New("exporter export timeout"))
	} else {
		ctx, cancel = context.WithCancel(parent) //nolint:gosec // Canceled by caller.
	}

	if c.metadata.Len() > 0 {
		md := c.metadata.Copy()
		if outMD, ok := metadata.FromOutgoingContext(ctx); ok {
			md = metadata.Join(outMD, md)
			// Enforce a single well-formed development version value per OTLP spec.
			md.Set(profilesDevelopmentVersionKey, profilesDevelopmentVersion)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx, cancel
}

type noopClient struct{}

func newNoopClient() *noopClient { return &noopClient{} }

func (*noopClient) UploadProfiles(context.Context, *colprofilepb.ExportProfilesServiceRequest) error {
	return nil
}

func (*noopClient) Shutdown(context.Context) error { return nil }

func retryable(err error) (bool, time.Duration) {
	s := status.Convert(err)
	switch s.Code() {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable,
		codes.DataLoss:
		_, d := throttleDelay(s)
		return true, d
	case codes.ResourceExhausted:
		return throttleDelay(s)
	}
	return false, 0
}

func throttleDelay(s *status.Status) (bool, time.Duration) {
	for _, detail := range s.Details() {
		if t, ok := detail.(*errdetails.RetryInfo); ok {
			return true, t.RetryDelay.AsDuration()
		}
	}
	return false, 0
}
