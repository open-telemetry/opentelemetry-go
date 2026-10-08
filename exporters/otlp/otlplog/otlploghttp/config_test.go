// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal/retry"
)

const (
	weakCertificate = `
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
	weakPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgN8HEXiXhvByrJ1zK
SFT6Y2l2KqDWwWzKf+t4CyWrNKehRANCAAS9nWSkmPCxShxnp43F+PrOtbGV7sNf
kbQ/kxzi9Ego0ZJdiXxkmv/C05QFddCW7Y0ZsJCLHGogQsYnWJBXUZOV
-----END PRIVATE KEY-----
`
	mismatchedPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgHBGLNXKfmMmxbPCc
8x4FlFAjxMVItDKT97/+z38BBfihRANCAAQPT2hf0xJYsDB3jTZW5qgN1KC62ByN
e2hJ8B5tcYFbwKqz3lqHFADw5oi9X9VgvQrhaObs7FrIJ0rbP9HdT2uP
-----END PRIVATE KEY-----
`
)

func newTLSConf(cert, key []byte) (*tls.Config, error) {
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM(cert); !ok {
		return nil, errors.New("failed to append certificate to the cert pool")
	}
	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	crts := []tls.Certificate{crt}
	return &tls.Config{RootCAs: cp, Certificates: crts}, nil
}

func TestNewConfig(t *testing.T) {
	orig := readFile
	readFile = func() func(name string) ([]byte, error) {
		index := map[string][]byte{
			"cert_path":    []byte(weakCertificate),
			"key_path":     []byte(weakPrivateKey),
			"invalid_cert": []byte("invalid certificate file."),
			"invalid_key":  []byte("invalid key file."),
		}
		return func(name string) ([]byte, error) {
			b, ok := index[name]
			if !ok {
				err := fmt.Errorf("file does not exist: %s", name)
				return nil, err
			}
			return b, nil
		}
	}()
	t.Cleanup(func() { readFile = orig })

	tlsCfg, err := newTLSConf([]byte(weakCertificate), []byte(weakPrivateKey))
	require.NoError(t, err, "testing TLS config")

	headers := map[string]string{"a": "A"}
	rc := retry.Config{}
	const (
		httpEndpointURL    = "http" + "://test:8080/path"
		httpsEndpointURL   = "https" + "://test:8080/path"
		httpsEndpointURL2  = "https" + "://test2/path2"
		ignoredEndpointURL = "http" + "://ignored:9090/alt"
	)

	testcases := []struct {
		name    string
		options []Option
		envars  map[string]string
		want    config
		errs    []string
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
				WithMaxRequestSize(1),
				WithMaxResponseSize(2),
				WithTimeout(time.Second),
				WithRetry(RetryConfig(rc)),
				// Do not test WithProxy. Requires func comparison.
			},
			want: config{
				endpoint:        newSetting("test"),
				path:            newSetting("/path"),
				insecure:        newSetting(true),
				tlsCfg:          newSetting(tlsCfg),
				headers:         newSetting(headers),
				compression:     newSetting(GzipCompression),
				maxRequestSize:  newSetting(1),
				maxResponseSize: newSetting(int64(2)),
				timeout:         newSetting(time.Second),
				retryCfg:        newSetting(rc),
			},
		},
		{
			name: "WithEndpointURL",
			options: []Option{
				WithEndpointURL(httpEndpointURL),
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
			name: "EndpointPrecedence",
			options: []Option{
				WithEndpointURL(httpsEndpointURL),
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
			name: "EndpointURLPrecedence",
			options: []Option{
				WithEndpoint("not-test:9090"),
				WithURLPath("/alt"),
				WithInsecure(),
				WithEndpointURL(httpsEndpointURL),
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
			name: "WithEndpointURL secure when Environment Endpoint is set insecure",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "http://env.endpoint:8080/prefix",
			},
			options: []Option{
				WithEndpointURL(httpsEndpointURL),
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
			name: "WithEndpointURL secure when Environment insecure is set false",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE": "true",
			},
			options: []Option{
				WithEndpointURL(httpsEndpointURL),
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
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "https://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "a=A",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "key_path",
			},
			want: config{
				endpoint:    newSetting("env.endpoint:8080"),
				path:        newSetting("/prefix"),
				insecure:    newSetting(false),
				tlsCfg:      newSetting(tlsCfg),
				headers:     newSetting(headers),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
		{
			name: "LogEndpointEnvironmentVariablesDefaultPath",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "http://env.endpoint",
			},
			want: config{
				endpoint: newSetting("env.endpoint"),
				path:     newSetting("/"),
				insecure: newSetting(true),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "OTLPEnvironmentVariables",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":           "http://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_HEADERS":            "a=A",
				"OTEL_EXPORTER_OTLP_COMPRESSION":        "none",
				"OTEL_EXPORTER_OTLP_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "key_path",
			},
			want: config{
				endpoint:    newSetting("env.endpoint:8080"),
				path:        newSetting("/prefix/v1/logs"),
				insecure:    newSetting(true),
				tlsCfg:      newSetting(tlsCfg),
				headers:     newSetting(headers),
				compression: newSetting(NoCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
		{
			name: "OTLPEnvironmentVariablesTrailingSlash",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env.endpoint:8080/prefix/",
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix/v1/logs"),
				insecure: newSetting(true),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "OTLPEndpointEnvironmentVariablesDefaultPath",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env.endpoint",
			},
			want: config{
				endpoint: newSetting("env.endpoint"),
				path:     newSetting(defaultPath),
				insecure: newSetting(true),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
		},
		{
			name: "EnvironmentVariablesPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":           ignoredEndpointURL,
				"OTEL_EXPORTER_OTLP_HEADERS":            "b=B",
				"OTEL_EXPORTER_OTLP_COMPRESSION":        "none",
				"OTEL_EXPORTER_OTLP_TIMEOUT":            "30000",
				"OTEL_EXPORTER_OTLP_CERTIFICATE":        "invalid_cert",
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "invalid_cert",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "invalid_key",

				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "https://env.endpoint:8080/path",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "a=A",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "key_path",
			},
			want: config{
				endpoint:    newSetting("env.endpoint:8080"),
				path:        newSetting("/path"),
				insecure:    newSetting(false),
				tlsCfg:      newSetting(tlsCfg),
				headers:     newSetting(headers),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
		{
			name: "OptionsPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":           ignoredEndpointURL,
				"OTEL_EXPORTER_OTLP_HEADERS":            "b=B",
				"OTEL_EXPORTER_OTLP_COMPRESSION":        "none",
				"OTEL_EXPORTER_OTLP_TIMEOUT":            "30000",
				"OTEL_EXPORTER_OTLP_CERTIFICATE":        "invalid_cert",
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "invalid_cert",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "invalid_key",

				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "https://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "a=A",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "key_path",
			},
			options: []Option{
				WithEndpoint("test"),
				WithEndpointURL(httpsEndpointURL2),
				WithURLPath("/path"),
				WithInsecure(),
				WithTLSClientConfig(tlsCfg),
				WithCompression(GzipCompression),
				WithHeaders(headers),
				WithTimeout(time.Second),
				WithRetry(RetryConfig(rc)),
			},
			want: config{
				endpoint:    newSetting("test2"),
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
			name: "InvalidEnvironmentVariables",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "%invalid",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "invalid key=value",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "xz",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "100 seconds",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "invalid_cert",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "invalid_cert",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "invalid_key",
			},
			want: config{
				endpoint: newSetting(defaultEndpoint),
				path:     newSetting(defaultPath),
				timeout:  newSetting(defaultTimeout),
				retryCfg: newSetting(defaultRetryCfg),
			},
			errs: []string{
				`invalid OTEL_EXPORTER_OTLP_LOGS_ENDPOINT value %invalid: parse "%invalid": invalid URL escape "%in"`,
				`failed to load TLS:`,
				`certificate not added`,
				`tls: failed to find any PEM data in certificate input`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_HEADERS value invalid key=value: invalid header key: invalid key`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_COMPRESSION value xz: unknown compression: xz`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_TIMEOUT value 100 seconds: strconv.Atoi: parsing "100 seconds": invalid syntax`,
			},
		},
		{
			name: "with percent-encoded headers",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "https://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "user%2Did=42,user%20name=alice%20smith",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "key_path",
			},
			want: config{
				endpoint:    newSetting("env.endpoint:8080"),
				path:        newSetting("/prefix"),
				insecure:    newSetting(false),
				tlsCfg:      newSetting(tlsCfg),
				headers:     newSetting(map[string]string{"user%2Did": "42", "user%20name": "alice smith"}),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
		{
			name: "with invalid header key",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "https://env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "valid-key=value,invalid key=value",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "gzip",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "15000",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "cert_path",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "key_path",
			},
			want: config{
				endpoint:    newSetting("env.endpoint:8080"),
				path:        newSetting("/prefix"),
				insecure:    newSetting(false),
				tlsCfg:      newSetting(tlsCfg),
				compression: newSetting(GzipCompression),
				timeout:     newSetting(15 * time.Second),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
		{
			name: "OptionEndpointURLWithoutScheme",
			options: []Option{
				WithEndpointURL("//env.endpoint:8080/prefix"),
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix"),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
		{
			name: "EnvEndpointWithoutScheme",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "//env.endpoint:8080/prefix",
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix"),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
		{
			name: "DefaultEndpointWithEnvInsecure",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE": "true",
			},
			want: config{
				endpoint: newSetting(defaultEndpoint),
				path:     newSetting(defaultPath),
				insecure: newSetting(true),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
		{
			name: "EnvEndpointWithoutSchemeWithEnvInsecure",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "//env.endpoint:8080/prefix",
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE": "true",
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix"),
				insecure: newSetting(true),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
		{
			name: "OptionEndpointURLWithoutSchemeWithEnvInsecure",
			options: []Option{
				WithEndpointURL("//env.endpoint:8080/prefix"),
			},
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE": "true",
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting("/prefix"),
				insecure: newSetting(true),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
		{
			name: "OptionEndpointWithEnvInsecure",
			options: []Option{
				WithEndpoint("env.endpoint:8080"),
			},
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_INSECURE": "true",
			},
			want: config{
				endpoint: newSetting("env.endpoint:8080"),
				path:     newSetting(defaultPath),
				insecure: newSetting(true),
				retryCfg: newSetting(defaultRetryCfg),
				timeout:  newSetting(defaultTimeout),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}

			var err error
			t.Cleanup(func(orig otel.ErrorHandler) func() {
				otel.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) {
					err = errors.Join(err, e)
				}))
				return func() { otel.SetErrorHandler(orig) }
			}(otel.GetErrorHandler()))
			c := newConfig(tc.options)
			if !tc.want.maxRequestSize.Set {
				tc.want.maxRequestSize = newSetting(64 * 1024 * 1024)
			}
			if !tc.want.maxResponseSize.Set {
				tc.want.maxResponseSize = newSetting(defaultMaxResponseSize)
			}

			// Do not compare pointer values.
			assertTLSConfig(t, tc.want.tlsCfg, c.tlsCfg)
			var emptyTLS setting[*tls.Config]
			c.tlsCfg, tc.want.tlsCfg = emptyTLS, emptyTLS

			// Cannot compare funcs, see TestWithProxy.
			c.proxy = setting[HTTPTransportProxyFunc]{}

			assert.Equal(t, tc.want, c)

			for _, errMsg := range tc.errs {
				assert.ErrorContains(t, err, errMsg)
			}
		})
	}
}

func TestNewConfigTLSClientEnvironment(t *testing.T) {
	for _, key := range []string{
		"OTEL_EXPORTER_OTLP_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_CLIENT_KEY",
		"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY",
	} {
		t.Setenv(key, "")
		require.NoError(t, os.Unsetenv(key))
	}

	crt, err := tls.X509KeyPair([]byte(weakCertificate), []byte(weakPrivateKey))
	require.NoError(t, err)
	explicitTLS := &tls.Config{ServerName: "explicit"}

	for _, tc := range []struct {
		name    string
		envars  map[string]string
		options []Option
		reads   []string
		err     string
	}{
		{name: "Unset"},
		{
			name: "GenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "generic.key",
			},
			reads: []string{"generic.pem", "generic.key"},
		},
		{
			name: "LogsPairPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE":      "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "logs.pem",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "logs.key",
			},
			reads: []string{"logs.pem", "logs.key"},
		},
		{
			name: "LogsCertificateWithGenericKey",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "logs.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
			},
			reads: []string{"logs.pem", "generic.key"},
		},
		{
			name: "GenericCertificateWithLogsKey",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "generic.pem",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":    "logs.key",
			},
			reads: []string{"generic.pem", "logs.key"},
		},
		{
			name: "LogsCertificateOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE":      "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "logs.pem",
			},
			reads: []string{"logs.pem", "generic.key"},
		},
		{
			name: "LogsKeyOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":    "logs.key",
			},
			reads: []string{"generic.pem", "logs.key"},
		},
		{
			name: "EmptyLogsCertificateFallsBack",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE":      "generic.pem",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "logs.key",
			},
			reads: []string{"generic.pem", "logs.key"},
		},
		{
			name: "EmptyLogsKeyFallsBack",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "logs.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "",
			},
			reads: []string{"logs.pem", "generic.key"},
		},
		{
			name: "CertificateWithoutKey",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "logs.pem",
			},
		},
		{
			name: "KeyWithoutCertificate",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY": "logs.key",
			},
		},
		{
			name: "InvalidLogsCertificateOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE":      "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "invalid.pem",
			},
			reads: []string{"invalid.pem", "generic.key"},
			err:   "tls: failed to find any PEM data in certificate input",
		},
		{
			name: "InvalidLogsKeyOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":    "invalid.key",
			},
			reads: []string{"generic.pem", "invalid.key"},
			err:   "tls: failed to find any PEM data in key input",
		},
		{
			name: "MismatchedLogsKeyOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE": "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":         "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":    "mismatched.key",
			},
			reads: []string{"generic.pem", "mismatched.key"},
			err:   "tls: private key does not match public key",
		},
		{
			name: "MissingLogsCertificateOverridesGenericPair",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE":      "generic.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "generic.key",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "missing.pem",
			},
			reads: []string{"missing.pem"},
			err:   "file does not exist: missing.pem",
		},
		{
			name: "OptionPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "missing.pem",
				"OTEL_EXPORTER_OTLP_CLIENT_KEY":              "invalid.key",
			},
			options: []Option{WithTLSClientConfig(explicitTLS)},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}

			origReadFile := readFile
			t.Cleanup(func() { readFile = origReadFile })
			var reads []string
			readFile = func(name string) ([]byte, error) {
				reads = append(reads, name)
				switch name {
				case "generic.pem", "logs.pem":
					return []byte(weakCertificate), nil
				case "generic.key", "logs.key":
					return []byte(weakPrivateKey), nil
				case "mismatched.key":
					return []byte(mismatchedPrivateKey), nil
				case "invalid.pem", "invalid.key":
					return []byte("invalid PEM"), nil
				default:
					return nil, fmt.Errorf("file does not exist: %s", name)
				}
			}

			origHandler := otel.GetErrorHandler()
			t.Cleanup(func() { otel.SetErrorHandler(origHandler) })
			var gotErr error
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) {
				gotErr = errors.Join(gotErr, e)
			}))

			cfg := newConfig(tc.options)
			assert.Equal(t, tc.reads, reads)
			var want setting[*tls.Config]
			switch {
			case len(tc.options) != 0:
				want = newSetting(explicitTLS)
			case tc.err != "":
				assert.ErrorContains(t, gotErr, "failed to load TLS:")
				assert.ErrorContains(t, gotErr, tc.err)
			case len(tc.reads) != 0:
				want = newSetting(&tls.Config{Certificates: []tls.Certificate{crt}})
			}
			if tc.err == "" {
				assert.NoError(t, gotErr)
			}
			assertTLSConfig(t, want, cfg.tlsCfg)
			if len(tc.options) != 0 {
				assert.Equal(t, explicitTLS.ServerName, cfg.tlsCfg.Value.ServerName)
			}
		})
	}
}

