// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tlsCfg := &tls.Config{}
	headers := map[string]string{"a": "A"}
	rc := RetryConfig{}

	testcases := []struct {
		name    string
		options []Option
		envars  map[string]string
		want    config
	}{
		{
			name: "Defaults",
			want: config{
				endpoint: newSetting(defaultEndpoint),
				path:     newSetting(defaultPath),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "Options",
			options: []Option{
				WithEndpoint("test"),
				WithURLPath("/path"),
				WithInsecure(),
				WithTLSClientConfig(tlsCfg),
				WithCompression(GzipCompression),
				WithHeaders(headers),
				WithTimeout(time.Second),
				WithRetry(rc),
				// Do not test WithProxy. Requires func comparison.
			},
			want: config{
				endpoint:    newSetting("test"),
				path:        newSetting("/path"),
				insecure:    newSetting(true),
				tlsCfg:      newSetting(tlsCfg),
				headers:     newSetting(headers),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(time.Second),
				retryCfg:    newSetting(rc),
			},
		},
		{
			name: "WithEndpointURL",
			options: []Option{
				WithEndpointURL("http://test:8080/path"),
			},
			want: config{
				endpoint: newSetting("test:8080"),
				path:     newSetting("/path"),
				insecure: newSetting(true),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "EndpointPrecidence",
			options: []Option{
				WithEndpointURL("https://test:8080/path"),
				WithEndpoint("not-test:9090"),
				WithURLPath("/alt"),
				WithInsecure(),
			},
			want: config{
				endpoint: newSetting("not-test:9090"),
				path:     newSetting("/alt"),
				insecure: newSetting(true),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "EndpointURLPrecidence",
			options: []Option{
				WithEndpoint("not-test:9090"),
				WithURLPath("/alt"),
				WithInsecure(),
				WithEndpointURL("https://test:8080/path"),
			},
			want: config{
				endpoint: newSetting("test:8080"),
				path:     newSetting("/path"),
				insecure: newSetting(false),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "LogEnvironmentVariables",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":    "https://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE":    "true",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":     "h1=v1,h2=v2",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION": "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":     "15000",
				// TODO: TLS.
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix/v1/traces"),
				insecure: newSetting(false),
				headers: newSetting(map[string]string{
					"h1": "v1",
					"h2": "v2",
				}),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			c := newConfig(tc.options)
			// Cannot compare funcs
			c.proxy = setting[HTTPTransportProxyFunc]{}
			assert.Equal(t, tc.want, c)
		})
	}
}

func TestWithProxy(t *testing.T) {
	proxy := func(*http.Request) (*url.URL, error) { return nil, nil }
	opts := []Option{WithProxy(HTTPTransportProxyFunc(proxy))}
	c := newConfig(opts)

	assert.True(t, c.proxy.Set)
	assert.NotNil(t, c.proxy.Value)
}
