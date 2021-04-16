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

package otlpconfig

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp"

	"github.com/stretchr/testify/assert"
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
	return nil, errors.New("File not found")
}

func TestConfigs(t *testing.T) {
	tlsCert, err := CreateTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []GenericOption
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *Config, grpcOption bool)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "localhost:4317", c.Traces.Endpoint)
				assert.Equal(t, "localhost:4317", c.Metrics.Endpoint)
				assert.Equal(t, otlp.NoCompression, c.Traces.Compression)
				assert.Equal(t, otlp.NoCompression, c.Metrics.Compression)
				assert.Equal(t, map[string]string(nil), c.Traces.Headers)
				assert.Equal(t, map[string]string(nil), c.Metrics.Headers)
				assert.Equal(t, 10*time.Second, c.Traces.Timeout)
				assert.Equal(t, 10*time.Second, c.Metrics.Timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []GenericOption{
				WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "someendpoint", c.Traces.Endpoint)
				assert.Equal(t, "someendpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test With Signal Specific Endpoint",
			opts: []GenericOption{
				WithEndpoint("overrode_by_signal_specific"),
				WithTracesEndpoint("traces_endpoint"),
				WithMetricsEndpoint("metrics_endpoint"),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, "metrics_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Traces.Endpoint)
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":  "env_traces_endpoint",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []GenericOption{
				WithTracesEndpoint("traces_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
			},
		},

		// Certificate tests
		{
			name: "Test With Certificate",
			opts: []GenericOption{
				WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					//TODO: make sure gRPC's credentials actually works
					assert.NotNil(t, c.Traces.GRPCCredentials)
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Metrics.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test With Signal Specific Certificate",
			opts: []GenericOption{
				WithTLSClientConfig(&tls.Config{}),
				WithTracesTLSClientConfig(tlsCert),
				WithMetricsTLSClientConfig(&tls.Config{RootCAs: x509.NewCertPool()}),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {

				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
					assert.Equal(t, 0, len(c.Metrics.TLSCfg.RootCAs.Subjects()))
				}
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
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Metrics.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Environment Signal Specific Certificate",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_CERTIFICATE":  "cert_path",
				"OTEL_EXPORTER_OTLP_METRICS_CERTIFICATE": "invalid_cert",
			},
			fileReader: fileReader{
				"cert_path":    []byte(WeakCertificate),
				"invalid_cert": []byte("invalid certificate file."),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
					assert.Nil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
					assert.Equal(t, (*tls.Config)(nil), c.Metrics.TLSCfg)
				}
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []GenericOption{
				WithMetricsTLSClientConfig(&tls.Config{RootCAs: x509.NewCertPool()}),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
					assert.Equal(t, 0, len(c.Metrics.TLSCfg.RootCAs.Subjects()))
				}
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []GenericOption{
				WithHeaders(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Metrics.Headers)
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Traces.Headers)
			},
		},
		{
			name: "Test With Signal Specific Headers",
			opts: []GenericOption{
				WithHeaders(map[string]string{"overrode": "by_signal_specific"}),
				WithMetricsHeaders(map[string]string{"m1": "mv1"}),
				WithTracesHeaders(map[string]string{"t1": "tv1"}),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.Metrics.Headers)
				assert.Equal(t, map[string]string{"t1": "tv1"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_HEADERS":  "h1=v1,h2=v2",
				"OTEL_EXPORTER_OTLP_METRICS_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []GenericOption{
				WithMetricsHeaders(map[string]string{"m1": "mv1"}),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.Metrics.Headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []GenericOption{
				WithCompression(otlp.GzipCompression),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, otlp.GzipCompression, c.Traces.Compression)
				assert.Equal(t, otlp.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test With Signal Specific Compression",
			opts: []GenericOption{
				WithCompression(otlp.NoCompression), // overrode by signal specific configs
				WithTracesCompression(otlp.GzipCompression),
				WithMetricsCompression(otlp.GzipCompression),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, otlp.GzipCompression, c.Traces.Compression)
				assert.Equal(t, otlp.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, otlp.GzipCompression, c.Traces.Compression)
				assert.Equal(t, otlp.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION":  "gzip",
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, otlp.GzipCompression, c.Traces.Compression)
				assert.Equal(t, otlp.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []GenericOption{
				WithTracesCompression(otlp.NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION":  "gzip",
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, otlp.NoCompression, c.Traces.Compression)
				assert.Equal(t, otlp.GzipCompression, c.Metrics.Compression)
			},
		},

		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []GenericOption{
				WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, 5*time.Second, c.Traces.Timeout)
				assert.Equal(t, 5*time.Second, c.Metrics.Timeout)
			},
		},
		{
			name: "Test With Signal Specific Timeout",
			opts: []GenericOption{
				WithTimeout(time.Duration(5 * time.Second)),
				WithTracesTimeout(time.Duration(13 * time.Second)),
				WithMetricsTimeout(time.Duration(14 * time.Second)),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, 13*time.Second, c.Traces.Timeout)
				assert.Equal(t, 14*time.Second, c.Metrics.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Metrics.Timeout, 15*time.Second)
				assert.Equal(t, c.Traces.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT":  "27000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Traces.Timeout, 27*time.Second)
				assert.Equal(t, c.Metrics.Timeout, 28*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT":  "27000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			opts: []GenericOption{
				WithTracesTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Traces.Timeout, 5*time.Second)
				assert.Equal(t, c.Metrics.Timeout, 28*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := EnvOptionsReader{
				GetEnv:   tt.env.getEnv,
				ReadFile: tt.fileReader.readFile,
			}

			// Tests Generic options as HTTP Options
			cfg := NewDefaultConfig()
			e.ApplyHTTPEnvConfigs(&cfg)
			for _, opt := range tt.opts {
				opt.ApplyHTTPOption(&cfg)
			}
			tt.asserts(t, &cfg, false)

			// Tests Generic options as gRPC Options
			cfg = NewDefaultConfig()
			e.ApplyGRPCEnvConfigs(&cfg)
			for _, opt := range tt.opts {
				opt.ApplyGRPCOption(&cfg)
			}
			tt.asserts(t, &cfg, true)
		})
	}
}
