// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	colprofilepb "go.opentelemetry.io/proto/slim/otlp/collector/profiles/v1development"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal/retry"
)

const (
	// profilesDevelopmentVersionHeader identifies the OTLP Profiles v1development
	// schema version per opentelemetry-proto/docs/specification.md.
	profilesDevelopmentVersionHeader = "OTLP-Profiles-Development-Version"
	profilesDevelopmentVersion       = "1"
)

type client struct {
	uploadProfiles func(context.Context, *colprofilepb.ExportProfilesServiceRequest) error
	stop           func()
}

func (c *client) UploadProfiles(ctx context.Context, req *colprofilepb.ExportProfilesServiceRequest) error {
	if c.uploadProfiles != nil {
		return c.uploadProfiles(ctx, req)
	}
	return nil
}

func newNoopClient() *client {
	return &client{}
}

var errInsecureEndpointWithTLS = errors.New("insecure HTTP endpoint cannot use TLS client configuration")

func newHTTPClient(cfg config) (*client, error) {
	if cfg.insecure.Value && cfg.tlsCfg.Value != nil {
		return nil, errInsecureEndpointWithTLS
	}

	hc := cfg.httpClient
	var exportTimeout time.Duration
	var closeIdle func()
	if hc == nil {
		exportTimeout = cfg.timeout.Value
		hc = &http.Client{
			Transport: ourTransport,
			Timeout:   exportTimeout,
		}

		if cfg.tlsCfg.Value != nil || cfg.proxy.Set {
			clonedTransport := ourTransport.Clone()
			hc.Transport = clonedTransport
			closeIdle = clonedTransport.CloseIdleConnections

			if cfg.tlsCfg.Value != nil {
				clonedTransport.TLSClientConfig = cfg.tlsCfg.Value
			}
			if cfg.proxy.Set {
				clonedTransport.Proxy = cfg.proxy.Value
			}
		}
	}

	u := &url.URL{
		Scheme:  "https",
		Host:    cfg.endpoint.Value,
		Path:    cfg.path.Value.Path,
		RawPath: cfg.path.Value.RawPath,
	}
	if cfg.insecure.Value {
		u.Scheme = "http"
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "OTel Go OTLP over HTTP/protobuf profiles exporter/"+Version())
	for k, v := range cfg.headers.Value {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set(profilesDevelopmentVersionHeader, profilesDevelopmentVersion)

	c := &httpClient{
		compression:     cfg.compression.Value,
		maxRequestSize:  cfg.maxRequestSize.Value,
		maxResponseSize: cfg.maxResponseSize.Value,
		req:             req,
		requestFunc:     cfg.retryCfg.Value.RequestFunc(evaluate),
		client:          hc,
		exportTimeout:   exportTimeout,
	}
	return &client{uploadProfiles: c.uploadProfiles, stop: closeIdle}, nil
}

type httpClient struct {
	req             *http.Request
	compression     Compression
	maxRequestSize  int
	maxResponseSize int64
	requestFunc     retry.RequestFunc
	client          *http.Client
	exportTimeout   time.Duration
}

// Local copy of net/http.DefaultTransport so caller mutations to
// http.DefaultTransport do not affect exporter transport settings.
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

func (c *httpClient) uploadProfiles(
	ctx context.Context,
	pbRequest *colprofilepb.ExportProfilesServiceRequest,
) (uploadErr error) {
	if maxSize := c.maxRequestSize; maxSize > 0 && proto.Size(pbRequest) > maxSize {
		return fmt.Errorf("request body too large: exceeded %d bytes", maxSize)
	}

	body, err := proto.Marshal(pbRequest)
	if err != nil {
		return err
	}
	request, err := c.newRequest(ctx, body)
	if err != nil {
		return err
	}

	if c.exportTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeoutCause(ctx, c.exportTimeout, errors.New("exporter export timeout"))
		defer cancel()
	}

	return errors.Join(uploadErr, c.requestFunc(ctx, func(iCtx context.Context) error {
		select {
		case <-iCtx.Done():
			return iCtx.Err()
		default:
		}

		request.reset(iCtx)
		//nolint:gosec // URL is constructed from validated OTLP endpoint configuration.
		resp, err := c.client.Do(request.Request)
		var urlErr *url.Error
		if errors.As(err, &urlErr) && urlErr.Temporary() {
			return newResponseError(http.Header{}, err)
		}
		if err != nil {
			return err
		}
		if resp != nil && resp.Body != nil {
			defer func() {
				if closeErr := resp.Body.Close(); closeErr != nil {
					uploadErr = errors.Join(uploadErr, closeErr)
				}
			}()
		}

		respData, err := c.readResponse(resp)
		if err != nil {
			return err
		}
		if sc := resp.StatusCode; sc >= 200 && sc <= 299 {
			if len(respData) == 0 {
				return nil
			}
			if isProtobufContentType(resp.Header.Get("Content-Type")) {
				var respProto colprofilepb.ExportProfilesServiceResponse
				if err := proto.Unmarshal(respData, &respProto); err != nil {
					return err
				}
				if respProto.PartialSuccess != nil {
					msg := respProto.PartialSuccess.GetErrorMessage()
					n := respProto.PartialSuccess.GetRejectedProfiles()
					if n != 0 || msg != "" {
						uploadErr = errors.Join(uploadErr, internal.ProfilePartialSuccessError(n, msg))
					}
				}
			}
			return nil
		}

		respStr := strings.TrimSpace(string(respData))
		if respStr == "" {
			respStr = "(empty)"
		}
		bodyErr := fmt.Errorf("body: %s", respStr)

		switch resp.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return newResponseError(resp.Header, bodyErr)
		default:
			return fmt.Errorf("failed to send profiles to %s: %s (%w)", request.URL, resp.Status, bodyErr)
		}
	}))
}

