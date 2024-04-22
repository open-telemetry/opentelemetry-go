// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oconf

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/envconfig"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	WeakCertificate = `
-----BEGIN CERTIFICATE-----
MIIBhzCCASygAwIBAgIRANHpHgAWeTnLZpTSxCKs0ggwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHb3RlbC1nbzAeFw0yMTA0MDExMzU5MDNaFw0yMTA0MDExNDU5MDNa
MBIxEDAOBgNVBAoTB290ZWwtZ28wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAS9
nWSkmPCxShxnp43F+PrOtbGV7sNfkbQ/kxzi9Ego0ZJdiXxkmv/C05QFddCW7Y0Z
sJCLHGogQsYnWJBXUZOVo2MwYTAOBgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYI
KwYBBQUHAwEwDAYDVR0TAQH/BAIwADAsBgNVHREEJTAjgglsb2NhbGhvc3SHEAAA
AAAAAAAAAAAAAAAAAAGHBH8AAAEwCgYIKoZIzj0EAwIDSQAwRgIhANwZVVKvfvQ/
1HXsTvgH+xTQswOwSSKYJ1cVHQhqK7ZbAiEAus8NxpTRnp5DiTMuyVmhVNPB+bVH
Lhnm4N/QDk5rek0=
-----END CERTIFICATE-----
`
	WeakPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgN8HEXiXhvByrJ1zK
SFT6Y2l2KqDWwWzKf+t4CyWrNKehRANCAAS9nWSkmPCxShxnp43F+PrOtbGV7sNf
kbQ/kxzi9Ego0ZJdiXxkmv/C05QFddCW7Y0ZsJCLHGogQsYnWJBXUZOV
-----END PRIVATE KEY-----
`
)

type env map[string]string

func (e *env) getEnv(env string) string {
	return (*e)[env]
}

type fileReader map[string][]byte

func (f *fileReader) readFile(filename string) ([]byte, error) {
	if b, ok := (*f)[filename]; ok {
		return b, nil
	}
	return nil, errors.New("file not found")
}

