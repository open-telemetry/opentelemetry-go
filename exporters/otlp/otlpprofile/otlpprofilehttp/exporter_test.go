// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp_test

import (
	"compress/gzip"
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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func sampleResourceProfiles() *profiledata.ResourceProfiles {
	loc := &profiledata.Location{
		Address: 0x401000,
		Lines: []profiledata.Line{
			{Function: profiledata.Function{Name: "main.main", Filename: "main.go"}, Line: 42},
		},
	}
	return &profiledata.ResourceProfiles{
		Resource: resource.NewSchemaless(attribute.String("service.name", "otlp-http-test")),
		ScopeProfiles: []profiledata.ScopeProfiles{
			{
				Scope: instrumentation.Scope{Name: "go.opentelemetry.io/otel/sdk/profile/runtime/heap"},
				Profiles: []profiledata.Profile{
					{
						SampleType: profiledata.ValueType{Type: "inuse_space", Unit: "bytes"},
						Samples:    []profiledata.Sample{{Stack: []*profiledata.Location{loc}, Values: []int64{4096}}},
						Time:       time.Unix(1_700_000_000, 0).UTC(),
					},
				},
			},
		},
	}
}

func TestExporterEndToEndWithRuntimeProducer(t *testing.T) {
	var (
		mu                   sync.Mutex
		receivedPath         string
		receivedContentType  string
		receivedDevVersion   string
		receivedEncoding     string
		receivedCustomHeader string
		receivedUserAgent    string
		receivedReq          *colprofilepb.ExportProfilesServiceRequest
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader = r.Body
		encoding := r.Header.Get("Content-Encoding")
		if encoding == "gzip" {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer func() { _ = gr.Close() }()
			reader = gr
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req := new(colprofilepb.ExportProfilesServiceRequest)
		if err := proto.Unmarshal(body, req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		receivedPath = r.URL.Path
		receivedContentType = r.Header.Get("Content-Type")
		receivedDevVersion = r.Header.Get("OTLP-Profiles-Development-Version")
		receivedEncoding = encoding
		receivedCustomHeader = r.Header.Get("X-Custom")
		receivedUserAgent = r.Header.Get("User-Agent")
		receivedReq = req
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	exp, err := otlpprofilehttp.New(
		t.Context(),
		otlpprofilehttp.WithEndpoint(u.Host),
		otlpprofilehttp.WithInsecure(),
		otlpprofilehttp.WithCompression(otlpprofilehttp.GzipCompression),
		otlpprofilehttp.WithHeaders(map[string]string{"X-Custom": "hello"}),
	)
	require.NoError(t, err)

	mr := profile.NewManualReader()
	pp := profile.NewProfilerProvider(
		profile.WithResource(resource.NewSchemaless(attribute.String("service.name", "e2e-svc"))),
		profile.WithProducer(profile.NewRuntimeProducer()),
		profile.WithReader(mr),
	)
	t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

	var rp profiledata.ResourceProfiles
	require.NoError(t, mr.Collect(t.Context(), &rp))
	require.NotEmpty(t, rp.ScopeProfiles)

	require.NoError(t, exp.Export(t.Context(), &rp))

	mu.Lock()
	assert.Equal(t, "/v1development/profiles", receivedPath)
	assert.Equal(t, "application/x-protobuf", receivedContentType)
	assert.Equal(t, "1", receivedDevVersion)
	assert.Equal(t, "gzip", receivedEncoding)
	assert.Equal(t, "hello", receivedCustomHeader)
	assert.Equal(t, "OTel Go OTLP over HTTP/protobuf profiles exporter/"+otlpprofilehttp.Version(), receivedUserAgent)
	require.NotNil(t, receivedReq)
	require.Len(t, receivedReq.ResourceProfiles, 1)
	assert.NotEmpty(t, receivedReq.ResourceProfiles[0].ScopeProfiles)
	require.NotNil(t, receivedReq.Dictionary)
	assert.NotEmpty(t, receivedReq.Dictionary.StringTable)
	assert.Empty(t, receivedReq.Dictionary.StringTable[0])
	assert.Contains(t, receivedReq.Dictionary.StringTable, "e2e-svc")
	mu.Unlock()

	require.NoError(t, exp.ForceFlush(t.Context()))
	require.NoError(t, exp.Shutdown(t.Context()))
	require.NoError(t, exp.Shutdown(t.Context())) // Idempotent.
	assert.ErrorIs(t, exp.Export(t.Context(), &rp), profile.ErrExporterShutdown)
}

func TestExporterPartialSuccessAndRetries(t *testing.T) {
	t.Run("partial_success", func(t *testing.T) {
		resp := &colprofilepb.ExportProfilesServiceResponse{
			PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{
				RejectedProfiles: 2,
				ErrorMessage:     "dropped 2 profiles",
			},
		}
		respBytes, err := proto.Marshal(resp)
		require.NoError(t, err)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/x-protobuf")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(respBytes)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.ErrorIs(t, err, internal.PartialSuccess{})
		assert.Contains(t, err.Error(), "dropped 2 profiles")
	})

	t.Run("retries_transient_status", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if attempts.Add(1) == 1 {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"),
			otlpprofilehttp.WithRetry(otlpprofilehttp.RetryConfig{
				Enabled:         true,
				InitialInterval: time.Millisecond,
				MaxInterval:     5 * time.Millisecond,
				MaxElapsedTime:  time.Second,
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(2), attempts.Load())
	})

	t.Run("non_retryable_400_fails_immediately", func(t *testing.T) {
		var attempts atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			http.Error(w, "unsupported version", http.StatusBadRequest)
		}))
		defer srv.Close()

		exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithEndpointURL(srv.URL+"/v1development/profiles"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "400 Bad Request")
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("max_request_size_enforced", func(t *testing.T) {
		exp, err := otlpprofilehttp.New(
			t.Context(),
			otlpprofilehttp.WithEndpoint("localhost:4318"),
			otlpprofilehttp.WithInsecure(),
			otlpprofilehttp.WithMaxRequestSize(8),
		)
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "request body too large")
	})

	t.Run("empty_resource_profiles_noop", func(t *testing.T) {
		exp, err := otlpprofilehttp.New(t.Context(), otlpprofilehttp.WithInsecure())
		require.NoError(t, err)
		t.Cleanup(func() { _ = exp.Shutdown(t.Context()) })

		require.NoError(t, exp.Export(t.Context(), nil))
		require.NoError(t, exp.Export(t.Context(), &profiledata.ResourceProfiles{}))
	})
}
