// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colprofilepb "go.opentelemetry.io/proto/slim/otlp/collector/profiles/v1development"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal"
	"go.opentelemetry.io/otel/sdk/profile"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type errCloseBody struct {
	io.Reader
	closeErr error
}

func (b errCloseBody) Close() error { return b.closeErr }

type tempErr string

func (e tempErr) Error() string { return string(e) }
func (tempErr) Temporary() bool { return true }
func (tempErr) Timeout() bool   { return false }

func TestClientTLSAndProxyOptions(t *testing.T) {
	t.Run("https_with_tls_client_config", func(t *testing.T) {
		var called atomic.Bool
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			called.Store(true)
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		pool := x509.NewCertPool()
		pool.AddCert(srv.Certificate())
		u, err := url.Parse(srv.URL)
		require.NoError(t, err)

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpoint(u.Host),
			otlpprofilehttp.WithTLSClientConfig(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.True(t, called.Load())
	})

	t.Run("custom_proxy_invoked", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		var proxyCalled atomic.Bool
		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithProxy(func(*http.Request) (*url.URL, error) {
				proxyCalled.Store(true)
				return nil, nil
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.True(t, proxyCalled.Load())
	})

	t.Run("with_proxy_nil_disables_env_proxy", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		// Point HTTP_PROXY at an unreachable address; WithProxy(nil) must bypass it.
		t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
		t.Setenv("NO_PROXY", "")

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithProxy(nil),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
	})

	t.Run("with_http_client_precedence", func(t *testing.T) {
		var called atomic.Bool
		customClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				// Sleep longer than WithTimeout(10ms) to verify WithHTTPClient overrides it.
				time.Sleep(25 * time.Millisecond)
				called.Store(true)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader(nil)),
				}, nil
			}),
		}

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpoint("localhost:4318"),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithTimeout(10*time.Millisecond),
			otlpprofilehttp.WithProxy(func(*http.Request) (*url.URL, error) {
				return nil, errors.New("proxy should not be called when WithHTTPClient is set")
			}),
			otlpprofilehttp.WithHTTPClient(customClient),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.True(t, called.Load())
	})
}

