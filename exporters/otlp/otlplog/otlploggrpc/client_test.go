// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"github.com/stretchr/testify/require"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestNewClient(t *testing.T) {
	newGRPCClientSwap := newGRPCClient
	t.Cleanup(func() {
		newGRPCClient = newGRPCClientSwap
	})

	// The gRPC connection created by newClient.
	conn, err := grpc.NewClient("test", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	newGRPCClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return conn, nil
	}

	// The gRPC connection created by users.
	userConn, err := grpc.NewClient("test 2", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	testCases := []struct {
		name string
		cfg  config
		cli  *client
	}{
		{
			name: "empty config",
			cli: &client{
				ourConn: true,
				conn:    conn,
				lsc:     collogpb.NewLogsServiceClient(conn),
			},
		},
		{
			name: "with headers",
			cfg: config{
				headers: newSetting(map[string]string{
					"key": "value",
				}),
			},
			cli: &client{
				ourConn:  true,
				conn:     conn,
				lsc:      collogpb.NewLogsServiceClient(conn),
				metadata: map[string][]string{"key": {"value"}},
			},
		},
		{
			name: "with gRPC connection",
			cfg: config{
				gRPCConn: newSetting(userConn),
			},
			cli: &client{
				ourConn: false,
				conn:    userConn,
				lsc:     collogpb.NewLogsServiceClient(userConn),
			},
		},
		{
			// It is not possible to compare grpc dial options directly, so we just check that the client is created
			// and no panic occurs.
			name: "with dial options",
			cfg: config{
				serviceConfig:      newSetting("service config"),
				gRPCCredentials:    newSetting(credentials.NewTLS(nil)),
				compression:        newSetting(GzipCompression),
				reconnectionPeriod: newSetting(10 * time.Second),
			},
			cli: &client{
				ourConn: true,
				conn:    conn,
				lsc:     collogpb.NewLogsServiceClient(conn),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cli, err := newClient(tc.cfg)
			require.NoError(t, err)

			assert.Equal(t, tc.cli.metadata, cli.metadata)
			assert.Equal(t, tc.cli.exportTimeout, cli.exportTimeout)
			assert.Equal(t, tc.cli.ourConn, cli.ourConn)
			assert.Equal(t, tc.cli.conn, cli.conn)
			assert.Equal(t, tc.cli.lsc, cli.lsc)
		})
	}
}
