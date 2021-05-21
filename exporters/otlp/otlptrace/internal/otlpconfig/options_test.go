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
	"errors"
	"testing"
	"time"

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
				assert.Equal(t, NoCompression, c.Traces.Compression)
				assert.Equal(t, map[string]string(nil), c.Traces.Headers)
				assert.Equal(t, 10*time.Second, c.Traces.Timeout)
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
			},
		},
		{
			name: "Test With Signal Specific Endpoint",
			opts: []GenericOption{
				WithEndpoint("overrode_by_signal_specific"),
				WithTracesEndpoint("traces_endpoint"),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "traces_endpoint", c.Traces.Endpoint)
			},
		},
		{
			name: "Test Environment Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "env_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Traces.Endpoint)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "env_traces_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Traces.Endpoint)
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
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Traces.Endpoint)
				assert.Equal(t, true, c.Traces.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTP scheme and leading & trailingspaces",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "      http://env_endpoint    ",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Traces.Endpoint)
				assert.Equal(t, true, c.Traces.Insecure)
			},
		},
		{
			name: "Test Environment Endpoint with HTTPS scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://env_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_endpoint", c.Traces.Endpoint)
				assert.Equal(t, false, c.Traces.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "http://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://env_traces_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, true, c.Traces.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint #2",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "http://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://env_traces_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, true, c.Traces.Insecure)
			},
		},
		{
			name: "Test Environment Signal Specific Endpoint with uppercase scheme",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":        "HTTP://overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "HtTp://env_traces_endpoint",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, "env_traces_endpoint", c.Traces.Endpoint)
				assert.Equal(t, true, c.Traces.Insecure)
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
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test With Signal Specific Certificate",
			opts: []GenericOption{
				WithTLSClientConfig(&tls.Config{}),
				WithTracesTLSClientConfig(tlsCert),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {

				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
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
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
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
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
				}
			},
		},
		{
			name: "Test Mixed Environment and With Certificate",
			opts: []GenericOption{},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_CERTIFICATE": "cert_path",
			},
			fileReader: fileReader{
				"cert_path": []byte(WeakCertificate),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				if grpcOption {
					assert.NotNil(t, c.Traces.GRPCCredentials)
				} else {
					assert.Equal(t, tlsCert.RootCAs.Subjects(), c.Traces.TLSCfg.RootCAs.Subjects())
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
				assert.Equal(t, map[string]string{"h1": "v1"}, c.Traces.Headers)
			},
		},
		{
			name: "Test With Signal Specific Headers",
			opts: []GenericOption{
				WithHeaders(map[string]string{"overrode": "by_signal_specific"}),
				WithTracesHeaders(map[string]string{"t1": "tv1"}),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"t1": "tv1"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Environment Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Environment Signal Specific Headers",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_HEADERS":        "overrode_by_signal_specific",
				"OTEL_EXPORTER_OTLP_TRACES_HEADERS": "h1=v1,h2=v2",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},
		{
			name: "Test Mixed Environment and With Headers",
			env:  map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "h1=v1,h2=v2"},
			opts: []GenericOption{},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, map[string]string{"h1": "v1", "h2": "v2"}, c.Traces.Headers)
			},
		},

		// Compression Tests
		{
			name: "Test With Compression",
			opts: []GenericOption{
				WithCompression(GzipCompression),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, GzipCompression, c.Traces.Compression)
			},
		},
		{
			name: "Test With Signal Specific Compression",
			opts: []GenericOption{
				WithCompression(NoCompression), // overrode by signal specific configs
				WithTracesCompression(GzipCompression),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, GzipCompression, c.Traces.Compression)
			},
		},
		{
			name: "Test Environment Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, GzipCompression, c.Traces.Compression)
			},
		},
		{
			name: "Test Environment Signal Specific Compression",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, GzipCompression, c.Traces.Compression)
			},
		},
		{
			name: "Test Mixed Environment and With Compression",
			opts: []GenericOption{
				WithTracesCompression(NoCompression),
			},
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION": "gzip",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, NoCompression, c.Traces.Compression)
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
			},
		},
		{
			name: "Test With Signal Specific Timeout",
			opts: []GenericOption{
				WithTimeout(time.Duration(5 * time.Second)),
				WithTracesTimeout(time.Duration(13 * time.Second)),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, 13*time.Second, c.Traces.Timeout)
			},
		},
		{
			name: "Test Environment Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT": "15000",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Traces.Timeout, 15*time.Second)
			},
		},
		{
			name: "Test Environment Signal Specific Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":        "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT": "27000",
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Traces.Timeout, 27*time.Second)
			},
		},
		{
			name: "Test Mixed Environment and With Timeout",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_TIMEOUT":        "15000",
				"OTEL_EXPORTER_OTLP_TRACES_TIMEOUT": "27000",
			},
			opts: []GenericOption{
				WithTracesTimeout(5 * time.Second),
			},
			asserts: func(t *testing.T, c *Config, grpcOption bool) {
				assert.Equal(t, c.Traces.Timeout, 5*time.Second)
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
