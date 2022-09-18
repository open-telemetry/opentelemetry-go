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

package envconfig_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/internal/envconfig"
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

func TestTraceConfigs(t *testing.T) {
	tlsCert, err := envconfig.CreateTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []envconfig.GenericOption
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *envconfig.Config, grpcOption bool)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.Equal(t, "localhost:4317", c.Sc.Endpoint)
				} else {
					assert.Equal(t, "localhost:4318", c.Sc.Endpoint)
				}
				assert.Equal(t, envconfig.NoCompression, c.Sc.Compression)
				assert.Equal(t, map[string]string(nil), c.Sc.Headers)
				assert.Equal(t, 10*time.Second, c.Sc.Timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []envconfig.GenericOption{
				envconfig.WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "someendpoint", c.Sc.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env.endpoint/prefix",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.False(t, c.Sc.Insecure)
				if grpcOption {
					assert.Equal(t, "env.endpoint/prefix", c.Sc.Endpoint)
				} else {
					assert.Equal(t, "env.endpoint", c.Sc.Endpoint)
					assert.Equal(t, "/prefix/v1/traces", c.Sc.URLPath)
				}
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "https://overrode.by.signal.specific/env/var",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://env.traces.endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.True(t, c.Sc.Insecure)
				assert.Equal(t, "env.traces.endpoint", c.Sc.Endpoint)
				if !grpcOption {
					assert.Equal(t, "/", c.Sc.URLPath)
				}
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []envconfig.GenericOption{
				envconfig.WithEndpoint("traces_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "traces_endpoint", c.Sc.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme and leading & trailingspaces",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "      http://env_endpoint    ",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTPS scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, false, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint with uppercase scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "HTTPS://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "HtTp://env_traces_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},

		// Certificate tests
		{
			name: "Test Default Certificate",
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					assert.Nil(t, c.Sc.TLSCfg)
				}
			},
		},
		{
			name: "Test With Certificate",
			opts: []envconfig.GenericOption{
				envconfig.WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					//TODO: make sure gRPC's credentials actually works
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
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
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Environment Signal Specific Certificate",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE":        "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path":    []byte(WeakCertificate),
				"invalid_cert": []byte("invalid certificate file."),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []envconfig.GenericOption{},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
				}
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []envconfig.GenericOption{
				envconfig.WithHeader(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":        "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []envconfig.GenericOption{},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Sc.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []envconfig.GenericOption{
				envconfig.WithCompression(envconfig.GzipCompression),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []envconfig.GenericOption{
				envconfig.WithCompression(envconfig.NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.NoCompression, c.Sc.Compression)
			},
		},

		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []envconfig.GenericOption{
				envconfig.WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, 5*time.Second, c.Sc.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":        "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT": "27000",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 27*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":        "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT": "27000",
			},
			opts: []envconfig.GenericOption{
				envconfig.WithTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 5*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origEOR := envconfig.DefaultEnvOptionsReader
			envconfig.DefaultEnvOptionsReader = envconfig.EnvOptionsReader{
				GetEnv:    tt.env.getEnv,
				ReadFile:  tt.fileReader.readFile,
				Namespace: "OTEL_EXPORTER_OTLP",
			}
			t.Cleanup(func() { envconfig.DefaultEnvOptionsReader = origEOR })

			// Tests Generic options as HTTP Options
			cfg := envconfig.NewHTTPTraceConfig(asHTTPOptions(tt.opts)...)
			tt.asserts(t, &cfg, false)

			// Tests Generic options as gRPC Options
			cfg = envconfig.NewGRPCTraceConfig(asGRPCOptions(tt.opts)...)
			tt.asserts(t, &cfg, true)
		})
	}
}

func TestMetricsConfigs(t *testing.T) {
	tlsCert, err := envconfig.CreateTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		opts       []envconfig.GenericOption
		env        env
		fileReader fileReader
		asserts    func(t *testing.T, c *envconfig.Config, grpcOption bool)
	}{
		{
			name: "Test default configs",
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.Equal(t, "localhost:4317", c.Sc.Endpoint)
				} else {
					assert.Equal(t, "localhost:4318", c.Sc.Endpoint)
				}
				assert.Equal(t, envconfig.NoCompression, c.Sc.Compression)
				assert.Equal(t, map[string]string(nil), c.Sc.Headers)
				assert.Equal(t, 10*time.Second, c.Sc.Timeout)
			},
		},

		// Endpoint Tests
		{
			name: "Test With Endpoint",
			opts: []envconfig.GenericOption{
				envconfig.WithEndpoint("someendpoint"),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "someendpoint", c.Sc.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env.endpoint/prefix",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.False(t, c.Sc.Insecure)
				if grpcOption {
					assert.Equal(t, "env.endpoint/prefix", c.Sc.Endpoint)
				} else {
					assert.Equal(t, "env.endpoint", c.Sc.Endpoint)
					assert.Equal(t, "/prefix/v1/metrics", c.Sc.URLPath)
				}
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "https://overrode.by.signal.specific/env/var",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "http://env.metrics.endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.True(t, c.Sc.Insecure)
				assert.Equal(t, "env.metrics.endpoint", c.Sc.Endpoint)
				if !grpcOption {
					assert.Equal(t, "/", c.Sc.URLPath)
				}
			},
		},
		{
			name: "Test Mixed Environment and With Endpoint",
			opts: []envconfig.GenericOption{
				envconfig.WithEndpoint("metrics_endpoint"),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "metrics_endpoint", c.Sc.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme and leading & trailingspaces",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "      http://env_endpoint    ",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTPS scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Sc.Endpoint)
				assert.Equal(t, false, c.Sc.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint with uppercase scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":         "HTTPS://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "HtTp://env_metrics_endpoint",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, "env_metrics_endpoint", c.Sc.Endpoint)
				assert.Equal(t, true, c.Sc.Insecure)
			},
		},

		// Certificate tests
		{
			name: "Test Default Certificate",
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					assert.Nil(t, c.Sc.TLSCfg)
				}
			},
		},
		{
			name: "Test With Certificate",
			opts: []envconfig.GenericOption{
				envconfig.WithTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					//TODO: make sure gRPC's credentials actually works
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
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
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
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
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Sc.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []envconfig.GenericOption{},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Sc.GRPCCredentials)
				} else {
					// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
					assert.Equal(t, 1, len(c.Sc.TLSCfg.RootCAs.Subjects()))
				}
			},
		},

		// Headers tests
		{
			name: "Test With Headers",
			opts: []envconfig.GenericOption{
				envconfig.WithHeader(map[string]string{"h1": "v1"}),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":         "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_METRICS_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Sc.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []envconfig.GenericOption{
				envconfig.WithHeader(map[string]string{"m1": "mv1"}),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"m1": "mv1"}, c.Sc.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []envconfig.GenericOption{
				envconfig.WithCompression(envconfig.GzipCompression),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.GzipCompression, c.Sc.Compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []envconfig.GenericOption{
				envconfig.WithCompression(envconfig.NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, envconfig.NoCompression, c.Sc.Compression)
			},
		},

		// Timeout Tests
		{
			name: "Test With Timeout",
			opts: []envconfig.GenericOption{
				envconfig.WithTimeout(time.Duration(5 * time.Second)),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, 5*time.Second, c.Sc.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 28*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":         "15000",
				"OTEL_EXPORTER_OTLP_METRICS_TIMEOUT": "28000",
			},
			opts: []envconfig.GenericOption{
				envconfig.WithTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *envconfig.Config, grpcOption bool) {
				assert.Equal(t, c.Sc.Timeout, 5*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origEOR := envconfig.DefaultEnvOptionsReader
			envconfig.DefaultEnvOptionsReader = envconfig.EnvOptionsReader{
				GetEnv:    tt.env.getEnv,
				ReadFile:  tt.fileReader.readFile,
				Namespace: "OTEL_EXPORTER_OTLP",
			}
			t.Cleanup(func() { envconfig.DefaultEnvOptionsReader = origEOR })

			// Tests Generic options as HTTP Options
			cfg := envconfig.NewHTTPMetricsConfig(asHTTPOptions(tt.opts)...)
			tt.asserts(t, &cfg, false)

			// Tests Generic options as gRPC Options
			cfg = envconfig.NewGRPCMetricsConfig(asGRPCOptions(tt.opts)...)
			tt.asserts(t, &cfg, true)
		})
	}
}

func asHTTPOptions(opts []envconfig.GenericOption) []envconfig.HTTPOption {
	converted := make([]envconfig.HTTPOption, len(opts))
	for i, o := range opts {
		converted[i] = envconfig.NewHTTPOption(o.ApplyHTTPOption)
	}
	return converted
}

func asGRPCOptions(opts []envconfig.GenericOption) []envconfig.GRPCOption {
	converted := make([]envconfig.GRPCOption, len(opts))
	for i, o := range opts {
		converted[i] = envconfig.NewGRPCOption(o.ApplyGRPCOption)
	}
	return converted
}
