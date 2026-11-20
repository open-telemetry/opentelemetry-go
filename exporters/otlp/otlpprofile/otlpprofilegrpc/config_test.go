// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel"
)

func TestConfigDefaultsAndEnvVars(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := newConfig(nil)
		assert.Equal(t, "localhost:4317", cfg.endpoint.Value)
		assert.False(t, cfg.insecure.Value)
		assert.Equal(t, noCompression, cfg.compression.Value)
		assert.Equal(t, 10*time.Second, cfg.timeout.Value)
		assert.Equal(t, 64*1024*1024, cfg.maxRequestSize.Value)
	})

	t.Run("general_and_signal_env_precedence", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://env.endpoint:4317")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "https://overrode.by.signal.specific:9443")
		t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "a=1")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_HEADERS", "b=2,c=v%203")
		t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "none")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_COMPRESSION", "GZIP")
		t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "1000")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", "4500")

		cfg := newConfig(nil)
		assert.Equal(t, "overrode.by.signal.specific:9443", cfg.endpoint.Value)
		assert.False(t, cfg.insecure.Value)
		assert.Equal(t, gzipCompression, cfg.compression.Value)
		assert.Equal(t, 4500*time.Millisecond, cfg.timeout.Value)
		assert.Equal(t, map[string]string{"b": "2", "c": "v 3"}, cfg.headers.Value)
	})

	t.Run("insecure_env_when_no_endpoint_scheme", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "false")
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_INSECURE", "true")
		cfg := newConfig(nil)
		assert.True(t, cfg.insecure.Value)
	})

	t.Run("options_override_env", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT", "https://env.endpoint:4317")
		creds := insecure.NewCredentials()
		cfg := newConfig([]Option{
			WithEndpointURL("http://someendpoint:5317"),
			WithCompressor("gzip"),
			WithHeaders(map[string]string{"x-key": "val"}),
			WithTimeout(2 * time.Second),
			WithMaxRequestSize(4096),
			WithReconnectionPeriod(3 * time.Second),
			WithServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			WithTLSCredentials(creds),
			WithDialOption(grpc.WithAuthority("override.example")),
		})
		assert.Equal(t, "someendpoint:5317", cfg.endpoint.Value)
		assert.True(t, cfg.insecure.Value)
		assert.Equal(t, gzipCompression, cfg.compression.Value)
		assert.Equal(t, map[string]string{"x-key": "val"}, cfg.headers.Value)
		assert.Equal(t, 2*time.Second, cfg.timeout.Value)
		assert.Equal(t, 4096, cfg.maxRequestSize.Value)
		assert.Equal(t, 3*time.Second, cfg.reconnectionPeriod.Value)
		assert.JSONEq(t, `{"loadBalancingPolicy":"round_robin"}`, cfg.serviceConfig.Value)
		assert.Equal(t, creds, cfg.gRPCCredentials.Value)
		assert.Len(t, cfg.dialOptions.Value, 1)

		dialOpts := newGRPCDialOptions(cfg)
		assert.NotEmpty(t, dialOpts)
	})

	t.Run("invalid_env_values_fall_back", func(t *testing.T) {
		errs := captureOTelErrors(t)
		t.Setenv("OTEL_EXPORTER_OTLP_COMPRESSION", "brotli")
		t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "notanint")
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "notabool")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "://bad")
		t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "invalid,bad key=1,=empty,k=%ZZ,ok=val")

		cfg := newConfig(nil)
		assert.Equal(t, "localhost:4317", cfg.endpoint.Value)
		assert.Equal(t, noCompression, cfg.compression.Value)
		assert.Equal(t, 10*time.Second, cfg.timeout.Value)
		assert.False(t, cfg.insecure.Value)
		assert.Nil(t, cfg.headers.Value)

		require.Len(t, *errs, 5)
		assert.ErrorContains(t, (*errs)[0], "OTEL_EXPORTER_OTLP_ENDPOINT")
		assert.ErrorContains(t, (*errs)[1], "OTEL_EXPORTER_OTLP_INSECURE")
		assert.ErrorContains(t, (*errs)[2], "OTEL_EXPORTER_OTLP_HEADERS")
		assert.ErrorContains(t, (*errs)[3], "OTEL_EXPORTER_OTLP_COMPRESSION")
		assert.ErrorContains(t, (*errs)[4], "OTEL_EXPORTER_OTLP_TIMEOUT")
	})

	t.Run("invalid_option_values_fall_back", func(t *testing.T) {
		errs := captureOTelErrors(t)
		cfg := newConfig([]Option{
			nil,
			WithEndpointURL("://bad"),
			WithCompressor("unknown"),
		})
		assert.Equal(t, "localhost:4317", cfg.endpoint.Value)
		assert.Equal(t, noCompression, cfg.compression.Value)
		require.Len(t, *errs, 1)
		assert.ErrorContains(t, (*errs)[0], "unknown compression: unknown")
	})

	t.Run("compression_case_insensitive", func(t *testing.T) {
		for _, tc := range []struct {
			in   string
			want compression
		}{
			{in: "none", want: noCompression},
			{in: "NONE", want: noCompression},
			{in: "", want: noCompression},
			{in: "gzip", want: gzipCompression},
			{in: "GZIP", want: gzipCompression},
		} {
			got, err := convCompression(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		}
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
		assert.NotEmpty(t, newGRPCDialOptions(cfg))
	})

	t.Run("env_tls_errors_ignored_gracefully", func(t *testing.T) {
		t.Run("invalid_pem", func(t *testing.T) {
			errs := captureOTelErrors(t)
			badFile := filepath.Join(t.TempDir(), "bad.pem")
			require.NoError(t, os.WriteFile(badFile, []byte("not-a-pem"), 0o600))
			t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", badFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", badFile)
			t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", badFile)

			cfg := newConfig(nil)
			assert.Nil(t, cfg.tlsCfg.Value)
			require.Len(t, *errs, 1)
			assert.ErrorContains(t, (*errs)[0], "failed to load TLS")
		})
		t.Run("missing_files", func(t *testing.T) {
			missing := filepath.Join(t.TempDir(), "does-not-exist.pem")
			certPath, _ := writeTestCertAndKey(t)

			t.Run("missing_ca", func(t *testing.T) {
				errs := captureOTelErrors(t)
				t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", missing)
				cfg := newConfig(nil)
				assert.Nil(t, cfg.tlsCfg.Value)
				require.Len(t, *errs, 1)
			})
			t.Run("missing_client_cert", func(t *testing.T) {
				errs := captureOTelErrors(t)
				t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", missing)
				t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", certPath)
				cfg := newConfig(nil)
				assert.Nil(t, cfg.tlsCfg.Value)
				require.Len(t, *errs, 1)
			})
			t.Run("missing_client_key", func(t *testing.T) {
				errs := captureOTelErrors(t)
				t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", certPath)
				t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", missing)
				cfg := newConfig(nil)
				assert.Nil(t, cfg.tlsCfg.Value)
				require.Len(t, *errs, 1)
			})
		})
	})

	t.Run("env_tls_skipped_when_insecure_or_explicit_creds", func(t *testing.T) {
		certPath, keyPath := writeTestCertAndKey(t)
		t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_CLIENT_KEY", keyPath)

		t.Run("with_insecure", func(t *testing.T) {
			cfg := newConfig([]Option{WithInsecure()})
			assert.True(t, cfg.insecure.Value)
			assert.Nil(t, cfg.tlsCfg.Value)
		})
		t.Run("with_http_endpoint_url", func(t *testing.T) {
			cfg := newConfig([]Option{WithEndpointURL("http://localhost:4317")})
			assert.True(t, cfg.insecure.Value)
			assert.Nil(t, cfg.tlsCfg.Value)
		})
		t.Run("with_tls_credentials", func(t *testing.T) {
			creds := credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})
			cfg := newConfig([]Option{WithTLSCredentials(creds)})
			assert.Equal(t, creds, cfg.gRPCCredentials.Value)
			assert.Nil(t, cfg.tlsCfg.Value)
		})
	})

	t.Run("endpoint_option_ignores_env_endpoint_scheme", func(t *testing.T) {
		t.Run("with_endpoint_not_downgraded_by_http_env", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
			cfg := newConfig([]Option{WithEndpoint("prod.example.com:4317")})
			assert.Equal(t, "prod.example.com:4317", cfg.endpoint.Value)
			assert.False(t, cfg.insecure.Value)
		})
		t.Run("schemeless_endpoint_url_not_downgraded_by_http_env", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
			cfg := newConfig([]Option{WithEndpointURL("dns:///prod.example.com:4317")})
			assert.Equal(t, "dns:///prod.example.com:4317", cfg.endpoint.Value)
			assert.False(t, cfg.insecure.Value)
		})
		t.Run("with_endpoint_still_honors_env_insecure_over_https_env_endpoint", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://remote.example.com:4317")
			t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_INSECURE", "true")
			cfg := newConfig([]Option{WithEndpoint("localhost:4317")})
			assert.Equal(t, "localhost:4317", cfg.endpoint.Value)
			assert.True(t, cfg.insecure.Value)
		})
	})

	t.Run("header_env_errors_redact_values", func(t *testing.T) {
		errs := captureOTelErrors(t)

		const secret = "super-secret-bearer-token"
		t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "Authorization=Bearer "+secret+",badentry,X-Esc=%ZZ"+secret)

		cfg := newConfig(nil)
		assert.Nil(t, cfg.headers.Value)
		require.Len(t, *errs, 1)
		assert.NotContains(t, (*errs)[0].Error(), secret)
	})

	t.Run("timeout_env_bounds", func(t *testing.T) {
		t.Run("zero_allowed", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", "0")
			cfg := newConfig(nil)
			assert.Equal(t, time.Duration(0), cfg.timeout.Value)
		})
		t.Run("max_valid_ms", func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", "9223372036854")
			cfg := newConfig(nil)
			assert.Equal(t, 9223372036854*time.Millisecond, cfg.timeout.Value)
		})
		for _, bad := range []string{"-1", "9223372036855", "9223372036854775807"} {
			t.Run("reject_"+bad, func(t *testing.T) {
				errs := captureOTelErrors(t)
				t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", "3000")
				t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT", bad)
				cfg := newConfig(nil)
				assert.Equal(t, 3*time.Second, cfg.timeout.Value)
				require.Len(t, *errs, 1)
			})
			t.Run("reject_generic_"+bad, func(t *testing.T) {
				errs := captureOTelErrors(t)
				t.Setenv("OTEL_EXPORTER_OTLP_TIMEOUT", bad)
				cfg := newConfig(nil)
				assert.Equal(t, defaultTimeout, cfg.timeout.Value)
				require.Len(t, *errs, 1)
			})
		}
	})
}

