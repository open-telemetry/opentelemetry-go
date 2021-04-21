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

package otlpgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/exporters/otlp"
)

func doRequest(ctx context.Context, fn func(context.Context) error, rs otlp.RetrySettings, stopCh chan struct{}) error {
	expBackoff := otlp.NewExponentialConfig(rs)

	for {
		err := fn(ctx)
		if err == nil {
			// request succeeded.
			return nil
		}

		if !rs.Enabled {
			return err
		}

		// We have an error, check gRPC status code.
		st := status.Convert(err)
		if st.Code() == codes.OK {
			// Not really an error, still success.
			return nil
		}

		// Now, this is this a real error.

		if !shouldRetry(st.Code()) {
			// It is not a retryable error, we should not retry.
			return err
		}

		// Need to retry.
		var delay time.Duration

		// Respect server throttling.
		if throttle := getThrottleDuration(st); throttle != 0 {
			delay = throttle
		} else {
			backoffDelay := expBackoff.NextBackOff()
			if backoffDelay == backoff.Stop {
				// throw away the batch
				err = fmt.Errorf("max elapsed time expired %w", err)
				return err
			}
			delay = backoffDelay
		}

		// back-off, but get interrupted when shutting down or request is cancelled or timed out.
		select {
		case <-ctx.Done():
			return fmt.Errorf("request is cancelled or timed out %w", err)
		case <-stopCh:
			return fmt.Errorf("interrupted due to shutdown %w", err)
		case <-time.After(delay):
		}
	}
}

func shouldRetry(code codes.Code) bool {
	switch code {
	case codes.OK:
		// Success. This function should not be called for this code, the best we
		// can do is tell the caller not to retry.
		return false

	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable,
		codes.DataLoss:
		// These are retryable errors.
		return true

	case codes.Unknown,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.FailedPrecondition,
		codes.Unimplemented,
		codes.Internal:
		// These are fatal errors, don't retry.
		return false

	default:
		// Don't retry on unknown codes.
		return false
	}
}

func getThrottleDuration(status *status.Status) time.Duration {
	// See if throttling information is available.
	for _, detail := range status.Details() {
		if t, ok := detail.(*errdetails.RetryInfo); ok {
			if t.RetryDelay.Seconds > 0 || t.RetryDelay.Nanos > 0 {
				// We are throttled. Wait before retrying as requested by the server.
				return time.Duration(t.RetryDelay.Seconds)*time.Second + time.Duration(t.RetryDelay.Nanos)*time.Nanosecond
			}
			return 0
		}
	}
	return 0
}
