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

package connection

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/retry"
)

func TestThrottleDuration(t *testing.T) {
	c := codes.ResourceExhausted
	testcases := []struct {
		status   *status.Status
		expected time.Duration
	}{
		{
			status:   status.New(c, "no retry info"),
			expected: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "single retry info").WithDetails(
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(15 * time.Millisecond),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			expected: 15 * time.Millisecond,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "error info").WithDetails(
					&errdetails.ErrorInfo{Reason: "no throttle detail"},
				)
				require.NoError(t, err)
				return s
			}(),
			expected: 0,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "error and retry info").WithDetails(
					&errdetails.ErrorInfo{Reason: "with throttle detail"},
					&errdetails.RetryInfo{
						RetryDelay: durationpb.New(13 * time.Minute),
					},
				)
				require.NoError(t, err)
				return s
			}(),
			expected: 13 * time.Minute,
		},
		{
			status: func() *status.Status {
				s, err := status.New(c, "double retry info").WithDetails(
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
			expected: 13 * time.Minute,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.status.Message(), func(t *testing.T) {
			require.Equal(t, tc.expected, throttleDelay(tc.status))
		})
	}
}

func TestEvaluate(t *testing.T) {
	retryable := map[codes.Code]bool{
		codes.OK:                 false,
		codes.Canceled:           true,
		codes.Unknown:            false,
		codes.InvalidArgument:    false,
		codes.DeadlineExceeded:   true,
		codes.NotFound:           false,
		codes.AlreadyExists:      false,
		codes.PermissionDenied:   false,
		codes.ResourceExhausted:  true,
		codes.FailedPrecondition: false,
		codes.Aborted:            true,
		codes.OutOfRange:         true,
		codes.Unimplemented:      false,
		codes.Internal:           false,
		codes.Unavailable:        true,
		codes.DataLoss:           true,
		codes.Unauthenticated:    false,
	}

	for c, want := range retryable {
		got, _ := evaluate(status.Error(c, ""))
		assert.Equalf(t, want, got, "evaluate(%s)", c)
	}
}

func TestDoRequest(t *testing.T) {
	ev := func(error) (bool, time.Duration) { return false, 0 }

	c := new(Connection)
	c.requestFunc = retry.Config{}.RequestFunc(ev)
	c.stopCh = make(chan struct{})

	ctx := context.Background()
	assert.NoError(t, c.DoRequest(ctx, func(ctx context.Context) error {
		return nil
	}))
	assert.NoError(t, c.DoRequest(ctx, func(ctx context.Context) error {
		return status.Error(codes.OK, "")
	}))
	assert.ErrorIs(t, c.DoRequest(ctx, func(ctx context.Context) error {
		return assert.AnError
	}), assert.AnError)
}
