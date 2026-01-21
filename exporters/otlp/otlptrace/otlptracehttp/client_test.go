// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/counter"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/observ"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlptracetest"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

const (
	relOtherTracesPath = "post/traces/here"
	otherTracesPath    = "/post/traces/here"
)

var (
	testHeaders = map[string]string{
		"Otel-Go-Key-1": "somevalue",
		"Otel-Go-Key-2": "someothervalue",
	}

	customUserAgentHeader = map[string]string{
		"user-agent": "custom-user-agent",
	}

	customProxyHeader = map[string]string{
		"header-added-via-proxy": "proxy-value",
	}
)

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name            string
		opts            []otlptracehttp.Option
		mcCfg           mockCollectorConfig
		tls             bool
		withURLEndpoint bool
	}{
		{
			name: "no extra options",
			opts: nil,
		},
		{
			name:            "with URL endpoint",
			withURLEndpoint: true,
		},
		{
			name: "with gzip compression",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			},
		},
		{
			name: "retry",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{503, 429},
			},
		},
		{
			name: "retry with gzip compression",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{503, 502},
			},
		},
		{
			name: "retry with throttle",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{504},
				InjectResponseHeader: []map[string]string{
					{"Retry-After": "10"},
				},
			},
		},
		{
			name: "with empty paths (forced to defaults)",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithURLPath(""),
			},
		},
		{
			name: "with relative paths",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithURLPath(relOtherTracesPath),
			},
			mcCfg: mockCollectorConfig{
				TracesURLPath: otherTracesPath,
			},
		},
		{
			name: "with TLS",
			opts: nil,
			mcCfg: mockCollectorConfig{
				WithTLS: true,
			},
			tls: true,
		},
		{
			name: "with extra headers",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithHeaders(testHeaders),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: testHeaders,
			},
		},
		{
			name: "with custom user agent",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithHeaders(customUserAgentHeader),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: customUserAgentHeader,
			},
		},
		{
			name: "with custom proxy",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithProxy(func(r *http.Request) (*url.URL, error) {
					for k, v := range customProxyHeader {
						r.Header.Set(k, v)
					}
					return r.URL, nil
				}),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: customProxyHeader,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mc := runMockCollector(t, tc.mcCfg)
			defer mc.MustStop(t)
			allOpts := []otlptracehttp.Option{}

			if tc.withURLEndpoint {
				allOpts = append(allOpts, otlptracehttp.WithEndpointURL("http://"+mc.Endpoint()))
			} else {
				allOpts = append(allOpts, otlptracehttp.WithEndpoint(mc.Endpoint()))
			}
			if tc.tls {
				tlsConfig := mc.ClientTLSConfig()
				require.NotNil(t, tlsConfig)
				allOpts = append(allOpts, otlptracehttp.WithTLSClientConfig(tlsConfig))
			} else {
				allOpts = append(allOpts, otlptracehttp.WithInsecure())
			}
			allOpts = append(allOpts, tc.opts...)
			client := otlptracehttp.NewClient(allOpts...)
			ctx := t.Context()
			exporter, err := otlptrace.New(ctx, client)
			if assert.NoError(t, err) {
				defer func() {
					assert.NoError(t, exporter.Shutdown(ctx))
				}()
				otlptracetest.RunEndToEndTest(ctx, t, exporter, mc)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollector(t, mockCollectorConfig{})
	defer func() {
		_ = mc.Stop()
	}()

	<-time.After(5 * time.Millisecond)

	otlptracetest.RunExporterShutdownTest(t, func() otlptrace.Client {
		return otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(mc.endpoint),
		)
	})
}

func TestTimeout(t *testing.T) {
	delay := make(chan struct{})
	mcCfg := mockCollectorConfig{Delay: delay}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	defer func() { close(delay) }()
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(time.Nanosecond),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{Enabled: false}),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, client)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.ErrorContains(t, err, "Client.Timeout exceeded while awaiting headers")
}

func TestNoRetry(t *testing.T) {
	mc := runMockCollector(t, mockCollectorConfig{
		InjectHTTPStatus: []int{http.StatusBadRequest},
		Partial: &coltracepb.ExportTracePartialSuccess{
			ErrorMessage: "missing required attribute aaa",
		},
	})
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Nanosecond,
			MaxInterval:     1 * time.Nanosecond,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "))

	msg := fmt.Sprintf("failed to send to http://%s/v1/traces: 400 Bad Request", mc.endpoint)
	assert.ErrorContains(t, err, msg)

	msg = "missing required attribute aaa"
	assert.ErrorContains(t, err, msg)

	assert.Empty(t, mc.GetSpans())
}