func TestConfigs(t *testing.T) {
	tlsCert, err := CreateTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []GRPCOption
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *Config)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "localhost:4317", c.Metrics.Endpoint)
				assert.Equal(t, NoCompression, c.Metrics.Compression)
				assert.Equal(t, map[string]string(nil), c.Metrics.Headers)
				assert.Equal(t, 10*time.Second, c.Metrics.Timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []GRPCOption{
				WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "someendpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test With Endpoint URL",
			opts: []GRPCOption{
				WithEndpointURL("http://someendpoint/somepath"),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "someendpoint", c.Metrics.Endpoint)
				assert.Equal(t, "/somepath", c.Metrics.URLPath)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},
		{
			name: "Test With Secure Endpoint URL",
			opts: []GRPCOption{
				WithEndpointURL("https://someendpoint/somepath"),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "someendpoint", c.Metrics.Endpoint)
				assert.Equal(t, "/somepath", c.Metrics.URLPath)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},
		{
			name: "Test With Invalid Endpoint URL",
			opts: []GRPCOption{
				WithEndpointURL("%invalid"),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "localhost:4317", c.Metrics.Endpoint)
				assert.Equal(t, "/v1/metrics", c.Metrics.URLPath)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env.endpoint/prefix",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.False(t, c.Metrics.Insecure)
				assert.Equal(t, "env.endpoint/prefix", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "https://overrode.by.signal.specific/env/var",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "http://env.metrics.endpoint",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.True(t, c.Metrics.Insecure)
				assert.Equal(t, "env.metrics.endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []GRPCOption{
				WithEndpoint("metrics_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "metrics_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env_endpoint",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme and leading & trailingspaces",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "      http://env_endpoint    ",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTPS scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env_endpoint",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint with uppercase scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "HTTPS://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "HtTp://env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},

		// Certificate tests
		{
			name: "Test Default Certificate",
			asserts: func(t *testing.T, c *Config) {
				assert.NotNil(t, c.Metrics.GRPCCredentials)
			},
		},
		{
			name: "Test With Certificate",
			opts: []GRPCOption{
				WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *Config) {
				// TODO: make sure gRPC's credentials actually works
				assert.NotNil(t, c.Metrics.GRPCCredentials)
			},
		},
		{
			name: "Test Environment Certificate",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.NotNil(t, c.Metrics.GRPCCredentials)
			},
		},
		{
			name: "Test Environment Signal Specific Certificate",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path":    []byte(WeakCertificate),
				"invalid_cert": []byte("invalid certificate file."),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.NotNil(t, c.Metrics.GRPCCredentials)
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []GRPCOption{},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.NotNil(t, c.Metrics.GRPCCredentials)
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []GRPCOption{
				WithHeaders(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []GRPCOption{
				WithHeaders(map[string]string{"m1": "mv1"}),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.Metrics.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []GRPCOption{
				WithCompression(GzipCompression),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []GRPCOption{
				WithCompression(NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, NoCompression, c.Metrics.Compression)
			},
		},

		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []GRPCOption{
				WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, 5*time.Second, c.Metrics.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, c.Metrics.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, c.Metrics.Timeout, 28*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			opts: []GRPCOption{
				WithTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.Equal(t, c.Metrics.Timeout, 5*time.Second)
			},
		},

		// Temporality Selector Tests
		{
			name: "WithTemporalitySelector",
			opts: []GRPCOption{
				WithTemporalitySelector(deltaSelector),
			},
			asserts: func(t *testing.T, c *Config) {
				// Function value comparisons are disallowed, test non-default
				// behavior of a TemporalitySelector here to ensure our "catch
				// all" was set.
				var undefinedKind metric.InstrumentKind
				got := c.Metrics.TemporalitySelector
				assert.Equal(t, metricdata.DeltaTemporality, got(undefinedKind))
			},
		},

		// Aggregation Selector Tests
		{
			name: "WithAggregationSelector",
			opts: []GRPCOption{
				WithAggregationSelector(dropSelector),
			},
			asserts: func(t *testing.T, c *Config) {
				// Function value comparisons are disallowed, test non-default
				// behavior of a AggregationSelector here to ensure our "catch
				// all" was set.
				var undefinedKind metric.InstrumentKind
				got := c.Metrics.AggregationSelector
				assert.Equal(t, metric.AggregationDrop{}, got(undefinedKind))
			},
		},

		// Proxy Tests
		{
			name: "Test With Proxy",
			opts: []GRPCOption{
				WithProxy(func(r *http.Request) (*url.URL, error) {
					return url.Parse("http://proxy.com")
				}),
			},
			asserts: func(t *testing.T, c *Config) {
				assert.NotNil(t, c.Metrics.Proxy)
				proxyURL, err := c.Metrics.Proxy(&http.Request{})
				assert.NoError(t, err)
				assert.Equal(t, "http://proxy.com", proxyURL.String())
			},
		},
		{
			name: "Test Without Proxy",
			opts: []GRPCOption{},
			asserts: func(t *testing.T, c *Config) {
				assert.Nil(t, c.Metrics.Proxy)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origEOR := DefaultEnvOptionsReader
			DefaultEnvOptionsReader = envconfig.EnvOptionsReader{
				GetEnv:    tt.env.getEnv,
				ReadFile:  tt.fileReader.readFile,
				Namespace: "OTEL_EXPORTER_OTLP",
			}
			t.Cleanup(func() { DefaultEnvOptionsReader = origEOR })

			// Tests Generic options as gRPC Options
			cfg := NewGRPCConfig(tt.opts...)
			tt.asserts(t, &cfg)
		})
	}
}

func dropSelector(metric.InstrumentKind) metric.Aggregation {
	return metric.AggregationDrop{}
}

func deltaSelector(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func TestCleanPath(t *testing.T) {
	type args struct {
		urlPath     string
		defaultPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "clean empty path",
			args: args{
				urlPath:     "",
				defaultPath: "DefaultPath",
			},
			want: "DefaultPath",
		},
		{
			name: "clean metrics path",
			args: args{
				urlPath:     "/prefix/v1/metrics",
				defaultPath: "DefaultMetricsPath",
			},
			want: "/prefix/v1/metrics",
		},
		{
			name: "clean traces path",
			args: args{
				urlPath:     "https://env_endpoint",
				defaultPath: "DefaultTracesPath",
			},
			want: "/https:/env_endpoint",
		},
		{
			name: "spaces trimmed",
			args: args{
				urlPath: " /dir",
			},
			want: "/dir",
		},
		{
			name: "clean path empty",
			args: args{
				urlPath:     "dir/..",
				defaultPath: "DefaultTracesPath",
			},
			want: "DefaultTracesPath",
		},
		{
			name: "make absolute",
			args: args{
				urlPath: "dir/a",
			},
			want: "/dir/a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanPath(tt.args.urlPath, tt.args.defaultPath); got != tt.want {
				t.Errorf("CleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