func TestMaxResponseSize(t *testing.T) {
	for _, size := range []int64{0, -1} {
		cfg := newConfig([]Option{WithMaxResponseSize(size)})
		assert.Equal(t, defaultMaxResponseSize, cfg.maxResponseSize.Value)
	}

	cfg := newConfig([]Option{WithMaxResponseSize(2), WithMaxResponseSize(0)})
	assert.Equal(t, int64(2), cfg.maxResponseSize.Value)
}

func assertTLSConfig(t *testing.T, want, got setting[*tls.Config]) {
	t.Helper()

	assert.Equal(t, want.Set, got.Set, "setting Set")
	if !want.Set {
		return
	}

	if want.Value == nil {
		assert.Nil(t, got.Value, "*tls.Config")
		return
	}
	require.NotNil(t, got.Value, "*tls.Config")

	if want.Value.RootCAs == nil {
		assert.Nil(t, got.Value.RootCAs, "*tls.Config.RootCAs")
	} else if assert.NotNil(t, got.Value.RootCAs, "RootCAs") {
		assert.True(t, want.Value.RootCAs.Equal(got.Value.RootCAs), "RootCAs equal")
	}
	assert.Equal(t, want.Value.Certificates, got.Value.Certificates, "Certificates")
}

func TestWithProxy(t *testing.T) {
	proxy := func(*http.Request) (*url.URL, error) { return nil, nil }
	opts := []Option{WithProxy(HTTPTransportProxyFunc(proxy))}
	c := newConfig(opts)

	assert.True(t, c.proxy.Set)
	assert.NotNil(t, c.proxy.Value)
}

