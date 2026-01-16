// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/counter"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/observ"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlptracetest"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

var roSpans = tracetest.SpanStubs{{Name: "Span 0"}}.Snapshots()

func contextWithTimeout(
	parent context.Context,
	t *testing.T,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
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
				otlptracegrpc.WithDialOption(
					grpc.WithConnectParams(grpc.ConnectParams{
						Backoff:           backoff.DefaultConfig,
						MinConnectTimeout: time.Second,
					})),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func TestWithEndpointURL(t *testing.T) {
	mc := runMockCollector(t)

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
	exp := newGRPCExporter(t, ctx, "", []otlptracegrpc.Option{
		otlptracegrpc.WithEndpointURL("http://" + mc.endpoint),
	}...)
	t.Cleanup(func() {
		ctx, cancel := contextWithTimeout(ctx, t, 10*time.Second)
		defer cancel()

		require.NoError(t, exp.Shutdown(ctx))
	})

	// RunEndToEndTest closes mc.
	otlptracetest.RunEndToEndTest(ctx, t, exp, mc)
}

func newGRPCExporter(
	tb testing.TB,
	ctx context.Context,
	endpoint string,
	additionalOpts ...otlptracegrpc.Option,
) *otlptrace.Exporter {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	client := otlptracegrpc.NewClient(opts...)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		tb.Fatalf("failed to create a new collector exporter: %v", err)
	}
	return exp
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlptracegrpc.Option) {
	mc := runMockCollector(t)

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
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
		)
	}
	otlptracetest.RunExporterShutdownTest(t, factory)
}

func TestNewInvokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	// Invoke Start numerous times, should return errAlreadyStarted
	for i := range 10 {
		if err := exp.Start(ctx); err == nil || !strings.Contains(err.Error(), "already started") {
			t.Fatalf("#%d unexpected Start error: %v", i, err)
		}
	}

	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to Shutdown the exporter: %v", err)
	}
	// Invoke Shutdown numerous times
	for i := range 10 {
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
	ctx := t.Context()
	exp := newGRPCExporter(t, ctx, endpoint)
	require.NoError(t, exp.Shutdown(ctx))
}

func TestNewWithEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := t.Context()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	require.NoError(t, exp.Shutdown(ctx))
}

func TestNewWithHeaders(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
	additionalKey := "additional-custom-header"
	ctx = metadata.AppendToOutgoingContext(ctx, additionalKey, "additional-value")
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithHeaders(map[string]string{"header1": "value1"}))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	headers := mc.getHeaders()
	require.Regexp(t, "OTel OTLP Exporter Go/1\\..*", headers.Get("user-agent"))
	require.Len(t, headers.Get("header1"), 1)
	require.Len(t, headers.Get(additionalKey), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}

func TestExportSpansTimeoutHonored(t *testing.T) {
	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
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

	unwrapped := errors.Unwrap(err)
	require.Equal(t, codes.DeadlineExceeded, status.Convert(unwrapped).Code())
	require.True(t, strings.HasPrefix(err.Error(), "traces export: "), "%+v", err)
}

func TestNewWithMultipleAttributeTypes(t *testing.T) {
	mc := runMockCollector(t)

	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
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

func TestEmptyData(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
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

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	err := exp.ExportSpans(ctx, roSpans)
	want := internal.TracePartialSuccessError(0, "")
	assert.ErrorIs(t, err, want)
}

func TestCustomUserAgent(t *testing.T) {
	customUserAgent := "custom-user-agent"
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlptracegrpc.WithDialOption(grpc.WithUserAgent(customUserAgent)))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.ExportSpans(ctx, roSpans))

	headers := mc.getHeaders()
	require.Contains(t, headers.Get("user-agent")[0], customUserAgent)
}

