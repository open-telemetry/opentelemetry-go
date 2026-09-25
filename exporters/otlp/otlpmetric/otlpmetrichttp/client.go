// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/counter"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/observ"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/retry"
)

type client struct {
	// req is cloned for every upload the client makes.
	req             *http.Request
	compression     Compression
	maxRequestSize  int
	maxResponseSize int64
	requestFunc     retry.RequestFunc
	httpClient      *http.Client

	inst *observ.Instrumentation
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

var errInsecureEndpointWithTLS = errors.New("insecure HTTP endpoint cannot use TLS client configuration")

// newClient creates a new HTTP metric client.
func newClient(cfg oconf.Config) (*client, error) {
	if cfg.Metrics.Insecure && cfg.Metrics.TLSCfg != nil {
		return nil, errInsecureEndpointWithTLS
	}

	if Compression(cfg.Metrics.Compression) == ZstdCompression && !zstdSupported {
		// Fail fast at construction rather than at first export: this
		// build was compiled with the nozstd tag, so the zstd codec is
		// unavailable.
		return nil, errors.New(
			"otlpmetrichttp: zstd compression configured but not available (built with the nozstd tag?)",
		)
	}

	httpClient := cfg.Metrics.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: ourTransport,
			Timeout:   cfg.Metrics.Timeout,
		}

		if cfg.Metrics.TLSCfg != nil || cfg.Metrics.Proxy != nil {
			clonedTransport := ourTransport.Clone()
			httpClient.Transport = clonedTransport

			if cfg.Metrics.TLSCfg != nil {
				clonedTransport.TLSClientConfig = cfg.Metrics.TLSCfg
			}
			if cfg.Metrics.Proxy != nil {
				clonedTransport.Proxy = cfg.Metrics.Proxy
			}
		}
	}

	u := &url.URL{
		Scheme: "https",
		Host:   cfg.Metrics.Endpoint,
		Path:   cfg.Metrics.URLPath,
	}
	if cfg.Metrics.Insecure {
		u.Scheme = "http"
	}
	// Body is set when this is cloned during upload.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	userAgent := "OTel Go OTLP over HTTP/protobuf metrics exporter/" + Version()
	req.Header.Set("User-Agent", userAgent)

	if n := len(cfg.Metrics.Headers); n > 0 {
		for k, v := range cfg.Metrics.Headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-protobuf")

	// Initialize the instrumentation.
	inst, err := observ.NewInstrumentation(counter.NextExporterID(), cfg.Metrics.Endpoint)

	return &client{
		compression:     Compression(cfg.Metrics.Compression),
		maxRequestSize:  cfg.Metrics.MaxRequestSize,
		maxResponseSize: cfg.Metrics.MaxResponseSize,
		req:             req,
		requestFunc:     cfg.RetryConfig.RequestFunc(evaluate),
		httpClient:      httpClient,
		inst:            inst,
	}, err
}

// Shutdown shuts down the client, freeing all resources.
func (c *client) Shutdown(ctx context.Context) error {
	// The otlpmetric.Exporter synchronizes access to client methods and
	// ensures this is called only once. The only thing that needs to be done
	// here is to release any computational resources the client holds.

	c.requestFunc = nil
	c.httpClient = nil
	return ctx.Err()
}

// UploadMetrics sends protoMetrics to the connected endpoint.
//
// Retryable errors from the server will be handled according to any
// RetryConfig the client was created with.
func (c *client) UploadMetrics(ctx context.Context, protoMetrics *metricpb.ResourceMetrics) (uploadErr error) {
	// The otlpmetric.Exporter synchronizes access to client methods, and
	// ensures this is not called after the Exporter is shutdown. Only thing
	// to do here is send data.

	pbRequest := &colmetricpb.ExportMetricsServiceRequest{
		ResourceMetrics: []*metricpb.ResourceMetrics{protoMetrics},
	}
	body, err := proto.Marshal(pbRequest)
	if err != nil {
		return err
	}
	if maxSize := c.maxRequestSize; maxSize > 0 && len(body) > maxSize {
		return fmt.Errorf("request body too large: exceeded %d bytes", maxSize)
	}
	request, err := c.newRequest(ctx, body)
	if err != nil {
		return err
	}

	var statusCode int
	if c.inst != nil {
		op := c.inst.ExportMetrics(ctx, protoMetrics)
		defer func() { op.End(uploadErr, statusCode) }()
	}

	return errors.Join(uploadErr, c.requestFunc(ctx, func(iCtx context.Context) error {
		select {
		case <-iCtx.Done():
			return iCtx.Err()
		default:
		}

		var partialErr error
		var err error
		statusCode, err, partialErr = c.doHTTPAttempt(iCtx, &request)
		uploadErr = errors.Join(uploadErr, partialErr)
		if (statusCode == http.StatusUnsupportedMediaType || statusCode == http.StatusBadRequest) &&
			c.compression == ZstdCompression {
			// A 415/400 for a zstd-compressed export is non-retryable (see
			// the zstd OTEP), so without this the batch would be dropped
			// outright on any receiver that doesn't support zstd. Retry
			// this same batch once with gzip, which every conformant OTLP
			// receiver is required to support.
			otel.Handle(fmt.Errorf(
				"otlpmetrichttp: receiver does not support zstd compression (status %d), retrying with gzip: %w",
				statusCode, err,
			))
			if gzErr := request.recompressWithGzip(); gzErr != nil {
				return errors.Join(err, gzErr)
			}
			statusCode, err, partialErr = c.doHTTPAttempt(iCtx, &request)
			uploadErr = errors.Join(uploadErr, partialErr)
		}
		return err
	}))
}

