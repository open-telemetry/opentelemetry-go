// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
)

func TestConfigDefaultsAndEnvVars(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := newConfig(nil)
		assert.Equal(t, "localhost:4318", cfg.endpoint.Value)
		assert.Equal(t, "/v1development/profiles", cfg.path.Value.Path)
		assert.False(t, cfg.insecure.Value)
		assert.Equal(t, NoCompression, cfg.compression.Value)
		assert.Equal(t, 10*time.Second, cfg.timeout.Value)
		assert.Equal(t, 64*1024*1024, cfg.maxRequestSize.Value)
		assert.Equal(t, int64(4*1024*1024), cfg.maxResponseSize.Value)
	})

	t.Run("general_otlp_env_appends_signal_path", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://env.endpoint:4318/base")
		t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "k1=v1,k2=v%202")
		t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "GZIP")
		t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "2500")

		cfg := newConfig(nil)
		assert.Equal(t, "env.endpoint:4318", cfg.endpoint.Value)
		assert.Equal(t, "/base/v1development/profiles", cfg.path.Value.Path)
		assert.True(t, cfg.insecure.Value)
		assert.Equal(t, GzipCompression, cfg.compression.Value)
		assert.Equal(t, 2500*time.Millisecond, cfg.timeout.Value)
		assert.Equal(t, map[string]string{"k1": "v1", "k2": "v 2"}, cfg.headers.Value)
	})

	t.Run("signal_specific_env_takes_precedence", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://env.endpoint:4318")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "https://overrode.by.signal.specific:9443/custom/profiles")
		t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "a=1")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_HEADERS", "b=2")
		t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "none")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_COMPRESSION", "gzip")
		t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "1000")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", "5000")

		cfg := newConfig(nil)
		assert.Equal(t, "overrode.by.signal.specific:9443", cfg.endpoint.Value)
		assert.Equal(t, "/custom/profiles", cfg.path.Value.Path)
		assert.False(t, cfg.insecure.Value)
		assert.Equal(t, GzipCompression, cfg.compression.Value)
		assert.Equal(t, 5*time.Second, cfg.timeout.Value)
		assert.Equal(t, map[string]string{"b": "2"}, cfg.headers.Value)
	})

	t.Run("insecure_env_when_no_endpoint_scheme", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "false")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_INSECURE", "true")
		cfg := newConfig(nil)
		assert.True(t, cfg.insecure.Value)
	})

	t.Run("options_override_env", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "https://env.endpoint:4318/env")
		customClient := &http.Client{Timeout: 5 * time.Second}
		proxyFn := func(*http.Request) (*url.URL, error) { return nil, nil }
		cfg := newConfig([]Option{
			WithEndpointURL("http://someendpoint:8080/opt"),
			WithURLPath("/custom/path"),
			WithCompression(GzipCompression),
			WithHeaders(map[string]string{"x-test": "ok"}),
			WithTimeout(3 * time.Second),
			WithMaxRequestSize(1024),
			WithMaxResponseSize(2048),
			WithProxy(proxyFn),
			WithHTTPClient(customClient),
		})
		assert.Equal(t, "someendpoint:8080", cfg.endpoint.Value)
		assert.Equal(t, "/custom/path", cfg.path.Value.Path)
		assert.True(t, cfg.insecure.Value)
		assert.Equal(t, GzipCompression, cfg.compression.Value)
		assert.Equal(t, map[string]string{"x-test": "ok"}, cfg.headers.Value)
		assert.Equal(t, 3*time.Second, cfg.timeout.Value)
		assert.Equal(t, 1024, cfg.maxRequestSize.Value)
		assert.Equal(t, int64(2048), cfg.maxResponseSize.Value)
		assert.NotNil(t, cfg.proxy.Value)
		assert.Same(t, customClient, cfg.httpClient)
	})

	t.Run("invalid_env_and_option_values_fall_back", func(t *testing.T) {
		var mu sync.Mutex
		var handled []error
		orig := otel.GetErrorHandler()
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			mu.Lock()
			handled = append(handled, err)
			mu.Unlock()
		}))
		t.Cleanup(func() { otel.SetErrorHandler(orig) })

		t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "brotli")
		t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "notanint")
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "notabool")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "://bad")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "://also-bad")
		t.Setenv(
			"OTEL_EXPORTER_OTLP_HEADERS",
			"Authorization=Bearer top-secret-token,invalid,bad key=1,=emptykey,k=%ZZ",
		)

		cfg := newConfig([]Option{
			nil,
			WithEndpoint("kept.example:4318"),
			WithEndpointURL("://bad"),
			WithMaxResponseSize(2048),
			WithMaxResponseSize(0),
			WithMaxResponseSize(-1),
		})
		assert.Equal(t, "kept.example:4318", cfg.endpoint.Value)
		assert.Equal(t, NoCompression, cfg.compression.Value)
		assert.Equal(t, 10*time.Second, cfg.timeout.Value)
		assert.False(t, cfg.insecure.Value)
		assert.Nil(t, cfg.headers.Value)
		assert.Equal(t, int64(2048), cfg.maxResponseSize.Value)

		mu.Lock()
		defer mu.Unlock()
		assert.NotEmpty(t, handled)
		for _, err := range handled {
			assert.NotContains(t, err.Error(), "top-secret-token")
		}
	})

	t.Run("timeout_env_bounds_and_fallback", func(t *testing.T) {
		for _, tc := range []struct {
			name    string
			general string
			signal  string
			want    time.Duration
		}{
			{name: "zero", general: "0", want: 0},
			{name: "negative_uses_default", general: "-1", want: defaultTimeout},
			{name: "max_valid", general: "9223372036854", want: 9223372036854 * time.Millisecond},
			{name: "first_overflow", general: "9223372036855", want: defaultTimeout},
			{name: "max_int64_overflow", general: "9223372036854775807", want: defaultTimeout},
			{name: "invalid_signal_keeps_general", general: "2500", signal: "-100", want: 2500 * time.Millisecond},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", tc.general)
				t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", tc.signal)
				cfg := newConfig(nil)
				assert.Equal(t, tc.want, cfg.timeout.Value)
			})
		}
	})

	t.Run("endpoint_paths_compression_and_scheme_precedence", func(t *testing.T) {
		t.Run("signal_endpoint_without_path_uses_root", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "https://collector:4318")
			t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_INSECURE", "true")
			t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "none")

			cfg := newConfig(nil)
			assert.Equal(t, "collector:4318", cfg.endpoint.Value)
			assert.Equal(t, "/", cfg.path.Value.Path)
			assert.False(t, cfg.insecure.Value)
			assert.Equal(t, NoCompression, cfg.compression.Value)
		})

		t.Run("general_endpoint_trailing_slash_and_raw_path", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://collector:4318/a%2Fb/")
			t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "false")

			cfg := newConfig(nil)
			assert.Equal(t, "/a/b/v1development/profiles", cfg.path.Value.Path)
			assert.Equal(t, "/a%2Fb/v1development/profiles", cfg.path.Value.RawPath)
			assert.True(t, cfg.insecure.Value)
		})

		t.Run("schemeless_endpoint_preserves_insecure_env", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "//collector:4318/base")
			t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")

			cfg := newConfig(nil)
			assert.Equal(t, "collector:4318", cfg.endpoint.Value)
			assert.True(t, cfg.insecure.Value)
		})
	})

	t.Run("env_tls_certificates", func(t *testing.T) {
		certPath, keyPath := writeTestCertAndKey(t)
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CLIENT_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CLIENT_KEY", keyPath)

		cfg := newConfig(nil)
		require.NotNil(t, cfg.tlsCfg.Value)
		assert.NotNil(t, cfg.tlsCfg.Value.RootCAs)
		assert.Len(t, cfg.tlsCfg.Value.Certificates, 1)

		// Explicit WithTLSClientConfig(nil) suppresses env TLS config and allows WithInsecure.
		exp, err := New(t.Context(), WithInsecure(), WithTLSClientConfig(nil))
		require.NoError(t, err)
		require.NoError(t, exp.Shutdown(t.Context()))
	})

	t.Run("env_tls_errors_ignored_gracefully", func(t *testing.T) {
		badFile := filepath.Join(t.TempDir(), "bad.pem")
		require.NoError(t, os.WriteFile(badFile, []byte("not-a-pem"), 0o600))
		missingFile := filepath.Join(t.TempDir(), "missing.pem")

		t.Run("invalid_pem", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", badFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", badFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", badFile)
			assert.Nil(t, newConfig(nil).tlsCfg.Value)
		})

		t.Run("missing_files", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", missingFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", missingFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", badFile)
			assert.Nil(t, newConfig(nil).tlsCfg.Value)
		})

		t.Run("missing_key_file", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", badFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", missingFile)
			assert.Nil(t, newConfig(nil).tlsCfg.Value)
		})
	})

	t.Run("insecure_with_tls_rejected", func(t *testing.T) {
		_, err := New(t.Context(), WithInsecure(), WithTLSClientConfig(&tls.Config{MinVersion: tls.VersionTLS12}))
		require.ErrorIs(t, err, errInsecureEndpointWithTLS)
	})
}

