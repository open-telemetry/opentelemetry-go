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

package otlptracehttp_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlptracetest"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

const (
	relOtherTracesPath = "post/traces/here"
	otherTracesPath    = "/post/traces/here"
)

var (
	testHeaders = map[string]string{
		"Otel-Go-Key-1": "somevalue",
		"Otel-Go-Key-2": "someothervalue",
	}

	customUserAgentHeader = map[string]string{
		"user-agent": "custome-user-agent",
	}
)

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name  string
		opts  []otlptracehttp.Option
		mcCfg mockCollectorConfig
		tls   bool
	}{
		{
			name: "no extra options",
			opts: nil,
		},
		{
			name: "with gzip compression",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			},
		},
		{
			name: "retry",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{503, 429},
			},
		},
		{
			name: "retry with gzip compression",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{503, 502},
			},
		},
		{
			name: "retry with throttle",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: time.Nanosecond,
					MaxInterval:     time.Nanosecond,
					// Do not stop trying.
					MaxElapsedTime: 0,
				}),
			},
			mcCfg: mockCollectorConfig{
				InjectHTTPStatus: []int{504},
				InjectResponseHeader: []map[string]string{
					{"Retry-After": "10"},
				},
			},
		},
		{
			name: "with empty paths (forced to defaults)",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithURLPath(""),
			},
		},
		{
			name: "with relative paths",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithURLPath(relOtherTracesPath),
			},
			mcCfg: mockCollectorConfig{
				TracesURLPath: otherTracesPath,
			},
		},
		{
			name: "with TLS",
			opts: nil,
			mcCfg: mockCollectorConfig{
				WithTLS: true,
			},
			tls: true,
		},
		{
			name: "with extra headers",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithHeaders(testHeaders),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: testHeaders,
			},
		},
		{
			name: "with custom user agent",
			opts: []otlptracehttp.Option{
				otlptracehttp.WithHeaders(customUserAgentHeader),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: customUserAgentHeader,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mc := runMockCollector(t, tc.mcCfg)
			defer mc.MustStop(t)
			allOpts := []otlptracehttp.Option{
				otlptracehttp.WithEndpoint(mc.Endpoint()),
			}
			if tc.tls {
				tlsConfig := mc.ClientTLSConfig()
				require.NotNil(t, tlsConfig)
				allOpts = append(allOpts, otlptracehttp.WithTLSClientConfig(tlsConfig))
			} else {
				allOpts = append(allOpts, otlptracehttp.WithInsecure())
			}
			allOpts = append(allOpts, tc.opts...)
			client := otlptracehttp.NewClient(allOpts...)
			ctx := context.Background()
			exporter, err := otlptrace.New(ctx, client)
			if assert.NoError(t, err) {
				defer func() {
					assert.NoError(t, exporter.Shutdown(ctx))
				}()
				otlptracetest.RunEndToEndTest(ctx, t, exporter, mc)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollector(t, mockCollectorConfig{})
	defer func() {
		_ = mc.Stop()
	}()

	<-time.After(5 * time.Millisecond)

	otlptracetest.RunExporterShutdownTest(t, func() otlptrace.Client {
		return otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(mc.endpoint),
		)
	})
}

func TestTimeout(t *testing.T) {
	delay := make(chan struct{})
	mcCfg := mockCollectorConfig{Delay: delay}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	defer func() { close(delay) }()
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(time.Nanosecond),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{Enabled: false}),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, client)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.ErrorContains(t, err, "retry-able request failure")
}

func TestNoRetry(t *testing.T) {
	mc := runMockCollector(t, mockCollectorConfig{
		InjectHTTPStatus: []int{http.StatusBadRequest},
	})
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Nanosecond,
			MaxInterval:     1 * time.Nanosecond,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	unwrapped := errors.Unwrap(err)
	assert.Equal(t, fmt.Sprintf("failed to send to http://%s/v1/traces: 400 Bad Request", mc.endpoint), unwrapped.Error())
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "))
	assert.Empty(t, mc.GetSpans())
}

func TestEmptyData(t *testing.T) {
	mcCfg := mockCollectorConfig{}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	assert.NoError(t, err)
	err = exporter.ExportSpans(ctx, nil)
	assert.NoError(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestCancelledContext(t *testing.T) {
	mcCfg := mockCollectorConfig{}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx, cancel := context.WithCancel(context.Background())
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()
	cancel()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestDeadlineContext(t *testing.T) {
	statuses := make([]int, 0, 5)
	for i := 0; i < cap(statuses); i++ {
		statuses = append(statuses, http.StatusTooManyRequests)
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Hour,
			MaxInterval:     1 * time.Hour,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.Error(t, err)
	assert.Empty(t, mc.GetSpans())
}

func TestStopWhileExportingConcurrentSafe(t *testing.T) {
	statuses := make([]int, 0, 5)
	for i := 0; i < cap(statuses); i++ {
		statuses = append(statuses, http.StatusTooManyRequests)
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Hour,
			MaxInterval:     1 * time.Hour,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	doneCh := make(chan struct{})
	go func() {
		err := exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
		assert.Error(t, err)
		assert.Empty(t, mc.GetSpans())
		close(doneCh)
	}()
	<-time.After(time.Second)
	err = exporter.Shutdown(ctx)
	assert.NoError(t, err)
	<-doneCh
}

func TestPartialSuccess(t *testing.T) {
	mcCfg := mockCollectorConfig{
		Partial: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 2,
			ErrorMessage:  "partially successful",
		},
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()

	errs := []error{}
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errs = append(errs, err)
	}))
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.NoError(t, err)

	require.Equal(t, 1, len(errs))
	require.Contains(t, errs[0].Error(), "partially successful")
	require.Contains(t, errs[0].Error(), "2 spans rejected")
}

func TestOtherHTTPSuccess(t *testing.T) {
	for code := 201; code <= 299; code++ {
		t.Run(fmt.Sprintf("status_%d", code), func(t *testing.T) {
			mcCfg := mockCollectorConfig{
				InjectHTTPStatus: []int{code},
			}
			mc := runMockCollector(t, mcCfg)
			defer mc.MustStop(t)
			driver := otlptracehttp.NewClient(
				otlptracehttp.WithEndpoint(mc.Endpoint()),
				otlptracehttp.WithInsecure(),
			)
			ctx := context.Background()
			exporter, err := otlptrace.New(ctx, driver)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, exporter.Shutdown(context.Background()))
			}()

			errs := []error{}
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
				errs = append(errs, err)
			}))
			err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
			assert.NoError(t, err)

			assert.Equal(t, 0, len(errs))
		})
	}
}

func TestCollectorRespondingNonProtobufContent(t *testing.T) {
	mcCfg := mockCollectorConfig{
		InjectContentType: "application/octet-stream",
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(mc.Endpoint()),
		otlptracehttp.WithInsecure(),
	)
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()
	err = exporter.ExportSpans(ctx, otlptracetest.SingleReadOnlySpan())
	assert.NoError(t, err)
	assert.Len(t, mc.GetSpans(), 1)
}
