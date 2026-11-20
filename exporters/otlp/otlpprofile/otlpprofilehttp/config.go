// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
)

// Default values.
var (
	defaultEndpoint                               = "localhost:4318"
	defaultPath                                   = "/v1development/profiles"
	defaultTimeout                                = 10 * time.Second
	defaultMaxRequestSize                         = 64 * 1024 * 1024
	defaultMaxResponseSize                        = int64(4 * 1024 * 1024)
	defaultProxy           HTTPTransportProxyFunc = http.ProxyFromEnvironment
	defaultRetryCfg                               = retry.DefaultConfig
)

// Environment variable keys.
var (
	envEndpoint = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	}
	envInsecure = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_INSECURE",
		"OTEL_EXPORTER_OTLP_INSECURE",
	}

	envPathSignal = []string{"OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT"}
	envPathOTLP   = []string{"OTEL_EXPORTER_OTLP_ENDPOINT"}

	envHeaders = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_HEADERS",
		"OTEL_EXPORTER_OTLP_HEADERS",
	}

	envCompression = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_COMPRESSION",
		"OTEL_EXPORTER_OTLP_COMPRESSION",
	}

	envTimeout = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT",
		"OTEL_EXPORTER_OTLP_TIMEOUT",
	}

	envTLSCert = []string{
		"OTEL_EXPORTER_OTLP_PROFILES_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_CERTIFICATE",
	}
	envTLSClient = []struct {
		Certificate string
		Key         string
	}{
		{
			"OTEL_EXPORTER_OTLP_PROFILES_CLIENT_CERTIFICATE",
			"OTEL_EXPORTER_OTLP_PROFILES_CLIENT_KEY",
		},
		{
			"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE",
			"OTEL_EXPORTER_OTLP_CLIENT_KEY",
		},
	}
)

// Option applies an option to the Exporter.
type Option interface {
	applyHTTPOption(config) config
}

type fnOpt func(config) config

func (f fnOpt) applyHTTPOption(c config) config { return f(c) }

type endpointPath struct {
	Path    string
	RawPath string
}

type config struct {
	endpoint        setting[string]
	path            setting[endpointPath]
	insecure        setting[bool]
	tlsCfg          setting[*tls.Config]
	headers         setting[map[string]string]
	compression     setting[Compression]
	maxRequestSize  setting[int]
	maxResponseSize setting[int64]
	timeout         setting[time.Duration]
	proxy           setting[HTTPTransportProxyFunc]
	retryCfg        setting[retry.Config]
	httpClient      *http.Client
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		if opt != nil {
			c = opt.applyHTTPOption(c)
		}
	}

	c.endpoint = c.endpoint.Resolve(
		getenv(envEndpoint, convEndpoint),
		fallback(defaultEndpoint),
	)
	c.path = c.path.Resolve(
		getenv(envPathSignal, convPathExact),
		getenv(envPathOTLP, convPath),
		fallback(endpointPath{Path: defaultPath}),
	)
	c.insecure = c.insecure.Resolve(
		loadInsecureFromEnvEndpoint(envEndpoint),
		getenv(envInsecure, convInsecure),
	)
	c.tlsCfg = c.tlsCfg.Resolve(
		loadEnvTLS[*tls.Config](),
	)
	c.headers = c.headers.Resolve(
		getenvRedacted(envHeaders, convHeaders),
	)
	c.compression = c.compression.Resolve(
		getenv(envCompression, convCompression),
	)
	c.timeout = c.timeout.Resolve(
		getenv(envTimeout, convDuration),
		fallback(defaultTimeout),
	)
	c.maxRequestSize = c.maxRequestSize.Resolve(
		fallback(defaultMaxRequestSize),
	)
	c.maxResponseSize = c.maxResponseSize.Resolve(
		fallback(defaultMaxResponseSize),
	)
	c.proxy = c.proxy.Resolve(
		fallback(defaultProxy),
	)
	c.retryCfg = c.retryCfg.Resolve(
		fallback(defaultRetryCfg),
	)

	return c
}

// WithEndpoint sets the target host and optional port (without scheme or path).
func WithEndpoint(endpoint string) Option {
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(endpoint)
		return c
	})
}

// WithEndpointURL sets the target endpoint URL (scheme, host, and path).
// If rawURL cannot be parsed, the option is ignored and an error is reported.
func WithEndpointURL(rawURL string) Option {
	u, err := url.Parse(rawURL)
	if err != nil {
		global.Error(err, "otlpprofile: parse endpoint url", "url", rawURL)
		return fnOpt(func(c config) config { return c })
	}
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(u.Host)
		c.path = newSetting(endpointPath{Path: u.Path, RawPath: u.RawPath})
		c.insecure = insecureFromScheme(c.insecure, u.Scheme)
		return c
	})
}

