// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetricgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/oconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/otest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestThrottleDelay(t *testing.T) {
	c := codes.ResourceExhausted
	testcases := []struct {
		status       *status.Status
		wantOK       bool
		wantDuration time.Duration
	}{
		{
			status:       status.New(c, "NoRetryInfo"),
			wantOK:       false,
			wantDuration: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "SingleRetryInfo").WithDetails(
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(15 * time.Millisecond),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 15 * time.Millisecond,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "ErrorInfo").WithDetails(
					&errdetails.ErrorInfo{Reason: "no throttle detail"},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       false,
			wantDuration: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "ErrorAndRetryInfo").WithDetails(
					&errdetails.ErrorInfo{Reason: "with throttle detail"},
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(13 * time.Minute),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 13 * time.Minute,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "DoubleRetryInfo").WithDetails(
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(13 * time.Minute),
					},
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(15 * time.Minute),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			wantOK:       true,
			wantDuration: 13 * time.Minute,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.status.Message(), func(t *testing.T) {
			ok, d := throttleDelay(tc.status)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantDuration, d)
		})
	}
}

func TestRetryable(t *testing.T) {
	retryableCodes := map[codes.Code]bool{
		codes.OK:                 false,
		codes.Canceled:           true,
		codes.Unknown:            false,
		codes.InvalidArgument:    false,
		codes.DeadlineExceeded:   true,
		codes.NotFound:           false,
		codes.AlreadyExists:      false,
		codes.PermissionDenied:   false,
		codes.ResourceExhausted:  false,
		codes.FailedPrecondition: false,
		codes.Aborted:            true,
		codes.OutOfRange:         true,
		codes.Unimplemented:      false,
		codes.Internal:           false,
		codes.Unavailable:        true,
		codes.DataLoss:           true,
		codes.Unauthenticated:    false,
	}

	for c, want := range retryableCodes {
		got, _ := retryable(status.Error(c, ""))
		assert.Equalf(t, want, got, "evaluate(%s)", c)
	}
}

func TestRetryableGRPCStatusResourceExhaustedWithRetryInfo(t *testing.T) {
	delay := 15 * time.Millisecond
	s, err := status.New(codes.ResourceExhausted, "WithRetryInfo").WithDetails(
		&errdetails.RetryInfo{
			RetryDelay: durationpb.New(delay),
		},
	)
	require.NoError(t, err)

	ok, d := retryableGRPCStatus(s)
	assert.True(t, ok)
	assert.Equal(t, delay, d)
}

type clientShim struct {
	*client
}

func (clientShim) Temporality(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (clientShim) Aggregation(metric.InstrumentKind) metric.Aggregation {
	return nil
}

func (clientShim) ForceFlush(ctx context.Context) error {
	return ctx.Err()
}

func TestClient(t *testing.T) {
	factory := func(rCh <-chan otest.ExportResult) (otest.Client, otest.Collector) {
		coll, err := otest.NewGRPCCollector("", rCh)
		require.NoError(t, err)

		ctx := context.Background()
		addr := coll.Addr().String()
		opts := []Option{WithEndpoint(addr), WithInsecure()}
		cfg := oconf.NewGRPCConfig(asGRPCOptions(opts)...)
		client, err := newClient(ctx, cfg)
		require.NoError(t, err)
		return clientShim{client}, coll
	}

	t.Run("Integration", otest.RunClientTests(factory))
}

func TestConfig(t *testing.T) {
	factoryFunc := func(rCh <-chan otest.ExportResult, o ...Option) (metric.Exporter, *otest.GRPCCollector) {
		coll, err := otest.NewGRPCCollector("", rCh)
		require.NoError(t, err)

		ctx := context.Background()
		opts := append([]Option{
			WithEndpoint(coll.Addr().String()),
			WithInsecure(),
		}, o...)
		exp, err := New(ctx, opts...)
		require.NoError(t, err)
		return exp, coll
	}

	t.Run("WithEndpointURL", func(t *testing.T) {
		coll, err := otest.NewGRPCCollector("", nil)
		require.NoError(t, err)
		t.Cleanup(coll.Shutdown)

		ctx := context.Background()
		exp, err := New(ctx, WithEndpointURL("http://"+coll.Addr().String()))
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithHeaders", func(t *testing.T) {
		key := "my-custom-header"
		headers := map[string]string{key: "custom-value"}
		exp, coll := factoryFunc(nil, WithHeaders(headers))
		t.Cleanup(coll.Shutdown)

		ctx := context.Background()
		additionalKey := "additional-custom-header"
		ctx = metadata.AppendToOutgoingContext(ctx, additionalKey, "additional-value")
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Regexp(t, "OTel Go OTLP over gRPC metrics exporter/[01]\\..*", got)
		require.Contains(t, got, key)
		require.Contains(t, got, additionalKey)
		assert.Equal(t, []string{headers[key]}, got[key])
	})

	t.Run("WithTimeout", func(t *testing.T) {
		// Do not send on rCh so the Collector never responds to the client.
		rCh := make(chan otest.ExportResult)
		t.Cleanup(func() { close(rCh) })
		exp, coll := factoryFunc(
			rCh,
			WithTimeout(time.Millisecond),
			WithRetry(RetryConfig{Enabled: false}),
		)
		t.Cleanup(coll.Shutdown)
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, &metricdata.ResourceMetrics{})
		assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
	})

	t.Run("WithCustomUserAgent", func(t *testing.T) {
		key := "user-agent"
		customerUserAgent := "custom-user-agent"
		exp, coll := factoryFunc(nil, WithDialOption(grpc.WithUserAgent(customerUserAgent)))
		t.Cleanup(coll.Shutdown)
		ctx := context.Background()
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		assert.Contains(t, got[key][0], customerUserAgent)
	})
}
