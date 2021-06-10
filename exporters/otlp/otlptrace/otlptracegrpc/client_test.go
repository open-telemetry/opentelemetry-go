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
package otlptracegrpc_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/otlptracetest"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

var roSpans = tracetest.SpanStubs{{Name: "Span 0"}}.Snapshots()

func TestNew_endToEnd(t *testing.T) {
	tests := []struct {
		name           string
		additionalOpts []otlptracegrpc.Option
	}{
		{
			name: "StandardExporter",
		},
		{
			name: "WithCompressor",
			additionalOpts: []otlptracegrpc.Option{
				otlptracegrpc.WithCompressor(gzip.Name),
			},
		},
		{
			name: "WithServiceConfig",
			additionalOpts: []otlptracegrpc.Option{
				otlptracegrpc.WithServiceConfig("{}"),
			},
		},
		{
			name: "WithDialOptions",
			additionalOpts: []otlptracegrpc.Option{
				otlptracegrpc.WithDialOption(grpc.WithBlock()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func newGRPCExporter(t *testing.T, ctx context.Context, endpoint string, additionalOpts ...otlptracegrpc.Option) *otlptrace.Exporter {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	client := otlptracegrpc.NewClient(opts...)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	return exp
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlptracegrpc.Option) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint, additionalOpts...)
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	otlptracetest.RunEndToEndTest(ctx, t, exp, mc)
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")
	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	otlptracetest.RunExporterShutdownTest(t, func() otlptrace.Client {
		return otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(mc.endpoint),
			otlptracegrpc.WithReconnectionPeriod(50*time.Millisecond))
	})
}

func TestNew_invokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		if err := exp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	// Invoke Start numerous times, should return errAlreadyStarted
	for i := 0; i < 10; i++ {
		if err := exp.Start(ctx); err == nil || !strings.Contains(err.Error(), "already started") {
			t.Fatalf("#%d unexpected Start error: %v", i, err)
		}
	}

	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to Shutdown the exporter: %v", err)
	}
	// Invoke Shutdown numerous times
	for i := 0; i < 10; i++ {
		if err := exp.Shutdown(ctx); err != nil {
			t.Fatalf(`#%d got error (%v) expected none`, i, err)
		}
	}
}

func TestNew_collectorConnectionDiesThenReconnectsWhenInRestMode(t *testing.T) {
	mc := runMockCollector(t)

	reconnectionPeriod := 20 * time.Millisecond
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithRetry(otlptracegrpc.RetrySettings{Enabled: false}),
		otlptracegrpc.WithReconnectionPeriod(reconnectionPeriod))
	defer func() { require.NoError(t, exp.Shutdown(ctx)) }()

	// Wait for a connection.
	mc.ln.WaitForConn()

	// We'll now stop the collector right away to simulate a connection
	// dying in the midst of communication or even not existing before.
	require.NoError(t, mc.stop())

	// first export, it will send disconnected message to the channel on export failure,
	// trigger almost immediate reconnection
	require.Error(t, exp.ExportSpans(ctx, roSpans))

	// second export, it will detect connection issue, change state of exporter to disconnected and
	// send message to disconnected channel but this time reconnection gouroutine will be in (rest mode, not listening to the disconnected channel)
	require.Error(t, exp.ExportSpans(ctx, roSpans))

	// as a result we have exporter in disconnected state waiting for disconnection message to reconnect

	// resurrect collector
	nmc := runMockCollectorAtEndpoint(t, mc.endpoint)

	// make sure reconnection loop hits beginning and goes back to waiting mode
	// after hitting beginning of the loop it should reconnect
	nmc.ln.WaitForConn()

	n := 10
	for i := 0; i < n; i++ {
		// when disconnected exp.ExportSpans doesnt send disconnected messages again
		// it just quits and return last connection error
		require.NoError(t, exp.ExportSpans(ctx, roSpans))
	}

	nmaSpans := nmc.getSpans()

	// Expecting 10 spans that were sampled, given that
	if g, w := len(nmaSpans), n; g != w {
		t.Fatalf("Connected collector: spans: got %d want %d", g, w)
	}

	dSpans := mc.getSpans()
	// Expecting 0 spans to have been received by the original but now dead collector
	if g, w := len(dSpans), 0; g != w {
		t.Fatalf("Disconnected collector: spans: got %d want %d", g, w)
	}

	require.NoError(t, nmc.Stop())
}

