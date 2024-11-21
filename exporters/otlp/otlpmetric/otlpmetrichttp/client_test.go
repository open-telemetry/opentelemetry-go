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

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/otest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
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
	ctx := context.Background()
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
	ctx := context.Background()
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

		ctx := context.Background()
		exp, err := New(ctx, opts...)
		require.NoError(t, err)
		return exp, coll
	}

	t.Run("WithEndpointURL", func(t *testing.T) {
		coll, err := otest.NewHTTPCollector("", nil)
		require.NoError(t, err)
		ctx := context.Background()

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
		ctx := context.Background()
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
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		// Push this after Shutdown so the HTTP server doesn't hang.
		t.Cleanup(func() { close(rCh) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		err := exp.Export(ctx, &metricdata.ResourceMetrics{})
		assert.ErrorAs(t, err, new(retryableError))
	})

	t.Run("WithCompressionGZip", func(t *testing.T) {
		exp, coll := factoryFunc("", nil, WithCompression(GzipCompression))
		ctx := context.Background()
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
		ctx := context.Background()
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
		ctx := context.Background()
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
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithTLSClientConfig", func(t *testing.T) {
		ePt := "https://localhost:0"
		tlsCfg := &tls.Config{InsecureSkipVerify: true}
		exp, coll := factoryFunc(ePt, nil, WithTLSClientConfig(tlsCfg))
		ctx := context.Background()
		t.Cleanup(func() { require.NoError(t, coll.Shutdown(ctx)) })
		t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
		assert.NoError(t, exp.Export(ctx, &metricdata.ResourceMetrics{}))
		assert.Len(t, coll.Collect().Dump(), 1)
	})

	t.Run("WithCustomUserAgent", func(t *testing.T) {
		key := http.CanonicalHeaderKey("user-agent")
		headers := map[string]string{key: "custom-user-agent"}
		exp, coll := factoryFunc("", nil, WithHeaders(headers))
		ctx := context.Background()
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
		exp, coll := factoryFunc("", nil, WithProxy(func(r *http.Request) (*url.URL, error) {
			r.Header.Set(headerKeySetInProxy, headerValueSetInProxy)
			return r.URL, nil
		}))
		ctx := context.Background()
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
		ctx := context.Background()
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
