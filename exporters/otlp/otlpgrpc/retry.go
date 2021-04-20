package otlpgrpc

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"go.opentelemetry.io/otel/exporters/otlp"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func doRequest(ctx context.Context, fn func(context.Context) error, rs otlp.RetrySettings, stopCh chan struct{}) error {
	expBackoff := otlp.NewExponentialConfig(rs)

	retryNum := 0
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

		// Check if server returned throttling information.
		throttle := getThrottleDuration(st)

		backoffDelay := expBackoff.NextBackOff()
		if backoffDelay == backoff.Stop {
			// throw away the batch
			err = fmt.Errorf("max elapsed time expired %w", err)
			fmt.Println("Exporting failed. No more retries left. Dropping data.", err)
			return err
		}

		// Respect server throttling.
		if throttle > backoffDelay {
			backoffDelay = throttle
		}

		retryNum++

		fmt.Println("Attempt #", retryNum, ": retrying in ", backoffDelay.String())

		// back-off, but get interrupted when shutting down or request is cancelled or timed out.
		select {
		case <-ctx.Done():
			return fmt.Errorf("request is cancelled or timed out %w", err)
		case <-stopCh:
			return fmt.Errorf("interrupted due to shutdown %w", err)
		case <-time.After(backoffDelay):
		}

	}

}

// Send a trace or metrics request to the server. "perform" function is expected to make
// the actual gRPC unary call that sends the request. This function implements the
// common OTLP logic around request handling such as retries and throttling.
func processError(err error) (e error, throttleDur time.Duration) {
	if err == nil {
		// Request is successful, we are done.
		return nil, 0
	}

	// We have an error, check gRPC status code.

	st := status.Convert(err)
	if st.Code() == codes.OK {
		// Not really an error, still success.
		return nil, 0
	}

	// Now, this is this a real error.

	if !shouldRetry(st.Code()) {
		// It is not a retryable error, we should not retry.
		return err, 0
	}

	// Need to retry.

	// Check if server returned throttling information.
	return err, getThrottleDuration(st)
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
