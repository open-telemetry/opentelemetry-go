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

package otlphttp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/internal/otlpconfig"

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
	tlsCert, err := otlpconfig.CreateTLSConfig([]byte(otlpconfig.WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []Option
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *config)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "localhost:4317", c.traces.endpoint)
				assert.Equal(t, "localhost:4317", c.metrics.endpoint)
				assert.Equal(t, NoCompression, c.traces.compression)
				assert.Equal(t, NoCompression, c.metrics.compression)
				assert.Equal(t, map[string]string(nil), c.traces.headers)
				assert.Equal(t, map[string]string(nil), c.metrics.headers)
				assert.Equal(t, 10*time.Second, c.traces.timeout)
				assert.Equal(t, 10*time.Second, c.metrics.timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []Option{
				WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "someendpoint", c.traces.endpoint)
				assert.Equal(t, "someendpoint", c.metrics.endpoint)
			},
		},
		{
			name: "Test With Signal Specific Endpoint",
			opts: []Option{
				WithEndpoint("overrode_by_signal_specific"),
				WithTracesEndpoint("traces_endpoint"),
				WithMetricsEndpoint("metrics_endpoint"),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "traces_endpoint", c.traces.endpoint)
				assert.Equal(t, "metrics_endpoint", c.metrics.endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "env_endpoint", c.traces.endpoint)
				assert.Equal(t, "env_endpoint", c.metrics.endpoint)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":  "env_traces_endpoint",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "env_traces_endpoint", c.traces.endpoint)
				assert.Equal(t, "env_metrics_endpoint", c.metrics.endpoint)
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []Option{
				WithTracesEndpoint("traces_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, "traces_endpoint", c.traces.endpoint)
				assert.Equal(t, "env_endpoint", c.metrics.endpoint)
			},
		},

		// Certificate tests
		{
			name: "Test With Certificate",
			opts: []Option{
				WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.traces.tlsCfg.RootCAs.Subjects())
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.metrics.tlsCfg.RootCAs.Subjects())
			},
		},
		{
			name: "Test With Signal Specific Endpoint",
			opts: []Option{
				WithTLSClientConfig(&tls.Config{}),
				WithTracesTLSClientConfig(tlsCert),
				WithMetricsTLSClientConfig(&tls.Config{RootCAs: x509.NewCertPool()}),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.traces.tlsCfg.RootCAs.Subjects())
				assert.Equal(t, 0, len(c.metrics.tlsCfg.RootCAs.Subjects()))
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(otlpconfig.WeakCertificate),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.traces.tlsCfg.RootCAs.Subjects())
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.metrics.tlsCfg.RootCAs.Subjects())
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_CERTIFICATE":  "cert_path",
				"OTEL_EXPORTER_OTLP_METRICS_CERTIFICATE": "invalid_cert",
			},
			fileReader: fileReader{
				"cert_path":    []byte(otlpconfig.WeakCertificate),
				"invalid_cert": []byte("invalid certificate file."),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.traces.tlsCfg.RootCAs.Subjects())
				assert.Equal(t, (*tls.Config)(nil), c.metrics.tlsCfg)
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []Option{
				WithMetricsTLSClientConfig(&tls.Config{RootCAs: x509.NewCertPool()}),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(otlpconfig.WeakCertificate),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, tlsCert.RootCAs.Subjects(), c.traces.tlsCfg.RootCAs.Subjects())
				assert.Equal(t, 0, len(c.metrics.tlsCfg.RootCAs.Subjects()))
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []Option{
				WithHeaders(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.metrics.headers)
				assert.Equal(t, map[string]string{"h1": "v1"}, c.traces.headers)
			},
		},
		{
			name: "Test With Signal Specific Headers",
			opts: []Option{
				WithHeaders(map[string]string{"overrode": "by_signal_specific"}),
				WithMetricsHeaders(map[string]string{"m1": "mv1"}),
				WithTracesHeaders(map[string]string{"t1": "tv1"}),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.metrics.headers)
				assert.Equal(t, map[string]string{"t1": "tv1"}, c.traces.headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.metrics.headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.traces.headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_HEADERS":  "h1=v1,h2=v2",
				"OTEL_EXPORTER_OTLP_METRICS_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.metrics.headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.traces.headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []Option{
				WithMetricsHeaders(map[string]string{"m1": "mv1"}),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.metrics.headers)
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.traces.headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []Option{
				WithCompression(GzipCompression),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, GzipCompression, c.traces.compression)
				assert.Equal(t, GzipCompression, c.metrics.compression)
			},
		},
		{
			name: "Test With Signal Specific Compression",
			opts: []Option{
				WithCompression(NoCompression), // overrode by signal specific configs
				WithTracesCompression(GzipCompression),
				WithMetricsCompression(GzipCompression),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, GzipCompression, c.traces.compression)
				assert.Equal(t, GzipCompression, c.metrics.compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, GzipCompression, c.traces.compression)
				assert.Equal(t, GzipCompression, c.metrics.compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION":  "gzip",
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, GzipCompression, c.traces.compression)
				assert.Equal(t, GzipCompression, c.metrics.compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []Option{
				WithTracesCompression(NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION":  "gzip",
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, NoCompression, c.traces.compression)
				assert.Equal(t, GzipCompression, c.metrics.compression)
			},
		},

		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []Option{
				WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, 5*time.Second, c.traces.timeout)
				assert.Equal(t, 5*time.Second, c.metrics.timeout)
			},
		},
		{
			name: "Test With Signal Specific Timeout",
			opts: []Option{
				WithTimeout(time.Duration(5 * time.Second)),
				WithTracesTimeout(time.Duration(13 * time.Second)),
				WithMetricsTimeout(time.Duration(14 * time.Second)),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, 13*time.Second, c.traces.timeout)
				assert.Equal(t, 14*time.Second, c.metrics.timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, c.metrics.timeout, 15*time.Second)
				assert.Equal(t, c.traces.timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT":  "27000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, c.traces.timeout, 27*time.Second)
				assert.Equal(t, c.metrics.timeout, 28*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT":  "27000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			opts: []Option{
				WithTracesTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *config) {
				assert.Equal(t, c.traces.timeout, 5*time.Second)
				assert.Equal(t, c.metrics.timeout, 28*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newDefaultConfig()

			e := envOptionsReader{
				getEnv:   tt.env.getEnv,
				readFile: tt.fileReader.readFile,
			}
			e.applyEnvConfigs(&cfg)

			for _, opt := range tt.opts {
				opt.Apply(&cfg)
			}
			tt.asserts(t, &cfg)
		})
	}
}
