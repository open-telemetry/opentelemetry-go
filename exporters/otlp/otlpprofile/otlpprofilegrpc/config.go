// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal/grpcendpoint"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
)

// Default values.
var (
	defaultEndpoint       = "localhost:4317"
	defaultTimeout        = 10 * time.Second
	defaultMaxRequestSize = 64 * 1024 * 1024
	defaultRetryCfg       = retry.DefaultConfig
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

type fnOpt func(config) config

func (f fnOpt) applyOption(c config) config { return f(c) }

// Option applies an option to the Exporter.
type Option interface {
	applyOption(config) config
}

type config struct {
	endpoint       setting[string]
	insecure       setting[bool]
	tlsCfg         setting[*tls.Config]
	headers        setting[map[string]string]
	compression    setting[compression]
	maxRequestSize setting[int]
	timeout        setting[time.Duration]
	retryCfg       setting[retry.Config]

	gRPCCredentials    setting[credentials.TransportCredentials]
	serviceConfig      setting[string]
	reconnectionPeriod setting[time.Duration]
	dialOptions        setting[[]grpc.DialOption]
	gRPCConn           setting[*grpc.ClientConn]
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		if opt != nil {
			c = opt.applyOption(c)
		}
	}

	endpointFromOpt := c.endpoint.Set
	c.endpoint = c.endpoint.Resolve(
		getEnv(envEndpoint, convEndpoint),
		fallback(defaultEndpoint),
	)
	if !endpointFromOpt {
		// Only derive transport security from an env endpoint when that env
		// endpoint actually supplied c.endpoint, so overriding the endpoint via
		// WithEndpoint or schemeless WithEndpointURL cannot be downgraded.
		c.insecure = c.insecure.Resolve(
			loadInsecureFromEnvEndpoint(envEndpoint),
		)
	}
	c.insecure = c.insecure.Resolve(
		getEnv(envInsecure, convInsecure),
	)
	// Skip loading env TLS certificates when transport security is disabled or
	// overridden by explicit credentials or a caller-provided connection.
	if !c.insecure.Value && c.gRPCCredentials.Value == nil && c.gRPCConn.Value == nil {
		c.tlsCfg = c.tlsCfg.Resolve(
			loadEnvTLS[*tls.Config](),
		)
	}
	c.headers = c.headers.Resolve(
		getEnv(envHeaders, convHeaders),
	)
	c.compression = c.compression.Resolve(
		getEnv(envCompression, convCompression),
	)
	c.timeout = c.timeout.Resolve(
		getEnv(envTimeout, convDuration),
		fallback(defaultTimeout),
	)
	c.maxRequestSize = c.maxRequestSize.Resolve(
		fallback(defaultMaxRequestSize),
	)
	c.retryCfg = c.retryCfg.Resolve(
		fallback(defaultRetryCfg),
	)

	return c
}

// RetryConfig defines configuration for retrying failed profile exports.
type RetryConfig retry.Config

// WithInsecure disables client transport security for the gRPC connection.
// This option has no effect if WithGRPCConn is used.
func WithInsecure() Option {
	return fnOpt(func(c config) config {
		c.insecure = newSetting(true)
		return c
	})
}

// WithEndpoint sets the target gRPC endpoint.
// This option has no effect if WithGRPCConn is used.
func WithEndpoint(endpoint string) Option {
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(endpoint)
		return c
	})
}

// WithEndpointURL sets the target gRPC endpoint URL (scheme and host).
// This option has no effect if WithGRPCConn is used.
func WithEndpointURL(rawURL string) Option {
	u, err := url.Parse(rawURL)
	if err != nil {
		global.Error(err, "otlpprofile: parse endpoint url", "url", rawURL)
		return fnOpt(func(c config) config { return c })
	}
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(grpcendpoint.Target(u))
		if insecure, ok := grpcendpoint.Insecure(u); ok {
			c.insecure = newSetting(insecure)
		}
		return c
	})
}

// WithReconnectionPeriod sets the minimum amount of time between connection
// attempts to the target endpoint. This option has no effect if WithGRPCConn is used.
func WithReconnectionPeriod(rp time.Duration) Option {
	return fnOpt(func(c config) config {
		c.reconnectionPeriod = newSetting(rp)
		return c
	})
}

type compression int

const (
	noCompression compression = iota
	gzipCompression
)

// WithCompressor sets the gRPC compressor ("none" or "gzip", case-insensitive).
// This option has no effect if WithGRPCConn is used.
func WithCompressor(compressor string) Option {
	return fnOpt(func(c config) config {
		c.compression = newSetting(compressorToCompression(compressor))
		return c
	})
}