func TestClientHelpers(t *testing.T) {
	rc, err := bodyReaderErr([]byte("abc"))()
	require.NoError(t, err)
	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, "abc", string(data))

	assert.True(t, isProtobufContentType("application/x-protobuf; charset=utf-8"))
	assert.False(t, isProtobufContentType("text/plain"))
	assert.False(t, isProtobufContentType(";;;bad"))

	assert.Equal(t, 5*time.Second, retryAfterDuration("5"))
	assert.Equal(t, maxRetryAfter, retryAfterDuration("9223372036854775807"))
	assert.Equal(t, time.Duration(0), retryAfterDuration("-5"))
	assert.Equal(t, time.Duration(0), retryAfterDuration("not-a-date"))
	past := time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)
	assert.Equal(t, time.Duration(0), retryAfterDuration(past))
	future := time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)
	d := retryAfterDuration(future)
	assert.Greater(t, d, 59*time.Minute)
	assert.LessOrEqual(t, d, time.Hour)

	assert.NoError(t, newNoopClient().UploadProfiles(t.Context(), nil))

	wrapped := errors.New("boom")
	rErr := retryableError{throttle: time.Second, err: wrapped}
	assert.Contains(t, rErr.Error(), "boom")
	assert.ErrorIs(t, rErr, wrapped)
	assert.Equal(t, "retry-able request failure", (retryableError{}).Error())

	ok, d := evaluate(nil)
	assert.False(t, ok)
	assert.Zero(t, d)
	ok, d = evaluate(rErr)
	assert.True(t, ok)
	assert.Equal(t, time.Second, d)
}

func writeTestCertAndKey(t *testing.T) (certPath, keyPath string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IsCA:         true,
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	require.NoError(t, err)
	keyDER, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)

	dir := t.TempDir()
	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	require.NoError(t, os.WriteFile(certPath, certPEM, 0o600))
	require.NoError(t, os.WriteFile(keyPath, keyPEM, 0o600))
	return certPath, keyPath
}
