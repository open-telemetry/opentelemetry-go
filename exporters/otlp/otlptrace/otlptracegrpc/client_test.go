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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	// TODO: Fix this test #1527
	t.Skip("This test is flaky and needs to be rewritten")
	mc := runMockCollector(t)

	reconnectionPeriod := 20 * time.Millisecond
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{Enabled: false}),
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

func TestNew_collectorConnectionDiesThenReconnects(t *testing.T) {
	// TODO: Fix this test #1527
	t.Skip("This test is flaky and needs to be rewritten")
	mc := runMockCollector(t)

	reconnectionPeriod := 50 * time.Millisecond
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{Enabled: false}),
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
			exp := newGRPCExporter(t, ctx, mc.endpoint, otlptracegrpc.WithTimeout(tt.timeout), otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{Enabled: false}))
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
