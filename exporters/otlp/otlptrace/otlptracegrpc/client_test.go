// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/otlptracetest"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

var roSpans = tracetest.SpanStubs{{Name: "Span 0"}}.Snapshots()

func contextWithTimeout(parent context.Context, t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	d, ok := t.Deadline()
	if !ok {
		d = time.Now().Add(timeout)
	} else {
		d = d.Add(-1 * time.Millisecond)
		now := time.Now()
		if d.Sub(now) > timeout {
			d = now.Add(timeout)
		}
	}
	return context.WithDeadline(parent, d)
}

func TestNewEndToEnd(t *testing.T) {
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
	mc := runMockCollector(t)

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint, additionalOpts...)
	t.Cleanup(func() {
		ctx, cancel := contextWithTimeout(ctx, t, 10*time.Second)
		defer cancel()

		require.NoError(t, exp.Shutdown(ctx))
	})

	// RunEndToEndTest closes mc.
	otlptracetest.RunEndToEndTest(ctx, t, exp, mc)
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:0")
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	factory := func() otlptrace.Client {
		return otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(mc.endpoint),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
	}
	otlptracetest.RunExporterShutdownTest(t, factory)
}

func TestNewInvokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

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

// This test takes a long time to run: to skip it, run tests using: -short.
func TestNewCollectorOnBadConnection(t *testing.T) {
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

func TestNewWithEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNewWithHeaders(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithHeaders(map[string]string{"header1": "value1"}))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	headers := mc.getHeaders()
	require.Regexp(t, "OTel OTLP Exporter Go/1\\..*", headers.Get("user-agent"))
	require.Len(t, headers.Get("header1"), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}

func TestExportSpansTimeoutHonored(t *testing.T) {
	ctx, cancel := contextWithTimeout(context.Background(), t, 1*time.Minute)
	t.Cleanup(cancel)

	mc := runMockCollector(t)
	exportBlock := make(chan struct{})
	mc.traceSvc.exportBlock = exportBlock
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	exp := newGRPCExporter(
		t,
		ctx,
		mc.endpoint,
		otlptracegrpc.WithTimeout(1*time.Nanosecond),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{Enabled: false}),
	)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	err := exp.ExportSpans(ctx, roSpans)
	// Release the export so everything is cleaned up on shutdown.
	close(exportBlock)

	require.Equal(t, codes.DeadlineExceeded, status.Convert(err).Code())
}

func TestNewWithMultipleAttributeTypes(t *testing.T) {
	mc := runMockCollector(t)

	<-time.After(5 * time.Millisecond)

	ctx, cancel := contextWithTimeout(context.Background(), t, 10*time.Second)
	t.Cleanup(cancel)

	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	t.Cleanup(func() { require.NoError(t, tp.Shutdown(ctx)) })

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
		ctx, cancel := contextWithTimeout(ctx, t, 10*time.Second)
		defer cancel()
		require.NoError(t, tp.Shutdown(ctx))
	}()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	require.NoError(t, exp.Shutdown(ctx))

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	require.NoError(t, mc.stop())

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

func TestStartErrorInvalidAddress(t *testing.T) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		// Validate the connection in Start (which should return the error).
		otlptracegrpc.WithDialOption(
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		),
		otlptracegrpc.WithEndpoint("invalid"),
		otlptracegrpc.WithReconnectionPeriod(time.Hour),
	)
	err := client.Start(context.Background())
	assert.EqualError(t, err, `connection error: desc = "transport: error while dialing: dial tcp: address invalid: missing port in address"`)
}

func TestEmptyData(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	<-time.After(5 * time.Millisecond)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	assert.NoError(t, exp.ExportSpans(ctx, nil))
}

func TestPartialSuccess(t *testing.T) {
	mc := runMockCollectorWithConfig(t, &mockConfig{
		partial: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 2,
			ErrorMessage:  "partially successful",
		},
	})
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	errors := []error{}
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errors = append(errors, err)
	}))
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	require.Equal(t, 1, len(errors))
	require.Contains(t, errors[0].Error(), "partially successful")
	require.Contains(t, errors[0].Error(), "2 spans rejected")
}

func TestCustomUserAgent(t *testing.T) {
	customUserAgent := "custom-user-agent"
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithDialOption(grpc.WithUserAgent(customUserAgent)))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	headers := mc.getHeaders()
	require.Contains(t, headers.Get("user-agent")[0], customUserAgent)
}
