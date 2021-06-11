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

package otlpconfig_test

import (
	"errors"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpconfig"

	"github.com/stretchr/testify/assert"
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
	return nil, errors.New("File not found")
}

func TestConfigs(t *testing.T) {
	tlsCert, err := otlpconfig.CreateTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []otlpconfig.GenericOption
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *otlpconfig.Config, grpcOption bool)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "localhost:4317", c.Metrics.Endpoint)
				assert.Equal(t, otlpconfig.NoCompression, c.Metrics.Compression)
				assert.Equal(t, map[string]string(nil), c.Metrics.Headers)
				assert.Equal(t, 10*time.Second, c.Metrics.Timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "someendpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithEndpoint("metrics_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "metrics_endpoint", c.Metrics.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme and leading & trailingspaces",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "      http://env_endpoint    ",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, true, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTPS scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "http://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "https://env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint #2",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "http://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint with uppercase scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "HTTP://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_metrics_endpoint", c.Metrics.Endpoint)
				assert.Equal(t, false, c.Metrics.Insecure)
			},
		},

		// Certificate tests
		{
			name: "Test With Certificate",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				if grpcOption {
					//TODO: make sure gRPC's credentials actually works
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Metrics.TLSCfg.RootCAs.Subjects())
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
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Metrics.TLSCfg.RootCAs.Subjects())
				}
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
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Metrics.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []otlpconfig.GenericOption{},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Metrics.GRPCCredentials)
				} else {
					assert.Equal(t, 1, len(c.Metrics.TLSCfg.RootCAs.Subjects()))
				}
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithHeaders(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Metrics.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithHeaders(map[string]string{"m1": "mv1"}),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.Metrics.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithCompression(otlpconfig.GzipCompression),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, otlpconfig.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, otlpconfig.GzipCompression, c.Metrics.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, otlpconfig.GzipCompression, c.Metrics.Compression)
			},
		},
		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, 5*time.Second, c.Metrics.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Metrics.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Metrics.Timeout, 28*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			opts: []otlpconfig.GenericOption{
				otlpconfig.WithTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *otlpconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Metrics.Timeout, 5*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := otlpconfig.EnvOptionsReader{
				GetEnv:   tt.env.getEnv,
				ReadFile: tt.fileReader.readFile,
			}

			// Tests Generic options as HTTP Options
			cfg := otlpconfig.NewDefaultConfig()
			e.ApplyHTTPEnvConfigs(&cfg)
			for _, opt := range tt.opts {
				opt.ApplyHTTPOption(&cfg)
			}
			tt.asserts(t, &cfg, false)

			// Tests Generic options as gRPC Options
			cfg = otlpconfig.NewDefaultConfig()
			e.ApplyGRPCEnvConfigs(&cfg)
			for _, opt := range tt.opts {
				opt.ApplyGRPCOption(&cfg)
			}
			tt.asserts(t, &cfg, true)
		})
	}
}
