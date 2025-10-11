// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/counter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/observ"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/otest"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

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
		coll, err := otest.NewHTTPCollector("", rCh)
		require.NoError(t, err)

		addr := coll.Addr().String()
		opts := []Option{WithEndpoint(addr), WithInsecure()}
		cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)
		client, err := newClient(cfg)
		require.NoError(t, err)
		return clientShim{client}, coll
	}

	t.Run("Integration", otest.RunClientTests(factory))
}

func TestClientWithHTTPCollectorRespondingPlainText(t *testing.T) {
	ctx := t.Context()
	coll, err := otest.NewHTTPCollector("", nil, otest.WithHTTPCollectorRespondingPlainText())
	require.NoError(t, err)

	addr := coll.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)
	client, err := newClient(cfg)
	require.NoError(t, err)

	require.NoError(t, client.UploadMetrics(ctx, &mpb.ResourceMetrics{}))
	require.NoError(t, client.Shutdown(ctx))
	got := coll.Collect().Dump()
	require.Len(t, got, 1, "upload of one ResourceMetrics")
}

func TestNewWithInvalidEndpoint(t *testing.T) {
	ctx := t.Context()
	exp, err := New(ctx, WithEndpoint("host:invalid-port"))
	assert.Error(t, err)
	assert.Nil(t, exp)
}

func TestConfig(t *testing.T) {
	factoryFunc := func(ePt string, rCh <-chan otest.ExportResult, o ...Option) (metric.Exporter, *otest.HTTPCollector) {
		coll, err := otest.NewHTTPCollector(ePt, rCh)
		require.NoError(t, err)

		opts := []Option{WithEndpoint(coll.Addr().String())}
		if !strings.HasPrefix(strings.ToLower(ePt), "https") {
			opts = append(opts, WithInsecure())
		}
		opts = append(opts, o...)

		ctx := t.Context()
		exp, err := New(ctx, opts...)
		require.NoError(t, err)
		return exp, coll
	}

	t.Run("WithEndpointURL", func(t *testing.T) {
		coll, err := otest.NewHTTPCollector("", nil)
		require.NoError(t, err)
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.

		exp, err := New(ctx, WithEndpointURL("http://"+coll.Addr().String()))
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithHeaders", func(t *testing.T) {
		key := http.CanonicalHeaderKey("my-custom-header")
		headers := map[string]string{key: "custom-value"}
		exp, coll := factoryFunc("", nil, WithHeaders(headers))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Regexp(t, "OTel Go OTLP over HTTP/protobuf metrics exporter/[01]\\..*", got)
		require.Contains(t, got, key)
		assert.Equal(t, []string{headers[key]}, got[key])
	})

	t.Run("WithTimeout", func(t *testing.T) {
		// Do not send on rCh so the Collector never responds to the client.
		rCh := make(chan otest.ExportResult)
		exp, coll := factoryFunc(
			"",
			rCh,
			WithTimeout(time.Millisecond),
			WithRetry(RetryConfig{Enabled: false}),
		)
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, &metricdata.ResourceMetrics{})
		assert.ErrorAs(t, err, new(retryableError))
	})

	t.Run("WithCompressionGZip", func(t *testing.T) {
		exp, coll := factoryFunc("", nil, WithCompression(GzipCompression))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithRetry", func(t *testing.T) {
		emptyErr := errors.New("")
		rCh := make(chan otest.ExportResult, 5)
		header := http.Header{http.CanonicalHeaderKey("Retry-After"): {"10"}}
		// All retryable errors.
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusServiceUnavailable,
			Err:    emptyErr,
			Header: header,
		}}
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusTooManyRequests,
			Err:    emptyErr,
		}}
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusGatewayTimeout,
			Err:    emptyErr,
		}}
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusBadGateway,
			Err:    emptyErr,
		}}
		rCh <- otest.ExportResult{}
		exp, coll := factoryFunc("", rCh, WithRetry(RetryConfig{
			Enabled:         true,
			InitialInterval: time.Nanosecond,
			MaxInterval:     time.Millisecond,
			MaxElapsedTime:  time.Minute,
		}))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}), "failed retry")
		assert.Empty(t, rCh, "failed HTTP responses did not occur")
	})

	t.Run("WithRetryAndExporterErr", func(t *testing.T) {
		exporterErr := errors.New("rpc error: code = Unavailable desc = service.name not found in resource attributes")
		rCh := make(chan otest.ExportResult, 1)
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusTooManyRequests,
			Err:    exporterErr,
		}}
		exp, coll := factoryFunc("", rCh, WithRetry(RetryConfig{
			Enabled: false,
		}))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, &metricdata.ResourceMetrics{})
		assert.ErrorContains(t, err, exporterErr.Error())

		// To test the `Unwrap` and `As` function of retryable error
		var retryErr *retryableError
		assert.ErrorAs(t, err, &retryErr)
		assert.ErrorIs(t, err, *retryErr)
	})

	t.Run("WithURLPath", func(t *testing.T) {
		path := "/prefix/v2/metrics"
		ePt := fmt.Sprintf("http://localhost:0%s", path)
		exp, coll := factoryFunc(ePt, nil, WithURLPath(path))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithTLSClientConfig", func(t *testing.T) {
		ePt := "https://localhost:0"
		tlsCfg := &tls.Config{InsecureSkipVerify: true}
		exp, coll := factoryFunc(ePt, nil, WithTLSClientConfig(tlsCfg))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithCustomUserAgent", func(t *testing.T) {
		key := http.CanonicalHeaderKey("user-agent")
		headers := map[string]string{key: "custom-user-agent"}
		exp, coll := factoryFunc("", nil, WithHeaders(headers))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Contains(t, got, key)
		assert.Equal(t, []string{headers[key]}, got[key])
	})

	t.Run("WithProxy", func(t *testing.T) {
		headerKeySetInProxy := http.CanonicalHeaderKey("X-Using-Proxy")
		headerValueSetInProxy := "true"
		exp, coll := factoryFunc("", nil, WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					r.Header.Set(headerKeySetInProxy, headerValueSetInProxy)
					return r.URL, nil
				},
			},
		}))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Contains(t, got, headerKeySetInProxy)
		assert.Equal(t, []string{headerValueSetInProxy}, got[headerKeySetInProxy])
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		headerKeySetInProxy := http.CanonicalHeaderKey("X-Using-Proxy")
		headerValueSetInProxy := "true"
		exp, coll := factoryFunc("", nil, WithProxy(func(r *http.Request) (*url.URL, error) {
			r.Header.Set(headerKeySetInProxy, headerValueSetInProxy)
			return r.URL, nil
		}))
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		require.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		// Ensure everything is flushed.
		require.NoError(t, exp.Shutdown(ctx))

		got := coll.Headers()
		require.Contains(t, got, headerKeySetInProxy)
		assert.Equal(t, []string{headerValueSetInProxy}, got[headerKeySetInProxy])
	})

	t.Run("non-retryable errors are propagated", func(t *testing.T) {
		exporterErr := errors.New("missing required attribute aaa")
		rCh := make(chan otest.ExportResult, 1)
		rCh <- otest.ExportResult{Err: &otest.HTTPResponseError{
			Status: http.StatusBadRequest,
			Err:    exporterErr,
		}}
		exp, coll := factoryFunc("", rCh)
		ctx := context.Background() //nolint:usetesting // required to avoid getting a canceled context at cleanup.
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		exCtx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		err := exp.Export(exCtx, &metricdata.ResourceMetrics{})
		assert.ErrorContains(t, err, exporterErr.Error())

		assert.NoError(t, exCtx.Err())
	})
}