// WithHeaders sends the provided key-value pairs as gRPC metadata on each export.
func WithHeaders(headers map[string]string) Option {
	return fnOpt(func(c config) config {
		c.headers = newSetting(headers)
		return c
	})
}

// WithTLSCredentials sets the gRPC transport credentials.
// This option has no effect if WithGRPCConn is used.
func WithTLSCredentials(credential credentials.TransportCredentials) Option {
	return fnOpt(func(c config) config {
		c.gRPCCredentials = newSetting(credential)
		return c
	})
}

// WithServiceConfig sets the default gRPC service config JSON.
// This option has no effect if WithGRPCConn is used.
func WithServiceConfig(serviceConfig string) Option {
	return fnOpt(func(c config) config {
		c.serviceConfig = newSetting(serviceConfig)
		return c
	})
}

// WithDialOption appends explicit [grpc.DialOption] values when dialing.
// This option has no effect if WithGRPCConn is used.
func WithDialOption(opts ...grpc.DialOption) Option {
	return fnOpt(func(c config) config {
		c.dialOptions = newSetting(opts)
		return c
	})
}

// WithGRPCConn uses conn directly for all gRPC communication. The caller
// remains responsible for closing conn; [Exporter.Shutdown] will not close it.
func WithGRPCConn(conn *grpc.ClientConn) Option {
	return fnOpt(func(c config) config {
		c.gRPCConn = newSetting(conn)
		return c
	})
}

// WithTimeout sets the maximum total duration the Exporter waits for an
// export, including retries. This takes precedence over [WithRetry].
// Values <= 0 disable the exporter timeout. By default, 10 seconds is used
// unless overridden by OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT or
// OTEL_EXPORTER_OTLP_TIMEOUT.
func WithTimeout(duration time.Duration) Option {
	return fnOpt(func(c config) config {
		c.timeout = newSetting(duration)
		return c
	})
}

// WithMaxRequestSize sets the maximum size in bytes of a serialized export
// request. Values <= 0 disable the limit.
func WithMaxRequestSize(size int) Option {
	return fnOpt(func(c config) config {
		c.maxRequestSize = newSetting(size)
		return c
	})
}

// WithRetry sets the retry policy for transient retryable gRPC errors.
func WithRetry(rc RetryConfig) Option {
	return fnOpt(func(c config) config {
		c.retryCfg = newSetting(retry.Config(rc))
		return c
	})
}

func convCompression(s string) (compression, error) {
	switch strings.ToLower(s) {
	case "gzip":
		return gzipCompression, nil
	case "none", "":
		return noCompression, nil
	}
	return noCompression, fmt.Errorf("unknown compression: %s", s)
}

func convEndpoint(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return grpcendpoint.Target(u), nil
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
					// Parse errors are already reported while resolving envEndpoint.
					continue
				}
				if insecure, ok := grpcendpoint.Insecure(u); ok {
					s.Set = true
					s.Value = insecure
				}
				return s
			}
		}
		return s
	}
}

func convHeaders(s string) (map[string]string, error) {
	out := make(map[string]string)
	var err error
	for header := range strings.SplitSeq(s, ",") {
		rawKey, rawVal, found := strings.Cut(header, "=")
		if !found {
			err = errors.Join(err, fmt.Errorf("invalid header: %s", header))
			continue
		}
		key := strings.TrimSpace(rawKey)
		if !isValidHeaderKey(key) {
			err = errors.Join(err, fmt.Errorf("invalid header key: %s", rawKey))
			continue
		}
		escVal, e := url.PathUnescape(rawVal)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid header value: %s", rawVal))
			continue
		}
		out[key] = strings.TrimSpace(escVal)
	}
	return out, err
}

func convDuration(s string) (time.Duration, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if d < 0 || d > math.MaxInt64/int64(time.Millisecond) {
		return 0, fmt.Errorf("duration out of range: %d", d)
	}
	return time.Duration(d) * time.Millisecond, nil
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

func loadCertPool(path string) (*x509.CertPool, error) {
	b, err := readFile(path)
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

func compressorToCompression(compressor string) compression {
	c, err := convCompression(compressor)
	if err != nil {
		otel.Handle(fmt.Errorf("%w, using no compression as default", err))
		return noCompression
	}
	return c
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

func getEnv[T any](keys []string, conv func(string) (T, error)) resolver[T] {
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
				otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, vStr, err))
			}
		}
		return s
	}
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