// Compression describes the compression used for exported payloads.
type Compression int

const (
	// NoCompression indicates that no compression is used.
	NoCompression Compression = iota
	// GzipCompression indicates that gzip compression is used.
	GzipCompression
)

// WithCompression sets the compression strategy for HTTP request bodies.
func WithCompression(compression Compression) Option {
	return fnOpt(func(c config) config {
		c.compression = newSetting(compression)
		return c
	})
}

// WithURLPath sets the URL path to which export requests are sent.
// By default, "/v1development/profiles" is used.
func WithURLPath(urlPath string) Option {
	return fnOpt(func(c config) config {
		c.path = newSetting(endpointPath{Path: urlPath})
		return c
	})
}

// WithTLSClientConfig sets the TLS configuration for HTTP requests.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return fnOpt(func(c config) config {
		c.tlsCfg = newSetting(tlsCfg.Clone())
		return c
	})
}

// WithInsecure disables client transport security for the HTTP connection.
func WithInsecure() Option {
	return fnOpt(func(c config) config {
		c.insecure = newSetting(true)
		return c
	})
}

// WithHeaders sends the provided headers with each HTTP request.
func WithHeaders(headers map[string]string) Option {
	return fnOpt(func(c config) config {
		c.headers = newSetting(headers)
		return c
	})
}

// WithTimeout sets the maximum duration for an export, across retries.
// This takes precedence over any retry settings defined by [WithRetry].
func WithTimeout(duration time.Duration) Option {
	return fnOpt(func(c config) config {
		c.timeout = newSetting(duration)
		return c
	})
}

// WithMaxRequestSize sets the maximum size in bytes of an uncompressed
// serialized request. Values <= 0 disable the limit.
func WithMaxRequestSize(size int) Option {
	return fnOpt(func(c config) config {
		c.maxRequestSize = newSetting(size)
		return c
	})
}

// WithMaxResponseSize sets the maximum size in bytes of a decompressed HTTP
// response body. Values <= 0 are ignored; the default is 4 MiB.
func WithMaxResponseSize(size int64) Option {
	return fnOpt(func(c config) config {
		if size > 0 {
			c.maxResponseSize = newSetting(size)
		}
		return c
	})
}

// RetryConfig defines configuration for retrying failed profile exports.
type RetryConfig retry.Config

// WithRetry sets the retry policy for transient retryable HTTP failures.
func WithRetry(rc RetryConfig) Option {
	return fnOpt(func(c config) config {
		c.retryCfg = newSetting(retry.Config(rc))
		return c
	})
}

// HTTPTransportProxyFunc resolves the proxy URL for an HTTP request.
type HTTPTransportProxyFunc func(*http.Request) (*url.URL, error)

// WithProxy sets the proxy resolver used by the HTTP transport.
func WithProxy(pf HTTPTransportProxyFunc) Option {
	return fnOpt(func(c config) config {
		c.proxy = newSetting(pf)
		return c
	})
}

// WithHTTPClient sets a custom [http.Client] for the exporter. It takes
// precedence over [WithProxy], [WithTimeout], and [WithTLSClientConfig].
func WithHTTPClient(c *http.Client) Option {
	return fnOpt(func(cfg config) config {
		cfg.httpClient = c
		return cfg
	})
}

type setting[T any] struct {
	Value T
	Set   bool
}

func newSetting[T any](value T) setting[T] {
	return setting[T]{Value: value, Set: true}
}

type resolver[T any] func(setting[T]) setting[T]

func (s setting[T]) Resolve(fn ...resolver[T]) setting[T] {
	for _, f := range fn {
		s = f(s)
	}
	return s
}

func loadEnvTLS[T *tls.Config]() resolver[T] {
	return func(s setting[T]) setting[T] {
		if s.Set {
			return s
		}

		var rootCAs *x509.CertPool
		var err error
		for _, key := range envTLSCert {
			if v := os.Getenv(key); v != "" {
				rootCAs, err = loadCertPool(v)
				break
			}
		}

		var certs []tls.Certificate
		for _, pair := range envTLSClient {
			cert := os.Getenv(pair.Certificate)
			key := os.Getenv(pair.Key)
			if cert != "" && key != "" {
				var e error
				certs, e = loadCertificates(cert, key)
				err = errors.Join(err, e)
				break
			}
		}

		if err != nil {
			otel.Handle(fmt.Errorf("failed to load TLS: %w", err))
		} else if rootCAs != nil || certs != nil {
			s.Set = true
			s.Value = &tls.Config{RootCAs: rootCAs, Certificates: certs}
		}
		return s
	}
}

var readFile = os.ReadFile

func loadCertPool(filePath string) (*x509.CertPool, error) {
	b, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM(b); !ok {
		return nil, errors.New("certificate not added")
	}
	return cp, nil
}