// doHTTPAttempt sends request once and interprets the response. It returns
// the response status code, an error if the attempt failed (nil on
// success), and any partial-success error, which is returned separately
// since it can be non-nil even when the attempt itself succeeded.
func (c *client) doHTTPAttempt(ctx context.Context, request *request) (statusCode int, err, partialErr error) {
	request.reset(ctx)
	// nolint:gosec // URL is constructed from validated OTLP endpoint configuration
	resp, doErr := c.httpClient.Do(request.Request)
	var urlErr *url.Error
	if errors.As(doErr, &urlErr) && urlErr.Temporary() {
		return 0, newResponseError(http.Header{}, doErr), nil
	}
	if doErr != nil {
		return 0, doErr, nil
	}
	if resp != nil {
		statusCode = resp.StatusCode
		if resp.Body != nil {
			defer func() {
				if closeErr := resp.Body.Close(); closeErr != nil {
					err = errors.Join(err, closeErr)
				}
			}()
		}
	}

	if statusCode >= 200 && statusCode <= 299 {
		// Success, do not retry.

		// Read the partial success message, if any.
		var respData bytes.Buffer
		if err := internal.CopyResponseBody(&respData, resp.Body, c.maxResponseSize); err != nil {
			return statusCode, err, nil
		}
		if respData.Len() == 0 {
			return statusCode, nil, nil
		}

		if resp.Header.Get("Content-Type") == "application/x-protobuf" {
			var respProto colmetricpb.ExportMetricsServiceResponse
			if err := proto.Unmarshal(respData.Bytes(), &respProto); err != nil {
				return statusCode, err, nil
			}

			if respProto.PartialSuccess != nil {
				msg := respProto.PartialSuccess.GetErrorMessage()
				n := respProto.PartialSuccess.GetRejectedDataPoints()
				if n != 0 || msg != "" {
					partialErr = internal.MetricPartialSuccessError(n, msg)
				}
			}
		}
		return statusCode, nil, partialErr
	}
	// Error cases.

	// server may return a message with the response
	// body, so we read it to include in the error
	// message to be returned. It will help in
	// debugging the actual issue.
	var respData bytes.Buffer
	if err := internal.CopyResponseBody(&respData, resp.Body, c.maxResponseSize); err != nil {
		return statusCode, err, nil
	}
	respStr := strings.TrimSpace(respData.String())
	if respStr == "" {
		respStr = "(empty)"
	}
	bodyErr := fmt.Errorf("body: %s", respStr)

	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		// Retryable failure.
		return statusCode, newResponseError(resp.Header, bodyErr), nil
	default:
		// Non-retryable failure.
		return statusCode, fmt.Errorf("failed to send metrics to %s: %s (%w)", request.URL, resp.Status, bodyErr), nil
	}
}

var gzPool = sync.Pool{
	New: func() any {
		w := gzip.NewWriter(io.Discard)
		return w
	},
}