func TestConvHeaders(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "simple test",
			value:   "userId=alice",
			want:    map[string]string{"userId": "alice"},
			wantErr: false,
		},
		{
			name:    "simple test with spaces",
			value:   " userId = alice  ",
			want:    map[string]string{"userId": "alice"},
			wantErr: false,
		},
		{
			name:    "simple header conforms to RFC 3986 spec",
			value:   " userId = alice+test ",
			want:    map[string]string{"userId": "alice+test"},
			wantErr: false,
		},
		{
			name:  "multiple headers encoded",
			value: "userId=alice,serverNode=DF%3A28,isProduction=false",
			want: map[string]string{
				"userId":       "alice",
				"serverNode":   "DF:28",
				"isProduction": "false",
			},
			wantErr: false,
		},
		{
			name:  "multiple headers encoded per RFC 3986 spec",
			value: "userId=alice+test,serverNode=DF%3A28,isProduction=false,namespace=localhost/test",
			want: map[string]string{
				"userId":       "alice+test",
				"serverNode":   "DF:28",
				"isProduction": "false",
				"namespace":    "localhost/test",
			},
			wantErr: false,
		},
		{
			name:    "invalid headers format",
			value:   "userId:alice",
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name:  "invalid key",
			value: "%XX=missing,userId=alice",
			want: map[string]string{
				"%XX":    "missing",
				"userId": "alice",
			},
			wantErr: false,
		},
		{
			name:  "invalid value",
			value: "missing=%XX,userId=alice",
			want: map[string]string{
				"userId": "alice",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyValues, err := convHeaders(tt.value)
			assert.Equal(t, tt.want, keyValues)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "expected no error but got one")
			}
		})
	}
}