func TestEmptyData(t *testing.T) {
	mcCfg := mockCollectorConfig{}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	assert.NoError(t, err)
	err = exporter.ExportSpans(ctx, nil)
	assert.NoError(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestCancelledContext(t *testing.T) {
	mcCfg := mockCollectorConfig{}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx, cancel := context.WithCancel(t.Context())
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(t.Context()))
	}()
	cancel()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestDeadlineContext(t *testing.T) {
	statuses := make([]int, 0, 5)
	for i := 0; i < cap(statuses); i++ {
		statuses = append(statuses, http.StatusTooManyRequests)
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Hour,
			MaxInterval:     1 * time.Hour,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(t.Context()))
	}()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestStopWhileExportingConcurrentSafe(t *testing.T) {
	statuses := make([]int, 0, 5)
	for i := 0; i < cap(statuses); i++ {
		statuses = append(statuses, http.StatusTooManyRequests)
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Hour,
			MaxInterval:     1 * time.Hour,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	doneCh := make(chan struct{})
	go func() {
		err := exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
		assert.Error(t, err)
		assert.Empty(t, mc.GetSpans())
		close(doneCh)
	}()
	<-time.After(time.Second)
	err = exporter.Shutdown(ctx)
	assert.NoError(t, err)
	<-doneCh
}

func TestPartialSuccess(t *testing.T) {
	mcCfg := mockCollectorConfig{
		Partial: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 2,
			ErrorMessage:  "partially successful",
		},
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(t.Context()))
	}()

	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	want := internal.TracePartialSuccessError(0, "")
	assert.ErrorIs(t, err, want)
}

func TestOtherHTTPSuccess(t *testing.T) {
	for code := 201; code <= 299; code++ {
		t.Run(fmt.Sprintf("status_%d", code), func(t *testing.T) {
			mcCfg := mockCollectorConfig{
				InjectHTTPStatus: []int{code},
			}
			mc := runMockCollector(t, mcCfg)
			defer mc.MustStop(t)
			driver := otlptracehttp.NewClient(
				otlptracehttp.WithEndpoint(mc.Endpoint()),
				otlptracehttp.WithInsecure(),
			)
			ctx := t.Context()
			exporter, err := otlptrace.New(ctx, driver)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, exporter.Shutdown(t.Context()))
			}()

			errs := []error{}
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
				errs = append(errs, err)
			}))
			err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
			assert.NoError(t, err)

			assert.Empty(t, errs)
		})
	}
}

func TestCollectorRespondingNonProtobufContent(t *testing.T) {
	mcCfg := mockCollectorConfig{
		InjectContentType: "application/octet-stream",
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := t.Context()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(t.Context()))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.NoError(t, err)
	assert.Len(t, mc.GetSpans(), 1)
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
	mc := runMockCollector(t, mockCollectorConfig{
		InjectHTTPStatus: []int{400},
		Partial: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: n,
			ErrorMessage:  msg,
		},
	})
	t.Cleanup(func() { require.NoError(t, mc.Stop()) })

	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(t.Context(), driver)
	require.NoError(t, err)

	err = exporter.ExportSpans(t.Context(), otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)

	require.NoError(t, exporter.Shutdown(t.Context()))
	var got metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &got))

	attrs := observ.BaseAttrs(id, mc.endpoint)

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
							otelconv.SDKExporterOperationDuration{}.AttrHTTPResponseStatusCode(400),
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

func BenchmarkExporterExportSpans(b *testing.B) {
	const n = 10

	run := func(b *testing.B) {
		mc := runMockCollector(b, mockCollectorConfig{
			Partial: &coltracepb.ExportTracePartialSuccess{
				RejectedSpans: 5,
				ErrorMessage:  "partially successful",
			},
		})
		b.Cleanup(func() { require.NoError(b, mc.Stop()) })

		c := otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(mc.Endpoint()),
			otlptracehttp.WithInsecure(),
		)
		exp, err := otlptrace.New(b.Context(), c)
		require.NoError(b, err)
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
