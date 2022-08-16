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

//go:build go1.18
// +build go1.18

package otlpmetrichttp // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/internal/retry"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/oconf"
	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

const contentTypeProto = "application/x-protobuf"

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

// Keep it in sync with golang's DefaultTransport from net/http! We
// have our own copy to avoid handling a situation where the
// DefaultTransport is overwritten with some different implementation
// of http.RoundTripper or it's modified by another package.
var ourTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type client struct {
	url         string
	compression Compression
	headers     map[string]string
	requestFunc retry.RequestFunc
	client      *http.Client
}

// NewClient creates a new HTTP metric client.
func NewClient(ctx context.Context, opts ...Option) (otlpmetric.Client, error) {
	cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)

	httpClient := &http.Client{
		Transport: ourTransport,
		Timeout:   cfg.Metrics.Timeout,
	}
	if cfg.Metrics.TLSCfg != nil {
		transport := ourTransport.Clone()
		transport.TLSClientConfig = cfg.Metrics.TLSCfg
		httpClient.Transport = transport
	}

	format := "https://%s%s"
	if cfg.Metrics.Insecure {
		format = "http://%s%s"
	}
	rawURL := fmt.Sprintf(format, cfg.Metrics.Endpoint, cfg.Metrics.URLPath)
	// Ensure target URL is valid at start, instead of every export call.
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Make a copy of headers so the underlying value does not change.
	var h map[string]string
	if n := len(cfg.Metrics.Headers); n > 0 {
		h = make(map[string]string, n)
		for k, v := range cfg.Metrics.Headers {
			h[k] = v
		}
	}

	return &client{
		url:         u.String(),
		compression: Compression(cfg.Metrics.Compression),
		headers:     h,
		requestFunc: cfg.RetryConfig.RequestFunc(evaluate),
		client:      httpClient,
	}, nil
}

// ForceFlush does nothing, the client holds no state.
func (c *client) ForceFlush(ctx context.Context) error { return ctx.Err() }

// Shutdown shuts down the client, freeing all resources.
func (c *client) Shutdown(ctx context.Context) error {
	// The otlpmetric.Exporter synchronizes access to client methods and
	// ensures this is called only once. The only thing that needs to be done
	// here is to release any computational resources the client holds.

	c.requestFunc = nil
	c.client = nil
	return ctx.Err()
}

// UploadMetrics sends protoMetrics to the connected endpoint.
//
// Retryable errors from the server will be handled according to any
// RetryConfig the client was created with.
func (c *client) UploadMetrics(ctx context.Context, protoMetrics *metricpb.ResourceMetrics) error {
	// The otlpmetric.Exporter synchronizes access to client methods, and
	// ensures this is not called after the Exporter is shutdown. Only thing
	// to do here is send data.

	pbRequest := &colmetricpb.ExportMetricsServiceRequest{
		ResourceMetrics: []*metricpb.ResourceMetrics{protoMetrics},
	}
	rawRequest, err := proto.Marshal(pbRequest)
	if err != nil {
		return err
	}

	request, err := c.newRequest(rawRequest)
	if err != nil {
		return err
	}

	return c.requestFunc(ctx, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request.reset(ctx)
		resp, err := c.client.Do(request.Request)
		if err != nil {
			return err
		}

		var rErr error
		switch resp.StatusCode {
		case http.StatusOK:
			// Success, do not retry.
		case http.StatusTooManyRequests,
			http.StatusServiceUnavailable:
			// Retry-able failure.
			rErr = newResponseError(resp.Header)

			// Going to retry, drain the body to reuse the connection.
			if _, err := io.Copy(io.Discard, resp.Body); err != nil {
				_ = resp.Body.Close()
				return err
			}
		default:
			rErr = fmt.Errorf("failed to send metrics to %s: %s", request.URL, resp.Status)
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}
		return rErr
	})
}

func (c *client) newRequest(body []byte) (request, error) {
	r, err := http.NewRequest(http.MethodPost, c.url, nil)
	if err != nil {
		return request{Request: r}, err
	}

	for k, v := range c.headers {
		r.Header.Set(k, v)
	}
	r.Header.Set("Content-Type", contentTypeProto)

	req := request{Request: r}
	switch c.compression {
	case NoCompression:
		r.ContentLength = (int64)(len(body))
		req.bodyReader = bodyReader(body)
	case GzipCompression:
		// Ensure the content length is not used.
		r.ContentLength = -1
		r.Header.Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		var b bytes.Buffer
		gz.Reset(&b)

		if _, err := gz.Write(body); err != nil {
			return req, err
		}
		// Close needs to be called to ensure body if fully written.
		if err := gz.Close(); err != nil {
			return req, err
		}

		req.bodyReader = bodyReader(b.Bytes())
	}

	return req, nil
}

// bodyReader returns a closure returning a new reader for buf.
func bodyReader(buf []byte) func() io.ReadCloser {
	return func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(buf))
	}
}

// request wraps an http.Request with a resettable body reader.
type request struct {
	*http.Request

	// bodyReader allows the same body to be used for multiple requests.
	bodyReader func() io.ReadCloser
}

// reset reinitializes the request Body and uses ctx for the request.
func (r *request) reset(ctx context.Context) {
	r.Body = r.bodyReader()
	r.Request = r.Request.WithContext(ctx)
}

// retryableError represents a request failure that can be retried.
type retryableError struct {
	throttle int64
}

// newResponseError returns a retryableError and will extract any explicit
// throttle delay contained in headers.
func newResponseError(header http.Header) error {
	var rErr retryableError
	if s, ok := header["Retry-After"]; ok {
		if t, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			rErr.throttle = t
		}
	}
	return rErr
}

func (e retryableError) Error() string {
	return "retry-able request failure"
}

// evaluate returns if err is retry-able. If it is and it includes an explicit
// throttling delay, that delay is also returned.
func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	rErr, ok := err.(retryableError)
	if !ok {
		return false, 0
	}

	return true, time.Duration(rErr.throttle)
}