func TestExporterExportFailureAndRecoveryModes(t *testing.T) {
	tts := []struct {
		name   string
		errors []error
		rs     otlptracegrpc.RetrySettings
		fn     func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector)
		opts   []otlptracegrpc.Option
	}{
		{
			name: "Do not retry if succeeded",
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				require.NoError(t, exp.ExportSpans(ctx, roSpans))

				span := mc.getSpans()

				require.Len(t, span, 1)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 success request.")
			},
		},
		{
			name: "Do not retry if 'error' is ok",
			errors: []error{
				status.Error(codes.OK, ""),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				require.NoError(t, exp.ExportSpans(ctx, roSpans))

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 error OK request.")
			},
		},
		{
			name: "Fail three times and succeed",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  300 * time.Millisecond,
				InitialInterval: 2 * time.Millisecond,
				MaxInterval:     10 * time.Millisecond,
			},
			errors: []error{
				status.Error(codes.Unavailable, "backend under pressure"),
				status.Error(codes.Unavailable, "backend under pressure"),
				status.Error(codes.Unavailable, "backend under pressure"),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				require.NoError(t, exp.ExportSpans(ctx, roSpans))

				span := mc.getSpans()

				require.Len(t, span, 1)
				require.Equal(t, 4, mc.traceSvc.requests, "trace service must receive 3 failure requests and 1 success request.")
			},
		},
		{
			name: "Permanent error should not be retried",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  300 * time.Millisecond,
				InitialInterval: 2 * time.Millisecond,
				MaxInterval:     10 * time.Millisecond,
			},
			errors: []error{
				status.Error(codes.InvalidArgument, "invalid arguments"),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				require.Error(t, exp.ExportSpans(ctx, roSpans))

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 error requests.")
			},
		},
		{
			name: "Test all transient errors and succeed",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  500 * time.Millisecond,
				InitialInterval: 1 * time.Millisecond,
				MaxInterval:     2 * time.Millisecond,
			},
			errors: []error{
				status.Error(codes.Canceled, ""),
				status.Error(codes.DeadlineExceeded, ""),
				status.Error(codes.ResourceExhausted, ""),
				status.Error(codes.Aborted, ""),
				status.Error(codes.OutOfRange, ""),
				status.Error(codes.Unavailable, ""),
				status.Error(codes.DataLoss, ""),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				require.NoError(t, exp.ExportSpans(ctx, roSpans))

				span := mc.getSpans()

				require.Len(t, span, 1)
				require.Equal(t, 8, mc.traceSvc.requests, "trace service must receive 9 failure requests and 1 success request.")
			},
		},
		{
			name: "Retry should honor server throttling",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  time.Minute,
				InitialInterval: time.Nanosecond,
				MaxInterval:     time.Nanosecond,
			},
			opts: []otlptracegrpc.Option{
				otlptracegrpc.WithTimeout(time.Millisecond * 100),
			},
			errors: []error{
				newThrottlingError(codes.ResourceExhausted, time.Second*30),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				err := exp.ExportSpans(ctx, roSpans)
				require.Error(t, err)
				require.Equal(t, "context deadline exceeded", err.Error())

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 failure requests and 1 success request.")
			},
		},
		{
			name: "Retry should fail if server throttling is higher than the MaxElapsedTime",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  time.Millisecond * 100,
				InitialInterval: time.Nanosecond,
				MaxInterval:     time.Nanosecond,
			},
			errors: []error{
				newThrottlingError(codes.ResourceExhausted, time.Minute),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				err := exp.ExportSpans(ctx, roSpans)
				require.Error(t, err)
				require.Equal(t, "max elapsed time expired when respecting server throttle: rpc error: code = ResourceExhausted desc = ", err.Error())

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 failure requests and 1 success request.")
			},
		},
		{
			name: "Retry stops if takes too long",
			rs: otlptracegrpc.RetrySettings{
				Enabled:         true,
				MaxElapsedTime:  time.Millisecond * 100,
				InitialInterval: time.Millisecond * 50,
				MaxInterval:     time.Millisecond * 50,
			},
			errors: []error{
				status.Error(codes.Unavailable, "unavailable"),
				status.Error(codes.Unavailable, "unavailable"),
				status.Error(codes.Unavailable, "unavailable"),
				status.Error(codes.Unavailable, "unavailable"),
				status.Error(codes.Unavailable, "unavailable"),
				status.Error(codes.Unavailable, "unavailable"),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				err := exp.ExportSpans(ctx, roSpans)
				require.Error(t, err)

				require.Equal(t, "max elapsed time expired: rpc error: code = Unavailable desc = unavailable", err.Error())

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.LessOrEqual(t, 1, mc.traceSvc.requests, "trace service must receive at least 1 failure requests.")
			},
		},
		{
			name: "Disabled retry",
			rs: otlptracegrpc.RetrySettings{
				Enabled: false,
			},
			errors: []error{
				status.Error(codes.Unavailable, "unavailable"),
			},
			fn: func(t *testing.T, ctx context.Context, exp *otlptrace.Exporter, mc *mockCollector) {
				err := exp.ExportSpans(ctx, roSpans)
				require.Error(t, err)

				require.Equal(t, "rpc error: code = Unavailable desc = unavailable", err.Error())

				span := mc.getSpans()

				require.Len(t, span, 0)
				require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 failure requests.")
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mc := runMockCollectorWithConfig(t, &mockConfig{
				errors: tt.errors,
			})

			opts := []otlptracegrpc.Option{
				otlptracegrpc.WithRetry(tt.rs),
			}

			if len(tt.opts) != 0 {
				opts = append(opts, tt.opts...)
			}

			exp := newGRPCExporter(t, ctx, mc.endpoint, opts...)

			tt.fn(t, ctx, exp, mc)

			require.NoError(t, mc.Stop())
			require.NoError(t, exp.Shutdown(ctx))
		})
	}

}