func TestClientInstrumentation(t *testing.T) {
	// Enable instrumentation for this test.
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Reset client ID to be deterministic.
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
	rCh := make(chan otest.ExportResult, 1)
	// Test partial success - return HTTP 200 with partial success info
	rCh <- otest.ExportResult{
		Response: &colmetricpb.ExportMetricsServiceResponse{
			PartialSuccess: &colmetricpb.ExportMetricsPartialSuccess{
				RejectedDataPoints: n,
				ErrorMessage:       msg,
			},
		},
	}

	coll, err := otest.NewHTTPCollector("", rCh)
	require.NoError(t, err)
	t.Cleanup(func() {
		//nolint:usetesting // required to avoid getting a canceled context at cleanup.
		require.NoError(t, coll.Shutdown(context.Background()))
	})
	t.Cleanup(func() { close(rCh) })

	addr := coll.Addr().String()
	opts := []Option{WithEndpoint(addr), WithInsecure()}
	ctx := t.Context()
	exp, err := New(ctx, opts...)
	require.NoError(t, err)

	// Export some test data
	err = exp.Export(ctx, &metricdata.ResourceMetrics{
		Resource: resource.NewWithAttributes(semconv.SchemaURL, attribute.String("service.name", "test")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "test"},
				Metrics: []metricdata.Metrics{
					{
						Name: "test-metric",
						Data: metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{Value: 42}}},
					},
				},
			},
		},
	})

	// Should get partial success error
	wantErr := internal.MetricPartialSuccessError(n, msg)
	require.ErrorIs(t, err, wantErr, "Expected partial success error")

	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
	require.NoError(t, exp.Shutdown(context.Background()))
	var got metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &got))

	attrs := observ.BaseAttrs(id, addr)

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      observ.ScopeName,
			Version:   observ.Version,
			SchemaURL: observ.SchemaURL,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...)},
					},
					Temporality: metricdata.CumulativeTemporality,
				},
			},
			{
				Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...)},
						{Attributes: attribute.NewSet(append(
							attrs,
							otelconv.SDKExporterMetricDataPointExported{}.AttrErrorType("*errors.joinError"),
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
							otelconv.SDKExporterOperationDuration{}.AttrHTTPResponseStatusCode(200),
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

func BenchmarkExporterExportMetrics(b *testing.B) {
	const n = 10

	run := func(b *testing.B) {
		coll, err := otest.NewHTTPCollector("", nil)
		require.NoError(b, err)
		b.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			require.NoError(b, coll.Shutdown(context.Background()))
		})

		opts := []Option{WithEndpoint(coll.Addr().String()), WithInsecure()}
		ctx := b.Context()
		exp, err := New(ctx, opts...)
		require.NoError(b, err)
		b.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(b, exp.Shutdown(context.Background()))
		})

		// Generate realistic test metric data with multiple metrics.
		now := time.Now()
		rm := &metricdata.ResourceMetrics{
			ScopeMetrics: []metricdata.ScopeMetrics{
				{
					Scope: instrumentation.Scope{
						Name:    "test",
						Version: "v1.0.0",
					},
					Metrics: make([]metricdata.Metrics, n),
				},
			},
		}

		for i := range rm.ScopeMetrics[0].Metrics {
			rm.ScopeMetrics[0].Metrics[i] = metricdata.Metrics{
				Name:        fmt.Sprintf("test_counter_%d", i),
				Description: fmt.Sprintf("A test counter %d", i),
				Unit:        "1",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								attribute.String("test", "value"),
								attribute.Int("counter", i),
							),
							StartTime: now,
							Time:      now,
							Value:     int64(i * 10),
						},
					},
				},
			}
		}

		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			err = exp.Export(b.Context(), rm)
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
