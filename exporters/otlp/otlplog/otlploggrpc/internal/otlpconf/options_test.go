// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpconf // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/otlpconf"

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
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
)

// This is only for testing, as the package that is using this utility has its own newConfig function.
func newConfig(options []Option) Config {
	var c Config
	for _, opt := range options {
		c = opt.ApplyOption(c)
	}

	c = LoadConfig(c)

	return c
}

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

	testcases := []struct {
		name    string
		options []Option
		envars  map[string]string
		want    Config
		errs    []string
	}{
		{
			name: "Defaults",
			want: Config{
				Endpoint: NewSetting(defaultEndpoint),
				Path:     NewSetting(defaultPath),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
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
			},
			want: Config{
				Endpoint:    NewSetting("test"),
				Path:        NewSetting("/path"),
				Insecure:    NewSetting(true),
				TLSCfg:      NewSetting(tlsCfg),
				Headers:     NewSetting(headers),
				Compression: NewSetting(GzipCompression),
				Timeout:     NewSetting(time.Second),
				RetryCfg:    NewSetting(rc),
			},
		},
		{
			name: "WithEndpointURL",
			options: []Option{
				WithEndpointURL("http://test:8080/path"),
			},
			want: Config{
				Endpoint: NewSetting("test:8080"),
				Path:     NewSetting("/path"),
				Insecure: NewSetting(true),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
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
			want: Config{
				Endpoint: NewSetting("not-test:9090"),
				Path:     NewSetting("/alt"),
				Insecure: NewSetting(true),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
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
			want: Config{
				Endpoint: NewSetting("test:8080"),
				Path:     NewSetting("/path"),
				Insecure: NewSetting(false),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
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
			want: Config{
				Endpoint:    NewSetting("env.endpoint:8080"),
				Path:        NewSetting("/prefix"),
				Insecure:    NewSetting(false),
				TLSCfg:      NewSetting(tlsCfg),
				Headers:     NewSetting(headers),
				Compression: NewSetting(GzipCompression),
				Timeout:     NewSetting(15 * time.Second),
				RetryCfg:    NewSetting(defaultRetryCfg),
			},
		},
		{
			name: "LogEnpointEnvironmentVariablesDefaultPath",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "http://env.endpoint",
			},
			want: Config{
				Endpoint: NewSetting("env.endpoint"),
				Path:     NewSetting("/"),
				Insecure: NewSetting(true),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
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
			want: Config{
				Endpoint:    NewSetting("env.endpoint:8080"),
				Path:        NewSetting("/prefix/v1/logs"),
				Insecure:    NewSetting(true),
				TLSCfg:      NewSetting(tlsCfg),
				Headers:     NewSetting(headers),
				Compression: NewSetting(NoCompression),
				Timeout:     NewSetting(15 * time.Second),
				RetryCfg:    NewSetting(defaultRetryCfg),
			},
		},
		{
			name: "OTLPEnpointEnvironmentVariablesDefaultPath",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://env.endpoint",
			},
			want: Config{
				Endpoint: NewSetting("env.endpoint"),
				Path:     NewSetting(defaultPath),
				Insecure: NewSetting(true),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
			},
		},
		{
			name: "EnvironmentVariablesPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":           "http://ignored:9090/alt",
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
			want: Config{
				Endpoint:    NewSetting("env.endpoint:8080"),
				Path:        NewSetting("/path"),
				Insecure:    NewSetting(false),
				TLSCfg:      NewSetting(tlsCfg),
				Headers:     NewSetting(headers),
				Compression: NewSetting(GzipCompression),
				Timeout:     NewSetting(15 * time.Second),
				RetryCfg:    NewSetting(defaultRetryCfg),
			},
		},
		{
			name: "OptionsPrecedence",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT":           "http://ignored:9090/alt",
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
				WithURLPath("/path"),
				WithInsecure(),
				WithTLSClientConfig(tlsCfg),
				WithCompression(GzipCompression),
				WithHeaders(headers),
				WithTimeout(time.Second),
				WithRetry(rc),
			},
			want: Config{
				Endpoint:    NewSetting("test"),
				Path:        NewSetting("/path"),
				Insecure:    NewSetting(true),
				TLSCfg:      NewSetting(tlsCfg),
				Headers:     NewSetting(headers),
				Compression: NewSetting(GzipCompression),
				Timeout:     NewSetting(time.Second),
				RetryCfg:    NewSetting(rc),
			},
		},
		{
			name: "InvalidEnvironmentVariables",
			envars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":           "%invalid",
				"OTEL_EXPORTER_OTLP_LOGS_HEADERS":            "a,%ZZ=valid,key=%ZZ",
				"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION":        "xz",
				"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT":            "100 seconds",
				"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE":        "invalid_cert",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE": "invalid_cert",
				"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY":         "invalid_key",
			},
			want: Config{
				Endpoint: NewSetting(defaultEndpoint),
				Path:     NewSetting(defaultPath),
				Timeout:  NewSetting(defaultTimeout),
				RetryCfg: NewSetting(defaultRetryCfg),
			},
			errs: []string{
				`invalid OTEL_EXPORTER_OTLP_LOGS_ENDPOINT value %invalid: parse "%invalid": invalid URL escape "%in"`,
				`failed to load TLS:`,
				`certificate not added`,
				`tls: failed to find any PEM data in certificate input`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_HEADERS value a,%ZZ=valid,key=%ZZ:`,
				`invalid header: a`,
				`invalid header key: %ZZ`,
				`invalid header value: %ZZ`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_COMPRESSION value xz: unknown compression: xz`,
				`invalid OTEL_EXPORTER_OTLP_LOGS_TIMEOUT value 100 seconds: strconv.Atoi: parsing "100 seconds": invalid syntax`,
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

			// Do not compare pointer values.
			assertTLSConfig(t, tc.want.TLSCfg, c.TLSCfg)
			var emptyTLS Setting[*tls.Config]
			c.TLSCfg, tc.want.TLSCfg = emptyTLS, emptyTLS

			assert.Equal(t, tc.want, c)

			for _, errMsg := range tc.errs {
				assert.ErrorContains(t, err, errMsg)
			}
		})
	}
}

func assertTLSConfig(t *testing.T, want, got Setting[*tls.Config]) {
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
	} else {
		if assert.NotNil(t, got.Value.RootCAs, "RootCAs") {
			assert.True(t, want.Value.RootCAs.Equal(got.Value.RootCAs), "RootCAs equal")
		}
	}
	assert.Equal(t, want.Value.Certificates, got.Value.Certificates, "Certificates")
}