func TestPermanentErrorsShouldNotBeRetried(t *testing.T) {
	permanentErrors := []*status.Status{
		status.New(codes.Unknown, "Unknown"),
		status.New(codes.InvalidArgument, "InvalidArgument"),
		status.New(codes.NotFound, "NotFound"),
		status.New(codes.AlreadyExists, "AlreadyExists"),
		status.New(codes.FailedPrecondition, "FailedPrecondition"),
		status.New(codes.Unimplemented, "Unimplemented"),
		status.New(codes.Internal, "Internal"),
		status.New(codes.PermissionDenied, ""),
		status.New(codes.Unauthenticated, ""),
	}

	for _, sts := range permanentErrors {
		t.Run(sts.Code().String(), func(t *testing.T) {
			ctx := context.Background()

			mc := runMockCollectorWithConfig(t, &mockConfig{
				errors: []error{sts.Err()},
			})

			exp := newGRPCExporter(t, ctx, mc.endpoint)

			err := exp.ExportSpans(ctx, roSpans)
			require.Error(t, err)
			require.Len(t, mc.getSpans(), 0)
			require.Equal(t, 1, mc.traceSvc.requests, "trace service must receive 1 permanent error requests.")

			require.NoError(t, mc.Stop())
			require.NoError(t, exp.Shutdown(ctx))
		})
	}
}