func captureOTelErrors(t *testing.T) *[]error {
	t.Helper()
	var errs []error
	origHandler := otel.GetErrorHandler()
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errs = append(errs, err)
	}))
	t.Cleanup(func() { otel.SetErrorHandler(origHandler) })
	return &errs
}

func TestClientHelpers(t *testing.T) {
	noop := newNoopClient()
	assert.NoError(t, noop.UploadProfiles(t.Context(), nil))
	assert.NoError(t, noop.Shutdown(t.Context()))

	c := &client{
		metadata: metadata.Pairs("x-custom", "1", profilesDevelopmentVersionKey, profilesDevelopmentVersion),
	}
	parent := metadata.NewOutgoingContext(t.Context(), metadata.Pairs(
		"x-parent", "2",
		profilesDevelopmentVersionKey, "999",
	))
	ctx, cancel := c.exportContext(parent)
	defer cancel()

	outMD, ok := metadata.FromOutgoingContext(ctx)
	require.True(t, ok)
	assert.Equal(t, []string{"2"}, outMD.Get("x-parent"))
	assert.Equal(t, []string{"1"}, outMD.Get("x-custom"))
	assert.Equal(t, []string{profilesDevelopmentVersion}, outMD.Get(profilesDevelopmentVersionKey))

	cfg := newConfig([]Option{WithTLSCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13}))})
	assert.NotEmpty(t, newGRPCDialOptions(cfg))
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