func loadCertificates(certPath, keyPath string) ([]tls.Certificate, error) {
	cert, err := readFile(certPath)
	if err != nil {
		return nil, err
	}
	key, err := readFile(keyPath)
	if err != nil {
		return nil, err
	}
	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return []tls.Certificate{crt}, nil
}

func getenv[T any](keys []string, conv func(string) (T, error)) resolver[T] {
	return resolveEnv(keys, conv, false)
}

// getenvRedacted behaves like getenv but omits raw environment variable values
// from error messages so credentials in headers are not logged.
func getenvRedacted[T any](keys []string, conv func(string) (T, error)) resolver[T] {
	return resolveEnv(keys, conv, true)
}

func resolveEnv[T any](keys []string, conv func(string) (T, error), redact bool) resolver[T] {
	return func(s setting[T]) setting[T] {
		if s.Set {
			return s
		}
		for _, key := range keys {
			if vStr := os.Getenv(key); vStr != "" {
				v, err := conv(vStr)
				if err == nil {
					s.Value = v
					s.Set = true
					break
				}
				if redact {
					otel.Handle(fmt.Errorf("invalid %s: %w", key, err))
				} else {
					otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, vStr, err))
				}
			}
		}
		return s
	}
}

func convEndpoint(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

func convPathExact(s string) (endpointPath, error) {
	u, err := url.Parse(s)
	if err != nil {
		return endpointPath{}, err
	}
	if u.Path == "" {
		return endpointPath{Path: "/"}, nil
	}
	return endpointPath{Path: u.Path, RawPath: u.RawPath}, nil
}

func convPath(s string) (endpointPath, error) {
	u, err := url.Parse(s)
	if err != nil {
		return endpointPath{}, err
	}
	suffix := defaultPath
	if strings.HasSuffix(u.EscapedPath(), "/") {
		suffix = strings.TrimPrefix(suffix, "/")
	}
	p := endpointPath{Path: u.Path + suffix, RawPath: u.RawPath}
	if p.RawPath != "" {
		p.RawPath += suffix
	}
	return p, nil
}

func convInsecure(s string) (bool, error) {
	s = strings.ToLower(s)
	if s != "true" && s != "false" {
		return false, fmt.Errorf("can't convert %q to bool", s)
	}
	return s == "true", nil
}

func loadInsecureFromEnvEndpoint(keys []string) resolver[bool] {
	return func(s setting[bool]) setting[bool] {
		if s.Set {
			return s
		}
		for _, key := range keys {
			if vStr := os.Getenv(key); vStr != "" {
				u, err := url.Parse(vStr)
				if err != nil {
					otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, vStr, err))
					continue
				}
				return insecureFromScheme(s, u.Scheme)
			}
		}
		return s
	}
}

func insecureFromScheme(prev setting[bool], scheme string) setting[bool] {
	if scheme == "https" {
		return newSetting(false)
	}
	if scheme != "" {
		return newSetting(true)
	}
	return prev
}

func convHeaders(s string) (map[string]string, error) {
	out := make(map[string]string)
	var err error
	for header := range strings.SplitSeq(s, ",") {
		rawKey, rawVal, found := strings.Cut(header, "=")
		if !found {
			err = errors.Join(err, errors.New("invalid header: missing '='"))
			continue
		}
		key := strings.TrimSpace(rawKey)
		if !isValidHeaderKey(key) {
			err = errors.Join(err, fmt.Errorf("invalid header key: %q", key))
			continue
		}
		escVal, e := url.PathUnescape(rawVal)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid header value encoding for key %q", key))
			continue
		}
		out[key] = strings.TrimSpace(escVal)
	}
	return out, err
}

func convCompression(s string) (Compression, error) {
	switch strings.ToLower(s) {
	case "gzip":
		return GzipCompression, nil
	case "none", "":
		return NoCompression, nil
	}
	return NoCompression, fmt.Errorf("unknown compression: %s", s)
}

func convDuration(s string) (time.Duration, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if d < 0 {
		return 0, errors.New("negative duration")
	}
	if d > math.MaxInt64/int64(time.Millisecond) {
		return 0, errors.New("duration overflows time.Duration")
	}
	return time.Duration(d) * time.Millisecond, nil
}

func fallback[T any](val T) resolver[T] {
	return func(s setting[T]) setting[T] {
		if !s.Set {
			s.Value = val
			s.Set = true
		}
		return s
	}
}

func isValidHeaderKey(key string) bool {
	if key == "" {
		return false
	}
	for _, c := range key {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c rune) bool {
	return c <= unicode.MaxASCII && (unicode.IsLetter(c) ||
		unicode.IsDigit(c) ||
		c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*' ||
		c == '+' || c == '-' || c == '.' || c == '^' || c == '_' || c == '`' || c == '|' || c == '~')
}