func newThrottlingError(code codes.Code, duration time.Duration) error {
	s := status.New(code, "")

	s, _ = s.WithDetails(&errdetails.RetryInfo{RetryDelay: durationpb.New(duration)})

	return s.Err()
}

func TestNew_collectorConnectionDiesThenReconnects(t *testing.T) {
	mc := runMockCollector(t)

	reconnectionPeriod := 50 * time.Millisecond
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithRetry(otlptracegrpc.RetrySettings{Enabled: false}),
		otlptracegrpc.WithReconnectionPeriod(reconnectionPeriod))
	defer func() { require.NoError(t, exp.Shutdown(ctx)) }()

	mc.ln.WaitForConn()

	// We'll now stop the collector right away to simulate a connection
	// dying in the midst of communication or even not existing before.
	require.NoError(t, mc.stop())

	// In the test below, we'll stop the collector many times,
	// while exporting traces and test to ensure that we can
	// reconnect.
	for j := 0; j < 3; j++ {

		// No endpoint up.
		require.Error(t, exp.ExportSpans(ctx, roSpans))

		// Now resurrect the collector by making a new one but reusing the
		// old endpoint, and the collector should reconnect automatically.
		nmc := runMockCollectorAtEndpoint(t, mc.endpoint)

		// Give the exporter sometime to reconnect
		nmc.ln.WaitForConn()

		n := 10
		for i := 0; i < n; i++ {
			require.NoError(t, exp.ExportSpans(ctx, roSpans))
		}

		nmaSpans := nmc.getSpans()
		// Expecting 10 spans that were sampled, given that
		if g, w := len(nmaSpans), n; g != w {
			t.Fatalf("Round #%d: Connected collector: spans: got %d want %d", j, g, w)
		}

		dSpans := mc.getSpans()
		// Expecting 0 spans to have been received by the original but now dead collector
		if g, w := len(dSpans), 0; g != w {
			t.Fatalf("Round #%d: Disconnected collector: spans: got %d want %d", j, g, w)
		}

		// Disconnect for the next try.
		require.NoError(t, nmc.stop())
	}
}

// This test takes a long time to run: to skip it, run tests using: -short
func TestNew_collectorOnBadConnection(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping this long running test")
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to grab an available port: %v", err)
	}
	// Firstly close the "collector's" channel: optimistically this endpoint won't get reused ASAP
	// However, our goal of closing it is to simulate an unavailable connection
	_ = ln.Close()

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	endpoint := fmt.Sprintf("localhost:%s", collectorPortStr)
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNew_withEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNew_withHeaders(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithHeaders(map[string]string{"header1": "value1"}))
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	defer func() {
		_ = exp.Shutdown(ctx)
	}()

	headers := mc.getHeaders()
	require.Len(t, headers.Get("header1"), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}

func TestNew_WithTimeout(t *testing.T) {
	tts := []struct {
		name    string
		fn      func(exp *otlptrace.Exporter) error
		timeout time.Duration
		spans   int
		code    codes.Code
		delay   bool
	}{
		{
			name: "Timeout Spans",
			fn: func(exp *otlptrace.Exporter) error {
				return exp.ExportSpans(context.Background(), roSpans)
			},
			timeout: time.Millisecond * 100,
			code:    codes.DeadlineExceeded,
			delay:   true,
		},

		{
			name: "No Timeout Spans",
			fn: func(exp *otlptrace.Exporter) error {
				return exp.ExportSpans(context.Background(), roSpans)
			},
			timeout: time.Minute,
			spans:   1,
			code:    codes.OK,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {

			mc := runMockCollector(t)
			if tt.delay {
				mc.traceSvc.delay = time.Second * 10
			}
			defer func() {
				_ = mc.stop()
			}()

			ctx := context.Background()
			exp := newGRPCExporter(t, ctx, mc.endpoint, otlptracegrpc.WithTimeout(tt.timeout), otlptracegrpc.WithRetry(otlptracegrpc.RetrySettings{Enabled: false}))
			defer func() {
				_ = exp.Shutdown(ctx)
			}()

			err := tt.fn(exp)

			if tt.code == codes.OK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			s := status.Convert(err)
			require.Equal(t, tt.code, s.Code())

			require.Len(t, mc.getSpans(), tt.spans)
		})
	}
}