func TestClientInstrumentation(t *testing.T) {
	// Enable instrumentation for this test.
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Reset client ID to be deterministic
	const id = 0
	counter.SetExporterID(id)

	// Save original meter provider and restore at end of test.
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	// Create a new meter provider to capture metrics.
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	const n, msg = 2, "partially successful"
	mc := runMockCollectorWithConfig(t, &mockConfig{
		endpoint: "localhost:0", // Determine canonical endpoint.
		partial: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: n,
			ErrorMessage:  msg,
		},
	})
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	exp := newGRPCExporter(t, t.Context(), mc.endpoint)
	err := exp.ExportSpans(t.Context(), roSpans)
	assert.ErrorIs(t, err, internal.TracePartialSuccessError(n, msg))
	require.NoError(t, exp.Shutdown(t.Context()))

	var got metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &got))

	attrs := observ.BaseAttrs(id, canonical(t, mc.endpoint))
	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      observ.ScopeName,
			Version:   observ.Version,
			SchemaURL: observ.SchemaURL,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        otelconv.SDKExporterSpanInflight{}.Name(),
				Description: otelconv.SDKExporterSpanInflight{}.Description(),
				Unit:        otelconv.SDKExporterSpanInflight{}.Unit(),
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...)},
					},
					Temporality: metricdata.CumulativeTemporality,
				},
			},
			{
				Name:        otelconv.SDKExporterSpanExported{}.Name(),
				Description: otelconv.SDKExporterSpanExported{}.Description(),
				Unit:        otelconv.SDKExporterSpanExported{}.Unit(),
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...)},
						{Attributes: attribute.NewSet(append(
							attrs,
							otelconv.SDKExporterSpanExported{}.AttrErrorType("*errors.joinError"),
						)...)},
					},
					Temporality: 0x1,
					IsMonotonic: true,
				},
			},
			{
				Name:        otelconv.SDKExporterOperationDuration{}.Name(),
				Description: otelconv.SDKExporterOperationDuration{}.Description(),
				Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
				Data: metricdata.Histogram[float64]{
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{Attributes: attribute.NewSet(append(
							attrs,
							otelconv.SDKExporterOperationDuration{}.AttrErrorType("*errors.joinError"),
							attribute.Int64("rpc.grpc.status_code", int64(codes.OK)),
						)...)},
					},
					Temporality: 0x1,
				},
			},
		},
	}
	require.Len(t, got.ScopeMetrics, 1)
	opt := []metricdatatest.Option{
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreExemplars(),
		metricdatatest.IgnoreValue(),
	}
	metricdatatest.AssertEqual(t, want, got.ScopeMetrics[0], opt...)
}

func canonical(t *testing.T, endpoint string) string {
	t.Helper()

	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	c, err := grpc.NewClient(endpoint, opt) // Used to normaliz endpoint.
	if err != nil {
		t.Fatalf("failed to create grpc client: %v", err)
	}
	out := c.CanonicalTarget()
	_ = c.Close()

	return out
}

func BenchmarkExporterExportSpans(b *testing.B) {
	const n = 10

	run := func(b *testing.B) {
		mc := runMockCollectorWithConfig(b, &mockConfig{
			endpoint: "localhost:0",
			partial: &coltracepb.ExportTracePartialSuccess{
				RejectedSpans: 5,
				ErrorMessage:  "partially successful",
			},
		})
		b.Cleanup(func() { require.NoError(b, mc.stop()) })

		exp := newGRPCExporter(b, b.Context(), mc.endpoint)
		b.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(b, exp.Shutdown(context.Background()))
		})

		stubs := make([]tracetest.SpanStub, n)
		for i := range stubs {
			stubs[i].Name = fmt.Sprintf("Span %d", i)
		}
		spans := tracetest.SpanStubs(stubs).Snapshots()

		b.ReportAllocs()
		b.ResetTimer()

		var err error
		for b.Loop() {
			err = exp.ExportSpans(b.Context(), spans)
		}
		_ = err
	}

	b.Run("Observability", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		run(b)
	})

	b.Run("NoObservability", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
		run(b)
	})
}