func TestClientResponsesAndLimits(t *testing.T) {
	t.Run("gzip_response_and_partial_success_variants", func(t *testing.T) {
		for _, tc := range []struct {
			name    string
			ps      *colprofilepb.ExportProfilesPartialSuccess
			gzipped bool
			wantErr bool
		}{
			{name: "empty_partial_success_ok", ps: &colprofilepb.ExportProfilesPartialSuccess{}, wantErr: false},
			{name: "message_only", ps: &colprofilepb.ExportProfilesPartialSuccess{ErrorMessage: "warn"}, wantErr: true},
			{name: "rejected_only_gzipped", ps: &colprofilepb.ExportProfilesPartialSuccess{RejectedProfiles: 1}, gzipped: true, wantErr: true},
		} {
			t.Run(tc.name, func(t *testing.T) {
				raw, err := proto.Marshal(&colprofilepb.ExportProfilesServiceResponse{PartialSuccess: tc.ps})
				require.NoError(t, err)
				payload := raw
				if tc.gzipped {
					var buf bytes.Buffer
					gw := gzip.NewWriter(&buf)
					_, err = gw.Write(raw)
					require.NoError(t, err)
					require.NoError(t, gw.Close())
					payload = buf.Bytes()
				}

				// Use a custom RoundTripper so Transport does not auto-strip Content-Encoding: gzip.
				hc := &http.Client{
					Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
						hdr := make(http.Header)
						hdr.Set("Content-Type", "application/x-protobuf")
						if tc.gzipped {
							hdr.Set("Content-Encoding", "gzip")
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     hdr,
							Body:       io.NopCloser(bytes.NewReader(payload)),
						}, nil
					}),
				}
				exp, err := otlpprofilehttp.New(
					t.Context(),
					otlpprofilehttp.WithInsecure(),
					otlpprofilehttp.WithHTTPClient(hc),
				)
				require.NoError(t, err)
				t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

				err = exp.Export(t.Context(), sampleResourceProfiles())
				if tc.wantErr {
					require.ErrorIs(t, err, internal.PartialSuccess{})
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	t.Run("non_protobuf_2xx_ignored", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
	})

	t.Run("malformed_protobuf_2xx_fails_without_retry", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			w.Header().Set("Content-Type", "application/x-protobuf")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte{0xff, 0xff, 0xff})
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.Error(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("max_response_size_plain_and_gzip", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(bytes.Repeat([]byte("a"), 128))
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithMaxResponseSize(32),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "response body too large")
		assert.Equal(t, int32(1), attempts.Load())

		// Malformed gzip body also fails without retry.
		hc := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				hdr := make(http.Header)
				hdr.Set("Content-Encoding", "gzip")
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     hdr,
					Body:       io.NopCloser(bytes.NewReader([]byte("not-gzip"))),
				}, nil
			}),
		}
		exp2, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(hc),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp2.Shutdown(t.Context()) })
		require.Error(t, exp2.Export(t.Context(), sampleResourceProfiles()))
	})

	t.Run("max_request_size_zero_disables_and_enforced_sends_zero_requests", func(t *testing.T) {
		var requests atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			requests.Add(1)
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		limited, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithMaxRequestSize(8),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = limited.Shutdown(t.Context()) })
		require.Error(t, limited.Export(t.Context(), sampleResourceProfiles()))
		assert.Zero(t, requests.Load())

		unlimited, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithMaxRequestSize(0),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = unlimited.Shutdown(t.Context()) })
		require.NoError(t, unlimited.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(1), requests.Load())
	})

	t.Run("header_precedence", func(t *testing.T) {
		var mu sync.Mutex
		var gotUA, gotCT, gotVer string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			gotUA = r.Header.Get("User-Agent")
			gotCT = r.Header.Get("Content-Type")
			gotVer = r.Header.Get("OTLP-Profiles-Development-Version")
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithHeaders(map[string]string{
				"User-Agent":                        "custom-agent/1.0",
				"Content-Type":                      "application/json",
				"OTLP-Profiles-Development-Version": "999",
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		mu.Lock()
		assert.Equal(t, "custom-agent/1.0", gotUA)
		assert.Equal(t, "application/x-protobuf", gotCT)
		assert.Equal(t, "1", gotVer)
		mu.Unlock()
	})
}

func TestClientRetriesTimeoutsAndEdgeErrors(t *testing.T) {
	fastRetry := otlpprofilehttp.RetryConfig{
		Enabled:         true,
		InitialInterval: time.Millisecond,
		MaxInterval:     5 * time.Millisecond,
		MaxElapsedTime:  time.Second,
	}

	t.Run("retryable_status_codes_and_body_replay", func(t *testing.T) {
		for _, status := range []int{
			http.StatusTooManyRequests,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		} {
			for _, comp := range []otlpprofilehttp.Compression{
				otlpprofilehttp.NoCompression,
				otlpprofilehttp.GzipCompression,
			} {
				var attempts atomic.Int32
				var bodies [][]byte
				var mu sync.Mutex
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					b, _ := io.ReadAll(r.Body)
					mu.Lock()
					bodies = append(bodies, b)
					mu.Unlock()
					if attempts.Add(1) == 1 {
						w.Header().Set("Retry-After", "0")
						w.WriteHeader(status)
						return
					}
					w.WriteHeader(http.StatusOK)
				}))

				exp, err := otlpprofilehttp.New(
					t.Context(),
					otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
					otlpprofilehttp.WithCompression(comp),
					otlpprofilehttp.WithRetry(fastRetry),
				)
				require.NoError(t, err)

				require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
				assert.Equal(t, int32(2), attempts.Load())
				mu.Lock()
				require.Len(t, bodies, 2)
				assert.NotEmpty(t, bodies[0])
				assert.Equal(t, bodies[0], bodies[1])
				mu.Unlock()

				_ = exp.Shutdown(t.Context())
				srv.Close()
			}
		}
	})

	t.Run("redirect_307_replays_body", func(t *testing.T) {
		for _, comp := range []otlpprofilehttp.Compression{
			otlpprofilehttp.NoCompression,
			otlpprofilehttp.GzipCompression,
		} {
			var finalBody []byte
			var mu sync.Mutex
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1development/profiles" {
					w.Header().Set("Location", "/redirected")
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}
				b, _ := io.ReadAll(r.Body)
				mu.Lock()
				finalBody = b
				mu.Unlock()
				w.WriteHeader(http.StatusOK)
			}))

			exp, err := otlpprofilehttp.New(
				t.Context(),
				otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
				otlpprofilehttp.WithCompression(comp),
			)
			require.NoError(t, err)

			require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
			mu.Lock()
			assert.NotEmpty(t, finalBody)
			mu.Unlock()

			_ = exp.Shutdown(t.Context())
			srv.Close()
		}
	})

	t.Run("retry_disabled_fails_after_one_attempt", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithRetry(otlpprofilehttp.RetryConfig{Enabled: false}),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.Error(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("huge_retry_after_respects_max_elapsed_time", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			w.Header().Set("Retry-After", "9223372036854775807")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithRetry(fastRetry),
			otlpprofilehttp.WithTimeout(0),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "max retry time would elapse")
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("temporary_and_permanent_transport_errors", func(t *testing.T) {
		var attempts atomic.Int32
		hc := &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if attempts.Add(1) == 1 {
					return nil, &url.Error{Op: "Post", URL: r.URL.String(), Err: tempErr("transient")}
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader(nil)),
				}, nil
			}),
		}
		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(hc),
			otlpprofilehttp.WithRetry(fastRetry),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(2), attempts.Load())

		// Permanent transport error does not retry.
		var permAttempts atomic.Int32
		permClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				permAttempts.Add(1)
				return nil, errors.New("connection refused")
			}),
		}
		exp2, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(permClient),
			otlpprofilehttp.WithRetry(fastRetry),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp2.Shutdown(t.Context()) })

		require.Error(t, exp2.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(1), permAttempts.Load())
	})

	t.Run("body_close_error_reset_on_retry_and_returned_on_final", func(t *testing.T) {
		var attempts atomic.Int32
		closeErr := errors.New("close boom")
		hc := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				if attempts.Add(1) == 1 {
					return &http.Response{
						StatusCode: http.StatusServiceUnavailable,
						Header:     http.Header{"Retry-After": []string{"0"}},
						Body:       errCloseBody{Reader: bytes.NewReader(nil), closeErr: closeErr},
					}, nil
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader(nil)),
				}, nil
			}),
		}
		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(hc),
			otlpprofilehttp.WithRetry(fastRetry),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		// Attempt 1's close error must not leak when attempt 2 succeeds.
		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))

		// Final attempt's close error is returned.
		hcFinal := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       errCloseBody{Reader: bytes.NewReader(nil), closeErr: closeErr},
				}, nil
			}),
		}
		exp2, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(hcFinal),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp2.Shutdown(t.Context()) })
		require.ErrorIs(t, exp2.Export(t.Context(), sampleResourceProfiles()), closeErr)
	})

	t.Run("nil_response_and_nil_body_do_not_panic", func(t *testing.T) {
		nilRespClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return nil, nil
			}),
		}
		exp1, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(nilRespClient),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp1.Shutdown(t.Context()) })
		require.Error(t, exp1.Export(t.Context(), sampleResourceProfiles()))

		nilBodyClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}, nil
			}),
		}
		exp2, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithHTTPClient(nilBodyClient),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp2.Shutdown(t.Context()) })
		require.NoError(t, exp2.Export(t.Context(), sampleResourceProfiles()))
	})

	t.Run("timeout_and_context_cancellation", func(t *testing.T) {
		// Pre-canceled context does not hit the server.
		var hits atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			select {
			case <-time.After(2 * time.Second):
			case <-r.Context().Done():
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithTimeout(50*time.Millisecond),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		canceled, cancel := context.WithCancel(t.Context())
		cancel()
		require.ErrorIs(t, exp.Export(canceled, sampleResourceProfiles()), context.Canceled)
		assert.Zero(t, hits.Load())

		// Slow server hits exporter timeout.
		require.Error(t, exp.Export(t.Context(), sampleResourceProfiles()))
	})

	t.Run("invalid_endpoint_compression_and_zero_value_exporter", func(t *testing.T) {
		_, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpoint("host:invalid-port"),
			otlpprofilehttp.WithInsecure(),
		)
		require.Error(t, err)

		expBadComp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithCompression(otlpprofilehttp.Compression(99)),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = expBadComp.Shutdown(t.Context()) })
		err = expBadComp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported compression")

		var zeroExp otlpprofilehttp.Exporter
		require.Error(t, zeroExp.Export(t.Context(), sampleResourceProfiles()))
		require.NoError(t, zeroExp.ForceFlush(t.Context()))
		require.NoError(t, zeroExp.Shutdown(t.Context()))
		require.NoError(t, zeroExp.ForceFlush(t.Context()))
	})
}

func TestExporterConcurrentSafe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"))
	require.NoError(t, err)

	rp := sampleResourceProfiles()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 10 {
				err := exp.Export(t.Context(), rp)
				if err != nil {
					assert.ErrorIs(t, err, profile.ErrExporterShutdown)
				}
				assert.NoError(t, exp.ForceFlush(t.Context()))
			}
		})
	}

	time.Sleep(time.Millisecond)
	assert.NoError(t, exp.Shutdown(t.Context()))
	wg.Wait()
	assert.NoError(t, exp.ForceFlush(t.Context()))
	assert.ErrorIs(t, exp.Export(t.Context(), rp), profile.ErrExporterShutdown)
}