func TestNew_withInvalidSecurityConfiguration(t *testing.T) {
	mc := runMockCollector(t)
	defer func() {
		_ = mc.stop()
	}()

	ctx := context.Background()
	driver := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(mc.endpoint))
	exp, err := otlptrace.New(ctx, driver)
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}

	err = exp.ExportSpans(ctx, roSpans)

	expectedErr := fmt.Sprintf("traces exporter is disconnected from the server %s: grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)", mc.endpoint)

	require.Error(t, err)
	require.Equal(t, expectedErr, err.Error())

	defer func() {
		_ = exp.Shutdown(ctx)
	}()
}

func TestNew_withMultipleAttributeTypes(t *testing.T) {
	mc := runMockCollector(t)

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)

	defer func() {
		_ = exp.Shutdown(ctx)
	}()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	defer func() { _ = tp.Shutdown(ctx) }()

	tr := tp.Tracer("test-tracer")
	testKvs := []attribute.KeyValue{
		attribute.Int("Int", 1),
		attribute.Int64("Int64", int64(3)),
		attribute.Float64("Float64", 2.22),
		attribute.Bool("Bool", true),
		attribute.String("String", "test"),
	}
	_, span := tr.Start(ctx, "AlwaysSample")
	span.SetAttributes(testKvs...)
	span.End()

	// Flush and close.
	func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a tracer provider: %v", err)
		}
	}()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to stop the exporter: %v", err)
	}

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	_ = mc.stop()

	// Now verify that we only got one span
	rss := mc.getSpans()
	if got, want := len(rss), 1; got != want {
		t.Fatalf("resource span count: got %d, want %d\n", got, want)
	}

	expected := []*commonpb.KeyValue{
		{
			Key: "Int",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 1,
				},
			},
		},
		{
			Key: "Int64",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 3,
				},
			},
		},
		{
			Key: "Float64",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_DoubleValue{
					DoubleValue: 2.22,
				},
			},
		},
		{
			Key: "Bool",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_BoolValue{
					BoolValue: true,
				},
			},
		},
		{
			Key: "String",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: "test",
				},
			},
		},
	}

	// Verify attributes
	if !assert.Len(t, rss[0].Attributes, len(expected)) {
		t.Fatalf("attributes count: got %d, want %d\n", len(rss[0].Attributes), len(expected))
	}
	for i, actual := range rss[0].Attributes {
		if a, ok := actual.Value.Value.(*commonpb.AnyValue_DoubleValue); ok {
			e, ok := expected[i].Value.Value.(*commonpb.AnyValue_DoubleValue)
			if !ok {
				t.Errorf("expected AnyValue_DoubleValue, got %T", expected[i].Value.Value)
				continue
			}
			if !assert.InDelta(t, e.DoubleValue, a.DoubleValue, 0.01) {
				continue
			}
			e.DoubleValue = a.DoubleValue
		}
		assert.Equal(t, expected[i], actual)
	}
}

func TestDisconnected(t *testing.T) {
	ctx := context.Background()
	// The endpoint is whatever, we want to be disconnected. But we
	// setting a blocking connection, so dialing to the invalid
	// endpoint actually fails.
	exp := newGRPCExporter(t, ctx, "invalid",
		otlptracegrpc.WithReconnectionPeriod(time.Hour),
		otlptracegrpc.WithDialOption(
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		),
	)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.Error(t, exp.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan()))
}

func TestEmptyData(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:56561")

	defer func() {
		_ = mc.stop()
	}()

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	defer func() {
		assert.NoError(t, exp.Shutdown(ctx))
	}()

	assert.NoError(t, exp.ExportSpans(ctx, nil))
}