func (c *client) newRequest(ctx context.Context, body []byte) (request, error) {
	r := c.req.Clone(ctx)
	req := request{Request: r, rawBody: body}

	switch c.compression {
	case NoCompression:
		r.ContentLength = int64(len(body))
		req.bodyReader = bodyReader(body)
		req.GetBody = bodyReaderErr(body)
	case GzipCompression:
		// Ensure the content length is not used.
		r.ContentLength = -1
		r.Header.Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer func() {
			gz.Reset(io.Discard)
			gzPool.Put(gz)
		}()

		var b bytes.Buffer
		gz.Reset(&b)

		if _, err := gz.Write(body); err != nil {
			return req, err
		}
		// Close needs to be called to ensure body is fully written.
		if err := gz.Close(); err != nil {
			return req, err
		}

		req.bodyReader = bodyReader(b.Bytes())
		req.GetBody = bodyReaderErr(b.Bytes())
	case ZstdCompression:
		if !zstdSupported {
			return req, errors.New(
				"otlpmetrichttp: zstd compression configured but not available (built with the nozstd tag?)",
			)
		}
		// Ensure the content length is not used.
		r.ContentLength = -1
		r.Header.Set("Content-Encoding", "zstd")

		compressed, err := compressZstd(body)
		if err != nil {
			return req, err
		}

		req.bodyReader = bodyReader(compressed)
		req.GetBody = bodyReaderErr(compressed)
	default:
		return req, fmt.Errorf("otlpmetrichttp: unsupported compression: %v", c.compression)
	}

	return req, nil
}

// bodyReader returns a closure returning a new reader for buf.
func bodyReader(buf []byte) func() io.ReadCloser {
	return func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(buf))
	}
}

// bodyReaderErr returns a closure returning a new reader for buf.
func bodyReaderErr(buf []byte) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(buf)), nil
	}
}

// request wraps an http.Request with a resettable body reader.
type request struct {
	*http.Request

	// bodyReader allows the same body to be used for multiple requests.
	bodyReader func() io.ReadCloser

	// rawBody is the uncompressed request body, retained so a rejected
	// zstd export can be recompressed with gzip and retried.
	rawBody []byte
}

// reset reinitializes the request Body and uses ctx for the request.
func (r *request) reset(ctx context.Context) {
	r.Body = r.bodyReader()
	r.Request = r.WithContext(ctx)
}

// recompressWithGzip replaces the request's body with a gzip-compressed
// copy of rawBody. It's used to fall back from zstd on a receiver that
// rejects it, since every conformant OTLP receiver is required to support
// gzip.
func (r *request) recompressWithGzip() error {
	gz := gzPool.Get().(*gzip.Writer)
	defer func() {
		gz.Reset(io.Discard)
		gzPool.Put(gz)
	}()

	var b bytes.Buffer
	gz.Reset(&b)
	if _, err := gz.Write(r.rawBody); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	r.Header.Set("Content-Encoding", "gzip")
	r.ContentLength = -1
	r.bodyReader = bodyReader(b.Bytes())
	r.GetBody = bodyReaderErr(b.Bytes())
	return nil
}

// retryableError represents a request failure that can be retried.
type retryableError struct {
	throttle time.Duration
	err      error
}

// newResponseError returns a retryableError and will extract any explicit
// throttle delay contained in headers. The returned error wraps wrapped
// if it is not nil.
func newResponseError(header http.Header, wrapped error) error {
	var rErr retryableError
	if v := header.Get("Retry-After"); v != "" {
		rErr.throttle = retryAfterDuration(v)
	}

	rErr.err = wrapped
	return rErr
}

func retryAfterDuration(v string) time.Duration {
	if t, err := strconv.ParseInt(v, 10, 64); err == nil && t >= 0 {
		const maxRetryAfterSeconds = int64(1<<63-1) / int64(time.Second)
		if t > maxRetryAfterSeconds {
			return time.Duration(1<<63 - 1)
		}
		return time.Duration(t) * time.Second
	}

	if date, err := http.ParseTime(v); err == nil {
		return max(time.Until(date), 0)
	}

	return 0
}

func (e retryableError) Error() string {
	if e.err != nil {
		return "retry-able request failure: " + e.err.Error()
	}

	return "retry-able request failure"
}

func (e retryableError) Unwrap() error {
	return e.err
}

func (e retryableError) As(target any) bool {
	if e.err == nil {
		return false
	}

	switch v := target.(type) {
	case **retryableError:
		*v = &e
		return true
	default:
		return false
	}
}

// evaluate returns if err is retry-able. If it is and it includes an explicit
// throttling delay, that delay is also returned.
func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	// Do not use errors.As here, this should only be flattened one layer. If
	// there are several chained errors, all the errors above it will be
	// discarded if errors.As is used instead.
	rErr, ok := err.(retryableError) //nolint:errorlint
	if !ok {
		return false, 0
	}

	return true, rErr.throttle
}