func (c *httpClient) readResponse(resp *http.Response) ([]byte, error) {
	var respData bytes.Buffer
	if err := internal.CopyHTTPResponseBody(&respData, resp, c.maxResponseSize); err != nil {
		return nil, err
	}
	return respData.Bytes(), nil
}

var gzPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

func (c *httpClient) newRequest(ctx context.Context, body []byte) (request, error) {
	r := c.req.Clone(ctx)
	req := request{Request: r}

	switch c.compression {
	case NoCompression:
		r.ContentLength = int64(len(body))
		req.bodyReader = bodyReader(body)
		req.GetBody = bodyReaderErr(body)
	case GzipCompression:
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
		if err := gz.Close(); err != nil {
			return req, err
		}

		req.bodyReader = bodyReader(b.Bytes())
		req.GetBody = bodyReaderErr(b.Bytes())
	}
	return req, nil
}

func bodyReader(buf []byte) func() io.ReadCloser {
	return func() io.ReadCloser {
		return io.NopCloser(bytes.NewReader(buf))
	}
}

func bodyReaderErr(buf []byte) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(buf)), nil
	}
}

type request struct {
	*http.Request

	bodyReader func() io.ReadCloser
}

func (r *request) reset(ctx context.Context) {
	r.Body = r.bodyReader()
	r.Request = r.WithContext(ctx)
}

type retryableError struct {
	throttle time.Duration
	err      error
}

func newResponseError(header http.Header, wrapped error) error {
	var rErr retryableError
	if v := header.Get("Retry-After"); v != "" {
		rErr.throttle = retryAfterDuration(v)
	}
	rErr.err = wrapped
	return rErr
}

func isProtobufContentType(contentType string) bool {
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.EqualFold(mediatype, "application/x-protobuf")
}

func retryAfterDuration(v string) time.Duration {
	if t, err := strconv.ParseInt(v, 10, 64); err == nil && t >= 0 {
		const maxRetryAfterSeconds = math.MaxInt64 / int64(time.Second)
		if t > maxRetryAfterSeconds {
			return time.Duration(math.MaxInt64)
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
		return fmt.Sprintf("retry-able request failure: %v", e.err.Error())
	}
	return "retry-able request failure"
}

func (e retryableError) Unwrap() error {
	return e.err
}

func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}
	rErr, ok := err.(retryableError) //nolint:errorlint // Only match direct retryableError wrapper.
	if !ok {
		return false, 0
	}
	return true, rErr.throttle
}
