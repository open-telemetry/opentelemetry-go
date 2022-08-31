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

package otlpmetrichttp_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpmetrictest"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	relOtherMetricsPath = "post/metrics/here"
	otherMetricsPath    = "/post/metrics/here"
)

var (
	oneRecord = otlpmetrictest.OneRecordReader()

	testResource = resource.Empty()
)

var (
	testHeaders = map[string]string{
		"Otel-Go-Key-1": "somevalue",
		"Otel-Go-Key-2": "someothervalue",
	}
)

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name  string
		opts  []otlpmetrichttp.Option
		mcCfg mockCollectorConfig
		tls   bool
	}{
		{
			name: "no extra options",
			opts: nil,
		},
		{
			name: "with gzip compression",
			opts: []otlpmetrichttp.Option{
				otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
			},
		},
		{
			name: "with empty paths (forced to defaults)",
			opts: []otlpmetrichttp.Option{
				otlpmetrichttp.WithURLPath(""),
			},
		},
		{
			name: "with relative paths",
			opts: []otlpmetrichttp.Option{
				otlpmetrichttp.WithURLPath(relOtherMetricsPath),
			},
			mcCfg: mockCollectorConfig{
				MetricsURLPath: otherMetricsPath,
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
			opts: []otlpmetrichttp.Option{
				otlpmetrichttp.WithHeaders(testHeaders),
			},
			mcCfg: mockCollectorConfig{
				ExpectedHeaders: testHeaders,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mc := runMockCollector(t, tc.mcCfg)
			defer mc.MustStop(t)
			allOpts := []otlpmetrichttp.Option{
				otlpmetrichttp.WithEndpoint(mc.Endpoint()),
			}
			if tc.tls {
				tlsConfig := mc.ClientTLSConfig()
				require.NotNil(t, tlsConfig)
				allOpts = append(allOpts, otlpmetrichttp.WithTLSClientConfig(tlsConfig))
			} else {
				allOpts = append(allOpts, otlpmetrichttp.WithInsecure())
			}
			allOpts = append(allOpts, tc.opts...)
			client := otlpmetrichttp.NewClient(allOpts...)
			ctx := context.Background()
			exporter, err := otlpmetric.New(ctx, client)
			if assert.NoError(t, err) {
				defer func() {
					assert.NoError(t, exporter.Shutdown(ctx))
				}()
				otlpmetrictest.RunEndToEndTest(ctx, t, exporter, mc)
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

	otlpmetrictest.RunExporterShutdownTest(t, func() otlpmetric.Client {
		return otlpmetrichttp.NewClient(
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpoint(mc.endpoint),
		)
	})
}

func TestTimeout(t *testing.T) {
	delay := make(chan struct{})
	mcCfg := mockCollectorConfig{Delay: delay}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	defer func() { close(delay) }()
	client := otlpmetrichttp.NewClient(
		otlpmetrichttp.WithEndpoint(mc.Endpoint()),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithTimeout(time.Nanosecond),
	)
	ctx := context.Background()
	exporter, err := otlpmetric.New(ctx, client)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	err = exporter.Export(ctx, testResource, oneRecord)
	assert.Equalf(t, true, os.IsTimeout(err), "expected timeout error, got: %v", err)
}

func TestEmptyData(t *testing.T) {
	mcCfg := mockCollectorConfig{}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlpmetrichttp.NewClient(
		otlpmetrichttp.WithEndpoint(mc.Endpoint()),
		otlpmetrichttp.WithInsecure(),
	)
	ctx := context.Background()
	exporter, err := otlpmetric.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	assert.NoError(t, err)
	err = exporter.Export(ctx, testResource, oneRecord)
	assert.NoError(t, err)
	assert.NotEmpty(t, mc.GetMetrics())
}

func TestCancelledContext(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlpmetrichttp.NewClient(
		otlpmetrichttp.WithEndpoint(mc.Endpoint()),
		otlpmetrichttp.WithInsecure(),
	)
	ctx, cancel := context.WithCancel(context.Background())
	exporter, err := otlpmetric.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()
	cancel()
	_ = exporter.Export(ctx, testResource, oneRecord)
	assert.Empty(t, mc.GetMetrics())
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
	driver := otlpmetrichttp.NewClient(
		otlpmetrichttp.WithEndpoint(mc.Endpoint()),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithBackoff(time.Minute),
	)
	ctx := context.Background()
	exporter, err := otlpmetric.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(context.Background()))
	}()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err = exporter.Export(ctx, testResource, oneRecord)
	assert.Error(t, err)
	assert.Empty(t, mc.GetMetrics())
}

func TestStopWhileExporting(t *testing.T) {
	statuses := make([]int, 0, 5)
	for i := 0; i < cap(statuses); i++ {
		statuses = append(statuses, http.StatusTooManyRequests)
	}
	mcCfg := mockCollectorConfig{
		InjectHTTPStatus: statuses,
	}
	mc := runMockCollector(t, mcCfg)
	defer mc.MustStop(t)
	driver := otlpmetrichttp.NewClient(
		otlpmetrichttp.WithEndpoint(mc.Endpoint()),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithBackoff(time.Minute),
	)
	ctx := context.Background()
	exporter, err := otlpmetric.New(ctx, driver)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, exporter.Shutdown(ctx))
	}()
	doneCh := make(chan struct{})
	go func() {
		err := exporter.Export(ctx, testResource, oneRecord)
		assert.Error(t, err)
		assert.Empty(t, mc.GetMetrics())
		close(doneCh)
	}()
	<-time.After(time.Second)
	err = exporter.Shutdown(ctx)
	assert.NoError(t, err)
	<-doneCh
}
